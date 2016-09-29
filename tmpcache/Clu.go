package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	nn "github.com/zaddone/collection/bpnn"
	"github.com/zaddone/collection/common"
	"fmt"
	"encoding/json"
)
type Clu struct {
	Nn     *nn.InputLevel
	RawPatterns []*tmpdata.Val
//	Dss   []*Ds
	Counts [2]float64
	CountsY [3]float64
	Lock bool
	Core  *tmpdata.Val
//	CoreD  []*Distance

	start  chan bool
}
type TmpDa struct {
	Cl *Clu
	Nntmp *nn.TmpData
}
func InputData( d []byte) (*Clu,error) {
	tmp := new(TmpDa)
	err := json.Unmarshal(d,tmp)
	if err != nil {
		return nil, err
	}
	if tmp.Nntmp != nil {
		Nn := new(nn.InputLevel)
		err = Nn.SetTmp(tmp.Nntmp,true)
		if err != nil {
			return tmp.Cl,err
		}
		tmp.Cl.Nn = Nn
	}
	return tmp.Cl,nil
}
func (self *Clu) OutputData() ([]byte,error) {
	tmp := new(TmpDa)
	tmp.Cl = new(Clu)
	tmp.Cl.RawPatterns = self.RawPatterns
	tmp.Cl.Core = self.Core
	tmp.Cl.Counts = self.Counts
	tmp.Cl.CountsY = self.CountsY
	if self.Nn != nil {
		tmp.Nntmp = self.Nn.GetTmp(0,0,0)
	}
	return json.Marshal(tmp)
}
func (self Clu) GetCountsKey() (key [2]int) {
	for i,k := range self.Counts {
		if k == 0{
			key[i] = 0
		}else{
			key[i] = 1
		}
	}
	return key
}
func (self *Clu) NnCheck() ([]byte,chan bool) {
	if len(self.RawPatterns) <10 {
		return nil,nil
	}
	if !self.Check() {
		return nil,nil
	}
	return self.GetNnTrData(),self.start
}
func (self *Clu) Integ( clu *Clu) (isI bool) {
	isI = false
	if clu.Core != nil {
		p := self.RawPatterns[len(self.RawPatterns)-1]
		for i:=len(clu.RawPatterns)-1;i>=0;i-- {
//		for i,cl := range clu.RawPatterns {
			cl := clu.RawPatterns[i]
			if cl == clu.Core {
				continue
			}
			d1 := new(Distance)
			d1.Init(clu.Core,cl,i)

			d2 := new(Distance)
			d2.Init(p,cl,i)
			if d1.dis > d2.dis {
				self.Append(cl)
//				self.RawPatterns = append(self.RawPatterns,cl)
				clu.RawPatterns  = append(clu.RawPatterns[:i],clu.RawPatterns[i+1 :]...)
				isI = true
			}
		}
		if isI {
			clu.Update()
		}
	}
	return isI
}

func (self *Clu) Clear() {
	self.RawPatterns = nil
//	self.Dss = nil
	self.Counts = [2]float64{0,0}
	self.CountsY = [3]float64{0,0,0}
	self.Lock = true
}
func (self Clu) FindSortList(v *tmpdata.Val) (dis []*Distance) {
	L := len(self.RawPatterns)
	self.RawPatterns = append(self.RawPatterns,v)
	pat :=self.getScorePat()
	tv := pat[L]
	for i,pa := range pat {
		if i == L {
			continue
		}
		d := new(Distance)
		d.Init(pa,tv,i)
		dis,_ = sortDis(dis,d)
	}
	return dis
}
func (self Clu) FindSortVal(v *tmpdata.Val) (vs []int) {
	dis := self.FindSortList(v)
	vs= make([]int,len(dis))
	for j,d := range dis {
		vs[j] = d.i
	}
	return vs
}
func (self Clu) FindMinVal(v *tmpdata.Val) int {
	L := len(self.RawPatterns)
	self.RawPatterns = append(self.RawPatterns,v)
//	fmt.Println(len(self.RawPatterns))
//	self.updateDss(v)
	pat :=self.getScorePat()
	tv := pat[L]
	var dis []*Distance
	for i,pa := range pat {
		if i == L {
			continue
		}
		d := new(Distance)
//		d.Init(pa,tv,pa.GetD())
		d.Init(pa,tv,i)
		dis,_ = sortDis(dis,d)
	}
//	for j,d := range dis {
//		fmt.Println(j,d.i)
//	}
	return dis[0].i
}
func (self *Clu) getYVal() int {
	max:=0.0
	mi := -1
	for i,y := range self.CountsY {
		if max <= y {
			max = y
			mi = i
		}
	}
	return mi
}
func (self *Clu) InitNn(d []byte) {
	tmp := new(nn.TmpData)
	err := json.Unmarshal(d,tmp)
	if err != nil {
		panic(err)
		fmt.Println(err)
		return
	}
	Nn := new(nn.InputLevel)
	err = Nn.SetTmp(tmp,true)
	if err != nil {
		panic(err)
		fmt.Println(err)
	}else{
		self.Nn = Nn
	}
}

func (self *Clu) GetNnTrData() []byte {
	self.Update()
	k := self.getYVal()
	if k == 2 {
		return nil
	}
	tr := new(tmpdata.NnTrData)
	tr.Val = k

	tr.Data = self.getScorePat()
	if self.Nn != nil {
		d,err:=json.Marshal(self.Nn.GetTmp(0,0,0))
		if err != nil {
			return nil
		}
		tr.NnInfo = d
	}
	return tr.GetJsonData()
}
func (self Clu) Gaugedh(v *tmpdata.Val) (float64,float64) {
	L := len(self.RawPatterns)
	self.RawPatterns = append(self.RawPatterns,v)
//	self.updateDss(v)
	pat :=self.getScorePat()
	tv := pat[L]
	var dis []*Distance
	for i,pa := range pat {
		if i == L {
			continue
		}
		d := new(Distance)
		d.Init(pa,tv,i)
		dis,_ = sortDis(dis,d)
	}


	Ldis :=dis[0]
	Min :=Ldis.a
//	if v.C != Min.C {
//		return false
//	}
//	dis = nil
	sum := 0.0
	num := 0.0
	for i,pa := range pat {
		if i == L {
			continue
		}
		if pa == Min {
			continue
		}
		d := new(Distance)
		d.Init(pa,Min,i)
		sum += d.dis
		num ++
//		dis = sortDis(dis,d)
	}
	return sum/num , Ldis.dis
}
func (self Clu) Gauged(v *tmpdata.Val) bool {

	L := len(self.RawPatterns)
	self.RawPatterns = append(self.RawPatterns,v)
//	self.updateDss(v)
	pat :=self.getScorePat()
	tv := pat[L]
	var dis []*Distance
	for i,pa := range pat {
		if i == L {
			continue
		}
		d := new(Distance)
		d.Init(pa,tv,i)
		dis,_ = sortDis(dis,d)
	}


	Ldis :=dis[0]
	Min :=Ldis.a
	if v.C != Min.C {
		return false
	}
//	dis = nil
	sum := 0.0
	num := 0.0
	for i,pa := range pat {
		if i == L {
			continue
		}
		if pa == Min {
			continue
		}
		d := new(Distance)
		d.Init(pa,Min,i)
		sum += d.dis
		num ++
//		dis = sortDis(dis,d)
	}
	return sum/num > Ldis.dis

}
func (self *Clu) getDistance(pa []*tmpdata.Val) (dis [][]*Distance,mdis [][]*Distance) {
	L := len(pa)
	mdis = make([][]*Distance,L)
	for i,_ := range mdis {
		mdis[i] = make([]*Distance,L)
	}
	dis = make([][]*Distance,L)
	for i,P := range pa[:L-1]{
		for j,p := range pa[i+1:] {
			if P.GetD() == p.GetD() {
				panic(fmt.Sprint("test 1",i,j,"----"))
			}
			d := new(Distance)
			d.Init(P,p,i+j)
			dis[P.GetD()],_ = sortDis(dis[P.GetD()],d)
			dis[p.GetD()],_ = sortDis(dis[p.GetD()],d)
			mdis[P.GetD()][p.GetD()] = d
//			mdis[p.GetD()][P.GetD()] = d
//			fmt.Println(P.GetD(),p.GetD())
//			dis = sortDis(dis,d)
		}
	}
//	for i,md := range mdis {
//		fmt.Println(i,md)
//	}
//	panic("test ----")
	return dis,mdis
}
func (self *Clu) getScorePat() []*tmpdata.Val {
	pat := make([]*tmpdata.Val,len(self.RawPatterns))
//	copy(pat,self.RawPatterns)
	for i,v := range self.RawPatterns {
		pat[i] = v
		pat[i].SetD(i)

//		pat[i] = new(tmpdata.Val)
//		pat[i].X = make([]float64,len(v.X))
//		pat[i].Y = v.Y
//		pat[i].SetD(i)
//		for j,x := range v.X {
//			pat[i].X[j] = self.Dss[j].GetScore(x)
//		}
	}
	return pat
}
func (self *Clu) FindOutPat() ([]int )  {
	pat :=self.getScorePat()
//	Ls := len(pat)
	dis,mdis:=self.getDistance(pat)
//	var pas1 []*tmpdata.Val = make([]*tmpdata.Val,Ls)
//	var pas []*tmpdata.Val
	var maxdis []*Distance
	var ls,T int
	for t,d := range dis {
		if d == nil {
			panic(fmt.Sprintf("%d %d test",t,len(dis)))
		}
		maxdis,ls = sortDis(maxdis,d[len(d)-1])
		if ls == 0 {
			T = t
//			pat[t],pat[ls]= pat[ls],pat[t]
		}
	}
	co :=pat[T]
	self.Core = co
//	self.CoreD = dis[T]
//	c := co.GetD()
//	Lo := -1
//	var out *tmpdata.Val
	var oint int = -1
//	var outs []*tmpdata.Val
	var outint []int
	Lc:= len(self.RawPatterns)
	for _,p := range dis[T]{
		var hh *tmpdata.Val
		if p.a == co {
			hh = p.b
		}else{
			hh = p.a
		}
	//	fmt.Println(k,hh.GetD(),"___")
		if oint == -1 {
			if self.RawPatterns[hh.GetD()].C != self.RawPatterns[T].C {
				oint = hh.GetD()
			}
			if Lc > 50 {
				if self.RawPatterns[hh.GetD()].Y != self.RawPatterns[T].Y {
					oint = hh.GetD()
				}
			}
		}else{
			outint = append(outint,hh.GetD())
		}
	}
	if oint == -1 {
		return nil
	}
//	fmt.Println(T,oint,outint)
//	panic("test")
	var t1,t2 *Distance
	tmpOut:= []int{oint}
	tmpOutG:= []int{oint}
	for _,a := range outint {

		if a == oint {
			panic(fmt.Sprintf("%d %d %d ----",a,T,oint))
		}
		t1 = getDis(a,T,mdis)
		if self.RawPatterns[a].C  == self.RawPatterns[oint].C {
			tmpOutG = append(tmpOutG,a)
			tmpOut = appendSortInt(tmpOut,a)
			continue
		}
		if Lc > 50 {
			if self.RawPatterns[T].Y != self.RawPatterns[a].Y {
				tmpOutG = append(tmpOutG,a)
				tmpOut = appendSortInt(tmpOut,a)
				continue
			}
		}
		for _,b := range tmpOutG{
			t2 = getDis(a,b,mdis)
			if t1.dis > t2.dis{
				tmpOut = appendSortInt(tmpOut,a)
				break
			}
		}
	}
	if len(self.RawPatterns) - len(tmpOut) < 10 {
//		self.CoreD = nil
		self.Core = nil
	}
	return tmpOut
}
func getDis(a,b int, mdis [][]*Distance) (dis *Distance) {
	dis = mdis[a][b]
	if dis != nil {
		return dis
	}
	dis = mdis[b][a]
	if dis == nil {
		panic(fmt.Sprintf("__%d %d___",a,b))
	}
	return dis
}
func appendSortInt(li []int ,a int) []int {
	L := len(li)
	li = append(li,a)
	for i:= L -1 ;i >=0;i--{
//		l = li[i]
//	for i,l := range li[:L] {
		if li[i]<a {
			li[i],li[L] = li[L],li[i]
			L = i
		}else{
			break
		}
	}
	return li
}

func (self *Clu) Update()  {
//	self.Dss,self.Counts,self.CountsY = GetDss(self.RawPatterns)
	self.Counts = [2]float64{0,0}
	self.CountsY = [3]float64{0,0,0}
	for _,p := range self.RawPatterns {
		self.Counts[p.C] ++
		self.CountsY[p.Y] ++
	}
}
func (self *Clu) Init(vs... *tmpdata.Val)  {
//	if self.Dss == nil {
		self.RawPatterns = vs
		self.Update()
		if self.start != nil {
			close(self.start)
		}
		self.start = make(chan bool)
//	}else{
//		for _,v := range vs {
//			self.RawPatterns = append(self.RawPatterns,v)
//			self.updateDss(v)
//		}
//	}
//	if self.Dss == nil {
//		panic(len(vs))
//	}
}
func (self *Clu) Append(v *tmpdata.Val)  {
	if self.RawPatterns == nil {
		self.RawPatterns = []*tmpdata.Val{v}
//		self.Dss,self.Counts,self.CountsY = GetDss(self.RawPatterns)
	}else{
		self.RawPatterns,_ = AppendSortVal(self.RawPatterns,v)
//		self.updateDss(v)
	}
	self.Counts[v.C]++
	self.CountsY[v.Y]++
	if self.start != nil {
		close(self.start)
	}
	self.start = make(chan bool)
}
func (self *Clu) GetStart() chan bool {
	return self.start
}
func (self *Clu) SetStart(s chan bool) {
	self.start = s
}

//func (self *Clu) updateDss(v *tmpdata.Val) {
//	for j,x := range v.X {
//		self.Dss[j].SetWeight(x)
//	}
//}
func (self *Clu) Check() bool {
	w := float64(len(self.RawPatterns)) * 0.1
	for _,c :=range self.Counts {
		if c < w {
			return true
		}
	}
	return false
}
func (self *Clu) Forecast(v *tmpdata.Val) int {
//	if len(self.RawPatterns) < 10 {
//		return -1
//	}
	if !self.Check() {
		return -1
	}
	if self.Gauged(v) {
		max:=0.0
		mi := -1
		count := 0.0
		for i,y := range self.CountsY {
			count+=y
			if max <= y {
				max = y
				mi = i
			}
		}
//		if mi == 2 {
//			return -1
//		}
		if self.Nn == nil {
			if max / count < 0.5 {
				return -1
			}
			return mi
		}
		vs,err := self.Nn.Update(common.GetX(v.X))
		if err != nil || len(vs) == 0 {
			return -1
		}
		f,err := common.GetForVal(vs)
		if err != nil {
			return -1
		}
		if f == 0 {
			return -1
		}
		return mi
	}
	return -1
}
