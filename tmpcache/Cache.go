package tmpcache
import (
	"github.com/zaddone/collection/tmpdata"
	"github.com/zaddone/collection/dbscan"
//	"github.com/zaddone/collection/common"
	"fmt"
	"net"
	"time"
//	"encoding/json"
//	"sync"
)
const (
	POOL int = 100
	KeyLen int = 1
)
func  FindSame(list []*dbscan.SED,a *dbscan.SED) bool {
	for _,_a := range list {
		if _a == a {
			return true
		}
	}
	return false
}
//type SedsTmp struct {
//	s *Clu
//	st chan bool
//	d []byte
//}

type Cache struct {

//	SedsMap map[[KeyLen]int]*Clusters
	Sedslist *Clus
//	CluMap map[string]*Clu
//	ConnPool  chan int
//	SedsPool  chan *SedsTmp
//	SedsPool  chan *Clu

	Len int
//	Lis *net.Listener
	same float64
	er float64

}
func (self *Cache) TestInfo(isok bool) {
	self.same ++
	if !isok {
		self.er ++
	}
}
//func (self *Cache) GetCache() *Cache {
//	return self
//}

func (self *Cache) getConn(Lis *net.Listener) net.Conn {
	L := *Lis
	conn,err := L.Accept()
	if err == nil {
		wait := time.Tick(100*time.Millisecond)
		for {
			select {
			case <-wait:
				conn,err := L.Accept()
				if err == nil {
					return conn
				}
			}
		}
	}
	return conn

}
//func (self *Cache) Conns(Lis *net.Listener) {
////	i:=0
//	for {
//		for k,sed:= range self.CluMap {
//			d := sed.GetNnTrData()
//			if d != nil {
//				conn := self.getConn(Lis)
//				go  handleConn( conn,sed,sed.GetStart(),d)
//			}
//			delete(self.CluMap,k)
//		}
//		time.Sleep(time.Second)
//	}
//}
//func handleConn(conn net.Conn,sed *Clu,State chan bool,da []byte) error {
//	defer conn.Close()
//	var begindata [common.SocketLen]byte
//	_,err:=conn.Read(begindata[0:])
//	if err != nil {
//		return err
//	}
//	if begindata[0] != 1 {
//		return fmt.Errorf("head err")
//	}
//	err = common.ConnWrite(conn,da)
//	if err != nil {
//		return err
//	}
//	for {
//		select {
//			case <-State :
////				fmt.Println("kill")
//				_,err = conn.Write([]byte{11})
//				if err != nil {
//					return err
//				}
//				break
//			default :
//
//				_,err = conn.Write([]byte{0})
//				if err != nil {
//					return err
//				}
//				_,err = conn.Read(begindata[0:])
//				if err != nil {
//					return err
//				}
//				if begindata[0] == 1 {
//					_,err = conn.Write([]byte{1})
//					if err != nil {
//						return err
//					}
//					status,err :=common.ConnRead(conn)
//					if err != nil {
//						return err
//					}
//					sed.InitNn(status)
//					return nil
//				}
//		}
//	}
//	return nil
//
//}
//func (self *Cache) SyncServer() {
//	beart := time.Tick(time.Second)
////	i := 0
//	for {
//		select {
//		case <- beart :
//			sy := new(sync.WaitGroup)
//			for _,clus := range self.SedsMap {
//				sy.Add(1)
//				go clus.SyncMergeServer(sy)
//			}
//			sy.Wait()
//		//	i ++
//		//	fmt.Printf("%d\r",i)
//		}
//	}
//}
func (self *Cache)Init(listener *net.Listener) *Cache {

//	self.SedsMap = make(map[[KeyLen]int]*Clusters)
	self.Sedslist = new(Clus)
//	self.CluMap  = make(map[string]*Clu)
//	self.TmpMap = make(map[[KeyLen]int][]*tmpdata.Val)
//	self.ConnPool = make(chan int,POOL)
//	self.SedsPool = make(chan *Clu,POOL)
	self.Len = 10
//	self.Lis = listener
//	for i:=0;i<10;i++{

//	go self.SyncGetConn()
//	go self.Conns(listener)

//	}
//	go self.SyncServer()
	return self

}
func (self *Cache) ShowInfo() string {
	return fmt.Sprintf("%.5f %.1f %.1f",self.er/self.same,self.er,self.same)
}
func GetMapKey(k []int) (key [KeyLen]int) {
	if KeyLen == 1 {
		key[0] = k[len(k) -1]
		return key
	}
	n :=KeyLen/2
	for i:=0; i<n;i++ {
		key[i] = k[i]
	}
	for i,L:= KeyLen -1,len(k) - 1;i>= n; i,L=i-1,L-1 {
		key[i] = k[L]
	}
	return key
}
func GetBL(vals []*tmpdata.Val) (o,e float64) {
	for _,v := range vals {
		if v.Y > 0 {
			o++
		}else{
			e++
		}
	}
	return o,e
}
func (self *Cache) Forecast(v *tmpdata.Val) (int,error) {
	return 0,fmt.Errorf("is nil")

//	key := GetMapKey(v.K)
//	Clus :=self.SedsMap[key]
//	if Clus == nil{
//		return -1,fmt.Errorf("is nil")
//	}
//
//
//	S1,_:= Clus.GetSameSed(v)
//	if S1 == nil {
//		return -1,fmt.Errorf("is nil")
//	}
//	f := S1.Forecast(v)
//	if f == -1 {
//		return -1,fmt.Errorf("check is false")
//	}else {
//		return f,nil
//	}
//	return -1,fmt.Errorf("is err %d",f)

}
func (self *Cache) Input(v *tmpdata.Val) (error) {
//	fmt.Println(v.X)
	self.Sedslist.AppendVal(v,1)
	return nil
//	key := GetMapKey(v.K)
//	Clus :=self.SedsMap[key]
//	if Clus == nil{
//		Clus = new(Clusters)
//		Clus.Init(self)
//		Clus.AppendVal(v)
//		self.SedsMap[key]=Clus
//		return nil
//	}else{
//		Clus.AppendVal(v)
//	}
//	return nil
}
//func (self *Cache) AppendWaits(S1 *Clu) {
//	if len(S1.RawPatterns) < self.Len {
//		return
//	}
//
////	self.CluMap[fmt.Sprintf("%p",S1)] = S1
//
//}
func (self *Cache) AppendWait(S1 *dbscan.SED,isNn bool,vs []*tmpdata.Val) {
//	self.ConnPool <- len(self.ConnPool)
//	fmt.Println("conn",len(self.ConnPool))
//	for i,v := range S1.RawPatterns{
//		fmt.Println(v)
//	}
//	fmt.Println(string(S1.GetValData()))
//	panic("--")
//	S1.State = make(chan bool)
//	go self.GetConn(S1,S1.State,self.ConnPool,isNn,vs)

//	if !FindSame(self.WaitSeds,S1) {
//		self.WaitSeds = append(self.WaitSeds,S1)
//		self.SedsReady ++
//	}

}
