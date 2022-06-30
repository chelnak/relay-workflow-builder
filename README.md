# Relay Workflow Builder

## What is Relay Workflow Builder

Relay Workflow Builder will help you programatically create workflows that contain many steps and triggers.

The package exports a Workflow interface that provides methods for adding Triggers and Steps.

## Usage

`main.go` serves as a good example of using this package.

The `Print` method writes directly to `stdout`, therefore output from the tool can be piped in to a new file:

```bash
go run . > my_workflow.yaml
```

Alternatively you can pass any `io.Writer` to handle the output.

Workflows that have been created with this tool can be pushed in to the Relay service with the CLI:

```bash
relay workflow save "my-workflow" -f ./my_workflow.yaml
```
