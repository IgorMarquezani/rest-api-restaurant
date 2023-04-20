package models

import (
	"github.com/api/database"
)

// type RelationalObject interface {
//   Tab
// }

type Inserter interface {
	InsertQuery() string
	InsertParam() []any
}

func Insert(inserter Inserter) error {
	db := database.GetConnection()

	_, err := db.Query(inserter.InsertQuery(), inserter.InsertParam()...)
	if err != nil {
		return err
	}

	return nil
}

type Updater interface {
	UpdateQuery() string
	UpdateParam() []any
}

func Update(updater Updater) error {
	db := database.GetConnection()

	_, err := db.Query(updater.UpdateQuery(), updater.UpdateParam()...)
	if err != nil {
		return err
	}

	return nil
}

type Deleter interface {
	DeleteQuery() string
	DeleteParam() []any
}

func Delete(deleter Deleter) error {
	db := database.GetConnection()

	_, err := db.Query(deleter.DeleteQuery(), deleter.DeleteParam()...)
	if err != nil {
		return err
	}

	return nil
}

type Selector interface {
	SelectQuery() string
	Param() []any
	NumberOfFields() int
}

// func Select[T RelationalObject](s Selector) ([]T, error) {
// 	tabs := make([]T, 0)
// 	db := database.GetConnection()

// 	rows, err := db.Query(s.SelectQuery(), s.Param()...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for i := 0; rows.Next(); i++ {
//     tabs = append(tabs, T{})

// 		err := rows.Scan(tabs[i].FieldsAdress()...)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return tabs, nil
// }
