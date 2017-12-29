# boltsec
A boltdb wrapper to encrypt and decrypt the values stored in the boltdb via AES Cryptor, and also provides common 
db operations such as GetOne, GetByPrefix, GetKeyList, Save, Delete and etc. 

Boltdb file is always open in the file system unless the DB.Close() is called, which cause inconvenience 
if you want to do some file operations to the db file while the program is running. This package provides the parameter: batchMode to 
control whether to close the db after each db operation, this has performance impact but could be a useful feature.

## Features
1. [x] Support AES to encrypt and decrypt the values
1. [x] Batch mode option to control whether to close the db after each db operation 
1. [x] Initialize db file and cryptor

## Performance
The below is the benchmark data for the DB related operations.
```golang
BenchmarkDBMOps-8               	    2000	   1114202 ns/op
BenchmarkDBMOpsBatchMode-8      	    5000	    383995 ns/op
BenchmarkDBMOpsNoEncryption-8   	    1000	   1029510 ns/op
```

## LICENSE
This project is under license [MIT](LICENSE)

### Usage Example

```golang
var err error
bucketName := "article"

dbm, err := NewDBManager("test.dat", "example", "secret", false, []string{bucketName})

data := Article{
	ID: "ID-0001", 
	Title: "input with more than 16 characters", 
}

if err=dbm.Save(bucketName, data.ID, data); err != nil {
	t.Errorf("TestDBMCreate save data return err: %s", err)
}

var bytes []byte
if bytes, err = dbm.GetOne(bucketName, data.ID); err != nil {
	t.Errorf("TestDBMCreate GetOne return err: %s", err)
}

resNew := new(Article)
if err = json.Unmarshal(bytes, resNew); err != nil {
    t.Errorf("json.Unmarshal return err: %s", err)
}

if resNew.Title != data.Title {
	t.Errorf("returned Title is not equal: new:%s, org: %s", resNew.Title, data.Title)
}

if err=dbm.Delete(bucketName, data.ID); err != nil {
	t.Errorf("TestDBMCreate Delete return err: %s", err)
}
```

A example code can be found in [example](example/) folder on how to utilize this package
