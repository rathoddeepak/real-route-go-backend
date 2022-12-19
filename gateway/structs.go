package gateway

import "github.com/gorilla/websocket"

// UserStruct is used for sending users with socket id
type UserStruct struct {
	CityId   string `json:"city_id"`
	UserId   string `json:"user_id"`
}

// SocketEventStruct struct of socket events
type SocketEventStruct struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub                 *Hub
	webSocketConnection *websocket.Conn
	send                chan SocketEventStruct
	CityId            	string
	UserId              string
}

// JoinDisconnectPayload will have struct for payload of join disconnect
type JoinDisconnectPayload struct {
	Users  []UserStruct `json:"users"`
	UserId string       `json:"user_id"`
}