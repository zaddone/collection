package signal

import (
	"encoding/json"
	"fmt"
)

type Atlias struct {
	Diff   float64   `json:"diff"`
	Time   uint16  `json:"time"`
	Volume float64 `json:"volume"`
	Msc    uint64  `json:"msc"`
	Last   float64   `json:"last"`
	Begin  float64   `json:"begin"`
	Ask   float64   `json:"ask"`
	Bid   float64   `json:"bid"`
}

func (at *Atlias) GetAtlas(data []byte) error {
	err := json.Unmarshal(data, &at)
	if err != nil {
		fmt.Println(err)
	}
	at.Diff = -at.Diff
	return err
}

func (self *Atlias) CompSame(at *Atlias) bool {
	return *self == *at
}
