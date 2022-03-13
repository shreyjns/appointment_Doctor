package Appointment

import (
	//"encoding/json"

	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       uint   `json:"id"gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}

var DB *gorm.DB

func Connect() {
	//Conecting DataBase
	connection, err := gorm.Open(mysql.Open("root:2205@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	DB = connection
	connection.AutoMigrate(&User{}, &Doctor{}, &Appoinment{})
	//connection.Migrator().CreateConstraint(&Doctor{}, "Doctor")
	//connection.Migrator().CreateConstraint(&Doctor{}, "fk_availability_doctor")
}

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	//Registration
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	//Performing Login
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user User

	DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "incorrect email",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
func UserL(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user User

	DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}
func Logout(c *fiber.Ctx) error {
	//Logout
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

/////////////
type Doctor struct {
	Id           int    `json:"id"gorm:"primary_key"`
	FullName     string `json:"name"`
	Availability string `json:"Availability"`
	//Availabilitys []Availability `json:"availabilitys`
}
type Appoinment struct {
	Id           int    `json:"id"gorm:"primary_key"`
	Availability string `json:"Availability"`
}

func CreateDoctor(c *fiber.Ctx) error {
	var data Doctor

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	DB.Create(&data)
	return c.JSON(fiber.Map{
		"message": "success",
	})

}
func GetAvailability(c *fiber.Ctx) error {
	var aval []Doctor

	DB.Find(&aval)

	return c.JSON(&aval)
}

func BookAppointment(c *fiber.Ctx) error {
	var data1 Appoinment

	if err := c.BodyParser(&data1); err != nil {
		return err
	}
	DB.Create(&data1)
	return c.JSON(fiber.Map{
		"message": "Appointment Booked",
	})

}
