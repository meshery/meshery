package meshes

import (
	context "context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AdapterClient represents a gRPC adapter client
type AdapterClient struct {
	MClient AdapterServiceClient
	conn    *grpc.ClientConn
}

// CreateClient creates a AdapterClient for the given params
func CreateClient(_ context.Context, adapterLocationURL string) (*AdapterClient, error) {
	var opts []grpc.DialOption
	// creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to create TLS credentials %v", err)
	// 	}
	// 	opts = append(opts, grpc.WithTransportCredentials(creds))
	// } else {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// }
	conn, err := grpc.NewClient(adapterLocationURL, opts...)
	if err != nil {
		return nil, err
	}

	mClient := NewAdapterServiceClient(conn)

	return &AdapterClient{
		conn:    conn,
		MClient: mClient,
	}, nil
}

// Close closes the AdapterClient
func (m *AdapterClient) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}
