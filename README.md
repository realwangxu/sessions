# sessions
golang sessions

### Example   
```go
package main

import (
	"fmt"
	"github.com/realwangxu/sessions"
	"net/http"
	"time"
	"log"
)

func main() {
	manager, _ := sessions.NewCookieManager("memory", time.Minute * 30, "sessionid", 1800)
	sessions.WithBackground(manager)
	sessions.GC(time.Minute * 30)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session := sessions.Start(w, r)
		fmt.Fprintf(w, "session: %v", sessions.SessionID())
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessions.Destory(w, r)
		fmt.Fprintf(w, "clear session")
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
```