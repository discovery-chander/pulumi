package profile

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/repository"
	"github.com/EurosportDigital/global-transcoding-platform/model"
	"github.com/EurosportDigital/global-transcoding-platform/model/gormmodel"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Repository provides CRUD operations on model.Profile types.
type Repository interface {
	// Get retrieves the model.Profile with the specified ID.
	Get(id int) (*model.Profile, error)

	// GetMany retrieves all model.Profiles with the specified IDs.
	GetMany(ids []int) ([]*model.Profile, error)

	// GetByName retrieves the model.Profile with the specified name.
	GetByName(name string) (*model.Profile, error)

	// GetManyByName retrieves all model.Profiles with the specified names.
	GetManyByName(names []string) ([]*model.Profile, error)

	// Create adds the specified model.Profile to the database.
	Create(profile *model.Profile) error

	// Update updates an existing record in the database.
	Update(profile *model.Profile) error

	// Delete removes the model.Profile with the specified ID.
	Delete(id int) error

	// All retrieves all model.Profile within the database.
	All() ([]*model.Profile, error)
}

type gormRepository struct {
	db *gorm.DB
}

// New constructs a new instance of the profile repository.
func New(db *gorm.DB) *gormRepository {
	return &gormRepository{db}
}

func (profileRepo *gormRepository) Get(id int) (*model.Profile, error) {
	return profileRepo.getFirstProfile(id)
}

func (profileRepo *gormRepository) GetMany(ids []int) ([]*model.Profile, error) {
	return profileRepo.getManyProfiles("id IN (?)", ids)
}

func (profileRepo *gormRepository) GetByName(name string) (*model.Profile, error) {
	return profileRepo.getFirstProfile("name = ?", name)
}

func (profileRepo *gormRepository) GetManyByName(names []string) ([]*model.Profile, error) {
	return profileRepo.getManyProfiles("name IN (?)", names)
}

func (profileRepo *gormRepository) Create(profile *model.Profile) error {
	gormProfile := gormmodel.ToGormProfile(profile)
	if config := gormProfile.EncConfig; config != nil && config.ID != 0 {
		// The ID should only be specified if it's referring to an existing config.
		// If we do not set the EncConfig to nil then it will update that row, which we do not want on a Create.
		gormProfile.EncConfig = nil
	}
	result := profileRepo.db.Create(gormProfile)
	if result.Error != nil {
		return errors.Wrapf(result.Error, "unable to create profile %v", profile)
	}
	*profile = *gormmodel.ToProfile(gormProfile)
	return nil
}

func (profileRepo *gormRepository) Update(profile *model.Profile) error {
	gormProfile := gormmodel.ToGormProfile(profile)
	result := profileRepo.db.Model(&gormProfile).Update(gormProfile)
	if result.Error != nil {
		return errors.Wrapf(result.Error, "unable to update profile %v", profile)
	}

	return nil
}

func (profileRepo *gormRepository) Delete(id int) error {
	result := profileRepo.db.Delete(&gormmodel.Profile{ID: id})
	if result.Error != nil {
		return errors.Wrapf(result.Error, "unable to delete profile %v", id)
	}
	if result.RowsAffected == 0 {
		return errors.Wrapf(repository.ErrEntityNotFound, "did not find profile %v", id)
	}

	return nil
}

func (profileRepo *gormRepository) All() ([]*model.Profile, error) {
	var profiles []*model.Profile
	err := profileRepo.db.Find(&profiles).Error
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve all profiles")
	}
	return profiles, nil
}

func (profileRepo *gormRepository) getFirstProfile(where ...interface{}) (*model.Profile, error) {
	var profile gormmodel.Profile
	err := repository.EvaluateError(profileRepo.db.Preload("EncConfig.Encoder").First(&profile, where...).Error)
	if err != nil {
		return nil, errors.Wrapf(err, "did not find profile where %v", where)
	}

	return gormmodel.ToProfile(&profile), nil
}

func (profileRepo *gormRepository) getManyProfiles(where ...interface{}) ([]*model.Profile, error) {
	var profiles []*gormmodel.Profile
	err := repository.EvaluateError(profileRepo.db.Preload("EncConfig.Encoder").Find(&profiles, where...).Error)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find profiles where %v", where)
	}

	modelProfiles := make([]*model.Profile, len(profiles))
	for i := range profiles {
		modelProfiles[i] = gormmodel.ToProfile(profiles[i])
	}

	return modelProfiles, nil
}
