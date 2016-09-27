package signal

import (
	"fmt"
	"os"
)

type GraphicSignal struct {

//	DB       *graphic.DBGraphic
	//LastDB   *DBGraphic
	//Cl       *DBCluster
	atlias   *Atlias
	state    bool
	FilePath string

}

func (self *GraphicSignal) Init(at *Atlias) {
//	self.DB = db
	self.atlias = at
}
func (self *GraphicSignal) Create(diff int,kill int) {

//	Wave:=len(self.DB.Trend)
//	tr := self.DB.Trend[Wave-1]
//	L := len(tr)
//	if tr[L-1].Begin != self.atlias.Last {
//		err := fmt.Errorf("%d,%d", int(tr[0].Begin), int(self.atlias.Last))
//		check(err)
//	}
//	diff := 0
//	if Si {
//		diff = 1
//	} else {
//		diff = -1
//	}
	info := fmt.Sprintf("push %d %d %d %d %d %d\n", diff, 0, int(self.atlias.Last), int(self.atlias.Begin), self.atlias.Msc,kill)
	fmt.Println(info)
	self.WriteFile(info)
}

func (self *GraphicSignal) WriteFile(str string) {
	f, err := os.OpenFile(self.FilePath, os.O_SYNC|os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	check(err)
	defer f.Close()
	_, err = f.WriteString(str)
	check(err)
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
