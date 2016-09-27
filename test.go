package main
import(
	"fmt"
)
func main () {
	t := make([]int,10)
	for i,_t := range t[5:]{
		fmt.Println(i,_t)
	}
}
