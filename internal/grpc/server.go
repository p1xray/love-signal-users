package usersserver

import (
	"context"
	"love-signal-users/internal/dto"

	userspb "github.com/p1xray/love-signal-protos/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	emptyValue = 0
)

// UsersService is a service for working with user data.
type UsersService interface {
	// UserInfo returns information about a user by their external identifier.
	UserInfo(
		ctx context.Context,
		userExternalId int64,
	) (*dto.UserInfo, error)

	// UserProfileCard returns the user profile card.
	UserProfileCard(
		ctx context.Context,
		userId int64,
	) (*dto.UserProfileCard, error)
}

type serverAPI struct {
	userspb.UnimplementedUsersServer
	users UsersService
}

// Register registers the implementation of the API service with the gRPC server.
func Register(gRPC *grpc.Server, users UsersService) {
	userspb.RegisterUsersServer(gRPC, &serverAPI{users: users})
}

// UserInfo returns information about a user by their external identifier.
func (s *serverAPI) UserInfo(
	ctx context.Context,
	req *userspb.UserInfoRequest,
) (*userspb.UserInfoResponse, error) {
	if err := validateUserInfoRequest(req); err != nil {
		return nil, err
	}

	userInfo, err := s.users.UserInfo(ctx, req.GetUserExternalId())
	if err != nil {
		return nil, internalError("failed to get user info")
	}

	userInfoResponse := &userspb.UserInfoResponse{Id: userInfo.Id, Name: userInfo.Name}
	return userInfoResponse, nil
}

// UserProfileCard returns the user profile card.
func (s *serverAPI) UserProfileCard(
	ctx context.Context,
	req *userspb.UserProfileCardRequest,
) (*userspb.UserProfileCardResponse, error) {
	if err := validateUserProfileCardRequest(req); err != nil {
		return nil, err
	}

	card, err := s.users.UserProfileCard(ctx, req.GetId())
	if err != nil {
		return nil, internalError("failed to get user profile card")
	}

	cardResponse := &userspb.UserProfileCardResponse{
		Id:          card.Id,
		Name:        card.Name,
		DateOfBirth: timestamppb.New(card.DateOfBirth),
		Gender:      userspb.Gender(card.Gender),
	}

	return cardResponse, nil
}

func validateUserInfoRequest(req *userspb.UserInfoRequest) error {
	if req.GetUserExternalId() == emptyValue {
		return invalidArgumentError("user_external_id is empty")
	}

	return nil
}

func validateUserProfileCardRequest(req *userspb.UserProfileCardRequest) error {
	if req.GetId() == emptyValue {
		return invalidArgumentError("id is empty")
	}

	return nil
}

func invalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

func internalError(msg string) error {
	return status.Error(codes.Internal, msg)
}
