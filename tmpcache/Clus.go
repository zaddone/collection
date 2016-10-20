package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
//	"github.com/zaddone/collection/curves"
	"encoding/json"
//	"sync"
	"fmt"
	"math"
	"time"
)
const (
	LongLen int = 1000
	MaxVal  int = 10000
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
func (self *Clus) FindSame(v *tmpdata.Val) *tmpdata.Val {
	clus,tmps := self.getTmpVals(nil)
//	fmt.Println(len(clus))
	if len(clus) < 2 {
		return nil
	}
	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
	outList := make([][]int,Long)
	tmpDisSort:=make( []*Distance,Long)
	sortlist:=make( []int,Long)
	tmpClus :=make( []*Cl,Long)
	isOut := 0
	H:
	for j,d := range dis[:Long] {
		oic := clus[d.i]
		if oic.lock {
			for {
				if oic.RawPatterns == nil {
					continue H
				}
				if !oic.lock{
					break
				}
				time.Sleep(100*time.Millisecond)
			}
		}
		cl,L,num:=oic.TmpAppendVal(v)
		cl.SetOic(oic)
		out := cl.OutputCheck(L,num)
		if out == nil {
			if isOut >2 {
				Long = j
				break
			}
			isOut ++
		}else{
			isOut = 0
		}

		tmpClus[j] = cl
		outList[j] = out
		tmpDisSort[j] = cl.DisSort[len(cl.DisSort)-1][0]
		sortlist[j] = j
		appendSort(tmpDisSort,sortlist,j)
	}
	var cls *Cl
	if 0 == sortlist[0] {
		return tmpDisSort[0].a
		cls =tmpClus[0].oic
//		cls =tmpClus[0].oic
//		cls.UpdateCore()
//		cls.oic.Copy(cls)

		a:=float64(cls.CountY[0])
		b:=float64(cls.CountY[1])
		sum := a+b
//		if sum >10 {
			df := a-b
			if math.Abs(df)/sum > 0.75 {
				y := 0
				if df <0 {
					y = 1
				}
				val :=tmpDisSort[0].a
				if val.Y == y {
					return val
				}

			}
//		}
	}
//	return -1,fmt.Errorf("is nil")
	return nil

}
func (self *Clus) AppendVal(v *tmpdata.Val,isBack int)  {

	clus,tmps := self.getTmpVal()
	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
	outList := make([][]int,Long)
	tmpDisSort:=make( []*Distance,Long)
	sortlist:=make( []int,Long)
	tmpClus :=make( []*Cl,Long)
	isOut := 0
	for j,d := range dis[:Long] {
		oic := clus[d.i]
		cl,L,num:=oic.TmpAppendVal(v)
		cl.SetOic(oic)
		out := cl.OutputCheck(L,num)
		if out == nil {
			if isOut >2 {
				Long = j
				break
			}
			isOut ++
		}else{
			isOut = 0
		}

		tmpClus[j] = cl
		outList[j] = out
		tmpDisSort[j] = cl.DisSort[len(cl.DisSort)-1][0]
		sortlist[j] = j
		appendSort(tmpDisSort,sortlist,j)
	}
	var cls *Cl
	if 0 == sortlist[0] {
		cls =tmpClus[0]
		cls.UpdateCore()

		if cls.oic.GetPG() < 1 {
		sameVal := tmpDisSort[0].a
		if sameVal.Y == 0 && sameVal.C == v.C {
			self.ca.same++
			if sameVal.Y != v.Y {
				self.ca.er ++
			}
		}
		}

	//	if cls.oic.GetPG() < 1 {
	//		a:=float64(cls.oic.CountY[0])
	//		b:=float64(cls.oic.CountY[1])
	//		sum := a+b
//	//		if sum >10 {
	//		df := a-b
	//		if math.Abs(df)/sum > 0.65 {
	//			y := 0
	//			if df <0 {
	//				y = 1
	//			}
	//			if tmpDisSort[0].a.Y == y {
	//				self.ca.same++
//	//				fmt.Println(cls.CountY,cls.oic.CountY)
	//				if y != v.Y {
	//					self.ca.er ++
	//				}
	//			}

	//		}
	//	}
//		}
		cls.oic.Copy(cls)
	}else{
		cls = new(Cl)
		cls.Append(v)
		self.Clu = append(self.Clu,cls)
		for _,J := range sortlist[:Long]{
//			J :=sortlist[0]
			lis := outList[J]
			if lis == nil {
				continue
			}
			fc := tmpClus[J].oic
			fc.lock = true
//			fmt.Println(lis,len(fc.RawPatterns),J)
			Lf := len(fc.RawPatterns)
			if float64(len(lis))/float64(Lf) > 0.6 {
				for i:= Lf-1;i>=0;i--{
//				for _,i := range lis {
					err := self.AppendEdit(fc.RawPatterns[i],fc)
					if err == nil {
						_,err = fc.DeleteVal(i)
						if err != nil {
							panic(err)
						}
						fc.UpdateCore()
					}
				}
			}else{
				for _,i := range lis {
					err := self.AppendEdit(fc.RawPatterns[i],fc)
					if err == nil {
						_,err = fc.DeleteVal(i)
						if err != nil {
							panic(err)
						}
						fc.UpdateCore()
					}
				}

			}
			fc.lock = false
//			fc.UpdateCore()
//			if len(fc.RawPatterns) == 1 {
//				err := self.AppendEdit(fc.RawPatterns[0],fc)
//				if err == nil {
//					fc.Clear()
//				}
//			}
		}
//		cls.UpdateCore()
		if 1 == len(cls.RawPatterns) {
			err := self.AppendEdit(cls.RawPatterns[0],cls)
			if err == nil {
//				fmt.Println(len(self.Clu))
				self.Clu = self.Clu[:len(self.Clu)-1]
//				fmt.Println("is add ok")
//				fmt.Println(len(self.Clu))
			}else{
				fmt.Printf("count: %d %d\r",len(self.Clu),self.ValCount)
			}
		}
	}

}
func (self *Clus) AppendEdit(v *tmpdata.Val,clu *Cl) error {

	clus,tmps := self.getTmpVals(clu)
	if len(clus) < 3 {
		return fmt.Errorf("len is m")
	}
	dis := (&Cl{RawPatterns:tmps}).FindSortVal(v)
	Long := len(dis)
	if Long > LongLen {
		Long = LongLen
	}
//	outList := make([][]int,Long)
	tmpDisSort:=make( []*Distance,Long)
	sortlist:=make( []int,Long)
	tmpClus :=make( []*Cl,Long)
	isOut := 0
	for j,d := range dis[:Long] {
		oic := clus[d.i]
		cl,L,num:=oic.TmpAppendVal(v)
		cl.SetOic(oic)
		out := cl.OutputCheck(L,num)
		if out == nil {
			if isOut >2 {
				Long = j
				break
			}
			isOut ++
		}else{
			isOut = 0
		}

		tmpClus[j] = cl
//		outList[j] = out
		tmpDisSort[j] = cl.DisSort[len(cl.DisSort)-1][0]
		sortlist[j] = j
		appendSort(tmpDisSort,sortlist,j)
	}
	var cls *Cl
	if 0 == sortlist[0] {
		cls =tmpClus[0]
		cls.UpdateCore()
		cls.oic.Copy(cls)
		return nil
	}
	return fmt.Errorf("is err")

}
func (self *Clus) getTmpVals(clu *Cl) (tmpc []*Cl,tmp []*tmpdata.Val) {
	L := len(self.Clu)
	tmpc = make([]*Cl,L)
	tmp = make([]*tmpdata.Val,L)
	j := 0
	for _,_c := range self.Clu{
//	for i := L -1; i>=0; i--{
//		_c := self.Clu[i]
//		le := len(_c.RawPatterns)

		if _c.lock {
			continue
		}

		if _c.RawPatterns == nil {
//			self.Clu = append(self.Clu[:i],self.Clu[i+1:]...)
			continue
		}
		if _c == clu {
			continue
		}
		if len(_c.RawPatterns) > _c.Core {
			tmpc[j] = _c
			tmp[j] = _c.RawPatterns[_c.Core]
			j ++
		}
	}
	return tmpc[:j],tmp[:j]
}
func (self *Clus) getTmpVal() (tmpc []*Cl,tmp []*tmpdata.Val) {

	L := len(self.Clu)
	tmpc = make([]*Cl,L)
	tmp = make([]*tmpdata.Val,L)
//	hj := self.ValCount - MaxVal
	j := 0
	for i := L -1;i>=0;i--{
		_c := self.Clu[i]
		le := len(_c.RawPatterns)
		if le < 10 {
			_c.lock = true
			if le == 0 {
				self.Clu = append(self.Clu[:i],self.Clu[i+1:]...)
				continue
			}
			for j := le-1;j>=0;j -- {
				err := self.AppendEdit(_c.RawPatterns[j],_c)
				if err == nil {
					if len(_c.RawPatterns) >1 {
						_,err = _c.DeleteVal(j)
						if err != nil {
							panic(err)
						}
					}else{
						_c.Clear()
						break
					}
				}
			}
			if len(_c.RawPatterns) == 0 {
				self.Clu = append(self.Clu[:i],self.Clu[i+1:]...)
				continue
			}
//			if hj > 0 {
//				last := _c.RawPatterns[le-1]
//				if last.GetH() < hj {
//					self.Clu = append(self.Clu[:i],self.Clu[i+1:]...)
//					continue
//				}
//			}
			_c.UpdateCore()
			_c.lock = false
		}
		if len(_c.RawPatterns) > _c.Core {
			tmpc[j] = _c
			tmp[j] = _c.RawPatterns[_c.Core]
			j ++
		}
	}
	return tmpc[:j],tmp[:j]

}
