package main // import "webdav-server"

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func main() {
	storagePath := "./storage"
	certPath := "./pem/cert.pem"
	keyPath := "./pem/key.pem"

	// When use "/" not "/webdav" in http.HandleFunc, srv.Prefix should be removed.
	srv := &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(storagePath),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				fmt.Printf("WebDAV %s: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				fmt.Printf("WebDAV %s: %s \n", r.Method, r.URL)
			}
		},
	}

	// Slash must be inputed to end of path in http.HandleFunc
	http.HandleFunc("/webdav/", func(w http.ResponseWriter, r *http.Request) {
		// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		if username != "davuser" || password != "pass" {
			srv.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="BASIC REALM"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	})

	// Check cert files are existing. For HTTPS only.
	_, errCert := os.Stat(certPath)
	_, errKey := os.Stat(keyPath)
	if errCert != nil || errKey != nil {
		log.Fatal("cert.pem or key.pem is not found. Check pem directory.")
		return
	}

	// For HTTP and HTTPS
	// go http.ListenAndServeTLS(":8443", certPath, keyPath, nil)
	// http.ListenAndServe(":8080", nil)

	// For HTTPS only
	http.ListenAndServeTLS(":8443", certPath, keyPath, nil)
}
