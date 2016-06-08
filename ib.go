package ib

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Portfolio struct {
}

func (p *Portfolio) UnMarshall(d string) {

}

// IB is the main communication and initialization struct to receive and send communication to interactive brokers api.
type IB struct {
	PortFolio    chan *Portfolio
	CreateBackup bool
	BackupFile   string
	Conn         net.Conn
	Quit         chan bool
	ClientId     int64
	Rid          int64
	OutStream    *bytes.Buffer
}

// processPortfolioMsg messages received which contain portfolio information.
func (b *IB) processPortfolioMsg(d string) {
	s := &Portfolio{}
	s.UnMarshall(d)

	b.PortFolio <- s
}

// ProcessReceiver is one of the main reciever functions that interprets data received by IQFeed and processes it in sub functions.
func (b *IB) processReceiver(d string) {
	/*data := d[2:]
	switch d[0] {

	case 0x50: // Start letter is P, indicating a summary message.
	*/
	data := d
	b.processPortfolioMsg(data)
	//}

}

// Read function does as expected and reads data from the network stream.
func (b *IB) read() {
	r := bufio.NewReader(b.Conn)
	for {
		select {
		case <-b.Quit:
			log.Println("Client quitting")
			b.Conn.Close()
			os.Exit(0)
			break
		default:
			str, err := r.ReadString(DELIM_BYTE)
			if err == nil {
				d := strings.TrimRight(str, DELIM_STR)
				fmt.Printf(d)
				b.processReceiver(d)
			}
			if err != nil {
				log.Printf("Could not read: %s", err)
				b.Quit <- true
			}

			//b.processReceiver(d)
		}
	}

}

// Start function will start the concurrent functions to read and write data to the and from the network stream.
func (b *IB) Start(connectString string) *IB {
	b.OutStream = bytes.NewBuffer(make([]byte, 0, 4096))
	err := b.connect(connectString)
	if err != nil {
		log.Printf("Could not connect to IB: %s", err)
	}
	b.PortFolio = make(chan *Portfolio)

	return b
}
