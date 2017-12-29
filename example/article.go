/*MIT License

Copyright (c) 2017 linkthings

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.*/

// A sample code to describe how to use the boltsec and use json.Unmarshal to 
// convert the []byte values into Article object 
package example

import (
	"time"
	"encoding/json"
	"sort"
	"fmt"
	"errors"
	"math/rand"	
	"log"
	"os"
	"boltsec"
)

var Logger = log.New(os.Stdout, "[ArticleManager] ", log.LstdFlags)

type Article struct {
	ID     		string	`json:"id"`
	Name 		string	`json:"title"`
	Tags 	 [] string 	`json:"tags"`
	Content     string `json:"content"`
	ContentToMark string `json:"contentToMark"`
	CreatedAt 	time.Time  `json:"createdAt"`
	UpdatedAt	time.Time  `json:"updatedAt"`
}

type ArticleSearch struct {
	Keywords 	 	string 	`json:"keywords"`
}

type ArticleSortByUpdateTime []*Article

func (a ArticleSortByUpdateTime) Len() int           { return len(a) }
func (a ArticleSortByUpdateTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArticleSortByUpdateTime) Less(i, j int) bool { return a[i].UpdatedAt.After(a[j].UpdatedAt) }

type ArticleManager struct {
	dbm * boltsec.DBManager
	bucket string
	articlePrefix string
}

func NewArticleManager(name, fullpath, secret string)(* ArticleManager, error){
	dbm, err := boltsec.NewDBManager(name, fullpath, secret, false, []string{"al-article"})
	var am * ArticleManager
	if dbm != nil {
		am = &ArticleManager{
			dbm: dbm, 
			bucket: "al-article", 
			articlePrefix : "a-",
		}
	}

	InitRand()
	return am, err
}

func (am * ArticleManager) NewArticle (name, content, contentToMark string, tags [] string) (* Article, error) {
	article := & Article{
		Name: name,
		Content: content,
		ContentToMark: contentToMark,
		Tags: tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := am.Save(article)
	return article, err 
}


func (am * ArticleManager) Seek() (results []*Article, err error) {
	_func := "Seek"
	if am == nil || am.dbm == nil {
		return 
	}
	results = make([]*Article, 0)

	var byteResults [][]byte

	if byteResults, err = am.dbm.GetByPrefix(am.bucket,am.articlePrefix); err != nil {
		Logger.Printf("%s fsm.dbm.GetByPrefix return err: %s", _func, err)
		return nil, err
	}

	for _, iter := range byteResults {
		resNew := new(Article)
		if err = json.Unmarshal([]byte(iter), resNew); err != nil {
	        Logger.Printf("%s json.Unmarshal return err: %s", _func, err)
	    } else {
	    	results = append(results, resNew)
	    }
	}

	sort.Sort(ArticleSortByUpdateTime(results))

	return
}

func (am * ArticleManager)GetByID(id string) (result *Article, err error) {
	var _func = "GetByID"
	if id == "" {
		return nil, errors.New("id is empty")
	}
	
	var key string
	key = fmt.Sprintf("%s%s", am.articlePrefix, id)
	Logger.Printf("%s get using key %s", _func, key)

	var byt []byte
	byt, err = am.dbm.GetOne(am.bucket, key)

	result = new(Article)
	if err = json.Unmarshal([]byte(byt), result); err != nil {
        Logger.Printf("%s json.Unmarshal[%s] return err: %s", _func, id, err)
    }

	return
}

func (am * ArticleManager) Save(record * Article) error {
	if record == nil {
		return errors.New("record is nil")
	}
	if record.ID == "" || record.ID == "0" {
		record.ID = RandStringRunes(12)
	}

	record.UpdatedAt = time.Now()

	key := fmt.Sprintf("%s%s", am.articlePrefix, record.ID)
	return am.dbm.Save(am.bucket, key, record)
}

func (am * ArticleManager) Update(record * Article) error {
	if record == nil || record.ID == "0" {
		return errors.New("record is nil")
	}

	var oldRecord * Article
	var err error

	if oldRecord, err = am.GetByID(record.ID); err != nil {
		return err
	}
	//update the time values from the old record
	record.CreatedAt = oldRecord.CreatedAt
	record.UpdatedAt = time.Now()

	return am.Save(record)
}

func (am * ArticleManager) Delete(id string) error {
	if id == "" {
		return boltsec.ErrKeyInvalid
	}
	key := fmt.Sprintf("%s%s", am.articlePrefix, id)
	return am.dbm.Delete(am.bucket, key)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func InitRand() {
    rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
    b := make([]rune, n)

    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}
