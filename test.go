package main
import(
	"fmt"
)
func main () {
	t := make([]int,10)
	for i,_t := range t{
		fmt.Println(i,_t)
		t[i] = i
	}
//	fmt.Println(t)
	for i:= 9;i>=0;i-- {
		t = append(t[:i],t[i+1:]...)
		fmt.Println(i,t)
	}
}
