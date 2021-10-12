package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func home(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Check password and username
	if username != "Egor" || password != "123123" {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		"Egor S",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	// Create cookie with jwt
	c.SetCookie(&http.Cookie{
		Name:  "Authorization",
		Value: t,
	})

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func sayhello(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.Render(http.StatusOK, "restricted", fmt.Sprintf("Welcome %s !", name))
}

func secretinfo(c echo.Context) error {
	return c.Render(http.StatusOK, "restricted", "Information only for those who are logged in")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	t := &Template{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}

	e.Renderer = t

	//Accesible routes
	e.GET("/", home)
	e.POST("/login", login)

	// Restricted group
	r := e.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:      &jwtCustomClaims{},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:Authorization",
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("/sayhello", sayhello)
	r.GET("/secretinfo", secretinfo)

	e.Logger.Fatal(e.Start(":1323"))
}
