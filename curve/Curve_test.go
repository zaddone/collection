package curve
import (

	"fmt"
	"testing"

)

func TestInverse(t *testing.T){

	cur:=new(Curve)
	cur.G = [][]float64{[]float64{9,52,381},[]float64{52,381,3017},[]float64{381,3017,25317}}
	cur.y = []float64{32,147,1025}
	Y:=make([][]float64,len(cur.y))
	for i,_:=range Y {
		Y[i]=[]float64{cur.y[i]}
	}
	Inv_G:=MatrixInverse(cur.G)
	val:=MatrixMul(Inv_G,Y)
	fmt.Println(Inv_G,MatrixInverse(Inv_G))
	fmt.Println(val)
}
func TestInverse2(t *testing.T){

	cur:=new(Curve)
	cur.G = [][]float64{[]float64{-276225,-2796921},[]float64{105507,1503063}}
	cur.y = []float64{129921,-24363}
	Y:=make([][]float64,len(cur.y))
	for i,_:=range Y {
		Y[i]=[]float64{cur.y[i]}
	}
	Inv_G:=MatrixInverse(cur.G)
	val:=MatrixMul(Inv_G,Y)
	fmt.Println(Inv_G,MatrixInverse(Inv_G))
	fmt.Println(val)
}
