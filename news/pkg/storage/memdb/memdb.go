package memdb

import "GoNews/pkg/storage"

// Хранилище данных в памяти, используемая для тестов.
type Store struct{}

// Конструктор объекта хранилища.
func New() *Store {
	return new(Store)
}

func (s *Store) PostSearchILIKE(keyWord string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
	return nil, storage.Pagination{}, nil
}

func (s *Store) PostDetal(id int) (storage.Post, error) {
	return storage.Post{}, nil
}

func (s *Store) Posts(limit, offset int) ([]storage.Post, error) {
	return posts[0:0], nil
}

func (s *Store) AddPost(storage.Post) error {
	return nil
}
func (s *Store) PostsMany([]storage.Post) error {
	return nil
}
func (s *Store) UpdatePost(storage.Post) error {
	return nil
}
func (s *Store) DeletePost(storage.Post) error {
	return nil
}

var posts = []storage.Post{
	{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
		PubTime: 0,
		Link:    "Test Link"},
	{
		ID:      2,
		Title:   "Test Title2",
		Content: "Test Content2",
		PubTime: 0,
		Link:    "Test Link2"},
	{
		ID:      3,
		Title:   "Test Title3",
		Content: "Test Content3",
		PubTime: 0,
		Link:    "Test Link3"},
	{
		ID:      4,
		Title:   "Test Title4",
		Content: "Test Content4",
		PubTime: 0,
		Link:    "Test Link4"},
}
