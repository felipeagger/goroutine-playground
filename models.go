package main

import "time"

type chanResp struct {
	User UserInfo
	Err error
}

type UserInfo struct {
	Login       string      `json:"login"`
	ID          int         `json:"id"`
	AvatarURL   string      `json:"avatar_url"`
	URL         string      `json:"url"`
	ReposURL    string      `json:"repos_url"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Company     interface{} `json:"company"`
	Blog        string      `json:"blog"`
	PublicRepos int         `json:"public_repos"`
	Followers   int         `json:"followers"`
	Following   int         `json:"following"`
	CreatedAt   time.Time   `json:"created_at"`
}