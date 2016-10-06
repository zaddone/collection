package route
import (
	"os"
//	"log"
	"path/filepath"
	"fmt"
	"time"
	"net"
//	"strconv"
//	"strings"
	"github.com/zaddone/collection/tmpcache"
	"encoding/json"
)
type Info struct {
	Pi *PathInfo
	start  chan *PathInfo
	stop   chan bool
}
func (self *Info) init(){
	self.start = make(chan *PathInfo,50)
	self.stop = make(chan bool)
}
type Route struct {

	Path string
	OutPath string

	LogPath string
	dbPath string

	LastInfo     map[string]*Info
	IsUpdate  bool

	listener  net.Listener

}
func (self *Route) Init (path string, OutPath string,lastFile string,Port string) error {
	self.IsUpdate = false
	var err error
	self.Path,err = filepath.Abs(path)
	if err != nil {
		return err
	}
	self.OutPath,err = filepath.Abs(OutPath)
	if err != nil {
		return err
	}

	self.LogPath,err = filepath.Abs(lastFile)
	if err != nil {
		return err
	}
	self.dbPath = filepath.Join(self.LogPath,"data")

	self.LogPath = filepath.Join( self.LogPath,"lastInfo.log")
	self.LastInfo = make(map[string]*Info)
	err=self.SetLastInfo()
	if err!= nil {
		fmt.Println(err)
	}

	self.listener,err = net.Listen("tcp",":1301")
	if err != nil {
		panic(err)
	}
	go self.Console()
	go self.bakServer()
	return nil
}
func (self *Route) SetLastInfo() error {
	f,err := os.Open(self.LogPath)
	if err != nil {
		return err
	}
	de := json.NewDecoder(f)
	err  = de.Decode(self.LastInfo)
	if err != nil {
		return err
	}
	defer f.Close()

//	bys,err:=self.ReadSaveFile()
//	if err!= nil {
//		return err
//	}
//	for _,b:=range bys{
//		pi := new(PathInfo)
//		pi.Init(b,nil)
//		self.appendPathInfo(pi)
//	}
	err = self.ReadCacheDB()
	if err!= nil {
		return err
	}
	return nil
}
func (self *Route) ReadPath() {
	filepath.Walk(self.Path, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		pi := new(PathInfo)
		pi.Init(path,fi)
		if pi.Name != "USDJPY" {
			return nil
		}

//		fmt.Println(path)
		last:=self.appendPathInfo(pi)
//		self.SetPathInfoCache(pi)
		if last!=nil {
		last.start <- pi
		}
		return nil
	})
	self.endSync()
	self.UpdateLastFile()
}
func (self *Route) bakServer() {
	beat := time.Tick(30*time.Minute)
	for {
		select {
		case <-beat:
			self.UpdateLastFile()
			self.BakCache()
		}
	}
}
func (self *Route) BakCache() {
	for name,last := range self.LastInfo {
		ca:=last.Pi.tr.GetCache()
		db := ca.Sedslist.Getbakdb()
		self.SaveCacheDB(filepath.Join(self.dbPath,name),db)

	}

}
func (self *Route) Console() {
	var str string
//	beat := time.Tick(30*time.Minute)
	for {
		fmt.Scanf("%s\n", &str)
		fmt.Println(str)
		self.QueryCache(str)
	}
}

func (self *Route) SaveCacheDB(name string,db []byte) {
	f,e:=os.Create(name)
	if e != nil {
		fmt.Println(e)
		return
	}
	_,e=f.Write(db)
	if e != nil {
		fmt.Println(e)
		return
	}
}
func (self *Route) ReadCacheDB() (err error){

	return filepath.Walk(self.dbPath, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
//		ns:=strings.Split(fi.Name(),"_")
		last:=self.LastInfo[fi.Name()]
		if last == nil {
			return fmt.Errorf("Fount not info")
		}
		ca:=last.Pi.tr.GetCache()
//		k,err:=strconv.Atoi(ns[1])
//		if err != nil {
//			return err
//		}
		sed:=ca.Sedslist
//		sed:=ca.SedsMap[[1]int{k}]
//		if sed == nil {
//			sed = new(tmpcache.Clusters)
//			sed.Init(ca)
//			ca.SedsMap[[tmpcache.KeyLen]int{k}] = sed
//		}
		f,err := os.Open(path)
		if err != nil {
			return err
		}
		db:=make([]byte,fi.Size())
		_,err = f.Read(db)
		if err != nil {
			return err
		}
		defer f.Close()
		return sed.LoadData(db)
	})
}

func (self *Route) QueryCache(name string) {

	last:=self.LastInfo[name]
	if last == nil {
		fmt.Println("find nil")
		return
	}
	sed:=last.Pi.tr.GetCache().Sedslist
	kcon:=0
	con:=0
	sedLen:=10
	seds := make([]*tmpcache.Cl,sedLen)
	for _,c := range sed.Clu {
		L := len(c.RawPatterns)
		con += L
		if L > 15 {
			if sedLen > 0 {
				sedLen -- 
				seds[sedLen]=c
			}
			kcon+=L
		}
	}
	fmt.Printf("%d %d %d %d\r\n",sed.ValCount,con,kcon,len(sed.Clu))
	var s1 string
	for {
		fmt.Println("Wait input")
//		fmt.Scanf("%s %s\n",&s1,&s2)
		fmt.Scanf("%s\n",&s1)
		if s1 == "x" {
			self.QueryCache(name)
			return
		}
		if s1 == "xx" {
			break
		}
		for _,cl := range seds {
			if cl != nil {
				fmt.Println(cl.CountY,cl.GetPG())
			}
		}
	}
	return
//	seds := make(map[[1]int][]*tmpcache.Clu)
//	for k,v:=range ca.SedsMap {
////		seds[k] := new(tmpMsg)
//		var count,countN int
//		var db []byte
//		seds[k],count,countN,db = v.GetCountSeds()
//		fmt.Println(k,len(v.Clu),len(seds[k]),count,countN)
//		self.SaveCacheDB(filepath.Join(self.dbPath,fmt.Sprintf("%s_%d",name,k[0])),db)
//	}
//	var key [1]int
//	var s1 string
////	var s2 string
//	var err error
//	for {
//		fmt.Println("Wait input")
////		fmt.Scanf("%s %s\n",&s1,&s2)
//		fmt.Scanf("%s\n",&s1)
//		if s1 == "x" {
//			self.QueryCache(name)
//			return
//		}
//		if s1 == "xx" {
//			break
//		}
//		key[0],err = strconv.Atoi(s1)
//		if err !=nil {
//			fmt.Println(err)
//			continue
//		}
////		key[1],err = strconv.Atoi(s2)
////		if err != nil {
////			fmt.Println(err)
////			continue
////		}
//		vs :=seds[key]
//		if vs == nil {
//			fmt.Println("found not")
//			continue
//		}
//		for i,sed := range vs {
//			if sed != nil {
////				sed.Equs()
//				fmt.Printf("%d %d %.1f %.1f\r\n",i,len(sed.RawPatterns),sed.Counts,sed.CountsY)
//			}
//		}
////		fmt.Println(key,len(vs.Seds))
//	}
}
func (self *Route) appendPathInfo(pi *PathInfo) *Info{

	last:=self.LastInfo[pi.Name]
	if last==nil {
		last = new(Info)
		last.init()
	//	host:=self.InfoRing.Value.(*DbServer)
		pi.UpdateMT(self.OutPath,&self.listener)
		last.Pi = pi
		self.LastInfo[pi.Name] =last
//		self.InfoRing = self.InfoRing.Next()
		go self.SyncRead(last.start,last.stop)
	}else{
		//p := last.Pi
		tr:=last.Pi.GetTrend()
		if tr== nil {
			pi.UpdateMT(self.OutPath,&self.listener)
		}else{
			pi.SetTrend(tr)
		}
//		pi.MT = last.Pi.MT
		if last.Pi.Comp(pi){
			last.Pi = pi
		}else{
			return nil
		}
//		if p.date.Before( pi.date) {
//			last.Pi = pi
//		}else if p.date == pi.date && p.hour < pi.hour {
//			last.Pi = pi
//		}
	}

	return last
}
func (self *Route) endSync(){

	for _,pi:=range self.LastInfo {
//		pi.start  <- nil
		close(pi.start)
		<-pi.stop
	}

}
func  (self *Route)SyncRead(start chan *PathInfo,stop chan bool) {

	for {
		it,ok:= <-start
		if !ok || it == nil{
			stop<-true
			return
		}
		err:=it.StartFile()
		if err!=nil {
			fmt.Println(err)
		}
//		log.Fatal(it.path)
	}

}

func (self *Route) RunUpdate(){
	self.IsUpdate = true
	h:=make(chan bool)
	for k,v := range self.LastInfo {
		fmt.Println(k,v.Pi)
		go self.SyncUpdate(v)
	}
	<-h
}
func (self *Route) SyncUpdate(info *Info) {
	heartbeat:=time.Tick(100*time.Millisecond)
	for {
		select {

		case <-heartbeat:
			pis,err:=self.GetLastData(info.Pi)
			if err!=nil {
//				fmt.Println(err)
			}else{
//				if pis!=info.Pi {
			//		panic(0)
//					fmt.Println("New hour ",pis)
//				}
				fmt.Println("New hour ",pis)
				info.Pi=pis
			}
		}
	}

}
func (self *Route) GetLastData(LastFile *PathInfo) (*PathInfo,error) {

	LastFile.GetTrend().SetUpdate(true)
	sizeDiff,err:=LastFile.UpdateFile()
	if err != nil {
		fmt.Println(err)
		return LastFile,err
	}
	if sizeDiff > 0 {
	//	fmt.Println(sizeDiff)
		return LastFile,nil
	}
//	if !OutTime(LastFile.LastTime) {
//		return LastFile,fmt.Errorf("info is not update")
//	}
	NewPi, err := LastFile.GetNext()
	if err != nil {
		return LastFile,err
	}

//	self.WritePathFile(NewPi.Path,nil)
	self.appendPathInfo(NewPi)
	err = NewPi.StartFile()
	if err != nil {
		fmt.Println(err)
		return NewPi,err
	}
//	NewPi.MT.IsUpdate = false
	return NewPi,nil

}
func OutTime(oldTime int64) bool {

	if time.Now().Unix()-oldTime > 60 {
		return true
	}
	return false

}
func (self *Route)ReadSaveFile() (bufs []string,err error) {

	fileInfo, err := os.Stat(self.LogPath)
	if err != nil {
		return nil,err
	}
	fi, err := os.Open(self.LogPath)
	defer fi.Close()
	if err != nil {
		return nil,err
	}
	buffer:=make([]byte,fileInfo.Size())
	_, err = fi.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	lastIndex := 0
	for i, b := range buffer {
		if b == byte('\n') {
			bufs=append(bufs,string(buffer[lastIndex:i]))
			lastIndex = i + 1
		}
	}
	return bufs,nil

}
func (self *Route)UpdateLastFile(){
	f,err := os.Create(self.LogPath)
	defer f.Close()
	if err!= nil {
		panic(err)
	}
	en := json.NewEncoder(f)
	err = en.Encode(self.LastInfo)
	if err!= nil {
		panic(err)
	}

//	f,err := os.OpenFile(self.LogPath,os.O_SYNC|os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC,0777)
//	if err!= nil {
//		panic(err)
//	}
//	for _,v:=range self.LastInfo{
//		self.WritePathFile(v.Pi.Path,f)
//	}

}
func (self *Route)WritePathFile(str string,f *os.File) error {

	var err error
	if f==nil {
		f,err = os.OpenFile(self.LogPath,os.O_SYNC|os.O_RDWR|os.O_CREATE|os.O_APPEND,0777)
		if err!=nil {
			return err
		}
		defer f.Close()
	}
	_,err = f.WriteString(fmt.Sprintf("%s\n",str))
	return err

}
