package main

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/opentelemetry/opentelemetry-demo/src/product-catalog/genproto/oteldemo"
)

const bufSize = 1024 * 1024

// startTestServer starts an in-memory gRPC server with the real productCatalog service registered
// and returns a connected client + cleanup function.
func startTestServer(t *testing.T) (pb.ProductCatalogServiceClient, grpc_health_v1.HealthClient, func()) {
	t.Helper()

	listener := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	svc := &productCatalog{} // uses global `catalog` loaded by init() from products/products.json

	// Register servers
	pb.RegisterProductCatalogServiceServer(s, svc)
	grpc_health_v1.RegisterHealthServer(s, svc)

	// Serve in background
	go func() {
		if err := s.Serve(listener); err != nil {
			t.Fatalf("gRPC Serve error: %v", err)
		}
	}()

	// Dialer for bufconn
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		cancel()
		t.Fatalf("gRPC Dial error: %v", err)
	}

	catalogClient := pb.NewProductCatalogServiceClient(conn)
	healthClient := grpc_health_v1.NewHealthClient(conn)

	cleanup := func() {
		conn.Close()
		s.GracefulStop()
		cancel()
	}

	return catalogClient, healthClient, cleanup
}

func TestHealthCheck(t *testing.T) {
	_, healthClient, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("health Check error: %v", err)
	}
	if got, want := resp.GetStatus(), grpc_health_v1.HealthCheckResponse_SERVING; got != want {
		t.Fatalf("health status = %v, want %v", got, want)
	}
}

func TestListProducts(t *testing.T) {
	client, _, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.ListProducts(ctx, &pb.Empty{})
	if err != nil {
		t.Fatalf("ListProducts error: %v", err)
	}

	// products/products.json in the repo has 10 items
	if got, wantMin := len(resp.GetProducts()), 10; got < wantMin {
		t.Fatalf("ListProducts returned %d products, want at least %d", got, wantMin)
	}
}

func TestGetProduct_Success_DefaultFlag(t *testing.T) {
	// NOTE: checkProductFailure only fails for id "OLJCESPC7Z" *if* feature flag is enabled.
	// In tests we don't configure OpenFeature provider → default is false → should succeed.
	client, _, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &pb.GetProductRequest{Id: "OLJCESPC7Z"}
	got, err := client.GetProduct(ctx, req)
	if err != nil {
		t.Fatalf("GetProduct error: %v", err)
	}
	if got.GetId() != "OLJCESPC7Z" {
		t.Fatalf("GetProduct id = %s, want %s", got.GetId(), "OLJCESPC7Z")
	}
	if got.GetName() == "" {
		t.Fatalf("GetProduct name should not be empty")
	}
}

func TestSearchProducts_Telescope(t *testing.T) {
	client, _, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &pb.SearchProductsRequest{Query: "telescope"}
	resp, err := client.SearchProducts(ctx, req)
	if err != nil {
		t.Fatalf("SearchProducts error: %v", err)
	}

	// Repo data contains multiple items with "Telescope" in the name.
	if len(resp.GetResults()) < 2 {
		t.Fatalf("SearchProducts results = %d, want >= 2", len(resp.GetResults()))
	}
}
