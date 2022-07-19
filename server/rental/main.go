package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile"
	profiledao "coolcar/rental/profile/dao"
	"coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profClient "coolcar/rental/trip/client/profile"
	tripdao "coolcar/rental/trip/dao"
	mgo "coolcar/shared/mongo"
	"coolcar/shared/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()

	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(mgo.MongoURL))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	blobConn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	profService := &profile.Service{
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		PhotoGetExpire:    time.Hour,
		PhotoUploadExpire: time.Hour,
		Mongo:             profiledao.NewMongo(mongoClient.Database("coolcar")),
		Logger:            logger,
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "shared/auth/public.key",
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				CarManager:   &car.Manager{},
				POIManager:   &poi.Manager{},
				DistanceCalc: nil,
				Mongo:        tripdao.NewMongo(mongoClient.Database("coolcar")),
				Logger:       logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))
}
