package meshes

import (
	context "context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MeshClient represents a gRPC adapter client
type MeshClient struct {
	MClient MeshServiceClient
	conn    *grpc.ClientConn
}

// CreateClient creates a MeshClient for the given params
func CreateClient(_ context.Context, meshLocationURL string) (*MeshClient, error) {
	var opts []grpc.DialOption
	// creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to create TLS credentials %v", err)
	// 	}
	// 	opts = append(opts, grpc.WithTransportCredentials(creds))
	// } else {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// }
	conn, err := grpc.Dial(meshLocationURL, opts...)
	if err != nil {
		return nil, err
	}

	mClient := NewMeshServiceClient(conn)

	return &MeshClient{
		conn:    conn,
		MClient: mClient,
	}, nil
}

// Close closes the MeshClient
func (m *MeshClient) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}
