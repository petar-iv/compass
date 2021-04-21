package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"testing"
	"time"
)

type Server struct {
	*http.Server
}

func TestRunMain(t *testing.T) {
	go main()

	srv := &Server{Server: &http.Server{
		Addr:    ":5003",
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
			}
		}()
	})
	srv.Handler = router

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("TEST TEST")
}
