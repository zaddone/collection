package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"math"
//	"fmt"
)

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
	i  int
}
func (self *Distance) Init (a,b *tmpdata.Val,i int) {
//	self.dis = EucDistance(a.X,b.X)

	self.dis = (b.Cur.GetErrOther(a.X)+a.Cur.GetErrOther(b.X))*0.5

//	self.dis = b.Cur.GetErrOther(a.X)
//	dis1 :=	a.Cur.GetErrOther(b.X)
//	fmt.Println(self.dis,dis1)
	self.a = a
	self.b = b
	self.i = i
}
func (self *Distance) Find(a *tmpdata.Val) bool {
	if self.a == a || self.b == a{
		return true
	}
	return false
}
func sortDis(dis []*Distance,d *Distance) ([]*Distance,int) {
	Ls:=len(dis)
//	L := Ls-1
	dis  = append(dis,d)
	for i := Ls -1 ;i>=0 ;i-- {
		if dis[i].dis < d.dis {
			break
		}
		dis[i],dis[Ls] = dis[Ls],dis[i]
		Ls = i
	}
	return dis,Ls
}
