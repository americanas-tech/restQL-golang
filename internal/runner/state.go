package runner

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
)

// State tracks the status of the statements to be resolved
// being resolved and done.
type State struct {
	todo      domain.Resources
	requested domain.Resources
	done      domain.Resources
}

// NewState constructs a State value.
func NewState(todo domain.Resources) *State {
	return &State{todo: todo, requested: make(domain.Resources), done: make(domain.Resources)}
}

// Available return all statements to be resolved
// that have no dependency or which dependency was resolved.
func (s *State) Available() domain.Resources {
	available := make(domain.Resources)
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
	for _, v := range statement.With.Values {
		if !s.isValueResolved(v) {
			return false
		}
	}

	for _, v := range statement.Headers {
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

		_, found := s.done[domain.ResourceID(resourceTarget)]
		return found
	case domain.Function:
		return s.isValueResolved(value.Target())
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

// UpdateDone set a being resolved Resource as done.
func (s *State) UpdateDone(resourceID domain.ResourceID, response interface{}) {
	s.done[resourceID] = response
	delete(s.requested, resourceID)
}

// Done returns all Resources already resolved
func (s *State) Done() domain.Resources {
	d := s.done

	return d
}

// SetAsRequest define an to be resolved Resource
// into a being resolved Resource.
func (s *State) SetAsRequest(resourceID domain.ResourceID) {
	statement := s.todo[resourceID]
	s.requested[resourceID] = statement
	delete(s.todo, resourceID)
}

// Requested returns all Resources being resolved
func (s *State) Requested() domain.Resources {
	return s.requested
}

// HasFinished returns true if all Resources are set as done.
func (s *State) HasFinished() bool {
	return len(s.todo) == 0 && len(s.requested) == 0
}
