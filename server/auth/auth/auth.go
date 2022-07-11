package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"go.uber.org/zap"
)

// Service implements auth service.
type Service struct {
	authpb.UnimplementedAuthServiceServer
	Logger *zap.Logger
}

//// OpenIDResolver resolves an authorization code
//// to an open id.
//type OpenIDResolver interface {
//	Resolve(code string) (string, error)
//}
//
//// TokenGenerator generates a token for the specified account.
//type TokenGenerator interface {
//	GenerateToken(accountID string, expire time.Duration) (string, error)
//}

// Login logs a user in.
func (s *Service) Login(_ context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", req.Code))
	return &authpb.LoginResponse{
		AccessToken: "token for " + req.Code,
		ExpiresIn:   req.Time,
	}, nil
	//openID, err := s.OpenIDResolver.Resolve(req.Code)
	//if err != nil {
	//	return nil, status.Errorf(codes.Unavailable,
	//		"cannot resolve openid: %v", err)
	//}
	//
	//accountID, err := s.Mongo.ResolveAccountID(c, openID)
	//if err != nil {
	//	s.Logger.Error("cannot resolve account id", zap.Error(err))
	//	return nil, status.Error(codes.Internal, "")
	//}
	//
	//tkn, err := s.TokenGenerator.GenerateToken(accountID.String(), s.TokenExpire)
	//if err != nil {
	//	s.Logger.Error("cannot generate token", zap.Error(err))
	//	return nil, status.Error(codes.Internal, "")
	//}
	//
	//return &authpb.LoginResponse{
	//	AccessToken: tkn,
	//	ExpiresIn:   int32(s.TokenExpire.Seconds()),
	//}, nil
}
