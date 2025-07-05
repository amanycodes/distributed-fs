package p2p

type Peer interface {
	Close() error
}

type Transport interface {
	ListenAndAccept() error
	Consumer() <-chan RPC
}
