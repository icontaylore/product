package repository

import (
	"fmt"
	"gorm.io/gorm"
	"product/internal/infrastructure/models"
)

type RepoInterface interface {
	Create(product *models.Product) error
	Read(id uint) (*models.Product, error)
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

func (p *ProductRepository) Read(id uint) (*models.Product, error) {
	var prod models.Product         // Структура для хранения найденного продукта
	result := p.db.First(&prod, id) // Получаем первый продукт с указанным ID

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("repo: продукт с id %d не найден", id)
		}
		return nil, fmt.Errorf("repo: ошибка при получении продукта с id %d: %v", id, result.Error)
	}

	return &prod, nil // Возвращаем указатель на продукт, если он найден
}

func (p *ProductRepository) Update(product *models.Product) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		var prod models.Product
		result := tx.First(&prod, product.Id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return fmt.Errorf("продукт с id %d не найден", product.Id) // Исправлено на product.Id
			}
			return fmt.Errorf("не удалось найти продукт с id %d: %v", product.Id, result.Error)
		}

		// Обновляем продукт
		result = tx.Save(product)
		if result.Error != nil {
			return fmt.Errorf("не удалось обновить продукт: %v", result.Error) // Добавлен правильный формат
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
		return fmt.Errorf("repo:не удалось удалить продукт: %w", deleteResult.Error)
	}

	return nil
}
