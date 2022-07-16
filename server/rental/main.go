package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	"coolcar/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Fatal("cannot listen", zap.Error(err))
	}

	in, err := auth.Interceptor("shared/auth/public.key")
	if err != nil {
		logger.Fatal("cannot create auth interceptor", zap.Error(err))
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(in))
	rentalpb.RegisterTripServiceServer(s, &trip.Service{Logger: logger})

	err = s.Serve(lis)
	logger.Fatal("cannot server", zap.Error(err))
}
