package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/forms"
	"github.com/sabrodigan/bookings-app/internal/models"
	"github.com/sabrodigan/bookings-app/internal/render"
	"log"
	"net/http"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home page hit\n")
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	log.Println("About page hit\n")
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	log.Println("Reservation page hit\n")
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	log.Println("PostReservation page hit\n")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// using the required to validate entry
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)

	if !form.Valid() {
		var data = make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	log.Println("Generals page hit\n")
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	log.Println("Majors page hit\n")
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	log.Println("Availability page hit\n")
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	log.Println("PostAvailability page hit\n")
	TypeWriter("PostAvailability page hit\n", 50)
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	log.Println("Contact page hit\n")
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
func TypeWriter(phrase string, delayMs int) {
	// Simulate a typewriter effect with a custom delay in milliseconds
	text := phrase
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}
func Spinner(delay time.Duration, done chan bool) {
	// Spinner characters
	for {
		for _, r := range `|/-\` {
			select {
			case <-done:
				fmt.Print("\r") // Clear the spinner
				return
			default:
				fmt.Printf("\r%c", r) // Print spinner character
				time.Sleep(delay)
			}
		}
	}
}
