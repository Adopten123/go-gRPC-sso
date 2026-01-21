package suite

import (
	"context"
	ssov1 "github.com/Adopten123/go-protobufcontract-sso/gen/go/sso"
	"go-gRPC-sso/internal/config"
	"testing"
)

type Suite struct {
	*testing.T
	Config     *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	co
}
