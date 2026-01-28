package repository

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

type IPostRepository interface {
	// Define methods for post repository
	FindBySlug(ctx context.Context, slug string) (*domain.Post, error)
	Create(ctx context.Context, data *domain.Post) error
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Post, error)
	Update(ctx context.Context, data *domain.Post) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	Count(ctx context.Context) (int64, error)
	GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Post, error)
	CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error)
}

type PostRepositoryImpl struct {
	db *gorm.DB
}

func PostRepositoryProvider(db *gorm.DB) IPostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, data *domain.Post) error {
	// Use Association to properly handle many-to-many relationship
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	// If there are categories, use Association API to ensure join table is populated
	if len(data.Categories) > 0 {
		if err := r.db.WithContext(ctx).Model(data).Association("Categories").Replace(data.Categories); err != nil {
			return err
		}
	}

	return nil
}

func (r *PostRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	var data domain.Post
	if err := r.db.WithContext(ctx).Preload("Categories").First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *PostRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	var data []domain.Post
	err := r.db.WithContext(ctx).Preload("Categories").Limit(limit).Offset(offset).Find(&data).Error
	return data, err
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Post{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *PostRepositoryImpl) DeleteAll(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no posts found with the provided IDs")
	}
	return nil
}

func (r *PostRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Post{}).Count(&count).Error
	return count, err
}

func (r *PostRepositoryImpl) Update(ctx context.Context, data *domain.Post) error {
	// Update post fields
	result := r.db.WithContext(ctx).Model(&domain.Post{}).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"title":   data.Title,
		"content": data.Content,
		"slug":    data.Slug,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found or no changes made")
	}

	// Update categories association if provided
	if data.Categories != nil {
		if err := r.db.WithContext(ctx).Model(data).Association("Categories").Replace(data.Categories); err != nil {
			return err
		}
	}

	return nil
}

func (r *PostRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	var data domain.Post
	if err := r.db.WithContext(ctx).Preload("Categories").Where("slug = ?", slug).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &data, nil
}

// GetWithQuery retrieves posts with query parameters
func (r *PostRepositoryImpl) GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Post, error) {
	var data []domain.Post
	tx := r.db.WithContext(ctx).Preload("Categories").Limit(queryParams.Limit).Offset(queryParams.Offset)

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

// CountByQuery returns the count of posts matching the query parameters
func (r *PostRepositoryImpl) CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Post{})

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
