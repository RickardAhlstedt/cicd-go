package executor

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/RickardAhlstedt/cicd-go/cfg"
	"github.com/RickardAhlstedt/cicd-go/vars"
)

func RunBuild(config *cfg.BuildConfig, buildFile string, changedFile string, eventType string) {
	fmt.Println("🚀 Starting pipeline...")

	ctx := vars.BuildContext(config.Variables, "", buildFile, changedFile, eventType)

	runSteps(config.Setup, ctx)

	runSteps(config.Steps, ctx)

	runParallelSteps(config.Parallel, ctx)

	runSteps(config.PostBuild, ctx)

	fmt.Println("✅ Pipeline completed successfully!")
}

func runCommand(command string, ctx vars.RuntimeContext) {
	var cmd *exec.Cmd
	interpolateCmd := vars.Interpolate(command, ctx)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", interpolateCmd)
	} else {
		cmd = exec.Command("sh", "-c", interpolateCmd)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Error executing '%s': %v\n", interpolateCmd, err)
		fmt.Println("🚨 Step failed! Exiting...")
		os.Exit(1) // Exit with a failure code
	}
}

func shouldRunStep(condition string, ctx vars.RuntimeContext) bool {
	if condition == "" {
		return true
	}
	interpolated := vars.Interpolate(condition, ctx)
	interpolated = strings.TrimSpace(interpolated)

	operators := []string{"==", "!=", "^=", "$=", "*=", "~="}

	for _, op := range operators {
		if strings.Contains(interpolated, op) {
			parts := strings.SplitN(interpolated, op, 2)
			if len(parts) != 2 {
				return false
			}
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			switch op {
			case "==":
				return left == right
			case "!=":
				return left != right
			case "^=":
				return strings.HasPrefix(left, right)
			case "$=":
				return strings.HasSuffix(left, right)
			case "*=":
				return strings.Contains(left, right)
			case "~=":
				matched, err := regexp.MatchString(right, left)
				if err != nil {
					return false
				}
				return matched
			}
		}
	}

	return strings.ToLower(strings.TrimSpace(interpolated)) == "true"
}

func runSteps(steps []cfg.BuildStep, ctx vars.RuntimeContext) {
	for _, step := range steps {
		stepCtx := ctx
		stepCtx.StepName = step.Name
		if !shouldRunStep(step.If, stepCtx) {
			fmt.Printf("⏭️ Skipping step: %s as condition was not met\n", step.Name)
			continue
		}
		fmt.Printf("🔧 Running step: %s\n", step.Name)
		runCommand(step.Command, stepCtx)
	}
}

func runParallelSteps(parallelSteps []cfg.ParallelStep, ctx vars.RuntimeContext) {
	var wg sync.WaitGroup
	for _, p := range parallelSteps {
		fmt.Printf("🚀 Running parallel task: %s\n", p.Name)
		for _, cmd := range p.Commands {
			wg.Add(1)
			go func(command string, name string) {
				defer wg.Done()
				parallelCtx := ctx
				parallelCtx.StepName = name
				runCommand(command, parallelCtx)
			}(cmd, p.Name)
		}
	}
	wg.Wait()
}
