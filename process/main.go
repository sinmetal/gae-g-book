package process

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Guest struct {
	Name string
	Date time.Time
}

type Guest_View struct {
	Name string
	Date string
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/write", write)
	http.HandleFunc("/list", list)
}

const inputForm = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>名前の登録</title>
</head>
<body>
<form method="POST" action="write">
	<label>お名前<input type="text" name="name" /></label>
	<input type="submit">
</form>
</body>
</html>
`

const guestTemplateHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>登録者リスト</title>
</head>
<body>
	<table border="1">
		<tr>
			<th>名前</th>
			<th>登録日時</th>
		</tr>
		{{range .}}
			<tr>
				{{if .Name}}
					<td>{{.Name|html}}</td>
				{{else}}
					<td>名無し</td>
				{{end}}
				{{if .Date}}
					<td>{{.Date|html}}</td>
				{{else}}
					<td> - </td>
				{{end}}
			</tr>
		{{end}}
	</table>
</body>
</html>
`

var guestTemplate = template.Must(template.New("guest").Parse(guestTemplateHTML))

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", inputForm)
}

func write(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Not Found")
		return
	}

	c := appengine.NewContext(r)

	// Datastoreへの書き込み
	var g Guest
	g.Name = r.FormValue("name")
	g.Date = time.Now()
	if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Guest", nil), &g); err != nil {
		http.Error(w, "Internal Server Error : "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func list(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Guest").Order("Date")
	count, err := q.Count(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guests := make([]Guest, 0, count)
	if _, err := q.GetAll(c, &guests); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("guests len=%d", len(guests))

	guest_views := make([]Guest_View, count)
	for pos, guest := range guests {
		guest_views[pos].Name = fmt.Sprintf("%s", guest.Name)

		//localTime := time.SecondsToLocalTIme(int64(guest.Date) / 1000000)
		localTime := guest.Date.Format(time.ANSIC)

		//guest_views[pos].Date = fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", localTime.Year, localTIme.Month, localTime.Day, localTime.Hour, localTime.Minute, localTIme.Second)
		guest_views[pos].Date = localTime
	}

	if err := guestTemplate.Execute(w, guest_views); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
