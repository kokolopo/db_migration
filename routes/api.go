package routes

import (
	tablenamehandlers "db_migration/Handlers"

	"github.com/gofiber/fiber/v2"
)

func APIRoutes(app *fiber.App, tablenameHandler *tablenamehandlers.TablenameHandler) {
	app.Get("source/tablesname", tablenameHandler.TablesNameDB1)
	app.Get("target/tablesname", tablenameHandler.TablesNameDB2)
	app.Get(":source/describe-table/:table_name", tablenameHandler.DescribeColumn)
	app.Get(":source/data/:table_name", tablenameHandler.GetDataInTable)
	// app.Post("source/field/:tablename", tablenameHandler.GetDataInTable)
	app.Post("migration/:table_source/to/:table_target", tablenameHandler.MigrationExec)
	app.Delete("migration/releted_data", tablenameHandler.DeleteReletedData)
	app.Post("migration/releted_data", tablenameHandler.MigrateRelatedData)

}
