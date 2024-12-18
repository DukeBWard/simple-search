package routes

import (
	"dukebward/search/db"
	"dukebward/search/utils"
	"dukebward/search/views"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type loginForm struct {
	Email    string `form:email`
	Password string `form:password`
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := loginForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Invalid input</h2>")
	}
	// pointer to user
	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Unauthorized login</h2>")
	}
	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Something went wrong</h2>")
	}
	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	// set context cookie to a pointer to the cookie
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

func LogoutHandler(c *fiber.Ctx) error {
	// clear the cookie and redirect to login
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

// AdminClaims is a struct that will be used to parse the token
type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	// cookie is actually just a string right now
	// its the value from admin (token string)
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}
	// Parse JWT token from cookie and verify signature
	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	// Handle invalid tokens
	if err != nil {
		return c.Redirect("/login", 302)
	}
	// Verify claims and token validity
	// If both pass, allows request to continue via c.Next()
	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}
	// if we get here the cookie or jwt token is invalid
	return c.Redirect("/login", 302)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Can't get settings</h2>")
	}
	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}
