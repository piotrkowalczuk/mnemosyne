package mnemosyned

import (
	"github.com/boltdb/bolt"
	"fmt"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"time"
	"encoding/json"
)

type boltStorage struct {
	db	*bolt.DB
	bucketName []byte
}

func newBoltStorage(db *bolt.DB, bucketName []byte) storage {
	return &boltStorage{
		db:	db,
		bucketName: bucketName,
	}
}

func (bs *boltStorage) Setup() error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bs.bucketName)
		if err != nil {
			return err
		}

		return nil
	})
}

func (bs *boltStorage) TearDown() error {
	return bs.db.Update(func (tx *bolt.Tx) error {
		return tx.DeleteBucket(bs.bucketName)
	});
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
	});

	return ent.session()

}

func (bs *boltStorage) Abandon(string) (bool, error) {

}

func (bs *boltStorage) Get(accessToken string) (*mnemosynerpc.Session, error) {
	var b []byte
	var se sessionEntity

	return bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bs.bucketName)
		b = bucket.Get(accessToken)

		copy(se, b)

		return nil
	})

	err := json.Unmarshal(b, &se)
	if err != nil {
		return nil, err
	}

	return se, nil
}

func (bs *boltStorage) List(int64, int64, *time.Time, *time.Time) ([]*mnemosynerpc.Session, error) {

}

func (bs *boltStorage) Exists(string) (bool, error) {

}

func (bs *boltStorage) Delete(string, *time.Time, *time.Time) (int64, error) {

}

func (bs *boltStorage) SetValue(string, string, string) (map[string]string, error) {

}

