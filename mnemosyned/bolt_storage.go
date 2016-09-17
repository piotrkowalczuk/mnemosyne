package mnemosyned

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"time"
)

type boltStorage struct {
	db         *bolt.DB
	bucketName []byte
}

func newBoltStorage(db *bolt.DB, bucketName []byte) storage {
	return &boltStorage{
		db:         db,
		bucketName: bucketName,
	}
}

func (bs *boltStorage) Setup() error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bs.bucketName)
		if err != nil {
			return err
		}

		return nil
	})
}

func (bs *boltStorage) TearDown() error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bs.bucketName)
	})
}

func (bs *boltStorage) Start(sid, sc string, b map[string]string) (*mnemosynerpc.Session, error) {
	at, err := mnemosynerpc.RandomAccessToken(tmpKey)
	if err != nil {
		return nil, err
	}

	ent := &sessionEntity{
		AccessToken:   at,
		SubjectID:     sid,
		SubjectClient: sc,
		Bag:           bag(b),
	}

	bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.bucketName)
		b, err := json.Marshal(ent)
		if err != nil {
			return err
		}

		return b.Put([]byte(at), b)
	})

	return ent.session()

}

func (bs *boltStorage) Abandon(accessToken string) (bool, error) {
	bs.db.Update(func(tx *bolt.Tx) {
		err := tx.Bucket(bs.bucketName).Delete(accessToken)
		if err != nil {
			return false, err
		}

		return true, nil
	})
}

func (bs *boltStorage) Get(accessToken string) (*mnemosynerpc.Session, error) {
	var b []byte
	var se sessionEntity
	var err error

	err = bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bs.bucketName)
		b = bucket.Get(accessToken)

		copy(se, b)

		return nil
	})

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &se)
	if err != nil {
		return nil, err
	}

	return se, nil
}

func (bs *boltStorage) List(offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*mnemosynerpc.Session, error) {
	return errors.New("@todo")
}

func (bs *boltStorage) Exists(accessToken string) (bool, error) {
	var exists bool
	var err error

	err = bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bs.bucketName)
		b = bucket.Get(accessToken)
		exists = b != nil

		return nil
	})

	if err != nil {
		return err
	}
}

func (bs *boltStorage) Delete(accessToken string, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	// todo expiredAtFrom + expiredAtTo

	if accessToken == "" {
		return 0, errors.New("session cannot be deleted, no accessToken provided")
	}

	if expiredAtFrom != nil || expiredAtTo != nil {
		return 0, errors.New("expiredAtFrom and exoiredAtTo not supported")
	}

	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bs.bucketName)
		return b.Delete(accessToken)
	})
}

func (bs *boltStorage) SetValue(accessToken string, key, value string) (map[string]string, error) {
	if accessToken == "" {
		return nil, errMissingAccessToken
	}

	return bs.db.Update(func(tx *bolt.Tx) {
		var se sessionEntity
		var err error
		b, err := bs.Get(accessToken)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &se)
		if err != nil {
			return err
		}

		se[key] = value

		b, err = json.Marshal(se)

		if err != nil {
			return err
		}
	})
}
