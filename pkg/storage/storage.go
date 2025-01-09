package storage

// Post - публикация.
type Post struct {
	ID          int    `bson:"_id,omitempty"`
	Title       string `bson:"title"`
	Content     string `bson:"content"`
	AuthorID    int    `bson:"authorId"`
	AuthorName  string `bson:"authorName"`
	CreatedAt   int64  `bson:"createdAt"`
	PublishedAt int64  `bson:"publishedAt"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
