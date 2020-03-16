package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

type State struct {
	todo      Resources
	requested Resources
	done      Resources
}

func NewState(todo Resources) *State {
	return &State{todo: todo, requested: make(Resources), done: make(Resources)}
}

func (s *State) Available() Resources {
	available := make(Resources)
	for key, stmt := range s.todo {
		switch stmt := stmt.(type) {
		case domain.Statement:
			if s.canRequest(stmt) {
				available[key] = stmt
			}
		case []interface{}:
			if s.canRequestMultiplexedGroup(stmt) {
				available[key] = stmt
			}
		}
	}

	return available
}

func (s *State) canRequestMultiplexedGroup(statements []interface{}) bool {
	can := false
	for _, stmt := range statements {
		switch stmt := stmt.(type) {
		case domain.Statement:
			can = s.canRequest(stmt)
		case []interface{}:
			can = s.canRequestMultiplexedGroup(stmt)
		}
	}

	return can
}

func (s *State) canRequest(statement domain.Statement) bool {
	for _, v := range statement.With {
		if !s.isValueResolved(v) {
			return false
		}
	}

	return true
}

func (s *State) isValueResolved(value interface{}) bool {
	switch value := value.(type) {
	case domain.Chain:
		resourceTarget, ok := value[0].(string)
		if !ok {
			return false
		}

		_, found := s.done[ResourceId(resourceTarget)]
		return found
	case map[string]interface{}:
		for _, v := range value {
			if !s.isValueResolved(v) {
				return false
			}
		}
		return true
	case []interface{}:
		for _, v := range value {
			if !s.isValueResolved(v) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func (s *State) UpdateDone(resourceId ResourceId, response interface{}) {
	s.done[resourceId] = response
	delete(s.requested, resourceId)
}

func (s *State) Done() Resources {
	return s.done
}

func (s *State) SetAsRequest(resourceId ResourceId) {
	statement := s.todo[resourceId]
	s.requested[resourceId] = statement
	delete(s.todo, resourceId)
}

func (s *State) Requested() Resources {
	return s.requested
}

func (s *State) HasFinished() bool {
	return len(s.todo) == 0 && len(s.requested) == 0
}
