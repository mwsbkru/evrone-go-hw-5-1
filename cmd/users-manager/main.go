package main

import (
	"fmt"
	"hw_2_1/internal/entity"
	"hw_2_1/internal/repo"
	"hw_2_1/internal/usecase"
)

func main() {
	repos := []repo.UserRepository{repo.NewMockUserRepo(), repo.NewInMemoryUserRepo()}

	for _, repository := range repos {
		service := usecase.NewUserService(repository)

		fmt.Printf("=====Repository %T=====\n", repository)

		fmt.Println("\n-----ListUsers()-----")
		fmt.Println(service.ListUsers())

		fmt.Println("\n-----CreateUser()-----")
		fmt.Println(service.CreateUser("Vasia", "vasia@email.com", entity.UserRoleUser))
		fmt.Println(service.CreateUser("Petia", "petia@email.com", entity.UserRoleAdmin))
		fmt.Println(service.CreateUser("Vova", "vova@email.com", entity.UserRoleGuest))

		fmt.Println("\n-----GetUser()-----")
		fmt.Println(service.GetUser("vasia@email.com"))

		fmt.Println("\n-----RemoveUser()-----")
		fmt.Println(service.RemoveUser("vasia@email.com"))

		fmt.Println("\n-----FindByRole()-----")
		fmt.Println(service.FindByRole(entity.UserRoleGuest))

		fmt.Println("\n-----ListUsers()-----")
		fmt.Println(service.ListUsers())

		fmt.Println("\n\n\n")
	}
}
