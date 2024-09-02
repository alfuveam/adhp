package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"time"

	"github.com/alfuveam/tcc/backend/config"
	"github.com/alfuveam/tcc/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type ErrorResponseCreateAccount struct {
	Completename string `json:"error_completename"`
	// Cpf          string `json:"cpf"`
	Phone           string `json:"error_phone"`
	Email           string `json:"error_email"`
	Password        string `json:"error_password"`
	Country         string `json:"error_country"`
	OnCreateAccount string `json:"error_on_create_account"`
}

type MyCustomClaims struct {
	User models.User
	jwt.RegisteredClaims
}

type ErroOnLogin struct {
	Error string `json:"error_on_login"`
}

type DashboardResponse struct {
	Name string `json:"name"`
}

func ValidateJWT(tokenString string) (models.User, bool, error) {
	var user models.User
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.MySigningKey), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		// log.Fatal(err)
		return user, false, err
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		return claims.User, true, nil
	} else {
		// log.Fatal("unknown claims type, cannot proceed")
		return user, false, errors.New("unknown claims type, cannot proceed")
	}
}

func CreateAndSignJWT(user *models.User) (string, error) {

	// Create claims while leaving out some of the optional fields
	claims := MyCustomClaims{
		*user,
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(config.MySigningKey))
	_, _, err := ValidateJWT(ss)

	if err != nil {
		log.Println("CreateAndSignJWT: ", err)
	}

	return ss, nil
}

func OnLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid request payload",
		})
		return
	}

	// userID, hasEmail, err := req.OnLoginCheckEmailAndPasswordId(db)
	// if !hasEmail {
	// 	res := ErroOnLogin{
	// 		Error: err.Error(),
	// 	}
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }

	user, hasEmail, err := req.OnLoginCheckEmailAndPasswordUser(db)
	if !hasEmail {
		res := ErroOnLogin{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid request payload",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(config.SaltDB+req.Password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid password",
		})
		return
	}

	//	put no salt password in public token
	user.Password = req.Password

	token, err := CreateAndSignJWT(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid token",
		})
		return
	}

	res := LoginResponse{
		Message: "Login successful",
		Token:   token,
	}
	json.NewEncoder(w).Encode(res)
}

func RegisterUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req models.User
	var plainPassword string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid request payload",
		})
		return
	}

	var errorRes ErrorResponseCreateAccount
	hasError := false
	if len(req.Completename) > 50 || len(req.Completename) < 5 {
		errorRes.Completename = "The correct size of name it's betwent 50 and 5."
		hasError = true
	}

	// if len(req.Cpf) > 11 {
	// 	errorRes.Cpf = "Incorrect size of cpf."
	// 	hasError = true
	// }

	if len(req.Phone) > 50 && len(req.Phone) < 6 {
		errorRes.Phone = "Incorrect size of phone."
		hasError = true
	}

	sampleRegexp := regexp.MustCompile(`\d+`)
	match := sampleRegexp.MatchString(req.Phone)
	if !match {
		errorRes.Phone = "Incorrect values in phone number."
		hasError = true
	}

	if len(req.Email) > 100 {
		errorRes.Email = "Incorrect size of Email."
		hasError = true
	}

	if len(req.Password) > 100 || len(req.Password) < 9 {
		errorRes.Password = "The correct size of password it's betwent 100 and 9."
		hasError = true
	}

	plainPassword = req.Password
	req.Password = config.SaltDB + req.Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword: %v", err)
		errorRes.OnCreateAccount = "Error on create account contact the administrator"
		hasError = true
	}

	req.Password = string(bytes)

	if len(req.Country) != 3 {
		errorRes.Country = "Incorrect country."
		hasError = true
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		errorRes.Email = err.Error()
		hasError = true
	}

	if hasError {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	req.CreateAt = time.Now()
	req.LastLogin = time.Now()

	if err := req.CreateUser(db); err != nil {
		log.Println("req.CreateUser: %v", err)
		errorRes.OnCreateAccount = "Error on create account contact the administrator"
		hasError = true
	}

	if hasError {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	req.Password = plainPassword
	token, err := CreateAndSignJWT(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid token",
		})
		return
	}

	res := LoginResponse{
		Message: "Login successful",
		Token:   token,
	}
	json.NewEncoder(w).Encode(res)
}

func OnLoadDashBoard(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	user, ok := r.Context().Value(config.MySigningKey).(models.User)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid user ID",
		})
		return
	}

	// var req models.User
	// res := LoginResponse{
	// 	// Message: "Login successful",
	// 	// Token:   "example_token",
	// 	userID: userID,
	// }
	log.Println("show ID of user: ", user)

	res := DashboardResponse{
		Name: "Josh",
	}
	// json.NewEncoder(w).Encode("userID: " + userID)
	json.NewEncoder(w).Encode(res)
}
