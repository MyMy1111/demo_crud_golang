package main

import (
	"gorm/models"
	"gorm/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Page struct {
	Body				string		`json:"body"`
	Page_Id			string		`json:"page_id"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreatePage(context *fiber.Ctx) error {
	page := Page{}

	err := context.BodyParser(&page)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
			return err
	}

	err = r.DB.Create(&page).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "count not create a page"})
			return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "page created successfully"})

	return nil
}

func (r *Repository) GetPages(context *fiber.Ctx) error {
	pageModels := &[]models.Pages{}

	err := r.DB.Find(pageModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "count not get a page"})
		return err
	}
	
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "fetched successfully",
		"data": pageModels,
	})

	return nil
}

func (r *Repository) UpdatePage(context *fiber.Ctx) error {
	pageModel := models.Pages{}
	id := context.Params(("id"))
	body := models.UpdatePages{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
			return nil
	}

	err := context.BodyParser(&body)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
			return err
	}

	res := r.DB.First(&pageModel, id)
	if res.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete page",
		})
		return res.Error
	}
	
	pageModel.Body = body.Body
	pageModel.Page_id = body.Page_id

	r.DB.Save(&pageModel)
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "page updated successfully"})
		
	return nil 
}

func (r *Repository) DeletePage(context *fiber.Ctx) error {
	pageModel := models.Pages{}
	id := context.Params(("id"))
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
			return nil
	}

	err := r.DB.Delete(pageModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete page",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "page delete successfully",
	})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/pages", r.CreatePage)
	api.Get("/pages", r.GetPages)
	api.Put("/pages/:id", r.UpdatePage)
	api.Delete("/pages/:id", r.DeletePage)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config {
		Host:	os.Getenv("DB_HOST"),
		Port:	os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("couldn't connect to DB")
	}

	err = models.MigratePages(db)
	if err != nil {
		log.Fatal("couldn't migrate db")
	}

	r := Repository {
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":1111")
}

