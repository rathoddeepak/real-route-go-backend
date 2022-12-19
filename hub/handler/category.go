/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 05 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Hub Microservice  <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (
	"context"

	"go-micro.dev/v4/logger"

	pb "hubservice/proto"
	db "hubservice/database"
)

//Category Related Functions
func (hs *HubService) CreateCategory (ctx context.Context, in *pb.CreateCategoryRequest, out *pb.CreateCategoryResponse) (error) {
	category := &db.Category {
		HubId: in.HubId,
		Name: in.Name,
	}
	category_id, err := db.InsertCategory(category);
	if err != nil {
		logger.Info(err);
		return err
	}
	out.CategoryId = *category_id;
	return nil;
}

func (hs *HubService) UpdateCategory (ctx context.Context, in *pb.UpdateCategoryRequest, out *pb.UpdateCategoryResponse) (error) {
	category := &db.Category {
		Id: 	in.CategoryId,
		Name: 	in.Name,
	}
	err := db.UpdateCategory(category);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateCategoryStatus (ctx context.Context, in *pb.UpdateCategoryStatusRequest, out *pb.UpdateCategoryResponse) (error) {
	category := &db.Category {
		Id: 	in.CategoryId,
		Status: in.Status,
	}
	err := db.UpdateCategoryStatus(category);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) GetHubCategories (ctx context.Context, in *pb.GetCategoryRequest, out *pb.GetCategoriesResponse) (error) {
	mCategories, err := db.GetCategoriesOfHub(in.HubId);
	if err != nil {
		return err
	}
	var categories []*pb.Category;
	for _, mCategory := range *mCategories {
	   category := makeProtoCategory(mCategory);
       categories = append(categories, category);
    }
    out.Categories = categories;
	return nil;	
}

func (hs *HubService) GetCategoryById (ctx context.Context, in *pb.GetCategoryRequest, out *pb.GetCategoryResponse) (error) {
	mCategory, err := db.GetCategoryById(in.CategoryId);
	if err != nil {
		return err
	}
	category := makeProtoCategory(mCategory);
    out.Category = category;
	return nil;	
}

func (hs *HubService) UpdateCategoryImage (ctx context.Context, in *pb.UpdateCategoryImageRequest, out *pb.UpdateCategoryResponse) (error) {	
	category := &db.Category {
		Id: 	  in.CategoryId,
		Image: 	  in.Image,
		Blurhash: in.Blurhash,
	}
	if err := db.UpdateCategoryImage(category); err != nil {
		return err
	}
	out.Status = 200;
	out.Message = in.Image;
	return nil;
}

func makeProtoCategory(mCategory *db.Category) *pb.Category {
	return &pb.Category {
	   	Id 	   	 : 	mCategory.Id,
	   	HubId 	 : 	mCategory.HubId,
	   	Name 	 : 	mCategory.Name,
	   	Status 	 : 	mCategory.Status,
	   	Image 	 : 	mCategory.Image,
	   	Blurhash : 	mCategory.Blurhash,
	}
}