# sessions
golang sessions

### Example Code    
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
	sessions.Register("memory", sessions.NewMemoryStore(time.Second * 1800))
	manager, _ := sessions.NewCookieManager("memory", "sessionid", 1800)
	manager.GC(time.Second * 1800)
	sessions.WithBackground(manager)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session := sessions.Start(w, r)
		session.Set("token", sessions.NewUUID())
		token, _ := session.Get("token")
		fmt.Fprintf(w, "session: %v, token: %v", sessions.SessionID(), token)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessions.Destory(w, r)
		fmt.Fprintf(w, "clear session")
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
```