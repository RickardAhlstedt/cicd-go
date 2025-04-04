package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/RickardAhlstedt/cicd-go/cfg"
	"github.com/RickardAhlstedt/cicd-go/internal/executor"
	"github.com/fsnotify/fsnotify"
)

func StartWatching(rootDir string) {
	configPath := filepath.Join(rootDir, "build.yaml")
	config, err := cfg.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load buildy.yaml: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var watchDir func(string) error
	watchDir = func(dir string) error {
		if config.ShouldIgnore(dir) {
			fmt.Println("Ignoring: ", dir)
		}

		err := watcher.Add(dir)
		if err != nil {
			return err
		}
		fmt.Println("Watching: ", dir)

		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				subDir := filepath.Join(dir, entry.Name())
				err = watchDir(subDir)
				if err != nil {
					log.Println("Error watching subdirectory: ", subDir, err)
				}
			}
		}
		return nil
	}

	err = watcher.Add(rootDir)
	if err != nil {
		log.Fatal("Error setting up watchers: ", err)
	}

	debounceTimer := time.NewTimer(0)
	<-debounceTimer.C

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if config.ShouldIgnore(event.Name) {
				continue
			}

			if event.Op&fsnotify.Create != 0 {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					err := watchDir(event.Name)
					if err != nil {
						log.Println("Error adding new directory: ", err)
					}
				}
			}

			if event.Op&(fsnotify.Write|fsnotify.Rename) != 0 {
				fmt.Println("File changed: ", event.Name)
				debounceTimer.Reset(1 * time.Second)
			}
		case <-debounceTimer.C:
			fmt.Println("Triggering build due to file changes...")
			// configPath := filepath.Join(dir, "build.yaml")
			config, err := cfg.LoadConfig(configPath)
			if err != nil {
				log.Println("Error loading build.yaml: ", err)
				continue
			}
			executor.RunBuild(config)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error: ", err)
		}
	}
}
