package comms

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"os"

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
	serviceList, err := bbc.listAndPunchIn(myservice)
	if err != nil {
		// fmt.Println("Cannot proceed with the execution of Block:", BlockName, err.Error())
		panic(err)
		// return nil
	}
	bbc.services = serviceList

	// notify your presence to all other servers by calling /new/service
	// of the client api

	return nil
}

// Send implements communictor send interface
func (bbc BlockBasicComm) Send(target string, data interface{}) error {
	return nil
}

func (bbc BlockBasicComm) listAndPunchIn(s service) ([]service, error) {
	var services []service
	err := bbc.dbConn.Update(func(tx *bolt.Tx) error {
		nameKey := []byte("services")
		listKey := []byte("list")

		b := tx.Bucket(nameKey)
		if b == nil {
			var err error
			b, err = tx.CreateBucket(nameKey)
			if err != nil {
				return err
			}
		}

		listb := b.Get(listKey)
		// if there is nothing inside list insert the incoming service
		if listb == nil {
			services = append(services, s)

			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			err := enc.Encode(services)
			if err != nil {
				return err
			}
			err = b.Put(listKey, buf.Bytes())
			if err != nil {
				return err
			}
			return nil
		}

		// if list bytes are not nil then it means that we have some
		// services punched in before
		listBuf := bytes.NewBuffer(listb)
		dec := gob.NewDecoder(listBuf)
		derr := dec.Decode(&services)
		if derr != nil {
			return derr
		}
		fmt.Println("coming?")

		// the for loop makes sure that the execution does not come here
		// if the service name is already present in services list
		for _, ls := range services {
			if ls.Name == s.Name {
				return nil
			}
		}

		// lets append the service
		services = append(services, s)
		bbc.dbConn.Update(func(tx *bolt.Tx) error {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			err := enc.Encode(services)
			if err != nil {
				return err
			}
			b := tx.Bucket([]byte("services"))
			b.Put([]byte("list"), buf.Bytes())
			return nil
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}

func init() {
	addRoutes()
	r.Use(msgpRouter)

	r.Attach(BlockName, BlockBasicComm{})
}
