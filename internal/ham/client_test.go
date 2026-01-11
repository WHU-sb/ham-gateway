package ham

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestAddAuth(t *testing.T) {
	client := &Client{
		appID:     "test-app-id",
		appSecret: "test-app-secret",
	}

	ctx := context.Background()
	authCtx := client.addAuth(ctx)

	md, ok := metadata.FromOutgoingContext(authCtx)
	assert.True(t, ok, "metadata should be present")
	assert.Equal(t, []string{"test-app-id"}, md.Get("x-app-id"))
	assert.Equal(t, []string{"test-app-secret"}, md.Get("x-app-secret"))
}

func TestNewClient_InvalidCertificate(t *testing.T) {
	client, err := NewClient(
		"localhost:4443",
		"app-id",
		"app-secret",
		"invalid-cert.crt",
		"invalid-key.key",
	)

	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to load x509 key pair")
}

// Note: Full integration tests with real gRPC server would require
// setting up a test gRPC server, which is beyond unit testing scope.
// These tests verify the client initialization and auth context setup.
