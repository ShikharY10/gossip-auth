package models

import "errors"

type SignupRequest struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	AvatarData string `json:"avatarData"`
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

type RequestLoginRequest struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (lr *RequestLoginRequest) Examine() error {
	if lr.Type == "" {
		return errors.New("type not found")
	} else if lr.Email == "" && lr.Username == "" {
		return errors.New(lr.Type + "not found")
	}
	return nil
}
