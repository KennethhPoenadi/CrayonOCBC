package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Kenneth/crayon/repo/users"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use: "db",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := sql.Open("sqlite3", "file:db.sqlite")
			if err != nil {
				log.Println(err)
				return
			}
			defer db.Close()

			ctx := cmd.Context()
			udb := users.NewRepo(db)

			initTable(ctx, udb)
			listTables(ctx, udb)
		},
	}

	rootCmd.AddCommand(cmd)
}

func listTables(ctx context.Context, udb *users.Repo) {
	userList, err := udb.List(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, user := range userList {
		fmt.Println(user.Id, user.Name, user.Age)
	}
}

func initTable(ctx context.Context, udb *users.Repo) {
	err := udb.CreateTable(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	err = udb.Insert(ctx, "Alvin", 21)
	if err != nil {
		log.Println(err)
		return
	}

	err = udb.Insert(ctx, "John", 21)
	if err != nil {
		log.Println(err)
		return
	}
}
