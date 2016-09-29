package main

import (
	"encoding/json"
	"fmt"
	zk "github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
)

type zkDB struct {
	Con     *zk.Conn
	Eve     zk.Event
	isSetup bool
}

func New() *zkDB {
	return &zkDB{isSetup: false}
}

func (db *zkDB) IsSetup() bool {
	return db.isSetup
}

func (db *zkDB) Set(Key string, Value string) error {
    //TODO get the value of status_global-path
	//	_, err := db.Kapi.Set(db.Ctx, Key, string(Value), nil)
	//    _, err := db.Con.Set(Key,[]byte{},-1)
	globalstatus := make(map[string]string)
	globalstatus[Key] = Value
	if statusbytes, err := json.Marshal(globalstatus); err == nil {
		if _, err := db.Con.Set(status_global, statusbytes, -1); err == zk.ErrNoNode {
			db.Con.Create(status_global, statusbytes, DEF_FLAGS, DEF_ACL)
		}
	}

	return err

}

/*
//TODO
func (db *etcdDB) Get(Key string) (string, error) {

	resp, err := db.Kapi.Get(db.Ctx, Key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (db *etcdDB) IsDir(Key string) (error, bool) {
	resp, err := db.Kapi.Get(db.Ctx, Key, nil)

	if err != nil {
		return err, false
	}
	return nil, resp.Node.Dir
}
*/

func (db *zkDB) IsKey(Key string) (bool, error) {
	if by, _, err := db.Con.Get(Key); err != nil {
		fmt.Printf("Get failed on node 2: %+v", err)
	} else if string(by) != "foo-cluster" {
		fmt.Printf("Wrong data for node 2")
	}
	return true, nil
}

//CreateSection will create a directory in zookeeper store
func (db *zkDB) CreateSection(Key string) error {
	//	_, err := db.Kapi.Set(db.Ctx, Key, "", &cli.SetOptions{Dir: true, PrevExist: cli.PrevNoExist})
	_, err := db.Con.Set(Key, []byte{}, -1)
	//    fmt.Printf("Create Section")
	if err != nil {
		return err
	}
	return nil
}

func (db *zkDB) Login() error {
	var err error
	//TODO Get the ip from config file
	//&db.Con,_,err = zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	db.Con, _, err = zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	//TODO need to remove children according to the need
	children, stat, ch, err := db.Con.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)

	return nil
}

func (db *zkDB) Setup() error {
	var err error
	/*	db.Cfg = cli.Config{
			Endpoints: []string{config},
			Transport: cli.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
	*/
	err = db.Login()
	if err != nil {
		return err
	}

	err = db.CreateSection(ETC_BASE_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	err = db.CreateSection(ETC_INST_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	err = db.CreateSection(ETC_CONF_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	db.isSetup = true
	return nil
}

func main() {
	var zi zkDB
	zi.Setup()
	//zi.isKey("nodeName")
	if ok, _ := zi.IsKey("nodeName"); !ok {
		//return nil
		//fmt.Printf("Shivam")
	}

	//zi.Login()
	//zi.CreateSection(ETC_BASE_DIR)
}
