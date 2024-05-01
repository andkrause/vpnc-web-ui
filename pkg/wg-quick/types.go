package wgquick

import (
	"regexp"
	"sync"
)

type WgQuick struct {
	wgQuickCommand                string
	configFolder                  string
	configFolderCheckRegexp       *regexp.Regexp
	wireguardNetworkInterfaceName string
	wgQuickConfigSearchDir        string
	mu                            sync.Mutex
}
