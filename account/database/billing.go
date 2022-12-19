/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 17 August 2022
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

	sqlBuilder "github.com/doug-martin/goqu/v9"
	pb "accountservice/proto"
)

func InsertCompanyBilling (companyBilling *pb.CompanyBilling) (int64, error){
	query, _, err := sqlBuilder.Insert(COMPANY_BILLING_TABLE).Rows(
		sqlBuilder.Record {
			"rp_subscription_id":companyBilling.RazorpaySubscriptionId,
			"company_id":companyBilling.CompanyId,
			"rp_plan_id":companyBilling.RazorpayPlanId,
			"amount":companyBilling.Amount,
			"start_at":companyBilling.StartAt,
			"end_at":companyBilling.EndAt,
			"status":companyBilling.Status,
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

func GetCompanyBillings (company_id int64) (*[]*pb.CompanyBilling, error){
	query, _, _ := sqlBuilder.From(COMPANY_BILLING_TABLE).Select(GetCompanyBillingColumns()...).Where (
		sqlBuilder.Ex{"company_id": company_id},
	).Order(sqlBuilder.I("id").Desc()).ToSQL();		
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var txns []*pb.CompanyBilling;
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			continue
		}
		txn, err := StructCompanyBilling(rows);
		txns = append(txns, txn);
	}
	return &txns, nil;
}

func UpdateCompanyBillingStatus (company *pb.MarkCompanyBillingStatusRequest) (error) {
	query, _, err := sqlBuilder.Update(COMPANY_BILLING_TABLE).Set(sqlBuilder.Record {"status":company.Status}).Where(
		sqlBuilder.Ex {"rp_subscription_id": company.RazorpaySubscriptionId},
	).ToSQL();
	if err != nil {
		return err;
	}
	rows, err := db.Query(query);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func StructCompanyBilling (row *sql.Rows) (*pb.CompanyBilling, error) {
	var u pb.CompanyBilling;
	err := row.Scan(
		&u.Id,&u.RazorpaySubscriptionId,&u.RazorpayPlanId,&u.Amount,&u.StartAt,&u.EndAt,&u.Status,
	);	
	if err != nil {
		return nil, err
	}
	return &u, nil;
}