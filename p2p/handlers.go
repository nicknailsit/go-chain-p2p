package p2p

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"os"
	"sync"
)


type TypeHandler struct {
	sync.Mutex
	StreamType string
	DataWriter bufio.Writer
	DataReader bufio.Reader
	endpoint string
}

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}


func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err !=nil {
			fmt.Println("error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}

		if str != "\n" {

			//todo read logic for blockchain here


			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}

}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func handleAuthStream(stream network.Stream) {

}


func handleSyncStream(stream network.Stream) {
	fmt.Println("Got a new blockchain sync request!")
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}