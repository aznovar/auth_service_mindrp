package models

type User struct {
	ID       int64
	SocialId string
	Email    string
	PassHash []byte
}
