package shared

//User Recime User
type User struct {
	Email    string   `json:"email"`
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Company  string   `json:"company"`
	Config   []Config `json:"config"`
}
