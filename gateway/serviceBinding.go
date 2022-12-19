package gateway

import (
	"context"

	"errors"
	"strconv"
	"bytes"

	"net/http"

	hubService     "justify_backend/proto/hub"
	logisticsService     "justify_backend/proto/logistics"
)

func ProcessHubServiceFileRequest(requestType *string, file *bytes.Buffer, fileSize int64, r *http.Request) (*string, error) {
	imageType 	   := r.Form.Get("image_type");
	if(*requestType == "uploadImage"){
		id      			 := r.Form.Get("product_id");
		imagePath, err := saveImage(file.Bytes(), imageType);
		if err != nil {
			return nil, err
		}
		productId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		_, err = hubServiceClient.UpdateProductImage(context.Background(), &hubService.UpdateProductImageRequest{
			Image: *imagePath,
			ProductId: int64(productId),
			Blurhash: "none",
		});
		if err != nil {
			return nil, err
		}
		return imagePath, nil;
	}else if(*requestType == "updateCatImage"){
		id      			 := r.Form.Get("category_id");		
		imagePath, err := saveImage(file.Bytes(), imageType);
		if err != nil {
			return nil, err
		}
		category_id, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		res, err := hubServiceClient.UpdateCategoryImage(context.Background(), &hubService.UpdateCategoryImageRequest{
			Image: *imagePath,
			CategoryId: int64(category_id),
			Blurhash: "none",
		});
		if err != nil {
			return nil, err
		}
		return &res.Message, nil;
	}else if(*requestType == "updateSubCatImage"){
		id      			 := r.Form.Get("sub_category_id");		
		imagePath, err := saveImage(file.Bytes(), imageType);
		if err != nil {
			return nil, err
		}
		sub_category_id, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		res, err := hubServiceClient.UpdateSubCategoryImage(context.Background(), &hubService.UpdateSubCategoryImageRequest{
			Image: *imagePath,
			SubCategoryId: int64(sub_category_id),
			Blurhash: "none",
		});
		if err != nil {
			return nil, err
		}
		return &res.Message, nil;
	}else{
		return nil, errors.New("Invalid type!");
	}
}

func ProcessLogisticsServiceFileRequest(requestType *string, file *bytes.Buffer, fileSize int64, r *http.Request) (*string, error) {
	imageType 	   := r.Form.Get("image_type");
	if(*requestType == "uploadAgentAvatar"){
		id      			 := r.Form.Get("agent_id");
		imagePath, err := saveImage(file.Bytes(), imageType);
		if err != nil {
			return nil, err
		}
		agentId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		_, err = logisticsServiceClient.UpdateAgentAvatar(context.Background(), &logisticsService.UpdateAgentRequest {
			Avatar: *imagePath,
			AgentId: int64(agentId),
		});
		if err != nil {
			return nil, err
		}
		return imagePath, nil;
	}else if(*requestType == "uploadVehiclePhoto"){
		id      			 := r.Form.Get("vehicle_id");
		imagePath, err := saveImage(file.Bytes(), imageType);
		if err != nil {
			return nil, err
		}
		vehicle_id, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		_, err = logisticsServiceClient.UpdateVehicleImage(context.Background(), &logisticsService.UpdateVehicleRequest {
			Image: *imagePath,
			VehicleId: int64(vehicle_id),
		});
		if err != nil {
			return nil, err
		}
		return imagePath, nil;
	}else{
		return nil, errors.New("Invalid type!");
	}
}