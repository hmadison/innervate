package main

import (
	"flag"
	"github.com/eliasgs/mdns"
	"github.com/fsnotify/fsnotify"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
)

var configDir string
var proxyPort int
var childStartPort int
var tld string
var restartOnConfigChanges bool
var wg sync.WaitGroup

func init() {
	user, err :=  user.Current()

	if err != nil {
		panic(err)
	}
	
	flag.IntVar(&proxyPort, "port", 8080, "port to run the proxy server on")
	flag.StringVar(&configDir, "config", user.HomeDir + "/.config/innervate", "path to the configuration directory")
	flag.IntVar(&childStartPort, "start-port", 10000, "port to start running child processes on")
	flag.StringVar(&tld, "tld", "localhost", "the tld to use")
	flag.BoolVar(&restartOnConfigChanges, "restart", false, "should the app restart when there is a config change")
	flag.Parse()
}

func main() {
	cps := childStartPort
	apps := ParseConfig(configDir, &cps)

	zone, err := mdns.New()

	if err != nil {
		panic(err)
	}

	for name, app := range apps {
		domain := name + "." + tld
		log.Printf(" '%s' @ 'http://%s:%d'", app.Name, domain, proxyPort)
		zone.Publish(domain + " 60 IN A 127.0.0.1")
		app.StartChildren()
	}

	wg.Add(1)
	go func() {
		proxy := newReverseProxy(&apps)
		http.ListenAndServe(":"+strconv.Itoa(proxyPort), proxy)
		wg.Done()
	}()

	if restartOnConfigChanges {
		log.Printf("Watching %s for config changes", configDir)
		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			panic(err)
		}

		watcher.Add(configDir)

		wg.Add(1)
		go watchConfigDirForModification(watcher)
	}

	wg.Wait()
}

func getAppNameFromHost(input string) (appName string) {
	host := input

	if strings.Contains(input, ":") {
		splitHost, _, err := net.SplitHostPort(input)
		host = splitHost

		if err != nil {
			panic(err)
		}
	}

	appName = strings.TrimSuffix(host, "."+tld)
	return
}

func newReverseProxy(apps *map[string]Application) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		appName := getAppNameFromHost(req.Host)

		app, ok := (*apps)[appName]

		if ok {
			req.URL.Scheme = "http"
			req.URL.Host = "localhost:" + strconv.Itoa(app.Port)
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

func watchConfigDirForModification(watcher *fsnotify.Watcher) {
loop:
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == 0 {
				log.Printf("Restarting")
				os.Exit(0)
			}
		case err := <-watcher.Errors:
			if err == nil {
				break loop
			}
		}
	}
	wg.Done()
}
