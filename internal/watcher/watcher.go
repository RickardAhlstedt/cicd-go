package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/RickardAhlstedt/cicd-go/cfg"
	"github.com/RickardAhlstedt/cicd-go/internal/executor"
	"github.com/fsnotify/fsnotify"
)

func StartWatching(rootDir string, buildFile *string) {
	println("Loading config", *buildFile)
	configPath := filepath.Join(rootDir, *buildFile)
	config, err := cfg.LoadConfig(configPath)
	println("Loaded config")
	if err != nil {
		log.Fatalf("Failed to load pipeline: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("⚠️ Failed to close watcher: %v", err)
		}
	}()

	var (
		lastChangedFile string
		lastEventType   string
		debounceMu      sync.Mutex
		debounceTimer   *time.Timer
	)

	watchedDir := make(map[string]struct{})

	var watchDir func(string) error
	watchDir = func(dir string) error {
		absDir, _ := filepath.Abs(dir)
		if _, watched := watchedDir[absDir]; watched {
			return nil
		}
		if config.ShouldIgnore(dir) {
			return nil
		}

		err := watcher.Add(dir)
		if err != nil {
			return err
		}
		watchedDir[absDir] = struct{}{}
		fmt.Println("Watching: ", absDir)

		entries, err := os.ReadDir(absDir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				subDir := filepath.Join(absDir, entry.Name())
				err = watchDir(subDir)
				if err != nil {
					log.Println("Error watching subdirectory: ", subDir, err)
				}
			}
		}
		return nil
	}

	err = watchDir(rootDir)
	if err != nil {
		log.Fatal("Error setting up watchers: ", err)
	}

	triggerBuild := func() {
		debounceMu.Lock()
		defer debounceMu.Unlock()
		if lastChangedFile != "" {
			fmt.Println("🔁 Triggering pipeline due to file changes...")
			executor.RunBuild(config, configPath, lastChangedFile, lastEventType)
			lastChangedFile = ""
			lastEventType = ""
		}
	}

	debounceTimer = time.NewTimer(time.Hour)
	debounceTimer.Stop()
	// <-debounceTimer.C

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
					if err := watchDir(event.Name); err != nil {
						log.Printf("⚠️ Failed to watch new dir %s: %v", event.Name, err)
					}
				}
			}

			if event.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
				absPath, _ := filepath.Abs(event.Name)
				if _, ok := watchedDir[absPath]; ok {
					delete(watchedDir, absPath)
					if err := watcher.Remove(absPath); err != nil {
						log.Printf("⚠️ Failed to remove watcher for %s: %v", absPath, err)
					}
					fmt.Println("🗑️ Unwatched: ", absPath)
				}
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Remove) != 0 {
				lastChangedFile = event.Name
				lastEventType = event.Op.String()

				debounceMu.Lock()
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(1*time.Second, triggerBuild)
				debounceMu.Unlock()
			}
		case <-debounceTimer.C:
			triggerBuild()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("🚨 Watcher error: ", err)
		}
	}
}
