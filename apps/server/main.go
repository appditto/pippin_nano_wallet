package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/appditto/pippin_nano_wallet/apps/server/controller"
	"github.com/appditto/pippin_nano_wallet/libs/config"
	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/go-chi/chi/v5"
	"k8s.io/klog/v2"
)

var Version = "dev"

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	// Server options
	flag.Usage = usage
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	flag.Set("v", "3")
	// if utils.GetEnv("ENVIRONMENT", "development") == "development" {
	// 	flag.Set("stderrthreshold", "INFO")
	// 	flag.Set("v", "3")
	// }
	version := flag.Bool("version", false, "Display the version")
	flag.Parse()

	if *version {
		fmt.Printf("Pippin version: %s\n", Version)
		os.Exit(0)
	}

	// Read yaml configuration
	conf, err := config.ParsePippinConfig()
	if err != nil {
		klog.Fatalf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	// Setup database conn
	ctx := context.Background()
	fmt.Println("🏡 Connecting to database...")
	dbconn, err := database.GetSqlDbConn(false)
	if err != nil {
		klog.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	entClient, err := database.NewEntClient(dbconn)
	if err != nil {
		klog.Fatalf("Failed to create ent client: %v", err)
		os.Exit(1)
	}
	defer entClient.Close()

	// Run migrations
	klog.Infoln("🦋 Running migrations...")
	if err := entClient.Schema.Create(ctx); err != nil {
		klog.Fatalf("Failed to run migrations: %v", err)
		os.Exit(1)
	}

	// Create app
	app := chi.NewRouter()

	// Setup controller
	hc := controller.HttpController{}

	// HTTP Routes
	app.Post("/", hc.HandleAction)

	http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port), app)
}
