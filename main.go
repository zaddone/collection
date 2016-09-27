package main
import (
	"flag"
	"runtime"
	"github.com/zaddone/collection/route"
)
var (
	Port = flag.String("p",":1301","IP port to listen on")
//	InPath = flag.String("i","/home/dimon/pycode/files/win/","in path")
	InPath = flag.String("i","/home/dimon/code/go/data/","in path")
	OutPath = flag.String("o","/home/dimon/sbFiles/","out path")
	Db = flag.String("d","db/","db")
)
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	P := new(route.Route)
	err := P.Init(*InPath,*OutPath,*Db,*Port)
	if err != nil {
		panic(err)
	}
	P.ReadPath()
	P.RunUpdate()
}
