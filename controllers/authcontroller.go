package controllers

import (
	"errors"
	"golang_session_login/config"
	"golang_session_login/entities"
	"golang_session_login/libraries"
	"golang_session_login/models"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

var userModel = models.NewUserModel()
var validation = libraries.NewValidation()

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)
	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			data := map[string]interface{}{
				"nama_lengkap": session.Values["nama_lengkap"],
			}

			t, _ := template.ParseFiles("views/index.html")
			t.Execute(w, data)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Proses login
		r.ParseForm()
		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
		errorMessages := validation.Struct(UserInput)
		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
			}
			t, _ := template.ParseFiles("views/login.html")
			t.Execute(w, data)
		} else {

			var user entities.User
			userModel.Where(&user, "username", UserInput.Username)

			var message error
			if user.Username == "" {
				message = errors.New("invalid username or password")
			} else {
				// Check password
				errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
				if errPassword != nil {
					message = errors.New("invalid username or password")
				}
			}
			if message != nil {
				data := map[string]interface{}{
					"error": message,
				}

				t, _ := template.ParseFiles("views/login.html")
				t.Execute(w, data)
			} else {
				// set session
				session, _ := config.Store.Get(r, config.SESSION_ID)

				session.Values["loggedIn"] = true
				session.Values["email"] = user.Email
				session.Values["username"] = user.Username
				session.Values["nama_lengkap"] = user.NamaLengkap

				session.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)
	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("views/register.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Register process

		// Mengambil form input
		r.ParseForm()
		user := entities.User{
			NamaLengkap:     r.Form.Get("nama_lengkap"),
			Email:           r.Form.Get("email"),
			Username:        r.Form.Get("username"),
			Password:        r.Form.Get("password"),
			ConfirmPassword: r.Form.Get("confirm_password"),
		}

		errorMessages := validation.Struct(user)
		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
				"user":       user,
			}
			t, _ := template.ParseFiles("views/register.html")
			t.Execute(w, data)
		} else {
			// hash password
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashPassword)

			// insert ke database
			userModel.Create(user)
			data := map[string]interface{}{
				"pesan": "Registration success",
			}
			t, _ := template.ParseFiles("views/register.html")
			t.Execute(w, data)
		}
	}
}
