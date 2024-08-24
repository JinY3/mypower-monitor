package checkdaily

type User struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	To       string `json:"to"`
	Homeid   string `json:"homeid"`
}
