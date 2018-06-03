package netty

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"

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
	newConn.state = PAUSED
	newConn.incomingData = make(chan map[string]string)
	newConn.isConnected = false

	addr, err := newConn.conf.FullLocalAddress();	if err != nil {return Connection{}, errors.Wrap(err, "could not determine local address")}
	listener, err := net.Listen(newConn.conf.Protocol, addr);	if err != nil {return Connection{}, errors.Wrap(err, "unable to listen on "+addr)}
	newConn.listener = listener

	return newConn, nil
}

func (c *Connection) Connect() error {
	if c.isConnected {
		return errors.New("connection already established")
	}
	remoteAddr, err := c.conf.FullRemoteAddress();	if err != nil {return errors.Wrap(err, "could not read remote addr")}
	connection, err := net.Dial(c.conf.Protocol, remoteAddr);	if err != nil {return errors.Wrap(err, "unable to connect to "+remoteAddr)}
	c.rw = bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection))
	c.isConnected = true
	return nil
}

func (c *Connection) Accept(errChan chan error) {
	go func() {
		for c.state != STOPPED {
			//TODO: Check if connection comes from specified remote
			conn, err := c.listener.Accept()
			if err != nil {
				pushErr(errChan, errors.Wrap(err, "failed accepting connection request"))
			} else if !c.isConnected {
				c.isConnected = true
				c.rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
				c.handleConnection(errChan)
				conn.Close()
			} else {
				pushErr(errChan, errors.New("connection already established"))
			}
		}
	}()
}

func (c *Connection) Start() {
	c.state = RUNNING
}
func (c *Connection) Pause() {
	c.state = PAUSED
}
func (c *Connection) Stop() {
	c.state = STOPPED
}

func (c *Connection) Send(key string, data interface{}) error {
	if !c.isConnected {
		return errors.New("connection not yet established")
	}

	dataString, err := serialize(data);	if err != nil {return errors.Wrap(err, "unable to send data")}
	sanitizedDataStr, err := sanitize(dataString, EOT_RUNE);	if err != nil {return errors.Wrap(err, "could not sanitize data")}
	sanitizedKey, err := sanitize(key, EOK_RUNE);	if err != nil {return errors.Wrap(err, "could not sanitize key")}
	//Append escape sequence
	sanitizedStr := sanitizedKey + string(EOK_RUNE) + sanitizedDataStr + string(EOT_RUNE) + " "
	_, err = c.rw.Write([]byte(sanitizedStr));	if err != nil {return errors.Wrap(err, "unable to send data")}
	err = c.rw.Flush();	if err != nil {return errors.Wrap(err, "unable to send data")}

	return nil
}

func (c *Connection) handleConnection(errChan chan error) {
	if !c.isConnected {
		pushErr(errChan, errors.New("connection not yet established"))
	}
	for c.state != STOPPED {
		for c.state == RUNNING {
			//one "character" (or more if waiting for another escape Rune)
			var token []rune
			key := ""
			strData := ""
			dataIsKey := true
			for {
				thisRune, _, err := c.rw.ReadRune()
				token = append(token, thisRune)
				if err != nil {
					//TODO: send notification of faulty msg to client
					pushErr(errChan, err)
					strData = ""
					token = []rune{}
				}
				//Token is of max length 2
				if dataIsKey {
					token, strData, err = readTokenWithEscapeRune(token, strData, EOK_RUNE)
					if err != nil {
						if err == io.EOF {
							key = strData[:len(strData)-1]
							//reset strData buffer to first rune of value strData
							strData = string(strData[len(strData)-1])
							dataIsKey = false
						} else {
							pushErr(errChan, err)
						}
					}
				} else {
					token, strData, err = readTokenWithEscapeRune(token, strData, EOT_RUNE)
					if err != nil {
						if err == io.EOF {
							dataMap := make(map[string]string)
							dataMap[key] = strData[:len(strData)-1]
							fmt.Println(dataMap)
							go func() { c.incomingData <- dataMap }()
							strData = ""
							key = ""
						} else {
								pushErr(errChan, err)
						}
						dataIsKey = true
					}
				}
			}
		}
	}
}

func sanitize(str string, escapeRunes ...rune) (string, error) {
	sanitizedStr := ""
	for _, char := range str {
		for _, escapeRune := range escapeRunes {
			if char == escapeRune {
				sanitizedStr += string(escapeRune)
			}
		}
		sanitizedStr += string(char)
	}
	return sanitizedStr, nil
}

func serialize(variable interface{}) (string, error) {
	var result string
	var err error

	switch t := variable.(type) {
	case bool:
		if t {
			result = "bool:true"
			} else {
				result = "bool:false"
			}
		case int32:
			result = "int:" + strconv.Itoa(int(t))
		case int:
			result = "int:" + strconv.Itoa(t)
		case float32:
			result = "float:" + strconv.FormatFloat(float64(t), 'G', -1, 64)
		case float64:
			result = "float:" + strconv.FormatFloat(t, 'G', -1, 64)
		case string:
			result = "str:" + t
		default:
			err = errors.New(fmt.Sprintf("Can not serialize %#v of type %T", t, t))
	}
	return result, err
}


func readTokenWithEscapeRune(token []rune, data string, escapeRune rune) ([]rune, string, error) {
	var err error = nil
	if len(token) == 1 && token[0] != escapeRune {
		//Single rune
		data += string(token[0])
		token = []rune{}

	} else if len(token) == 2 && token[0] == escapeRune {
		if token[0] == escapeRune && token[1] == escapeRune {
		//Escaped escape rune
		token = []rune{escapeRune}
		} else {
			//Recieved single escape rune as end statement
			err = io.EOF
			token = []rune{token[1]}
		}
		data += string(token[0])
		token = []rune{}

	} else {
		//Token not correctly formatted
		if len(token) > 2 || len(token) == 0 {
			err = errors.New("token of unusable size")
		} else if len(token) == 2 && token[0] != escapeRune {
			err = errors.New("token longer than 1 but does not begin with escape rune")
		}
	}

	return token, data, err
}
