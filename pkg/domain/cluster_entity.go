package domain

type ClusterInfo struct {
	ID      uint64 `json:"ID"`
	Address string `json:"Address"`
	Role    uint8  `json:"Role"`
	Leader  bool   `json:"Leader"`
}
