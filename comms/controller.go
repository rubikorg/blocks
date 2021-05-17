package comms

import (
	"bytes"
	"encoding/gob"

	r "github.com/rubikorg/rubik"
	bolt "go.etcd.io/bbolt"
)

func createBucketIfNotExist(conn *bolt.DB) {
	conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			tx.CreateBucket([]byte("services"))
		}
		return nil
	})
}

func getServiceList(conn *bolt.DB) ([]service, error) {
	var services []service
	err := conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		listb := b.Get([]byte("list"))
		// if this is the first time then return no error
		// and empty service list
		if listb == nil {
			return nil
		}

		buf := bytes.NewBuffer(listb)
		dec := gob.NewDecoder(buf)
		err := dec.Decode(&services)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (bbc BlockBasicComm) listCtl(req *r.Request) {
	services, err := getServiceList(bbc.dbConn)

	if err != nil {
		req.Throw(500, err, r.Type.JSON)
	}
	req.Respond(services, r.Type.JSON)
}

func (bbc BlockBasicComm) newServiceCtl(req *r.Request) {
	services, err := getServiceList(bbc.dbConn)

	if err != nil {
		req.Throw(500, err, r.Type.JSON)
	}

	// we just query the db and save the new list to our bbc instance
	bbc.services = services
	req.Respond("success")
}
