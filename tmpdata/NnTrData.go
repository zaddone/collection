package tmpdata
import (
//	nn "github.com/zaddone/collection/bpnn"
	"encoding/json"
)
type NnTrData struct {
	Data  []*Val
	NnInfo []byte
	Val  int
}
func (self *NnTrData) GetJsonData () []byte {
	d,e:=json.Marshal(self)
	if e != nil {
		panic(e)
	}
	return d
}
