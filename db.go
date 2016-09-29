package main
import (
"fmt"
"time"
zk "github.com/samuel/go-zookeeper/zk"
)

const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
)

type zkDB struct{
    Con     zk.Conn
    Eve     zk.Event
    isSetup bool
}
func New() *zkDB {
	return &zkDB{isSetup: false}
}

//CreateSection will create a directory in zookeeper store
func (db *zkDB) CreateSection(Key string) error {

//	_, err := db.Kapi.Set(db.Ctx, Key, "", &cli.SetOptions{Dir: true, PrevExist: cli.PrevNoExist})

    _,err := zk.(*Conn).Set(Key,"",-1)

	if err != nil {
		return err
	}
	return nil
}

func (db *zkDB)Login() error {
var err error
//&db.Con,_,err = zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
Co,_,err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
if err != nil {
panic(err)
}
db.Con=*Co
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
var zi zkDB
zi.Login()
zi.CreateSection(ETC_BASE_DIR)
}
