package models

import (
	"errors"
	"fmt"
)

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
	Type     string `json:"type" binding:"required"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (lr *RequestLoginRequest) Examine() error {
	fmt.Println("type: ", lr.Type)
	fmt.Println("email: ", lr.Email)
	fmt.Println("username: ", lr.Username)

	if lr.Type == "" {
		return errors.New("type not found")
	} else if lr.Type == "username" && lr.Username == "" {
		return errors.New("type set to {" + lr.Type + "}, but value not provided")
	} else if lr.Type == "email" && lr.Email == "" {
		return errors.New("type set to {" + lr.Type + "}, but value not provided")
	}
	return nil
}

type AvatarUpdateRequest struct {
	AvatarData string `json:"avatarData"`
	AvatarExt  string `json:"avatarExt"`
}

func (ar *AvatarUpdateRequest) Examine() error {
	if ar.AvatarData == "" {
		return errors.New("avatarData not found")
	} else if ar.AvatarExt == "" {
		return errors.New("avatarExt not found")
	}
	return nil
}
