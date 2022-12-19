package database;

import (		
	"fmt"
	"strings"
	"errors"
	"database/sql"	

	_ "github.com/lib/pq"

	pb "hubservice/proto"

	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Product struct {
	Id  		  int64 	`json:id`
	Name    	  string	`json:name`
	CategoryId    int64 	`json:category_id`
	SubCategoryId int64  	`json:sub_category_id`
	HubId   	  int64  	`json:hub_id`
	CityId   	  int64  	`json:city_id`
	Price         float64 	`json:price`
	HigherPrice   float64 	`json:higher_price`
	Qty           int64 	`json:qty`
	MaxLimit      int64 	`json:max_limit`
	Status        int32 	`json:status`
	Image         string  	`json:image`
	Blurhash      string 	`json:blurhash`
}

func InsertProduct(product *Product) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(PRODUCT_TABLE).Rows(
		sqlBuilder.Record {
			"name": 			product.Name,
			"category_id": 		product.CategoryId,
			"sub_category_id":  product.SubCategoryId,
			"city_id": 			product.CityId,
			"hub_id": 			product.HubId,
			"price":  			product.Price,
			"higher_price":  	product.HigherPrice,
			"qty":  			product.Qty,
			"max_limit":		product.MaxLimit,
			"status": 			0,
			"image" : 			"",
			"blurhash":			"",
		},
	).ToSQL();
	var productId int64;
	err := db.QueryRow(sql + " RETURNING id").Scan(&productId);	
	if err != nil {
		return nil, err;
	}
	return &productId, nil
}

func UpdateProduct(product *Product) (error){	
	sql, _, _ := sqlBuilder.Update(PRODUCT_TABLE).Set(
		sqlBuilder.Record {
			"name":product.Name,
			"price":product.Price,
			"higher_price":product.HigherPrice,
			"qty":product.Qty,
			"max_limit": product.MaxLimit,
		},
	).Where(sqlBuilder.Ex {"id": product.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateProductStatus(product *Product) (error){	
	sql, _, _ := sqlBuilder.Update(PRODUCT_TABLE).Set(
		sqlBuilder.Record {			
			"status" : product.Status,
		},
	).Where(sqlBuilder.Ex {"id": product.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateProductQty(product *Product) (error){	
	sql, _, _ := sqlBuilder.Update(PRODUCT_TABLE).Set(
		sqlBuilder.Record {			
			"qty": product.Qty,
			"max_limit": product.MaxLimit,
		},
	).Where(sqlBuilder.Ex {"id": product.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateProductImage(product *Product) (error){	
	// product, err := GetProductById(info.ProductId);
	// if err != nil {
	// 	return err
	// }
	// if(len(product.Image) != 0){
	// 	os.Remove(product.Image);
	// }
	sql, _, _ := sqlBuilder.Update(PRODUCT_TABLE).Set(
		sqlBuilder.Record {			
			"image": product.Image,
			"blurhash": product.Blurhash,
		},
	).Where(sqlBuilder.Ex {"id": product.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetProductOfCategory(category_id int64) (*[]*Product, error) {
	var products []*Product;
	query := fmt.Sprintf(selectProductSQL, "category_id", category_id);
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
		product, err := StructProduct(rows);
		products = append(products, product);
	}
	return &products, nil;
}

func GetProductOfSubCategory(sub_category_id int64) (*[]*Product, error) {
	var products []*Product;
	query := fmt.Sprintf(selectProductSQL, "sub_category_id", sub_category_id);
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
		product, err := StructProduct(rows);
		products = append(products, product);
	}
	return &products, nil;
}

func FilterProducts(filterMap *pb.ProductSearchRequest) (*[]*Product, error) {
	var products []*Product;
	dbSQL := sqlBuilder.From(PRODUCT_TABLE).Select("id","name","category_id","sub_category_id","hub_id","price","higher_price","qty","max_limit","status","image","blurhash","city_id");

	if filterMap.CityId == 0 {
		return nil, errors.New("City Id is required!");
	}else{
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"city_id": filterMap.CityId})
	}

	if len(filterMap.Name) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"name": "%"+filterMap.Name+"%"})
	}

	if filterMap.HubId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"hub_id": filterMap.HubId})
	}

	if filterMap.MainCatId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"category_id": filterMap.MainCatId})
	}

	if filterMap.SubCatId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"sub_category_id": filterMap.SubCatId})
	}
	
	if filterMap.ProductId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"id": filterMap.ProductId})
	}

	if filterMap.Status != 9 {//TODO: 9 is temporary change later
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"status": filterMap.Status})
	}

	if filterMap.PriceRange == true {//TODO: 9 is temporary change later
		dbSQL = dbSQL.Where(sqlBuilder.And(
			sqlBuilder.C("price").Gte(filterMap.MinPrice),
			sqlBuilder.C("price").Lte(filterMap.MaxPrice),
		))
	}
	
	if filterMap.Limit == 0 {
		dbSQL = dbSQL.Limit(10);
	}else{
		dbSQL = dbSQL.Limit(uint(filterMap.Limit));
	}

	if filterMap.Offset == 0 {
		dbSQL = dbSQL.Offset(0);
	}else{
		dbSQL = dbSQL.Offset(uint(filterMap.Offset));
	}

	query, _, _ := dbSQL.Order(sqlBuilder.I("id").Desc()).ToSQL();
	
	if len(filterMap.Name) > 0 {
		query = strings.Replace(query, "name\" =", "name\" ILIKE" , 1);
	}

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
		product, err := StructProduct(rows);
		products = append(products, product);
	}
	return &products, nil;
}

func GetProductById(product_id int64) (*Product, error) {	
	query := fmt.Sprintf(selectProductSQL + " LIMIT 1", "id", product_id);
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
	product, err := StructProduct(row);
	if err != nil {
		return nil, err;
	}
	return product, nil;
}

func StructProduct (row *sql.Rows) (*Product, error) {
	var product Product;
	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.CategoryId,
		&product.SubCategoryId,
		&product.HubId,
		&product.Price,
		&product.HigherPrice,
		&product.Qty,
		&product.MaxLimit,
		&product.Status,
		&product.Image,
		&product.Blurhash,
		&product.CityId,
	);
	if err != nil {
		return nil, err
	}
	return &product, nil;
}