package curve
import (
	"math"
	"fmt"
	"github.com/zaddone/collection/signal"
//	"os"
)
const (
	ORDER int = 4
	Errs  float64 = 0.55
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
//			fmt.Println(t,Tx[i][j])
		}
	}
}

type Curve struct {
	G [][]float64
	a []float64
	y []float64
	err  float64
}
func (self *Curve) show() {
	for _,g := range self.G {
		fmt.Println(g)
	}
}
func (self *Curve) GetWeights() []float64 {

//	fmt.Println(self.a)
	return self.a

}

func (self *Curve) AppendAtlias(ats []*signal.Atlias) error {

//	var lastAt *signal.Atlias
	if self.G == nil {
		self.init()
	}
	var msc float64 = 1
	var lastAr *signal.Atlias = nil
	L:=len(ats)
//	var X [][]float64 = make([][]float64,L)
	var TX [][]float64 = make([][]float64,L)
	var Y [][]float64 = make([][]float64,L)
//	var Dss []*Ds = make([]*Ds,L)
	var diff float64
	for i,at:= range ats {
		if lastAr == nil {
			lastAr = at
		}else {
			msc+= float64(at.Msc- lastAr.Msc)/1000
		}
		diff+=at.Diff
//		TX[i] = []float64{1}
		for t:=1; t <= ORDER; t++ {
			TX[i] = append(TX[i],math.Pow(diff,float64(t)))
			TX[i] = append(TX[i],math.Pow(at.Volume,float64(t)))
			TX[i] = append(TX[i],math.Pow(float64(at.Time),float64(t)))
			TX[i] = append(TX[i],math.Pow(float64(at.Ask-at.Bid),float64(t)))
		}

	//	Y[i]= make([]float64,1)
	//	Y[i][0] = msc
		Y[i] = []float64{msc}
	}
	ExchangeDs(TX)
	for j,x := range TX {
		TX[j] = append(x,1)
	}
//	fmt.Println(TX)
	ExchangeDs(Y)
	Xs:=Transpose(TX...)
	XX:=MatrixInverse(MatrixMul(Xs,TX))
	XY:=MatrixMul(Xs,Y)
//	fmt.Println(XX,XY,Y)
	val :=MatrixMul( XX,XY)
	A := Transpose(val...)
	self.a = A[0]
//	fmt.Println(self.a)
//	os.Exit(0)
//	fmt.Println(self.a)
	E := make([]float64,L)
	var Yerr float64
	var YSum float64
	for i,x := range TX {

		y:=Y[i][0]
		YSum += y*y
		E[i] =  y
		for j,_x := range x {
			E[i]-=self.a[j]*_x
		}
//		fmt.Println(E[i])
		Yerr += E[i]*E[i]
	}
	R:=(YSum-Yerr)/YSum
//	fmt.Println(R,Yerr)
//	os.Exit(0)
	if R > Errs {
		return nil
	}else{
		return fmt.Errorf("err")
	}

//	fmt.Println("R",R,YSum,Yerr,L)

//	Inv_G:=MatrixInverse(self.G)
//	val := MatrixMul(Inv_G,Transpose(self.y))
//	//fmt.Println(val)
//	for i,v:=range val {
//		fmt.Println(v)
//		self.a[i]=v[0]
//	}
	

}

func Transpose(data ...[]float64) (val [][]float64) {

	val = make([][]float64,len(data[0]))
	for i,_ := range val {
		val[i]=make([]float64,len(data))
	}

	for i,d:= range data {
		for j,_d := range d{
			val[j][i] = _d
		}
	}
	return val

}

func (self *Curve) init() {
	self.G = make([][]float64,ORDER)
	self.a = make([]float64,ORDER)
	self.y = make([]float64,ORDER)
	for i,_ := range self.G {
		self.G[i]=make([]float64,ORDER)
	}
}
func (self *Curve) SetY (y float64) {
	for i,g := range self.G {
		self.y[i] += g[i]*y
	}
}
func (self *Curve) SetG (x float64) {
	for i:=0;i < ORDER*2;i++ {
		self.findX(math.Pow(x,float64(i)),i)
	}
}
func (self *Curve)findX(x float64,num int){
	var val float64 = 0
	for i:=0;i< ORDER;i++{
		for j:=0; j<ORDER;j++ {
			if i+j != num {
				continue
			}
			if val != 0 {
				self.G[i][j] = val
			}else{
				self.G[i][j] += x
				val = self.G[i][j]
			}
		}
	}
}

func MatrixMul(a, b [][]float64) (c [][]float64) {
	var sum float64 = 0

	col:=len(a)
	row:=len(b[0])
	c= make([][]float64,col)
	for i:=0;i<col;i++{
		c[i] = make([]float64,row)
		for j:=0; j<row;j++{
			for k:=0; k<len(b); k++ {
				sum += a[i][k] * b[k][j]
		//		fmt.Println(i,j,k,sum)
			}
		//	fmt.Println(i,j)
			c[i][j] = sum
			sum=0
		}
	}
	return c
}
func MatrixInverse(G [][]float64) (Inv [][]float64) {
	M:=len(G)
	N:=M*2
	Inv_G := make([][]float64,M)
	Inv = make([][]float64,M)
	for i,_ := range Inv_G {
		Inv_G[i] = make([]float64,N)
//		fmt.Println(Inv_G[i],self.G[i])
		copy(Inv_G[i],G[i])
		Inv_G[i][M+i] = 1
	}
	for i,g := range Inv_G {
		if Inv_G[i][i] == 0 {
			for k:=i;k<M;k++ {
				if Inv_G[k][i] != 0 {
					for j:=0; j<N;j++ {
						Inv_G[i][j],Inv_G[k][j] = Inv_G[k][j],Inv_G[i][j]
					}
					break
				}
				if k==M {
					fmt.Println("matrix is not inverse !!!")
					return nil
				}
			}
		}
		for j:=N-1; j>=i;j--{
			Inv_G[i][j]/=g[i]
		}
		for k:=0; k<M;k++ {
			if k == i {
				continue
			}
			temp:=Inv_G[k][i]
			for j:=0;j<N;j++ {
				Inv_G[k][j]-=temp*Inv_G[i][j]
			}
		}
	}
	for i,_:= range Inv {
		Inv[i] = Inv_G[i][M:]
	}
	return Inv

}
func (self *Curve) ErrorSum(y float64) {
	
}
