package boltsec

import (
	"encoding/json"
	"testing"
	proto "github.com/golang/protobuf/proto"
)

type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
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

func TestJsonCreate(t *testing.T) {
	// marshal := func (input interface{}) (out interface{}) {
	// 	value, err := json.Marshal(input)
	// 	if err != nil {
	// 		t.Errorf()
	// 	}
	// 	return value
	// }

	bucketName := "jsonArticle"
	dbm, err := NewDBManager("json.dat", ".", "secret", false, []string{bucketName})

	data := Article{
		ID:    "JSON-0001",
		Title: "How much JSON could a jchuck chuck if a jchuck could chuck json?",
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Json Serialization Failed: %s", err)
	}

	if err = dbm.Save(bucketName, data.ID, rawData); err != nil {
		t.Errorf("TestJsonCreate save data returned: %s", err)
	}

}

func TestProtobufCreate(t *testing.T) {
	bucketName := "protobufArticle"
	dbm, err := NewDBManager("protobuf.data", ".", "secret", false, []string{bucketName})

	data := &ProtobufArticle{
		ID:		"PROTOBUF-0001",
		Title:  "We have a protobuf here...",
	}

	rawData, err := proto.Marshal(data)
	if err = dbm.Save(bucketName, data.ID, rawData); err != nil {
		t.Errorf("TestProtobufCreate save data returned: %s", err)
	}

	var bytes []byte
	if bytes, err = dbm.GetOne(bucketName, data.ID); err != nil {
		t.Errorf("dbm.GetOne error: %s", err)
	}

	data2 := &ProtobufArticle{}
	if err = proto.Unmarshal(bytes, data2); err != nil {
		t.Errorf("protobuf unmarshal error: %s", err)
	}

	if ! proto.Equal(data, data2) {
		t.Errorf("Marshal/Save/Get/Unmarshal != source data")
	}
}

func BenchmarkDBMOps(b *testing.B) {
	var err error
	bucketName := "article"

	dbm, err := NewDBManager("test.dat", "example", "secret", false, []string{bucketName})

	data := Article{
		ID:    "ID-0001",
		Title: "input with more than 16 characters",
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		b.Errorf("Benchmark Error - Marshaling Json...: %s", err)
	}

	for n := 0; n < b.N; n++ {
		if err = dbm.Save(bucketName, data.ID, rawData); err != nil {
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

		if err = dbm.Delete(bucketName, data.ID); err != nil {
			b.Errorf("TestDBMCreate Delete return err: %s", err)
		}
	}
}

func BenchmarkDBMOpsBatchMode(b *testing.B) {
	var err error
	bucketName := "article"

	dbm, err := NewDBManager("test.dat", "example", "", true, []string{bucketName})

	data := Article{
		ID:    "ID-0001",
		Title: "input with more than 16 characters",
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		b.Errorf("Benchmark Error - Marshaling Json...: %s", err)
	}

	for n := 0; n < b.N; n++ {
		if err = dbm.Save(bucketName, data.ID, rawData); err != nil {
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

		if err = dbm.Delete(bucketName, data.ID); err != nil {
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
		ID:    "ID-0001",
		Title: "input with more than 16 characters",
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		b.Errorf("Benchmark Error - Marshaling Json...: %s", err)
	}

	for n := 0; n < b.N; n++ {
		if err = dbm.Save(bucketName, data.ID, rawData); err != nil {
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

		if err = dbm.Delete(bucketName, data.ID); err != nil {
			b.Errorf("TestDBMCreate Delete return err: %s", err)
		}
	}
}
