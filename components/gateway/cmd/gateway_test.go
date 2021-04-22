package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/vrischmann/envconfig"
	"net/http"
	"testing"
	"time"
)

type Config struct {
	Port string `envconfig:"default=:3002,APP_COVERAGE_PORT"`
}

type Server struct {
	*http.Server
}

func TestRunMain(t *testing.T) {
	go main()

	cfg := &Config{}
	if err := envconfig.Init(cfg); err != nil {
		log.D().Fatal(err)
	}

	srv := &Server{Server: &http.Server{
		Addr:    cfg.Port,
		Handler: nil,
	}}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		go func() {
			time.Sleep(time.Second * 3)
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
			defer cancelFunc()
			err := srv.Shutdown(ctx)
			if err != nil {
				log.D().Fatal(err)
			}
		}()
	})
	srv.Handler = router

	err := srv.ListenAndServe()
	if err != nil {
		log.D().Fatal(err)
	}
}
