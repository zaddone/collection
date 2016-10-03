package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
//	"fmt"
)
type Cl struct {
	RawPatterns []*tmpdata.Val
	DisSort  [][]*Distance
	Core   int
}
func (self *Cl) CopyCl(c *Cl) {
	self.RawPatterns = c.RawPatterns
	self.DisSort = c.DisSort
	self.Core = c.Core
}
func (self Cl) TmpAppendVal(v *tmpdata.Val) (*Cl, []int) {
	sortlist:=make([]int,len(self.RawPatterns))
//	self.RawPatterns = append(self.RawPatterns,v)
	var dis []*Distance
	for i,pa := range self.RawPatterns {
//		pa.SetD(i)
		d := new(Distance)
		d.Init(pa,v,i)
		self.DisSort[i],sortlist[i] = sortDis(self.DisSort[i],d)
		dis,_ = sortDis(dis,d)
	}
	self.RawPatterns = append(self.RawPatterns,v)
	self.DisSort = append(self.DisSort,dis)
	return &self,sortlist
}
func (self *Cl) DeleteVal(i int) {
//	v := self.RawPatterns[i]
	vd := self.DisSort[i]
	self.DisSort = append(self.DisSort[:i],self.DisSort[i+1 :]...)
	self.RawPatterns = append(self.RawPatterns[:i],self.RawPatterns[i+1 :]...)
	return
	for _,ds := range self.DisSort {
		J:
		for j,d := range ds {
			for _j,_d := range vd {
				if _d != d {
					continue
				}
				ds = append(ds[:j],ds[j+1:]...)
				vd = append(vd[:_j],vd[_j+1:]...)
				break J
			}
		}
	}
//	self.UpdateCore()
}
func (self *Cl) OutputCheck(L []int) ([]int) {
	Le := len(self.RawPatterns)-1
	if Le < 3 {
		return nil
	}
	Lc := L[self.Core]+1
	if Lc > Le {
		return nil
	}
	lastDis := self.DisSort[Le]
	core := self.RawPatterns[self.Core]
	last := self.RawPatterns[Le]
	var tmpint []int
	AppendSort := func(li []int,a int) []int {
		L := len(li)
		li = append(li,a)
		for i:= L -1 ;i >=0;i--{
			if li[i]<a {
				li[i],li[L] = li[L],li[i]
				L = i
			}else{
				break
			}
		}
		return li
	}
//	fmt.Println(self.Core,len(self.DisSort),Le,Lc)
	for _,d1 := range self.DisSort[self.Core][Lc:] {
		v1:=d1.a
		if v1 == core {
			v1 = d1.b
		}
//		fmt.Printf("%p %p %p \r\n",d1,v1,core)
		for i,d2 := range lastDis {
			v2 := d2.a
			if v2 == last {
				v2 = d2.b
			}
			if v2 == v1 {

//			if d2.a == v1 {
			//	fmt.Printf("%p %p %p %p\r\n",d1,d2,v1,v2)
				if d1.dis > d2.dis {
					tmpint = AppendSort(tmpint,d2.i)
					lastDis = append(lastDis[:i],lastDis[i+1 :]...)
				}
				break
			}
		}
	}
	return tmpint
}
func (self *Cl) UpdateCore() {
	var maxdis []*Distance
	var ls int
	for t,d := range self.DisSort {
		maxdis,ls = sortDis(maxdis,d[len(d)-1])
		if ls == 0 {
			self.Core = t
		}
	}
//	self.Core =self.RawPatterns[T]
}
func (self *Cl) FindSortVal(v *tmpdata.Val) (dis []*Distance) {

	for i,pa := range self.RawPatterns {
		d := new(Distance)
		d.Init(pa,v,i)
		dis,_ = sortDis(dis,d)
	}
	return dis

}
func (self *Cl) Append(v *tmpdata.Val) (L int)  {
//	for _,val := range self.RawPatterns {
//		if val == v {
//			panic(0)
//		}
//	}
	if self.RawPatterns == nil {
		self.RawPatterns = []*tmpdata.Val{v}
		self.DisSort = make([][]*Distance,1)
	}else{
		var dis []*Distance
		for i,pa := range self.RawPatterns {
			d := new(Distance)
			d.Init(pa,v,i)
			self.DisSort[i],_ = sortDis(self.DisSort[i],d)
			dis,_ = sortDis(dis,d)
		}
		var L int
		self.RawPatterns,L = AppendSortVal(self.RawPatterns,v)
		self.DisSort = append(append(self.DisSort[:L],dis),self.DisSort[L+1 :]...)
	}
	return L
}
func AppendSortVal(vs []*tmpdata.Val,v *tmpdata.Val) ([]*tmpdata.Val,int) {
	L := len(vs)
	vs = append(vs,v)
	for i:=L-1;i>=0;i-- {
		if vs[i].GetH() > v.GetH() {
			vs[i],vs[L] = vs[L],vs[i]
			L = i
		}else{
			break
		}
	}
	return vs,L
}
