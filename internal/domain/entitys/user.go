package entitys

const (
	PRIVATE_INFORMATION_MESSAGE = "private_information"
)

type (
	Location struct {
		Lat     float64 `json:"lat"`
		Ing     float64 `json:"ing"`
		Address string  `json:"address"`
		City    string  `json:"city"`
		State   string  `json:"state"`
		ZipCode int64   `json:"zip_code"`
	}
	User struct {
		ID       uint64   `json:"id"`
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Login    string   `json:"login"`
		Password string   `json:"password"`
		Location Location `json:"location"`
	}
)

func (u *User) ClearSensitiveInformation() {
	u.ID = 0
	u.Email = "private_information"
	u.Login = "private_information"
	u.Password = "private_information"
}
