package server

import (
	"ex02/postDB"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

const (
	validUsername = "admin"
	validPassword = "admin"
)

type HtmlResponsePosts struct {
	Total      int
	Posts      []postDB.Post
	PrevPage   int
	NextPage   int
	LastPage   int
	IsLastPage bool
}

type HtmlResponseOnePost struct {
	Title        string
	Content      string
	PubDate      time.Time
	PreviousPage int
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	db := postDB.ConnectDB()
	defer db.Close()

	count, err := postDB.CountPosts(db)
	if err != nil {
		http.Error(w, "Error counting posts", http.StatusInternalServerError)
		return
	}
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
	}
	if err != nil || page < 1 || page > (int(count)+postDB.PageSize-1)/postDB.PageSize {
		http.Error(w, "Invalid 'page' value", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		posts, err := postDB.GetPosts(db, page)
		if err != nil {
			http.Error(w, "Error receiving data from DB", http.StatusInternalServerError)
			return
		}

		response := HtmlResponsePosts{
			Total:      int(count),
			Posts:      posts,
			PrevPage:   page - 1,
			NextPage:   page + 1,
			LastPage:   (int(count)/postDB.PageSize + 1),
			IsLastPage: page == (int(count)/postDB.PageSize + 1),
		}

		tmpl, err := template.ParseFiles("templates/posts.html")
		if err != nil {
			http.Error(w, "Template load error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for i := range response.Posts {
			response.Posts[i].Link = "posts/post?id=" + strconv.Itoa(int(response.Posts[i].ID)) + "&page=" + strconv.Itoa(page)
		}

		if err := tmpl.Execute(w, response); err != nil {
			http.Error(w, "HTML generation error", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Only method GET is allowed", http.StatusMethodNotAllowed)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	pageStr := r.URL.Query().Get("page")
	db := postDB.ConnectDB()
	defer db.Close()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' value", http.StatusBadRequest)
		return
	}

	previousPage := 1
	if pageStr != "" {
		previousPage, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid 'page' value", http.StatusBadRequest)
			return
		}
	}

	if r.Method == http.MethodGet {
		post, err := postDB.GetOnePost(db, id)
		if err != nil {
			http.Error(w, "Error receiving data from DB", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/one_post.html")
		if err != nil {
			http.Error(w, "Template load error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		response := HtmlResponseOnePost{
			Title:        post.Title,
			Content:      post.Content,
			PubDate:      post.PubDate,
			PreviousPage: previousPage,
		}

		if err := tmpl.Execute(w, response); err != nil {
			http.Error(w, "HTML generation error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Only method GET is allowed", http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == validUsername && password == validPassword {
			session, _ := store.Get(r, "session-name")
			session.Values["authenticated"] = true
			session.Save(r, w)

			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		http.Error(w, "Wrong login or password", http.StatusUnauthorized)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		var post postDB.Post
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")

		db := postDB.ConnectDB()
		defer db.Close()

		postDB.InsertPost(db, &post)

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Template load error", http.StatusInternalServerError)
	}
}
