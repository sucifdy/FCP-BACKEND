package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
	"errors"
	"fmt"
)

type CategoryRepository interface {
	Store(Category *model.Category) error         // Stores a new category
	Update(id int, category model.Category) error // Updates an existing category by ID
	Delete(id int) error                          // Deletes a category by ID
	GetByID(id int) (*model.Category, error)      // Retrieves a category by its ID
	GetList() ([]model.Category, error)           // Retrieves all categories
}

type categoryRepository struct {
	filebasedDb *filebased.Data // File-based database connection
}

// NewCategoryRepo initializes and returns a new CategoryRepository instance
func NewCategoryRepo(filebasedDb *filebased.Data) *categoryRepository {
	return &categoryRepository{filebasedDb}
}

// Store adds a new category to the database
func (c *categoryRepository) Store(category *model.Category) error {
	return c.filebasedDb.StoreCategory(*category)
}

// Update modifies an existing category in the database based on its ID
func (c *categoryRepository) Update(id int, category model.Category) error {
	// Check if the category with the given ID exists
	existingCategory, err := c.filebasedDb.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// Update the existing category's data
	existingCategory.Name = category.Name
	return c.filebasedDb.UpdateCategory(id, *existingCategory)
}

// Delete removes a category from the database based on its ID
func (c *categoryRepository) Delete(id int) error {
	// Delete category by ID
	return c.filebasedDb.DeleteCategory(id)
}

// GetByID retrieves a category by its ID from the database
func (c *categoryRepository) GetByID(id int) (*model.Category, error) {
	return c.filebasedDb.GetCategoryByID(id)
}

// GetList retrieves all categories from the database
func (c *categoryRepository) GetList() ([]model.Category, error) {
	// Fetch all categories
	categories, err := c.filebasedDb.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("error retrieving categories: %v", err)
	}
	return categories, nil
}
