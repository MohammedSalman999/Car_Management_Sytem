package main

import (
	"log" // Add this import
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.tmpl", nil)
}

// ye admin dashboard ko handle karega
func adminDashboard(w http.ResponseWriter, req *http.Request) {
	loggedIn, role := alreadyLoggedIn(req)
	if !loggedIn || role != "admin" {
		log.Println("Not logged in as admin, redirecting to /admin/login")
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}

	cars, err := getAllCars()
	if err != nil {
		http.Error(w, "Unable to fetch cars", http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "dashboard.tmpl", cars)
}

func personHandler(w http.ResponseWriter, req *http.Request) {
	loggedIn, role := alreadyLoggedIn(req)
	if !loggedIn || role != "user" {
		log.Println("Not logged in as admin, redirecting to /admin/login")
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}

	cars, err := getAllCars()
	if err != nil {
		http.Error(w, "Unable to fetch cars", http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "person.tmpl", cars)
}

//

func getAllCarsHandler(w http.ResponseWriter, req *http.Request) {
	// Check for allowed request method
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check if session cookie is present and valid
	cookie, err := req.Cookie("session")
	if err != nil {
		// No session cookie found, redirect to login page
		log.Println("No session cookie found, redirecting to /user/login")
		http.Redirect(w, req, "/user/login", http.StatusSeeOther)
		return
	}

	username, ok := dbSessions[cookie.Value]
	if !ok {
		// Invalid session ID, redirect to login page
		log.Println("Invalid session ID, redirecting to /user/login")
		http.Redirect(w, req, "/user/login", http.StatusSeeOther)
		return
	}

	// Fetch all cars from the database
	cars, err := getAllCars()
	if err != nil {
		http.Error(w, "Unable to fetch cars", http.StatusInternalServerError)
		return
	}

	// Render the template with the cars data
	err = tpl.ExecuteTemplate(w, "seecars.tmpl", cars)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}

	// Log successful retrieval and display
	log.Printf("Cars successfully retrieved and displayed for user: %s", username)
}

func carHandler(w http.ResponseWriter, req *http.Request) {
	loggedIn, role := alreadyLoggedIn(req)
	if !loggedIn || role != "admin" {
		log.Println("Not logged in as admin, redirecting to /admin/login")
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}

	switch req.Method {
	case http.MethodGet:
		// Handle GET request (show car records)
		cars, err := getAllCars()
		if err != nil {
			http.Error(w, "Unable to fetch cars", http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "cars.tmpl", cars)

	case http.MethodPost:
		// Add New Car
		name := req.FormValue("name")
		manufacturingYear := req.FormValue("manufacturing_year")
		price := req.FormValue("price")

		_, err := db.Exec("INSERT INTO cars (name, manufacturing_year, price) VALUES (?, ?, ?)", name, manufacturingYear, price)
		if err != nil {
			log.Printf("Error inserting car: %v", err) // Detailed error logging
			http.Error(w, "Unable to create car", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
	}
}

func updateCarHandler(w http.ResponseWriter, req *http.Request) {
	// Check if user is logged in and has admin role
	loggedIn, role := alreadyLoggedIn(req)
	if !loggedIn || role != "admin" {
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}

	// Handle only POST method for updating
	if req.Method == http.MethodPost {
		// Extract form values
		id := req.FormValue("id")
		name := req.FormValue("name")
		manufacturingYear := req.FormValue("manufacturing_year")
		price := req.FormValue("price")

		// Execute the update query
		_, err := db.Exec("UPDATE cars SET name = ?, manufacturing_year = ?, price = ? WHERE id = ?", name, manufacturingYear, price, id)
		if err != nil {
			log.Printf("Error updating car: %v", err)
			http.Error(w, "Unable to update car", http.StatusInternalServerError)
			return
		}

		// Redirect to dashboard after successful update
		http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
	} else {
		// If not a POST request, return method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func deleteCarHandler(w http.ResponseWriter, req *http.Request) {
	loggedIn, role := alreadyLoggedIn(req)
	if !loggedIn || role != "admin" {
		log.Println("Not logged in as admin, redirecting to /admin/login")
		http.Redirect(w, req, "/admin/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		id := req.FormValue("id")
		_, err := db.Exec("DELETE FROM cars WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting car: %v", err) // Detailed error logging
			http.Error(w, "Unable to delete car", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
	}
}
