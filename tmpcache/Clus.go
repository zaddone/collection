package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"github.com/zaddone/collection/curves"
	"encoding/json"
	"sync"
	"fmt"
)
const (
	LongLen int = 10
	LongIsBack int = 1
	MaxLong  int = 10000
)
type Clus struct {
	Clu []*Cl
	ca *Cache
//	Len int
	ValCount int
	ErrMin *Distance
	OkMin *Distance
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
//	fmt.Printf("%p\r\n",v)
	v.SetH(self.ValCount)
	self.ValCount++
//	fmt.Println(self.ValCount)
	if len(self.Clu) < 2 {
		clu := new(Cl)
		clu.Append(v)
		self.Clu = append(self.Clu,clu)
		return
	}
	self.AppendVal(v,1)
}
func (self *Clus) FindSame(v *tmpdata.Val) *Distance {
	clus,tmps := self.getTmpVal()
	if clus == nil {
		return nil
	}
	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
	tmpDisSort:=make([]*Distance,Long)
//	sortlist:=make([]int,Long)
	var ls,out int
	for j,d := range dis[:Long] {
		cl,_:=clus[d.i].TmpAppendVal(v)
		tmpDisSort[j] = cl.DisSort[len(cl.DisSort)-1][0]
//		sortlist[j] = j
		ls = appendDis(tmpDisSort,j)
		if ls == 0 {
			out = 0
		}else{
			out ++
			if out >10 {
				Long = j+1
				break
			}
		}
	}
	return tmpDisSort[0]
}
func (self *Clus) AppendVal1(v *tmpdata.Val,isBack int)  {

//	clus,tmps := self.getTmpVal()
//	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
//	Long := len(dis)
//	if Long > LongLen {
//		Long = LongLen
//	}
//	for j,d := range dis[:Long] {
//		oic := clus[d.i]
//		cl,L:=oic.TmpAppendVal(v)
//		cl.SetOic(oic)
//		ls := cl.OutputCheck(L)
//	}


}
func (self *Clus) AppendVal(v *tmpdata.Val,isBack int)  {

	clus,tmps := self.getTmpVal()

	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
	tmpDisSort:=make( []*Distance,Long)
	sortlist:=make( []int,Long)
	var tmpClu []*Cl
	var vals []*tmpdata.Val
	var firstCl *Cl
	isB := 0
	for j,d := range dis[:Long] {
//		fmt.Println(i,d.i)
		oic := clus[d.i]
		cl,_:=oic.TmpAppendVal(v)
		cl.SetOic(oic)
		if firstCl == nil && vals == nil {
			tmpClu = append(tmpClu,cl)
//			fmt.Println(tmpClu)
		}
		L := len(oic.RawPatterns)
		tmpDisSort[j] = cl.DisSort[L][0]
		sortlist[j] = d.i
		appendSort(tmpDisSort,sortlist,j)
		//appendDis(tmpDisSort,j)
		var delist []int
//		var vs []*tmpdata.Val
		for _,d := range cl.DisSort[L] {
			if d.dis < curves.Errs {
//				vs = append(vs,d.a)
				delist= AppendSort(delist,d.i)
			}
		}
		if delist == nil {
			if isB <3 {
				isB ++
				continue
			}else{
				break
			}

		}
		if float64(len(delist))/float64(L) >= 0.5 {
		//	cl.Clear()
			if firstCl == nil {
				firstCl = cl
			}else{
				vals = append(vals,cl.RawPatterns...)
				oic.Clear()
				cl = nil
			}
		}else{
//			fmt.Println(delist,L)
			for _,d_i := range delist {
//				fmt.Printf("%d \r\n",d_i)
				_v,err := oic.DeleteVal(d_i)
				if err != nil {
					fmt.Println(delist)
					fmt.Println(oic.RawPatterns)
					fmt.Println(cl.RawPatterns)
					fmt.Printf("%p %p\r\n",_v,v)
					panic(err)
				}
				vals = append(vals,_v)
			}
			oic.UpdateCore()
			cl = nil
		}
	}
	if isBack == 1 {
		minCl := clus[sortlist[0]]
		kv := float64(minCl.CountY[0])/float64(minCl.CountY[1])
		if kv > 1.5 || kv < 0.5 {
			Y := 0
			if minCl.CountY[0] < minCl.CountY[1] {
				Y =1
			}
//		if float64(minCl.CountY[0])/float64(minCl.CountY[1]) < 0.5 {
			di := tmpDisSort[0]
			if di.a.Y != Y {
		if di.dis < curves.Errs {
//			self.ca.same1++
//			if di.a.C != v.C {
//				self.ca.er1 ++
//			}else{
				self.ca.same++
				if di.a.Y != v.Y {
					self.ca.er ++
				}
			}
			}
		}
	}
	if firstCl != nil {
		for _,_v := range vals{
			firstCl.Append(_v)
		}
		firstCl.UpdateCore()
		firstCl.oic.Copy(firstCl)
	}else if vals != nil{
		firstCl = new(Cl)
		firstCl.Append(v)
		for _,_v := range vals{
			firstCl.Append(_v)
		}
		firstCl.UpdateCore()
		self.Clu = append(self.Clu,firstCl)
	}else{
		firstCl = tmpClu[0]
		firstCl.UpdateCore()
		firstCl.oic.Copy(firstCl)
	}

}
func (self *Clus) AppendValBak(v *tmpdata.Val,isBack int)  {

	clus,tmps := self.getTmpVal()

	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
	tmpDisSort:=make( []*Distance,Long)
	var tmpc []*Cl = make([]*Cl,Long)
	var tmpli [][]int = make([][]int,Long)
	var minT int = -1
	var ls int
	sortlist:=make( []int,Long)
	out := 0
	for j,d := range dis[:Long] {
		cl,L:=clus[d.i].TmpAppendVal(v)
		cl.SetOic(clus[d.i])
		tmpc[j] = cl
		tmpli[j] = L
		tmpDisSort[j] = cl.DisSort[len(cl.DisSort)-1][0]
		sortlist[j] = j
		ls = appendSort(tmpDisSort,sortlist,j)
		if ls == 0 {
			out = 0
			minT = j
		}else{
			out ++
			if out >10 {
				Long = j+1
				break
			}
		}
	}
	var isdiff *Distance
//	fmt.Println(tmpDisSort[0])
//	fmt.Printf("%p\r\n",v)
	if isBack == 1 {
		self.ca.same++
		if tmpDisSort[0].a.C != v.C {
			self.ca.er ++
		}
		if tmpDisSort[0].dis < curves.Errs {
			self.ca.same1 ++
			if tmpDisSort[0].a.C != v.C {
				self.ca.er1 ++
			}
		}
	}
	if tmpDisSort[0].a.C != v.C {
		isdiff = tmpDisSort[0]
//		newClu := new(Cl)
//		newClu.Append(v)
//		self.Clu = append(self.Clu,newClu)
//		minT = -1
	}else{
	//	fmt.Println(tmpDisSort[0].a)
	//	fmt.Println(v)
		minT = sortlist[0]
		cp :=tmpc[minT]

		isFind:= false
		for _,val := range cp.RawPatterns {
			if val == tmpDisSort[0].a {
				isFind = true
				break
			}
		}
		if !isFind {
			fmt.Println("find not")
			panic(99)
		}

		cp.UpdateCore()
		cp.oic.Copy(cp)
//		clus[dis[minT].i].Copy(cp)
	}
	if isdiff != nil {
		for j,d := range tmpDisSort[1:Long] {
			if d.a.C == v.C {
				testd := new(Distance)
				testd.Init(d.a,isdiff.a,0)
				if testd.dis > d.dis {
					minT = sortlist[1:Long][j]
					cp :=tmpc[minT]

					isFind:= false
					for _,val := range cp.RawPatterns {
						if val == d.a {
							isFind = true
							break
						}
					}
					if !isFind {
						fmt.Println("find not")
						panic(99)
					}

					cp.UpdateCore()
					cp.oic.Copy(cp)
//					clus[dis[minT].i].Copy(cp)
					isdiff = nil
					break
				}
			}
		}
		if isdiff != nil {
			newClu := new(Cl)
			newClu.Append(v)
			self.Clu = append(self.Clu,newClu)
			minT = -1
		}
	}
//	return
	if isBack > LongIsBack {
		return
	}
	var vals []*tmpdata.Val
	var outCount int
	for j,d := range dis[:Long] {
		if j == minT {
			continue
		}
		ls := tmpc[j].OutputCheck(tmpli[j])
		if ls == nil {
			if outCount <3{
				outCount++
				continue
			}else{
				break
			}
		}else{
			outCount = 0
		}
		cp := clus[d.i]
//		fmt.Println(ls,len(cp.RawPatterns))
		if len(cp.RawPatterns)-len(ls) < 2 {
			vals = append(vals,cp.RawPatterns...)
			cp.Clear()
		}else{
//			fmt.Println(ls)
//			fmt.Println(tmpc[j])
//			fmt.Println(cp)
			for _,_i := range ls {
//				val:=cp.RawPatterns[_i]
				vals = append(vals, cp.RawPatterns[_i])
				cp.DeleteVal(_i)
			}
			cp.UpdateCore()
		}
		if clus[d.i] != cp {
			panic(999)
		}
		clus[d.i] = cp
		if len(vals) >10{
			break
		}
	}
//	fmt.Println(len(vals),isBack)
	if vals == nil {
		return
	}
	for _,val := range vals {
		self.AppendVal(val,isBack+1)
	}
//	walk := new(sync.WaitGroup)
//	walk.Add(1)
//	go self.syncAppendVal(vals,isBack,walk)
//	walk.Wait()

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
//	o := L - MaxLong
//	if o<0 {
//		o = 0
//	}
	o := 0
	for i := L -1;i>=o;i--{
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
