package common
import (
	"math"
	"net"
	"fmt"
)
const (
	SocketLen  int = 512
	ORDER int = 1
)
func ByteToInt(b []byte) int {
	s:=SocketLen/2
	num:=0
	for i:=len(b)-1;i>=0;i--{
		num = num*s +int(b[i])
//		fmt.Println(num)
	}
	return num
}
func IntToByte(l int,b []byte)[]byte{
	s:=SocketLen/2
	k:=l/s
	h:=l%s
//	fmt.Println(k,h)
	b=append(b,uint8(h))
	if k>s {
		return IntToByte(k,b)
	}
	return append(b,uint8(k))
}
func ConnWrite(conn net.Conn,jsonData []byte) (error) {
	Len := len(jsonData)
	L := IntToByte(Len,nil)
	_,err:=conn.Write(L)
	if err != nil {
		return err
	}
	var begindata [SocketLen]byte
	n,err:=conn.Read(begindata[0:])
	if err != nil {
		return err
	}
	if n != 1 && begindata[0] != 1 {
		return fmt.Errorf("code 1")
	}
	for i:=0;i<Len;i+=SocketLen {
		end:=i+SocketLen
		if end>Len {
			end=Len
		}
		_,err = conn.Write(jsonData[i:end])
		if err != nil {
			return err
		}
		n,err = conn.Read(begindata[0:])
		if err != nil {
			return err
		}
		if n != 1 && begindata[0] != 1 {
			return fmt.Errorf("code 1")
		}
	}
	return nil
}
func ConnRead(conn net.Conn) ([]byte,error) {
	var begindata [SocketLen]byte

	n,err := conn.Read(begindata[0:])
	if err != nil {
		return nil,err
	}
	_,err=conn.Write([]byte{1})
	if err != nil {
		return nil,err
	}
	L := ByteToInt(begindata[:n])
	status := make([]byte,L)
	for i:=0;i<L; {
		num,err:=conn.Read(begindata[0:])
		if err!= nil {
			return nil,err
		}
		end:=i+num
		if end>L {
			return nil,fmt.Errorf("%d %d",end,L)
		}
		copy(status[i:end],begindata[0:num])
		i = end
		_,err=conn.Write([]byte{1})
		if err != nil {
			return nil,err
		}
	}
	return status,nil

}
func GetX(w []float64) (f []float64) {

	f=[]float64{1}
	for t:=1; t <= ORDER; t++ {
		for _,x := range w {
			f = append(f,math.Pow(x,float64(t)))
		}
	}

	return f

}
func Swap(num int, NodeSort []int, Nodes  []float64, step int) int {

	if num == 0{
		return step
	}
	upNum:=num-1
	a:=NodeSort[num]
	b:=NodeSort[upNum]
	if Nodes[a] > Nodes[b] {
		NodeSort[num], NodeSort[upNum] = NodeSort[upNum], NodeSort[num]
		step++
		return Swap(upNum,NodeSort,Nodes,step)
	}
	return step

}
func GetMaxSort(num int, NodeSort []int, Nodes  []float64, step int) int {

	NodeSort[num]=num
	if num!=0 {
		step=Swap(num,NodeSort,Nodes,step)
	}
	num++
	if int(num)==len(NodeSort){
		return step
	}
	return GetMaxSort(num,NodeSort,Nodes,step)

}
func IntToString(a []int) (str string) {

	str=""
	for _,_a := range a {
		str=fmt.Sprintf("%s%d",str,_a)
	}
	return str

}
func cos(a []float64,b []float64) float64 {

	var sum float64
	var pA float64
	var pB float64
	for i,_a := range a {
		_b := b[i]
		sum+=_a*b[i]
		pA += _a*_a
		pB += _b*_b
	}
	return sum/(math.Sqrt(pA)*math.Sqrt(pB))

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
		//		fmt.Println(i,j,k,sum)
				if len(a[i])> k && len(b[k]) > j {
					sum += a[i][k] * b[k][j]
				}else{
					fmt.Println(len(a),len(a[0]))
					fmt.Println(len(b),len(b[0]))
					panic(0)
				}
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
func GetForVal(vs []float64) ( int, error) {
	if len(vs) == 1 {
		if vs[0] >0.5 {
			return 1,nil
		}else{
			return 0,nil
		}
	}
	f := -1
	for i,v := range vs {
		if v >0.5 {
			if f == -1 {
				f = i
			}else{
				return f,fmt.Errorf("is err")
			}
		}
	}
	return f,nil
}
