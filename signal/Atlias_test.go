package signal
import (
	"fmt"
	"testing"
)
func TestCompSame(t *testing.T){
	at1 := new(Atlias)
	at1.Ask=2
	at2 := new(Atlias)
	b:=at1.CompSame(at2)
	fmt.Println(b)

}
