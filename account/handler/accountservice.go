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

package handler

import (
	
	"fmt"

	"strconv"

	"math/rand"
	
	"context"

	"errors"

	"encoding/json"

	"go-micro.dev/v4/logger"

	"go-micro.dev/v4/broker"

	pb "accountservice/proto"

	db "accountservice/database"
)

const (
	topic = "go.micro.topic.sendotp"
	errorText = "Unable to create user!";
	sessionError = "Invalid Session!";
	userNotFound = "User not found!";
	notAdminError = "User is not admin!";

	
	UPDATED_MSG = "Updated!";
)


type SMSRequest struct {
	Phone string `json:"phone"`
	Body string `json:"body"`
}

type AccountService struct {
	SessionService pb.SessionService
	LogisticsService pb.LogisticsService
}

func (service *AccountService) CreateUser (ctx context.Context, in *pb.CreateUserRequest, out *pb.CreateUserResponse) error {
	var otp string;
	//Deprecated Methods
	//TODO Remove in new version
	if user, _ := db.GetUserByPhone(in.Phone, int64(0)); user == nil {
		id, err := db.InsertUser(in.Phone, otp);
		if err != nil {
			logger.Info("User Not Inserted!")
			out.UserId = 0;
			out.Error = errorText;			
			return err;
		}else{
			out.UserId = int64(id);
			out.Error = "";
		}
	}else{
		out.UserId = user.Id;
		out.Error = "";
	}
	if(in.SendOtp == true){
		otp = fmt.Sprint(generateOTP());
		body, _ := json.Marshal(&SMSRequest {
			Phone: in.Phone,
			Body: "Your Verification code is " + otp,
		});
		_ = db.UpdateUserOTP(out.UserId, otp);
		if err := broker.Publish(topic, &broker.Message{Body: body}); err != nil {
			logger.Fatal(err);
		}
	}
	return nil;
}

func (service *AccountService) VerifyOTP (ctx context.Context, in *pb.VerifyOTPRequest, out *pb.VerifyOTPResponse) error {
	user, err := db.GetUserById(in.UserId);
	if err != nil {
		return err
	}
	otp, err := strconv.Atoi(in.Otp);
	if err != nil {
		return err;
	}
	if(user.OTP == int32(otp)){
		logger.Info("Verified!");
	}else{
		logger.Info("Not Verified!");
		return errors.New("Invalid OTP");
	}
	response, err := service.SessionService.CreateSession(ctx, &pb.CreateSessionRequest {
		UserId:in.UserId,
		ClientData:in.Client,
	});
	if err != nil {
		return err;
	}
	out.Token = response.SessionKey
	return nil;
}

func (service *AccountService) UpdateUser (ctx context.Context, in *pb.UpdateUserRequest, out *pb.CommonResponse) error {
	sessionData, err := service.SessionService.GetSession(ctx, &pb.GetSessionRequest {
		SessionKey: in.Token,
	});
	if err != nil {
		return errors.New(sessionError);
	}
	err = db.UpdateUserData(sessionData.UserId, in.Name);
	if err != nil {
		return err;
	}
	out.Status = 200;
	out.Message = "User Updated!";
	return nil;
}

func (service *AccountService) GetUserById (ctx context.Context, in *pb.GetUserRequest, out *pb.GetUserReponse) error {
	user, err := db.GetUserById(in.UserId);
	if err != nil {
		return err;
	}
	out.User = makeProtoUser(user);
	return nil;
}

func (service *AccountService) SearchUsers (ctx context.Context, in *pb.SearchUserRequest, out *pb.GetUsersReponse) error {
	mUsers, err := db.FilterUsers(in);
	if err != nil {
		return err
	}
	var users []*pb.User;
	for _, mUser := range *mUsers {
	   user := makeProtoUser(mUser);
	   response, err := service.LogisticsService.DetermineWalletBalance(ctx, &pb.DetermineWalletBalanceRequest {
	   	CurrentAmount: user.Wallet,
	   	UserId: user.Id,
	   });
	   if err == nil {
	   	user.Balance = response.Balance;
	   }else{
	   	logger.Info(err);
	   	user.Balance = 0;
	   }
	   users = append(users, user);
	}
	out.Users = users;
	return nil;	
}

func (service *AccountService) AdminLogin (ctx context.Context, in *pb.AdminLoginRequest, out *pb.AdminLoginResponse) error {
	//Deprecated Methods
	//TODO Remove in new version
	if user, _ := db.GetUserByPhone(in.Phone, int64(0)); user == nil {
		return errors.New(userNotFound);
	}else{
		if(user.Admin == true){
			out.UserId = user.Id;
		}else{
			return errors.New(notAdminError)
		}
	}
	return nil;
}

func (service *AccountService) AdminVerify (ctx context.Context, in *pb.AdminVerifyRequest, out *pb.AdminVerifyResponse) error {
	user, _ := db.GetUserById(in.UserId);
	if user == nil {
		return errors.New(userNotFound);
	}
	if in.Passcode == user.Passcode {
		if user.Admin == false {
			return errors.New(notAdminError);
		}
		response, err := service.SessionService.CreateSession(ctx, &pb.CreateSessionRequest {
			UserId:in.UserId,
			ClientData:in.Client,
		});
		if err != nil {
			logger.Info(err);
			return errors.New(sessionError);
		}
		out.Token = response.SessionKey;
		return nil;
	}else{
		return errors.New("Invalid Passcode");
	}	
}

func (service *AccountService) AdminCreateUser (ctx context.Context, in *pb.AdminCreateUserRequest, out *pb.CreateUserResponse) error {	
	user, err := db.GetUserByPhone(in.Phone, in.CompanyId);
	if user != nil {
		return errors.New("User Already Created!");
	}
	id, err := db.InsertUserData(&db.User {
		Name  :in.Name,
		Phone :in.Phone,
		CompanyId :in.CompanyId,
	});
	if err != nil {
		return err;
	}
	out.UserId = id;
	out.Error = "";
	return nil;
}
func (service *AccountService) AdminUpdateUser (ctx context.Context, in *pb.AdminUpdateUserRequest, out *pb.CommonResponse) error {
	err := db.UpdateUserData(in.UserId, in.Name);
	if err != nil {
		return err;
	}
	out.Status = 200;
	out.Message = "User Updated!";
	return nil;
}


func generateOTP () int {
	return rand.Intn(99999-10000) + 10000
}

func makeProtoUser(user *db.User) *pb.User {
	return &pb.User {	   	
		Id  		: user.Id,		
		Name 		: user.Name,    
		Phone 		: user.Phone,     
		Otp 		: user.OTP,
		VerifyCount : int32(user.VerifyCount),
		LastVerified: user.LastVerified,
		CreatedAt 	: user.CreatedAt,
		Wallet 		: user.Wallet,
	}
}