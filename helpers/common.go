package helpers

import "os"

func GetPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}
