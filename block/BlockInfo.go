package block
import (
//	"encoding/json"
	"github.com/zaddone/collection/curves"
	"github.com/zaddone/collection/signal"

	"github.com/zaddone/collection/tmpcache"
	"github.com/zaddone/collection/tmpdata"

//	"os"
//	"fmt"
	"math"
)
//type Matrix interface{
//	GetWeights() []float64
////	GetKey() []int
////	SetCovs(val ...float64) int
//	AppendAtlias(at []*signal.Atlias) error
//
//}
type BlockInfo struct {

	atlias   []*signal.Atlias
	begin  int
	end int
	IsBuy  int     `json:"isbuy"`
//	covs   *cov.CovMatrix
//	covs   Matrix
	covs   *curves.Curve
	Sort  []int  `json:"sort"`
	Weight  []float64 `json:"Weight"`
	Next *BlockInfo
	Timeout int
	TestBuy int

}
func (self *BlockInfo) UpdateAts(Ats []*signal.Atlias) bool {
	if self.atlias[len(self.atlias)-1].Begin != Ats[0].Last {
//		fmt.Println("b l is err")
		return false
	}
	diff :=0.0
	for _,a := range Ats {
		diff += a.Diff
	}
	if (diff >0) == (self.IsBuy >1) {
		self.atlias = append(self.atlias,Ats...)
		return true
	}else{
		return false
	}
}
func (self *BlockInfo) SetCovs() bool {
	self.covs = new(curves.Curve)
	err := self.covs.AppendAtlias(self.atlias)
	if err != nil {
		return false
	}
	return true
}
func (self *BlockInfo) Forecast(Cache *tmpcache.Cache) int {
//	self.TestBuy = self.IsBuy
//	return self.TestBuy
	b,err :=Cache.Forecast(&tmpdata.Val{X:self.covs.GetWeights(),K:[]int{0},Y:2,C:self.IsBuy,Cur:self.covs})
	if err != nil {
		self.TestBuy = 3
	}else{
		self.TestBuy = b
	}
	return self.TestBuy
}
func (self *BlockInfo) SendData(Cache *tmpcache.Cache) int {
	Y:=2
//	C:=1
	if self.Next != nil {
		Y = self.Next.IsBuy
//		C = 0
	}
	Cache.Input(&tmpdata.Val{X:self.covs.GetWeights(),K:[]int{0},Y:Y,C:self.IsBuy,Cur:self.covs})
//	if self.TestBuy < 2 {
	if self.TestBuy < 3 {
//	if self.TestBuy < 2 && Y < 2 {
		Cache.TestInfo(self.TestBuy == Y)
	}
	return Y
}
func (self *BlockInfo) CheckA() bool {
//	if self.atlias[self.end].Msc-self.atlias[self.begin].Msc > 10000 {
//		return false
//	}
	if self.begin != self.end {
		Bdiff := 0.0
		Max := 0.0
		diff := 0.0
//		fmt.Println(self.begin,self.end)
		for _,at := range self.atlias[self.begin:self.end] {
			Bdiff += at.Diff
			if self.IsBuy == 1 {
				Max = math.Max(Max,Bdiff)
			}else{
				Max = math.Min(Max,Bdiff)

			}
		}
		for _,at := range self.atlias[self.end:] {
			diff += at.Diff
		}
		if Max/(Bdiff+diff) > 0.2 {
			return false
		}
	}
	return true
}
func (self *BlockInfo) TimeOut( out uint64) bool {
	if self.Timeout == 0 {
		if self.atlias[self.end].Msc-self.atlias[0].Msc > out {
			self.Timeout = 2
			self.Truncated(out/3)
		}else{
			self.Timeout = 1
		}
	}
	if self.Timeout == 1 {
		return false
	}else{
		return true
	}
}
func (self *BlockInfo) Truncated(d uint64) bool {
	if self.end ==0 {
		return false
	}
	for i := self.end-1;i>=0;i--{
		if self.atlias[self.end].Msc-self.atlias[i].Msc > d {
			b :=i+1
			self.atlias = self.atlias[b:]
			if self.begin<b{
				self.begin = 0
			}else{
				self.begin = self.begin -b
			}
			self.end = self.end -b
			return true
		}
	}
	return true
}
//func (self *BlockInfo) CheckB() bool {
//	if self.atlias[self.end].Msc-self.atlias[0].Msc > 15000  {
//		if self.begin ==0 {
//			return false
//		}
////		isf := false
//		for i := self.begin-1;i>=0;i--{
//			if self.atlias[self.begin].Msc-self.atlias[i].Msc > 5000 {
//				b :=i+1
//		//		if b == self.begin {
//		//			return false
//		//		}
//				self.atlias = self.atlias[b:]
//				self.begin = self.begin -b
//				self.end = self.end -b
//		//		isf =true
//				break
//			}
//		}
////		if !isf {
////			return false
////		}
//	}
//	return true
//}
