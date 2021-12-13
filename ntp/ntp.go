package ntp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
	modeClient     = 3
	modeServer     = 4
	leapUnsync     = 3
)

func GetTime(host string) (time.Time, error) {
	msg, err := queryServer(host)
	if err != nil {
		fmt.Println(" is broken")
		return ntpEpoch, err
	}
	fmt.Println(" is a stratum", msg.Stratum, "NTP server")
	return parseTime(msg.ReceiveTime), nil
}

func queryServer(host string) (*msg, error) {
	addr := net.JoinHostPort(host, "123")
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	defer con.Close()
	con.SetDeadline(time.Now().Add(defaultTimeout))

	recvMsg := new(msg)
	sendMsg := new(msg)
	sendMsg.setLeap(leapUnsync)
	sendMsg.setVersion(4)
	sendMsg.setMode(modeClient)

	err = binary.Write(con, binary.BigEndian, sendMsg)
	if err != nil {
		return nil, err
	}

	err = binary.Read(con, binary.BigEndian, recvMsg)
	if err != nil {
		return nil, err
	}

	if recvMsg.getVersion() != 4 {
		return nil, errors.New("invalid version in response")
	}
	if recvMsg.getMode() != modeServer {
		return nil, errors.New("invalid mode in response")
	}
	if recvMsg.TransmitTime == 0 {
		return nil, errors.New("invalid transmit time in response")
	}
	if recvMsg.OriginTime != sendMsg.TransmitTime {
		return nil, errors.New("server response mismatch")
	}
	if recvMsg.ReceiveTime > recvMsg.TransmitTime {
		return nil, errors.New("server clock ticked backwards")
	}

	return recvMsg, nil
}
