package netty

import (
	"bufio"
	"net"

	"github.com/pkg/errors"
)

type Connection struct {
	conf         ConnectionConf
	listener     net.Listener
	incomingData chan map[string]string
	state        Runlevel
	isConnected  bool
	rw           *bufio.ReadWriter
}

func NewConnection(nc ConnectionConf) (Connection, error) {
	var err error
	newConn := Connection{}

	ok, err := nc.Valid()
	if !ok || err != nil {
		return Connection{}, errors.Wrap(err, "invalid network configuration")
	}

	newConn.conf = nc
	newConn.state = STOPPED
	newConn.incomingData = make(chan map[string]string)
	newConn.isConnected = false

	addr, err := newConn.conf.FullLocalAddress();	if err != nil {return Connection{}, errors.Wrap(err, "could not determine local address")}
	listener, err := net.Listen(newConn.conf.Protocol, addr);	if err != nil {return Connection{}, errors.Wrap(err, "unable to listen on "+addr)}
	newConn.listener = listener

	return newConn, nil
}

func (c *Connection) Connect() error {
	remoteAddr, err := c.conf.FullRemoteAddress();	if err != nil {return errors.Wrap(err, "could not read remote addr")}
	connection, err := net.Dial(c.conf.Protocol, remoteAddr);	if err != nil {return errors.Wrap(err, "unable to connect to "+remoteAddr)}
	c.rw = bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection))
	c.isConnected = true
	return nil
}

/*
TODO:
I only need one connection!
User listener!
func (c *Connection) TestWrite() error {
	n, err := c.rw.Write([]byte("test;"))
	c.rw.Flush();	if err != nil {return err}
	return nil
}
func (c *Connection) TestRead() (rune, error) {
	thisRune, _, err := c.rw.ReadRune();	if err != nil {return thisRune, err}
	return thisRune, nil
}*/
