package main

import (
	entity "belajar-golang-fiber/entity"
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	"github.com/stretchr/testify/assert"
)

var engine = mustache.New("./template", ".mustache")


var app = fiber.New(fiber.Config{
	Views: engine,
	ErrorHandler: func (ctx *fiber.Ctx, err error) error  {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.SendString("Error : " + err.Error())
	},
})

func TestRoutingHelloWorld(t *testing.T) {
	app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(bytes))
}

func TestCtx(t *testing.T) {
	app.Get("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name", "Guest")
        return ctx.SendString("Hello " + name)
    })

	request := httptest.NewRequest("GET", "/hello?name=Bagus", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Bagus", string(bytes))


	request = httptest.NewRequest("GET", "/hello", nil)
	response, err = app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Guest", string(bytes))
}


func TestHttpRequest(t *testing.T) {
	app.Get("/request", func(ctx *fiber.Ctx) error {
		first := ctx.Get("firstname")
		last := ctx.Cookies("lastname")
        return ctx.SendString("Hello "+ first + " "+ last)
    })

	request := httptest.NewRequest("GET", "/request", nil)
	request.Header.Set("firstname", "Bagus")
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "Wicaksono"})
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Bagus Wicaksono", string(bytes))
}

func TestRouteParamater(t *testing.T) {
	app.Get("/users/:userId/orders/:orderId", func(ctx *fiber.Ctx) error {
		userId := ctx.Params("userId")
		orderId := ctx.Params("orderId")
        return ctx.SendString("Get Order "+ orderId + " From User "+ userId	)
    })

	request := httptest.NewRequest("GET", "/users/Bagus/orders/10", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Get Order 10 From User Bagus", string(bytes))
}

func TestFormRequest(t *testing.T) {
	app.Post("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
        return ctx.SendString("Hello " + name)
    })

	body := strings.NewReader("name=Bagus")
	request := httptest.NewRequest("POST", "/hello", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Bagus", string(bytes))
}


//go:embed source/contoh.txt
var contohFile []byte
func TestFormUpload(t *testing.T) {
	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		
		if err!= nil {
            return err
        }
		
		err = ctx.SaveFile(file, "./target/"+file.Filename)
		if err!= nil {
			return err
		}

		return ctx.SendString("Upload Success")
    })

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", "contoh.txt")
	assert.Nil(t, err)
	file.Write(contohFile)
	writer.Close()

	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
    bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Upload Success", string(bytes))
}


// Materi Request Body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRequestBody(t *testing.T) {
	app.Post("/login", func(ctx *fiber.Ctx) error {
		body := ctx.Body()

		request := new(LoginRequest)
		err := json.Unmarshal(body, request)
		if err != nil {
			return err
		}

		return ctx.SendString("Hello " + request.Username)
    })

	body := strings.NewReader(`{"username":"Bagus", "password":"rahasia"}`)
	request := httptest.NewRequest("POST", "/login", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Bagus", string(bytes))
}

type RegisterRequest struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
    Name string `json:"name" xml:"name" form:"name"`
}

func TestBodyParser(t *testing.T) {
	app.Post("/register", func(ctx *fiber.Ctx) error {
		request := new(RegisterRequest) //pointer ke register request
		err := ctx.BodyParser(request) // parse request body ke register request
		fmt.Println(request)
		if err!= nil {
            return err
        }

		return ctx.SendString("Register Success " + request.Username)
    })
}

func TestBodyParserJSON(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(`{"username":"Bagus", "password":"rahasia", "name": "Bagus Wicaksono"}`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Bagus", string(bytes))
}

func TestBodyParserForm(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(`username=Bagus&password=rahasia&name=Bagus Wicaksono`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Bagus", string(bytes))
}

func TestBodyParserXml(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(
		`<RegisterRequest>
			<username>Bagus</username>
			<password>rahasia</password>
			<name>Bagus Wicaksono</name>
		</RegisterRequest>`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/xml")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Bagus", string(bytes))
}


func TestResponseJSON(t *testing.T) {
	app.Get("/user", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"username" : "Bagus",
			"name" : "Bagus Wicaksono",
		})
    })

	request := httptest.NewRequest("GET", "/user", nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"Bagus Wicaksono","username":"Bagus"}`, string(bytes))
}

func TestDownloadFile(t *testing.T) {
	app.Get("/download", func(ctx *fiber.Ctx) error {
		return ctx.Download("./source/contoh.txt", "contoh.txt")
    })

	request := httptest.NewRequest("GET", "/download", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, `attachment; filename="contoh.txt"`, response.Header.Get("Content-Disposition"))
	
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "ini data ku", string(bytes))
}


func TestRoutingGroup(t *testing.T) {
	helloWordl := func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	}

	api := app.Group("/api")
	api.Get("/hello", helloWordl)
	api.Get("/world", helloWordl)

	web := app.Group("/web")
	web.Get("/hello", helloWordl)
	web.Get("/world", helloWordl)

	request := httptest.NewRequest("GET", "/api/hello", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(bytes))
}

func TestStatic(t *testing.T) {
	app.Static("/public", "./source")

	request := httptest.NewRequest("GET", "/public/contoh.txt", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "ini data ku", string(bytes))
}

func TestErrorHandler(t *testing.T) {
	app.Get("/error", func(ctx *fiber.Ctx) error {
		return errors.New("ups")
	})

	request := httptest.NewRequest("GET", "/error", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 500, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t,"Error : ups", string(bytes))
}


func TestView(t *testing.T) {
	app.Get("/view", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"title" : "Hello Title",
			"header" : "Hello Header",
			"content" : "Hello Content",
		}) 
	})

	request := httptest.NewRequest("GET", "/view", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), "Hello Title")
	assert.Contains(t, string(bytes), "Hello Header")
	assert.Contains(t, string(bytes), "Hello Content")
}

func TestClient(t *testing.T) {
	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	agent := client.Get("https://example.com")
	status, response, errors := agent.String()
	assert.Nil(t, errors)
	assert.Equal(t, 200, status)
	assert.Contains(t, response, "Example Domain")

	fmt.Println(response)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := entity.UserLogs{
			UserId: "1",
			Action: "Test Action",
		}

		err := db.Create(&userLog).Error
		assert.Nil(t, err)

		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog.ID)
	}
}

