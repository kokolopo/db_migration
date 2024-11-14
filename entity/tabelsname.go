package entity

import (
	"db_migration/utils"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

var TablesNameSQL = `SELECT ROW_NUMBER() OVER() AS no, tablename FROM pg_tables WHERE schemaname = ?`
var TablesNameMySQL = `SELECT @rownum := @rownum + 1 AS no, table_name as tablename FROM information_schema.tables, (SELECT @rownum := 0) r WHERE table_schema = ?`

type TablesName struct {
	No        string `json:"no"`
	Tablename string `json:"tablename"`
}

type ForeignClient struct {
	ClientID string `json:"client_id"`
	UserID   int    `json:"user_id"`
	PlanID   int    `json:"plan_id"`
}

type ColumnDescribtion struct {
	Field   string `json:"Field"`
	Type    string `json:"Type"`
	Null    string `json:"Null"`
	Key     string `json:"Key"`
	Default string `json:"Default"`
	Extra   string `json:"Extra"`
}

type IRepository interface {
	GetTablesNameDB1() ([]TablesName, error)
	GetTablesNameDB2() ([]TablesName, error)
	GetDescColumnDB1(tabelName string) ([]ColumnDescribtion, error)
	GetDescColumnDB2(tabelName string) ([]ColumnDescribtion, error)
	ExecMigrationTable(tableSource string, tableTarget string, page int, limit int, email string) (bool, error)
	MigrateRelatedData(email string, offset, limit int) (bool, error)
	DeleteReletedData(email string) (bool, error)
	FetchDataTable(tablename string, page int, source string) ([]map[string]any, error)
}

type repository struct {
	DB1 *gorm.DB
	DB2 *gorm.DB
}

func NewTablesNameRepository(db1 *gorm.DB, db2 *gorm.DB) *repository {
	return &repository{db1, db2}
}

func (r *repository) GetTablesNameDB1() ([]TablesName, error) {
	var tablesname []TablesName
	var DB1_NAME = os.Getenv("DB1_NAME")
	err := r.DB1.Raw(TablesNameMySQL, DB1_NAME).Scan(&tablesname).Error
	if err != nil {
		return tablesname, err
	}

	return tablesname, nil
}

func (r *repository) GetTablesNameDB2() ([]TablesName, error) {
	var tablesname []TablesName
	var DB2_NAME = os.Getenv("DB2_NAME")
	err := r.DB2.Raw(TablesNameMySQL, DB2_NAME).Scan(&tablesname).Error
	if err != nil {
		return tablesname, err
	}

	return tablesname, nil
}

func (r *repository) GetDescColumnDB1(tabelName string) ([]ColumnDescribtion, error) {
	var columnDesc []ColumnDescribtion
	sql := fmt.Sprintf("SHOW COLUMNS FROM %s", tabelName)
	err := r.DB1.Raw(sql).Scan(&columnDesc).Error
	if err != nil {
		return columnDesc, err
	}

	return columnDesc, err
}

func (r *repository) GetDescColumnDB2(tabelName string) ([]ColumnDescribtion, error) {
	var columnDesc []ColumnDescribtion
	sql := fmt.Sprintf("SHOW COLUMNS FROM %s", tabelName)
	err := r.DB2.Raw(sql).Scan(&columnDesc).Error
	if err != nil {
		return columnDesc, err
	}

	return columnDesc, err
}

func (r *repository) ExecMigrationTable(tableSource string, tableTarget string, page int, limit int, email string) (bool, error) {
	var dataTableDB1 []map[string]any

	// Jumlah data per halaman
	offset := (page - 1) * limit

	// Start transaction for DB2 (target database)
	tx := r.DB2.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// Defer rollback in case of error - will be ignored if committed successfully
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := r.DB1.Table(tableSource).Limit(limit).Offset(offset).Find(&dataTableDB1).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// insert to target
	for _, v := range dataTableDB1 {
		log.Println(v)

		switch tableSource {
		case "client_tbl_kyc":
			delete(v, "UserID")
		case "web_tbl_user":
			t := time.Time{}
			if v["UserLockDate"] == t {
				v["UserLockDate"] = nil
			}
		case "client_tbl_detail_bankaccount":
			utils.ClientBankaccountRule(r.DB1, "client_tbl_kyc", v)
		case "client_tbl_plan_portfolio":
			utils.PlanPortfolioRule(v)
		case "web_tbl_transaction":
			utils.TransactionRule(v)
		}

		if err := r.DB2.Table(tableTarget).Create(&v).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	return true, nil
}

func (r *repository) MigrateRelatedData(email string, offset int, limit int) (bool, error) {
	// Jumlah data per halaman
	// offset := (page - 1) * limit

	// Start transaction for DB2 (target database)
	tx := r.DB2.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// Defer rollback in case of error - will be ignored if committed successfully
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	res, err := utils.MigrateRelatedData(email, limit, offset, tx, r.DB1, r.DB2)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (r *repository) DeleteReletedData(email string) (bool, error) {
	// Start transaction for DB2 (target database)
	tx := r.DB2.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// Defer rollback in case of error - will be ignored if committed successfully
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	res, err := utils.DeleteReletedDataRule(tx, r.DB2, email)

	return res, err
}

func (r *repository) FetchDataTable(tablename string, page int, source string) ([]map[string]any, error) {
	var results []map[string]any
	limit := 10                  // Jumlah data per halaman
	offset := (page - 1) * limit // Menghitung offset berdasarkan halaman

	// Query dengan paginasis
	switch source {
	case "source":
		if err := r.DB1.Table(tablename).Limit(limit).Offset(offset).Find(&results).Error; err != nil {
			return nil, err
		}
	case "target":
		if err := r.DB2.Table(tablename).Limit(limit).Offset(offset).Find(&results).Error; err != nil {
			return nil, err
		}
	}

	return results, nil
}
