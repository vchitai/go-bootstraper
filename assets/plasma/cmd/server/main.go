package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"

	"domain_dash/service_dash/configs"
	"domain_dash/service_dash/internal/services"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func initialize(cfg *configs.Config) error {
	// do any init
	return nil
}

func createServer(cfg *configs.Config) []services.ServiceServer {
	return []services.ServiceServer{
		// TODO: Your service server init goes here
		// example: services.NewServerNameTitle(cfg),
	}
}

func createGateway(cfg *configs.Config, mux *runtime.ServeMux) *http.Server {
	muxHttp := http.NewServeMux()
	muxHttp.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	muxHttp.HandleFunc("/debug/pprof/", pprof.Index)
	muxHttp.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	muxHttp.HandleFunc("/debug/pprof/profile", pprof.Profile)
	muxHttp.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	muxHttp.HandleFunc("/debug/pprof/trace", pprof.Trace)
	muxHttp.Handle("/", mux)

	return &http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTPAddress), Handler: muxHttp}
}

func run(args []string) error {
	var cfg = mustLoadConfig()
	ll, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("cannot init logger %w", err)
	}

	if err := initialize(cfg); err != nil {
		return fmt.Errorf("initialize %w", err)
	}

	var mux = runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	var ctx = context.TODO()

	grpc_zap.ReplaceGrpcLoggerV2(ll)

	var grpcServer = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(ll),
			grpc_validator.UnaryServerInterceptor()),
		grpc_middleware.WithStreamServerChain(grpc_prometheus.StreamServerInterceptor,
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(ll),
			grpc_validator.StreamServerInterceptor(),
		),
		grpc.MaxConcurrentStreams(cfg.MaxConcurrentStreams),
	)

	grpcListen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCAddress))
	if err != nil {
		return fmt.Errorf("cannot create grpc listerner %w", err)
	}
	var master masterRoutine
	defer master.shutdown()

	master.fork(
		func() {
			if err := grpcServer.Serve(grpcListen); err != nil {
				ll.Error("Serving gRPC", zap.Error(err))
			}
		},
		func(ctx context.Context) {
			grpcServer.GracefulStop()
		},
	)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%d", cfg.GRPCAddress), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("cannot dial grpc %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			ll.Error("Cannot close gRPC Dial", zap.Error(err))
		}
	}()

	for _, svc := range createServer(cfg) {
		svc.RegisterWithServer(grpcServer)
		if err := svc.RegisterWithHandler(ctx, mux, conn); err != nil {
			return err
		}
	}

	var srv = createGateway(cfg, mux)

	master.fork(func() {
		if err := srv.ListenAndServe(); err != nil {
			ll.Error("Serving", zap.Error(err))
		}
	}, func(ctx context.Context) {
		if err := srv.Shutdown(ctx); err != nil {
			ll.Error("Shutdown gateway server", zap.Error(err))
		}
	})

	master.shutdownWait()

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
