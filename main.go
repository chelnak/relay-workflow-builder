package main

import (
	"context"
	"fmt"
	"os"

	collar "github.com/chelnak/collar/pkg/modules"
	"github.com/chelnak/relay-workflow-builder/pkg/workflow"
)

const (
	moduleOwner  = "puppetlabs"
	imageName    = "ghcr.io/chelnak/cat-team-github-metrics:latest"
	scheduleCron = "0 0 * * *"
	scheduleType = "schedule"
)

func main() {
	modules, err := getSupportedModules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := workflow.NewWorkflow("A workflow for collecting GitHub metrics.")

	w.AddTrigger(
		workflow.Trigger{
			Name: "schedule",
			Source: map[string]string{
				"type":     scheduleType,
				"schedule": scheduleCron,
			},
		},
	)

	for _, module := range *modules {
		w.AddStep(
			workflow.Step{
				Name:  fmt.Sprintf("Metric collection: %s", module.Name),
				Image: imageName,
				Spec: map[string]string{
					"connection":           "${connections.gcp.'content-and-tooling-lab'}",
					"repo_owner":           moduleOwner,
					"repo_name":            module.Name,
					"github_token":         "${secrets.GITHUB_TOKEN}",
					"big_query_project_id": "${secrets.BIG_QUERY_PROJECT_ID}",
				},
			},
		)
	}

	err = w.Write(nil)
	if err != nil {
		if err != workflow.ErrValidation {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

func getSupportedModules() (*[]collar.Module, error) {
	client := collar.NewModuleClient(nil, "")
	ctx := context.Background()
	return client.GetSupportedModules(ctx)
}
