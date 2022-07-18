package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/server"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

//var addr = flag.String("addr", ":8080", "address to listen")
var authAddr = flag.String("auth_addr", "localhost:8081", "address for auth service")
var tripAddr = flag.String("trip_addr", "localhost:8082", "address for trip service")

var profileAddr = flag.String("profile_addr", "localhost:8082", "address for profile service")

//var carAddr = flag.String("car_addr", "localhost:8084", "address for car service")

func main() {
	lg, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}

	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux()

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         *authAddr,
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "trip",
			addr:         *tripAddr,
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
		{
			name:         "profile",
			addr:         *profileAddr,
			registerFunc: rentalpb.RegisterProfileServiceHandlerFromEndpoint,
		},
	}

	for _, s := range serverConfig {
		err := s.registerFunc(
			c, mux, s.addr,
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			lg.Sugar().Fatalf("cannot register service %s: %v", s.name, err)
		}
	}
	addr := ":8080"
	lg.Sugar().Infof("grpc gateway started at %s", addr)
	lg.Sugar().Fatal(http.ListenAndServe(addr, mux))
}
