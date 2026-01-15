package commands

import (
	"fmt"

	"github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/*
HashPasswordCommand creates a new user in the system
*/
func HashPasswordCommand(c *cli.Context, tx *gorm.DB) error {
	plainPassword := c.String("password")
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost, // currently 10
	)
	if err != nil {
		return err
	}
	fmt.Println("New hashed password: ", string(hashedBytes))
	return nil
}
