package main

import "github.com/ToshihiroOgino/elib/repository"

func main() {
	userRepository := repository.NewUserRepository()
	allUsers, _ := userRepository.FindAll()
	for _, user := range allUsers {
		println(user.Email) // Assuming User has an Email field
	}
}
