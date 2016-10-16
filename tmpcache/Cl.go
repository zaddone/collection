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
	oic   *Cl
}
func (self *Cl) Copy(c *Cl) {
	self.RawPatterns = c.RawPatterns
	self.DisSort = c.DisSort
	self.Core = c.Core
//	fmt.Println(self.CountY,c.CountY)
	self.CountY = c.CountY
}
func (self *Cl) Clear(){
	self.RawPatterns = nil
	self.DisSort = nil
	self.Core = -1
	self.CountY = [3]int{0,0,0}
}
func (self *Cl) SetOic(c *Cl) {
	self.oic = c
}
func (self Cl) TmpAppendVal(v *tmpdata.Val) (*Cl, []int) {
	L := len(self.RawPatterns)
	sortlist:=make([]int,L)
	dis:=make([]*Distance,L)
	disSort:= make([][]*Distance,L+1)
	RawPatterns:= make([]*tmpdata.Val,L+1)
	copy(RawPatterns,self.RawPatterns)
	RawPatterns[L] = v
	disSort[L] = dis
//	copy(disSort,self.DisSort)
//	fmt.Printf("append %p\r\n",v)
	for i,pa := range self.RawPatterns {
//		for u,_d := range self.DisSort[i] {
//			isF := false
//			for h,_v := range self.RawPatterns {
//				if h == i {
//					continue
//				}
//				if _d.a == _v || _d.b == _v {
//					isF = true
//				}
//			}
//			if !isF {
//				fmt.Println(_d,i,u,len(self.RawPatterns),len(self.DisSort[i]))
//				fmt.Println(self.RawPatterns)
//				panic(77)
//			}
//		}

		d := new(Distance)
		d.Init(pa,v,i)
//		ds,sortlist[i] = sortDis(ds,d)
//		disSort[i] = ds
		diss:= make([]*Distance,L)
		copy(diss,self.DisSort[i])
//		fmt.Println(diss)
		diss[L-1] = d
		sortlist[i] = appendDis(diss,L-1)
		disSort[i] = diss
//		disSort[i],sortlist[i] = sortDis(diss,d)

//		sortlist[i] = 0
//		disSort[i] = append(disSort[i],d)

		self.Sum+=d.dis
		self.DisCount ++
		dis[i]=d
		appendDis(dis,i)
//		dis,_ = sortDis(dis,d)
	}


	self.RawPatterns = RawPatterns
	self.DisSort = disSort
	self.CountY[v.Y]++
//	var L int
//	self.RawPatterns,L = AppendSortVal(RawPatterns,v)
//	self.DisSort = append(append(disSort[:L],dis),disSort[L:]...)
	return &self,sortlist
}
func (self *Cl) GetPG() float64{
	return self.Sum/self.DisCount
}
func (self *Cl) DeleteVal(I int) (*tmpdata.Val,error) {
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
//		var delD *Distance
		str := ""
		for j,d := range ds {
			str += fmt.Sprintln(d)
//			d := ds[j]
			if d.a == v || d.b == v {
//				delD = d
				self.DisSort[i] = append(ds[:j],ds[j+1:]...)
				self.Sum -= d.dis
				self.DisCount --
				isD = true
				break
			}
		}
//		fmt.Printf("e %d %d %p %p %p\r\n",i,I,delD,v,self.RawPatterns[i])
//		fmt.Println(delD)
//		fmt.Println(str)
		if !isD {
			panic(55)
			return v,fmt.Errorf("e:%d %d",len(self.DisSort[i]),I)
			errI = append(errI,i)
		}
	}
	if errI != nil {
		fmt.Println("err:",errI)
		panic(2)
	}
	self.CountY[v.Y]--
	return v,nil

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

func AppendSort (li []int,a int) []int {
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
func (self *Cl) UpdateCore() {
	maxdis:=make([]*Distance,len(self.DisSort))
	var ls int
	for t,d := range self.DisSort {
		maxdis[t] = d[len(d)-1]
		ls = appendDis(maxdis,t)
//		fmt.Println(t,len(d),"____")
//		maxdis,ls = sortDis(maxdis,d[len(d)-1])
		if ls == 0 {
			self.Core = t
		}
	}
//	self.Core =self.RawPatterns[T]
}
func (self *Cl) FindSortVal(v *tmpdata.Val) (dis []*Distance) {

	dis = make([]*Distance,len(self.RawPatterns))
	for i,pa := range self.RawPatterns {
		d := new(Distance)
		d.Init(pa,v,i)
		dis[i]=d
		appendDis(dis,i)
//		dis,_ = sortDis(dis,d)
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
		dis := make( []*Distance,len(self.RawPatterns))
		for i,pa := range self.RawPatterns {
//			for _,_d := range self.DisSort[i] {
//				isF := false
//				for h,_v := range self.RawPatterns {
//					if h == i {
//						continue
//					}
//					if _d.a == _v || _d.b == _v {
//						isF = true
//					}
//				}
//				if !isF {
//					panic(77)
//				}
//			}
			d := new(Distance)
			d.Init(pa,v,i)
			if self.DisSort[i] == nil {
				self.DisSort[i] = []*Distance{d}
			}else{
				self.DisSort[i],_ = sortDis(self.DisSort[i],d)
			}

			dis[i] = d
			appendDis(dis,i)
//			dis,_ = sortDis(dis,d)
			self.Sum+=d.dis
			self.DisCount ++
		}
		var L int
		Le:= len(self.RawPatterns)
		self.RawPatterns,L = AppendSortVal(self.RawPatterns,v)
		if L == Le {
			self.DisSort = append(self.DisSort,dis)
		}else if Le == 0 {
			self.DisSort = append([][]*Distance{dis},self.DisSort...)
		}else{
			tmpDis:=append([][]*Distance{},self.DisSort[L:]...)
			self.DisSort = append(append(self.DisSort[:L],dis),tmpDis...)
		}
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
