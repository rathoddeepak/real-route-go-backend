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

package handler

import (
		
	"context"

	"errors"

	"go-micro.dev/v4/logger"

	pb "accountservice/proto"

	db "accountservice/database"
)

func (service *AccountService) CreateCompany (ctx context.Context, in *pb.CreateCompanyRequest, out *pb.CreateCompanyResponse) error {
	if email, _ := db.GetCompanyByEmail(in.Email); email == nil {
		id, err := db.InsertCompany(in);
		if err != nil {
			logger.Info(err);
			return err;
		}else{
			out.CompanyId = id;
		}
	}else{
	    return errors.New("Email Already Used!");
	}
	return nil;
}

func (service *AccountService) LoginCompany (ctx context.Context, in *pb.LoginCompanyRequest, out *pb.LoginCompanyResponse) error {
	company_id, err := db.VerifyCompanyLogin(in.Email, in.Password);
	if err != nil {
		return err
	}
	if company_id == 0 {
		return errors.New("Incorrect Email/Password!")
	}
	company, err := db.GetCompanyById(company_id);

	out.CompanyId = company_id; 
	out.Name = company.Name; 
	out.Contact = company.Contact; 
	return nil;
}

func (service *AccountService) VerifyCompanyOTP (ctx context.Context, in *pb.VerifyCompanyOTPRequest, out *pb.VerifyCompanyOTPResponse) error {
	verified := db.VerifyCompanyOTP(in);
	var status int32;
	status = 200;
	if verified == false {
		status = 400
	}
	out.Status = status; 
	return nil;
}

func (service *AccountService) UpdateCompany (ctx context.Context, in *pb.UpdateCompanyRequest, out *pb.CommonResponse) error {
	err := db.UpdateCompany(in);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = UPDATED_MSG;
	return nil;
}

func (service *AccountService) UpdateCompanySetting (ctx context.Context, in *pb.UpdateCompanySettingRequest, out *pb.CommonResponse) error {
	err := db.SetCompanySettings(in);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = UPDATED_MSG;
	return nil;
}



func (service *AccountService) CompanyPasswordReset (ctx context.Context, in *pb.CompanyPasswordResetRequest, out *pb.CommonResponse) error {
	err := db.CompanyPasswordReset(in.Email);
	if err != nil {
		return err;
	}
	out.Status = 200;
	out.Message = UPDATED_MSG;
	return nil;
}

func (service *AccountService) GetCompany (ctx context.Context, in *pb.GetCompanyRequest, out *pb.GetCompanyResponse) error {
	company, err := db.GetCompanyById(in.CompanyId);
	if err != nil {
		return err
	}
	out.Company = company;
	return nil;
}

func (service *AccountService) GetCompanies (ctx context.Context, in *pb.GetCompanyRequest, out *pb.GetCompaniesResponse) error {
	companies, err := db.FilterCompanies(in);
	if err != nil {
		return err
	}
	out.Companies = *companies;
	return nil;
}

