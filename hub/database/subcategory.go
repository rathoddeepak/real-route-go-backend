package database;

import (
	"fmt"
	"errors"
	_ "github.com/lib/pq"
	"database/sql"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type SubCategory struct {
	Id 	   			int64 		`json:id`
	MainCatId 		int64 		`json:main_cat_id`
	Name  			string  	`json:name`
	Status  		int32   	`json:status`
	Image  			string   	`json:image`
	Blurhash  		string   	`json:blurhash`
}

func InsertSubCategory(subCategory *SubCategory) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(SUBCAT_TABLE).Rows(
		sqlBuilder.Record {
			"main_cat_id" : subCategory.MainCatId,
			"name" : subCategory.Name,
			"status" : 0,
			"image" : "",
			"blurhash" : "",
		},
	).ToSQL();
	var subCategoryId int64;
	err := db.QueryRow(sql + " RETURNING id").Scan(&subCategoryId);	
	if err != nil {
		return nil, err;
	}
	return &subCategoryId, nil
}

func UpdateSubCategory(subCategory *SubCategory) (error){	
	sql, _, _ := sqlBuilder.Update(SUBCAT_TABLE).Set(
		sqlBuilder.Record {			
			"name" : subCategory.Name,
		},
	).Where(sqlBuilder.Ex {"id": subCategory.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateSubCategoryStatus(subCategory *SubCategory) (error){	
	sql, _, _ := sqlBuilder.Update(SUBCAT_TABLE).Set(
		sqlBuilder.Record {			
			"status" : subCategory.Status,
		},
	).Where(sqlBuilder.Ex {"id": subCategory.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetSubCategoriesOfCategory(mainCatId int64) (*[]*SubCategory, error) {
	var subCategories []*SubCategory;
	query := fmt.Sprintf(selectSubCatSQL, "main_cat_id", mainCatId);
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			continue
		}
		subCategory, err := StructSubCategory(rows);
		subCategories = append(subCategories, subCategory);
	}
	return &subCategories, nil;
}

func GetSubCategoryById(subCategory_id int64) (*SubCategory, error) {	
	query := fmt.Sprintf(selectSubCatSQL + " LIMIT 1", "id", subCategory_id);
	row, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer row.Close();	
	present := row.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := row.Err(); err != nil {
		return nil, err;
	}
	category, err := StructSubCategory(row);
	if err != nil {
		return nil, err;
	}
	return category, nil;
}

func UpdateSubCategoryImage(subCategory *SubCategory) (error){	
	// mSubCategory, err := GetSubCategoryById(subCategory.Id);
	// if err != nil {
	// 	return err
	// }
	// if(len(mSubCategory.Image) != 0){
	// 	os.Remove(mSubCategory.Image);
	// }
	sql, _, _ := sqlBuilder.Update(SUBCAT_TABLE).Set(
		sqlBuilder.Record {			
			"image": subCategory.Image,
			"blurhash": subCategory.Blurhash,
		},
	).Where(sqlBuilder.Ex {"id": subCategory.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func StructSubCategory (row *sql.Rows) (*SubCategory, error) {
	var subCat SubCategory;
	err := row.Scan(
		&subCat.Id,
		&subCat.MainCatId,
		&subCat.Name,		
		&subCat.Status,
		&subCat.Image,
		&subCat.Blurhash,
	);
	if err != nil {
		return nil, err
	}
	return &subCat, nil;
}