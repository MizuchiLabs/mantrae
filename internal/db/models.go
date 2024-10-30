// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type Agent struct {
	ID         string      `json:"id"`
	Hostname   string      `json:"hostname"`
	PublicIp   *string     `json:"publicIp"`
	PrivateIps interface{} `json:"privateIps"`
	Containers interface{} `json:"containers"`
	LastSeen   *time.Time  `json:"lastSeen"`
}

type Config struct {
	ProfileID   int64       `json:"profileId"`
	Overview    interface{} `json:"overview"`
	Entrypoints interface{} `json:"entrypoints"`
	Routers     interface{} `json:"routers"`
	Services    interface{} `json:"services"`
	Middlewares interface{} `json:"middlewares"`
	Tls         interface{} `json:"tls"`
	Version     *string     `json:"version"`
}

type Middleware struct {
	ID        string      `json:"id"`
	ProfileID int64       `json:"profileId"`
	Name      string      `json:"name"`
	Provider  string      `json:"provider"`
	Type      string      `json:"type"`
	Protocol  string      `json:"protocol"`
	AgentID   *string     `json:"agentId"`
	Content   interface{} `json:"content"`
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
	ExternalIp string  `json:"externalIp"`
	ApiKey     string  `json:"apiKey"`
	ApiUrl     *string `json:"apiUrl"`
	ZoneType   *string `json:"zoneType"`
	Proxied    bool    `json:"proxied"`
	IsActive   bool    `json:"isActive"`
}

type Router struct {
	ID          string      `json:"id"`
	ProfileID   int64       `json:"profileId"`
	Name        string      `json:"name"`
	Provider    string      `json:"provider"`
	Protocol    string      `json:"protocol"`
	Status      *string     `json:"status"`
	AgentID     *string     `json:"agentId"`
	EntryPoints interface{} `json:"entryPoints"`
	Middlewares interface{} `json:"middlewares"`
	Rule        string      `json:"rule"`
	RuleSyntax  *string     `json:"ruleSyntax"`
	Service     string      `json:"service"`
	Priority    *int64      `json:"priority"`
	Tls         interface{} `json:"tls"`
	DnsProvider *int64      `json:"dnsProvider"`
	Errors      interface{} `json:"errors"`
}

type Service struct {
	ID           string      `json:"id"`
	ProfileID    int64       `json:"profileId"`
	Name         string      `json:"name"`
	Provider     string      `json:"provider"`
	Type         string      `json:"type"`
	Protocol     string      `json:"protocol"`
	AgentID      *string     `json:"agentId"`
	Status       *string     `json:"status"`
	ServerStatus interface{} `json:"serverStatus"`
	LoadBalancer interface{} `json:"loadBalancer"`
	Weighted     interface{} `json:"weighted"`
	Mirroring    interface{} `json:"mirroring"`
	Failover     interface{} `json:"failover"`
}

type Setting struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type User struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
	Type     string  `json:"type"`
}
