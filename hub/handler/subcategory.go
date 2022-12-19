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
 --->     Sub Category 	  <---
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
func (hs *HubService) CreateSubCategory (ctx context.Context, in *pb.CreateSubCategoryRequest, out *pb.CreateSubCategoryResponse) (error) {
	subcategory := &db.SubCategory {
		MainCatId: in.MainCatId,
		Name: in.Name,
	}
	subcategory_id, err := db.InsertSubCategory(subcategory);
	if err != nil {
		logger.Info(err);
		return err
	}
	out.SubCategoryId = *subcategory_id;
	return nil;
}

func (hs *HubService) UpdateSubCategory (ctx context.Context, in *pb.UpdateSubCategoryRequest, out *pb.UpdateSubCategoryResponse) (error) {
	subcategory := &db.SubCategory {
		Id: 	in.SubCategoryId,
		Name: 	in.Name,
	}
	err := db.UpdateSubCategory(subcategory);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateSubCategoryStatus (ctx context.Context, in *pb.UpdateSubCategoryStatusRequest, out *pb.UpdateSubCategoryResponse) (error) {
	subcategory := &db.SubCategory {
		Id: 	in.SubCategoryId,
		Status: in.Status,
	}
	err := db.UpdateSubCategoryStatus(subcategory);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) GetSubCategoriesOfCategory (ctx context.Context, in *pb.GetSubCategoryRequest, out *pb.GetSubCategoriesResponse) (error) {
	mSubCategories, err := db.GetSubCategoriesOfCategory(in.MainCatId);
	if err != nil {
		return err
	}
	var subcategories []*pb.SubCategory;
	for _, mSubCategory := range *mSubCategories {
	   subcategory := makeProtoSubCategory(mSubCategory);
       subcategories = append(subcategories, subcategory);
    }
    out.SubCategories = subcategories;
	return nil;	
}

func (hs *HubService) GetSubCategoryById (ctx context.Context, in *pb.GetSubCategoryRequest, out *pb.GetSubCategoryResponse) (error) {
	mSubCategory, err := db.GetSubCategoryById(in.SubCategoryId);
	if err != nil {
		return err
	}
	subcategory := makeProtoSubCategory(mSubCategory);
    out.SubCategory = subcategory;
	return nil;	
}

func (hs *HubService) UpdateSubCategoryImage (ctx context.Context, in *pb.UpdateSubCategoryImageRequest, out *pb.UpdateSubCategoryResponse) (error) {	
	subCategory := &db.SubCategory {
		Id: 	  in.SubCategoryId,
		Image: 	  in.Image,
		Blurhash: in.Blurhash,
	}
	if err := db.UpdateSubCategoryImage(subCategory); err != nil {
		return err
	}
	out.Status = 200;
	out.Message = in.Image;
	return nil;
}

func makeProtoSubCategory(mSubCategory *db.SubCategory) *pb.SubCategory {
	return &pb.SubCategory {
	   	Id 	   		: 	mSubCategory.Id,
	   	MainCatId 	: 	mSubCategory.MainCatId,
	   	Name 		: 	mSubCategory.Name,
	   	Status 		: 	mSubCategory.Status,   	
	   	Image 	 	: 	mSubCategory.Image,
		Blurhash 	: 	mSubCategory.Blurhash,
	}
}