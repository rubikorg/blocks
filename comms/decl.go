package comms

import (
	"bytes"
	"encoding/gob"
	"errors"
	"path/filepath"
	"time"

	"os"

	"github.com/rubikorg/rubik"
	r "github.com/rubikorg/rubik"
	bolt "go.etcd.io/bbolt"
)

const (
	// BlockName describes this rubik block
	BlockName = "BasicMessagePasser"
)

// BlockBasicComm is implementation of basic communicator
// for this basic message passing system we do not need
// any service discovery layer or framework to make it work
// this block should only be used for local services and
// message passing in between them
//
// This uses basic HTTP endpoints to post messages as GOBs
// for faster decoding and easier type translations
// accross services
type BlockBasicComm struct {
	dbConn   *bolt.DB
	services []service
}

type service struct {
	Name        string
	Location    string
	PunchInTime time.Time
}

// OnAttach function for block inteface
func (bbc BlockBasicComm) OnAttach(app *r.App) error {
	var conf map[string]string
	err := app.Decode("basicMsgPasser", &conf)
	if err != nil {
		return err
	}

	// create a database with name given in config.name || create own name
	// inside home/.rubik/ folder as msgp.db
	home, _ := os.UserHomeDir()
	folderPath := filepath.Join(home, ".rubik")

	os.MkdirAll(folderPath, 0666)

	dbPath := filepath.Join(folderPath, "msgp.db")
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return err
	}

	bbc.dbConn = db

	if conf["name"] == "" {
		return errors.New("No `name` key inside basicMsgPasser config")
	}

	myservice := service{
		Name:        conf["name"],
		Location:    app.CurrentURL,
		PunchInTime: time.Now(),
	}

	// punch-in your attendance in db
	// get list of all the other local services .. check if the current service
	// is present this will help determin wether to notify of new service or not
	serviceList, serr := bbc.listAndPunchIn(myservice)
	if serr != nil {
		return serr
	}

	bbc.services = serviceList

	// notify your presence to all other servers by calling /new/service
	// of the client api
	for _, s := range serviceList {
		if s.Name == myservice.Name {
			continue
		}
		rubcl := rubik.NewClient(s.Location, 30*time.Second)
		// when we just want to ping that route saying that a new service has
		// arrived we use blank request entity
		be := r.BlankRequestEntity{}
		be.PointTo = "/_msgp/new/service"
		go rubcl.Get(be)
	}

	// add routes that we need for message passing
	newPunchInRoute.Controller = bbc.newServiceCtl
	listServicesRoute.Controller = bbc.listCtl
	msgpRouter.Add(listServicesRoute)
	msgpRouter.Add(newPunchInRoute)
	r.Use(msgpRouter)

	return nil
}

// Send implements communictor send interface
func (bbc BlockBasicComm) Send(target string, data interface{}) error {
	return nil
}

func (bbc BlockBasicComm) listAndPunchIn(s service) ([]service, error) {
	createBucketIfNotExist(bbc.dbConn)
	services, err := getServiceList(bbc.dbConn)
	if err != nil {
		return nil, err
	}

	// the for loop makes sure that the execution does not transfer to
	// the next block if the service name is already present in services list
	for _, ls := range services {
		if ls.Name == s.Name {
			return services, nil
		}
	}

	// lets append the service
	services = append(services, s)
	uerr := bbc.dbConn.Update(func(tx *bolt.Tx) error {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(services)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("services"))
		perr := b.Put([]byte("list"), buf.Bytes())
		if perr != nil {
			return perr
		}
		return nil
	})

	if uerr != nil {
		return nil, uerr
	}

	return services, nil
}

func init() {
	r.Attach(BlockName, BlockBasicComm{})
}
