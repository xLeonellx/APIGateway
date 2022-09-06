package postgres

import (
	"GoNews/pkg/storage"
	"testing"
)

const dbURL = "postgres://postgres:postgrespw@localhost:55000/news"

func Test_postgres(t *testing.T) {

	db, err := New(dbURL)
	if err != nil {
		t.Fatal(err)
	}

	testCase := []storage.Post{{ID: 1,
		Title:   "Test Title",
		Content: "Test Content",
		PubTime: 0,
		Link:    "Test Link"}}

	err = db.PostsMany(testCase)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.Posts(2, 0)
	if err != nil {
		t.Fatal(err)
	}

	const wantLen = 1
	if len(news) < wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(news), wantLen)
	}

	_, _, err = db.PostSearchILIKE("", 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.PostDetal(2)
	if err != nil {
		t.Fatal(err)
	}

	err = db.UpdatePost(testCase[0])
	if err != nil {
		t.Fatal(err)
	}
	err = db.DeletePost(testCase[0])
	if err != nil {
		t.Fatal(err)
	}

}
