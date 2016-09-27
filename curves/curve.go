package curves
import (
	"math"
	"fmt"
	"github.com/zaddone/collection/signal"
	"github.com/zaddone/collection/common"
//	"os"
)
const (
//	ORDER float64 = 5
	Errs  float64 = 0.05
)
type Curve struct {
	x [][]float64
	y [][]float64
	Inx [][]float64
	Iny [][]float64
	t float64
}
func (self *Curve) Init(Len int) {
	self.x = make([][]float64,Len)
	self.y = make([][]float64,Len)
	self.Inx = make([][]float64,Len)
	self.Iny = make([][]float64,Len)
}
func (self *Curve) Chenge() {
	ExchangeDs(self.x)
	ExchangeDs(self.y)
}
func (self *Curve) SetData(t float64) error {
	if t >1 {
		return fmt.Errorf("is nil")
	}
	self.t = t
	self.Inx = make([][]float64,len(self.x))
	for j:=1.0;j<=t;j++ {
		for i,X:= range self.x{
			for _,_x := range X{
				self.Inx[i] = append(self.Inx[i],math.Pow(_x,j))
			}
		}
	}
	copy(self.Iny,self.y)
	return nil
}
func (self *Curve) GetWeight(t float64) ([]float64,error) {
	err :=self.SetData(t)
	if err != nil {
		return nil,err
	}
//	InX:=self.Inx
	InX:=make([][]float64,len(self.Inx))
	for i,_ := range InX {
		InX[i] = append([]float64{1},self.Inx[i]...)
	}
//	fmt.Println(InX)
//	panic(0)

//	for t:=1.0;t<=ORDER;t++ {
//		for i,X:= range self.x{
//			for _,_x := range X{
//				InX[i] = append(InX[i],math.Pow(_x,t))
//			}
//		}
//	}

	Xs:=common.Transpose(InX...)
	XX:=common.MatrixInverse(common.MatrixMul(Xs,InX))
	XY:=common.MatrixMul(Xs,self.Iny)
	val :=common.MatrixMul( XX,XY)
	A := common.Transpose(val...)[0]


	R := self.GetErr(InX,A)
	if R < Errs {
//		fmt.Println(A)
		return A,nil
	}
	t ++
	return self.GetWeight(t)
//	return nil,fmt.Errorf("is err")

}
func (self *Curve) GetWeightOther() ([]float64,error) {

	Xs:=common.Transpose(self.x...)
	XX:=common.MatrixInverse(common.MatrixMul(Xs,self.x))
	XY:=common.MatrixMul(Xs,self.y)
	val :=common.MatrixMul( XX,XY)
	A := common.Transpose(val...)[0]
	R := self.GetErrOther(A)
	if R < Errs {
		return A,nil
	}
	return nil,fmt.Errorf("is err")

}
func (self *Curve)GetErrOther(w []float64) float64 {
	var Yerr float64
	for i,x := range self.x {
		y:=self.y[i][0]
		for j,_x := range x {
			y -= w[j]*_x
		}
		Yerr += y*y
	}
	return Yerr/float64(len(self.x)-1)
}
func (self *Curve)GetErr(InX [][]float64,w []float64) float64 {

	var Yerr float64
	for i,x := range InX {
		y:=self.y[i][0]
		for j,_x := range x {
			y -= w[j]*_x
		}
		Yerr += y*y
	}
	return Yerr/float64(len(self.x)-1)

}
func (self *Curve) Append(x []float64,y float64,i int) {
//	var xs []float64
	xs := []float64{1}
	for i,_x := range x {
		I:= i+1
		if I < len(x){
			for _,_y := range x[I:] {
				z:= math.Sqrt(_x*_x+_y*_y)
				xs = append(xs,_x/z)
				xs = append(xs,_y/z)
			}
		}
	}
	self.x[i] = xs

	//self.x[i] = x
	self.y[i] = []float64{y}
}

type Curves struct {
	a []float64
	k []int
}

func (self *Curves) GetKey() []int {
	return self.k
}
func (self *Curves) GetWeights() []float64 {

	return self.a

}
func cleanUp(ats []*signal.Atlias) (Cur *Curve) {
	L := len(ats)
//	attr = make([][]float64,L)
	diff := 0.0
	vol := 0.0
	mac := 0.0
	Time := 0.0

	Cur = new(Curve)
	Cur.Init(L)
	lastMac:=0.0
	for i:=0; i<L;i++ {
		at:=ats[i]
		diff += at.Diff
		if lastMac == 0 {
			mac = 1
		}else{
//			mac = float64(at.Msc)-lastMac
			mac += float64(at.Msc)-lastMac
		}
		lastMac = float64(at.Msc)
		Time += float64(at.Time)
		ks := (mac/Time)
		t :=at.Ask-at.Bid
		if !math.Signbit(at.Diff) {
			vol -= at.Volume
		}else{
			vol += at.Volume
		}

		y :=t/math.Sqrt(Time*Time + t*t)
		Cur.Append([]float64{diff,vol,ks},y,i)

	}
	return Cur
}

func (self *Curves) AppendAtlias(ats []*signal.Atlias) (err error) {

//	Covs :=new(CovMatrix)
	Cur:=cleanUp(ats)
//	Cur.Chenge()

	self.a,err=Cur.GetWeightOther()
	if err != nil {
		return err
	}
	self.k = []int{int(Cur.t)}

	return nil
}
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

	return (x-self.Sum/self.Num)/(3*self.Weight)

}
func ExchangeDs(Tx [][]float64) {
	var D []*Ds
	for i,T := range Tx {
		if i == 0 {
			D = make([]*Ds,len(T))
			for j,_:= range D {
				D[j] = new(Ds)
			}
		}
		for j,t := range T {
			D[j].SetWeight(t)
//			fmt.Println(t,D[j].GetScore(t))
		}
	}
	for _,T := range Tx {
		for j,t := range T {
			T[j] = D[j].GetScore(t)
		}
//		fmt.Println(Tx[i])
	}
}
