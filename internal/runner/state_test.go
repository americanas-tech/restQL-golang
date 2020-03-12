package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"reflect"
	"testing"
)

func TestAvailableResources(t *testing.T) {
	t.Run("should return resource with no parameter or static parameters", func(t *testing.T) {
		statement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}

		query := domain.Query{Statements: []domain.Statement{statement}}

		expected := runner.AvailableResources{
			"hero": statement,
		}

		state := runner.NewState(query)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})

	t.Run("should return resource with no parameter or static parameters using alias as resource id", func(t *testing.T) {
		statement := domain.Statement{Method: "from", Resource: "hero", Alias: "h", With: map[string]interface{}{"id": "123456"}}

		query := domain.Query{Statements: []domain.Statement{statement}}

		expected := runner.AvailableResources{
			"h": statement,
		}

		state := runner.NewState(query)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})

	t.Run("should not return resource with unresolved dependency", func(t *testing.T) {
		staticStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}
		unresolvedStatement := domain.Statement{Method: "from", Resource: "sidekick", With: map[string]interface{}{"id": domain.Chain{"hero", "sidekick", "id"}}}

		query := domain.Query{Statements: []domain.Statement{staticStatement, unresolvedStatement}}

		expected := runner.AvailableResources{
			"hero": staticStatement,
		}

		state := runner.NewState(query)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})

	t.Run("should return multiplexed resource grouped", func(t *testing.T) {
		firstStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}
		secondStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "987654"}}

		query := domain.Query{Statements: []domain.Statement{firstStatement, secondStatement}}

		expected := runner.AvailableResources{
			"hero": []domain.Statement{firstStatement, secondStatement},
		}

		state := runner.NewState(query)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})
}

func TestSetAsRequested(t *testing.T) {
	t.Run("should add resource to requested and remove from available", func(t *testing.T) {
		sidekickStatement := domain.Statement{Method: "from", Resource: "sidekick", With: map[string]interface{}{"id": "123456"}}
		heroStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "987654"}}

		query := domain.Query{Statements: []domain.Statement{sidekickStatement, heroStatement}}

		state := runner.NewState(query)

		expectedInitialAvailable := runner.AvailableResources{"sidekick": sidekickStatement, "hero": heroStatement}
		expectedInitialRequested := runner.RequestedResources{}

		initialAvailable := state.Available()
		initialRequested := state.Requested()

		if !reflect.DeepEqual(initialAvailable, expectedInitialAvailable) {
			t.Fatalf("Available = %#+v, want = %#+v", initialAvailable, expectedInitialAvailable)
		}

		if !reflect.DeepEqual(initialRequested, expectedInitialRequested) {
			t.Fatalf("Available = %#+v, want = %#+v", initialRequested, expectedInitialRequested)
		}

		state.SetAsRequest("hero")

		expectedFinalAvailable := runner.AvailableResources{
			"sidekick": sidekickStatement,
		}

		expectedFinalRequested := runner.RequestedResources{"hero": heroStatement}

		finalAvailable := state.Available()
		finalRequested := state.Requested()

		if !reflect.DeepEqual(finalAvailable, expectedFinalAvailable) {
			t.Fatalf("Available = %#+v, want = %#+v", finalAvailable, expectedFinalAvailable)
		}

		if !reflect.DeepEqual(finalRequested, expectedFinalRequested) {
			t.Fatalf("Available = %#+v, want = %#+v", finalRequested, expectedFinalRequested)
		}
	})
}

func TestUpdateDone(t *testing.T) {
	t.Run("should add completed resource to done and remove from requested", func(t *testing.T) {
		doneStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}
		query := domain.Query{Statements: []domain.Statement{doneStatement}}

		expectedDoneRequests := runner.DoneResources{
			"hero": runner.DoneRequest{StatusCode: 200, Body: []byte{}},
		}
		expectedRequestedStatements := runner.RequestedResources{}

		state := runner.NewState(query)

		response := runner.DoneRequest{StatusCode: 200, Body: []byte{}}

		state.UpdateDone("hero", response)

		gotRequestedStatements := state.Requested()
		gotDoneRequests := state.Done()

		if !reflect.DeepEqual(gotRequestedStatements, expectedRequestedStatements) {
			t.Fatalf("state had the requested statements = %#+v, expected = %#+v", gotRequestedStatements, expectedRequestedStatements)
		}

		if !reflect.DeepEqual(gotDoneRequests, expectedDoneRequests) {
			t.Fatalf("state had the done statementes = %#+v, want = %#+v", gotDoneRequests, expectedDoneRequests)
		}
	})
}
