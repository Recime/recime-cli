package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//User Recime User
type User struct {
	Email    string   `json:"email"`
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Company  string   `json:"company"`
	Config   []Config `json:"config"`
	APIKey   string   `json:"apiKey"`
}

// CurrentUser gets the current user for the token.
func (u *User) CurrentUser(source string, token string) User {
	client := &http.Client{}

	url := fmt.Sprintf("%s/user", source)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	resp, err := client.Do(req)

	check(err)

	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)

	var result User

	json.Unmarshal(dat, &result)

	return result

}
