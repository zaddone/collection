package route
import (
	"fmt"
	"testing"
//	"time"
)
func TestNewDate(t *testing.T) {
//	str := "20161010"
//	var ti time.Time
//	for i:=0;i<10;i++ {
//		fmt.Println(str)
//		str,ti,week =NewDate(str)
//	}

	pi := new(PathInfo)
	pi.Init("/home/dimon/code/go/data/20160909/USDJPY/20",nil)
	for i:=0;i<20;i++ {
		fmt.Println(i,pi)
		p,err := pi.GetNext()
		if err != nil {
			fmt.Println(err)
			continue
		}
		pi = p
	}
}
