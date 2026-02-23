package services

import (
	"context"
	"errors"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type TemplateTaskService interface {
	CRUDService[models.TemplateTask]
	GetAllTemplateTasks(ctx context.Context) ([]models.TemplateTask, error)
	GetTemplateTaskByID(ctx context.Context, id int64) (*models.TemplateTask, error)
	CreateTemplateTask(ctx context.Context, templateTask *models.TemplateTask) error
	UpdateTemplateTask(ctx context.Context, templateTask *models.TemplateTask) error
	DeleteTemplateTask(ctx context.Context, id int64) error
	ListByRoleID(ctx context.Context, roleID int64) ([]models.TemplateTask, error)
}

type TemplateTaskAttributeService interface {
	CRUDService[models.TemplateTaskAttribute]
	GetAllTemplateTaskAttributes(ctx context.Context) ([]models.TemplateTaskAttribute, error)
	GetTemplateTaskAttributeByID(ctx context.Context, id int64) (*models.TemplateTaskAttribute, error)
	CreateTemplateTaskAttribute(ctx context.Context, attr *models.TemplateTaskAttribute) error
	UpdateTemplateTaskAttribute(ctx context.Context, attr *models.TemplateTaskAttribute) error
	DeleteTemplateTaskAttribute(ctx context.Context, id int64) error
	ListByTaskID(ctx context.Context, taskID int64) ([]models.TemplateTaskAttribute, error)
}

type templateTaskService struct {
	*baseService[models.TemplateTask]
	repo repository.TemplateTaskRepository
}

type templateTaskAttributeService struct {
	*baseService[models.TemplateTaskAttribute]
	repo repository.TemplateTaskAttributeRepository
}

func NewTemplateTaskService(repo repository.TemplateTaskRepository) TemplateTaskService {
	return &templateTaskService{
		baseService: newBaseService[models.TemplateTask](repo),
		repo:        repo,
	}
}

func NewTemplateTaskAttributeService(repo repository.TemplateTaskAttributeRepository) TemplateTaskAttributeService {
	return &templateTaskAttributeService{
		baseService: newBaseService[models.TemplateTaskAttribute](repo),
		repo:        repo,
	}
}

func (s *templateTaskService) GetAllTemplateTasks(ctx context.Context) ([]models.TemplateTask, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *templateTaskService) GetTemplateTaskByID(ctx context.Context, id int64) (*models.TemplateTask, error) {
	if err := validatePositiveID("template task", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *templateTaskService) CreateTemplateTask(ctx context.Context, templateTask *models.TemplateTask) error {
	if templateTask == nil {
		return errors.New("template task payload cannot be empty")
	}
	if err := validateName("template task", templateTask.TetaName, maxTemplateName); err != nil {
		return err
	}
	if err := validatePositiveID(defaultLookupRoleTag, templateTask.TetaRoleID); err != nil {
		return err
	}
	return s.repo.Create(ctx, templateTask)
}

func (s *templateTaskService) UpdateTemplateTask(ctx context.Context, templateTask *models.TemplateTask) error {
	if templateTask == nil {
		return errors.New("template task payload cannot be empty")
	}
	if err := validatePositiveID("template task", templateTask.TetaID); err != nil {
		return err
	}
	if err := validateName("template task", templateTask.TetaName, maxTemplateName); err != nil {
		return err
	}
	if err := validatePositiveID(defaultLookupRoleTag, templateTask.TetaRoleID); err != nil {
		return err
	}
	return s.repo.Update(ctx, templateTask.TetaID, map[string]interface{}{
		"teta_name":    strings.TrimSpace(templateTask.TetaName),
		"teta_role_id": templateTask.TetaRoleID,
	})
}

func (s *templateTaskService) DeleteTemplateTask(ctx context.Context, id int64) error {
	if err := validatePositiveID("template task", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *templateTaskService) ListByRoleID(ctx context.Context, roleID int64) ([]models.TemplateTask, error) {
	if err := validatePositiveID(defaultLookupRoleTag, roleID); err != nil {
		return nil, err
	}
	return s.repo.ListByRoleID(ctx, roleID)
}

func (s *templateTaskAttributeService) GetAllTemplateTaskAttributes(ctx context.Context) ([]models.TemplateTaskAttribute, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *templateTaskAttributeService) GetTemplateTaskAttributeByID(ctx context.Context, id int64) (*models.TemplateTaskAttribute, error) {
	if err := validatePositiveID("template task attribute", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *templateTaskAttributeService) CreateTemplateTaskAttribute(ctx context.Context, attr *models.TemplateTaskAttribute) error {
	if attr == nil {
		return errors.New("template task attribute payload cannot be empty")
	}
	if err := validateName("template task attribute", attr.TetatName, maxTemplateName); err != nil {
		return err
	}
	if err := validatePositiveID("template task", attr.TetatTetaID); err != nil {
		return err
	}
	return s.repo.Create(ctx, attr)
}

func (s *templateTaskAttributeService) UpdateTemplateTaskAttribute(ctx context.Context, attr *models.TemplateTaskAttribute) error {
	if attr == nil {
		return errors.New("template task attribute payload cannot be empty")
	}
	if err := validatePositiveID("template task attribute", attr.TetatID); err != nil {
		return err
	}
	if err := validateName("template task attribute", attr.TetatName, maxTemplateName); err != nil {
		return err
	}
	if err := validatePositiveID("template task", attr.TetatTetaID); err != nil {
		return err
	}
	return s.repo.Update(ctx, attr.TetatID, map[string]interface{}{
		"tetat_name":    strings.TrimSpace(attr.TetatName),
		"tetat_teta_id": attr.TetatTetaID,
	})
}

func (s *templateTaskAttributeService) DeleteTemplateTaskAttribute(ctx context.Context, id int64) error {
	if err := validatePositiveID("template task attribute", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *templateTaskAttributeService) ListByTaskID(ctx context.Context, taskID int64) ([]models.TemplateTaskAttribute, error) {
	if err := validatePositiveID("template task", taskID); err != nil {
		return nil, err
	}
	return s.repo.ListByTaskID(ctx, taskID)
}
