package entities

import (
	vo "chatapp/internal/app/valobjects"
)

// `nil` value of struct
var NilStruct = User{}

type User struct {
	Id       string   `json:"id"`
	Name     vo.Name  `json:"name"`
	Username vo.Email `json:"username"`
	Password string   `json:"password"`
}

func (this User) Email() vo.Email {
	return this.Username
}

// func (this User) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(map[string]interface{}{a.Name: a.Value})
// }

// func (this User) getName() vo.Name {
// 	name, _ := vo.NewName(this.Name)
// 	return name
// }

// func setName(name string) {
// 	n, _ := vo.NewName(name)
// 	Name = n
// }
