package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

type ResourceId string
type AvailableResources map[ResourceId]interface{}
type RequestedResources map[ResourceId]interface{}
type DoneResources map[ResourceId]interface{}

type DoneRequest domain.Response
type DoneRequests []DoneRequest

type State struct {
	todo      map[ResourceId]interface{}
	requested RequestedResources
	done      DoneResources
}

func NewState(query domain.Query) *State {
	todo := makeTodoResources(query)

	return &State{todo: todo, requested: make(RequestedResources), done: make(DoneResources)}
}

func makeTodoResources(query domain.Query) map[ResourceId]interface{} {
	todo := make(map[ResourceId]interface{})
	for _, stmt := range query.Statements {
		index := getResourceId(stmt)
		current := todo[index]

		switch current := current.(type) {
		case []domain.Statement:
			todo[index] = append(current, stmt)
		case domain.Statement:
			list := []domain.Statement{current, stmt}
			todo[index] = list
		default:
			todo[index] = stmt
		}

	}

	return todo
}

func (s *State) Available() AvailableResources {
	available := make(AvailableResources)
	for key, stmt := range s.todo {
		switch stmt := stmt.(type) {
		case domain.Statement:
			if s.canRequest(stmt) {
				available[key] = stmt
			}
		case []domain.Statement:
			if s.canRequest(stmt[0]) {
				available[key] = stmt
			}
		}
	}

	return available
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

func getResourceId(statement domain.Statement) ResourceId {
	if statement.Alias != "" {
		return ResourceId(statement.Alias)
	}

	return ResourceId(statement.Resource)
}
