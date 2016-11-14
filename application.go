package main

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Application struct {
	Name    string
	Port    int
	Dir     string
	Procs   []Proc
	mu      sync.Mutex
	watcher *fsnotify.Watcher
}

func ParseConfig(dir string, startingPort *int) (apps map[string]Application) {
	apps = make(map[string]Application)
	files, err := ioutil.ReadDir(dir + "/")

	if err != nil {
		panic(err)
	}

	for _, fileInfo := range files {
		path := dir + "/" + fileInfo.Name()

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			destPath, err := os.Readlink(path)

			if err != nil {
				panic(err)
			}

			name := filepath.Base(destPath)
			domain := fileInfo.Name()

			app := Application{Name: name, Dir: destPath, Port: *startingPort}

			hasProcfile, err := app.HasFile("Procfile")

			if err != nil {
				panic(err)
			}

			hasConfigRu, err := app.HasFile("config.ru")

			if err != nil {
				panic(err)
			}

			if hasProcfile {
				parsed, err := ParseProcfile(destPath + "/Procfile")

				if err != nil {
					panic(err)
				}

				for _, command := range parsed {
					command = strings.Replace(command, "$PORT", strconv.Itoa(*startingPort), -1)
					proc := Proc{Dir: app.Dir, Command: command}
					app.Procs = append(app.Procs, proc)
				}
			} else if hasConfigRu {
				command := "bundle exec rackup -p " + strconv.Itoa(*startingPort)
				proc := Proc{Dir: app.Dir, Command: command}
				app.Procs = append(app.Procs, proc)
			}

			apps[domain] = app
			*startingPort++
		}
	}
	return
}

func (app *Application) HasFile(path string) (bool, error) {
	_, err := os.Stat(app.Dir + "/" + path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (app *Application) StartChildren() {
	app.mu.Lock()
	defer app.mu.Unlock()

	hasFile, err := app.HasFile("tmp")

	if err != nil {
		panic(err)
	}

	if hasFile {
		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			panic(err)
		}

		wg.Add(1)
		go watchAppForRestartRequests(app, watcher)
		watcher.Add(app.Dir + "/tmp")
		app.watcher = watcher
	}

	for _, proc := range app.Procs {
		proc.Start()
	}
}

func (app *Application) RestartChildren() {
	app.mu.Lock()
	defer app.mu.Unlock()

	for _, proc := range app.Procs {
		proc.Restart()
	}
}

func (app *Application) StopChildren() {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.watcher != nil {
		app.watcher.Close()
		app.watcher = nil
	}

	for _, proc := range app.Procs {
		proc.Stop()
	}
}

func watchAppForRestartRequests(app *Application, watcher *fsnotify.Watcher) {
loop:
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod != 0 || event.Op&fsnotify.Create != 0 {
				name := filepath.Base(event.Name)
				if name == "restart.txt" {
					app.RestartChildren()
				}
			}
		case err := <-watcher.Errors:
			if err == nil {
				break loop
			}
		}
	}
	wg.Done()
}
