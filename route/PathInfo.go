package route
import (
	"os"
	"fmt"
	"time"
	"github.com/zaddone/collection/block"
	"github.com/zaddone/collection/signal"
	"github.com/zaddone/collection/tmpcache"
	"strconv"
	"strings"
	"path/filepath"
	"net"
)
type Blocks interface{
	Init(listener *net.Listener)
	SetSignal(path string)
	SetGraphic(Week int, Hour int, Name string)
	Analysis(at *signal.Atlias)
	SetUpdate(b bool)
	GetCache() *tmpcache.Cache
}

type PathInfo struct {
	Size     int64
	SizeDiff int64
	Path     string
	Paths    []string
	LastTime int64
	Week     int
	Hour     int
	Name string
//	tr    *trend.Trend
	tr    Blocks
	Date  time.Time
}
func (self *PathInfo) GetTrend() Blocks {
	return self.tr
}
func (self *PathInfo) SetTrend(tr Blocks)  {
	self.tr= tr
}

func (self *PathInfo) UpdateMT(path string,listener *net.Listener){

//	db := new(block.Block)
	db := new(block.Feature)
	db.Init(listener)
	db.SetSignal(filepath.Join(path,self.Name))
//	db.InitConn(host.GetAddrStr())
	self.tr = db

}
func (self *PathInfo) Init(path string, fi os.FileInfo){

	var err error
	if fi==nil {
		fi,err = os.Stat(path)
		if err!= nil {
			fmt.Println(err)
		}
	}
	self.Hour, err = strconv.Atoi(fi.Name())
	if err != nil {
		fmt.Println(path,fi)
		panic(err)
	//	return err
	}
	//self.Hour = hour
	self.Size = fi.Size()
	self.Path = path
	self.LastTime = time.Now().Unix()
	self.Paths =strings.Split(path, "/")
	L := len(self.Paths)
	if L<3 {
		self.Paths =strings.Split(path, "\\")
		L = len(self.Paths)
	}
	if L<3 {
		panic("path err")
	}
	self.Date, self.Week = GetWeek(self.Paths[L-3])
	self.Name = self.Paths[L-2]

}
func (self *PathInfo) Comp(p *PathInfo) bool {
	if self.Date.Before( p.Date) {
		return true
	}else if self.Date == p.Date {
		if self.Hour < p.Hour {
			return true
		}else if self.Hour == p.Hour {
			p.SizeDiff = self.Size
			return true
		}
	}
	return false
}
func (self *PathInfo) UpdateFile() (int,error) {

	fileInfo, err := os.Stat(self.Path)
	if err!=nil {
		return 0,err
	}
	size := fileInfo.Size()
	sizeDiff := int(size-self.Size)
	if sizeDiff==0 {
		return sizeDiff,nil
	}
	buffer := make([]byte,sizeDiff)
	err = self.ReadFild(self.Size,buffer)
	if err!=nil {
		return sizeDiff,err
	}
	self.Size = size
	self.ReadBytes(buffer)
	self.LastTime = time.Now().Unix()
	return sizeDiff,nil
}
func (self *PathInfo) StartFile() error {

//	var size int64 = 0
//	sizeDiff:=self.size
	self.tr.SetGraphic(self.Week, self.Hour, self.Name)
	buffer := make([]byte,self.Size)
	err := self.ReadFild(self.SizeDiff,buffer)
	if err!=nil {
		return err
	}
//	fmt.Println("buffer",len(buffer))
	self.ReadBytes(buffer)
	return nil

}
func (self *PathInfo) ReadFild(Size int64, buffer []byte) error {

	fi, err := os.Open(self.Path)
	defer fi.Close()
	if err != nil {
		return err
	}
	if Size > 0 {
		fi.Seek(Size, 0)
	}
	_, err = fi.Read(buffer)
	if err != nil {
//		fmt.Println(err)
		return err
	}
	return nil

}
func (self *PathInfo) ReadBytes(byt []byte) {
	//	oldtime := time.Now().UnixNano()
	//	oldMap := len(Map.Tag)
	lastIndex := 0
	for i, b := range byt {
		if b == byte('\n') {
			at := new(signal.Atlias)
			err := at.GetAtlas(byt[lastIndex:i])
			if err != nil {
				fmt.Println(err)
			}else{
			//fmt.Println(at)
				self.tr.Analysis(at)
			}
			lastIndex = i + 1
		}
	}
	//	fmt.Println(time.Now().UnixNano()-oldtime, " ", len(Map.Tag), " ", oldMap)
}
func GetWeek(day string) (time.Time,int) {
	t, err := time.Parse("20060102", day)
	if err != nil {
		panic(err)
	}
	return t,int(t.Weekday())
}
func (self PathInfo) GetNext() (*PathInfo, error) {
	//lastPath := PI.path
	L:=len(self.Paths)
	//hour, _ := strconv.Atoi(lastPath[L-1])
	self.Hour++
	if self.Hour > 23 {
//		self.Hour = 0
		oldP :=self.Paths[L-3]
		self.Paths[L-3],self.Date, self.Week = NewDate(oldP)
//		fmt.Println(filepath.Dir(self.Path),oldP,self.Paths[L-3])
		di :=filepath.Dir(self.Path)
		dir:=strings.Replace(di,oldP,self.Paths[L-3],-1)
		if di == dir {
			return nil,fmt.Errorf("Fount Not 2")
		}
//		self.Path = filepath.Join(dir,fmt.Sprintf("%02d", self.Hour))
//		self.UpdateMT()
//		Path := strings.Join(self.Paths[:L-1], "/")
//		Path := filepath.Join(self.Paths[:L-1]...)
//		fmt.Println(dir)
		var fiList []*os.FileInfo
		err:=filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			if nil == fi{
				return err
			}
			if fi.IsDir() {
				return nil
			}
			fiList = append(fiList,&fi)
			return nil

		})
		if err != nil {
			return nil,fmt.Errorf("Fount Not 2")
		}
		if fiList == nil {
			return nil,fmt.Errorf("Fount Not 2")
		}
		fi:= *fiList[0]
		fn :=fi.Name()
		self.Path = filepath.Join(dir,fn)
//		fmt.Println(self.Path)
//		panic(0)
		self.Paths[L-1] = fn
		self.Size = fi.Size()
		self.Hour,err =strconv.Atoi(fn)
		if err != nil {
			panic(err)
		}
		return &self,nil
	}else{
		file := fmt.Sprintf("%02d", self.Hour)
		self.Paths[L-1] = file
		self.Path = filepath.Join(filepath.Dir(self.Path),file)
		fileInfo, err := os.Stat(self.Path)
		if err != nil {
//			panic(err)
	//		fmt.Println(err)
			return nil, err
		}
		self.Size = fileInfo.Size()
		return &self,nil
	}
	return nil,fmt.Errorf("Same err")

}
func NewDate(lastPath string) (string,time.Time, int) {

	t, err := time.Parse("20060102", lastPath)
	if err != nil {
		panic(err)
	}

	t = t.Add(24 * time.Hour)
	week := t.Weekday()
	if week == 6 {
		t = t.Add(24 * 2 * time.Hour)
		week = 1
	} else if week == 0 {
		t = t.Add(24 * time.Hour)
		week = 1
	}
//	week = t.Weekday()
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day()),t, int(week)

}
