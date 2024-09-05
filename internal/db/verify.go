package db

import (
	"fmt"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
)

func (p *CreateProfileParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Url == "" {
		if !util.IsValidURL(p.Url) {
			return fmt.Errorf("url is not valid")
		}
	}
	return nil
}

func (p *UpdateProfileParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Url == "" {
		if !util.IsValidURL(p.Url) {
			return fmt.Errorf("url is not valid")
		}
	}
	return nil
}

func (u *CreateUserParams) Verify() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if u.Type == "" {
		return fmt.Errorf("user type cannot be empty")
	}

	hash, err := util.HashPassword(strings.TrimSpace(u.Password))
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}
	u.Password = hash
	return nil
}

func (u *UpdateUserParams) Verify() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if u.Type == "" {
		return fmt.Errorf("user type cannot be empty")
	}

	hash, err := util.HashPassword(strings.TrimSpace(u.Password))
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}
	u.Password = hash
	return nil
}

func (p *CreateProviderParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Type == "" {
		return fmt.Errorf("provider type cannot be empty")
	}
	if p.ExternalIp == "" {
		return fmt.Errorf("external ip cannot be empty")
	}
	if p.ApiKey == "" {
		return fmt.Errorf("api key cannot be empty")
	}
	return nil
}

func (p *UpdateProviderParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Type == "" {
		return fmt.Errorf("provider type cannot be empty")
	}
	if p.ExternalIp == "" {
		return fmt.Errorf("external ip cannot be empty")
	}
	if p.ApiKey == "" {
		return fmt.Errorf("api key cannot be empty")
	}
	return nil
}
