package types

type Post struct {
	ID         int64    `json:"id"`
	Created_by string   `json:"created_by"`
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	Likes      int      `json:"likes"`
	LikedBy    []string `json:"liked_by,omitempty"`
	Created_at string   `json:"created_at"`
	Updated_at string   `json:"updated_at"`
}

type UsersPosts struct {
	Posts []Post `json:"posts,omitempty"`
}
