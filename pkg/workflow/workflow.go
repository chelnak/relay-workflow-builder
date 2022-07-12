// Package workflow provides an interface for creating a Relay Workflow.
package workflow

import (
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

var ErrValidation = fmt.Errorf("validation error")

// Workflow is the interface that wraps the methods for adding steps and triggers.
type Workflow interface {
	AddParameter(key string, defaultValue string, description string)
	AddTrigger(trigger Trigger)
	AddStep(step Step)
	GetSteps() []Step
	Write(writer io.Writer) error
}

// Step represents a single step in a Relay Workflow.
type Step struct {
	Name      string            `yaml:"name" validate:"required"`
	DependsOn []string          `yaml:"dependsOn,omitempty"`
	Image     string            `yaml:"image" validate:"required"`
	Spec      map[string]string `yaml:"spec,omitempty"`
	When      string            `yaml:"when,omitempty"`
}

// TriggerBinding represents the binding property of a triggers.
type TriggerBinding struct {
	Key        string            `yaml:"key,omitempty"`
	Parameters map[string]string `yaml:"parameters,omitempty"`
}

// Trigger represents a single trigger in a Relay Workflow.
type Trigger struct {
	Name    string            `yaml:"name" validate:"required"`
	Source  map[string]string `yaml:"source" validate:"required"` // Can be one of [Schedule, Push, Webhook]
	Binding TriggerBinding    `yaml:"binding,omitempty"`
	When    string            `yaml:"when,omitempty"`
}

// Parameter represents a parameter in a Relay Workflow.
type Parameter struct {
	Default     string `yaml:"default,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type workflow struct {
	APIVersion  string                 `yaml:"apiVersion" validate:"required"`
	Kind        string                 `yaml:"kind" validate:"required"`
	Summary     string                 `yaml:"summary" validate:"required"`
	Description string                 `yaml:"description,omitempty"`
	Homepage    string                 `yaml:"homepage,omitempty"` // URI type
	Source      string                 `yaml:"source,omitempty"`   // URI type
	Tags        []string               `yaml:"tags,omitempty"`
	Locals      []map[string]string    `yaml:"locals,omitempty"`
	Parameters  []map[string]Parameter `yaml:"parameters,omitempty"` // maybe needs to be a struct
	Triggers    []Trigger              `yaml:"triggers"`
	Steps       []Step                 `yaml:"steps"`
}

// AddParameter adds a parameter to the parameter map.
func (w *workflow) AddParameter(key string, defaultValue string, description string) {
	p := map[string]Parameter{}
	p[key] = Parameter{
		Default:     defaultValue,
		Description: description,
	}

	w.Parameters = append(w.Parameters, p)
}

// AddTag adds a tag to the Workflow metadata.
func (w *workflow) AddTag(tag string) {
	w.Tags = append(w.Tags, tag)
}

// AddTrigger adds a trigger to the Workflow.
func (w *workflow) AddTrigger(trigger Trigger) {
	w.Triggers = append(w.Triggers, trigger)
}

// AddStep adds a step to the Workflow.
func (w *workflow) AddStep(step Step) {
	w.Steps = append(w.Steps, step)
}

// GetSteps returns a slice of Step.
func (w *workflow) GetSteps() []Step {
	return w.Steps
}

// Validate validates that certain properties are available in the struct that
// needs to be parsed.
// Maybe should use the api in the future https://github.com/puppetlabs/relay/blob/8fc48b290a327a4abcb31a188269717c72389120/pkg/client/revision.go#L13
func (w *workflow) Validate() error {
	validate := validator.New()
	err := validate.Struct(w)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return ErrValidation
		}

		fmt.Println("Could not generate workflow because of the following validation errors:")
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf(" -> Field %s [%s] %s\n", err.Field(), err.Type(), validationMap(err.Tag()))
		}
		return ErrValidation
	}

	return nil
}

// Write writes the YAML representation of the Workflow to the given writer.
// If no writer is provided, the output is written to stdout.
func (w *workflow) Write(writer io.Writer) error {
	err := w.Validate()
	if err != nil {
		return err
	}

	if writer == nil {
		writer = os.Stdout
	}

	b, err := yaml.Marshal(w)
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, string(b))
	return nil
}

// NewWorkflow creates an empty Relay Workflow.
// It returns a Workflow interface that can be used to add steps and triggers.
func NewWorkflow(summary string) Workflow {
	return &workflow{
		APIVersion: "v1",
		Kind:       "Workflow",
		Summary:    summary,
	}
}

// A simple vanity map for the tags returned by the validator
// to make the error messages more readable. It's almost certainly overkill.
func validationMap(tag string) string {
	reasons := map[string]string{
		"required": "was missing from the payload",
	}

	if reason, ok := reasons[tag]; ok {
		return reason
	} else {
		return tag
	}
}
