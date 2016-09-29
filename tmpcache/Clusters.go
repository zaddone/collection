package tmpcache
import (
//	"fmt"
	"github.com/zaddone/collection/tmpdata"
	"sync"
	"encoding/json"
)
type Clusters struct {
	Clu []*Clu
	Ca *Cache
	Len int
//	OldLong int
//	dellist []int
//	Stop chan bool
//	Stop1 chan bool
	ValCount int
//	SyncClear  bool
//	ClearList chan []*Clu
}
func (self *Clusters) SyncMergeServer(sy *sync.WaitGroup){
	defer sy.Done()
}
func getTmpClu(clu *Clu,clus []*Clu) (tmp []*tmpdata.Val,tmpc []*Clu) {
	for _,cl:= range clus {
		if cl == clu {
			continue
		}
		p:=cl.Core
		if cl.Core == nil {
			p=cl.RawPatterns[len(cl.RawPatterns)-1]
		}
		p.SetD(len(tmp))
		tmp = append(tmp,p)
		tmpc = append(tmpc,cl)
	}
	return tmp,tmpc
}
func (self *Clusters) mergeClu(clu *Clu,clus []*Clu,isP bool) bool {
	if len(clu.RawPatterns) < self.Len {
		return false
	}
//	if clu.RawPatterns !=nil {
//		for _,h :=range clu.Counts {
//			if h == 0 {
//				return false
//			}
//		}
//	}
	tmp,tmpc := getTmpClu(clu,clus)
	if len(tmp) < self.Len {
		return false
	}
	var out []int
	out=clu.FindOutPat()
	tmpClu := new(Clu)
	tmpClu.Init(tmp...)
	cL := len(clu.RawPatterns) - 200
	Gtmp := make(map[*Clu]bool)
	for _,i := range out {
		if cL >0 && i < cL {
			clu.RawPatterns = append(clu.RawPatterns[:i],clu.RawPatterns[i+1:]...)
			continue
		}
		p := clu.RawPatterns[i]
		pi := tmpClu.FindMinVal(p)
		mc :=tmpc[pi]
		if len(mc.RawPatterns) < 2 || mc.Gauged(p) {
			mc.Append(p)
			if isP {
				Gtmp[mc] = true
			}
			clu.RawPatterns = append(clu.RawPatterns[:i],clu.RawPatterns[i+1:]...)
		}
	}
	clu.Update()
	for k,_ := range Gtmp {
		self.mergeClu(k,tmpc,false)
	}
//	if  len(clu.RawPatterns) < 2 {
//		clu.Clear()
//		return true
//	}
	self.Ca.AppendWaits(clu)

	return false
}
func (self *Clusters) Getbakdb() (db []byte) {
	dbs:= make( [][]byte,len(self.Clu))
	i:=0
	for _,c := range self.Clu {
		if c.RawPatterns == nil {
			continue
		}
		d,err := c.OutputData()
		if err != nil {
			panic(err)
		}
		dbs[i] = d
		i++
	}
	tmp := &tmpClusters{Count:self.ValCount,Dbs:dbs[:i]}
	db,err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	return db
}
func (self *Clusters) GetCountSeds() (seds []*Clu,count int,countN int,db []byte) {
//	Le := len(self.Clu)
	var dbs [][]byte
//	for i:= Le-1;i>=0;i--{
//		c:=self.Clu[i]
	for _,c := range self.Clu {
		if c.Lock || c.RawPatterns == nil {
			continue
		}
		d,err := c.OutputData()
		if err != nil {
			panic(err)
		}
		dbs = append(dbs,d)
		L := len(c.RawPatterns)
	//	count += L
		if L > self.Len {
			if c.Check() {
				countN += L
				seds = append(seds,c)
			}
		}
	}
	count = self.ValCount
	tmp := &tmpClusters{Count:count,Dbs:dbs}
	db,err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	return seds,count,countN,db
}
func (self *Clusters) LoadData(data []byte) error {

	tmp:= new(tmpClusters)
	err :=json.Unmarshal(data,tmp)
	if err != nil {
		return err
	}
	self.Clu = make([]*Clu,len(tmp.Dbs))
	for i,d := range tmp.Dbs {
		self.Clu[i],err = InputData(d)
		if err != nil {
			panic(err)
		}
	}
	self.ValCount = tmp.Count
	return nil

}
type tmpClusters struct {
	Count int
	Dbs [][]byte
}
func (self *Clusters) Init(c *Cache) {
	self.Len = c.Len
	self.Ca = c
//	self.Stop1 = make(chan bool)
//	go self.SyncClear()
}
func (self *Clusters) getTmpCore(clus []*Clu) (tmp []*tmpdata.Val,clu []*Clu,tmp1 []*tmpdata.Val,clu1 []*Clu) {
	for _,cl := range clus {
		if cl.Core == nil {
			p :=cl.RawPatterns[len(cl.RawPatterns)-1]
			p.SetD(len(tmp))
			tmp = append(tmp,p)
			clu = append(clu,cl)
		}else{
			cl.Core.SetD(len(tmp1))
			tmp1 = append(tmp1,cl.Core)
			clu1 = append(clu1,cl)
		}
	}
	return tmp,clu,tmp1,clu1
}
func (self *Clusters) getTmpVals(y int) (tmp []*tmpdata.Val,clu []*Clu) {
	for i:= len(self.Clu)-1;i>=0;i--{
		c :=self.Clu[i]
		if c == nil || (c.Lock && c.RawPatterns == nil)  {
			self.Clu =append(self.Clu[:i],self.Clu[i+1:]...)
			continue
		}
		p := c.Core
		if p == nil {
			p =c.RawPatterns[len(c.RawPatterns)-1]
		}
		if p.C == y {
			p.SetD(len(tmp))
			tmp = append(tmp,p)
			clu = append(clu,c)
		}
	}
//	self.ClearClu()
	return tmp,clu
}
func (self *Clusters) getTmpVal(isout *Clu) (tmp []*tmpdata.Val,clu []*Clu) {
	for i:= len(self.Clu)-1;i>=0;i--{
		c :=self.Clu[i]
		if c == nil || c.Lock || isout == c {
			if c.RawPatterns == nil && c.Lock {
				self.Clu =append(self.Clu[:i],self.Clu[i+1:]...)
			}
			continue
		}
		p := c.Core
		if p == nil {
			p =c.RawPatterns[len(c.RawPatterns)-1]
		}
		p.SetD(len(tmp))
		tmp = append(tmp,p)
		clu = append(clu,c)
	}
//	self.ClearClu()
	return tmp,clu
}
func findSameSed(v *tmpdata.Val,tmp []*tmpdata.Val,clus []*Clu) (*Clu)  {
//	if len(clus) < 2 {
//		return nil
//	}
	tmpClu := new(Clu)
	tmpClu.Init(tmp...)
	i := tmpClu.FindMinVal(v)
	return clus[i]
}
func (self *Clusters) GetSameSed(v *tmpdata.Val) (*Clu,[]*Clu)  {
	if len(self.Clu) < 2 {
		return nil,nil
//	if len(tmp) < self.Len {
//		clu:=&Clu{RawPatterns:[]*tmpdata.Val{v}}
//		clu := new(Clu)
//		clu.Append(v)
//		self.Clu = append(self.Clu,clu)
//		return clu
	}
	tmp,clu:= self.getTmpVal(nil)
	if len(tmp) < 2 {
		return nil,nil
	}
	tmpClu := new(Clu)
	tmpClu.Init(tmp...)
	i := tmpClu.FindMinVal(v)
	return clu[i],clu
}
func (self *Clusters) AppendVal(v *tmpdata.Val)  {
	v.SetH(self.ValCount)
	self.ValCount++

//	tmp,clus:= self.getTmpVals(v.C)
	tmp,clus:= self.getTmpVal(nil)
	var clu *Clu
	if len(tmp)<self.Len {
		clu = new(Clu)
		clu.Append(v)
		self.Clu = append(self.Clu,clu)
		return
	}else{
		tmpClu := new(Clu)
		tmpClu.Init(tmp...)
		vs := tmpClu.FindSortVal(v)
		clu = checkGauged(clus,vs,v)
		if clu == nil {
//			clu = clus[vs[0]]
			clu = new(Clu)
			clu.Append(v)
			self.Clu = append(self.Clu,clu)
			Integration(clu,vs,clus)
			return
		}
		if len(clu.RawPatterns) < self.Len {
			clu.Append(v)
			return
		}
		if clu.Gauged(v) {
			clu.Append(v)
			self.mergeClu(clu,clus,true)
		}else{
			clu = new(Clu)
			clu.Append(v)
			self.Clu = append(self.Clu,clu)
			Integration(clu,vs,clus)
		}
		return
//	}
	}
	return
}
func checkGauged(clus []*Clu,vs []int,vl *tmpdata.Val) *Clu {
	min := 0.0
	I := -1
	er :=0
//	fmt.Println(vs)
//	var test  []int
	for _,v := range vs {
		cl := clus[v]
		if len(cl.RawPatterns) < 2 {
			if I == -1 {
				I = v
			}
			break
		}
		a,b := cl.Gaugedh(vl)
//		fmt.Println(a,b)
		if a>b {
//			test = append(test,v)
			if min == 0 || min > b {
				min = b
				I = v
			}
		}else{
			er ++
			if er > 15 {
				break
			}
		}
	}
//	fmt.Println(test)
	if I == -1 {
		return nil
	}
	return clus[I]
}
func Integration(clu *Clu,vs []int,clus []*Clu)  {
	isN := false
	for _,i := range vs[:5]{
		if clu.Integ(clus[i]) {
			isN = true
			break
		}
	}
	if !isN {
		return
	}
	for _,i := range vs{
		k:=0
		cl :=clus[i]
		L :=len(cl.RawPatterns)
		end := L - 100
		if end <0 {
			end = 0
		}
		for i:=L-1; i>=end;i--{
			v:= cl.RawPatterns[i]
//		for _,v:= range cl.RawPatterns{
			if clu.Gauged(v) {
				k ++
				cl.RawPatterns = append(cl.RawPatterns[:i],cl.RawPatterns[i+1 :]...)
				clu.Append(v)
			}
		}
		if k == 0 {
			break
		}else if k == L || k == 100 {
			cl.Clear()
		}else{
			cl.Update()
		}
	}
}
//
//func (self *Clusters) AppendVals(v *tmpdata.Val)  {
////	if self.Stop1 != nil {
////		<-self.Stop1
////	}
//	v.SetH(self.ValCount)
//	self.ValCount++
//
//	tmp,clus:= self.getTmpVal(nil)
//	var clu *Clu
//	if len(tmp)<50 {
//		if len(tmp)<2 {
//			clu = new(Clu)
//			clu.Append(v)
//			self.Clu = append(self.Clu,clu)
//			return
//		}else{
//			clu = findSameSed(v,tmp,clus)
//			if len(clu.RawPatterns) < self.Len {
//				clu.Append(v)
//				return
//			}
//			if clu.Gauged(v) {
//				clu.Append(v)
//				self.mergeClu(clu,clus,true)
//			}else{
//				clu = new(Clu)
//				clu.Append(v)
//				self.Clu = append(self.Clu,clu)
//			}
//			return
//		}
//	}else{
//		t1,c1,t2,c2 := self.getTmpCore(clus)
//		if len(t2)>2{
//			clu = findSameSed(v,t2,c2)
//			if clu.Gauged(v) {
//				clu.Append(v)
//				self.mergeClu(clu,clus,true)
//				return
//			}
//		}
//		if len(t1) >2 {
//			clu = findSameSed(v,t1,c1)
//			if clu.RawPatterns[len(clu.RawPatterns)-1].C == v.C {
//				clu.Append(v)
//				self.mergeClu(clu,clus,true)
//				return
//			}
//		}
//	}
//	clu = new(Clu)
//	clu.Append(v)
//	self.Clu = append(self.Clu,clu)
//}
