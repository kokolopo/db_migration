package main

import (
	tablenamehandlers "db_migration/Handlers"
	"db_migration/config"
	"db_migration/entity"
	"db_migration/routes"
	"db_migration/service"
	"db_migration/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func main() {
	// Memuat variabel environment dari file .env
	utils.LoadEnv()

	// Inisialisasi koneksi ke dua PostgreSQL database
	dbPostgres1, dbPostgres2, err := config.InitDBs()
	if err != nil {
		log.Fatal("Failed to connect to databases: ", err)
	}

	// var DB1_NAME1 = os.Getenv("DB1_NAME")
	handlers := initialApp(dbPostgres1, dbPostgres2)

	// Inisialisasi Fiber
	app := fiber.New()
	app.Use(cors.New())

	routes.APIRoutes(app, handlers.TablenameHandler)

	log.Fatal(app.Listen(":3000"))

}

type APP struct {
	TablenameHandler *tablenamehandlers.TablenameHandler
}

func initialApp(db1 *gorm.DB, db2 *gorm.DB) APP {

	var handlers = APP{}

	TNRepo := entity.NewTablesNameRepository(db1, db2)
	TNservice := service.NewUserService(TNRepo)
	TNHandler := tablenamehandlers.NewTablenameHandler(TNservice)

	handlers.TablenameHandler = TNHandler

	return handlers
}
