package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/fem_project/internal/app"
	"github.com/alireza-akbarzadeh/fem_project/internal/constants"
	"github.com/alireza-akbarzadeh/fem_project/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend  server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(fmt.Sprintf(constants.Red+"💥 Failed to initialize application: %v"+constants.Reset, err))
	}
	defer app.DB.Close()

	r := routes.SetupRoute(app)
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf(constants.Blue+"🚀 Server is running! Open in your browser: http://localhost:%d\n"+constants.Reset, port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf(constants.Red+"❌ Server failed: %v"+constants.Reset, err)
	}
}
