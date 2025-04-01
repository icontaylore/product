package repository

import (
	"fmt"
	"gorm.io/gorm"
	"product/internal/infrastructure/models"
)

type RepoInterface interface {
	Create(product *models.Product) error
	Read(id uint) error
	Update(product *models.Product) error
	Delete(id uint) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewRepoInterface(db *gorm.DB) RepoInterface {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) Create(product *models.Product) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		result := p.db.Create(product)
		if result.Error != nil {
			return fmt.Errorf("repo:не удалось создать продукт")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductRepository) Read(id uint) error {
	var prod *models.Product
	result := p.db.First(&prod, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("repo:продукт с таким id не найден")
		}
		return fmt.Errorf("repo:невозможно вытащить id")
	}

	return nil
}

func (p *ProductRepository) Update(product *models.Product) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		result := p.db.Save(product)
		if result.Error != nil {
			return fmt.Errorf("repo:не удалось обновить продукт")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
func (p *ProductRepository) Delete(id uint) error {
	var productDel *models.Product

	result := p.db.First(&productDel, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("repo:продукт с таким id не найден и не удален")
		}
		return fmt.Errorf("repo:не удалось найти продукт с таким id")
	}

	deleteResult := p.db.Delete(productDel)
	if deleteResult.Error != nil {
		return fmt.Errorf("не удалось удалить продукт: %w", deleteResult.Error)
	}

	return nil
}
