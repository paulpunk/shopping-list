package main

type List struct {
	ID         int
	Version    int
	User       string
	Name       string
	State      string
	SharedWith []*SharedUser
}
