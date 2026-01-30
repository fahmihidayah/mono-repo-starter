package service

import (
	"context"
	"errors"

	request "github.com/fahmihidayah/go-api-orchestrator/internal/data/request/posts"
	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-playground/validator/v10"
)

type IPostService interface {
	Create(ctx context.Context, post *request.CreatePostRequest, userID string) (*domain.Post, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Post, error)
	Update(ctx context.Context, id string, post *request.CreatePostRequest, userID string) (*domain.Post, error)
	Delete(ctx context.Context, id string, userID string) error
	DeleteAll(ctx context.Context, ids []string, userID string) error
	// React Admin specific methods
	GetAllReactAdmin(ctx context.Context, limit, offset int, sortField, sortOrder string, filters map[string]interface{}) ([]domain.Post, int64, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Post, error)
	UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}, userID string) ([]string, error)
	GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Post, *utils.PaginateInfo, error)
}

type PostServiceImpl struct {
	postRepository     repository.IPostRepository
	categoryRepository repository.ICategoryRepository
	validator          *validator.Validate
}

func PostServiceProvider(
	postRepository repository.IPostRepository,
	categoryRepository repository.ICategoryRepository,
) IPostService {
	return &PostServiceImpl{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
		validator:          validator.New(),
	}
}

func (s *PostServiceImpl) Create(ctx context.Context, post *request.CreatePostRequest, userID string) (*domain.Post, error) {
	// Validate request
	if err := s.validator.Struct(post); err != nil {
		return nil, err
	}

	var categories []*domain.Category = make([]*domain.Category, 0)

	if len(post.CategoryIDs) > 0 {
		result := s.categoryRepository.GetAllIds(ctx, post.CategoryIDs)

		// Validate that all requested categories exist
		if len(result) != len(post.CategoryIDs) {
			return nil, errors.New("one or more category IDs are invalid")
		}

		categories = result
	}

	// Map request to domain model
	data := &domain.Post{
		ID:         utils.GenerateUUID(),
		Slug:       utils.GenerateSlug(post.Title),
		Title:      post.Title,
		Content:    post.Content,
		UserID:     userID,
		Categories: categories,
	}

	// Save to repository
	if err := s.postRepository.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PostServiceImpl) GetAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get posts from repository
	posts, err := s.postRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.postRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *PostServiceImpl) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	return s.postRepository.GetByID(ctx, id)
}

func (s *PostServiceImpl) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	return s.postRepository.FindBySlug(ctx, slug)
}

func (s *PostServiceImpl) Update(ctx context.Context, id string, post *request.CreatePostRequest, userID string) (*domain.Post, error) {
	// Validate request
	if err := s.validator.Struct(post); err != nil {
		return nil, err
	}

	// Check if post exists and belongs to user
	existingPost, err := s.postRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if existingPost.UserID != userID {
		return nil, errors.New("unauthorized: you can only update your own posts")
	}

	// Update fields
	existingPost.Title = post.Title
	existingPost.Content = post.Content
	existingPost.Slug = utils.GenerateSlug(post.Title)

	// Update categories if provided
	if len(post.CategoryIDs) > 0 {
		result, err := s.categoryRepository.GetWithQuery(ctx, &utils.QueryParams{
			Filter: map[string]interface{}{
				"id": post.CategoryIDs,
			},
		})

		if err != nil {
			return nil, errors.New("invalid categories")
		}
		existingPost.Categories = result
	} else {
		// If no categories provided, clear existing categories
		existingPost.Categories = []*domain.Category{}
	}

	// Save to repository
	if err := s.postRepository.Update(ctx, existingPost); err != nil {
		return nil, err
	}

	return existingPost, nil
}

func (s *PostServiceImpl) Delete(ctx context.Context, id string, userID string) error {
	// Check if post exists and belongs to user
	existingPost, err := s.postRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if existingPost.UserID != userID {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	return s.postRepository.Delete(ctx, id)
}

func (s *PostServiceImpl) DeleteAll(ctx context.Context, ids []string, userID string) error {
	// Validate IDs array
	if len(ids) == 0 {
		return errors.New("IDs array cannot be empty")
	}

	// For bulk delete, we should verify ownership of all posts
	// For simplicity, we'll just call repository
	// In production, you'd want to check each post's ownership
	return s.postRepository.DeleteAll(ctx, ids)
}

// GetAllReactAdmin retrieves posts with React Admin parameters
func (s *PostServiceImpl) GetAllReactAdmin(ctx context.Context, limit, offset int, sortField, sortOrder string, filters map[string]interface{}) ([]domain.Post, int64, error) {
	// For now, use the basic GetAll method
	// In a more advanced implementation, you would apply sortField, sortOrder, and filters
	posts, err := s.postRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// GetByIDs retrieves multiple posts by their IDs
func (s *PostServiceImpl) GetByIDs(ctx context.Context, ids []string) ([]domain.Post, error) {
	if len(ids) == 0 {
		return []domain.Post{}, nil
	}

	posts := make([]domain.Post, 0, len(ids))
	for _, id := range ids {
		post, err := s.postRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip posts that don't exist
		}
		posts = append(posts, *post)
	}

	return posts, nil
}

// UpdateMany updates multiple posts with the same data
func (s *PostServiceImpl) UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}, userID string) ([]string, error) {
	if len(ids) == 0 {
		return nil, errors.New("IDs array cannot be empty")
	}

	updatedIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		post, err := s.postRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip posts that don't exist
		}

		// Check ownership
		if post.UserID != userID {
			continue // Skip posts user doesn't own
		}

		// Apply updates
		if title, ok := updates["title"].(string); ok && title != "" {
			post.Title = title
			post.Slug = utils.GenerateSlug(title)
		}
		if content, ok := updates["content"].(string); ok {
			post.Content = content
		}

		// Save changes
		if err := s.postRepository.Update(ctx, post); err != nil {
			continue // Skip if update fails
		}

		updatedIDs = append(updatedIDs, id)
	}

	if len(updatedIDs) == 0 {
		return nil, errors.New("no posts were updated")
	}

	return updatedIDs, nil
}

// GetWithQueryParams retrieves posts with React Admin parameters
func (s *PostServiceImpl) GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Post, *utils.PaginateInfo, error) {
	count, err := s.postRepository.CountByQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	paginateInfo := queryParams.ToPaginateInfo(count)
	posts, err := s.postRepository.GetWithQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	return posts, paginateInfo, nil
}
