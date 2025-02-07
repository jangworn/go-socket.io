package base

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Conn is a connection.
type Conn interface {
	FrameReader
	FrameWriter
	io.Closer
	URL() url.URL
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteHeader() http.Header
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

// ConnParameters is connection parameter of server.
type ConnParameters struct {
	PingInterval time.Duration
	PingTimeout  time.Duration
	SID          string
	Upgrades     []string
}

type jsonParameters struct {
	PingInterval int      `json:"pingInterval"`
	PingTimeout  int      `json:"pingTimeout"`
	SID          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
}

// ReadConnParameters reads ConnParameters from r.
func ReadConnParameters(r io.Reader) (ConnParameters, error) {
	var param jsonParameters
	if err := json.NewDecoder(r).Decode(&param); err != nil {
		return ConnParameters{}, err
	}

	return ConnParameters{
		SID:          param.SID,
		Upgrades:     param.Upgrades,
		PingInterval: time.Duration(param.PingInterval) * time.Millisecond,
		PingTimeout:  time.Duration(param.PingTimeout) * time.Millisecond,
	}, nil
}

// WriteTo writes to w with json format.
func (p ConnParameters) WriteTo(w io.Writer) (int64, error) {
	arg := jsonParameters{
		SID:          p.SID,
		Upgrades:     p.Upgrades,
		PingInterval: int(p.PingInterval / time.Millisecond),
		PingTimeout:  int(p.PingTimeout / time.Millisecond),
	}
	writer := writer{
		w: w,
	}
	err := json.NewEncoder(&writer).Encode(arg)
	return writer.i, err
}

type writer struct {
	i int64
	w io.Writer
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.i += int64(n)
	return n, err
}
