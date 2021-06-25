package services

import (
	"context"

	"github.com/joaoh82/buildingapi/entities"
	"github.com/joaoh82/buildingapi/interfaces"
	"github.com/joaoh82/buildingapi/utils"
)

type todoService struct {
	todoRepository interfaces.TodoRepository
}

func NewTodoService(todoRepository interfaces.TodoRepository) interfaces.TodoService {
	return &todoService{
		todoRepository: todoRepository,
	}
}

func (service *todoService) Find(ctx context.Context) ([]entities.Todo, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return service.todoRepository.Find(ctx)
}

func (service *todoService) FindByID(ctx context.Context, id uint) (entities.Todo, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return service.todoRepository.FindByID(ctx, id)
}

func (service *todoService) Create(ctx context.Context, name string) (entities.Todo, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	todo := entities.NewTodo(name)
	err := service.todoRepository.Create(ctx, &todo)
	return todo, err
}

func (service *todoService) Update(ctx context.Context, id uint, name string) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	updateData := entities.Todo{Name: name}
	return service.todoRepository.Update(ctx, id, updateData)
}

func (service *todoService) Delete(ctx context.Context, id uint) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return service.todoRepository.Delete(ctx, id)
}
