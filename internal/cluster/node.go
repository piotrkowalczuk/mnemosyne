package cluster

import "github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"

// Node ...
type Node struct {
	Addr   string
	Client mnemosynerpc.SessionManagerClient
}
