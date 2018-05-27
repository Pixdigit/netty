package netty

import (
	"github.com/pkg/errors"
	"strconv"
)

type NetConf struct {
	Protocol   string
	RemoteAddr string
	RemotePort uint16
	LocalAddr  string
	LocalPort  uint16
}

func (nc *NetConf) Valid() (bool, error) {
	for _, protocol := range protocols {
		if protocol == nc.Protocol {
			return true, nil
		}
	}
	return false, errors.New("unknown protocol \"" + nc.Protocol + "\"")
}

func (nc *NetConf) FullLocalAddress() (string, error) {
	return nc.LocalAddr + ":" + strconv.Itoa(int(nc.LocalPort)), nil
}
func (nc *NetConf) FullRemoteAddress() (string, error) {
	return nc.RemoteAddr + ":" + strconv.Itoa(int(nc.RemotePort)), nil
}
