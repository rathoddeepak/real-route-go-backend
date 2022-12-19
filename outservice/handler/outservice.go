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

package handler;


import (
	
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
	"encoding/json"

	"go-micro.dev/v4/logger"
	"context"
	pb "outservice/proto"
)

const OS_MANAGER_APP_KEY = "5a3162fc-def0-4cf9-b4af-5346f85329a7";
const OS_MANAGER_APP_ATH = "YzJkNzMyMmYtMzIxYi00OTY3LTg0Y2QtM2FiMWY4ZTg3MDY3";
const OS_MANAGER_APP_CHID = "4aa3ec2d-0b54-40ee-b5a6-1ffeb614c339";

const OS_AGENT_APP_KEY = "1427e7ac-aa4a-4c0f-8d67-2a95c677c3b4";
const OS_AGENT_APP_ATH = "MTFjNTIzMzctMDAxMS00ZjNkLTlhYTItNjQwZTc5ZDQ0NGE2";
const OS_AGENT_APP_CHID = "8fa7e40d-7683-4a0c-9a7e-fb1f63c88977";

type OutService struct {}

func (service *OutService) SendSMS (ctx context.Context, in *pb.SendSMSRequest, out *pb.SendSMSResponse) error {
	logger.Info("Message:", in.Body);
	logger.Info("Message Sent !");
	if out != nil {
		out.Status = 200;
		out.Msg = "Message Sent Successfully!";
	}	
	return nil;
}

func (service *OutService) SendMobileNotification (ctx context.Context, in *pb.SendNotificaitonRequest, out *pb.SendSMSResponse) error {
	var list []int64;
	var filters []map[string]interface {};
	key, app_key, app_auth, app_channel_id := "", "", "", "";

  if in.Title == "" || in.Content == "" {
     return nil;
  }

  if len(in.CityIds) != 0 {
		key = "city_id";
		list = in.CityIds;
		app_key = OS_MANAGER_APP_KEY;
		app_auth = OS_MANAGER_APP_ATH;      
    app_channel_id = OS_MANAGER_APP_CHID;
  }else if len(in.AgentIds) != 0 {
		key = "agent_id";
		list = in.AgentIds;
		app_key = OS_AGENT_APP_KEY;
		app_auth = OS_AGENT_APP_ATH;
		app_channel_id = OS_AGENT_APP_CHID;
  }else{
  	out.Status = 200;
		out.Msg = "No Notification Recepient!";
  	return nil;
  }

	if out != nil {
		out.Status = 200;
		out.Msg = "Notification Sent Successfully!";
	}

	listLength := len(list);

	if listLength > 0 {
    lastIdx := listLength - 1;
    for i := 0; i < listLength; i++ {
        filter := map[string] interface {} {
		      	"field" : "tag",
		      	"key" : key,
		      	"relation" : "=",
		      	"value" : list[i],
		    }	      
		    filters = append(filters, filter);
		    if lastIdx != i {
		    	filter = map[string] interface {} {
		      	"operator" : "OR",
		      }
		    	filters = append(filters, filter);  
		    }
    }
  }else{
		filter := map[string] interface {} {
			"field" : "tag",
			"key" : key,
			"relation" : "=",
			"value" : list[0],
		}	      
		filters = append(filters, filter);
  }

  var data map[string] interface {}; 
  if in.Data != "" {
  	data["data"] = in.Data;
  }
  fields := map[string]interface {} {
    "contents":map[string] interface {} {
    	"en" : in.Content,
    },
    "headings":map[string] interface {} {
    	"en" : in.Title,
    },
    "android_led_color":"FF0000FF",
    "android_sound":"notify",
    "app_id":app_key,
    "isChrome":false,
    "priority":10,
    "filters":filters,
    "data":data,
  };
    // if(app_channel_id == null){
    //   fields['android_channel_id'] = channel_id;
    // }else{
    //   fields['existing_android_channel_id'] = app_channel_id;
    // }
  fields["android_channel_id"] = app_channel_id;

	url := "https://onesignal.com/api/v1/notifications"   
	body, err := json.Marshal(fields);
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", app_auth));
	req.Header.Set("Content-Type", "application/json;charset=utf-8");
	client := &http.Client{};
	resp, err := client.Do(req);
	if err != nil {
	    panic(err)
	}

	defer resp.Body.Close();

	fmt.Println("response Status:", resp.Status);
	fmt.Println("response Headers:", resp.Header);
	body, _ = ioutil.ReadAll(resp.Body);
	fmt.Println("response Body:", string(body));

	// out.Status = 200;
	// out.Msg = "Notification Sent Successfully!";

	return nil;
}