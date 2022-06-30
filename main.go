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
	imageName    = "ghcr.io/chelnak/cat-github-metric-collector"
	scheduleCron = "0 0 * * *"
	scheduleType = "schedule"
)

func main() {
	modules, err := getSupportedModules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := workflow.NewWorkflow()
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
					"connection":          "${connections.gcp.'content-and-tooling-lab'}",
					"module_owner":        moduleOwner,
					"module_name":         module.Name,
					"github_token":        "${secrets.GITHUB_TOKEN}",
					"bigquery_project_id": "${secrets.BIGQUERY_PROJECT}",
				},
			},
		)
	}

	w.Print(nil)
}

func getSupportedModules() (*[]collar.Module, error) {
	client := collar.NewModuleClient(nil, "")
	ctx := context.Background()
	return client.GetSupportedModules(ctx)
}
