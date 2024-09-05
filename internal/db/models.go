// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

type Config struct {
	ProfileID   int64       `json:"profile_id"`
	Entrypoints interface{} `json:"entrypoints"`
	Routers     interface{} `json:"routers"`
	Services    interface{} `json:"services"`
	Middlewares interface{} `json:"middlewares"`
	Version     *string     `json:"version"`
}

type Profile struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Tls      bool    `json:"tls"`
}

type Provider struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	ExternalIp string  `json:"external_ip"`
	ApiKey     string  `json:"api_key"`
	ApiUrl     *string `json:"api_url"`
	IsActive   bool    `json:"is_active"`
}

type User struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
	Type     string  `json:"type"`
}