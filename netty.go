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
	rw           *bufio.ReadWriter
}

func NewConnector(nc NetConf) (Connector, error) {
	var err error
	newConn := Connector{}
	if ok, err := nc.Valid(); !ok {
		return Connector{}, errors.Wrap(err, "invalid network configuration")
	}
	newConn.netConf = nc
	newConn.state = STOPPED
	newConn.incomingData = make(chan map[string]string)
	addr, err := newConn.netConf.FullLocalAddress()
	if err != nil {
		return Connector{}, errors.Wrap(err, "could not determine local address")
	}
	listener, err := net.Listen(newConn.netConf.Protocol, addr)
	if err != nil {
		addr, err2 := nc.FullLocalAddress()
		if err2 != nil {
			return Connector{}, errors.Wrap(err2, "unable to determine local address")
		}
		return Connector{}, errors.Wrap(err, "unable to listen on "+addr)
	}
	newConn.listener = listener

	remoteAddr, _ := newConn.netConf.FullRemoteAddress()
	clientConn, err := net.Dial(newConn.netConf.Protocol, remoteAddr)
	if err != nil {
		addr, err := nc.FullRemoteAddress()
		if err != nil {
			return Connector{}, errors.Wrap(err, "unable to determine remote address")
		}
		return Connector{}, errors.Wrap(err, "unable to connect to "+addr)
	}

	newConn.rw = bufio.NewReadWriter(bufio.NewReader(clientConn), bufio.NewWriter(clientConn))

	return newConn, nil
}
