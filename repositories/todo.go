package repositories

import (
	"context"
	"errors"

	"github.com/joaoh82/buildingapi/entities"
	"github.com/joaoh82/buildingapi/interfaces"
	"github.com/joaoh82/buildingapi/utils"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func (repository *todoRepository) Find(ctx context.Context) ([]entities.Todo, error) {
	_, span := utils.StartSpan(ctx)
	defer span.End()

	var todos = []entities.Todo{}
	result := repository.db.Find(&todos)
	if result.Error != nil {
		return todos, entities.ErrorInternal(result.Error)
	}
	return todos, nil
}

func (repository *todoRepository) FindByID(ctx context.Context, id uint) (entities.Todo, error) {
	_, span := utils.StartSpan(ctx)
	defer span.End()

	var todo = entities.Todo{}
	result := repository.db.First(&todo, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return todo, entities.ErrorNotFound(result.Error)
		}
		return todo, entities.ErrorInternal(result.Error)
	}
	return todo, nil
}

func (repository *todoRepository) Create(ctx context.Context, todo *entities.Todo) error {
	_, span := utils.StartSpan(ctx)
	defer span.End()

	result := repository.db.Create(todo)
	if result.Error != nil {
		return entities.ErrorInternal(result.Error)
	}
	return nil
}

func (repository *todoRepository) Update(ctx context.Context, id uint, updateData entities.Todo) error {
	_, span := utils.StartSpan(ctx)
	defer span.End()

	result := repository.db.Model(entities.Todo{
		Model: gorm.Model{
			ID: id,
		},
	}).Updates(updateData)
	if result.Error != nil {
		return entities.ErrorInternal(result.Error)
	}
	return nil
}

func (repository *todoRepository) Delete(ctx context.Context, id uint) error {
	_, span := utils.StartSpan(ctx)
	defer span.End()

	result := repository.db.Delete(&entities.Todo{
		Model: gorm.Model{
			ID: id,
		},
	})
	if result.Error != nil {
		return entities.ErrorInternal(result.Error)
	}
	return nil
}

func NewTodoRepository(db *gorm.DB) interfaces.TodoRepository {
	return &todoRepository{
		db: db,
	}
}
