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

func (bbc BlockBasicComm) listCtl(en interface{}) r.ByteResponse {
	services, err := getServiceList(bbc.dbConn)

	if err != nil {
		return r.Failure(500, err, r.Type.JSON)
	}

	return r.Success(services, r.Type.JSON)
}

func (bbc BlockBasicComm) newServiceCtl(en interface{}) r.ByteResponse {
	services, err := getServiceList(bbc.dbConn)

	if err != nil {
		return r.Failure(500, err, r.Type.JSON)
	}

	// we just query the db and save the new list to our bbc instance
	bbc.services = services
	return r.Success("success")
}
