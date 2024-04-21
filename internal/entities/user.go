package entities

type User struct {
	ID          uint64   `json:"id"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Expressions []uint64 `json:"expressions"`
}

type UserData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
