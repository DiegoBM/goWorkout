package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/DiegoBM/goWorkout/internal/app"
	"github.com/DiegoBM/goWorkout/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	routes := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("Server started in port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
