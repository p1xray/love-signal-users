package usersserver

import (
	"context"
	"love-signal-users/internal/dto"

	userspb "github.com/p1xray/love-signal-protos/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	// FollowedUsers returns a list of users that the given user is followed to.
	FollowedUsers(
		ctx context.Context,
		userId int64,
	) ([]*dto.FollowedUser, error)

	// FollowUser adds the user with userIdToFollow to the list of followed users with userId.
	FollowUser(
		ctx context.Context,
		userId int64,
		userIdToFollow int64,
	) error

	// UnfollowUser removes a user from the follow list.
	UnfollowUser(
		ctx context.Context,
		followLinkId int64,
	) error
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

	var dateOfBirthPb *timestamppb.Timestamp
	if card.DateOfBirth != nil {
		dateOfBirthPb = timestamppb.New(*card.DateOfBirth)
	}

	genderPb := userspb.Gender(0)
	if card.Gender != nil {
		genderPb = userspb.Gender(*card.Gender)
	}

	cardResponse := &userspb.UserProfileCardResponse{
		Id:          card.Id,
		Name:        card.Name,
		DateOfBirth: dateOfBirthPb,
		Gender:      genderPb,
	}

	return cardResponse, nil
}

// FollowedUsers returns a list of users that the given user is followed to.
func (s *serverAPI) FollowedUsers(
	ctx context.Context,
	req *userspb.FollowedUsersRequest,
) (*userspb.FollowedUsersResponse, error) {
	if err := validateFollowedUsersRequest(req); err != nil {
		return nil, err
	}

	followedUsers, err := s.users.FollowedUsers(ctx, req.GetUserId())
	if err != nil {
		return nil, internalError("failed to get followed users")
	}

	followedUsersPb := make([]*userspb.FollowedUser, 0)
	for _, fu := range followedUsers {
		var avatarFileKeyPb *wrapperspb.StringValue
		if fu.AvatarFileKey == nil {
			avatarFileKeyPb = wrapperspb.String(*fu.AvatarFileKey)
		}

		followedUserPb := &userspb.FollowedUser{
			FollowLinkId:     fu.FollowLinkId,
			SendedLikesCount: fu.SendedLikesCount,
			UserId:           fu.UserId,
			Name:             fu.Name,
			AvatarFileKey:    avatarFileKeyPb,
		}
		followedUsersPb = append(followedUsersPb, followedUserPb)
	}

	followedUsersResponse := &userspb.FollowedUsersResponse{
		Users: followedUsersPb,
	}

	return followedUsersResponse, nil
}

// FollowUser adds the user with userIdToFollow to the list of followed users with userId.
func (s *serverAPI) FollowUser(
	ctx context.Context,
	req *userspb.FollowUserRequest,
) (*userspb.FollowUserResponse, error) {
	if err := validateFollowUserRequest(req); err != nil {
		return &userspb.FollowUserResponse{Success: false}, err
	}

	err := s.users.FollowUser(ctx, req.GetUserId(), req.GetUserIdToFollow())
	if err != nil {
		return &userspb.FollowUserResponse{Success: false}, internalError("failed to follow user")
	}

	return &userspb.FollowUserResponse{Success: true}, nil
}

// UnfollowUser removes a user from the follow list.
func (s *serverAPI) UnfollowUser(
	ctx context.Context,
	req *userspb.UnfollowUserRequest,
) (*userspb.UnfollowUserResponse, error) {
	if err := validateUnfollowUserRequest(req); err != nil {
		return &userspb.UnfollowUserResponse{Success: false}, err
	}

	err := s.users.UnfollowUser(ctx, req.GetFollowLinkId())
	if err != nil {
		return &userspb.UnfollowUserResponse{Success: false}, internalError("failed to follow user")
	}

	return &userspb.UnfollowUserResponse{Success: true}, nil
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

func validateFollowedUsersRequest(req *userspb.FollowedUsersRequest) error {
	if req.GetUserId() == emptyValue {
		return invalidArgumentError("user id is empty")
	}

	return nil
}

func validateFollowUserRequest(req *userspb.FollowUserRequest) error {
	if req.GetUserId() == emptyValue {
		return invalidArgumentError("user id is empty")
	}

	if req.GetUserIdToFollow() == emptyValue {
		return invalidArgumentError("user id to follow is empty")
	}

	return nil
}

func validateUnfollowUserRequest(req *userspb.UnfollowUserRequest) error {
	if req.GetFollowLinkId() == emptyValue {
		return invalidArgumentError("follow link id is empty")
	}

	return nil
}

func invalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

func internalError(msg string) error {
	return status.Error(codes.Internal, msg)
}
