//
// A boltdb wrapper to encrypt and decrypt the values stored in the boltdb via AES Cryptor, and also provides 
// db operations functions
//
package boltsec

import (
	"os"
	"log"
	"errors"
	"bytes"
	"github.com/boltdb/bolt"
	"path/filepath"
	"encoding/json"
)

type boltsecDB struct {
	*bolt.DB
}

type boltsecTx struct {
	*bolt.Tx
}

// The DBManager struct, all fields are not needed to be accessed by other packages
type DBManager struct {
	name string
	path string
	fullPath string
	secret string
	buckets []string
	batchMode bool
	cryptor * aesCryptor 
	db * boltsecDB
}

var Debug = false 
var Logger = log.New(os.Stdout, "[DB] ", log.LstdFlags)

// The key error messages generated in the package
var (
	ErrFileNameInvalid       = errors.New("invalid file name")
	ErrPathInvalid       	 = errors.New("invalid path name")
	ErrKeyInvalid			 = errors.New("invalid key or key is nil")
)

// The main function to initialize the the DB manager for all DB related operations 
// 	name: the db file name, such as mydb.dat, mytest.db
// 	path: the db file's path, can be "" or any other director
// 	secret: the secret value if you want to encrypt the values; if you don't want to encrypt the data, simply put it as ""
// 	batchMode: to control whether to close the db file after each db operation
// 	buckets: the buckets in the db file to be initialized if the db file does not existed
func NewDBManager(name, path, secret string, batchMode bool, buckets []string) (dbm * DBManager, err error) {
	var info os.FileInfo
	if path != "" {
		info, err = os.Stat(path)
		if err != nil || !info.Mode().IsDir() {
			err = ErrPathInvalid
			return 
		}
	}

	fullPath := filepath.Join(path, name)
	info, err = os.Stat(fullPath)
	if err == nil && !info.Mode().IsRegular() {
		err = ErrFileNameInvalid
		return 
	}

	if err != nil {
		if Debug {
			Logger.Printf("NewDBManager DB file %s does not exist, will be created", fullPath)
		}
	}

	dbm = & DBManager{
		name: name, 
		path: path,
		fullPath: fullPath,
		batchMode: batchMode, 
		buckets : buckets,
	}

	if err = dbm.SetSecret(secret); err != nil {
		return 
	}

	if err = dbm.openDB(); err != nil {
		return 
	}
	defer dbm.closeDB()

	return
}

// SetSecret is to set the AES Cryptor key, if the key is nil, the cryptor is not initialized; otherwise
// the cryptor is initialized, including the key and Cipher block that can be used directly for encrypt and decrypt functions
func (dbm * DBManager) SetSecret(secret string) (err error) {
	dbm.secret = secret
	if secret == "" {
		dbm.cryptor = nil 
	} else {
		dbm.cryptor, err = newAESCryptor([]byte(secret))
	}

	return err
}

// SetBatchMode is to set the batchMode for the boltdb. The boltdb file is always open in the file system unless the Close() is called.
// This cause inconvenience if you want to do some file operation to the db file while the program is running. Thus if the batchMode is
// set to false, the db will be closed after each db operation, this could reduce a certain performance. Thus if you have a lots of db 
// operations to execute, you can set the batchMode to be true before those operations. 
func (dbm * DBManager)SetBatchMode(mode bool) {
	dbm.batchMode = mode
	//if the batch mode is turned off, close DB directly 
	if !mode {
		dbm.closeDB()
	}
}

// This function creates the db file if it doesn't exist, and also initialize the buckets 
func (dbm * DBManager) openDB() (err error) {
	if dbm.batchMode && dbm.db != nil {
		return 
	}
	
	d, err := bolt.Open(dbm.fullPath, 0600, nil)
	if err != nil {
		return
	}

	db := &boltsecDB{d}

	initbuckets := func(tx * boltsecTx) error {
		for _, bname := range dbm.buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bname)); err != nil {
				return err
			}
		}
		return nil
	}

	if err = db.update(initbuckets); err != nil {
		db.Close()
		return
	}

	dbm.db = db
	return 
}

// The closeDB function closes the db when the dbm.db is not nil and the batchmode is false. 
// When the dbm batchmode is true, please set it to be false in order to close the DB.
func (dbm * DBManager) closeDB() {
	if !dbm.batchMode && dbm.db != nil  {
		dbm.db.Close()
		dbm.db = nil
	} 

	return 
}

// The view function is to retrieve the records
func (db *boltsecDB) view(fn func(*boltsecTx) error) error {
	wrapper := func(tx *bolt.Tx) error {
		return fn(&boltsecTx{tx})
	}
	return db.DB.View(wrapper)
}

// The update function applies changes to the database. There can be only one Update at a time.
func (db *boltsecDB) update(fn func(*boltsecTx) error) error {
	wrapper := func(tx *bolt.Tx) error {
		return fn(&boltsecTx{tx})
	}
	return db.DB.Update(wrapper)
}

// The GetByPrefix function returns the byte arrays for those records matched with specified Prefix. If the secret is set,
// the function returns the decrypted content.
func (dbm * DBManager) GetByPrefix(bucket, prefix string) ([][]byte, error) {
	var err error
	var results [][]byte
	if err = dbm.openDB(); err != nil {
		return nil, err
	}
	defer dbm.closeDB()

	results = make([][]byte, 0)

	seekPrefix := func(tx * boltsecTx) error {
		prefixKey := []byte(prefix)

		bkt := tx.Bucket([]byte(bucket))

		if bkt == nil {
			return bolt.ErrBucketNotFound
		}

		cursor := bkt.Cursor()
		for k, v := cursor.Seek(prefixKey); bytes.HasPrefix(k, prefixKey); k, v = cursor.Next() {

			if dbm.cryptor != nil {
				//secret key is set, decrypt the content before return
				content := make([]byte, len(v))
				copy(content, v)

				dec, err := dbm.cryptor.decrypt(content)
		        if err != nil {
		            return errors.New("Decrypt error from db")
		        }
			    results = append(results, dec)
			} else {
				results = append(results, v)
			}
		}
		return nil
	}

 	if err = dbm.db.view(seekPrefix); err != nil {
 		Logger.Printf("GetByPrefix return %s", err)
 	} 
 	
	return results, err
}

// The GetKeyList function returns the string array for keys with specified Prefix.  
func (dbm * DBManager) GetKeyList(bucket, prefix string) ([]string, error) {
	var err error
	var results []string
	if err = dbm.openDB(); err != nil {
		return nil, err
	}
	defer dbm.closeDB()

	results = make([]string, 0)

	seekPrefix := func(tx * boltsecTx) error {
		prefixKey := []byte(prefix)

		bkt := tx.Bucket([]byte(bucket))

		if bkt == nil {
			return bolt.ErrBucketNotFound
		}

		cursor := bkt.Cursor()
		for k, _ := cursor.Seek(prefixKey); k != nil && bytes.HasPrefix(k, prefixKey); k, _ = cursor.Next() {
		    results = append(results, string(k))
		}
		return nil
	}

 	if err = dbm.db.view(seekPrefix); err != nil {
 		Logger.Printf("GetByPrefix return %s", err)
 	} 
 	
	return results, err
}

// The GetOne function returns the first record containing the key, If the secret is set,
// the function returns the decrypted content.
func (dbm * DBManager)GetOne(bucket, key string) ([]byte, error) {
	var err error
	var result []byte

	if err = dbm.openDB(); err != nil {
		return nil, err
	}
	defer dbm.closeDB()

	if key == "" {
		return nil, ErrKeyInvalid
	}

	seek := func(tx * boltsecTx) error {
		prefixKey := []byte(key)

		bkt := tx.Bucket([]byte(bucket))

		if bkt == nil {
			return bolt.ErrBucketNotFound
		}

		cursor := bkt.Cursor()
		k, v := cursor.Seek(prefixKey)

		content := make([]byte, len(v))
		copy(content, v)

		if k != nil && bytes.HasPrefix(k, prefixKey) {
			if dbm.cryptor != nil {
				//secret key is set, decrypt the content before return 
				dec, err := dbm.cryptor.decrypt(content)
		        if err != nil {
		            return errors.New("Decrypt error from db")
		        }
		        result = dec
			} else {
				result = content
			}
		}

		return nil
	}

	if err := dbm.db.view(seek); err != nil {
		return nil, err
	}

	return result, nil
}

// The Save function stores the record into the db file. If the secret value is set, the function 
// encrypts the content before storing into the db.
func (dbm * DBManager)Save(bucket, key string, data interface{}) error {
	var err error

	if err = dbm.openDB(); err != nil {
		return err
	}
	defer dbm.closeDB()

	if data == nil {
		return errors.New("data is nil")
	}

	save := func(tx * boltsecTx) error {
		var err error
		bkt := tx.Bucket([]byte(bucket))

		value, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if dbm.cryptor != nil {
			//encrypt the content before store in the db 
			enc, err := dbm.cryptor.encrypt(value)
	        if err != nil {
	            return errors.New("Decrypt error from db")
	        }

			if err = bkt.Put([]byte(key), enc); err != nil {
				return err
			}			
		} else {
			if err = bkt.Put([]byte(key), value); err != nil {
				return err
			}
		} 

		return nil
	}

	return dbm.db.update(save)
}

// The Delete function deletes the record specified by the key.
func (dbm * DBManager)Delete(bucket, key string) error {
	var err error

	if err = dbm.openDB(); err != nil {
		return err
	}
	defer dbm.closeDB()

	if key == "" {
		return errors.New("cannot delete, key is nil")
	}

	delete := func(tx * boltsecTx) error {	
		bkt := tx.Bucket([]byte(bucket))
		if err := bkt.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	}

	return dbm.db.update(delete)
}
