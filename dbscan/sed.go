package dbscan
import (
	"math"
	"fmt"
	"github.com/zaddone/collection/tmpdata"
	"github.com/zaddone/collection/common"
	nn "github.com/zaddone/collection/bpnn"
	"encoding/json"
//	"github.com/zaddone/server/roam"
//	"time"
//	"crypto/sha1"
)
const (
//	ORDER int = 1
	Level int = 7
)
type Ds struct {
	Num float64
	Sum float64
	Fsum float64
	Weight float64
}
func (self *Ds) SetWeight(Xin float64) {
	self.Num ++
	self.Sum += Xin
	self.Fsum+= Xin*Xin
	self.Weight =  math.Sqrt(self.Fsum/self.Num - (self.Sum*self.Sum)/(self.Num*self.Num))
}

func (self *Ds) GetScore (x float64) float64 {

	return (x-self.Sum/self.Num)/self.Weight

}

func GetDss(vs []*tmpdata.Val) (Dss []*Ds) {
	for i,v := range vs {
		if i == 0 {
			Dss = make([]*Ds,len(v.X))
			for j,_ :=range Dss {
				Dss[j] = new(Ds)
			}
		}
		for j,x := range v.X {
//			if j == 0 {
//				continue
//			}
			Dss[j].SetWeight(x)
		}
	}
	return Dss
}
//func Exchange(vs []*tmpdata.Val) []*Ds{
//
//	Dss := GetDss(vs)
//	for i,v := range vs {
//		for j,x := range v.X {
//			if j == 0 {
//				continue
//			}
//			v.X[j] = Dss[j-1].GetScore(x)
//			fmt.Println(x,vs[i].X[j])
//		}
//	}
//	return Dss
//
//}
func EucDistance(a []float64,b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	sum := 0.0
	for i,_a := range a {
		sum += math.Pow(_a-b[i],2)
	}
	return math.Sqrt(sum)
}
type Distance struct {
	dis float64
	a *tmpdata.Val
	b *tmpdata.Val
}
func (self *Distance) Init (a,b *tmpdata.Val) {
	self.dis = EucDistance(a.X,b.X)
	self.a = a
	self.b = b
}
func (self *Distance) Find(a *tmpdata.Val) bool {
	if self.a == a || self.b == a{
		return true
	}
	return false
}
func sortDis(dis []*Distance,d *Distance) []*Distance {
	Ls:=len(dis)
	L := Ls-1
	dis  = append(dis,d)
	for i := L;i>=0 ;i-- {
		if dis[i].dis > d.dis {
			dis[i],dis[Ls] = dis[Ls],dis[i]
			Ls = i
		}
	}
	return dis
}
type SED struct {
	patterns []*tmpdata.Val
	patternsbak []*tmpdata.Val
	RawPatterns []*tmpdata.Val
	Dss   []*Ds
//	Ready  bool

	Nn     *nn.InputLevel
	Nn1     *nn.InputLevel

	Change int
	Count  float64
	Errs   float64
	State  chan bool
	Counts []float64
	CountsY []float64
	Dis  []*Distance
	Hot  int
}
func (self *SED) UpdateCounts() {
	self.Counts = make([]float64,2)
	self.CountsY = make([]float64,3)
	for _,R := range self.RawPatterns {
		self.Counts[R.C] ++
		self.CountsY[R.Y] ++
	}
}
func (self *SED) GetDiffCVal() []*tmpdata.Val {

	if self.Equs1() == -1 {
		return nil
	}
	cent := self.GetMin()
	var out  *tmpdata.Val
	var tmpPat []*tmpdata.Val
	var I int
	for i,v := range self.RawPatterns[1:] {
		if cent.Y != v.Y  || cent.C != v.C{
			out = v
			tmpPat = self.RawPatterns[i+1 :]
			I = i
			break
		}
	}
	if out == nil {
		return nil
	}
	var outP []*tmpdata.Val = []*tmpdata.Val{out}
	mapPat := make(map[*tmpdata.Val][2]float64)
	for _,d := range self.Dis {
		if d.a == out {
			pat,ok := mapPat[d.b]
			if !ok {
				pat = [2]float64{0,0}
				mapPat[d.a] = pat
			}
			for _,p := range tmpPat {
				if d.b == p {
					pat[0] = d.dis
					break
				}
			}
		}else if d.b == out {
			pat,ok := mapPat[d.a]
			if !ok {
				pat = [2]float64{0,0}
				mapPat[d.a] = pat
			}
			for _,p := range tmpPat {
				if d.a == p {
					pat[0] = d.dis
					break
				}
			}
		}else if d.a == cent {
			pat,ok := mapPat[d.b]
			if !ok {
				pat = [2]float64{0,0}
				mapPat[d.a] = pat
			}
			for _,p := range tmpPat {
				if d.b == p {
					pat[1] = d.dis
					break
				}
			}
		}else if d.b == cent {
			pat,ok := mapPat[d.a]
			if !ok {
				pat = [2]float64{0,0}
				mapPat[d.a] = pat
			}
			for _,p := range tmpPat {
				if d.a == p {
					pat[1] = d.dis
					break
				}
			}
		}
	}
	var outPat []*ValSort
	for k,v := range mapPat {
		if v[0] < v[1] {
			outP = append(outP,k)
			outPat = appendSortVal(outPat,&ValSort{val:k,long:v[1]})
		}
	}
	self.RawPatterns = append(self.RawPatterns[:I],self.RawPatterns[I+1 :]...)
	if len(outPat) == 0 {
		return outP
	}
	j :=0
	for i:= len(self.RawPatterns)-1; i>=0; i-- {
		if self.RawPatterns[i] == outPat[j].val {
			self.RawPatterns = append(self.RawPatterns[:i],self.RawPatterns[i+1 :]...)
			j ++
		}

	}
//	if len(self.RawPatterns) >1{
//		self.Refresh()
//	}
	if j != len(outPat){
		panic(j)
	}

	return outP

}
func appendSortVal(pat []*ValSort,p *ValSort) []*ValSort {
	L := len(pat)
	if L == 0 {
		pat = []*ValSort{p}
		return pat
	}
	pat = append(pat,p)
	swap(L,pat)
	return pat
}
func swap(num int,Node []*ValSort){
	if num == 0 {
		return
	}
	Next := num -1
	if Node[num].long < Node[Next].long {
		return
	}
	Node[num], Node[Next] = Node[Next],Node[num]
	swap(Next,Node)
}
type ValSort struct {
	val *tmpdata.Val
	long float64
}
func (self SED) GetSortValToDis(v *tmpdata.Val) (dis []*Distance) {
	L := len(self.RawPatterns)
	self.RawPatterns = append(self.RawPatterns,v)
	self.Dss = GetDss(self.RawPatterns)
	self.UpdatePatterns()
	tv := self.patterns[L]
//	var dis []*Distance
	for i,pa := range self.patterns {
		if i == L {
			continue
		}
		d := new(Distance)
		d.Init(pa,tv)
		dis = sortDis(dis,d)
	}
	return dis
}
func (self SED) FindDis(v *tmpdata.Val) bool {

	M :=self.RawPatterns[0]
	if v.Y != M.Y || v.C != M.C {
		return false
	}
	L := len(self.RawPatterns)
	self.RawPatterns = append(self.RawPatterns,v)
	self.Dss = GetDss(self.RawPatterns)
	self.UpdatePatterns()
	tv := self.patterns[L]
	var dis []*Distance
	for i,pa := range self.patterns {
		if i == L {
			continue
		}
		d := new(Distance)
		d.Init(pa,tv)
		dis = sortDis(dis,d)
	}

//	dis := self.GetSortValToDis(v)

	Ldis := dis[0]
	Min :=Ldis.a
	if v.C != Min.C {
		return false
	}
	dis = nil
	for i,pa := range self.patterns {
		if i == L {
			continue
		}
		if pa == Min {
			continue
		}
		d := new(Distance)
		d.Init(pa,tv)
		dis = sortDis(dis,d)
	}
	sum:=0.0
	for _,d := range dis {
		sum += d.dis
	}
	return sum/float64(len(dis)) > Ldis.dis

}
func (self *SED) IsReadys() (isS [2]int,is bool) {
	isS[0]=self.Equs1()
	isS[1]=self.Equs2()
	for _,i := range isS {
		if i != -1 {
			is = true
			break
		}
	}
	return isS,is

}
//func (self *SED) IsReady() (bool) {
//	if self.Equs1() > -1 {
//		return true
//	}
//	L :=len(self.RawPatterns)
//	for _,c := range self.CountsY {
//		if c == L {
//			return true
//		}else if c >0 {
//			break
//		}
//	}
//	return false
//}
func (self *SED) Merges(Sed *SED) (old []*tmpdata.Val) {
	Sm :=self.GetMin()
	dis := Sed.GetSortValToDis(Sm)
//	var old []*tmpdata.Val
	isOut:= false
	for _,d := range dis {
		if !isOut  && self.FindDis(d.a)  {
			self.RawPatterns = append(self.RawPatterns,d.a)
//			self.UpdateNone(d.a)
		}else{
			isOut = true
			old = append(old,d.a)
		}
	}
	self.Refresh()
	return old
//	if old == nil {
//		return nil
//	}
//	Sed.Init(old)
//	return Sed
}
func (self *SED) Merge(Sed *SED) bool {
	isl := false
	for _,pa := range Sed.RawPatterns {
		if self.Check(pa) {
			isl = true
			break
		}
	}
	if !isl {
		return false
	}
	self.RawPatterns = append(self.RawPatterns,Sed.RawPatterns...)
	self.Refresh()
	return true
}
func (self *SED) Equs2() int {
	L :=float64(len(self.RawPatterns))
	for i,_k := range self.CountsY {
		if _k == L  {
			return i
		}else if _k > 0 {
			break
		}
	}
	return -1
}
func (self *SED) Equs1() int {
	for i,_k := range self.Counts {
		if _k == 0  {
			return i
		}
	}
	return -1
}
func (self *SED) Equs() []float64 {
	L := float64(len(self.RawPatterns))
	if L < 5 {
		return nil
	}
	self.Counts = make([]float64,2)
	self.CountsY = make([]float64,3)
	for  _,p := range self.RawPatterns {
		self.Counts[p.C] ++
		self.CountsY[p.Y]++
	}
	for _,_k := range self.Counts {
		if _k == 0 {
			return self.Counts
		}
	}
	return nil
}
func (self *SED) Equ1() ([]*tmpdata.Val) {
	L := len(self.RawPatterns)
	if L < 10 {
		return nil
	}
	if self.CountsY[2] == 0 || self.CountsY[2] == float64(L) {
		return nil
	}
	patBak := make([]*tmpdata.Val,L)
	copy(patBak,self.patterns)
	for _,pa := range patBak {
		if pa.Y > 0  {
			pa.Y --
		}
	}
	return patBak
}
func (self *SED) Equ2() ([]*tmpdata.Val) {
	L := float64(len(self.patternsbak))
	if L < 10 {
		return nil
	}
	if self.CountsY[0] == 0 || self.CountsY[0] == L {
		return nil
	}
	return self.patternsbak
}
func (self *SED) Equ() bool {
	L := float64(len(self.RawPatterns))
	if L < 5 {
		return false
	}
	ks :=L*0.1
//	k := make([]float64,3)
//	for  _,p := range self.RawPatterns {
//		k[p.Y] ++
//	}
	for _,_k := range self.Counts {
		if _k < ks {
			return true
		}
	}
	return false
}
func (self *SED) GetValData() ([]byte,error) {

//	pa := make([]tmpdata.Val,len(self.patterns))
//	for i,p := range self.patterns {
//		pa[i] = *p
//		pa[i].Y = p.C
//	}
//	fmt.Println(len(self.patterns))

	return json.Marshal(self.patterns)

//	d,err:=json.Marshal(self.patterns)
//	if err != nil {
//		panic(err)
//	}
//	return d

}
func (self *SED) Forecast(v *tmpdata.Val) int {
	if !self.Equ() {
		return -1
	}
	if !self.Check(v) {
		return -1
	}
	p := self.GetCountsP()
	if p == 2 {
		return -1
	}
//	return p
	X:=common.GetX(v.X)
	if self.Nn1 != nil {
		vs,err := self.Nn1.Update(X)
		if err == nil && len(vs) > 0 {
			f,err := common.GetForVal(vs)
			if err == nil {
				if f == 1 {
					return -1
				}
			}

		}
	}
	if self.Nn != nil{
		vs,err := self.Nn.Update(X)
		if err == nil && len(vs) > 0 {
			f,err := common.GetForVal(vs)
			if err == nil {
				if f != p{
					return -1
				}
			}
		}
	}
	return p
}
func (self *SED) GetCountsP() int {
	maxc := 0.0
	j := 0
	for i,c := range self.CountsY {
		if c > maxc {
			maxc = c
			j = i
		}
	}
	return j

}
func (self *SED) NnCheck(v *tmpdata.Val) bool {
	if !self.Nn1Test(v) {
		return false
	}

	if v.Y == 2 {
		return true
	}else{
		return self.NnTest(v)
	}
	return false
}
func (self *SED) Nn1Test(v *tmpdata.Val) bool {
	if self.Nn1 == nil {
//		self.ClearNnTest()
		if v.Y != 2 {
			if self.CountsY[2] == 0 {
				return true
			}
		}else{
			if int(self.CountsY[2]) == len(self.RawPatterns) {

				return true
			}
		}
		return false
	}
	vs,err := self.Nn1.Update(common.GetX(v.X))
	if err != nil {
		return false
	}
	if len(vs) == 0 {
		return false
	}
	self.Count++
	f,err := common.GetForVal(vs)
	if err != nil {
		return false
	}
	if f>0 != ((v.Y-1) >0) {
		self.Errs ++
		return false
	}
	return true
}
func (self *SED) NnTest(v *tmpdata.Val) bool {
	if self.Nn == nil {
//		self.ClearNnTest()
		if v.Y == 1 && int(self.CountsY[1]) == len(self.patternsbak) {
			return true
		}else if v.Y == 0 && int(self.CountsY[0]) == len(self.patternsbak) {
			return true
		}
		return false
	}
	vs,err := self.Nn.Update(common.GetX(v.X))
	if err != nil {
		return false
	}
	if len(vs) == 0 {
		return false
	}
	f,err :=common.GetForVal(vs)
	if err != nil {
		return false
	}
	if f != v.Y {
		return false
	}

	return true

}
func (self *SED) ClearNnTest() {
	self.Count = 0
	self.Errs = 0
}
func (self *SED) ShowNnTest() float64 {

	return self.Errs/self.Count
}
func (self *SED) GetTest() string {
	o:=0.0
	e:=0.0
	for _,p := range self.RawPatterns {
		if p.C > 0 {
			o ++
		}else{
			e ++
		}
	}
	L :=float64(len(self.RawPatterns))
	return fmt.Sprintf("%d %.5f %.5f",int(L),o/L,e/L)
}
func (self *SED) ShowTest() {
	o:=0.0
	e:=0.0
	for _,p := range self.RawPatterns {
		if p.Y > 0 {
			o ++
		}else{
			e ++
		}
	}
	L :=float64(len(self.RawPatterns))
	fmt.Printf("%d %d %d %.5f %.5f\r\n",int(L),int(o),int(e),o/L,e/L)
}
func (self *SED) GetMax() *tmpdata.Val {
	return self.RawPatterns[len(self.RawPatterns)-1]
}
func (self *SED) GetMin() *tmpdata.Val {
	return self.RawPatterns[0]
}
func (self *SED) GetKey() []byte {
	return self.GetMin().GetKey()
}
func (self *SED) Init(pa []*tmpdata.Val){
	if len(pa)> 1 {
		self.SetPatterns(pa)
		self.GetCore(self.GetDistance(0,nil,nil))
//	self.Nn = new(nn.InputLevel)
	}else{
		self.RawPatterns = pa
	}
//	self.Ready = false
	self.State = make(chan bool)

}
func (self *SED) InitNn(d []byte,isNn bool) {
//	fmt.Println(string(d))
	tmp := new(nn.TmpData)
//	var tmps [2]*nn.TmpData
	err := json.Unmarshal(d,tmp)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isNn {
		self.Nn = new(nn.InputLevel)
		err = self.Nn.SetTmp(tmp,true)
		if err != nil {
			fmt.Println(err)
			self.Nn = nil
		}
//		fmt.Println(0,self.Nn)
	}else{
		self.Nn1 = new(nn.InputLevel)
		err = self.Nn1.SetTmp(tmp,true)
		if err != nil {
			fmt.Println(err)
			self.Nn1 = nil
		}
//		fmt.Println(1,self.Nn1)
	}
}
func (self *SED) SetPatterns (patterns []*tmpdata.Val){

	self.Dss = GetDss(patterns)
	self.RawPatterns = patterns
	self.UpdatePatterns()

}
func (self SED) Check(v *tmpdata.Val) bool {
	Pa:= []*tmpdata.Val{self.RawPatterns[0],self.RawPatterns[len(self.RawPatterns)-1],v}
	self.SetPatterns(Pa)
	d1 := new(Distance)
	d1.Init(self.patterns[0],self.patterns[1])

	d2 := new(Distance)
	d2.Init(self.patterns[0],self.patterns[2])
	if d2.dis < d1.dis {
	//	fmt.Println(d1.dis,d2.dis)
		return true
	}
	return false
}
func (self *SED) UpdateNone(v *tmpdata.Val) {

	self.UpdateDss(v)
	self.RawPatterns = append(self.RawPatterns,v)
//	self.Counts[v.C]++
//	self.CountsY[v.Y]++
	self.UpdatePatterns()
	self.GetCore(self.GetDistance(0,nil,nil))
	self.Change ++

	close(self.State)
	self.State = make(chan bool)

}
func (self *SED) FindPatterns(t *tmpdata.Val) int {
	for i,v := range self.RawPatterns {
		if v == t {
			return i
		}
	}
	return -1
}
func (self *SED) RmPatterns(t *tmpdata.Val) *tmpdata.Val{
	i := self.FindPatterns(t)
	if i<0 {
		return nil
	}
	self.RawPatterns = append(self.RawPatterns[:i],self.RawPatterns[i+1:]...)
	self.Refresh()
	return self.GetMin()
}
func (self *SED) Refresh(){
	self.Dss = GetDss(self.RawPatterns)
	self.UpdatePatterns()
	self.GetCore(self.GetDistance(0,nil,nil))

	close(self.State)
	self.State = make(chan bool)
}
//func (self *SED) Updates(v *tmpdata.Val) (ts []*tmpdata.Val) {
//
//	self.UpdateDss(v)
//	self.RawPatterns = append(self.RawPatterns,v)
//	self.UpdatePatterns()
//	dis :=self.GetDistance(0,nil,nil)
//	self.GetCore(dis)
//	
//}
func (self *SED) Update(v *tmpdata.Val) (t *tmpdata.Val) {
	self.UpdateNone(v)
	L:=len(self.RawPatterns)-1
	t = self.RawPatterns[L]
	self.RawPatterns = self.RawPatterns[:L]
	self.syncGetDss()
	return t
}
func (self *SED) OutputEnd() (t *tmpdata.Val){
//	L:=len(self.RawPatterns)-1
	return self.RawPatterns[len(self.RawPatterns)-1]
}
func (self *SED) DeleteVal(t *tmpdata.Val) bool {
	for i := len(self.RawPatterns)-1;i>=0;i-- {
		if self.RawPatterns[i] == t {
			self.RawPatterns = append(self.RawPatterns[:i],self.RawPatterns[i+1:]...)
			self.syncGetDss()
			return true
		}
	}
	return false
}

func (self *SED) Output() (t *tmpdata.Val){
	L:=len(self.RawPatterns)-1
	t = self.RawPatterns[L]
	self.RawPatterns = self.RawPatterns[:L]
	self.syncGetDss()
	return t
}
func (self *SED) syncGetDss() {
	self.Dss = GetDss(self.RawPatterns)
	self.UpdatePatterns()
	self.GetCore(self.GetDistance(0,nil,nil))

}
func (self *SED) GetSortInt(dis []*Distance) (out []int) {
	out = make([]int,len(dis))
	for i,d := range dis {
		_,out[i] = self.FindRaw(d.a)
	}
	return out
}
func (self *SED) GetSortDis(t int) (dis []*Distance) {
	pt := self.patterns[t]
//	var dis []*Distance
	for i,p := range self.patterns {
		if i == t {
			continue
		}else{
			d := new(Distance)
			d.Init(p,pt)
			dis = sortDis(dis,d)
			if len(dis) >10000 {
				break
			}
		}
	}
//	fmt.Println(len(dis))
	return dis
}
func (self *SED) FindMinVal(t int) (*tmpdata.Val,int) {

	dis:= self.GetSortDis(t)
	for _,d := range dis {
		if d.a.GetD()!= 0{
			continue
		}else{
			return self.FindRaw(d.a)
		}
	}
	return self.FindRaw(dis[0].a)

}

func (self *SED) FindRaw(c *tmpdata.Val) (*tmpdata.Val,int) {
	for i,p := range self.patterns {
		if p == c {
			return self.RawPatterns[i],i
		}
	}
	return nil,0
}

func FindSame(pas []*tmpdata.Val,pa *tmpdata.Val)([]*tmpdata.Val,bool) {
	for _,p := range pas {
		if p == pa {
			return pas,false
		}
	}
	return append(pas,pa),true
}
func (self *SED) GetCore(dis []*Distance){

	L := len(dis)
	Ls := len(self.patterns)
	var pas []*tmpdata.Val
	var b bool
	i:= L-1
	for ;i>=0;i-- {
		pas,b=FindSame(pas,dis[i].a)
		if b  {
			if len(pas) == Ls {
				break
			}
		}
		pas,b=FindSame(pas,dis[i].b)
		if b  {
			if len(pas) == Ls {
				break
			}
		}
	}

//	for i,p := range pas {
//		fmt.Println(i,p)
//	}
//	fmt.Println(Ls,len(pas))
	if len(pas) != Ls {
		panic(0)
	}
	self.SortPatterns(pas[Ls-1],dis)
//	self.Max = dis[i].dis

}
func (self *SED) SortPatterns(p *tmpdata.Val,dis []*Distance) {
	Np ,_ :=self.FindRaw(p)
	NewP := make([]*tmpdata.Val,len(self.patterns))
	NewP[0] = Np
	i:=1
//	NewP := []*tmpdata.Val{Np}
	for _,d := range dis {
		if d.a==p {
			Np,_ :=self.FindRaw(d.b)
			NewP[i] = Np
			i++
		}else if d.b == p {
			Np,_ :=self.FindRaw(d.a)
			NewP[i] = Np
			i++
		}
	}
	if i != len(NewP) {
		fmt.Println(i,len(NewP),len(self.patterns),len(self.RawPatterns))
		panic( "i != len %d %d")
	}
	self.RawPatterns = NewP

}

func (self *SED) GetDistance(t int,dis []*Distance,pa *tmpdata.Val) ([]*Distance) {
	if t == 0 {
		return self.GetDistance(1,dis,self.patterns[0])
	}
	var lastPa *tmpdata.Val
	for _,pat := range self.patterns[t:] {
		if lastPa == nil {
			lastPa = pat
		}
		d := new(Distance)
		d.Init(pa,pat)
		dis = sortDis(dis,d)
	}
	t++
	if t > len(self.patterns)-1 {
		self.Dis = dis
		return dis
	}
	return self.GetDistance(t,dis,lastPa)
}

func (self *SED) UpdateDss(v *tmpdata.Val) {
	for j,x := range v.X {
	//	self.Dss[i] 
//		if j == 0 {
//			continue
//		}
		self.Dss[j].SetWeight(x)
	}
}
func (self *SED) UpdatePatterns() {
	self.patterns = make([]*tmpdata.Val,len(self.RawPatterns))
	self.Counts = make([]float64,2)
	self.CountsY = make([]float64,3)
	self.Hot ++
	self.patternsbak = nil
	for i,v := range self.RawPatterns {
		self.Counts[v.C] ++
		self.CountsY[v.Y]++
		P:=new(tmpdata.Val)
		P.X = make([]float64,len(v.X))
		P.Y = v.Y
		for j,x := range v.X {
		//	if j == 0 {
		//		continue
		//	}
			P.X[j] = self.Dss[j].GetScore(x)
		}
		self.patterns[i] = P
		if P.Y != 2 {
			self.patternsbak = append(self.patternsbak,P)
		}
	}
}
