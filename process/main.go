package process

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"time"
)

type Guest struct {
	Name string
	Date time.Time
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/write", write)
}

const inputForm = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>名前の登録</title>
</head>
<body>
<form method "POST" action="write">
	<label>お名前<input type="text" name="name" /></label>
	<input type="submit">
</form>
</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", inputForm)
}

func write(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Not Fount")
		return
	}

	c := appengine.NewContext(r)

	// Datastoreへの書き込み
	var g Guest
	g.Name = r.FormValue("name")
	g.Date = time.Now()
	if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Guest", nil), &g); err != nil {
		http.Error(w, "Internal Server Error : " + err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}