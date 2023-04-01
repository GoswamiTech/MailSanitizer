package tokenListener

import (
	"fmt"
	"log"
	"net/http"
)

var (
	server = &http.Server{
		Addr:    ":80",
		Handler: nil,
	}
)

func Start(ch chan string) {
	log.Printf("start server:%v\n", server.Addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintf(w, "code :%s\n", code)
		fmt.Fprint(w, "Authorization is successful. Now go and grab a coffee and some snacks and watch the senitization of your Inbox.")
		ch <- code
	})
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("[ERROR] Lisnter not statred :%v\n", err)
		}
	}()
	log.Println("TokenListener server is started")
}

func Shutdown() {
	server.Close()
}
