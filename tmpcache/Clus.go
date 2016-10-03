package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"encoding/json"
	"sync"
//	"fmt"
)
const (
	LongLen int = 10
	LongIsBack int = 2
)
type Clus struct {
	Clu []*Cl
	ca *Cache
//	Len int
	ValCount int
}
func (self *Clus) Init(ca *Cache)  {
	self.ca = ca
}
func (self *Clus) Getbakdb() (db []byte) {
	db,err := json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return db
}
func (self *Clus) LoadData(db []byte) error {
	return json.Unmarshal(db,self)
}
func (self *Clus) Append(v *tmpdata.Val)  {
	v.SetH(self.ValCount)
	self.ValCount++
	if len(self.Clu) < 2 {
		clu := new(Cl)
		clu.Append(v)
		self.Clu = append(self.Clu,clu)
		return
	}
	self.AppendVal(v,1)
}
func (self *Clus) AppendVal(v *tmpdata.Val,isBack int)  {

	clus,tmps := self.getTmpVal()

	tmpClu := &Cl{RawPatterns:tmps}
	dis := tmpClu.FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
	var tmpDisSort []*Distance
	var tmpc []*Cl = make([]*Cl,Long)
	var tmpli [][]int = make([][]int,Long)
	var minT int = -1
	var ls int
	for j,d := range dis[:Long] {
		cl,L:=clus[d.i].TmpAppendVal(v)
		tmpc[j] = cl
		tmpli[j] = L
		md :=cl.DisSort[len(cl.DisSort)-1][0]

//		md.i = j
//		tmpDisSort,_ = sortDis(tmpDisSort,md)

		tmpDisSort,ls = sortDis(tmpDisSort,md)
		if ls == 0 {
			minT = j
		}
	}

	isUp := false
	for _,td := range tmpDisSort{
		if td.a.C == v.C {
			minT = td.i
			cp :=tmpc[minT]
			cp.UpdateCore()
			clus[dis[minT].i].CopyCl(cp)
			isUp = true
			break
		}
	}
	if !isUp {
		newClu := new(Cl)
		newClu.Append(v)
		self.Clu = append(self.Clu,newClu)
	}

//	if tmpDisSort[0].a.C != v.C {
//		newClu := new(Cl)
//		newClu.Append(v)
//		self.Clu = append(self.Clu,newClu)
//		minT = -1
//	}else{
//		cp :=tmpc[minT]
//		cp.UpdateCore()
//		clus[dis[minT].i].CopyCl(cp)
//	}
//	return
	if isBack > LongIsBack {
		return
	}
	var vals []*tmpdata.Val
	for j,d := range dis[:Long] {
		if j == minT {
			continue
		}
		ls := tmpc[j].OutputCheck(tmpli[j])
		if ls == nil {
			continue
		}
		cp := clus[d.i]
//		fmt.Println(ls,len(cp.RawPatterns))
		if len(cp.RawPatterns)-len(ls) < 2 {
			vals = append(vals,cp.RawPatterns...)
			cp.RawPatterns = nil
		}else{
			for _,_i := range ls {
//				val:=cp.RawPatterns[_i]
				vals = append(vals, cp.RawPatterns[_i])
				cp.DeleteVal(_i)
			}
			cp.UpdateCore()
		}
	}
//	fmt.Println(len(vals),isBack)
	if vals == nil {
		return
	}
	walk := new(sync.WaitGroup)
	walk.Add(1)
	go self.syncAppendVal(vals,isBack,walk)
	walk.Wait()

}
func (self *Clus) syncAppendVal(vals []*tmpdata.Val,isBack int,walk *sync.WaitGroup) {
//	fmt.Println("count:",len(vals),"____________")
	for _,val := range vals {
//		fmt.Println(i)
		self.AppendVal(val,isBack+1)
	}
	walk.Done()
}
func (self *Clus) getTmpVal() (tmpc []*Cl,tmp []*tmpdata.Val) {
	L := len(self.Clu)
	for i := L -1;i>=0;i--{
		_c := self.Clu[i]
//		le := len(_c.RawPatterns)
		if _c.RawPatterns == nil {
			self.Clu = append(self.Clu[:i],self.Clu[i+1:]...)
			continue
		}
		tmpc = append(tmpc,_c)
		tmp = append(tmp,_c.RawPatterns[_c.Core])
	}
	return tmpc,tmp
}
