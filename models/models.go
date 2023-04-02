package models

import (
	"github.com/api/database"
)

// type RelationalObject interface {
//   Tab  
// }

type Inserter interface {
	InsertQuery() string
	Param() []any
}

type Selector interface {
	SelectQuery() string
	Param() []any
	NumberOfFields() int
}

type Updater interface {
	UpdateQuery() string
	Param() []any
}

type Deleter interface {
	DeleteQuery() string
	Param() []any
}

func Insert(i Inserter) error {
	db := database.GetConnection()

	_, err := db.Query(i.InsertQuery(), i.Param()...)
	if err != nil {
		return err
	}

	return nil
}

func Update(u Updater) error {
	db := database.GetConnection()

	_, err := db.Query(u.UpdateQuery(), u.Param()...)
	if err != nil {
		return err
	}

	return nil
}

func Delete(d Deleter) error {
	db := database.GetConnection()

	_, err := db.Query(d.DeleteQuery(), d.Param()...)
	if err != nil {
		return err
	}

	return nil
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
