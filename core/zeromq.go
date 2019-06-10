package core

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"time"
)


func StartZeroMQWorkers(num int) {

	for i := 0; i != num; i = i + 1 {
		go worker()
	}


	context, _ := zmq.NewContext()
	defer context.Close()

	// Socket to talk to clients
	clients, _ := context.NewSocket(zmq.ROUTER)
	defer clients.Close()
	clients.Bind("tcp://*:5555")

	// Socket to talk to workers
	workers, _ := context.NewSocket(zmq.DEALER)
	defer workers.Close()
	workers.Bind("ipc://workers.ipc")

	// connect work threads to client threads via a queue
	zmq.Device(zmq.QUEUE, clients, workers)

}


func worker() {
	context, _ := zmq.NewContext()
	defer context.Close()

	// Socket to talk to dispatcher
	receiver, _ := context.NewSocket(zmq.REP)
	defer receiver.Close()
	receiver.Connect("ipc://workers.ipc")

	for true {
		received, _ := receiver.Recv(0)
		fmt.Printf("Received request [%s]\n", received)

		// Do some 'work'
		time.Sleep(time.Second)

		// Send reply back to client
		receiver.Send([]byte("World"), 0)
	}
}

