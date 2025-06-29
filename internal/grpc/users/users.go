package usersserver

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	lsuserspb "github.com/p1xray/love-signal-protos/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	server "love-signal-users/internal/grpc"
	"love-signal-users/internal/service"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	lsuserspb.UnimplementedUsersServer
	users server.UsersService
}

// Register registers the implementation of the API service with the gRPC server.
func Register(gRPC *grpc.Server, users server.UsersService) {
	lsuserspb.RegisterUsersServer(gRPC, &serverAPI{users: users})
}

// GetUserData returns information about a user by their ID.
func (s *serverAPI) GetUserData(
	ctx context.Context,
	req *lsuserspb.GetUserDataRequest,
) (*lsuserspb.UserDataResponse, error) {
	if err := validateGetUserDataRequest(req); err != nil {
		return nil, err
	}

	userData, err := s.users.UserData(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, server.NotFoundError("user not found")
		}

		return nil, server.InternalError("error getting user data")
	}

	var dateOfBirthPb *timestamppb.Timestamp
	if userData.DateOfBirth != nil {
		dateOfBirthPb = timestamppb.New(*userData.DateOfBirth)
	}

	genderPb := lsuserspb.Gender_GENDER_UNSPECIFIED
	if userData.Gender != nil {
		genderPb = lsuserspb.Gender(*userData.Gender)
	}

	var avatarFileKeyPb *wrappers.StringValue
	if userData.AvatarFileKey != nil {
		avatarFileKeyPb = &wrappers.StringValue{Value: *userData.AvatarFileKey}
	}

	userDataResponse := &lsuserspb.UserDataResponse{
		Id:            userData.ID,
		FullName:      userData.FullName,
		DateOfBirth:   dateOfBirthPb,
		Gender:        genderPb,
		AvatarFileKey: avatarFileKeyPb,
	}
	return userDataResponse, nil
}

// GetUserDataByExternalId returns information about a user by their external ID.
func (s *serverAPI) GetUserDataByExternalId(
	ctx context.Context,
	req *lsuserspb.GetUserDataByExternalIdRequest,
) (*lsuserspb.UserDataResponse, error) {
	if err := validateGetUserDataByExternalIdRequest(req); err != nil {
		return nil, err
	}

	userData, err := s.users.UserDataByExternalID(ctx, req.GetUserExternalId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, server.NotFoundError("user not found")
		}

		return nil, server.InternalError("error getting user data")
	}

	var dateOfBirthPb *timestamppb.Timestamp
	if userData.DateOfBirth != nil {
		dateOfBirthPb = timestamppb.New(*userData.DateOfBirth)
	}

	genderPb := lsuserspb.Gender_GENDER_UNSPECIFIED
	if userData.Gender != nil {
		genderPb = lsuserspb.Gender(*userData.Gender)
	}

	var avatarFileKeyPb *wrappers.StringValue
	if userData.AvatarFileKey != nil {
		avatarFileKeyPb = &wrappers.StringValue{Value: *userData.AvatarFileKey}
	}

	userDataResponse := &lsuserspb.UserDataResponse{
		Id:            userData.ID,
		FullName:      userData.FullName,
		DateOfBirth:   dateOfBirthPb,
		Gender:        genderPb,
		AvatarFileKey: avatarFileKeyPb,
	}
	return userDataResponse, nil
}

// GetFollowedUsers returns a list of users that the given user is followed to.
func (s *serverAPI) GetFollowedUsers(
	ctx context.Context,
	req *lsuserspb.GetFollowedUsersRequest,
) (*lsuserspb.FollowedUsersResponse, error) {
	if err := validateFollowedUsersRequest(req); err != nil {
		return nil, err
	}

	followedUsers, err := s.users.FollowedUsers(ctx, req.GetUserId())
	if err != nil {
		return nil, server.InternalError("failed to get followed users")
	}

	followedUsersPb := make([]*lsuserspb.FollowedUser, 0)
	for _, fu := range followedUsers {
		var avatarFileKeyPb *wrapperspb.StringValue
		if fu.AvatarFileKey != nil {
			avatarFileKeyPb = wrapperspb.String(*fu.AvatarFileKey)
		}

		followedUserPb := &lsuserspb.FollowedUser{
			FollowLinkId:  fu.FollowLinkID,
			NumberOfLikes: fu.NumberOfLikes,
			UserId:        fu.UserID,
			FullName:      fu.FullName,
			AvatarFileKey: avatarFileKeyPb,
		}
		followedUsersPb = append(followedUsersPb, followedUserPb)
	}

	followedUsersResponse := &lsuserspb.FollowedUsersResponse{
		Users: followedUsersPb,
	}

	return followedUsersResponse, nil
}

// FollowUser adds the user with userIdToFollow to the list of followed users with userId.
func (s *serverAPI) FollowUser(
	ctx context.Context,
	req *lsuserspb.FollowUserRequest,
) (*lsuserspb.FollowUserResponse, error) {
	if err := validateFollowUserRequest(req); err != nil {
		return &lsuserspb.FollowUserResponse{Success: false}, err
	}

	err := s.users.FollowUser(ctx, req.GetUserId(), req.GetUserIdToFollow())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, server.NotFoundError("user not found")
		}

		return &lsuserspb.FollowUserResponse{Success: false}, server.InternalError("error following user")
	}

	return &lsuserspb.FollowUserResponse{Success: true}, nil
}

// UnfollowUser removes a user from the follow list.
func (s *serverAPI) UnfollowUser(
	ctx context.Context,
	req *lsuserspb.UnfollowUserRequest,
) (*lsuserspb.UnfollowUserResponse, error) {
	if err := validateUnfollowUserRequest(req); err != nil {
		return &lsuserspb.UnfollowUserResponse{Success: false}, err
	}

	err := s.users.UnfollowUser(ctx, req.GetFollowLinkId())
	if err != nil {
		return &lsuserspb.UnfollowUserResponse{Success: false}, server.InternalError("error unfollowing user")
	}

	return &lsuserspb.UnfollowUserResponse{Success: true}, nil
}

func validateGetUserDataRequest(req *lsuserspb.GetUserDataRequest) error {
	if req.GetUserId() == emptyValue {
		return server.InvalidArgumentError("user id is empty")
	}

	return nil
}

func validateGetUserDataByExternalIdRequest(req *lsuserspb.GetUserDataByExternalIdRequest) error {
	if req.GetUserExternalId() == emptyValue {
		return server.InvalidArgumentError("user external id is empty")
	}

	return nil
}

func validateFollowedUsersRequest(req *lsuserspb.GetFollowedUsersRequest) error {
	if req.GetUserId() == emptyValue {
		return server.InvalidArgumentError("userid is empty")
	}

	return nil
}

func validateFollowUserRequest(req *lsuserspb.FollowUserRequest) error {
	if req.GetUserId() == emptyValue {
		return server.InvalidArgumentError("user id is empty")
	}

	if req.GetUserIdToFollow() == emptyValue {
		return server.InvalidArgumentError("user id to follow is empty")
	}

	return nil
}

func validateUnfollowUserRequest(req *lsuserspb.UnfollowUserRequest) error {
	if req.GetFollowLinkId() == emptyValue {
		return server.InvalidArgumentError("follow link id is empty")
	}

	return nil
}
