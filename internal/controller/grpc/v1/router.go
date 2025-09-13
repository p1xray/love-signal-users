package v1

import (
	"google.golang.org/grpc"
	"love-signal-users/internal/controller"
	"love-signal-users/internal/controller/grpc/v1/users"
)

// NewRoutes creates a new routes for the gRPC server controller of version 1.
func NewRoutes(
	server *grpc.Server,
	userDataUseCase controller.UserData,
	userDataByExternalIDUseCase controller.UserDataByExternalID,
	followedUsersUseCase controller.Followed,
	followUserUseCase controller.Follow,
	unfollowUserUseCase controller.Unfollow,
) {
	users.RegisterUsersServer(
		server,
		userDataUseCase,
		userDataByExternalIDUseCase,
		followedUsersUseCase,
		followUserUseCase,
		unfollowUserUseCase,
	)
}
