package tmpdata
import (
	"encoding/json"
	"crypto/sha1"
	"github.com/zaddone/collection/signal"
)


type Val struct {
	X   []float64
	K   []int
	Y   int
	C   int
	d   int
	h   int
	at  []*signal.Atlias
}
func (self *Val) SetAt(at  []*signal.Atlias) {
	self.at = at
}
func (self *Val) GetAt() []*signal.Atlias {
	return self.at
}

func (self *Val) SetH(d int){
	self.h = d
}
func (self *Val) GetH() int {
	return self.h
}
func (self *Val) SetD(d int){
	self.d = d
}
func (self *Val) GetD() int {
	return self.d
}

func (self *Val) GetKey() []byte {
	data:= self.GetBytes()
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}
func (self *Val) GetBytes() []byte {
	d,err:=json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return d
}
func (self *Val) Diff(v *Val) bool {

	if self.Y != v.Y {
		return false
	}
	if len(self.X) != len(v.X) {
		return false
	}
	for i,x := range self.X {
		if v.X[i] != x {
			return false
		}
	}
	return true

}
