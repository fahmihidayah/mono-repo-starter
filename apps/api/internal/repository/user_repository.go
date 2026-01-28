package repository

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

// IUserRepository defines the interface for user repository operations
type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetAllWithFilter(ctx context.Context, name, email string, limit, offset int) ([]domain.User, error)
	Count(ctx context.Context) (int64, error)
	CountWithFilter(ctx context.Context, name, email string) (int64, error)
	FindByResetPasswordToken(ctx context.Context, token string) (*domain.User, error)
	FindByEmailAndVerificationCode(ctx context.Context, email, code string) (*domain.User, error)
	UpdateLoginAttempts(ctx context.Context, id string, attempts int, lockUntil int64) error
	ResetPassword(ctx context.Context, id string, hashedPassword string) error
	GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]domain.User, error)
	CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error)
}

// UserRepositoryImpl implements IUserRepository
type UserRepositoryImpl struct {
	db *gorm.DB
}

// UserRepositoryProvider creates a new instance of UserRepositoryImpl
func UserRepositoryProvider(db *gorm.DB) IUserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create inserts a new user into the database
func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by their ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user's information
func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Delete removes a user from the database
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// DeleteAll deletes multiple users by their IDs
func (r *UserRepositoryImpl) DeleteAll(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no users found with the provided IDs")
	}
	return nil
}

// GetAll retrieves a paginated list of users
func (r *UserRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Count returns the total number of users
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&count).Error
	return count, err
}

// GetAllWithFilter retrieves a paginated list of users with optional name/email filtering
func (r *UserRepositoryImpl) GetAllWithFilter(ctx context.Context, name, email string, limit, offset int) ([]domain.User, error) {
	var users []domain.User
	query := r.db.WithContext(ctx).Model(&domain.User{})

	// Apply dynamic filters only if values are provided
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// CountWithFilter returns the total number of users matching the filter criteria
func (r *UserRepositoryImpl) CountWithFilter(ctx context.Context, name, email string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.User{})

	// Apply dynamic filters only if values are provided
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) FindByResetPasswordToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByVerificationCode retrieves a user by their verification code
func (r *UserRepositoryImpl) FindByEmailAndVerificationCode(ctx context.Context, email, code string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("verification_code = ?", code).Where("email = ? ", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateLoginAttempts updates the login attempts and lock status for a user
func (r *UserRepositoryImpl) UpdateLoginAttempts(ctx context.Context, id string, attempts int, lockUntil int64) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"login_attempts": attempts,
			"lock_until":     lockUntil,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// ResetPassword updates the user's password
func (r *UserRepositoryImpl) ResetPassword(ctx context.Context, id string, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"hashed_password":             hashedPassword,
			"reset_password_token":        "",
			"reset_password_token_expiry": 0,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetWithQuery retrieves users with query parameters
func (r *UserRepositoryImpl) GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]domain.User, error) {
	var data []domain.User
	tx := r.db.WithContext(ctx).Limit(queryParams.Limit).Offset(queryParams.Offset)

	if len(queryParams.Sort) >= 2 {
		tx = tx.Order(queryParams.Sort[0] + " " + queryParams.Sort[1])
	}

	for key, value := range queryParams.Filter {
		if key == "ids" || key == "id" {
			if ids, ok := queryParams.GetFilterIDs(); ok {
				tx = tx.Where("id IN ?", ids)
			}
		} else {
			if strValue, ok := value.(string); ok {
				tx = tx.Where(key+" LIKE ?", "%"+strValue+"%")
			}
		}
	}

	err := tx.Find(&data).Error

	return data, err
}

// CountByQuery returns the count of users matching the query parameters
func (r *UserRepositoryImpl) CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.User{})

	for key, value := range queryParams.Filter {
		if key == "ids" || key == "id" {
			if ids, ok := queryParams.GetFilterIDs(); ok {
				tx = tx.Where("id IN ?", ids)
			}
		} else {
			if strValue, ok := value.(string); ok {
				tx = tx.Where(key+" LIKE ?", "%"+strValue+"%")
			}
		}
	}

	err := tx.Count(&count).Error
	return count, err
}
