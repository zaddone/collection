package route
import (
	"fmt"
)
type DbServer struct {
	ip string
	port  int
}
func (self *DbServer) Init(ip string,port int) *DbServer {
	self.ip = ip
	self.port = port
	return self
}
func (self *DbServer) GetAddrStr() string {

	return fmt.Sprintf("%s:%d",self.ip,self.port)

}
