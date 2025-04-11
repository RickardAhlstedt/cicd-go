package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/RickardAhlstedt/cicd-go/cfg"
	"github.com/RickardAhlstedt/cicd-go/internal/executor"
	"github.com/RickardAhlstedt/cicd-go/internal/watcher"
)

const VERSION = "1.2"

func main() {

	buildFile := flag.String("file", "build.yaml", "Specify a custom file for pipeline")
	watch := flag.Bool("watch", false, "Watch the directory for changes")
	flag.Parse()

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Error getting working directory: %v", err)
	}

	configPath := filepath.Join(workingDir, *buildFile)
	config, err := cfg.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ Failed to load build.yaml:\n%v", err)
	}

	if config.ConfigVersion < VERSION {
		fmt.Printf("⚠️ %s is using an older version than the binary was compiled for.\nCurrent binary-version: %s, file-version: %s\n", *buildFile, VERSION, config.ConfigVersion)
	}

	config.Ignore = append(config.Ignore, "build.yaml")

	if *watch {
		fmt.Println("ℹ️ Watching for file change in: ", workingDir)
		watcher.StartWatching(workingDir, buildFile)
	} else {
		fmt.Printf("ℹ️ Working in: %s\n", workingDir)
		config, err := cfg.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("ℹ️ Failed to load build.yaml: %v", err)
		}
		executor.RunBuild(config, *buildFile, "", "")
	}
}
