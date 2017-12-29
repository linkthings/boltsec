package boltsec

import (
	"encoding/json"
	"testing"
)

type Article struct {
	ID     		string	`json:"id"`
	Title 		string	`json:"title"`
}

func ExampleNewDBManager() {
	var err error
	bucketName := "article"

	dbm, err := NewDBManager("test.dat", "example", "secret", false, []string{bucketName})

	if err != nil {
		// Handle the error
	}

	var bytes []byte
	if bytes, err = dbm.GetOne(bucketName, "data.ID"); err != nil {
		// handle the error
	}

	resNew := new(Article)
	if err = json.Unmarshal(bytes, resNew); err != nil {
        // handle the error
    }
}

func TestDBMCreate(t * testing.T) {
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

	return
}

func BenchmarkDBMOps(b *testing.B) {
	var err error
	bucketName := "article"

	dbm, err := NewDBManager("test.dat", "example", "secret", false, []string{bucketName})

	data := Article{
		ID: "ID-0001", 
		Title: "input with more than 16 characters", 
	}

    for n := 0; n < b.N; n++ {
		if err=dbm.Save(bucketName, data.ID, data); err != nil {
			b.Errorf("TestDBMCreate save data return err: %s", err)
		}

		var bytes []byte
		if bytes, err = dbm.GetOne(bucketName, data.ID); err != nil {
			b.Errorf("TestDBMCreate GetOne return err: %s", err)
		}

		resNew := new(Article)
		if err = json.Unmarshal(bytes, resNew); err != nil {
	        b.Errorf("json.Unmarshal return err: %s", err)
	    }

	    if resNew.Title != data.Title {
	    	b.Errorf("returned Title is not equal: new:%s, org: %s", resNew.Title, data.Title)
	    }

	    if err=dbm.Delete(bucketName, data.ID); err != nil {
	    	b.Errorf("TestDBMCreate Delete return err: %s", err)
	    }
    }
}

func BenchmarkDBMOpsBatchMode(b *testing.B) {
	var err error
	bucketName := "article"

	dbm, err := NewDBManager("test.dat", "example", "", true, []string{bucketName})

	data := Article{
		ID: "ID-0001", 
		Title: "input with more than 16 characters", 
	}

    for n := 0; n < b.N; n++ {
		if err=dbm.Save(bucketName, data.ID, data); err != nil {
			b.Errorf("TestDBMCreate save data return err: %s", err)
		}

		var bytes []byte
		if bytes, err = dbm.GetOne(bucketName, data.ID); err != nil {
			b.Errorf("TestDBMCreate GetOne return err: %s", err)
		}

		resNew := new(Article)
		if err = json.Unmarshal(bytes, resNew); err != nil {
	        b.Errorf("json.Unmarshal return err: %s", err)
	    }

	    if resNew.Title != data.Title {
	    	b.Errorf("returned Title is not equal: new:%s, org: %s", resNew.Title, data.Title)
	    }

	    if err=dbm.Delete(bucketName, data.ID); err != nil {
	    	b.Errorf("TestDBMCreate Delete return err: %s", err)
	    }
    }

    dbm.SetBatchMode(false)
}


func BenchmarkDBMOpsNoEncryption(b *testing.B) {
	var err error
	bucketName := "article"

	dbm, err := NewDBManager("test.dat", "example", "", false, []string{bucketName})

	data := Article{
		ID: "ID-0001", 
		Title: "input with more than 16 characters", 
	}

    for n := 0; n < b.N; n++ {
		if err=dbm.Save(bucketName, data.ID, data); err != nil {
			b.Errorf("TestDBMCreate save data return err: %s", err)
		}

		var bytes []byte
		if bytes, err = dbm.GetOne(bucketName, data.ID); err != nil {
			b.Errorf("TestDBMCreate GetOne return err: %s", err)
		}

		resNew := new(Article)
		if err = json.Unmarshal(bytes, resNew); err != nil {
	        b.Errorf("json.Unmarshal return err: %s", err)
	    }

	    if resNew.Title != data.Title {
	    	b.Errorf("returned Title is not equal: new:%s, org: %s", resNew.Title, data.Title)
	    }

	    if err=dbm.Delete(bucketName, data.ID); err != nil {
	    	b.Errorf("TestDBMCreate Delete return err: %s", err)
	    }
    }
}

