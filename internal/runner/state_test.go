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

		input := domain.Resources{"hero": statement}

		expected := domain.Resources{
			"hero": statement,
		}

		state := runner.NewState(input)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})

	t.Run("should return resource with no parameter or static parameters using alias as resource id", func(t *testing.T) {
		statement := domain.Statement{Method: "from", Resource: "hero", Alias: "h", With: map[string]interface{}{"id": "123456"}}

		input := domain.Resources{"h": statement}

		expected := domain.Resources{
			"h": statement,
		}

		state := runner.NewState(input)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})

	t.Run("should not return resource with unresolved dependency", func(t *testing.T) {
		heroStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}
		sidekickStatement := domain.Statement{Method: "from", Resource: "sidekick", With: map[string]interface{}{"id": domain.Chain{"hero", "sidekick", "id"}}}
		villainStatement := domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{domain.Chain{"hero", "villain", "id"}}}}
		crossoverStatement := domain.Statement{Method: "from", Resource: "crossover", With: map[string]interface{}{"id": map[string]interface{}{"heroes": domain.Chain{"hero", "id"}}}}

		input := domain.Resources{
			"hero":      heroStatement,
			"sidekick":  sidekickStatement,
			"villain":   villainStatement,
			"crossover": crossoverStatement,
		}

		expected := domain.Resources{
			"hero": heroStatement,
		}

		state := runner.NewState(input)
		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})

	t.Run("should return resource with resolved dependency inside complex param", func(t *testing.T) {
		heroStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}
		sidekickStatement := domain.Statement{Method: "from", Resource: "sidekick", With: map[string]interface{}{"id": domain.Chain{"hero", "sidekick", "id"}}}
		villainStatement := domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{domain.Chain{"hero", "villain", "id"}}}}
		crossoverStatement := domain.Statement{Method: "from", Resource: "crossover", With: map[string]interface{}{"id": map[string]interface{}{"heroes": domain.Chain{"hero", "id"}}}}

		input := domain.Resources{
			"hero":      heroStatement,
			"sidekick":  sidekickStatement,
			"villain":   villainStatement,
			"crossover": crossoverStatement,
		}

		expected := domain.Resources{
			"sidekick":  sidekickStatement,
			"villain":   villainStatement,
			"crossover": crossoverStatement,
		}

		state := runner.NewState(input)
		state.SetAsRequest("hero")
		state.UpdateDone("hero", nil)

		got := state.Available()

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Available = %#+v, want = %#+v", got, expected)
		}
	})
}

func TestSetAsRequested(t *testing.T) {
	t.Run("should add resource to requested and remove from available", func(t *testing.T) {
		heroStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "987654"}}
		sidekickStatement := domain.Statement{Method: "from", Resource: "sidekick", With: map[string]interface{}{"id": "123456"}}

		input := domain.Resources{"hero": heroStatement, "sidekick": sidekickStatement}

		state := runner.NewState(input)

		expectedInitialAvailable := domain.Resources{"hero": heroStatement, "sidekick": sidekickStatement}
		expectedInitialRequested := domain.Resources{}

		initialAvailable := state.Available()
		initialRequested := state.Requested()

		if !reflect.DeepEqual(initialAvailable, expectedInitialAvailable) {
			t.Fatalf("Initial Available = %#+v, want = %#+v", initialAvailable, expectedInitialAvailable)
		}

		if !reflect.DeepEqual(initialRequested, expectedInitialRequested) {
			t.Fatalf(" Initial Requested = %#+v, want = %#+v", initialRequested, expectedInitialRequested)
		}

		state.SetAsRequest("hero")

		expectedFinalAvailable := domain.Resources{
			"sidekick": sidekickStatement,
		}

		expectedFinalRequested := domain.Resources{"hero": heroStatement}

		finalAvailable := state.Available()
		finalRequested := state.Requested()

		if !reflect.DeepEqual(finalAvailable, expectedFinalAvailable) {
			t.Fatalf("Final Available = %#+v, want = %#+v", finalAvailable, expectedFinalAvailable)
		}

		if !reflect.DeepEqual(finalRequested, expectedFinalRequested) {
			t.Fatalf("Final Requested = %#+v, want = %#+v", finalRequested, expectedFinalRequested)
		}
	})
}

func TestUpdateDone(t *testing.T) {
	t.Run("should add completed resource to done and remove from requested", func(t *testing.T) {
		doneStatement := domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}}
		input := domain.Resources{"hero": doneStatement}

		expectedDoneRequests := domain.Resources{
			"hero": domain.DoneResource{Details: domain.Details{Status: 200}, Result: []byte{}},
		}
		expectedRequestedStatements := domain.Resources{}

		state := runner.NewState(input)

		response := domain.DoneResource{Details: domain.Details{Status: 200}, Result: []byte{}}

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
