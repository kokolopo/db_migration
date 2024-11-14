package service

import "db_migration/entity"

type IService interface {
	GetTablesNameDB1() ([]entity.TablesName, error)
	GetTablesNameDB2() ([]entity.TablesName, error)
	GetColumnsDescDB1(string) ([]entity.ColumnDescribtion, error)
	GetColumnsDescDB2(string) ([]entity.ColumnDescribtion, error)
	MigrateTable(string, string, int, int, string) (bool, error)
	MigrateReletedData(string, int, int) (bool, error)
	DeleteReletedData(string) (bool, error)
	GetDataTable(string, int, string) ([]map[string]any, error)
}

type service struct {
	repository entity.IRepository
}

func NewUserService(repository entity.IRepository) *service {
	return &service{repository}
}

func (s *service) GetTablesNameDB1() ([]entity.TablesName, error) {
	tablesname, err := s.repository.GetTablesNameDB1()
	if err != nil {
		return tablesname, err
	}

	return tablesname, nil
}

func (s *service) GetTablesNameDB2() ([]entity.TablesName, error) {
	tablesname, err := s.repository.GetTablesNameDB2()
	if err != nil {
		return tablesname, err
	}

	return tablesname, nil
}

func (s *service) GetColumnsDescDB1(columnName string) ([]entity.ColumnDescribtion, error) {
	data, err := s.repository.GetDescColumnDB1(columnName)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *service) GetColumnsDescDB2(columnName string) ([]entity.ColumnDescribtion, error) {
	data, err := s.repository.GetDescColumnDB1(columnName)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *service) MigrateTable(tableSource string, tableTarget string, page int, limit int, email string) (bool, error) {
	res, err := s.repository.ExecMigrationTable(tableSource, tableTarget, page, limit, email)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s *service) MigrateReletedData(email string, offset int, limit int) (bool, error) {
	res, err := s.repository.MigrateRelatedData(email, offset, limit)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (s *service) DeleteReletedData(email string) (bool, error) {
	res, err := s.repository.DeleteReletedData(email)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s *service) GetDataTable(tablename string, page int, source string) ([]map[string]any, error) {
	data, err := s.repository.FetchDataTable(tablename, page, source)
	if err != nil {
		return data, err
	}

	return data, nil
}
