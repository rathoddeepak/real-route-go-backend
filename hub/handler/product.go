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
	"errors"

	"context"

	"go-micro.dev/v4/logger"	

	pb "hubservice/proto"
	db "hubservice/database"
)

//Product Related Functions
func (hs *HubService) CreateProduct (ctx context.Context, in *pb.CreateProductRequest, out *pb.CreateProductResponse) (error) {	
   subCategory, err := db.GetSubCategoryById(in.SubCategoryId);
   if err != nil {
   	return errors.New("SubCategoryId Invalid!");
   }
   category, err := db.GetCategoryById(subCategory.MainCatId);
   if err != nil {
   	return errors.New("CategoryId Invalid!");
   }
   hub, err := db.GetHubById(category.HubId);
   if err != nil {
   	return errors.New("HubId Invalid!");
   }  
	product := &db.Product {
		Name: 			in.Name,
		CategoryId:    category.Id,
		SubCategoryId: subCategory.Id,
		HubId: 			hub.Id,
		CityId: 			hub.CityId,
		Price:  			in.Price,
		HigherPrice:  	in.HigherPrice,
		Qty:  			in.Qty,
		MaxLimit:		in.MaxLimit,
	}
	product_id, err := db.InsertProduct(product);
	if err != nil {
		logger.Info(err);
		return err
	}
	out.ProductId = *product_id;
	return nil;
}

func (hs *HubService) UpdateProduct (ctx context.Context, in *pb.UpdateProductRequest, out *pb.UpdateProductResponse) (error) {
	product := &db.Product {
		Id: 		 	 in.ProductId,
		Name: 		 in.Name,
		Price: 		 in.Price,
		HigherPrice: in.HigherPrice,
		Qty: 			 in.Qty,
		MaxLimit: 	 in.MaxLimit,
	}
	err := db.UpdateProduct(product);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateProductStatus (ctx context.Context, in *pb.UpdateProductStatusRequest, out *pb.UpdateProductResponse) (error) {
	product := &db.Product {
		Id: 	in.ProductId,
		Status: in.Status,
	}
	err := db.UpdateProductStatus(product);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateProductQty (ctx context.Context, in *pb.UpdateProductQtyRequest, out *pb.UpdateProductResponse) (error) {
	product := &db.Product {
		Id: 	in.ProductId,
		Qty: 	in.Qty,
		MaxLimit: in.MaxLimit,
	}
	err := db.UpdateProductQty(product);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateProductImage (ctx context.Context, in *pb.UpdateProductImageRequest, out *pb.UpdateProductResponse) (error) {	
	product := &db.Product {
		Id: 	  in.ProductId,
		Image: 	  in.Image,
		Blurhash: in.Blurhash,
	}
	if err := db.UpdateProductImage(product); err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) GetProductsOfCategory (ctx context.Context, in *pb.GetProductRequest, out *pb.GetProductsResponse) (error) {
	mProducts, err := db.GetProductOfCategory(in.CategoryId);
	if err != nil {
		return err
	}
	var products []*pb.Product;
	for _, mProduct := range *mProducts {
	   product := makeProtoProdcut(mProduct);
	   products = append(products, product);
	}
	out.Products = products;
	return nil;	
}

func (hs *HubService) GetProductsOfSubCategory (ctx context.Context, in *pb.GetProductRequest, out *pb.GetProductsResponse) (error) {
	mProducts, err := db.GetProductOfSubCategory(in.SubCategoryId);
	if err != nil {
		return err
	}
	var products []*pb.Product;
	for _, mProduct := range *mProducts {
	   product := makeProtoProdcut(mProduct);
       products = append(products, product);
    }
    out.Products = products;
	return nil;	
}

func (hs *HubService) GetProductById (ctx context.Context, in *pb.GetProductRequest, out *pb.GetProductResponse) (error) {
	mProduct, err := db.GetProductById(in.ProductId);
	if err != nil {
		return err
	}
	product := makeProtoProdcut(mProduct);
    out.Product = product;
	return nil;	
}

func (hs *HubService) NormalSearchProducts (ctx context.Context, in *pb.ProductSearchRequest, out *pb.NormalProductSearchResponse) (error) {
   mProducts, err := db.FilterProducts(in);
	if err != nil {
		return err
	}
	var products []*pb.Product;
	for _, mProduct := range *mProducts {
	   product := makeProtoProdcut(mProduct);
      products = append(products, product);
    }
    out.Products = products;
	return nil;
}

func (hs *HubService) SearchProducts (ctx context.Context, in *pb.ProductSearchRequest, out *pb.ProductSearchResponse) (error) {
   mProducts, err := db.FilterProducts(in);

	if err != nil {
		return err
	}
	var products []*pb.SearchProductModel;
	for _, mProduct := range *mProducts {
	   hub, err := db.GetHubById(mProduct.HubId);
	   if err != nil {
	   	continue
	   }
	   category, err := db.GetCategoryById(mProduct.CategoryId);
	   if err != nil {
	   	continue
	   }
	   subCategory, err := db.GetSubCategoryById(mProduct.SubCategoryId);
	   if err != nil {
	   	continue
	   }
	   product := makeProtoSearchProdcut(mProduct,hub,category,subCategory);
      products = append(products, product);
    }
    out.Products = products;
	return nil;
}

func makeProtoProdcut(product *db.Product) *pb.Product {
	return &pb.Product {	   	
	   Id  		  	  : product.Id,
		Name  		  : product.Name,
		CategoryId    : product.CategoryId,
		SubCategoryId : product.SubCategoryId,
		HubId  		  : product.HubId,
		CityId  		  : product.CityId,
		Price  		  : product.Price,
		HigherPrice   : product.HigherPrice,
		Qty  		  	  : product.Qty,
		MaxLimit  	  : product.MaxLimit,
		Status  	  	  : product.Status,
		Image  		  : product.Image,
		Blurhash  	  : product.Blurhash,
	}
}

func makeProtoSearchProdcut(product *db.Product, hub *db.Hub, category *db.Category, subcategory *db.SubCategory) *pb.SearchProductModel {
	return &pb.SearchProductModel {	   	
		Id 			: product.Id,
		Name 			: product.Name,
		Status 		: product.Status,
		Price 		: product.Price,
		HigherPrice : product.HigherPrice,
		Image 		: product.Image,
		Blurhash 	: product.Blurhash,
		CityId 		: product.CityId,

		HubId 		: hub.Id,
		HubName 		: hub.Name,

		MainCatId 	: category.Id,
		MainCatName : category.Name,

		SubCatId 	: subcategory.Id,
		SubCatName  : subcategory.Name,
	}
}