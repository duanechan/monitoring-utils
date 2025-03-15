// Copyright © 2025 Duane Matthew P. Chan

package email

import "fmt"

type User struct {
	Name  string
	Email string
}

func (u User) String() string {
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

func (u User) Exists() bool {
	return u.Name != "" && u.Email != ""
}
