package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"encoding/json"
	"fmt"
)
const (
	LongLen int = 100
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
func (self *Clus) AppendVal(v *tmpdata.Val,isBack int)  {

	v.SetH(self.ValCount)
	self.ValCount++
	if len(self.Clu) < 2 {
		clu := new(Cl)
		clu.Append(v)
		self.Clu = append(self.Clu,clu)
		return
	}
	clus,tmps := self.getTmpVal(nil)

//	tmpClu := new(Cl)
//	tmpClu.Init(tmps...)
	tmpClu := &Cl{RawPatterns:tmps}
	vs := tmpClu.FindSortVal(v)

	Long := len(vs)
	if Long > LongLen {
		Long = LongLen
	}
	var tmpDisSort []*Distance
	var tmpc []*Cl = make([]*Cl,Long)
	var tmpli [][]int = make([][]int,Long)
	var minT int
	var ls int
	for j,i := range vs[:Long] {
		cl,L:=clus[i].TmpAppendVal(v)
		tmpc[j] = cl
		tmpli[j] = L
		tmpDisSort,ls = sortDis(tmpDisSort,cl.DisSort[len(cl.DisSort)-1][0])
		fmt.Println(j,ls)
		if ls == 0 {
			minT = j
		}
	}
	if tmpDisSort[0].a.C != v.C {
		newClu := new(Cl)
		newClu.Append(v)
		self.Clu = append(self.Clu,newClu)
		minT = -1
	}else{
		cp :=tmpc[minT]
		cp.UpdateCore()
		clus[vs[minT]] = cp
	}
	if isBack > LongIsBack {
		return
	}
	for j,i := range vs[:Long] {
		if j == minT {
			continue
		}
		ls := tmpc[j].OutputCheck(tmpli[j])
		if ls == nil {
			continue
		}
		cp := clus[i]
		for _,_i := range ls {
			newBack := isBack+1
			self.AppendVal(cp.RawPatterns[_i],newBack)
			cp.DeleteVal(_i)
		}
		cp.UpdateCore()
	}

}
func (self *Clus) getTmpVal(c *Cl) (tmpc []*Cl,tmp []*tmpdata.Val) {
	L := len(self.Clu)
	for i := L -1;i>=0;i--{
		_c := self.Clu[i]
		le :=len(_c.RawPatterns)
		if le == 0 {
			self.Clu = append(self.Clu[:i],self.Clu[i+1:]...)
			continue
		}
		if _c == c {
			continue
		}
		tmpc = append(tmpc,_c)
		tmp = append(tmp,_c.RawPatterns[_c.Core])
//		tmp = append(tmp,_c.RawPatterns[0])
//		if _c.Core == nil {
//			tmp = append(tmp,_c.RawPatterns[le-1])
//		}else{
//			tmp = append(tmp,_c.RawPatterns[_c.Core])
//			//tmp = append(tmp,_c.Core)
//		}
	}
	return tmpc,tmp
}
