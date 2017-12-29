package example

import (
	"testing"
)

func TestArticleManager(t * testing.T) {
	var err error

	am, err := NewArticleManager("am.dat", "", "amsecret")
	if err != nil {
		t.Errorf("NewArticleManager return err: %s", err)
	}

	article, err := am.NewArticle("test-name", "test-content", "test-content", []string{"tag1", "tag2"})
	if err != nil {
		t.Errorf("NewArticle return err: %s", err)
	}

	if article.ID == "" {
		t.Errorf("NewArticle Id is nil")
	}

	retRecord, err := am.GetByID(article.ID)
	if err != nil {
		t.Errorf("GetByID return err: %s", err)
	}

	if article.Name != retRecord.Name ||  article.Content != retRecord.Content {
		t.Errorf("GetByID returned values are not equal")
	}

	return
}