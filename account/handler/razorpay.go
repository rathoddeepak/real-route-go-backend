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

func (service *AccountService) GetRazorpayPlans (ctx context.Context, in *pb.GetRazorPayPlansRequest, out *pb.GetRazorPayPlansResponse) error {
	plans, err := db.GetRazorpayPlans(in);
	if err != nil {
		return err
	}
	out.Plans = *plans;
	return nil;
}