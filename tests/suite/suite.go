package suite

import (
	"context"
	"net"
	"strconv"

	ssov1 "github.com/Adopten123/go-protobufcontract-sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-gRPC-sso/internal/config"
	"testing"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Config     *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	config := config.MustLoadByPath("../config/local_tests.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), config.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(config),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Config:     config,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(config *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(config.GRPC.Port))
}
