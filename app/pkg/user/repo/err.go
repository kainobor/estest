package repo

import "fmt"

type ErrNotFound struct {
	login string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("User with login %s not found", err.login)
}
