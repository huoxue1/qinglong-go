package system

type System struct {
	IsInitialized  bool   `json:"isInitialized"`
	Version        string `json:"version"`
	LastCommitTime string `json:"lastCommitTime"`
	LastCommitId   string `json:"lastCommitId"`
	Branch         string `json:"branch"`
}
