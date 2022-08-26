package channels

var (
	MeshSync = "meshsync"
)

func NewMeshSyncChannel() MeshSyncChannel {
	return make(chan struct{})
}

type MeshSyncChannel chan struct{}

func (ch MeshSyncChannel) Stop() {
	<-ch
}
