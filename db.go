package main
import (
"fmt"
"time"
zk "github.com/samuel/go-zookeeper/zk"
)

type zkDB struct{
    Con     zk.Conn
    Eve     zk.Event
//    isSetup bool
}
/*func New() *zkDB {
	return &zkDB{isSetup: false}
}*/
func (db *zkDB)Login() error {
var err error
//db.Con,db.Eve,err = zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
db.Con,_,err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
if err != nil {
panic(err)
}
//db.Con(Co)
children, stat, ch, err := db.Con.ChildrenW("/")
if err != nil {
panic(err)
}
fmt.Printf("%+v %+v\n", children, stat)
e := <-ch
fmt.Printf("%+v\n", e)

	return nil
}

func main() {
/*z1 = new(zkDB);
z1, _,err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
if err != nil {
panic(err)
}
children, stat, ch, err := z1.ChildrenW("/")
if err != nil {
panic(err)
}
fmt.Printf("%+v %+v\n", children, stat)
e := <-ch
fmt.Printf("%+v\n", e)
*/
var zi zkDB
zi.Login()
}
