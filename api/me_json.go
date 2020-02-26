package main

type IntraUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
	Staff    bool   `json:"staff?"`
}
