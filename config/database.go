package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Fungsi untuk menginisialisasi koneksi ke dua database PostgreSQL
// func InitDBs() (*gorm.DB, *gorm.DB, error) {

// 	// Koneksi ke PostgreSQL Database 1
// 	dsnPostgres1 := fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
// 		os.Getenv("DB1_HOST"), os.Getenv("DB1_USER"), os.Getenv("DB1_PASSWORD"), os.Getenv("DB1_NAME"), os.Getenv("DB1_PORT"),
// 	)
// 	dbPostgres1, err := gorm.Open(postgres.Open(dsnPostgres1), &gorm.Config{})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Koneksi ke PostgreSQL Database 2
// 	dsnPostgres2 := fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
// 		os.Getenv("DB2_HOST"), os.Getenv("DB2_USER"), os.Getenv("DB2_PASSWORD"), os.Getenv("DB2_NAME"), os.Getenv("DB2_PORT"),
// 	)
// 	dbPostgres2, err := gorm.Open(postgres.Open(dsnPostgres2), &gorm.Config{})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return dbPostgres1, dbPostgres2, nil
// }

func InitDBs() (*gorm.DB, *gorm.DB, error) {

    // Koneksi ke MySQL Database 1
    dsnMySQL1 := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB1_USER"), os.Getenv("DB1_PASSWORD"), os.Getenv("DB1_HOST"), os.Getenv("DB1_PORT"), os.Getenv("DB1_NAME"),
    )
    dbMySQL1, err := gorm.Open(mysql.Open(dsnMySQL1), &gorm.Config{})
    if err != nil {
        return nil, nil, err
    }

    // Koneksi ke MySQL Database 2
    dsnMySQL2 := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB2_USER"), os.Getenv("DB2_PASSWORD"), os.Getenv("DB2_HOST"), os.Getenv("DB2_PORT"), os.Getenv("DB2_NAME"),
    )
    dbMySQL2, err := gorm.Open(mysql.Open(dsnMySQL2), &gorm.Config{})
    if err != nil {
        return nil, nil, err
    }

    return dbMySQL1, dbMySQL2, nil
}
