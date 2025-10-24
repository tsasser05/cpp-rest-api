package godog

import (
	"testing"

	"github.com/cucumber/godog"
	"cpp-rest-api-tests/step_definitions"
)

func TestContactFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			step_definitions.InitializeScenario(ctx)
		},
		Options: &godog.Options{
			Format:   "progress,cucumber:report.json",
			Paths:    []string{"../features/contacts.feature"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("failed to run contact feature tests")
	}
}

