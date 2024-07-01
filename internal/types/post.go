package types

type UsersPosts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID         int64  `json:"id"`
	Created_by string `json:"created_by"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}
