/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 18 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 ---> LOGISTICS Microservice <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

// This page has functions realated to Subscription & One time order and most of the constants

package handler;

import (
    "fmt"
    
	"time"

	"context"

    "errors"

    "go-micro.dev/v4/logger"

	pb "logisticsService/proto"
    db "logisticsService/database"
)

const USER_ACTIVE int32 = 1;

//Delivery Instructions
const INHAND = 0;
const INBAG = 1;
const RING = 2;
const INBAG_RING = 3;

const AREA_NOT_SERVICABLE = "Area Not Servicable!";
const INSUFFICIENT_BALANCE = "Insufficient Balance, Recharge Your Wallet!";
const UPDATED_MSG = "Updated!";
const TASK_REQUIRED_ERR = "Task is required";
const AGENT_REQUIRED_ERR = "Agent is required";

const SUBSCRIPTION_DISPATCH = 0;

type LogisticsService struct {
	HubService pb.HubService;
	AccountService pb.AccountService;
	CityService pb.CityService;
}

const notificationTopic = "go.micro.topic.sendnotification"

func (lg *LogisticsService) CreateSubscription (ctx context.Context, in *pb.CreateSubscriptionRequest, out *pb.CreateSubscriptionResponse) (error) {
	now := time.Now();
    yyyy, mm, dd := now.Date();
    tomorrow := time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, now.Location());

    if in.StartFrom < tomorrow.Unix() {
    	return errors.New("Invalid Start From Date!");
    }

    running := db.CheckSubscriptionRunning(in.ProductId, in.UserId);
    if running == true {
        return errors.New("Subscription Already Created!");
    }

    p, err := lg.HubService.GetProductById(ctx, &pb.GetProductRequest{ProductId:in.ProductId});
    if err != nil {
    	return errors.New("Product Fetch Error!");
    }

    u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:in.UserId});
    if err != nil {
        return errors.New("User Not found!");
    }

    respAddress, err := lg.AccountService.GetAddressById(ctx, &pb.GetAddressRequest {AddressId:in.AddressId});
    if err != nil {
        return errors.New("Address Not found!");
    }

    mfence, err := lg.CityService.GetFenceByGeoPointAndId(ctx, &pb.GetFenceRequest {
        Lat   : respAddress.Address.Lat,
        Lng   : respAddress.Address.Lng,
        HubId : p.Product.HubId,
    });
    if err != nil {
        fmt.Println(err);
        return errors.New(AREA_NOT_SERVICABLE);
    }else if mfence.Fence.Status != 1 {        
        return errors.New(AREA_NOT_SERVICABLE)
    }

    allQtyZero := true;
    qtyList := []float64 {in.Sun, in.Mon, in.Tue, in.Wed, in.Thu, in.Fri, in.Sat};
    for _, val := range qtyList {
    	if val > float64(p.Product.MaxLimit) {
    		return errors.New("Order Max Limit Exceeded!");
    	}
    	if val != 0 {	    	
	    	allQtyZero = false;
	    }
    }
    if allQtyZero {
    	return errors.New("Qty Cannot be zero!");
    }

    walletBalace := CalculateWalletBalance(in.UserId, u.User.Wallet, ctx, lg);
    amount := CalculateSubscriptionCost(in.StartFrom, in.DlvDays, qtyList, p.Product.Price);
    if amount > walletBalace {
        return errors.New(INSUFFICIENT_BALANCE);
    }
    expEnd := time.Unix(in.StartFrom, 0).AddDate(0, 0, int(in.DlvDays)).Unix();
    subscription := &db.Subscription {
    	StartFrom : in.StartFrom,
    	ProductId : p.Product.Id,
    	UserId    : in.UserId,
    	Pattern   : in.Pattern,
    	Sun		  : in.Sun, 
    	Mon		  : in.Mon, 
    	Tue		  : in.Tue, 
    	Wed		  : in.Wed, 
    	Thu		  : in.Thu, 
    	Fri		  : in.Fri, 
    	Sat		  : in.Sat,
    	DlvIns	  : in.DlvIns,
    	DlvDays   : in.DlvDays,
        AddressId : in.AddressId,
        SlotId    : in.SlotId,
        CityId    : p.Product.CityId,
        HubId     : p.Product.HubId,
        ExpEnd    : expEnd,
    }
    subId, err := db.InsertSubscription(subscription);
    if err != nil {
    	return err;
    }

    out.SubscriptionId = *subId;
    return nil;    
}

func (lg *LogisticsService) UpdateSubscription (ctx context.Context, in *pb.UpdateSubscriptionRequest, out *pb.UpdateSubscriptionResponse) (error) {
	
    subscription, err := db.GetSubscriptionById(in.SubscriptionId);
    if err != nil {
    	return errors.New("Subscription not found!");
    }

    if subscription.StartFrom != in.StartFrom {
        now := time.Now();
        yyyy, mm, dd := now.Date();
        tomorrow := time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, now.Location());
        if in.StartFrom < tomorrow.Unix() {
            return errors.New("Invalid Start From Date!");
        }
    }    

    p, err := lg.HubService.GetProductById(ctx, &pb.GetProductRequest{ProductId:subscription.ProductId});
    if err != nil {
    	return errors.New("Product Fetch Error!");
    }

    u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:subscription.UserId});
    if err != nil {
        return errors.New("User Not found!");
    }

    respAddress, err := lg.AccountService.GetAddressById(ctx, &pb.GetAddressRequest {AddressId:subscription.AddressId});
    if err != nil {
        return errors.New("Address Not found!");
    }

    mfence, err := lg.CityService.GetFenceByGeoPointAndId(ctx, &pb.GetFenceRequest {
        Lat   : respAddress.Address.Lat,
        Lng   : respAddress.Address.Lng,
        HubId : p.Product.HubId,
    });
    if err != nil {
        return errors.New(AREA_NOT_SERVICABLE);
    }else if mfence.Fence.Status != 1 {
        return errors.New(AREA_NOT_SERVICABLE)
    }

    allQtyZero := true;
    qtyList := []float64 {in.Sun, in.Mon, in.Tue, in.Wed, in.Thu, in.Fri, in.Sat};
    for _, val := range qtyList {
    	if val > float64(p.Product.MaxLimit) {
    		return errors.New("Order Max Limit Exceeded!");
    	}
    	if val != 0 {	    	
	    	allQtyZero = false;
	    }
    }
    if allQtyZero {
    	return errors.New("Qty Cannot be zero!");
    }

    walletBalace := CalculateWalletBalance(subscription.UserId, u.User.Wallet, ctx, lg);
    amount := CalculateSubscriptionCost(in.StartFrom, in.DlvDays, qtyList, p.Product.Price);
    if amount > walletBalace {
    	return errors.New(INSUFFICIENT_BALANCE);
    }
    expEnd := time.Unix(in.StartFrom, 0).AddDate(0, 0, int(in.DlvDays)).Unix();
    mSubscription := &db.Subscription {
    	Id  	  : in.SubscriptionId,
    	StartFrom : in.StartFrom,
    	Pattern   : in.Pattern,
    	Sun		  : in.Sun, 
    	Mon		  : in.Mon, 
    	Tue		  : in.Tue, 
    	Wed		  : in.Wed, 
    	Thu		  : in.Thu, 
    	Fri		  : in.Fri, 
    	Sat		  : in.Sat,
    	DlvIns	  : in.DlvIns,
    	DlvDays   : in.DlvDays,
        AddressId : in.AddressId,
        SlotId    : in.SlotId,
        ExpEnd    : expEnd,
    }
    err = db.UpdateSubscription(mSubscription);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;    
}

func (lg *LogisticsService) UpdateSubscriptionStatus (ctx context.Context, in *pb.UpdateSubscriptionRequest, out *pb.UpdateSubscriptionResponse) (error) {
	subscription := &db.Subscription {
    	Id  	  : in.SubscriptionId,
    	Status    : in.Status,
    }
    err := db.UpdateSubscriptionStatus(subscription);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;    
}

func (lg *LogisticsService) GetSubscriptionById (ctx context.Context, in *pb.GetSubscriptionRequest, out *pb.GetSubscriptionResponse) (error) {
    mScp, err := db.GetSubscriptionById(in.SubscriptionId);
    if err != nil {
        return err
    }
    scp := makeProtoSubscription(mScp);
    out.Subscription = scp;
    return nil;
}

func (lg *LogisticsService) DetermineWalletBalance (ctx context.Context, in *pb.DetermineWalletBalanceRequest, out *pb.DetermineWalletBalanceResponse) (error) {
    balance := CalculateWalletBalance(in.UserId, in.CurrentAmount, ctx, lg);
    out.Balance = balance;
    return nil;    
}

func (lg *LogisticsService) GetInfoSubscriptions (ctx context.Context, in *pb.FilterSubscriptionRequest, out *pb.GetInfoSubscriptionsResponse) (error) {    
    mSubscriptions, _, err := db.FilterSubscriptions(in, false);
    if err != nil {
        return err;
    }
    var todayTime time.Time;
    if in.FromDate == 0 {
        todayTime = time.Now();
    }else{
        todayTime = time.Unix(in.FromDate, 0);
    }
    var subscriptions []*pb.SubscriptionInfo;
    for _, mSubscription := range *mSubscriptions {        
        var subscription pb.SubscriptionInfo;        

        u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:mSubscription.UserId});
        if err != nil {
            fmt.Println("account", err);
            continue
        }
        p, err := lg.HubService.GetProductById(ctx, &pb.GetProductRequest{ProductId:mSubscription.ProductId});
        if err != nil {
            fmt.Println("hub", err);
            continue
        }

        renewDate := time.Unix(mSubscription.StartFrom, 0).AddDate(0, 0, int(mSubscription.DlvDays)).Unix();
        currentStamp := time.Now().Unix();
        if renewDate < currentStamp {
            renewDate = 0;
        }

        balance := CalculateWalletBalance(u.User.Id, u.User.Wallet, ctx, lg);
        qty := getQtyFromSubscription(&todayTime, mSubscription);

        subscription.Status = &pb.SubsStatus {
            Code    : mSubscription.Status,
            Message : getSubscriptionStatusMessage(mSubscription.Status),
        }
        subscription.Customer = &pb.Customer {
            Id      : u.User.Id,
            Name    : u.User.Name,
            Wallet  : u.User.Wallet,
            Phone   : u.User.Phone,
            Balance : balance,
        }
        subscription.Product = &pb.ProductItem {
            Id   : p.Product.Id,
            Name : p.Product.Name,
            Qty  : qty,
        }
        subscription.RenewDate = renewDate;
        subscription.Created = mSubscription.Created;
        subscription.Id = mSubscription.Id;
        subscriptions = append(subscriptions, &subscription);
    }
    out.Subscriptions = subscriptions;
    return nil;    
}

func (lg *LogisticsService) GetSubscriptionSummary (ctx context.Context, in *pb.FilterSubscriptionRequest, out *pb.GetSubscriptionSummaryResponse) (error) {    
    _, subscriptionCount, err := db.FilterSubscriptions(in, true);
    if err != nil {
        return err;
    }
    var summary []*pb.KeyValue;
    summary = append(summary, &pb.KeyValue{
        Key: "Total Subscription",
        Value: fmt.Sprintf("%v", *subscriptionCount),
    });
    out.Summary = summary;
    return nil;    
}

func (lg *LogisticsService) GetInventoryData (ctx context.Context, in *pb.FilterSubscriptionRequest, out *pb.GetInvetoryDataResponse) (error) {    
    subscriptions, _, err := db.FilterSubscriptions(in, false);
    if err != nil {
        return err;
    }    
    findProductId := func(id int64, list *[]*pb.InventoryProduct) (int) {
        var idx = 0;
        for _, value := range *list {            
            if value.Id == id {
                return idx;
            }
            idx++;
        }
        return -1;
    }
    var todayTime time.Time;
    if in.FromDate == 0 {
        todayTime = time.Now();
    }else{
        todayTime = time.Unix(in.FromDate, 0);
    }
    var products []*pb.InventoryProduct;
    for _, subscription := range *subscriptions {
        idx := findProductId(subscription.ProductId, &products);
        qty := getQtyFromSubscription(&todayTime, subscription);
        if idx == -1 {
            p, err := lg.HubService.GetProductById(ctx, &pb.GetProductRequest {ProductId:subscription.ProductId});
            if err != nil {
                logger.Info(err);
                continue
            }
            h, err := lg.HubService.GetHubById(ctx, &pb.GetHubRequest {HubId:subscription.HubId})
            if err != nil {
                logger.Info(err);
                continue   
            }
            products = append(products, &pb.InventoryProduct {
                Id     : subscription.ProductId,
                Name   : p.Product.Name,
                HubId  : p.Product.HubId,
                HubName: h.Hub.Name,
                Qty    : qty,
            });
        } else {
            products[idx].Qty = products[idx].Qty + qty;
        }
    }
    out.Products = products;
    return nil;    
}

func (lg *LogisticsService) GetInventoryCustomers (ctx context.Context, in *pb.FilterSubscriptionRequest, out *pb.GetInvetoryCustomersResponse) (error) {    
    subscriptions, _, err := db.FilterSubscriptions(in, false);
    if err != nil {
        return err;
    }
    var todayTime time.Time;
    if in.FromDate == 0 {
        todayTime = time.Now();
    }else{
        todayTime = time.Unix(in.FromDate, 0);
    }
    var customers []*pb.InventoryCustomer;
    for _, subscription := range *subscriptions {
        u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:subscription.UserId});
        if err != nil {
            logger.Info(err);
            continue
        }
        qty := getQtyFromSubscription(&todayTime, subscription);
        customers = append(customers, &pb.InventoryCustomer {
            Id     : subscription.UserId,
            Name   : u.User.Name,
            Qty    : qty,
        });
    }
    out.Customers = customers;
    return nil;    
}

func CalculateSubscriptionCost (startFrom int64, dlvDays int64, weekList []float64, cost float64) (float64) {
	weekDay := time.Unix(startFrom, 0).Weekday();
	var amount float64;
	for i := int64(0); i < dlvDays; i++ {
		if weekDay == 6 {
			weekDay = 0;
		}else{
			weekDay++;
		}
		amount += weekList[weekDay] * cost;
	}
	return amount;
}

func CalculateWalletBalance(user_id int64, userWallet float64, ctx context.Context, lg *LogisticsService) (float64) {
    subscriptions, err := db.GetSubscriptions(user_id, 1);
    if err != nil {
        fmt.Println(err)
        return userWallet;
    }
    for _, mSub := range *subscriptions {        
        endDate := time.Unix(mSub.ExpEnd, 0);
        endDate = endDate.AddDate(0, 0, int(mSub.DlvDays));
        currentDate := time.Now();
        diff := endDate.Sub(currentDate);
        var remainingDays = int(diff.Hours()/24);
        if remainingDays <= 0 {
            continue
        }

        presp, err := lg.HubService.GetProductById(ctx, &pb.GetProductRequest {ProductId:mSub.ProductId});
        if err != nil {
            logger.Info(err);
            continue
        }

        amount := float64(remainingDays) * presp.Product.Price;
        userWallet -= amount;
    }
    if userWallet < 0 {
        userWallet = float64(0)
    }
    return userWallet;
}

func RemoveItem[T any](slice []T, s int) []T {
    return append(slice[:s], slice[s+1:]...)
}

func getSubscriptionStatusMessage (status int32) (string) {
    if status == 1 {
        return "Active";
    } else {
        return "Inactive";
    }
}

func makeProtoSubscription (scp *db.Subscription) (*pb.Subscription){
    return &pb.Subscription {
        Id        : scp.Id,
        StartFrom : scp.StartFrom,
        ProductId : scp.ProductId,
        UserId    : scp.UserId,
        Pattern   : scp.Pattern,
        Sun       : scp.Sun, 
        Mon       : scp.Mon, 
        Tue       : scp.Tue, 
        Wed       : scp.Wed, 
        Thu       : scp.Thu, 
        Fri       : scp.Fri, 
        Sat       : scp.Sat,
        DlvIns    : scp.DlvIns,
        DlvDays   : scp.DlvDays,
        AddressId : scp.AddressId,
        SlotId    : scp.SlotId,
        CityId    : scp.CityId,
        HubId     : scp.HubId,
        ExpEnd    : scp.ExpEnd,
        Status    : scp.Status,
    }
}

func getQtyFromSubscription (t *time.Time, scp *db.Subscription) (float64) {
    weekday := t.Weekday();
    switch weekday {
        case 0:
            return scp.Sun;
        case 1:
            return scp.Mon;
        case 2:
            return scp.Tue;
        case 3:
            return scp.Wed;
        case 4:
            return scp.Thu;
        case 5:
            return scp.Sat;
        case 6:
            return scp.Sun;
    }
    return float64(0);
}