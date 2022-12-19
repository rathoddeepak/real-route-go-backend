/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 11 Sep 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package database

import (
	"fmt"
	
	_ "github.com/lib/pq"
	"go-micro.dev/v4/logger"
	"database/sql"
	"encoding/json"
	"time"
	"errors"
	"strings"

	sqlBuilder "github.com/doug-martin/goqu/v9"
	pb "accountservice/proto"
)

const defaultCompnaySettings = "{\"map_setting\": {\"vehicle\": true,\"vehicle_type\": 5},\"notifications\": {\"on_start_task\": true,\"on_cancel_task\": true,\"on_complete_task\": true}}";

func InsertCompany (company *pb.CreateCompanyRequest) (int64, error){
	timeUnix := time.Now().Unix();
	otp := "082122"; 
	query, _, err := sqlBuilder.Insert(COMPANY_TABLE).Rows(
		sqlBuilder.Record {
			"email": company.Email,
			"password" : company.Password,
			"contact" : company.Contact,
			"name" : company.Name,
			"role" : company.Role,
			"created" : timeUnix,
			"otp" : otp,
			"settings" : defaultCompnaySettings,
		},
	).Returning("id").ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	var lastInsertId int64;
	err = db.QueryRow(query).Scan(&lastInsertId);
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil;
}

func VerifyCompanyOTP (in *pb.VerifyCompanyOTPRequest) bool {
	var count int;
	dbSQL := sqlBuilder.From(COMPANY_TABLE).Select(sqlBuilder.COUNT("*")).Where(
		sqlBuilder.Ex {"email": in.Email},
		sqlBuilder.Ex {"otp": in.Otp},
	);
	query, _, _ := dbSQL.ToSQL();	
	err := db.QueryRow(query).Scan(&count);
	if err != nil {
		return false;
	}
	if count == 0 {
		return false;
	} else {
		return true;
	}
}

func UpdateCompany(in *pb.UpdateCompanyRequest) error {
	query, _, err := sqlBuilder.Update(COMPANY_TABLE).Set(
		sqlBuilder.Record {
			"name"	  : in.Name,
			"contact" : in.Contact,
			"role"	  : in.Role,
		},
	).Where(
		sqlBuilder.Ex {"id": in.CompanyId},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	rows, err := db.Query(query);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func CompanyPasswordReset(email string) error {
	//Some Logic For Password Reset
	return nil;
}

func GetCompanyById (company_id int64) (*pb.Company, error){
	query, _, _ := sqlBuilder.From(COMPANY_TABLE).Select(GetCompanyColumns()...).Where (
		sqlBuilder.Ex{"id": company_id},
	).Limit(1).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	present := rows.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := rows.Err(); err != nil {
		return nil, err;
	}

	company, err := StructCompany(rows);
	if err != nil {
		return nil, err
	}
	return company, nil
}

func SetCompanySettings (in *pb.UpdateCompanySettingRequest) (error){
	setingBytes, err := json.Marshal(&in.Setting);
    if err != nil {
    	return err;
    }
	query, _, err := sqlBuilder.Update(COMPANY_TABLE).Set(
		sqlBuilder.Record {
			"settings"	  : string(setingBytes),
		},
	).Where(
		sqlBuilder.Ex {"id": in.CompanyId},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	rows, err := db.Query(query);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func VerifyCompanyLogin (email string, password string) (int64, error){
	query, _, _ := sqlBuilder.From(COMPANY_TABLE).Select("id").Where (
		sqlBuilder.Ex{"email": email},
		sqlBuilder.Ex{"password": password},
	).Limit(1).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return 0, err;
	}
	defer rows.Close();
	present := rows.Next();
	if present == false {
		return 0, errors.New(errText);
	}
	if err := rows.Err(); err != nil {
		return 0, err;
	}
	var company_id int64;
	err = rows.Scan(&company_id);
	if err != nil {
		return 0, err
	}
	return company_id, nil
}

func GetCompanyByEmail (email string) (*pb.Company, error){
	query, _, _ := sqlBuilder.From(COMPANY_TABLE).Select(GetCompanyColumns()...).Where (
		sqlBuilder.Ex{"email": email},
	).Limit(1).ToSQL();

	rows, err := db.Query(query);
	defer rows.Close();
	if err != nil {
		logger.Fatal(err)
		return nil, err;
	}
	present := rows.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := rows.Err(); err != nil {		
		return nil, err;
	}	

	company, err := StructCompany(rows);
	if err != nil {
		logger.Fatal(err);
		return nil, err
	}
	return company, nil
}

func FilterCompanies(filterMap *pb.GetCompanyRequest) (*[]*pb.Company, error) {
	var companies []*pb.Company;
	dbSQL := sqlBuilder.From(COMPANY_TABLE).Select(GetCompanyColumns()...);

	if len(filterMap.Name) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"name": "%"+filterMap.Name+"%"})
	}

	if len(filterMap.Contact) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"contact": "%"+filterMap.Contact+"%"})
	}

	if len(filterMap.Email) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"email": "%"+filterMap.Email+"%"})
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

	if len(filterMap.Contact) > 0 {
		query = strings.Replace(query, "contact\" =", "contact\" ILIKE" , 1);
	}

	if len(filterMap.Email) > 0 {
		query = strings.Replace(query, "email\" =", "email\" ILIKE" , 1);
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
		user, err := StructCompany(rows);
		if err != nil {
			return nil, err
		}
		companies = append(companies, user);
	}
	return &companies, nil;
}

func StructCompany (row *sql.Rows) (*pb.Company, error) {
	var companySettingJSONB []byte;
	var c pb.Company;
	err := row.Scan(
		&c.Id,&c.Name,&c.Email,&c.Contact,&c.Role,&c.Created,&companySettingJSONB,
	);
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(companySettingJSONB, &c.Settings);  
    if err != nil {
        fmt.Println(err);
    }
	return &c, nil;
}