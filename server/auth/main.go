package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"coolcar/shared/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var addr = flag.String("addr", ":8081", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri")
var privateKeyFile = flag.String("private_key_file", "shared/auth/private.key", "private key file")

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

	pkFile, err := os.Open(*privateKeyFile)
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name: "auth",
		Addr: *addr,
		RegisterFunc: func(g *grpc.Server) {
			authpb.RegisterAuthServiceServer(g, &auth.Service{
				OpenIDResolver: &wechat.Service{},
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
				Logger:         logger,
				TokenExpire:    20000 * time.Hour,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privKey),
			})
		},
		Logger: logger,
	}))
}
