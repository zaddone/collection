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
	Errs  float64 = 0.0001
	Pow  float64 = 2
)
type Curve struct {
	X [][]float64
	Y [][]float64

	inX [][]float64
	w []float64

}
func (self *Curve) SetWeight(w []float64) {
	self.w = w
}
func (self *Curve) Over(Len int) {

	if len(self.X) > Len {
		self.X = self.X[:Len]
		self.inX = self.inX[:Len]
		self.Y = self.Y[:Len]
	}
	copy(self.inX,self.X)

}
func (self *Curve) Init(Len int) {
	self.X = make([][]float64,Len)
	self.inX = make([][]float64,Len)
	self.Y = make([][]float64,Len)
}
func (self *Curve) GetWeightOther() ([]float64,error) {

	Xs:=common.Transpose(self.inX...)
	XX:=common.MatrixInverse(common.MatrixMul(Xs,self.inX))
	XY:=common.MatrixMul(Xs,self.Y)
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
	for i,x := range self.inX {
		y:=self.Y[i][0]
		for j,_x := range x {
			y -= w[j]*_x
		}
		Yerr += y*y
	}
	return Yerr/float64(len(self.X)-1)

}
func (self *Curve) Append(x []float64,y float64,t int) {

//	var xs []float64
	xs := []float64{1}
	for _,_x := range x {
//	for i,_x := range x {
		ax := math.Atan(_x)*2/math.Pi
		xs = append(xs,ax)

		for g := 2.0 ; g <= Pow; g ++ {
			xs = append(xs,math.Pow(ax,g))
		}

//		I:= i+1
//		if I < len(x){
//			for _,_y := range x[I:] {
//				z:= math.Sqrt(_x*_x+_y*_y)
//				sin :=_x/z
//				con :=_y/z
//				xs = append(xs,sin)
////				xs = append(xs,sin*sin)
////				xs = append(xs,sin*sin*sin)
//
//				xs = append(xs,con)
////				xs = append(xs,con*con)
////				xs = append(xs,con*con*sin)
//			}
//		}
	}
//	fmt.Println(y)
	if len(self.X)<=t {
		self.X = append(self.X,xs)
		self.Y = append(self.Y,[]float64{y})
	}else{
		self.X[t] = xs
		self.Y[t] = []float64{y}
	}

}

//type Curves struct {
//	a []float64
//	k []int
//}

//func (self *Curve) GetKey() []int {
//	return self.k
//}
func (self *Curve) GetWeights() []float64 {

	return self.w

}
func (self *Curve)cleanUp(ats []*signal.Atlias)  {
	L := len(ats)
//	attr = make([][]float64,L)
	diff := 0.0
	vol := 0.0
	mac := 0.0
	Time := 0.0

//	Cur = new(Curve)
	self.Init(L)
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
//		ks := (mac/Time)
		t :=at.Ask-at.Bid
		if !math.Signbit(at.Diff) {
			vol -= at.Volume
		}else{
			vol += at.Volume
		}

		y :=t/math.Sqrt(Time*Time + t*t)
//		Cur.Append([]float64{diff,vol,ks,Time},y,i)
		self.Append([]float64{diff,vol,Time},y,i)

	}
//	fmt.Println(Cur)
//	panic(0)
//	return Cur
}

func (self *Curve) AppendAtlias(ats []*signal.Atlias) (err error) {

//	Covs :=new(CovMatrix)
	self.cleanUp(ats)
//	Cur.Chenge()

	self.w,err=self.GetWeightOther()
	if err != nil {
		return err
	}
//	self.k = []int{0}

	return nil
}
