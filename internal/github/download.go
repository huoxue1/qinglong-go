package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"os"
)

func Download() {

	git.PlainClone("./data/faker2", false, &git.CloneOptions{URL: "https://github.com/shufflewzc/faker2.git", ProxyOptions: transport.ProxyOptions{
		URL:      "socks5://124.70.38.162:12785",
		Username: "admin",
		Password: "qqqq",
	}, Progress: os.Stdout})
}
