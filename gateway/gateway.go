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


package gateway

import (
	"errors"
	"bytes"
	"io"
	"os"
	"time"
	"github.com/google/uuid"

	"fmt"
	"log"

	"net/http"
	"encoding/json"

	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"	
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/cmd"

	hubService     "justify_backend/proto/hub"
	logisticsService "justify_backend/proto/logistics"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type EndPoints struct {
	AccountEndpoint		  *string
	SessionEndpoint 		*string
	CityEndpoint 				*string
	HubEndpoint 				*string
	LogisticsEndpoint 	*string
}

type EventBasicData struct {
	CityId int64 `json:"city_id"`
}

var hub *Hub;

const EVENT_MESSAGE = "message";
const WEBSOCKET_URL = "/ws/{city_id}";
const IMAGE_FOLDER string = "assets/images/%s/%s.%s";
const fileuploadURL string = "/api/files/upload";

var hubServiceClient hubService.HubServiceClient;
var logisticsServiceClient logisticsService.LogisticsServiceClient;

const REALTIME_TOPIC = "go.micro.topic.logi_live"

func SetupClients (endpoint *EndPoints, mux *runtime.ServeMux) error {	
	hubConn, err := grpc.Dial(*endpoint.HubEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("Hub Service Client Error!");
	}
	hubServiceClient = hubService.NewHubServiceClient(hubConn);

	logisticsConn, err := grpc.Dial(*endpoint.LogisticsEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("Logistics Service Client Error!");
	}
	logisticsServiceClient = logisticsService.NewLogisticsServiceClient(logisticsConn);

	err = mux.HandlePath("POST", fileuploadURL, HandleFileUpload);
	if err != nil {
		return err;
	}
	err = mux.HandlePath("GET", WEBSOCKET_URL, wsEndpoint);
	if err != nil {
		return err;
	}
	
	//Websocket handlers
	hub = NewHub()
	go hub.Run()

	//Kafka Setup
	cmd.Init();

	if err := broker.Init(); err != nil {
		log.Fatal(err);
	}

	if err := broker.Connect(); err != nil {
		log.Fatal(err); 
	}

	
	_, err = broker.Subscribe(REALTIME_TOPIC, func(p broker.Event) error {
		  var req EventBasicData;
			err := json.Unmarshal(p.Message().Body, &req);
			if err != nil {
				return err
			}
			if req.CityId == 0 {
				return errors.New("City Id is required!");
			}
		  payload := SocketEventStruct {
				EventName: EVENT_MESSAGE,
				EventPayload: string(p.Message().Body),
			}
			EmitToSpecificClient(hub, payload, fmt.Sprintf("%v", req.CityId))
			return nil
	});

	if err != nil {
		log.Fatal(err);
	}

	return nil
}


func HandleFileUpload(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	var response map[string]string;
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	multipartFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		response = map[string]string{"status": "400"}
		json.NewEncoder(w).Encode(response);
		return
	}
	fileSize := fileHeader.Size;
	fileBuffer := bytes.Buffer{}
	_, err = io.Copy(&fileBuffer, multipartFile);

	serviceName := r.Form.Get("issuer");
	requestType := r.Form.Get("request");
	var responseData *string;
	if(serviceName == "hub"){
		responseData, err = ProcessHubServiceFileRequest(&requestType, &fileBuffer, fileSize, r);
	}else if(serviceName == "logistics"){
		responseData, err = ProcessLogisticsServiceFileRequest(&requestType, &fileBuffer, fileSize, r);
	}else{
		err = errors.New("There was an error");
	}	
	if(err == nil){
		response = map[string]string{"status": "200", "data":*responseData}
	}else{
		response = map[string]string{"status": "400"}
	}	
	json.NewEncoder(w).Encode(response);
}

func ServeStatic() {
    route := mux.NewRouter()
    fs := http.FileServer(http.Dir("./assets/"))
    route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs));
    route.PathPrefix("/").Handler(http.FileServer(http.Dir("./track/")))
    fmt.Println("Serving Static on : 9191");
    log.Fatal(http.ListenAndServe(":9191", route))
}

func saveImage(file []byte, extension string) (*string, error) {
	dateFolder := time.Now().Format("200601");
	fileName, err := uuid.NewRandom()
	if err != nil {
		return nil, err;
	}
	imagePath := fmt.Sprintf(IMAGE_FOLDER, dateFolder, fileName, extension);
	ensureDir(imagePath);
	err = os.WriteFile(imagePath, file, 0644)
	if err != nil {
		return nil, err
	}
	return &imagePath, nil
}
func ensureDir(fileName string) {
  dirName := filepath.Dir(fileName)
  if _, serr := os.Stat(dirName); serr != nil {
    merr := os.MkdirAll(dirName, os.ModePerm)
    if merr != nil {
        panic(merr)
    }
  }
}

//Websocket
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request,  pathParams map[string]string) {
	// upgrade this connection to a WebSocket
	// connection
	city_id := r.URL.Query().Get("city_id")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	CreateNewSocketUser(hub, ws, city_id)
}