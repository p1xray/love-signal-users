package users

import (
	userspb "github.com/p1xray/love-signal-protos/gen/go/users"
	"google.golang.org/grpc"
)

type serverAPI struct {
	userspb.UnimplementedUsersServer
}

func Register(gRPC *grpc.Server) {
	userspb.RegisterUsersServer(gRPC, &serverAPI{})
}
