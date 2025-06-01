package utils

import (
	"fmt"
	"log/slog"
	"os/exec"
)

type ScriptStep struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

func executeScriptSteps(steps []ScriptStep, varsInput map[string]any, shell string) (string, error) {
	if shell == "" {
		return "", fmt.Errorf("shell must be specified")
	}
	// copy the vars to avoid mutation issues
	vars := make(map[string]any)
	for k, v := range varsInput {
		vars[k] = v
	}

	// Validate inputs
	if len(steps) == 0 {
		return "", fmt.Errorf("no steps provided to execute")
	}
	if _, ok := vars["StepOutputs"]; !ok {
		vars["StepOutputs"] = make(map[string]string)
	}
	if _, ok := vars["Arg"]; !ok {
		vars["Arg"] = "" // Ensure "Arg" is set in vars
	}
	if _, ok := vars["StepOutputs"].(map[string]string); !ok {
		return "", fmt.Errorf("StepOutputs must be a map[string]string")
	}
	slog.Debug("Executing script steps", "shell", shell, "stepsCount", len(steps))

	for i, step := range steps {

		// Check if "Arg" is set in vars and use it for the first step only
		if arg, ok := vars["Arg"].(string); ok && i == 0 {
			vars["Arg"] = arg
		} else {
			vars["Arg"] = ""
		}

		cmdStr, err := RenderTemplate(step.Script, vars)
		slog.Debug("Rendered script command for step", "step", step.Name, "cmd", cmdStr)
		if err != nil {
			return "", err
		}

		out, err := exec.Command(shell, "-c", cmdStr).CombinedOutput()
		slog.Debug("Executed script for step", "step", step.Name, "output", string(out), "error", err)

		if err != nil {
			slog.Error("Failed to execute script for step", "step", step.Name, "error", err)
			return "", fmt.Errorf("failed to execute script %s: %w", step.Name, err)
		}

		vars["StepOutputs"].(map[string]string)[step.Name] = string(out)
	}
	return vars["StepOutputs"].(map[string]string)[steps[len(steps)-1].Name], nil
}

func ExecuteBashScriptSteps(steps []ScriptStep, vars map[string]any) (string, error) {
	return executeScriptSteps(steps, vars, "bash")
}
