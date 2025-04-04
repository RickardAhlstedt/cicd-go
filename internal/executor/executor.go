package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/RickardAhlstedt/cicd-go/cfg"
)

func RunBuild(config *cfg.BuildConfig) {
	fmt.Println("🚀 Starting pipeline...")

	runSteps(config.Setup)

	runSteps(config.Steps)

	runParallelSteps(config.Parallel)

	runSteps(config.PostBuild)

	fmt.Println("✅ Pipeline completed successfully!")
}

func runCommand(command string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Error executing '%s': %v\n", command, err)
		fmt.Println("🚨 Step failed! Exiting...")
		os.Exit(1) // Exit with a failure code
	}
}

func shouldRunStep(condition string) bool {
	if condition == "" {
		return true
	}

	parts := strings.Split(condition, "==")
	if len(parts) != 2 {
		fmt.Println("⚠️ Invalid condition: ", condition)
		return false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	envValue := os.Getenv(strings.TrimPrefix(key, "env."))
	return envValue == strings.Trim(value, "'\"")

}

func runSteps(steps []cfg.BuildStep) {
	for _, step := range steps {
		if !shouldRunStep(step.If) {
			fmt.Printf("⏭️ Skipping step: %s as condition was not met\n", step.Name)
			continue
		}
		fmt.Printf("🔧 Running step: %s\n", step.Name)
		runCommand(step.Command)
	}
}

func runParallelSteps(parallelSteps []cfg.ParallelStep) {
	var wg sync.WaitGroup
	for _, p := range parallelSteps {
		fmt.Printf("🚀 Running parallel task: %s\n", p.Name)
		for _, cmd := range p.Commands {
			wg.Add(1)
			go func(command string) {
				defer wg.Done()
				runCommand(command)
			}(cmd)
		}
	}
	wg.Wait()
}
