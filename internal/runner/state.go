package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

type AvailableResources Resources
type RequestedResources Resources
type DoneResources Resources

type DoneRequest domain.Response
type DoneRequests []interface{}

type State struct {
	todo      Resources
	requested RequestedResources
	done      DoneResources
}

func NewState(todo Resources) *State {
	return &State{todo: todo, requested: make(RequestedResources), done: make(DoneResources)}
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
		switch v := v.(type) {
		case domain.Chain:
			resourceTarget, ok := v[0].(string)
			if !ok {
				return false
			}

			_, found := s.done[ResourceId(resourceTarget)]
			return found
		default:
			continue
		}
	}

	return true
}

func (s *State) UpdateDone(resourceId ResourceId, response interface{}) {
	s.done[resourceId] = response
	delete(s.requested, resourceId)
}

func (s *State) Done() DoneResources {
	return s.done
}

func (s *State) SetAsRequest(resourceId ResourceId) {
	statement := s.todo[resourceId]
	s.requested[resourceId] = statement
	delete(s.todo, resourceId)
}

func (s *State) Requested() RequestedResources {
	return s.requested
}

func (s *State) HasFinished() bool {
	return len(s.todo) == 0 && len(s.requested) == 0
}
