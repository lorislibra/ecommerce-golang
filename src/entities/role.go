package entities

import "fmt"

type Role int

const (
	_ Role = iota
	Admin
	Seller
	Teacher
	Student
)

func (r Role) In(list []Role) bool {
	for _, role := range list {
		if role == r {
			return true
		}
	}
	return false
}

func (r Role) String() string {
	switch r {
	case Admin:
		return "Admin"
	case Seller:
		return "Seller"
	case Teacher:
		return "Teacher"
	case Student:
		return "Student"
	default:
		return fmt.Sprintf("%d", r)
	}
}
