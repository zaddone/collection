package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"fmt"
)
type Cl struct {
	RawPatterns []*tmpdata.Val
	DisSort  [][]*Distance
	CountY   [3]int
	Sum   float64
	DisCount float64
	Core   int
}
func (self *Cl) Copy(c *Cl) {
	self.RawPatterns = c.RawPatterns
	self.DisSort = c.DisSort
	self.Core = c.Core
	self.CountY = c.CountY
}
func (self *Cl) Clear(){
	self.RawPatterns = nil
	self.DisSort = nil
	self.Core = -1
	self.CountY = [3]int{0,0,0}
}
func (self Cl) TmpAppendVal(v *tmpdata.Val) (*Cl, []int) {
	sortlist:=make([]int,len(self.RawPatterns))
	var dis []*Distance
	disSort:= make([][]*Distance,len(self.DisSort))
	RawPatterns:= make([]*tmpdata.Val,len(self.RawPatterns))
	copy(RawPatterns,self.RawPatterns)
	copy(disSort,self.DisSort)
//	fmt.Printf("append %p\r\n",v)
	for i,pa := range RawPatterns {
		d := new(Distance)
		d.Init(pa,v,i)
//		ds,sortlist[i] = sortDis(ds,d)
//		disSort[i] = ds
		diss:= make([]*Distance,len(disSort[i]))
		copy(diss,disSort[i])
		disSort[i],sortlist[i] = sortDis(diss,d)

//		sortlist[i] = 0
//		disSort[i] = append(disSort[i],d)

		self.Sum+=d.dis
		self.DisCount ++
		dis,_ = sortDis(dis,d)
	}
//	for i,ds := range disSort {
//
//		ds := disSort[i]
//		if len(ds) != len(self.RawPatterns) {
//			panic(98)
//		}
//		pa := self.RawPatterns[i]
//		for j,d := range ds {
//			v1 := d.a
//			if v1 == pa {
//				v1 = d.b
//			}
//			if v1 == v {
//				continue
//			}
//			isb:= false
//			for _,p := range self.RawPatterns {
//				if p == v1 {
//					isb = true
//					break
//				}
//			}
//			if !isb {
//				fmt.Println(j,d)
//				fmt.Println(self.RawPatterns)
//				fmt.Printf("%p \r\n",v1)
//				panic(99)
//			}
//		}
//	}

	self.RawPatterns = append(RawPatterns,v)
	self.DisSort = append(disSort,dis)
	self.CountY[v.Y]++
//	var L int
//	self.RawPatterns,L = AppendSortVal(RawPatterns,v)
//	self.DisSort = append(append(disSort[:L],dis),disSort[L:]...)
	return &self,sortlist
}
func (self *Cl) GetPG() float64{
	return self.Sum/self.DisCount
}
func (self *Cl) DeleteVal(I int) {
//	fmt.Println("del")
//	ld := len(self.DisSort)
//	lr := len(self.RawPatterns)
//	for t,di:= range self.DisSort {
//		li := len(di)+1
//		fmt.Println(t,li,ld,lr," ")
//		fmt.Println(di)
//		if ld  != li {
//			fmt.Println(di)
//			panic(10)
//		}
//	}

	v:=self.RawPatterns[I]
	self.DisSort = append(self.DisSort[:I],self.DisSort[I+1 :]...)
	self.RawPatterns = append(self.RawPatterns[:I],self.RawPatterns[I+1 :]...)
	var errI []int
	for i,ds := range self.DisSort {
//		L := len(ds)
		isD := false
//		fmt.Println("b",i,L,len(self.DisSort)," ")
//		for j:= L-1;j>=0;j-- {
		for j,d := range ds {
//			d := ds[j]
			if d.a == v || d.b == v {
				self.DisSort[i] = append(ds[:j],ds[j+1:]...)
				self.Sum -= d.dis
				self.DisCount --
				isD = true
				break
			}
		}
//		fmt.Println(self.DisSort[i])
		if !isD {
			errI = append(errI,i)
		}
//		fmt.Println("e",i,len(self.DisSort[i]))
	}
	if errI != nil {
		fmt.Println(errI)
		panic(2)
	}
	self.CountY[v.Y]--

//	ld := len(self.DisSort)
//	lr := len(self.RawPatterns)
//	for t,di:= range self.DisSort {
//		li := len(di)+1
////		fmt.Println(t,li,ld,lr," ")
//		if ld  != li {
//			fmt.Println(di)
//			panic(2)
//		}
//	}

//	panic(0)
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
//	last := self.RawPatterns[Le]
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
	for _,d1 := range self.DisSort[self.Core][Lc:] {
		v1:=d1.a
		if v1 == core {
			v1 = d1.b
		}
		isB := -1
		for j,d2 := range lastDis {
//			fmt.Println(d2)
			if d2.a == v1 {//|| d2.b == v1 {
				if d1.dis > d2.dis {
					tmpint = AppendSort(tmpint,d2.i)
				}
				isB = j
				break
			}
		}
//		fmt.Println("to:",d1)
//		fmt.Println("do:",self.RawPatterns)
//		fmt.Printf("%d %p %p\r\n",Lc,core,last)
		if isB<0 {
//			fmt.Println(v1)
//			fmt.Printf("%p %p %p \r\n",d1,v1,core)
			panic(3)
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
		self.DisSort = [][]*Distance{nil}
	}else{
		var dis []*Distance
		for i,pa := range self.RawPatterns {
			d := new(Distance)
			d.Init(pa,v,i)
			self.DisSort[i],_ = sortDis(self.DisSort[i],d)
			dis,_ = sortDis(dis,d)
			self.Sum+=d.dis
			self.DisCount ++
		}
		var L int
		self.RawPatterns,L = AppendSortVal(self.RawPatterns,v)
		self.DisSort = append(append(self.DisSort[:L],dis),self.DisSort[L:]...)
	}
	self.CountY[v.Y]++
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
