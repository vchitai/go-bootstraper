package services

import (
	"context"

	"domain_dash/service_dash/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type serverNameCamelInterface interface {
	/* range .Protos */
	pb.ProtoNameTitleServer
	/* end */
}

type serverNameCamelServer struct {
	serverNameCamelInterface
}

func newServerNameTitleServer(serverImpl serverNameCamelInterface) ServiceServer {
	return &serverNameCamelServer{
		serverNameCamelInterface: serverImpl,
	}
}

func (s *serverNameCamelServer) RegisterWithServer(server *grpc.Server) {
	/* range .Protos */
	pb.RegisterProtoNameTitleServer(server, s)
	/* end */
}

func (s *serverNameCamelServer) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	/* range .Protos */
	if err := pb.RegisterProtoNameTitleHandler(ctx, mux, conn); err != nil {
		return err
	}
	/* end */
	return nil
}

// Close ...
func (s *serverNameCamelServer) Close(ctx context.Context) {
}
