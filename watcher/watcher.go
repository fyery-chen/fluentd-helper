package watcher

import (
	"github.com/Sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/rancher/fluentd-helper/helper"
)

func Watcherfile(path string, done <-chan int) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Errorf("fsnotify NewWatcher failed %v", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				logrus.Infof("------- receive file event ----, name: %#v, op: %#v, event: %#v", event.Name, event.Op.String())
				if event.Op == fsnotify.Remove {
					watcher.Remove(event.Name)
					watcher.Add(event.Name)
					logrus.Infof("receive remove event: %#v, %#v", event.Name, event.Op.String())
					helper.ReloadFluentd()
				} else if event.Op == fsnotify.Write {
					logrus.Infof("receive write event: %#v, %#v", event.Name, event.Op.String())
					helper.RenewTicket()
					helper.ReloadFluentd()
				}
				// watch for errors
			case err := <-watcher.Errors:
				logrus.Errorf("file watcher get error: %v", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(path); err != nil {
		logrus.Errorf("file watcher add path error: %v", err)
	}

	<-done
}
