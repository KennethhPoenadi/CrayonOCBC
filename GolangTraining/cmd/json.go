package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address,omitempty"`
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use: "json",
		Run: func(cmd *cobra.Command, args []string) {
			decode()
		},
	})
}

func encode() {
	u := User{
		Name:    "John",
		Age:     17,
		Address: "BSD",
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}

func decode() {
	str := `{"name":"John","age":17,"address":"BSD"}`
	var u User
	err := json.Unmarshal([]byte(str), &u)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(u.Name, u.Age, u.Address)
}
