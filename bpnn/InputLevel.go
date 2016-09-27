package bpnn
import (
	"math"
//	"os"
	"fmt"
	"math/rand"
	"github.com/zaddone/collection/tmpdata"
	"encoding/json"
	"time"
)
const (
	Output  int = 1
)

type TmpData struct {
	Err  float64
	N  float64
	M  float64
	Base  []byte
	Weight []byte
	WeiLen int
//	LongLen int
//	InputLen int
//	OutPutLen int
}
type InputLevel struct {
	Ni  int

	Nhr  int
	Nhc  int

	No  int

	Alpht int

	ai  []float64
	Stop  float64
	IsTest bool
	level *HiddenLevel
	output  *HiddenLevel
	AllY int
}
func (self *InputLevel) Init(i ,hr ,hc,no int){
	self.Ni = i
	self.Nhr = hr
	self.Nhc = hc
//	self.No = no
	self.No = Output
	self.Alpht = 0
	self.IsTest = false

	self.Stop = 0.001
	self.load()

}
func (self *InputLevel) load() {
	self.ai = make([]float64,self.Ni)
	self.level = new(HiddenLevel)
	self.output= self.level.Init(self.Ni,self.Nhr,self.Nhc,self.No)

}
func (self *InputLevel) IsInit() bool {
	if self.ai == nil {
		return false
	}
	return true
}
func  GetAlpR(n float64,alpht int) (float64,int) {
	a:=math.Pow10(alpht)
	if n < a{
		alpht--
		return GetAlpR(n,alpht)
	}
//	fmt.Println(n,a)
	return n+a,alpht
}
func  GetAlp(n float64,alpht int) (float64,int) {
	a:=math.Pow10(alpht)
	if math.Floor(n/a) == 1  || n < a{
		alpht--
		return GetAlp(n,alpht)
	}
	return n-a,alpht
}
func (self *InputLevel) Update(inputs []float64) ([]float64,error) {
	if len(inputs) != self.Ni {
		fmt.Println("long",len(inputs) , self.Ni)
		return nil,fmt.Errorf("input is err")
	}
	for i,x:= range inputs {
		self.ai[i] = x
	}
	return self.level.update(inputs)
}

func (self *InputLevel) backPropagate(targets []float64,N,M float64) (error) {
	tar,err:=self.output.backPropagate(targets,N,M)
	if err != nil {
		return err
	}
	self.level.updateLevel(tar,self.ai,N,M)
//	e:=self.output.GetErr(targets)
//	return e,nil
	return nil
}

func (self *InputLevel) GetJsondate() []byte {

	b,err:=json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return b

}

func (self *InputLevel) SetJsondate(b []byte) error {

	return json.Unmarshal(b,self)
//	if err != nil {
//		panic(err)
//	}

}
func (self *InputLevel) GetTmp(N ...float64 ) *TmpData {

	tmp:=new(TmpData)
	tmp.Err =N[0]
	tmp.N = N[1]
	if len(N) >2 {
		tmp.M = N[2]
	}
	tmp.Base = self.GetJsondate()
	Wei:=self.level.GetWeight(nil)
	tmp.WeiLen = len(Wei)

//	tmp.LongLen = len(Wei[0])
	tmp.Weight,_ =json.Marshal( Wei )
	return tmp

}
func (self *InputLevel) SetTmp(tmp *TmpData,isload bool) error {

	err := self.SetJsondate(tmp.Base)
	if err != nil {
		return err
	}
	if isload {
		self.load()
	}
	var Wei [][][]float64 = make([][][]float64,tmp.WeiLen)
	err=json.Unmarshal(tmp.Weight,&Wei)
	if err != nil {
		return err
	//	panic(err)
	}
	return self.level.SetWeight(Wei,0)

}
func (self *InputLevel) Copy() (I *InputLevel) {

	I = new(InputLevel)
	I.Init(self.Ni,self.Nhr,self.Nhc,self.No)
	return I

}
func SortTmpData(Ts []*TmpData,T *TmpData) []*TmpData {

	if Ts == nil {
		return []*TmpData {T}
	}else{
		Ts = append(Ts,T)
	}
	L:=len(Ts)-1
	for i:=L-1;i>=0;i-- {
		if Ts[i].Err < T.Err {
			Ts[i],Ts[L] = Ts[L],Ts[i]
			L=i
		}else{
			break
		}
	}
	return Ts
}
func (self *InputLevel) Train(iteration int,patterns []*tmpdata.Val,N,M float64) (*TmpData, error) {
//	var lastTmp *TmpData
//	E:=0
	for i:=0;i<iteration;i++ {
//		fmt.Println(i)
		self.back(patterns,N,M)
//		Err := self.GetErr(patterns,N,M)
//		if Err < self.Stop {
//			self.IsTest = true
//			return self.GetTmp(Err,N,M),nil
//		}
//		if lastTmp == nil {
//
//			lastTmp = self.GetTmp(Err,N,M)
//			continue
//		}
//		k:=Err - lastTmp.Err
//		if  k < 0 {
//			lastTmp = self.GetTmp(Err,N,M)
//			if E > 0 {
//				E = 0
//			}
//		}else{
//			if k > self.Stop {
//				E++
//				if E >3{
//				break
//				}
//			}
//		}
	}
	Err := self.GetErr(patterns,N,M)
	if Err < self.Stop {
		self.IsTest = true
		return self.GetTmp(Err,N,M),nil
	}
	self.IsTest = false
	return self.GetTmp(Err,N,M),fmt.Errorf("err: %.5f",Err)
}
func (self *InputLevel) back(patterns []*tmpdata.Val,N,M float64)  {
	for _,p := range patterns {
		inputs := p.X
//		targets := 0.0
//		if p.Y {
//			targets = 1
//		}
		_,err := self.Update(inputs)
		if err != nil {
			panic(err)
			//return err
		}

		tar:=make([]float64,Output)
		if Output == 1 {
			if p.Y == self.AllY {
				tar[0] = 1
			}
			//tar[0] = float64(p.Y)
		}else{
			tar[p.Y] = 1
		}
//		tar[p.Y] = 1
		err = self.backPropagate(tar,N,M)
		if err != nil {
			panic(err)
			//return err
		}
//		e:=self.output.GetErr(tar)
//		Err += e
	}
}
func (self *InputLevel) GetErr(patterns []*tmpdata.Val,N,M float64) (Err float64) {

	for _,p := range patterns {
		inputs := p.X
//		targets := 0.0
//		if p.Y {
//			targets = 1
//		}
		_,err := self.Update(inputs)
		if err != nil {
			panic(err)
			//return err
		}

		tar:=make([]float64,Output)
		if Output == 1 {
			if p.Y == self.AllY {
				tar[0] = 1
			}
			//tar[0] = float64(p.Y)
		}else{
			tar[p.Y] = 1
		}
//		tar[p.Y] = 1
//		_,err = self.backPropagate(tar,N,M)
//		if err != nil {
//			panic(err)
//			//return err
//		}
		e:=self.output.GetErr(tar)
		Err += e
	}
//	return Err/float64(len(patterns))
	return Err

}
func RandFloat64(num int) (w []float64) {

	w=make([]float64, num)
	rand.Seed(time.Now().UnixNano())
	for i:=0;i<num;i++ {
//		time.Sleep(1*time.Microsecond)
//		w[i]=float64(rand.Intn(100))
//		n:=rand.Float64()
//		w[i]=(0.5-rand.Float64())*100//*(upper+lower)-lower
		w[i]=(0.5-rand.Float64())
//		w[i]=rand.Float64()
//		fmt.Println(w[i])
	}
	return w

}

func makeMatrix(I,J int,isRand bool) (m [][]float64){

	m=make([][]float64,I)
	for i:=0;i<I;i++ {
		if isRand {
			m[i] = RandFloat64(J)
		}else {
			m[i] = make([]float64,J)
		}

	}
	return m

}
func relu(x float64) float64 {
	return math.Max(0,x)
}
func reluder(x float64) float64 {
	if x<0 {
		return 0
	}else{
		return 1
	}
}
func softplus(x float64)float64{
	return math.Log(1+math.Exp(x))
}
func sigmoid(x float64) float64 {

	return 1.0/(1.0 + math.Exp(-x))
//	return math.Tanh(x)

}
func dsigmoid(y float64) float64 {

	return y*(1-y)
//	return 1-y*y
}
