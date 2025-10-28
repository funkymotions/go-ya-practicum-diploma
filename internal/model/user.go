package model

type User struct {
	ID    uint
	Login string
	// password will be hashed using md5 alogrithm
	// w/o using salt to keep it simple in an educational way
	PasswordHash string
}
