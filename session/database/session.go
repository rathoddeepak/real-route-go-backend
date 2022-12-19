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
	
	_ "github.com/lib/pq"

	"go-micro.dev/v4/logger"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/google/uuid"

	sqlBuilder "github.com/doug-martin/goqu/v9"
)

const TABLE string = "session";
const errText string = "No Records!";

var psqlInfo string = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",

    "localhost",5432,
    "postgres",
    "Rathod@123",
    "rr",
)

type Session struct {
	Id int64 `json:"id"`
	UserId int64 `json:"user_id"`
	SessionKey string `json:"SessionKey"` 
	ClientData string `json:"clientData"`
	CreatedAt int64 `json:"createdAt"`
}

var db *sql.DB;

func (session *Session) toJSON() (string, error) {
	res, err := json.Marshal(session);
	if err != nil {
		return "", err
	}
	return string(res), nil;
}

func init () {
	logger.Info("Opening Database Connection!");
	mdb, err := sql.Open("postgres", psqlInfo);
	if err != nil {
		logger.Fatal(err);
	}
	db = mdb;
	err = mdb.Ping();
	if err != nil {
		logger.Fatal(err);
	}
}


func InsertSession (user_id int64, clientData string) (int, string, error){
	uid, err := uuid.NewRandom();
	if err != nil {
		return 0, "", err
	}
	sessionKey := fmt.Sprint(user_id) + uid.String() + fmt.Sprint(time.Now().Unix());
	query, _, _ := sqlBuilder.Insert(TABLE).Rows(
		sqlBuilder.Record { 
			"user_id":user_id,
			"clientData":clientData,
			"sessionKey":sessionKey,
			"createdAt": time.Now().Unix(),
		},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
		return 0, "", err
	}
	var lastInsertId int;
	err = db.QueryRow(query + " RETURNING id").Scan(&lastInsertId);
	if err != nil {
		logger.Fatal(err);
		return 0, "", err
	}
	return lastInsertId, sessionKey, nil;
}

func GetSessionById (session_id int64) (*Session, error){
	query, _, _ := sqlBuilder.From(TABLE).Where (
		sqlBuilder.Ex{"id": session_id},
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
		return nil, err
	}
	session, err := StructSession(rows);
	if err != nil {
		return nil, err
	}
	return session, nil
}

func GetSession (sessionKey string) (*Session, error){
	query, _, _ := sqlBuilder.From(TABLE).Where (
		sqlBuilder.Ex{"sessionKey": sessionKey},
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
		return nil, err
	}
	session, err := StructSession(rows);
	if err != nil {
		return nil, err
	}
	return session, nil
}

func GetUserSession (user_id int64) (*[]*Session, error){
	query, _, _ := sqlBuilder.From(TABLE).Where (
		sqlBuilder.Ex{"user_id": user_id},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var sessions []*Session;
	for rows.Next() {
		session, err := StructSession(rows);
		if err != nil {
			continue
		}
		sessions = append(sessions, session);
	}
	return &sessions, nil
}

func StructSession (row *sql.Rows) (*Session, error) {
	var u Session;
	err := row.Scan(
		&u.Id,&u.ClientData,&u.UserId,&u.SessionKey,&u.CreatedAt,
	);
	if err != nil {
		return nil, err
	}
	return &u, nil;
}