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
	"context"
	db "sessionservice/database"
	pb "sessionservice/proto"
)

type SessionService struct {}

func (service *SessionService) GetSessionById(ctx context.Context, in *pb.GetSessionByIdRequest, out *pb.GetSessionResponse) error {
	session, err := db.GetSessionById(in.SessionId);
	if err != nil {
		return err;
	}
	out.SessionId = session.Id;
	out.SessionKey = session.SessionKey;
	return nil;
}

func (service *SessionService) GetSession(ctx context.Context, in *pb.GetSessionRequest, out *pb.GetSessionResponse) error {
	session, err := db.GetSession(in.SessionKey);
	if err != nil {
		return err;
	}
	out.SessionId = session.Id;
	out.UserId = session.UserId;
	out.SessionKey = session.SessionKey;
	return nil;
}

func (service *SessionService) GetUserSession(ctx context.Context, in *pb.GetSessionRequest, out *pb.GetSessionsResponse) error {
	mSessions, err := db.GetUserSession(in.UserId);
	if err != nil {
		return err;
	}
	var sessions []*pb.Session;
	for _, mSession := range *mSessions {
	   session := makeProtoSesseion(mSession);
	   sessions = append(sessions, session);
	}
	out.Sessions = sessions;
	return nil;
}

func (service *SessionService) CreateSession(ctx context.Context, in *pb.CreateSessionRequest, out *pb.GetSessionResponse) error {
	sessionId, sessionKey, err := db.InsertSession(in.UserId, in.ClientData);
	if err != nil {
		return err;
	}
	out.SessionId = int64(sessionId);
	out.SessionKey = sessionKey;
	return nil;
}

func makeProtoSesseion(session *db.Session) *pb.Session {
	return &pb.Session {	   	
		Id  		: session.Id,		
		SessionKey 	: session.SessionKey,    
		UserId 		: session.UserId,
		ClientData 	: session.ClientData,
		CreatedAt 	: session.CreatedAt,
	}
}