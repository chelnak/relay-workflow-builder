// package workflow provides an interface for creating a Relay Workflow.
package workflow

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Workflow is the interface that wraps the methods for adding steps and triggers.
type Workflow interface {
	AddTrigger(trigger Trigger)
	AddStep(step Step)
	Print(writer io.Writer)
}

// Step represents a single step in a Relay Workflow.
type Step struct {
	Name  string            `yaml:"name"`
	Image string            `yaml:"image"`
	Spec  map[string]string `yaml:"spec"`
}

// Trigger represents a single trigger in a Relay Workflow.
type Trigger struct {
	Name   string            `yaml:"name"`
	Source map[string]string `yaml:"source"`
}

type workflow struct {
	Triggers []Trigger `yaml:"triggers"`
	Steps    []Step    `yaml:"steps"`
}

// AddTrigger adds a trigger to the Workflow.
func (w *workflow) AddTrigger(trigger Trigger) {
	w.Triggers = append(w.Triggers, trigger)
}

// AddStep adds a step to the Workflow.
func (w *workflow) AddStep(step Step) {
	w.Steps = append(w.Steps, step)
}

// Print writes the YAML representation of the Workflow to the given writer.
// If no writer is provided, the output is written to stdout.
func (w *workflow) Print(writer io.Writer) {
	if writer == nil {
		writer = os.Stdout
	}

	b, err := yaml.Marshal(w)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(writer, string(b))
}

// NewWorkflow creates an empty Relay Workflow.
// It returns a Workflow interface that can be used to add steps and triggers.
func NewWorkflow() Workflow {
	return &workflow{}
}
