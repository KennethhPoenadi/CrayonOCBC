package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kenneth/crayon/repo/users"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// magic function yg auto jalan ketika package diimport
func init() {
	var address string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Running a server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config file used:", config)
			if address == "" {
				address = viper.GetString("server.address")
			}
			fmt.Println("Server is running on", address)

			db, err := sql.Open("sqlite3", "file:db.sqlite")
			if err != nil {
				log.Println(err)
				return
			}
			defer db.Close()

			udb := users.NewRepo(db) //bungkus koneksi db menjadi repository user.

			s := http.Server{
				Addr:    address,
				Handler: ChiRouter(udb),
			}

			if err := s.ListenAndServe(); err != nil {
				log.Println(err)
			}
		},
	}
	cmd.Flags().StringVarP(&address, "address", "a", "", "Address for server to listen to")

	rootCmd.AddCommand(cmd)
}

func GinRouter() http.Handler {
	handler := gin.New()
	handler.Use(func(ctx *gin.Context) {
		fmt.Println(ctx.Request.Method, ctx.Request.URL.Path)
	})

	handler.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello Wolrd!")
	})

	handler.GET("/hi", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hi")
	})

	return handler
}

// parameternya diubah jadi nerima udb, buat ngambil user dari database
func ChiRouter(user *users.Repo) http.Handler {
	handler := chi.NewMux()

	handler.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			log.Println("From middleware: ", r.Method, r.URL.Path)
		})
	})

	handler.Group(func(r chi.Router) {
		r.Use(AUthorizationMw)

		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			u := User{
				Name:    "John",
				Age:     17,
				Address: "BSD",
			}
			if err := json.NewEncoder(w).Encode(u); err != nil {
				log.Println(err)
			}
		})

		r.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
			var u User
			err := json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "unable to decode payload")
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		r.Get("/allusers", func(w http.ResponseWriter, r *http.Request) {
			userList, err := user.List(r.Context())
			if err != nil {
				log.Println(err)
				http.Error(w, "failed to get users", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(userList); err != nil {
				log.Println(err)
			}

			json.NewEncoder(w).Encode(userList)
		})

		r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
			var val struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}

			if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
			}

			if err := user.Insert(r.Context(), val.Name, val.Age); err != nil {
				http.Error(w, "insert error", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusCreated)
		})

		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(chi.URLParam(r, "id"))

			if err != nil {
				http.Error(w, "id not found", http.StatusBadRequest)
			}

			userId, error := user.Get(r.Context(), id)

			if error != nil {
				http.Error(w, "failed to get the user id", http.StatusBadRequest)
			}

			json.NewEncoder(w).Encode(userId)

		})
	})

	return handler
}

func AUthorizationMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
