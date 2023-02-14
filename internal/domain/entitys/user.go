package entitys

type User struct {
	ID       uint64 `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
