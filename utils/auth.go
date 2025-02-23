package utils

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkCredentials(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func checkCredentials(username, password string) bool {
	file, err := os.ReadFile(".htpasswd")
	if err != nil {
		return false
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == username {
			return bcrypt.CompareHashAndPassword([]byte(parts[1]), []byte(password)) == nil
		}
	}
	return false
}
