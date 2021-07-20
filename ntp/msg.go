package ntp

type msg struct {
	LeapVerMode    uint8
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	ReferenceTime  uint64
	OriginTime     uint64
	ReceiveTime    uint64
	TransmitTime   uint64
}

func (m *msg) getLeap() uint8 {
	return (m.LeapVerMode >> 6) & 0x03
}

func (m *msg) getVersion() int {
	return int((m.LeapVerMode >> 3) & 0x07)
}

func (m *msg) getMode() uint8 {
	return m.LeapVerMode & 0x07
}

func (m *msg) setLeap(leap uint8) {
	m.LeapVerMode = (m.LeapVerMode & 0x3f) | leap<<6
}

func (m *msg) setVersion(ver uint8) {
	m.LeapVerMode = (m.LeapVerMode & 0xc7) | ver<<3
}

func (m *msg) setMode(mode uint8) {
	m.LeapVerMode = (m.LeapVerMode & 0xf8) | mode
}
