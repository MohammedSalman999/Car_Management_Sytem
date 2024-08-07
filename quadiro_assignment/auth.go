package main

import (
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func alreadyLoggedIn(req *http.Request) (bool, string) {
	c, err := req.Cookie("session")
	if err != nil {
		log.Println("No session cookie found")
		return false, ""
	}
	sessionData, ok := dbSessions[c.Value]
	if !ok {
		log.Println("Session ID not found in dbSessions")
		return false, ""
	}
	log.Println("User logged in:", sessionData.Username, "with role:", sessionData.Role)
	return true, sessionData.Role
}

func createSession(w http.ResponseWriter, un, role string) {
	sID := uuid.NewV4().String()
	c := &http.Cookie{
		Name:  "session",
		Value: sID,
		Path:  "/",
	}
	http.SetCookie(w, c)
	dbSessions[c.Value] = SessionData{Username: un, Role: role}
	log.Println("Session created for user:", un, "with role:", role)
}

func userLogin(w http.ResponseWriter, req *http.Request) {
	loggedIn, role := alreadyLoggedIn(req)
	if loggedIn && role == "user" {
		log.Println("Already logged in as user, redirecting to /person")
		http.Redirect(w, req, "/person", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		var hashPassword []byte
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", un).Scan(&hashPassword)
		if err != nil {
			http.Error(w, "Username does not exist", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(hashPassword, []byte(p))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		createSession(w, un, "user")
		log.Println("Login successful, redirecting to /person")
		http.Redirect(w, req, "/person", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.tmpl", "Assignment for Quadiro Technologies - User Login")
}

func adminLogin(w http.ResponseWriter, req *http.Request) {
	loggedIn, role := alreadyLoggedIn(req)
	if loggedIn && role == "admin" {
		log.Println("Already logged in as admin, redirecting to /admin/dashboard")
		http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		var hashPassword []byte
		err := db.QueryRow("SELECT password FROM admins WHERE username = ?", un).Scan(&hashPassword)
		if err != nil {
			http.Error(w, "Username does not exist", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(hashPassword, []byte(p))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		createSession(w, un, "admin")
		log.Println("Login successful, redirecting to /admin/dashboard")
		http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.tmpl", "Assignment for Quadiro Technologies - Admin Login")
}

func signup(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		fn := req.FormValue("firstname")
		ln := req.FormValue("lastname")
		role := req.FormValue("role")

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
			return
		}

		if role == "admin" {
			_, err := db.Exec("INSERT INTO admins (username, password, firstname, lastname) VALUES (?, ?, ?, ?)", un, hashPassword, fn, ln)
			if err != nil {
				http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}
		} else {
			_, err := db.Exec("INSERT INTO users (username, password, firstname, lastname) VALUES (?, ?, ?, ?)", un, hashPassword, fn, ln)
			if err != nil {
				http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}
		}

		log.Println("Signup successful, redirecting to /login")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signin.tmpl", nil)
}

// lohoutHanlder handels the logging out handling
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn, role := alreadyLoggedIn(r)
	if !loggedIn && role == "admin" {
		log.Println("Already logged in as admin, redirecting to /admin/dashboard")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Cookie delete karna
	http.SetCookie(w, &http.Cookie{
		Name:   "session", // Apne session cookie ka naam yahan daalein
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Cookie ko expire karne ke liye
	})

	// Redirect karna
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
