package netty

import (
	"strconv"

	"github.com/pkg/errors"
)

type ConnectionConf struct {
	Protocol   string
	LocalAddr  string
	LocalPort  uint16
	RemoteAddr string
	RemotePort uint16
}

func (cc *ConnectionConf) Valid() (bool, error) {
	for _, protocol := range protocols {
		if protocol == cc.Protocol {
			return true, nil
		}
	}
	return false, errors.New("unknown protocol \"" + cc.Protocol + "\"")
}

func (cc *ConnectionConf) FullLocalAddress() (string, error) {
	return cc.LocalAddr + ":" + strconv.Itoa(int(cc.LocalPort)), nil
}
func (cc *ConnectionConf) FullRemoteAddress() (string, error) {
	return cc.RemoteAddr + ":" + strconv.Itoa(int(cc.RemotePort)), nil
}
