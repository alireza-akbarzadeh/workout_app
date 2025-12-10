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

// Package main Fem Project API
//
//  @title           Fem Project API Documentation
//  @version         1.0
//  @description     Fem Project is a social and personal empowerment platform built to support users in creating profiles, interacting with features, and managing their accounts seamlessly.
//  @description     This backend API provides secure endpoints for user registration, authentication, account management, and future social features.
//  @description     All responses are formatted using JSON and follow RESTful standards.
//
//  @termsOfService  https://femproject.com/terms
//
//  @contact.name     Fem Project Support Team
//  @contact.url      https://femproject.com/support
//  @contact.email    support@femproject.com
//
//  @license.name     MIT License
//  @license.url      https://opensource.org/licenses/MIT
//
//  @host       localhost:8080
//  @BasePath   /api/v1
//  @schemes    http https

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
