/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 17 July 2022
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
	"time"
	"errors"
	"strings"

	sqlBuilder "github.com/doug-martin/goqu/v9"
	pb "accountservice/proto"
)

type User struct {
	Id 			 int64 		`json:"id"`
	Name 	     string		`json:"name"`
	Email 	     string 	`json:"email"`
	Phone 	     string 	`json:"phone"`
	Gender 	     int16		`json:"gender"`
	OTP 	     int32 		`json:"otp"`
	Active 		 int16 		`json:"active"`
	VerifyCount  int16 		`json:"verifyCount"`	
	LastVerified int64 		`json:"lastVerified"`
	CreatedAt 	 int64 		`json:"createdAt"`
	Admin 	 	 bool 		`json:"admin"`
	Passcode 	 string 	`json:"passcode"`
	Wallet 	 	 float64 	`json:"wallet"`
	CompanyId 	 int64 		`json:"company_id"`
}

func InsertUser (phone string, otp string) (int, error){
	timeUnix := time.Now().Unix();
	query, _, err := sqlBuilder.Insert(USER_TABLE).Rows(
		sqlBuilder.Record {
			"phone":phone,
			"createdAt":timeUnix,
			"otp":otp,
			"verifyCount": 1,
			"lastVerified":timeUnix,
			"wallet":0,
		},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	var lastInsertId int;
	err = db.QueryRow(query + " RETURNING id").Scan(&lastInsertId);
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil;
}

func InsertUserData (user *User) (int64, error){
	timeUnix := time.Now().Unix();
	query, _, err := sqlBuilder.Insert(USER_TABLE).Rows(
		sqlBuilder.Record {
			"phone": user.Phone,
			"name" : user.Name,
			"active":1,
			"otp"  : 0,
			"createdAt":timeUnix,
			"verifyCount": 0,
			"lastVerified":timeUnix,
			"wallet":0,
			"company_id" : user.CompanyId,
		},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	var lastInsertId int64;
	err = db.QueryRow(query + " RETURNING id").Scan(&lastInsertId);
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil;
}

func UpdateUserData(user_id int64, name string) error {
	query, _, err := sqlBuilder.Update(USER_TABLE).Set(
		sqlBuilder.Record {
			"name":name,
		},
	).Where(
		sqlBuilder.Ex {"id": user_id},
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

func UpdateUserWallet(user_id int64, wallet float64) error {
	query, _, err := sqlBuilder.Update(USER_TABLE).Set(
		sqlBuilder.Record {
			"wallet":wallet,
		},
	).Where(
		sqlBuilder.Ex {"id": user_id},
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

func UpdateUserOTP(user_id int64, otp string) error {
	query, _, err := sqlBuilder.Update(USER_TABLE).Set(
		sqlBuilder.Record {
			"otp":otp,
		},
	).Where(
		sqlBuilder.Ex {"id": user_id},
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



func GetUserById (user_id int64) (*User, error){
	query, _, _ := sqlBuilder.From(USER_TABLE).Where (
		sqlBuilder.Ex{"id": user_id},
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

	user, err := StructUser(rows);
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByPhone (phone string, company_id int64) (*User, error){
	query, _, _ := sqlBuilder.From(USER_TABLE).Where (
		sqlBuilder.Ex{"phone": phone},
		sqlBuilder.Ex{"company_id": company_id},
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

	user, err := StructUser(rows);
	if err != nil {
		logger.Fatal(err);
		return nil, err
	}
	return user, nil
}

func FilterUsers(filterMap *pb.SearchUserRequest) (*[]*User, error) {
	var users []*User;
	dbSQL := sqlBuilder.From(USER_TABLE);

	if filterMap.CompanyId == 0 {
		return nil, errors.New("Company Id is required!");
	}else{
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"company_id": filterMap.CompanyId })
	}

	if filterMap.FilterName == true && len(filterMap.Name) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"name": "%"+filterMap.Name+"%"})
	}

	if filterMap.FilterPhone == true && len(filterMap.Phone) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"phone": "%"+filterMap.Phone+"%"})
	}

	if filterMap.WalletRange == true {
		dbSQL = dbSQL.Where(sqlBuilder.And(
			sqlBuilder.C("wallet").Gte(filterMap.MinPrice),
			sqlBuilder.C("wallet").Lte(filterMap.MaxPrice),
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
	
	if filterMap.FilterName == true && len(filterMap.Name) > 0 {
		query = strings.Replace(query, "name\" =", "name\" ILIKE" , 1);
	}

	if filterMap.FilterPhone == true && len(filterMap.Phone) > 0 {
		query = strings.Replace(query, "phone\" =", "phone\" ILIKE" , 1);
	}

	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	
	defer rows.Close();
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			fmt.Println(err);
			continue
		}
		user, err := StructUser(rows);
		if err != nil {
			return nil, err
		}
		users = append(users, user);
	}
	return &users, nil;
}

func StructUser (row *sql.Rows) (*User, error) {
	var u User;
	err := row.Scan(
		&u.Email,&u.Name,&u.OTP,&u.Id,&u.CreatedAt,&u.Phone,
		&u.Gender,&u.Active,&u.VerifyCount,&u.LastVerified,&u.Admin,
		&u.Passcode,&u.Wallet,&u.CompanyId,
	);	
	if err != nil {
		return nil, err
	}
	return &u, nil;
}