package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func Test(endpoint string, code string) (res *authpb.LoginResponse) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	authClient := authpb.NewAuthServiceClient(conn)
	r, err := authClient.Login(context.Background(), &authpb.LoginRequest{
		Code: code,
	})
	if err != nil {
		log.Println(err)
		return
	}
	return r
}
