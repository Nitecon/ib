package ib

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func NextClientId() int64 {
	t := time.Now().UnixNano()
	rand.Seed(t)
	CLIENT_ID_INCR = int64(rand.Intn(9999))

	return CLIENT_ID_INCR
}

func (b *IB) NextReqId() int64 {
	b.Rid += 1
	return b.Rid
}

func (b *IB) SendRequest() (int, error) {
	b.WriteString(DELIM_STR)
	output := b.OutStream.Bytes()
	log.Printf("Writing: %s\n", string(output))

	i, err := b.Conn.Write(output)

	b.OutStream.Reset()

	return i, err
}

func (b *IB) connect(cs string) error {
	if cs == "" {
		cs = "localhost:7497"
	}
	conn, err := net.Dial("tcp", cs)
	if err != nil {
		log.Fatal("Could not connect to Interactive Brokers API")
	}
	log.Printf("Connected... (%s)", cs)
	b.Conn = conn
	b.ClientId = NextClientId()
	go b.read()
	/*b.Conn.Write([]byte(fmt.Sprintf("%d%s", 63, DELIM_BYTE)))
	b.Conn.Write([]byte(fmt.Sprintf("%d%s", NextClientId(), DELIM_BYTE)))*/
	log.Println("Starting handshake...")
	err = b.ServerShake(63)

	if err != nil {
		log.Printf("Handshake failed : %s", err)
	}
	return err
	// FOrmatting ints : strconv.FormatInt(i, 10)

}
func (b *IB) ServerShake(version int64) error {
	b.WriteInt(version)
	b.WriteInt(b.ClientId)

	_, err := b.SendRequest()

	return err
}

func (b *IB) WriteString(s string) (int, error) {
	return b.OutStream.WriteString(s + DELIM_STR)
}

func (b *IB) WriteInt(i int64) (int, error) {
	return b.OutStream.WriteString(strconv.FormatInt(i, 10) + DELIM_STR)
}

func (b *IB) WriteFloat(f float64) (int, error) {
	return b.OutStream.WriteString(strconv.FormatFloat(f, 'g', 10, 64) + DELIM_STR)
}

func (b *IB) WriteBool(boo bool) (int, error) {
	if boo {
		return b.OutStream.WriteString("1" + DELIM_STR)
	}

	return b.OutStream.WriteString("0" + DELIM_STR)
}
