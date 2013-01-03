package main

import (
	"github.com/howeyc/fsnotify"
	"log"

	)

// Spin up a goroutine that watches two channels:
// watcher.Event for events of [Delete, Modify, Moved, New] and
// watcher.Error for any errors behind the scenes.
func watch() {
	// Straight from the author of github.com/howeyc/fsnotify:
	// Initialize our watcher and check for any errors...
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		select {
		case ev := <-watcher.Event:
	
			if *verbose {
				log.Printf("event:%v", ev)
			}
	
			if ev.IsModify() || ev.IsCreate() {
				push(ev.Name)
			}
	
			if ev.IsDelete() {
				drop(ev.Name)
			}

			// TODO: Need to work out how to know the old name and the new.
			// If it isn't supported in fsnotify then I'll need to make
			// a routine scan the present state and compare that to the
			// database.
			/*
				if ev.IsRename() {
					push(ev.Name)
				}
			*/

		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
	
	log.Print("setting up watcher")
	
	err = watcher.Watch(*content)
	if err != nil {
		log.Fatal(err)
	}
	
	defer watcher.Close()
}