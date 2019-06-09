package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"io"
	"swaggp2p/chain"
	"sync"
	"time"
)

const ID = "/chain/0.0.1"

var log = logging.Logger("swaggchain")

const ChunkSize = 1024


const TTL = 600; //how long before disconnect from inactive service 10 min



type ChainService struct {
	sync.Mutex
	Host host.Host
	Version uint8
	Chain chain.Swagg
	CurrentTimeUTC *UTCTimeService
	TTL int
}

func NewChainService(h host.Host, v uint8) *ChainService {

	cs := &ChainService{Host:h, Version:v, TTL:TTL}
	h.SetStreamHandler(ID, cs.ChainHandler)
	return cs

}

func (cs *ChainService) ChainHandler(s core.Stream) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	timer := time.NewTimer(time.Second *TTL)
	defer timer.Stop()

	buf := make([]byte, ChunkSize)


	go func(){
		select {
		case <-timer.C:
			case <-ctx.Done():
		}
		s.Close()
	}()

	for {
		_, err := io.ReadFull(s, buf)
		if err != nil {
			log.Debug(err)
			return
		}
		_, err = s.Write(buf)
		if err != nil {
			log.Debug(err)
			return
		}

		timer.Reset(TTL)
	}


}

func (cs *ChainService) SyncFull(ctx context.Context, p peer.ID) (<-chan []byte, error) {

	s, err := cs.Host.NewStream(ctx, p, ID)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)

	go func() {
		defer close(out)

		defer s.Close()

		for {
			select {
			case <-ctx.Done():
				return

			default:
				indexes, indexesLeft, err := getBlockIndexes(s, cs)
				if err != nil {
					log.Debugf("chain synchronization error: %s", err)
					return
				}

				fmt.Printf("Index #%s done, %d Indexes left to synchronize", string(indexes), indexesLeft)

				select {
				case out <- indexes:
				case <-ctx.Done():

					fmt.Printf("Indexes done, will now sync block headers........")

					headers, headersLeft, err := getBlockHeaders(s, cs)



					if err != nil {
						log.Debugf("chain synchronization error: %s", err)
						return
					}
					fmt.Printf("%d Block header left to synchronize", headersLeft)

					select {
					case out <- headers:
					case <-ctx.Done():

						fmt.Printf("Headers done, verifying chain integrity........")

						//todo verify chain integrity

						return
					}


				}
			}

		}

	}()

	return out, nil

}

func getBlockIndexes(s core.Stream, cs *ChainService) ([]byte, int, error) {


	blockLen := len(cs.Chain.Blocks)

	i := 0

	for block := range cs.Chain.Blocks {
		buf := make([]byte, len([]byte(block.Index)))
		buf = append(buf, block.Index)
		_, err := s.Write(buf)
		if err != nil {
			return []byte{}, 0, err
		}

		resBuf := make([]byte, 8)
		_, err = io.ReadFull(s, resBuf)
		if err != nil {
			return []byte{}, 0, err
		}

		if bytes.Compare([]byte("true"), resBuf) != 1 {
			return []byte{}, 0, errors.New("bad response from remote block index sync failed")
		}

		// how many index left to synchronize
		indexLeft := blockLen - i

		i = i+1

		//if all index are done let's close the stream (temporarily)
		if indexLeft == 0 {
			s.Close()
		}



	}

	return []byte{}, 0, nil
}


func getBlockHeaders(s core.Stream, cs *ChainService) ([]byte, int, error) {


	blockLen := len(cs.Chain.Blocks)

	i := 0

	for block := range cs.Chain.Blocks {
		buf := make([]byte, len([]byte(block.Header)))
		buf = append(buf, block.Header)
		_, err := s.Write(buf)
		if err != nil {
			return nil, 0, err
		}

		resBuf := make([]byte, 8)
		_, err = io.ReadFull(s, resBuf)
		if err != nil {
			return nil, 0, err
		}

		if bytes.Compare([]byte("true"), resBuf) != 1 {
			return nil, 0, errors.New("bad response from remote block headers sync failed")
		}

		// how many index left to synchronize
		headersLeft := blockLen - i

		i = i+1

		//if all index are done let's close the stream (temporarily)
		if headersLeft == 0 {
			s.Close()
			return nil, 0, nil
		}



	}

	return nil, 0, nil
}
