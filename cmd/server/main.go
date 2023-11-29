package main

import (
	"github.com/bekarys11/drive-uploader/pkg/gdrive"
	"github.com/fsnotify/fsnotify"
	"log"
	"math"
	"path/filepath"
	"time"
)

var fileName *string // video file's location path

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go waitForRecorder(watcher)

	p := filepath.Dir("/home/bekarys/Videos/")
	err = watcher.Add(p)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func waitForRecorder(w *fsnotify.Watcher) {
	var timer = time.NewTimer(math.MaxInt64) // create timer with max duration to wait for write operations

	for {
		select {
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Println("ERROR:", err)
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) { // if it is not create or write operation, then restart loop
				continue
			}

			if e.Has(fsnotify.Write) {
				time.Sleep(10 * time.Second)  // not reset immediately on write operation, because write op could be fired several times in 1 second
				timer.Reset(10 * time.Second) // reset enough time for next write operation to occur
				fileName = &e.Name
			}
		case <-timer.C: // if timer is up
			log.Println("timer is up")
			gdrive.ConnectToDrive(fileName)
		}
	}
}
