package runner

import "github.com/b2wdigital/restQL-golang/internal/domain"

type RunnerState struct {
	todo []domain.Statement
}

func NewRunnerState(query domain.Query) *RunnerState {
	return &RunnerState{todo: query.Statements}
}

func (rs *RunnerState) GetAvailableRequests() []domain.Statement {
	var available []domain.Statement
	for _, stmt := range rs.todo {
		if rs.canRequest(stmt) {
			available = append(available, stmt)
		}
	}

	return available
}

func (rs *RunnerState) canRequest(statement domain.Statement) bool {
	for _, v := range statement.With {
		switch v.(type) {
		case domain.Chain:
			return false
		default:
			continue
		}
	}

	return true
}
