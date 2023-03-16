package connection

type Connection interface {
	Discover() error
	Register() error
	Connect() error
	Ignore() error
	Disconnect() error
	Delete() error
}
