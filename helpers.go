package main

import (
	"context"

	"log"
	"path/filepath"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ansrivas/fwatcher/messages"
	"github.com/fsnotify/fsnotify"
)

func isFileTypeAllowed(ext string, allowedExt []string) bool {
	for _, e := range allowedExt {
		if ext == e {
			return true
		}
	}
	return false
}

func watchDirectory(ctx context.Context, dirToWatch string, allowedExt []string, pid *actor.PID) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(dirToWatch)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:

			log.Println(event.Op)
			if event.Op == fsnotify.Create {
				log.Println("File created...")
				fileext := filepath.Ext(event.Name)
				if isFileTypeAllowed(fileext, allowedExt) {
					pid.Tell(&messages.FileModified{Filepath: event.Name})
				} else {
					log.Printf("Skipping this as, file type %v not allowed", fileext)
				}
			}

		case err := <-watcher.Errors:
			log.Println("error:", err)

		case <-ctx.Done():
			return
		}
	}
}
