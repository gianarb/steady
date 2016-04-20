package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gianarb/lb/config"
	"github.com/gorilla/mux"
)

type Api struct {
}

func Start(c config.Configuration) {
	log.Printf("Start api system on %s:%d", c.RConf.Bind, c.RConf.Port)
	r := mux.NewRouter()
	r.HandleFunc("/ping", PingHandler()).Methods("GET")
	r.HandleFunc("/backup", BackupHandler(c)).Methods("GET")
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", c.RConf.Bind, c.RConf.Port), r)
	if err != nil {
		log.Fatalln(err)
	}
}