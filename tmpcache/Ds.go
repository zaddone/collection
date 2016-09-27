package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"math"
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

	return (x-self.Sum/self.Num)/(3*self.Weight)

}

func GetDss(vs []*tmpdata.Val) (Dss []*Ds,C [2]float64,Y [3]float64) {
	for i,v := range vs {
		C[v.C]++
		Y[v.Y]++
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
	return Dss,C,Y
}
