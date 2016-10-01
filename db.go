package zookeeper
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

var (
	DEF_ACL = zk.WorldACL(zk.PermAll)
)

type zkDB struct {
	Con     *zk.Conn
//	Eve     zk.Event
//	cfg     zk.ServerConfig
    cfg     string
	isSetup bool
}

func New() *zkDB {
	return &zkDB{isSetup: false}
}

func (db *zkDB) Login() error {
	var err error
//	db.Con,_,err = zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	db.Con, _, err = zk.Connect([]string{db.cfg}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
/*		children, stat, ch, err := db.Con.ChildrenW("/")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v %+v\n", children, stat)
		e := <-ch
		fmt.Printf("%+v\n", e)
*/	
	return nil
}

func (db *zkDB) IsSetup() bool {
	return db.isSetup
}

func (db *zkDB) Set(Key string, Value string) error {
	globalstatus := make(map[string]string)
	globalstatus[Key] = Value
	if statusbytes, err := json.Marshal(globalstatus); err == nil {
		if _, err := db.Con.Set(Key, statusbytes, -1); err != nil {
			db.Con.Create(Key,statusbytes, 0, DEF_ACL)
		}
	}
	return nil
}

func (db *zkDB) Get(Key string) (string, error) {
	globalstatus := make(map[string]string)
	if globalbytes, _, err := db.Con.Get(Key); err != nil {
		json.Unmarshal(globalbytes, &globalstatus)
	}
	return globalstatus[Key], nil
}

func (db *zkDB) IsDir(Key string) (error, bool) {
	//	resp, err := db.Kapi.Get(db.Ctx, Key, nil)

	globalstatus := make(map[string]string)
	if globalbytes, _, err := db.Con.Get(Key); err != nil {
		json.Unmarshal(globalbytes, &globalstatus)
		return err, false
	}
	/*
		if err != nil {
			return err, false
		}*/
	return nil, true
}

func (db *zkDB) IsKey(Key string) (bool, error) {

	globalstatus := make(map[string]string)
	if globalbytes, _, err := db.Con.Get(Key); err != nil {
		json.Unmarshal(globalbytes, &globalstatus)
	}
//	retStr := make([]string, len(globalstatus))
//	i := 0
	for k := range globalstatus {
		if(Key == k){
		    fmt.Printf("Matched")
		    return true, nil
        }
	}

/*	if by, _, err := db.Con.Get(Key); err != nil {
		//		t.Fatalf("Get failed on node 2: %+v", err)
		fmt.Printf("Get failed on node 2: %+v", err)
	} else if string(by) != "foo-cluster" {
		//		t.Fatal("Wrong data for node 2")
		fmt.Printf("Wrong data for node 2")
	}*/
	return false, nil
}

func (db *zkDB) Update(Key string, Value string, Lock bool) error {
	return nil
}

func (db *zkDB) Del(Key string) error {
	/*	if err := db.Con.Delete(Key, -1); err != nil && err != ErrNoNode {
			fmt.Printf("Delete returned error: %+v", err)
			return err
		}
		_, err := db.Kapi.Delete(db.Ctx, Key, nil)

		if err != nil {
			return err
		}
	*/return nil
}

//CreateSection will create a directory in zookeeper store
func (db *zkDB) CreateSection(Key string) error {
	globalstatus := make(map[string]string)
	globalstatus[Key] = ""
	if statusbytes, err := json.Marshal(globalstatus); err == nil {
		if _, err := db.Con.Set(Key, statusbytes, -1); err != nil {
		    db.Con.Create(Key,statusbytes, 0, DEF_ACL)
		}
	}
	return nil
}

func (db *zkDB) Setup(config string) error {
	var err error
	/*	db.Cfg = cli.Config{
			Endpoints: []string{config},
			Transport: cli.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
	*/
	
    i := strings.Index(config, "//")
    if i > -1 {
        db.cfg = config[i+2:]
    } else {
	    db.cfg = config
    }


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

func (db *zkDB) CleanSlate() error {
	//	_, err := db.Kapi.Delete(db.Ctx, ETC_BASE_DIR, &cli.DeleteOptions{Dir: true, Recursive: true})
	return nil
}

//DeleteSection section will delete a directory optionally delete
func (db *zkDB) DeleteSection(Key string) error {

	//	_, err := db.Kapi.Delete(db.Ctx, Key, &cli.DeleteOptions{Dir: true})
	return nil
}

//ListSection will list a directory
func (db *zkDB) ListSection(Key string, Recursive bool) ([]string, error) {

	globalstatus := make(map[string]string)
	if globalbytes, _, err := db.Con.Get(Key); err != nil {
		json.Unmarshal(globalbytes, &globalstatus)
	}
	//json.Unmarshal(globalbytes, &globalstatus)
	retStr := make([]string, len(globalstatus))
	i := 0
	for k := range globalstatus {
		retStr[i] = k
		i++
	}

	/*	resp, err := db.Kapi.Get(db.Ctx, Key, &cli.GetOptions{Sort: true})
		if err != nil {
			return nil, err
		}
		retStr := make([]string, len(resp.Node.Nodes))
		for i, n := range resp.Node.Nodes {
			retStr[i] = n.Key
		}
	*/
	return retStr, nil
}
