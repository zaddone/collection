package bpnn
import (
	"fmt"
)

type HiddenLevel struct {
	Wi  [][]float64
	ci  [][]float64

	ah  []float64
	sum  []float64
	Per *HiddenLevel
	Next *HiddenLevel
}
const (
	LINE  int =  8
)
func (self *HiddenLevel) Init (row,col,w,o int) *HiddenLevel {

//	fmt.Println(row)
	self.Wi = makeMatrix(row,col,true)
	self.ci = makeMatrix(row,col,false)
	rows:=col+1
	self.ah = make([]float64,rows)
	self.sum = make([]float64,rows)
	self.Next = new(HiddenLevel)
	self.Next.Per = self
	w--
	if w > 0 {
//		row = col++
//		return self.Next.Init(row,col,w,o)
		return self.Next.Init(rows,col,w,o)
	}
	self.Next.SetOutput(rows,o)
	return self.Next

}
func (self *HiddenLevel) SetOutput(row,col int) {

	self.Wi = makeMatrix(row,col,true)
	self.ci = makeMatrix(row,col,false)
	self.ah = make([]float64,col)
	self.sum = make([]float64,col)
}
func (self *HiddenLevel) GetWeight(ws [][][]float64)([][][]float64) {
	ws = append(ws,self.Wi)
	if self.Next != nil {
		return self.Next.GetWeight(ws)
	}
	return ws
}
func (self *HiddenLevel) SetWeight(ws [][][]float64,num int) error {
	if len(ws)<= num {
		return fmt.Errorf("is long %d %d",len(ws),num)
	}
	self.Wi = ws[num]
	num++
	if self.Next != nil {
		return self.Next.SetWeight(ws,num)
	}
	return nil
}
func (self *HiddenLevel) ShowInfo(num int) {
	fmt.Println(num)
	fmt.Println("wi",self.Wi)
	fmt.Println("ah",self.ah)
	fmt.Println("ci",self.ci)
	if self.Next != nil {
		num--
		self.Next.ShowInfo(num)
	}
}
type syncData struct {
	Var float64
	Num int
}
func (self *HiddenLevel) syncUpdate(X []float64,j int,pool chan int){

	sum := 0.0
	for i,w := range self.Wi{
		sum += X[i] * w[j]
	}
	self.ah[j] = sigmoid(sum)
//	self.ah[j] = softplus(sum)
//	if self.Next == nil {
//		self.ah[j] = sigmoid(sum)
//		self.ah[j] = softplus(sum)
//	}else{
//		self.ah[j] = relu(sum)
//	}
	self.sum[j] = sum

	pool <- j

}
func (self *HiddenLevel) update(X []float64) ([]float64,error) {
	L:=len(self.ah)-1
//	pool := make(chan int,LINE)
	count:=0
	for j:=0; j <= L; j++ {
		if j==L && self.Next != nil {
			self.ah[j]=1
			break
		}
//		go self.syncUpdate(X,j,pool)
		count++
		sum := 0.0
		for i,w := range self.Wi{
			sum += X[i] * w[j]
		}
		self.ah[j] = sigmoid(sum)

	}
//	for i:=0;i < count;i++ {
//		<-pool
//	}
//	close(pool)
	if self.Next != nil {
		return self.Next.update(self.ah)
	}

	return self.ah,nil
}
func (self *HiddenLevel) GetDeltas(tar []float64,j int,pool chan *syncData) {
	Err:=0.0
	for k,t := range tar {
		Err+= t * self.Next.Wi[j][k]
	}
	pool <- &syncData{Var:dsigmoid(self.ah[j]) * Err,Num:j}
//	pool <- &syncData{Var:sigmoid(self.sum[j]) * Err,Num:j}
//	pool <- &syncData{Var:reluder(self.sum[j]) * Err,Num:j}
}
func (self *HiddenLevel) syncDeltas(L int,targets []float64) (d []float64) {

	d = make([]float64,L)
	pool := make(chan *syncData,LINE)
	for j:=0; j<L; j++ {
		go self.GetDeltas(targets,j,pool)
	}
	for i:=0;i<L;i++{
		t := <-pool
		d[t.Num] = t.Var
	}
	close(pool)
	return d
}
func (self *HiddenLevel) backPropagate(targets []float64,N,M float64) ([]float64,error) {
	var output_deltas []float64

	if self.Next == nil {
		output_deltas = make([]float64,len(self.ah))
		if len(targets) != len(self.ah) {
			return nil,fmt.Errorf("targets is err %d %d",len(targets) , len(self.ah))
		}
		for k,tar:= range targets {
//			output_deltas[k] = reluder(self.sum[k])*(tar-self.ah[k])
			output_deltas[k] = dsigmoid(self.ah[k])*(tar-self.ah[k])
//			output_deltas[k] = sigmoid(self.sum[k])*(tar-self.ah[k])
		}
	}else{
//		output_deltas = self.syncDeltas(len(self.ah)-1,targets)
		output_deltas = make([]float64,len(self.ah)-1)
		for j,_ := range output_deltas {
			Err:=0.0
			for k,tar:= range targets {
				Err+= tar*self.Next.Wi[j][k]
			}
			output_deltas[j] = dsigmoid(self.ah[j]) * Err
		}
	}

	if self.Per!= nil {
		self.updateLevel(output_deltas,self.Per.ah,N,M)
		return self.Per.backPropagate(output_deltas,N,M)

	}
	return output_deltas,nil
}
func (self *HiddenLevel) syncUpdateLevel(targets []float64,a,N,M float64,j int,pool chan int ) {
	for k,tar := range targets {
		change := tar * a
		self.Wi[j][k] += N * change + M * self.ci[j][k]
		self.ci[j][k] = change
	}
	pool <- j
}
func (self *HiddenLevel) updateLevel(targets,ah []float64,N,M float64) {

//	pool:=make(chan int,LINE)
	for j,a := range ah{
//		go self.syncUpdateLevel(targets,a,N,M,j,pool)
//		go self.syncUpdateLevel(targets,softplus(a),N,M,j,pool)

		for k,tar := range targets {
			change := tar * a
			self.Wi[j][k] += N * change + M * self.ci[j][k]
			self.ci[j][k] = change
		}
	}
//
//	for i:=0;i<len(ah);i++ {
//		<-pool
//	}
//	close(pool)

}
func (self *HiddenLevel) GetErr(targets []float64) float64 {
	Err:=0.0
	for k,tar := range targets {
		t:=tar - self.ah[k]
		Err+=0.5*t*t
	}
	return Err
}
