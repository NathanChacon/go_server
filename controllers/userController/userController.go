package userController

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"test.com/events/model/userModel"
)

var idExpirationDate time.Time = time.Now().Add(24 * time.Hour)

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func generateLoginJwt(userId string) (string, error) {
	var secret = os.Getenv("JWT_SECRET")
	var jwtKey = []byte(secret)

	claims := &jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(idExpirationDate),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func CreateAccount(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	var newUser userModel.User

	errJson := decoder.Decode(&newUser)

	if errJson != nil {
		http.Error(response, "Wrong payload", http.StatusBadRequest)
		return
	}

	emailAlreadyRegistered, emailErr := userModel.EmailExists(newUser.Email)

	if emailErr != nil {
		http.Error(response, "Internal error email", http.StatusInternalServerError)
		return
	}

	if emailAlreadyRegistered {
		http.Error(response, "Email already in use", http.StatusBadRequest)
		return
	}

	if len(newUser.Name) < 6 || len(newUser.Name) > 20 {
		http.Error(response, "Name must be between 6 and 20 characters", http.StatusBadRequest)
		return
	}

	if len(newUser.Password) < 6 || len(newUser.Password) > 16 {
		http.Error(response, "Password must be between 6 and 16 characters", http.StatusBadRequest)
		return
	}

	if !isValidEmail(newUser.Email) {
		http.Error(response, "Invalid email format", http.StatusBadRequest)
		return
	}

	hashedPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if errBcrypt != nil {
		http.Error(response, "Internal error bcrypt", http.StatusInternalServerError)
		return
	}

	newUser.Password = string(hashedPassword)

	createdUserId, createUserError := userModel.CreateUser(newUser)

	if createUserError != nil {
		http.Error(response, "Internal error create user", http.StatusInternalServerError)
		return
	}

	tokenString, err := generateLoginJwt(createdUserId)

	if err != nil {
		http.Error(response, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(response, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  idExpirationDate,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	json.NewEncoder(response).Encode(map[string]interface{}{
		"user_id": createdUserId,
		"message": "Account created and logged in",
	})
}

func Login(response http.ResponseWriter, request *http.Request) {
	type loginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginData loginPayload

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&loginData)

	if err != nil {
		http.Error(response, "Invalid login payload", http.StatusBadRequest)
		return
	}

	user, getUserError := userModel.GetUserByEmail(loginData.Email)

	if getUserError != nil {
		http.Error(response, "User not found", http.StatusNotFound)
		return
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))

	if passwordError != nil {
		http.Error(response, "Wrong password", http.StatusUnauthorized)
		return
	}

	tokenString, tokenError := generateLoginJwt(user.ID)

	if tokenError != nil {
		http.Error(response, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(response, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  idExpirationDate,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	json.NewEncoder(response).Encode(map[string]interface{}{
		"user_id": user.ID,
		"message": "Account logged in",
	})

}
