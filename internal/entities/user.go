package entities

type User struct {
	ID          uint64   `json:"id"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Expressions []uint64 `json:"expressions"`
}
