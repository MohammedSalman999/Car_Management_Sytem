package main

import (
	"net/http"
)

func main() {
	InitDB()
	defer db.Close()
	// for handling the index page
	http.HandleFunc("/", index)

	// admin login and dashboard ko handle karne ke liye
	http.HandleFunc("/admin/login", adminLogin)
	http.HandleFunc("/admin/dashboard", adminDashboard)

	// user login and person page ko handle karne ke liye
	http.HandleFunc("/user/login", userLogin)
	http.HandleFunc("/person", personHandler)

	// logout handler
	http.HandleFunc("/logout", logoutHandler)

	// crud operation handler
	http.HandleFunc("/admin/car", carHandler)
	http.HandleFunc("/admin/car/update", updateCarHandler)
	http.HandleFunc("/admin/car/delete", deleteCarHandler)

	http.HandleFunc("/seecars", getAllCarsHandler)
	http.HandleFunc("/signin", signup)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)
}
