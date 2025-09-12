package grpc

import (
	"google.golang.org/grpc"
	"love-signal-users/internal/controller"
	v1 "love-signal-users/internal/controller/grpc/v1"
)

// NewRouter creates a new router for the gRPC server controller.
func NewRouter(
	server *grpc.Server,
	userDataUseCase controller.UserData,
	userDataByExternalIDUseCase controller.UserDataByExternalID,
) {
	v1.NewRoutes(
		server,
		userDataUseCase,
		userDataByExternalIDUseCase,
	)
}
