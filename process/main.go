package process

import (
	"appengine"
	"appengine/datastore"
	"data"
	"net/http"
	"template"
	"time"
	"util"
)

func init() {
	http.HandleFunc("/", topPage)
	http.HandleFunc("/new", newComment)
}

type Message_View struct {
	No int
	Text string		//本文
	Author string	//作成者
	Date string		//投稿日
}

func topPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Messgae")
	count, _ := q.Count(c)

	message := make([]data.Message, 0, count)

	if _, err := q.GetAll(c, &message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//表示用のデータの作成
	view_data := make([]Message_View, count)
	for pos, data := range message {
		view.data[pos].No = data.No
		view_data[pos].Text = data.Text
		view_data[pos].Author = data.Author
		view_data[pos].Date = util.DateToString(data.Date + (9 * 3600 * 1e6))
	}

	var t = template.Must(template.New("html").ParseFile("html/main.html"))

	if err := t.Execute(w, view_data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}