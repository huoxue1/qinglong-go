package models

type AuthFile struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Tokens   struct {
		Desktop string `json:"desktop"`
		Mobile  string `json:"mobile"`
	} `json:"tokens"`
	Lastlogon           int64  `json:"lastlogon"`
	Retries             int    `json:"retries"`
	Lastip              string `json:"lastip"`
	Lastaddr            string `json:"lastaddr"`
	Platform            string `json:"platform"`
	IsTwoFactorChecking bool   `json:"isTwoFactorChecking"`
	TwoFactorSecret     string `json:"twoFactorSecret"`
}
