package client

type client struct {
	channels []chan any
}

func (c *client) Write(p []byte) (n int, err error) {
	data := map[string]string{"type": "manuallyRunScript", "message": string(p)}
	for _, channel := range c.channels {
		select {
		case channel <- data:
		default:

		}

	}
	return len(p), nil
}

var (
	MyClient *client
)

func init() {
	MyClient = new(client)
}

func AddChan(c chan any) {
	MyClient.channels = append(MyClient.channels, c)
}
