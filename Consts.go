package netty

type Runlevel int32

const STOPPED Runlevel = 0
const RUNNING Runlevel = 1
const PAUSED Runlevel = 2

const TCP string = "tcp"
const TCP4 string = "tcp4"
const TCP6 string = "tcp6"
const UDP string = "udp"
const UDP4 string = "udp4"
const UDP6 string = "udp6"

var protocols = [...]string{TCP, TCP4, TCP6, UDP, UDP4, UDP6}

//End of key indicator
const EOK_RUNE = rune('/')

//End of value (and transmission) indicator
const EOT_RUNE = rune('!')

var LOCALNETCONF ConnectionConf

func init() {
	LOCALNETCONF = ConnectionConf{TCP, "localhost", 49125, "localhost", 49125}
}
