package entity

type User struct {
	ID      string
	Account Account
	Profile Profile
}

type Account struct {
	Email    string
	Password string
}

type Profile struct {
	Name   string
	Avatar string
	Bio    string
}

type Post struct {
	ID           string
	Body         string
	Timestamp    int
	RepliesCount int
	LikesCount   int
	IsLiked      bool
	Author       User
	User         User
}
