package services

import (
	"context"
	"github.com/beevik/ntp"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/opentracing/opentracing-go/log"

	"io"
	"strconv"
	"time"
)



var utclog = logging.Logger("utctimeservice")

const UTCSize = 64
const utcID = "utctime"

const TTLutc = 30; //how long before disconnect from inactive service 10 min
const poolAddress = "pool.swaggit.net"

func GetUTCTime() time.Time {
	utctime, err := ntp.Time("0.pool.ntp.org")
	if err != nil {
		panic(err)
	}

	return utctime
}

type UTCTimeService struct {
	Host host.Host
	UTCTime time.Time
	UTCTimeInt64 int64
}

func NewUTCTimeService(h host.Host) {

	s := &UTCTimeService{Host:h}
	h.SetStreamHandler(utcID, s.UTCTimeHandler)

}

func (u *UTCTimeService) UTCTimeHandler(s core.Stream) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	timer := time.NewTimer(time.Second *TTLutc)
	defer timer.Stop()

	buf := make([]byte, 64)


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
			log.Error(err)
			return
		}
		_, err = s.Write(buf)
		if err != nil {
			log.Error(err)
			return
		}

		timer.Reset(time.Second * TTLutc)
	}

}

func (u *UTCTimeService) GetTime(ctx context.Context, p peer.ID) (<-chan []byte, error) {
	s, err := u.Host.NewStream(ctx, p)
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
				utctime := GetUTCTime()


				timeStr := strconv.Itoa(int(utctime.UnixNano()))

				select {
				case out <- []byte(timeStr):
				case <-ctx.Done():
					return
				}

			}

		}
	}()

	return out, nil
}