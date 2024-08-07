package main

import "text/template"

// for car details
type Car struct {
	ID                int
	Name              string
	ManufacturingYear int
	Price             float64
}

// handling templates and sessions
var (
	tpl        *template.Template
	dbSessions = make(map[string]SessionData)
)

// to store the session data
type SessionData struct {
	Username string
	Role     string
}

// for parsing the template
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}
