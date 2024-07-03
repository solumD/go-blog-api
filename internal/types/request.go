package types

type UserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PostReq struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
