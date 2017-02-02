package main

import (
	"os"
	"errors"
	"net/url"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type RubyApplicationtSet struct {
	applications map[string]Application
}

func (set RubyApplicationtSet) ScanForApplications(dir string, startingPort *int, scannedSet ApplicationtSet) (ApplicationtSet, error) {
	if set.applications == nil {
		set.applications = make(map[string]Application)
	}

	files, err := ioutil.ReadDir(dir + "/")

	if err != nil {
		return set, err
	}

	for _, fileInfo := range files {
		path := dir + "/" + fileInfo.Name()

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			destPath, err := os.Readlink(path)

			if err != nil {
				return nil, err
			}

			name := filepath.Base(destPath)
			domain := fileInfo.Name()
			
			app := Application{Name: name, Dir: destPath, Port: *startingPort}
			hasConfigRu, err := app.HasFile("config.ru")

			if err != nil {
				return set, err
			}

			if scannedSet != nil && scannedSet.HasAppWithDomain(domain) {
				continue;
			}

			if hasConfigRu {
				command := "bundle exec rackup -p " + strconv.Itoa(*startingPort)
				proc := Proc{Dir: app.Dir, Command: command}
				app.Procs = append(app.Procs, proc)
				
				set.applications[domain] = app
				*startingPort++
			}
		}
	}
	
	return set, nil
}

func (set RubyApplicationtSet) PortFor(host string, url *url.URL, scannedSet ApplicationtSet) (int, error) {
	domain, err := HostWithoutPortOrTld(host)

	if err != nil {
		return 0, err
	}

	app, ok := set.applications[domain]

	if ok {
		return app.Port, nil
	} else {
		return 0, errors.New("No matching application found.")
	}
}

func (set RubyApplicationtSet) HasAppWithDomain(domain string) (bool) {
	_, ok := set.applications[domain]
	return ok
}

func (set RubyApplicationtSet) Applications() (map[string]Application) {
	return set.applications
}
