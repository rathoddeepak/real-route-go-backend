package database;

import (
	"fmt"
	"errors"
	_ "github.com/lib/pq"
	"database/sql"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Category struct {
	Id 	   			int64 		`json:id`
	HubId 			int64 		`json:hub_id`
	Name  			string  	`json:name`
	Status  		int32   	`json:status`
	Image  			string   	`json:image`
	Blurhash  		string   	`json:blurhash`
}

func InsertCategory(category *Category) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(CATEGORY_TABLE).Rows(
		sqlBuilder.Record {
			"hub_id" : category.HubId,
			"name" : category.Name,
			"status" : 0,
			"image"  : "",
			"blurhash" : "",
		},
	).ToSQL();
	var categoryId int64;
	err := db.QueryRow(sql + " RETURNING id").Scan(&categoryId);	
	if err != nil {
		return nil, err;
	}
	return &categoryId, nil
}

func UpdateCategory(category *Category) (error){	
	sql, _, _ := sqlBuilder.Update(CATEGORY_TABLE).Set(
		sqlBuilder.Record {			
			"name" : category.Name,
		},
	).Where(sqlBuilder.Ex {"id": category.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateCategoryStatus(category *Category) (error){	
	sql, _, _ := sqlBuilder.Update(CATEGORY_TABLE).Set(
		sqlBuilder.Record {			
			"status" : category.Status,
		},
	).Where(sqlBuilder.Ex {"id": category.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetCategoriesOfHub(hub_id int64) (*[]*Category, error) {
	var categories []*Category;
	query := fmt.Sprintf(selectCatSQL, "hub_id", hub_id);
	rows, err := db.Query(query);
	defer rows.Close();
	if err != nil {
		return nil, err;
	}
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			continue
		}
		category, err := StructCategory(rows);
		categories = append(categories, category);
	}
	return &categories, nil;
}

func GetCategoryById(category_id int64) (*Category, error) {	
	query := fmt.Sprintf(selectCatSQL + " LIMIT 1", "id", category_id);
	row, err := db.Query(query);
	defer row.Close();
	if err != nil {
		return nil, err;
	}	
	present := row.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := row.Err(); err != nil {
		return nil, err;
	}
	category, err := StructCategory(row);
	if err != nil {
		return nil, err;
	}
	return category, nil;
}

func UpdateCategoryImage(category *Category) (error){	
	// mCategory, err := GetCategoryById(category.Id);
	// if err != nil {
	// 	return err
	// }
	// if(len(mCategory.Image) != 0){
	// 	os.Remove(mCategory.Image);
	// }
	sql, _, _ := sqlBuilder.Update(CATEGORY_TABLE).Set(
		sqlBuilder.Record {			
			"image": category.Image,
			"blurhash": category.Blurhash,
		},
	).Where(sqlBuilder.Ex {"id": category.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func StructCategory (row *sql.Rows) (*Category, error) {
	var cat Category;
	err := row.Scan(
		&cat.Id,
		&cat.HubId,
		&cat.Name,		
		&cat.Status,
		&cat.Image,
		&cat.Blurhash,
	);
	if err != nil {
		return nil, err
	}
	return &cat, nil;
}