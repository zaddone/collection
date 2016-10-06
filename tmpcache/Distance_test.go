package tmpcache
import (
	"fmt"
	"testing"
	"math/rand"
//	"time"
)
func TestDistance(t *testing.T) {
//	diss:=make([][]*Distance,10)
	var diss []*Distance
	var l int
	for i:=0;i<10;i++ {
		ds := new(Distance)
		ds.dis = rand.Float64()
		fmt.Printf("%p \r\n",diss)
		diss,l = sortDis(diss,ds)
		fmt.Printf("%p \r\n",diss)
		fmt.Println(diss,len(diss),ds,l,i)
//		diss = dis
	}
}
