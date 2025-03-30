package repository

import "gorm.io/gorm"

type RepoInterface interface {
	Create() error
	Read() error
	Update() error
	Delete() error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProduct() RepoInterface {

}

func (p *ProductRepository) Create() {

}

func (p *ProductRepository) Read()  {

}

func (p *ProductRepository) ()  {

}
func (p *ProductRepository) () {}