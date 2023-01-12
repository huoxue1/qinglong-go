package client

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync"
)

type Clients struct {
	sync.Map
}

func (c *Clients) Write(p []byte) (n int, err error) {
	data := map[string]string{"type": "manuallyRunScript", "message": string(p)}
	content, _ := json.Marshal(data)
	var deleSlince []string
	c.Range(func(key, value any) bool {
		conn := value.(*Client).conn
		writer := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
		_, err := writer.Write(content)
		if err != nil {
			deleSlince = append(deleSlince, key.(string))
			return true
		}
		err = writer.Flush()
		if err != nil {
			deleSlince = append(deleSlince, key.(string))
			return true
		}
		return true
	})
	for _, s := range deleSlince {
		c.Delete(s)
	}
	return len(p), nil
}

var (
	MyClient *Clients
)

func init() {
	MyClient = new(Clients)
}

type Client struct {
	conn net.Conn
}

func AddWs(id string, conn net.Conn) {
	MyClient.Store(id, &Client{
		conn: conn,
	})

}
