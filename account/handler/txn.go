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

package handler

import (
	
	"context"

	"errors"

	pb "accountservice/proto"

	db "accountservice/database"
)

func (service *AccountService) CreditWallet (ctx context.Context, in *pb.WalletTxnRequest, out *pb.TxnResponse) error {
	if user, err := db.GetUserById(in.UserId); err == nil {
		newAmount := user.Wallet + in.Amount;
		txn := &db.Txn {
			UserId:    in.UserId,
			Amount:    in.Amount,
			TransType: db.CREDITED,
		}
		txn_id, err := db.CreateTXN(txn);
		if err != nil {
			return err;
		}
		db.UpdateUserWallet(user.Id, newAmount);
		out.TxnId = txn_id;
		return nil;
	}else{
		return errors.New(userNotFound);
	}
	return nil;
}

func (service *AccountService) DebitWallet (ctx context.Context, in *pb.WalletTxnRequest, out *pb.TxnResponse) error {
	if user, err := db.GetUserById(in.UserId); err == nil {
		newAmount := user.Wallet - in.Amount;
		if(newAmount < 0){
			newAmount = 0.0;
		}
		txn := &db.Txn {
			UserId:    in.UserId,
			Amount:    in.Amount,
			TransType: db.DEBITED,
		}
		txn_id, err := db.CreateTXN(txn);
		if err != nil {
			return err;
		}
		db.UpdateUserWallet(user.Id, newAmount);
		out.TxnId = txn_id;
		return nil;
	}else{
		return errors.New(userNotFound);
	}
	return nil;
}

func (service *AccountService) GetTxns (ctx context.Context, in *pb.GetTxnRequest, out *pb.GetTxnsResponse) error {
	var limit int64 = 10;
	if in.Limit != 0 {
		limit = in.Limit;
	}
	mTxns, err := db.GetUserTxn(in.UserId, limit, in.Offset);
	if err != nil {
		return err
	}
	var txns []*pb.Txn;
	for _, mTxn := range *mTxns {
	   txn := makeProtoTxn(mTxn);
	   txns = append(txns, txn);
	}
	out.Transactions = txns;
	return nil;	
}

func (service *AccountService) GetTxnById (ctx context.Context, in *pb.GetTxnRequest, out *pb.GetTxnResponse) error {
	mTxn, err := db.GetTxnById(in.UserId);
	if err != nil {
		return err
	}
	out.Transaction = makeProtoTxn(mTxn);
	return nil;	
}

func makeProtoTxn(txn *db.Txn) *pb.Txn {
	return &pb.Txn {	   	
	   Id		: txn.Id,
	   UserId	: txn.UserId,
	   TransType: txn.TransType,
	   Amount	: txn.Amount,
	   Created	: txn.Created,
	}
}