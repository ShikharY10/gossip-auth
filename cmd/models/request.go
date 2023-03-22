package models

import "errors"

type SignupRequest struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	AvatarData string `json:"imageData"`
	AvatarExt  string `json:"avatarExt"`
}

func (sr *SignupRequest) Examine() error {
	if sr.AvatarData == "" {
		return errors.New("avatarData not found")
	} else if sr.AvatarExt == "" {
		return errors.New("avatarExt not found")
	} else if sr.Name == "" {
		return errors.New("name not found")
	} else if sr.Username == "" {
		return errors.New("username not found")
	}
	return nil
}
