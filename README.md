# boltsec
A boltdb wrapper to encrypt and decrypt the values stored in the boltdb via AES Cryptor, and also provides common 
db operations such as `GetOne`, `GetByPrefix`, `GetKeyList`, `Save`, `Delete`, etc... 

Boltdb file is always open in the file system unless  `DB.Close()` is called, which cause inconvenience 
if you want to do some file operations to the db file while the program is running. This package provides the parameter: `batchMode` to 
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

## Usage Example

<details>
	<summary>header code</summary>

```golang
package main

import (
	"encoding/json"
	"fmt"

	"github.com/linkthings/boltsec"
)

type article struct {
	ID            string
	Name          string
}
```

</details>

```golang	bucketName := "article"
func main() {
	bucketName := "article"
	dbm, _ := boltsec.NewDBManager("test.dat", "", "secret", false, []string{bucketName})

	data := article{
		ID:   "ID-0001",
		Name: "input with more than 16 characters",
	}

	dbm.Save(bucketName, data.ID, data)

	var bytes []byte
	bytes, _ = dbm.GetOne(bucketName, data.ID)

	resNew := new(article)
	json.Unmarshal(bytes, resNew)

	if resNew.Name == data.Name {
		fmt.Printf("Decrypted data matches origional input:\nDecrypted:%s\nOriginal:%s\n", resNew.Name, data.Name)
	}

	dbm.Delete(bucketName, data.ID)
}
```

A complete example with error handling and testing can be found in the [example](example/) folder
