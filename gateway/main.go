package main

import (
	"context"
	"log"
	"net/http"
	"time"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/commons/discovery"
	"github.com/Far-sa/commons/discovery/consul"
	"github.com/Far-sa/gateway/gateway"
	"github.com/Far-sa/gateway/handler"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":3000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Unregister(ctx, instanceID, serviceName)

	orderGateway := gateway.NewGRPCGateway(registry)

	mux := http.NewServeMux()
	handler := handler.New(orderGateway)
	handler.RegisterRoutes(mux)

	log.Printf("starting HTTP server on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start server", err)
	}
}
