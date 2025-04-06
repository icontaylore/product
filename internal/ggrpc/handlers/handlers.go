package handlers

import (
	"context"
	"fmt"
	"log"
	pb "product/generated_proto/workspace/gen/product"
	"product/internal/infrastructure/models"
	"product/internal/repository"
)

type GetServiceHandler struct {
	repo repository.RepoInterface
}

func NewGetServiceHandler(repo repository.RepoInterface) *GetServiceHandler {
	return &GetServiceHandler{
		repo: repo,
	}
}

func (g *GetServiceHandler) CreateProduct(ctx context.Context, in *pb.CreateProductReq) (*pb.CreateProductResp, error) {
	newProduct := &models.Product{
		Name:          in.GetName(),
		Description:   in.GetDescription(),
		Price:         in.GetPrice(),
		StockQuantity: in.GetStockQuantity(),
	}
	if err := g.repo.Create(newProduct); err != nil {
		fmt.Println("handler proto:err create")
		return nil, err
	}

	return &pb.CreateProductResp{
		Success: true,
	}, nil
}

func (g *GetServiceHandler) GetProduct(ctx context.Context, id *pb.GetProductReq) (*pb.GetProductResp, error) {
	productID := id.GetId()

	product, err := g.repo.Read(uint(productID))
	if err != nil {
		log.Println("handler proto:read err")
		return nil, err
	}
	return &pb.GetProductResp{
		Id:            int32(product.Id),
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}, nil
}

func (g *GetServiceHandler) UpdateProduct(ctx context.Context, in *pb.UpdateProductReq) (*pb.UpdateProductResp, error) {
	id := in.GetId()
	updateData := &models.Product{
		Id:            int(id),
		Name:          in.GetName(),
		Description:   in.GetDescription(),
		Price:         in.GetPrice(),
		StockQuantity: in.GetStockQuantity(),
	}

	if err := g.repo.Update(updateData); err != nil {
		log.Println("handler proto:update err")
		return nil, err
	}

	return &pb.UpdateProductResp{
		Success: true,
		Message: "Продукт обновлен",
	}, nil
}

func (g *GetServiceHandler) DeleteProduct(ctx context.Context, in *pb.DeleteProductReq) (*pb.DeleteProductResp, error) {
	id := in.GetId()
	if err := g.repo.Delete(uint(id)); err != nil {
		log.Println("handler proto:delete err")
		return nil, err
	}

	return &pb.DeleteProductResp{Success: true}, nil
}
