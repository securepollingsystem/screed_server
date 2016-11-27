package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/securepollingsystem/tallyspider/screed"
)

var (
	screedBucket = []byte("screed")

	buckets = [][]byte{
		screedBucket,
	}

	ErrScreedNotFound = errors.New("screed not found")
	ErrScreedExists   = errors.New("A screed already exists for that user")
)

func mustInitDB() *bolt.DB {
	options := &bolt.Options{Timeout: 2 * time.Second}
	dbPath := path.Join(os.Getenv("BOLT_PATH"), "screed.db")

	// Open file
	db, err := bolt.Open(dbPath, 0600, options)
	if err != nil {
		log.Fatalf("Error opening bolt DB: %v", err)
	}

	// Create buckets
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bkt := range buckets {
			_, err := tx.CreateBucketIfNotExists(bkt)
			if err != nil {
				return fmt.Errorf("Error creating bucket `%s`: %v", bkt, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error creating '%v' boltDB bucket: %v\n", screedBucket,
			err)
	}

	return db
}

func GetScreedByPubkey(db *bolt.DB, pubkeyhex string) (s []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(screedBucket)
		pubkeyB := bkt.Get([]byte(pubkeyhex))
		if pubkeyB == nil {
			return ErrScreedNotFound
		}
		s = pubkeyB
		return nil
	})
	return s, err
}

func CreateScreedByPubkey(db *bolt.DB, s *screed.Screed) (pubkeyhex string, err error) {
	if s == nil {
		return "", errors.New("Cannot create nil *Screed!")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(screedBucket)

		if s.VoterPubKey() == nil {
			return errors.New("Cannot create *Screed with nil *VoterPubKey")
		}

		// Pubkey -> compress -> hex-encoded -> []byte
		pubkeyhex = fmt.Sprintf("%x", s.VoterPubKey().SerializeCompressed())
		pubkeyB := []byte(pubkeyhex)

		// Only save if user has no existing screed
		old := bkt.Get(pubkeyB)
		if old != nil {
			return ErrScreedExists
		}

		screedStr, err := s.Serialize()
		if err != nil {
			return err
		}

		return bkt.Put(pubkeyB, []byte(screedStr))
	})

	return pubkeyhex, err
}
