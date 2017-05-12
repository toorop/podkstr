package core

import (
	"io"
	"log"

	"github.com/spf13/viper"
	"github.com/toorop/gopenstack/context"
	"github.com/toorop/gopenstack/identity"
	"github.com/toorop/gopenstack/objectstorage/v1"
	"github.com/toorop/podkstr/logger"
)

// Storer is an interface for storage drivers
type Storer interface {
	Get(key string) (io.Reader, error)
	Put(key string, reader io.ReadSeeker) error
	Del(key string) error
}

// Store the store used by app
var Store Storer

////////////////////
// OpenStack swift
// for https://podkstr.com https://www.ovh.com/fr/public-cloud/storage/object-storage/

// OsStore is an Openstack Swift Storer
type OsStore struct {
}

// InitOsStore init Store as an OpenStack Storer
func InitOsStore() error {
	if err := context.InitKeyring(viper.GetString("openstack.user"), viper.GetString("openstack.password"), viper.GetString("openstack.tenant.name"), viper.GetString("openstack.authurl"), viper.GetString("openstack.tenant.id"), viper.GetString("openstack.region")); err != nil {
		return err
	}

	if err := identity.DoAuth(); err != nil {
		return err
	}

	// auto update Token each 30 minutes
	logger := log.New(logger.Log.Writer(), "osstore", 0)
	Store = Storer(OsStore{})
	///store = &storee
	identity.AutoUpdate(30, logger)
	return nil
}

// Get implements Storer.Get
func (o OsStore) Get(key string) (reader io.Reader, err error) {
	return
}

// Put implements Storer.Put
func (o OsStore) Put(key string, reader io.ReadSeeker) (err error) {
	object := &objectstorageV1.Object{
		Region:    viper.GetString("openstack.region"),
		Container: viper.GetString("openstack.container.name"),
		Name:      key,
		RawData:   reader,
	}
	return object.Put(nil)
}

// Del implements Storer.Del
func (o OsStore) Del(key string) (err error) {
	object := &objectstorageV1.Object{
		Region:    viper.GetString("openstack.region"),
		Container: viper.GetString("openstack.container.name"),
		Name:      key,
	}
	log.Println("OBJECT: ", object)
	return object.Delete(false)

}
