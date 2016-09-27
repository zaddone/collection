package curves
import (
	"math"
//	"fmt"
)
type Cov struct {

	Sum    float64 `json:"sum"`
	Num    float64 `json:"num"`

	Xsum   float64 `json:"xsum"`
	Xqsum  float64 `json:"xsum"`

	Ysum   float64 `json:"ysum"`
	Yqsum  float64 `json:"ysum"`

	Weight float64 `json:"weight"`
	//Tag  [2]int    `json:"tag"`
	X int
	Y int
}
func (self *Cov) CalWeight(X_val float64,Y_val float64) *Cov {

	self.Sum = self.Sum + X_val*Y_val
	self.Num = self.Num + 1

	self.Xsum = self.Xsum + X_val
	self.Xqsum = self.Xqsum + math.Pow(X_val,2)
//	Dx := math.Sqrt(self.Xqsum/self.Num - math.Pow(self.Xsum, 2)/math.Pow(self.Num, 2))

	self.Ysum = self.Ysum + Y_val
	self.Yqsum = self.Yqsum + math.Pow(Y_val,2)
//	Dy := math.Sqrt(self.Yqsum/self.Num - math.Pow(self.Ysum, 2)/math.Pow(self.Num, 2))

	Cov := self.Sum/self.Num - (self.Xsum*self.Ysum)/math.Pow(self.Num, 2)
	self.Weight = Cov///(Dx*Dy)
	return self

}
type CovMatrix struct {

	Covs   []*Cov

}

func (self *CovMatrix) Clear() {
	self.Covs = nil
}

func (self *CovMatrix) GetCovs() []*Cov {
	return self.Covs
}
func (self *CovMatrix) GetWeights() (wei []float64) {

	wei = make([]float64,len(self.Covs))
	for i,C:=range self.Covs {
		wei[i]=C.Weight
	}
	return wei

}
func (self *CovMatrix) SetCovs(val ...float64) int {

	var x int =0
	for i,n:=range val {
		for j,v:=range val {
			if j<=i {
				continue
			}
			if len(self.Covs)==x{
				c:=&Cov{X:i,Y:j}
//				fmt.Println(i,j)
				self.Covs=append(self.Covs,c)
				c.CalWeight(n,v)
			}else{
				c:=self.Covs[x]
				c.CalWeight(n,v)
			}
			x++
		}
	}
	return x
}
