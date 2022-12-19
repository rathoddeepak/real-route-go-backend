/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 14 Sep 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Account Service   <---
--------------------------------
*/

package database

import (
	"errors"

	razorpay "github.com/razorpay/razorpay-go"
	pb "accountservice/proto"
)

var client *razorpay.Client;
const RAZOR_SUBSCRIPTION_CREATED = 0;
const RAZOR_SUBSCRIPTION_AUTHENTICATED = 2;
const RAZOR_SUBSCRIPTION_ACTIVE = 1;
const RAZOR_SUBSCRIPTION_PENDING = 3;
const RAZOR_SUBSCRIPTION_HALTED = 4;
const RAZOR_SUBSCRIPTION_CANCELLED = 5;
const RAZOR_SUBSCRIPTION_COMPLETED = 6;
const RAZOR_SUBSCRIPTION_EXPIRED = 7;
const RAZOR_SUBSCRIPTION_UNKNOWN = 21;

const RAZORPAY_SECRET = "TDrVQI0dUjTRWzZ0j9rYLbm6";
const RAZORPAY_API_KEY = "rzp_live_TymCu5hNtINRGQ";

func init () {
	client = razorpay.NewClient(RAZORPAY_API_KEY, RAZORPAY_SECRET)
}

//Subscriptions and Plans
func GetRazorpayPlans (in *pb.GetRazorPayPlansRequest) (*[]*pb.RazorPayPlan, error){
	options:= map[string]interface{}{
	  "count": 1,
	  "skip": 0,
	}
	body, err := client.Plan.All(options, nil);
	if err != nil {
		return nil, err;
	}
	var plans []*pb.RazorPayPlan;


	count, ok := body["count"];
	if ok == false {
		return nil, errors.New(RESPONSE_PARSE_ERROR)
	}
	if count == 0 {
		return &plans, err;
	}

	response_plans, ok := body["items"].([]interface{});
	if ok == false {
		return nil, errors.New(RESPONSE_PARSE_ERROR)
	}

	for _, planInterface := range response_plans {
		plan, ok := planInterface.(map[string]interface {})
		if ok == false {
			continue
		}

		planItem, ok := plan["item"].(map[string]interface {})
		if ok == false {
			continue
		}
		description, ok := planItem["description"].(string);
		if ok == false {
			description = "";
		}
		plans = append(plans, &pb.RazorPayPlan {
			Id: plan["id"].(string),
			Name: planItem["name"].(string),
			Description: description,
			Amount: planItem["amount"].(float64),
			UnitAmount: planItem["unit_amount"].(float64),
		});
	}
	return &plans, nil;
}

func CreateRazorPaySubscription (in *pb.CreateCompanyBillingRequest) (*string, error){
	data := map[string]interface{}{
	  "plan_id":in.RazorpayPlanId,
	  "total_count":in.NumberOfBillingMonth,
	  "quantity": 1,
	  "customer_notify":1,
	  // "start_at":in.StartAt,
	  // "expire_by":in.ExpireBy,
	}
	body, err := client.Subscription.Create(data, nil);
	if err != nil {
		return nil, err
	}
	razorpay_subscription_id, ok := body["id"].(string);
	if ok == false {
		return nil, errors.New(RESPONSE_PARSE_ERROR);
	}
	return &razorpay_subscription_id, nil;
}

func GetRazorPaySubscription (razorpay_subscription_id string) (*pb.RazorpaySubscription, error){
	body, err := client.Subscription.Fetch(razorpay_subscription_id, nil, nil)
	if err != nil {
		return nil, err
	}
	return &pb.RazorpaySubscription {
		Id: body["id"].(string),
		Status : body["status"].(string),
		RazorPlanId : body["plan_id"].(string),
		CurrentStart : body["current_start"].(int64),
		CurrentEnd : body["current_end"].(int64),

		StartAt : body["start_at"].(int64),
		EndAt : body["end_at"].(int64),

		TotalCount : body["total_count"].(int32),
		PaidCount : body["paid_count"].(int64),
		RemainingCount : body["remaining_count"].(int32),

		ShortUrl : body["short_url"].(string),
	}, nil;
}

func ConvertRazorPaySubscriptionStatusToInt (status string) (int) {
	if status == "active" {
		return RAZOR_SUBSCRIPTION_ACTIVE;
	}else if status == "created" {
		return RAZOR_SUBSCRIPTION_CREATED;
	}else if status == "authenticated" {
		return RAZOR_SUBSCRIPTION_AUTHENTICATED;
	}else if status == "pending" {
		return RAZOR_SUBSCRIPTION_PENDING;
	}else if status == "halted" {
		return RAZOR_SUBSCRIPTION_HALTED;
	}else if status == "cancelled" {
		return RAZOR_SUBSCRIPTION_CANCELLED;
	}else if status == "completed" {
		return RAZOR_SUBSCRIPTION_COMPLETED;
	}else if status == "expired" {
		return RAZOR_SUBSCRIPTION_EXPIRED;
	}
	return RAZOR_SUBSCRIPTION_UNKNOWN;
}

//Subscription and Plans
