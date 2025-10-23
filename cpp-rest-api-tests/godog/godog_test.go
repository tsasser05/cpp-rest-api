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
			Tags:     "~@ignore",
			Paths:    []string{"../features"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run contact feature tests")
	}
}
