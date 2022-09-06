package storage

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    `json:"ID,omitempty"`      // номер записи
	Title   string `json:"title,omitempty"`   // заголовок публикации
	Content string `json:"content,omitempty"` // содержание публикации
	PubTime int64  `json:"pubTime,omitempty"` // время публикации
	Link    string `json:"link,omitempty"`    // ссылка на источник
}

type Pagination struct {
	NumOfPages int `json:"numOfPages,omitempty"`
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	PostSearchILIKE(keyWord string, limit, offset int) ([]Post, Pagination, error)
	PostDetal(id int) (Post, error)
	Posts(limit, offset int) ([]Post, error) // получение n-ого кол-ва публикаций
	AddPost(Post) error                      // создание новой публикации
	PostsMany([]Post) error                  // создание n-ого кол-ва публикаций
	UpdatePost(Post) error                   // обновление публикации
	DeletePost(Post) error                   // удаление публикации по ID
}
