package netty

import (
	"bufio"
	"net"

	"github.com/pkg/errors"
)

type Connector struct {
	netConf      NetConf
	listener     net.Listener
	incomingData chan map[string]string
	state        Runlevel
	isConnected  bool
	selfInitRw           *bufio.ReadWriter
	remoteInitRw           *bufio.ReadWriter
}

func NewConnector(nc NetConf) (Connector, error) {
	var err error
	newConn := Connector{}
	if ok, err := nc.Valid(); !ok {
		return Connector{}, errors.Wrap(err, "invalid metwork configuration")
	}
	newConn.netConf = nc
	newConn.state = STOPPED
	newConn.incomingData = make(chan map[string]string)
	addr, err := newConn.netConf.FullLocalAddress()
	if err != nil {
		return Connector{}, errors.Wrap(err, "could not determine local address")
	}
	newConn.listener, err = net.Listen(newConn.netConf.Protocol, addr)
	if err != nil {
		return Connector{}, errors.Wrap(err, "unable to listen on "+newConn.listener.Addr().String())
	}

	clientConn, err := net.Dial(newConn.NetConf.protocol, newConn.netConf.FullRemoteAddress())
	if err != nil {
		return Connector{}, errors.Wrap(err, "could not connect to remote address")
	}
	newConn.selfInitRw = bufio.NewReadWriter(bufio.NewReader(clientConn), bufio.NewWriter(clientConn))


	return newConn, nil
}
