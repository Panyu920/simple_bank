package gapi

import (
	"context"
	"database/sql"
	db "simple_bank/db/sqlc"
	"simple_bank/pb"
	"simple_bank/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "%v", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:                    refreshPayload.ID,
		Username:              refreshPayload.Username,
		UserAgent:             mtdt.UserAgent,
		UserIp:                mtdt.UserIp,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshPayload.ExpiresAt.Time,
		IsBlocked:             false,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	resp := &pb.LoginUserResponse{
		Session_Id:            session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
		// createUserResponse:    *fromDBUserTocreateUserResponse(&user),
		User: convertUser(&user),
	}

	return resp, nil
}
