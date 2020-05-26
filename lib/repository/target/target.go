package target

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/repository"
	"github.com/EurosportDigital/global-transcoding-platform/model"
	"github.com/EurosportDigital/global-transcoding-platform/model/gormmodel"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides CRUD operations on model.Target types.
type Repository interface {
	// Get retrieves the model with the specified ID.
	Get(id int) (*model.Target, error)

	// GetMany retrieves all objects with the specified IDs.
	GetMany(ids []int) ([]*model.Target, error)

	// Create adds the specified model.Target to the database.
	Create(target *model.Target) error

	// Update updates an existing record in the database.
	Update(target *model.Target) error

	// Delete removes the model.Target with the specified ID.
	Delete(id int) error

	// All retrieves all model.Targets within the database.
	All() ([]*model.Target, error)
}

type gormRepository struct {
	db *gorm.DB
}

// New constructs a new instance of the target repository.
func New(db *gorm.DB) *gormRepository {
	return &gormRepository{db}
}

func (targetRepo *gormRepository) Get(id int) (*model.Target, error) {
	target := gormmodel.Target{}
	err := repository.EvaluateError(targetRepo.db.First(&target, id).Error)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get target %v", id)
	}

	return gormmodel.ToTarget(&target), nil
}

func (targetRepo *gormRepository) GetMany(ids []int) ([]*model.Target, error) {
	var gormTargets []*gormmodel.Target
	err := repository.EvaluateError(targetRepo.db.Find(&gormTargets, "id IN (?)", ids).Error)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get targets %v", ids)
	}

	targets := make([]*model.Target, len(gormTargets))
	for i := range gormTargets {
		targets[i] = gormmodel.ToTarget(gormTargets[i])
	}
	return targets, nil
}

func (targetRepo *gormRepository) Create(target *model.Target) error {
	gormTarget := gormmodel.ToGormTarget(target)
	result := targetRepo.db.Create(gormTarget)
	if result.Error != nil {
		return errors.Wrapf(result.Error, "unable to create target %v", target)
	}
	*target = *gormmodel.ToTarget(gormTarget)
	return nil
}

func (targetRepo *gormRepository) Update(target *model.Target) error {
	gormTarget := gormmodel.ToGormTarget(target)
	result := targetRepo.db.Model(&gormTarget).Update(gormTarget)
	if result.Error != nil {
		return errors.Wrapf(result.Error, "unable to update target %v", target)
	}

	return nil
}

func (targetRepo *gormRepository) Delete(id int) error {
	result := targetRepo.db.Delete(&gormmodel.Target{ID: id})
	if result.Error != nil {
		return errors.Wrapf(result.Error, "unable to delete target %v", id)
	}
	if result.RowsAffected == 0 {
		return errors.Wrapf(repository.ErrEntityNotFound, "did not find target %v", id)
	}

	return nil
}

func (targetRepo *gormRepository) All() ([]*model.Target, error) {
	var gormTargets []*gormmodel.Target
	err := targetRepo.db.Find(&gormTargets).Error
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve all targets")
	}

	targets := make([]*model.Target, len(gormTargets))
	for i := range gormTargets {
		targets[i] = gormmodel.ToTarget(gormTargets[i])
	}
	return targets, nil
}
