package vpnc

import (
	"sync"
)

type VPNC struct {
	connectCommand           string
	disconnectCommand        string
	configFolder             string
	vpncNetworkInterfaceName string
	mu                       sync.Mutex
}
