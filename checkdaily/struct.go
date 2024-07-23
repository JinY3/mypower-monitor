package checkdaily

type User struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Homeid   string `json:"homeid"`
}
