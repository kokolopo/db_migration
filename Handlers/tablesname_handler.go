package handlers

import (
	"db_migration/service"

	"github.com/gofiber/fiber/v2"
)

type H map[string]interface{}

type TablenameHandler struct {
	tablenameService service.IService
}

func NewTablenameHandler(TablenameService service.IService) *TablenameHandler {
	return &TablenameHandler{TablenameService}
}

func (h *TablenameHandler) TablesNameDB1(c *fiber.Ctx) error {
	data, err := h.tablenameService.GetTablesNameDB1()
	if err != nil {

		return c.JSON(H{
			"error": err,
		})
	}
	return c.JSON(data)
}

func (h *TablenameHandler) TablesNameDB2(c *fiber.Ctx) error {
	data, err := h.tablenameService.GetTablesNameDB2()
	if err != nil {

		return c.JSON(H{
			"error": err,
		})
	}
	return c.JSON(data)
}

func (h *TablenameHandler) DescribeColumn(c *fiber.Ctx) error {
	tableName := c.Params("table_name")
	source := c.Params("source")

	if source == "source" {
		data, err := h.tablenameService.GetColumnsDescDB1(tableName)
		if err != nil {
			return c.JSON(H{
				"error": err,
			})
		}
		return c.JSON(data)
	}

	if source == "target" {
		data, err := h.tablenameService.GetColumnsDescDB2(tableName)
		if err != nil {
			return c.JSON(H{
				"error": err,
			})
		}
		return c.JSON(data)
	}

	return c.JSON(H{
		"error": "pilih sumber data dari source atau tager!",
	})
}

func (h *TablenameHandler) MigrationExec(c *fiber.Ctx) error {
	tabelSource := c.Params("table_source")
	tabelTarget := c.Params("table_target")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	email := c.Query("email")

	isMigrate, err := h.tablenameService.MigrateTable(tabelSource, tabelTarget, page, limit, email)
	if err != nil {
		return c.JSON(H{
			"error": err,
		})
	}

	return c.JSON(H{
		"is_migrate": isMigrate,
	})
}

func (h *TablenameHandler) DeleteReletedData(c *fiber.Ctx) error {
	email := c.Query("email")

	result, err := h.tablenameService.DeleteReletedData(email)
	if err != nil {
		return c.JSON(H{
			"error": err,
		})
	}

	return c.JSON(H{
		"is_deleted": result,
	})
}

func (h *TablenameHandler) MigrateRelatedData(c *fiber.Ctx) error {
	email := c.Query("email")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.tablenameService.MigrateReletedData(email, page, limit)
	if err != nil {
		return c.JSON(H{
			"error": err,
		})
	}

	return c.JSON(H{
		"is_migrated": result,
	})
}

func (h *TablenameHandler) GetDataInTable(c *fiber.Ctx) error {
	tablename := c.Params("table_name")
	source := c.Params("source")
	page := c.QueryInt("page", 1)
	data, err := h.tablenameService.GetDataTable(tablename, page, source)
	if err != nil {
		return c.JSON(H{
			"error": err,
		})
	}
	return c.JSON(data)
}
