package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile"
	profiledao "coolcar/rental/profile/dao"
	"coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profClient "coolcar/rental/trip/client/profile"
	tripdao "coolcar/rental/trip/dao"
	"coolcar/shared/server"
	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var addr = flag.String("addr", ":8082", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri")
var blobAddr = flag.String("blob_addr", "localhost:8083", "address for blob service")
var carAddr = flag.String("car_addr", "localhost:8084", "address for car service")
var authPublicKeyFile = flag.String("auth_public_key_file", "shared/auth/public.key", "public key file for auth")

func main() {
	flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()

	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	blobConn, err := grpc.Dial(*blobAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	profService := &profile.Service{
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		PhotoGetExpire:    time.Hour,
		PhotoUploadExpire: time.Hour,
		Mongo:             profiledao.NewMongo(mongoClient.Database("coolcar")),
		Logger:            logger,
	}

	carConn, err := grpc.Dial(*carAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              *addr,
		AuthPublicKeyFile: *authPublicKeyFile,
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				CarManager: &car.Manager{
					CarService: carpb.NewCarServiceClient(carConn),
				},
				POIManager:   &poi.Manager{},
				DistanceCalc: nil,
				Mongo:        tripdao.NewMongo(mongoClient.Database("coolcar")),
				Logger:       logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))
}
