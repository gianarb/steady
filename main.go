package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gianarb/lb/api"
	"github.com/gianarb/lb/config"
	"github.com/gianarb/lb/core"
	"github.com/gianarb/lb/proxy"
)

var conf config.Configuration

func run(fr *core.Frontend) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		newRequest := proxy.CreateNewRequest(req, fr.Nodes)

		resp, err := http.Get(newRequest.URL.String())
		if err != nil {
			fmt.Printf("$s", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("$s", err)
		}
		proxy.HydrateHeaders(resp, w)
		w.Write([]byte(body))
	}
}

func main() {
	var configPath string

	cmdFlags := flag.NewFlagSet("event", flag.ContinueOnError)
	cmdFlags.StringVar(&configPath, "c", "/etc/lb.config.json", "c")

	if err := cmdFlags.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalln(err)
	}

	conf.Parse(configPath)
	var wg sync.WaitGroup

	if conf.RConf.Admin != nil {
		go api.Start(conf)
	}

	for name, frontend := range conf.Frontends {
		wg.Add(1)
		go func(fr *core.Frontend, n string) {
			defer wg.Done()
			log.Printf("Start %s on %s:%d", n, fr.Bind, fr.Port)
			err := http.ListenAndServe(fmt.Sprintf("%s:%d", fr.Bind, fr.Port), run(fr))
			if err != nil {
				log.Fatalln(err)
			}
		}(frontend, name)
	}
	wg.Wait()
}
