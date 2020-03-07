package main

type User struct {
	ID             int
	Name           string
	Place          int
	Index          int `json:"index"`
	CurrentError   int
	LastErrorIndex int
	CurrentSpeed   float32
}
