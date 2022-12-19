/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 14 Sep 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler

import (

	"context"

	pb "accountservice/proto"

	db "accountservice/database"
)

func (service *AccountService) CreateCompanyBilling (ctx context.Context, in *pb.CreateCompanyBillingRequest, out *pb.CreateCompanyBillingResponse) error {
	razorpay_subscription_id, err := db.CreateRazorPaySubscription(in);
	if err != nil {
		return err
	}
	//TODO: HardCoded Start and End Date Must Be changed later
	billing_id, err := db.InsertCompanyBilling(&pb.CompanyBilling {
		CompanyId : in.CompanyId,
		RazorpaySubscriptionId : *razorpay_subscription_id,
		RazorpayPlanId : in.RazorpayPlanId,
		Amount : in.RazorpayPlanAmount,
		StartAt : int64(0),
		EndAt : int64(0),
		Status : db.RAZOR_SUBSCRIPTION_CREATED,
	});
	if err != nil {
		return err
	}
	out.BillingId = billing_id;
	out.RazorpaySubscriptionId = *razorpay_subscription_id;
	return nil;
}

func (service *AccountService) ValidateCompanyBilling (ctx context.Context, in *pb.ValidateCompanyBillingRequest, out *pb.ValidateCompanyBillingResponse) error {
	billings, err := db.GetCompanyBillings(in.CompanyId);
	if err != nil {
		return err
	}	
	for _, billing := range *billings {
		if billing.Status == db.RAZOR_SUBSCRIPTION_CREATED {
			out.RazorpaySubscriptionId = billing.RazorpaySubscriptionId;
			out.Active = false;
			out.Status = billing.Status;
		}else if billing.Status == db.RAZOR_SUBSCRIPTION_AUTHENTICATED || billing.Status == db.RAZOR_SUBSCRIPTION_ACTIVE {
			out.RazorpaySubscriptionId = billing.RazorpaySubscriptionId;
			out.Active = true;
			out.Status = billing.Status;
			return nil
		}
		// else if billing.Status == db.RAZOR_SUBSCRIPTION_PENDING || billing.Status == db.RAZOR_SUBSCRIPTION_HALTED|| billing.Status == db.RAZOR_SUBSCRIPTION_CANCELLED || billing.Status == db.RAZOR_SUBSCRIPTION_COMPLETED || billing.Status == db.RAZOR_SUBSCRIPTION_EXPIRED {
		// 	out.Active = false;
		// }else if  {
		// 	out.Active = false;
		// }
	}
	out.Active = false;
	return nil;
}

func (service *AccountService) MarkCompanyBillingStatus (ctx context.Context, in *pb.MarkCompanyBillingStatusRequest, out *pb.CommonResponse) error {
	err := db.UpdateCompanyBillingStatus(in);
	if err != nil {
		return err
	}	
	out.Status = 200;
	out.Message = UPDATED_MSG;
	return nil;
}


