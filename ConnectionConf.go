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

func (nc *ConnectionConf) Valid() (bool, error) {
	for _, protocol := range protocols {
		if protocol == nc.Protocol {
			return true, nil
		}
	}
	return false, errors.New("unknown protocol \"" + nc.Protocol + "\"")
}

func (nc *ConnectionConf) FullLocalAddress() (string, error) {
	return nc.LocalAddr + ":" + strconv.Itoa(int(nc.LocalPort)), nil
}
func (nc *ConnectionConf) FullRemoteAddress() (string, error) {
	return nc.RemoteAddr + ":" + strconv.Itoa(int(nc.RemotePort)), nil
}
