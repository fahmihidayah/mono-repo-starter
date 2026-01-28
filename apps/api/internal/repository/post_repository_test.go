package repository

import (
	"context"
	"strconv"
	"testing"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// PostRepositoryTestSuite defines the test suite for PostRepository
type PostRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository IPostRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *PostRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&domain.Post{}, &domain.Category{}, &domain.User{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = PostRepositoryProvider(db)
}

// SetupTest runs before each test
func (suite *PostRepositoryTestSuite) SetupTest() {
	// Clean up the database before each test
	suite.db.Exec("DELETE FROM post_categories")
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM categories")
	suite.db.Exec("DELETE FROM users")
}

// TearDownSuite runs once after all tests in the suite
func (suite *PostRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Helper function to create a test user
func (suite *PostRepositoryTestSuite) createTestUser(id, email string) *domain.User {
	user := &domain.User{
		ID:    id,
		Email: email,
		Name:  "Test User",
	}
	suite.db.Create(user)
	return user
}

// Helper function to create a test post
func (suite *PostRepositoryTestSuite) createTestPost(id, slug, title, userID string) *domain.Post {
	post := &domain.Post{
		ID:      id,
		Slug:    slug,
		Title:   title,
		Content: "Test content",
		UserID:  userID,
	}
	suite.db.Create(post)
	return post
}

// TestCreate tests the Create method
func (suite *PostRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	// Create a test user first
	user := suite.createTestUser("user-1", "test@example.com")

	// Test data
	post := &domain.Post{
		ID:      "post-1",
		Slug:    "test-post",
		Title:   "Test Post",
		Content: "This is a test post content",
		UserID:  user.ID,
	}

	// Execute
	err := suite.repository.Create(ctx, post)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the post was created
	var savedPost domain.Post
	suite.db.First(&savedPost, "id = ?", post.ID)
	assert.Equal(suite.T(), post.ID, savedPost.ID)
	assert.Equal(suite.T(), post.Slug, savedPost.Slug)
	assert.Equal(suite.T(), post.Title, savedPost.Title)
	assert.Equal(suite.T(), post.Content, savedPost.Content)
	assert.Equal(suite.T(), post.UserID, savedPost.UserID)
}

// TestCreate_DuplicateSlug tests creating a post with duplicate slug
func (suite *PostRepositoryTestSuite) TestCreate_DuplicateSlug() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create first post
	post1 := suite.createTestPost("post-1", "test-post", "Test Post 1", user.ID)

	// Try to create post with same slug
	post2 := &domain.Post{
		ID:      "post-2",
		Slug:    post1.Slug, // Same slug
		Title:   "Test Post 2",
		Content: "Different content",
		UserID:  user.ID,
	}

	// Execute
	err := suite.repository.Create(ctx, post2)

	// Assert - should return error due to unique constraint
	assert.Error(suite.T(), err)
}

// TestGetByID tests the GetByID method
func (suite *PostRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create test post
	expectedPost := suite.createTestPost("post-1", "test-post", "Test Post", user.ID)

	// Execute
	post, err := suite.repository.GetByID(ctx, expectedPost.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), post)
	assert.Equal(suite.T(), expectedPost.ID, post.ID)
	assert.Equal(suite.T(), expectedPost.Slug, post.Slug)
	assert.Equal(suite.T(), expectedPost.Title, post.Title)
	assert.Equal(suite.T(), expectedPost.Content, post.Content)
	assert.Equal(suite.T(), expectedPost.UserID, post.UserID)
}

// TestGetByID_NotFound tests getting a non-existent post
func (suite *PostRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	// Execute
	post, err := suite.repository.GetByID(ctx, "non-existent-id")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), post)
	assert.Equal(suite.T(), "post not found", err.Error())
}

// TestFindBySlug tests the FindBySlug method
func (suite *PostRepositoryTestSuite) TestFindBySlug() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create test post
	expectedPost := suite.createTestPost("post-1", "test-post-slug", "Test Post", user.ID)

	// Execute
	post, err := suite.repository.FindBySlug(ctx, expectedPost.Slug)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), post)
	assert.Equal(suite.T(), expectedPost.ID, post.ID)
	assert.Equal(suite.T(), expectedPost.Slug, post.Slug)
	assert.Equal(suite.T(), expectedPost.Title, post.Title)
}

// TestFindBySlug_NotFound tests finding a post with non-existent slug
func (suite *PostRepositoryTestSuite) TestFindBySlug_NotFound() {
	ctx := context.Background()
	// Execute
	post, err := suite.repository.FindBySlug(ctx, "non-existent-slug")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), post)
	assert.Equal(suite.T(), "post not found", err.Error())
}

// TestGetAll tests the GetAll method
func (suite *PostRepositoryTestSuite) TestGetAll() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create multiple test posts
	suite.createTestPost("post-1", "post-1", "Post 1", user.ID)
	suite.createTestPost("post-2", "post-2", "Post 2", user.ID)
	suite.createTestPost("post-3", "post-3", "Post 3", user.ID)

	// Execute
	posts, err := suite.repository.GetAll(ctx, 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 3)
}

// TestGetAll_WithPagination tests GetAll with limit and offset
func (suite *PostRepositoryTestSuite) TestGetAll_WithPagination() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create 5 test posts
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Post{
			ID:      "post-" + string(rune(i)),
			Slug:    "post-slug-" + string(rune(i)),
			Title:   "Post " + string(rune(i)),
			Content: "Content " + string(rune(i)),
			UserID:  user.ID,
		})
	}

	// Execute - Get 2 posts, skip first 2
	posts, err := suite.repository.GetAll(ctx, 2, 2)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 2)
}

// TestGetAll_Empty tests GetAll when no posts exist
func (suite *PostRepositoryTestSuite) TestGetAll_Empty() {
	ctx := context.Background()
	// Execute
	posts, err := suite.repository.GetAll(ctx, 10, 0)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), posts)
}

// TestUpdate tests the Update method
func (suite *PostRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create test post
	post := suite.createTestPost("post-1", "original-slug", "Original Title", user.ID)

	// Modify the post
	post.Title = "Updated Title"
	post.Content = "Updated content"
	post.Slug = "updated-slug"

	// Execute
	err := suite.repository.Update(ctx, post)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the update
	updatedPost, _ := suite.repository.GetByID(ctx, post.ID)
	assert.Equal(suite.T(), "Updated Title", updatedPost.Title)
	assert.Equal(suite.T(), "Updated content", updatedPost.Content)
	assert.Equal(suite.T(), "updated-slug", updatedPost.Slug)
}

// TestUpdate_NonExistent tests updating a non-existent post
func (suite *PostRepositoryTestSuite) TestUpdate_NonExistent() {
	ctx := context.Background()
	// Test data - post that doesn't exist
	post := &domain.Post{
		ID:      "non-existent-id",
		Slug:    "test-slug",
		Title:   "Test Title",
		Content: "Test Content",
		UserID:  "user-1",
	}

	// Execute
	err := suite.repository.Update(ctx, post)

	// Assert - should return error
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "post not found or no changes made", err.Error())
}

// TestDelete tests the Delete method
func (suite *PostRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create test post
	post := suite.createTestPost("post-1", "test-post", "Test Post", user.ID)

	// Execute
	err := suite.repository.Delete(ctx, post.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify the post is deleted
	deletedPost, err := suite.repository.GetByID(ctx, post.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), deletedPost)
}

// TestDelete_NotFound tests deleting a non-existent post
func (suite *PostRepositoryTestSuite) TestDelete_NotFound() {
	ctx := context.Background()
	// Execute
	err := suite.repository.Delete(ctx, "non-existent-id")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "post not found", err.Error())
}

// TestDeleteAll tests the DeleteAll method
func (suite *PostRepositoryTestSuite) TestDeleteAll() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create multiple test posts
	post1 := suite.createTestPost("post-1", "post-1", "Post 1", user.ID)
	post2 := suite.createTestPost("post-2", "post-2", "Post 2", user.ID)
	post3 := suite.createTestPost("post-3", "post-3", "Post 3", user.ID)

	// Execute - Delete first two posts
	ids := []string{post1.ID, post2.ID}
	err := suite.repository.DeleteAll(ctx, ids)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify only post3 remains
	posts, _ := suite.repository.GetAll(ctx, 10, 0)
	assert.Len(suite.T(), posts, 1)
	assert.Equal(suite.T(), post3.ID, posts[0].ID)
}

// TestDeleteAll_NoPostsFound tests DeleteAll with non-existent IDs
func (suite *PostRepositoryTestSuite) TestDeleteAll_NoPostsFound() {
	ctx := context.Background()
	// Execute
	ids := []string{"non-existent-1", "non-existent-2"}
	err := suite.repository.DeleteAll(ctx, ids)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "no posts found with the provided IDs", err.Error())
}

// TestDeleteAll_PartialMatch tests DeleteAll with mix of existing and non-existing IDs
func (suite *PostRepositoryTestSuite) TestDeleteAll_PartialMatch() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create test post
	post := suite.createTestPost("post-1", "post-1", "Post 1", user.ID)

	// Execute - Mix of existing and non-existing IDs
	ids := []string{post.ID, "non-existent-id"}
	err := suite.repository.DeleteAll(ctx, ids)

	// Assert - Should succeed because at least one ID exists
	assert.NoError(suite.T(), err)

	// Verify the existing post was deleted
	posts, _ := suite.repository.GetAll(ctx, 10, 0)
	assert.Len(suite.T(), posts, 0)
}

// TestCount tests the Count method
func (suite *PostRepositoryTestSuite) TestCount() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create multiple test posts
	suite.createTestPost("post-1", "post-1", "Post 1", user.ID)
	suite.createTestPost("post-2", "post-2", "Post 2", user.ID)
	suite.createTestPost("post-3", "post-3", "Post 3", user.ID)

	// Execute
	count, err := suite.repository.Count(ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)
}

// TestCount_Empty tests Count when no posts exist
func (suite *PostRepositoryTestSuite) TestCount_Empty() {
	ctx := context.Background()
	// Execute
	count, err := suite.repository.Count(ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

// TestPostWithCategories tests creating a post with categories
func (suite *PostRepositoryTestSuite) TestPostWithCategories() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create test categories
	category1 := &domain.Category{ID: "cat-1", Slug: "tech", Name: "Technology"}
	category2 := &domain.Category{ID: "cat-2", Slug: "programming", Name: "Programming"}
	suite.db.Create(category1)
	suite.db.Create(category2)

	// Create post with categories
	post := &domain.Post{
		ID:         "post-1",
		Slug:       "test-post",
		Title:      "Test Post",
		Content:    "Content",
		UserID:     user.ID,
		Categories: []*domain.Category{category1, category2},
	}

	// Execute
	err := suite.repository.Create(ctx, post)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify post and categories relationship
	var savedPost domain.Post
	suite.db.Preload("Categories").First(&savedPost, "id = ?", post.ID)
	assert.Len(suite.T(), savedPost.Categories, 2)
}

// TestSuite entry point
func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}

func (suite *PostRepositoryTestSuite) TestGetAllByQueryParams_Success() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create multiple test posts
	suite.createTestPost("post-1", "post-1", "Post 1", user.ID)
	suite.createTestPost("post-2", "post-2", "Post 2", user.ID)
	suite.createTestPost("post-3", "post-3", "Post 3", user.ID)

	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"title", "asc"},
		Filter: map[string]interface{}{
			"ids": []string{"post-1", "post-2"},
		},
	}

	posts, err := suite.repository.GetWithQuery(ctx, queryParameter)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 2)
}

func (suite *PostRepositoryTestSuite) TestGetAllByQueryParams_Page10Success() {
	ctx := context.Background()
	// Create a test user
	user := suite.createTestUser("user-1", "test@example.com")

	// Create multiple test posts
	for i := 0; i < 20; i++ {
		suite.createTestPost("post-"+strconv.Itoa(i), "post-"+strconv.Itoa(i), "Post "+strconv.Itoa(i), user.ID)
	}

	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"title", "asc"},
		Filter: map[string]interface{}{},
	}

	numberOfPost, err := suite.repository.CountByQuery(ctx, queryParameter)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int(20), int(numberOfPost))
}
