package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"

	"github.com/react-golang/database"
	"github.com/react-golang/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	database.DB.Create(&user)

	return
}

func Login(c *fiber.Ctx) {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return
}

func User(c *fiber.Ctx) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return
}

func Logout(c *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return
}
