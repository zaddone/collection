package block
import (
	"fmt"
	"net"
	"github.com/zaddone/collection/signal"
	"github.com/zaddone/collection/tmpcache"
//	"math"
//	"os"
//	"github.com/zaddone/collection/tmpdata"
)
const (
	Timeout   uint64   = 30000
	BB   float64   = 2
)
type Feature struct {

	Atlias   []*signal.Atlias

	lastAt   *signal.Atlias

	last   []*BlockInfo

	lastBak  []*BlockInfo
	lastBak1  []*BlockInfo

//	testCheck *WaitData

	Cache  *tmpcache.Cache
	Signal       *signal.GraphicSignal

	bk  float64
	kill float64

	IsUpdate bool
	Week   int
	Hour   int
	Name   string

	isCount float64
	isOk  float64
	isOk1  float64
	isOk2  float64
	vals [3]int

}
func (self *Feature) GetCache() *tmpcache.Cache{
	return self.Cache
}
func (self *Feature) Init(listener *net.Listener){

	self.Signal = new(signal.GraphicSignal)
	self.Cache = new(tmpcache.Cache)
	self.Cache.Init(listener)

}
func (self *Feature) SetSignal(path string){
	self.Signal.FilePath=path
}
func (self *Feature) GetUpdate() bool {
	return self.IsUpdate
}
func (self *Feature) SetUpdate( up bool) {
	self.IsUpdate = up
}
func (self *Feature) SetGraphic(week int, hour int, name string)  {

	if self.Week!=week{
		self.Week=week
		//fmt.Println(name,week,self.ShowInfo())
		fmt.Println(name,week,self.ShowTest(),self.Cache.ShowInfo())
//		self.ClearInfo()
//		self.block.Count=0
	}
//	self.Maps = make(map[string]*DataList)
	self.Hour = hour
	self.Name = name
//	if self.Week == 4 {
//		os.Exit(0)
//	}

}
func (self *Feature) ShowTest() string {
	msg:= fmt.Sprintf("%.5f %.1f %.1f %.1f %.1f", self.isOk/self.isCount,self.isOk,self.isCount,self.isOk1,self.isOk2)
	self.ClearInfo()
	return msg
}
func (self *Feature) ClearInfo(){
	self.isOk = 0
	self.isOk1 = 0
	self.isOk2 = 0
	self.isCount = 0
}
func (self *Feature) Clear(){

	self.Atlias = nil
	if self.last != nil {
//		for _,la := range self.last {
//			if la.SendData(self.Cache) {
//				self.isOk ++
//			}
//		}

		self.last = nil
	}
//	self.lastBak = nil

//	self.testCheck = nil
	self.lastBak = nil

}
func (self *Feature) Analysis(at *signal.Atlias) {

	if at.Bid == 0  {
		self.Clear()
		//self.lastAt = at
		return
	}

//	if at.Ask < at.Bid {
//		self.Clear()
//		return
//	}
//
//	if at.Begin <= at.Bid {
//		self.Clear()
//		return
//	}
//	if at.Last == 0 || at.Begin == 0 {
//		self.Clear()
//		return
//	}

	last:=self.lastAt
	if last == nil {
		self.lastAt = at
		self.Atlias = []*signal.Atlias{at}
		return
	}
	if at.Msc < last.Msc {
		return
	}
	if last.Begin != at.Last {
		if last.Msc == at.Msc {
			return
		}
		self.Clear()
		self.Atlias = []*signal.Atlias{at}
		self.lastAt = at
		return
	}
	if self.Atlias == nil {
		self.Atlias = []*signal.Atlias{at}
		self.lastAt = at
		return

	}
	self.appendAtlias(at)
	self.lastAt = at

}
func (self *Feature) GetAtDiff() (diff float64) {
	for _,a := range self.Atlias {
		diff +=a.Diff
	}
	return diff
}
func (self *Feature) appendAtlias(at *signal.Atlias) {

	if  len(self.Atlias)%2!=0 {
		self.Overlastbak(at)
		self.Atlias = append(self.Atlias,at)
		return
	}

	self.Signal.Init(at)
	self.bk = at.Ask-at.Bid
	self.kill = self.bk * BB
	last := self.GetBlockInfo(at)
	if last == nil {
		self.Atlias = append(self.Atlias,at)
		return
	}
	last.TimeOut(Timeout)
	if last.CheckA() {
		self.FindBak(last)
	}
	if last.SetCovs()  {
		self.GetSignal(last)
		self.lastBak = append(self.lastBak,last)
		self.isOk ++
	}
	self.last = append(self.last,last)
	self.isCount++
	end := last.end + int(float64(len(last.atlias[last.end:]))*0.7)
	self.Atlias = last.atlias[end:]
//	self.Atlias = []*signal.Atlias{at}

}
func  FindBakX(la,lb *BlockInfo) bool {
	for _,l := range la.atlias[:la.end] {
		for _,j := range lb.atlias[lb.end:] {
			if l == j {
				return true
			}
		}
	}
	return false
}
func (self *Feature) FindBak(la *BlockInfo) (begin int) {
	begin = -1
	for i,lb := range self.lastBak {
		if lb == la {
			begin = i-1
			break
		}
		if FindBakX(la,lb) {
			lb.Next = la
			begin = i
			self.isOk1 ++

		}
	}
	begin++
	if begin == 0 {
		return begin
	}
	for _,b := range self.lastBak[:begin]{
		b.SendData(self.Cache)
	}
	self.isOk2 += float64(begin)
	fmt.Printf("%.1f %.5f %d\r",self.isOk,self.isOk1/self.isOk2,len(self.Cache.CluMap))
	if begin >= len(self.lastBak) {
		self.lastBak = nil
	}else{
		self.lastBak = self.lastBak[begin :]
	}
	return begin
}
func (self *Feature) Overlastbak(at *signal.Atlias) {
	L := len(self.last)
	if L ==0  {
		return
	}
	begin := -1
//	isSend := false
//	bb := -1
	for i:=L-1;i>=0;i-- {
		if self.last[i].UpdateAts([]*signal.Atlias{self.lastAt,at}){
			begin = i
		}
	}
	if begin == -1 {
		self.last = nil
	}else{
		self.last = self.last[begin:]
	}
}

func (self *Feature) GetBlockInfo(at *signal.Atlias) *BlockInfo {

	L:= len(self.Atlias)
	end := -1
	price := -1

	buy:=0.0
//	sell:=0.0
//	dk :=at.Ask-at.Bid
	for i:=L-1;i>=0;i-- {
		a:= self.Atlias[i]
		buy = at.Begin - a.Last //+ at.Diff
//		if at.Msc - a.Msc   {
//			break
//		}
		if buy > self.kill {
			price = 1
			end = i
			break
		}else if buy < -self.kill {
			price = 0
			end = i
			break
		}
	}

	if end <= 0 {
		return nil
	}
	if (at.Begin - self.Atlias[end].Last >0) != (price >0 ) {
//		fmt.Println(at)
		return nil
	}
	b := &BlockInfo{atlias:append(self.Atlias,at),end:end,IsBuy:price,TestBuy:3}


	if price == 1{

		for i,a:=range self.Atlias{
			buy = at.Begin - a.Last// + at.Diff
			if buy > self.kill {
				b.begin = i
				break
			}
		}
	}else {
		for i,a:=range self.Atlias{
			buy = at.Begin - a.Last //+ at.Diff
			if buy < -self.kill {
				b.begin = i
				break
			}
		}

	}
	return b

}

func (self *Feature) GetSignal(la *BlockInfo){
	t := la.Forecast(self.Cache)

	if t< 0 {
	//	fmt.Println(err)
		return
	}

	if self.IsUpdate {
		kill := int(self.bk+self.kill)/2
		if t ==1 {
			self.Signal.Create(1,kill)
		}else if t == 0  {
			self.Signal.Create(-1,kill)
		}else{
			self.Signal.Create(0,kill)
		}
	}

}

