package meshes

import (
	context "context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AdapterClient represents a gRPC adapter client
type AdapterClient struct {
	AClient AdapterServiceClient
	conn    *grpc.ClientConn
}

// CreateClient creates a AdapterClient for the given params
func CreateClient(_ context.Context, meshLocationURL string) (*AdapterClient, error) {
	var opts []grpc.DialOption
	// creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to create TLS credentials %v", err)
	// 	}
	// 	opts = append(opts, grpc.WithTransportCredentials(creds))
	// } else {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// }
	conn, err := grpc.NewClient(meshLocationURL, opts...)
	if err != nil {
		return nil, err
	}

	aClient := NewAdapterServiceClient(conn)

	return &AdapterClient{
		conn:    conn,
		AClient: aClient,
	}, nil
}

// Close closes the AdapterClient
func (a *AdapterClient) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
