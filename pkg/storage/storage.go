package storage

// publication model
type Post struct {
	ID      int    `bson:"id" json:"ID"`            // post id
	Title   string `bson:"title" json:"Title"`      // post title
	Content string `bson:"content" json:"Content"`  // post content
	PubTime int64  `bson:"pub_time" json:"PubTime"` // publication time
	Link    string `bson:"link" json:"Link"`        // source link
}

// interface to work with db
type DBinterface interface {
	LastNews(n int) ([]Post, error) //gets n recent posts
	StoreNews([]Post) error         // creates posts array
	Close()                         // closes connection with db
}
