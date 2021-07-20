package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/cloudflare/roughtime/mjd"
	"github.com/cloudflare/roughtime/protocol"

	"golang.org/x/crypto/ed25519"
)

var (
	addr        string
	privKeyFile string
)

func init() {
	flag.StringVar(&addr, "a", ":2002", "address to listen on")
	flag.StringVar(&privKeyFile, "k", "priv.key", "file with private key")
	flag.Parse()
}

func main() {
	keyStr, err := os.ReadFile(privKeyFile)
	if err != nil {
		fmt.Printf("fatal: could not read %s: %v\n", privKeyFile, err)
		os.Exit(1)
	}

	rootKey, err := base64.StdEncoding.DecodeString(string(keyStr))
	if err != nil {
		fmt.Printf("fatal: could not decode private key: %v\n", err)
		os.Exit(1)
	}

	netAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("fatal: could not resolve %s: %v\n", netAddr, err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", netAddr)
	if err != nil {
		fmt.Printf("fatal: could not listen on %s: %v\n", addr, err)
		os.Exit(1)
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("fatal: could not generate key: %v\n", err)
		os.Exit(1)
	}

	now := mjd.Now()
	yesterday := mjd.New(now.Day()-1, now.Microseconds())
	tomorrow := mjd.New(now.Day()+1, now.Microseconds())
	cert, err := protocol.CreateCertificate(yesterday, tomorrow, pub, rootKey)
	if err != nil {
		fmt.Printf("fatal: could not generate certificate: %v\n", err)
		os.Exit(1)
	}

	query := make([]byte, 1280)
	for {
		queryLen, peer, err := conn.ReadFrom(query)
		if err != nil {
			fmt.Printf("error reading query: %v\n", err)
			continue
		}
		reply, err := protocol.CreateReply(query[:queryLen], getMjd(), 1000000, cert, priv)
		if err != nil {
			fmt.Printf("error generating reply: %v\n", err)
			continue
		}
		fmt.Printf("sending reply to %s\n", peer)
		conn.WriteTo(reply, peer)
	}
}
