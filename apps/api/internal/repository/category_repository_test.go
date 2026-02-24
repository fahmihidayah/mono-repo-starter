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

// CategoryRepositoryTestSuite defines the test suite for CategoryRepository
type CategoryRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository ICategoryRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *CategoryRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&domain.Category{}, &domain.Post{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = CategoryRepositoryProvider(db)
}

// SetupTest runs before each test
func (suite *CategoryRepositoryTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM post_categories")
	suite.db.Exec("DELETE FROM categories")
	suite.db.Exec("DELETE FROM posts")
}

// TearDownSuite runs once after all tests in the suite
func (suite *CategoryRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Helper function to create a test category
func (suite *CategoryRepositoryTestSuite) createTestCategory(id, slug, name string) *domain.Category {
	category := &domain.Category{
		ID:    id,
		Slug:  slug,
		Title: name,
	}
	suite.db.Create(category)
	return category
}

// TestCreate tests the Create method
func (suite *CategoryRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	category := &domain.Category{
		ID:    "cat-1",
		Slug:  "technology",
		Title: "Technology",
	}

	err := suite.repository.Create(ctx, category)

	assert.NoError(suite.T(), err)

	var savedCategory domain.Category
	suite.db.First(&savedCategory, "id = ?", category.ID)
	assert.Equal(suite.T(), category.ID, savedCategory.ID)
	assert.Equal(suite.T(), category.Slug, savedCategory.Slug)
	assert.Equal(suite.T(), category.Title, savedCategory.Title)
}

// TestCreate_DuplicateSlug tests creating a category with duplicate slug
func (suite *CategoryRepositoryTestSuite) TestCreate_DuplicateSlug() {
	ctx := context.Background()
	cat1 := suite.createTestCategory("cat-1", "technology", "Technology")

	cat2 := &domain.Category{
		ID:    "cat-2",
		Slug:  cat1.Slug,
		Title: "Tech",
	}

	err := suite.repository.Create(ctx, cat2)
	assert.Error(suite.T(), err)
}

// TestGetByID tests the GetByID method
func (suite *CategoryRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	expectedCategory := suite.createTestCategory("cat-1", "technology", "Technology")

	category, err := suite.repository.GetByID(ctx, expectedCategory.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), category)
	assert.Equal(suite.T(), expectedCategory.ID, category.ID)
	assert.Equal(suite.T(), expectedCategory.Slug, category.Slug)
	assert.Equal(suite.T(), expectedCategory.Title, category.Title)
}

// TestGetByID_NotFound tests getting a non-existent category
func (suite *CategoryRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	category, err := suite.repository.GetByID(ctx, "non-existent-id")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), category)
	assert.Equal(suite.T(), "category not found", err.Error())
}

func (suite *CategoryRepositoryTestSuite) TestGetAllIds() {
	ctx := context.Background()
	expectedCategory := suite.createTestCategory("cat-1", "technology", "Technology")

	category := suite.repository.GetAllIds(ctx, []string{expectedCategory.ID})

	assert.Equal(suite.T(), expectedCategory.ID, category[0].ID)
}

func (suite *CategoryRepositoryTestSuite) TestGetAllIds_NotFound() {
	ctx := context.Background()
	category := suite.repository.GetAllIds(ctx, []string{"non-existent-id"})

	assert.Equal(suite.T(), category, []*domain.Category{})
}

// TestFindBySlug tests the FindBySlug method
func (suite *CategoryRepositoryTestSuite) TestFindBySlug() {
	ctx := context.Background()
	expectedCategory := suite.createTestCategory("cat-1", "technology", "Technology")

	category, err := suite.repository.FindBySlug(ctx, expectedCategory.Slug)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), category)
	assert.Equal(suite.T(), expectedCategory.ID, category.ID)
	assert.Equal(suite.T(), expectedCategory.Slug, category.Slug)
}

// TestFindBySlug_NotFound tests finding a category with non-existent slug
func (suite *CategoryRepositoryTestSuite) TestFindBySlug_NotFound() {
	ctx := context.Background()
	category, err := suite.repository.FindBySlug(ctx, "non-existent-slug")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), category)
	assert.Equal(suite.T(), "category not found", err.Error())
}

// TestGetAll tests the GetAll method
func (suite *CategoryRepositoryTestSuite) TestGetAll() {
	ctx := context.Background()
	suite.createTestCategory("cat-1", "technology", "Technology")
	suite.createTestCategory("cat-2", "sports", "Sports")
	suite.createTestCategory("cat-3", "health", "Health")

	categories, err := suite.repository.GetAll(ctx, 10, 0)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), categories, 3)
}

// TestGetAll_WithPagination tests GetAll with limit and offset
func (suite *CategoryRepositoryTestSuite) TestGetAll_WithPagination() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Category{
			ID:    "cat-" + strconv.Itoa(i),
			Slug:  "slug-" + strconv.Itoa(i),
			Title: "Category " + strconv.Itoa(i),
		})
	}

	categories, err := suite.repository.GetAll(ctx, 2, 2)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), categories, 2)
}

// TestGetAll_Empty tests GetAll when no categories exist
func (suite *CategoryRepositoryTestSuite) TestGetAll_Empty() {
	ctx := context.Background()
	categories, err := suite.repository.GetAll(ctx, 10, 0)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), categories)
}

func (suite *CategoryRepositoryTestSuite) TestGetAllWithQueryParameter() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Category{
			ID:    "cat-" + strconv.Itoa(i),
			Slug:  "slug-" + strconv.Itoa(i),
			Title: "Category " + strconv.Itoa(i),
		})
	}
	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
	}

	categories, err := suite.repository.GetWithQuery(ctx, queryParameter)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), categories, 5)
}

func (suite *CategoryRepositoryTestSuite) TestGetAllWithQueryParameterFilter() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Category{
			ID:    "cat-" + strconv.Itoa(i),
			Slug:  "slug-" + strconv.Itoa(i),
			Title: "Category " + strconv.Itoa(i),
		})
	}
	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
		Filter: map[string]interface{}{
			"name": "1",
		},
	}

	categories, err := suite.repository.GetWithQuery(ctx, queryParameter)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), categories, 1)
}

func (suite *CategoryRepositoryTestSuite) TestCountByQuery() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Category{
			ID:    "cat-" + strconv.Itoa(i),
			Slug:  "slug-" + strconv.Itoa(i),
			Title: "Category " + strconv.Itoa(i),
		})
	}
	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
	}

	count, err := suite.repository.CountByQuery(ctx, queryParameter)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), count)
}

func (suite *CategoryRepositoryTestSuite) TestCountByQueryFilter() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Category{
			ID:    "cat-" + strconv.Itoa(i),
			Slug:  "slug-" + strconv.Itoa(i),
			Title: "Category " + strconv.Itoa(i),
		})
	}
	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
		Filter: map[string]interface{}{
			"name": "1",
		},
	}

	count, err := suite.repository.CountByQuery(ctx, queryParameter)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), count)
}

// TestUpdate tests the Update method
func (suite *CategoryRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	category := suite.createTestCategory("cat-1", "original-slug", "Original Name")

	category.Title = "Updated Name"
	category.Slug = "updated-slug"

	err := suite.repository.Update(ctx, category)

	assert.NoError(suite.T(), err)

	updatedCategory, _ := suite.repository.GetByID(ctx, category.ID)
	assert.Equal(suite.T(), "Updated Name", updatedCategory.Title)
	assert.Equal(suite.T(), "updated-slug", updatedCategory.Slug)
}

// TestDelete tests the Delete method
func (suite *CategoryRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	category := suite.createTestCategory("cat-1", "technology", "Technology")

	err := suite.repository.Delete(ctx, category.ID)

	assert.NoError(suite.T(), err)

	deletedCategory, err := suite.repository.GetByID(ctx, category.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), deletedCategory)
}

// TestDelete_NotFound tests deleting a non-existent category
func (suite *CategoryRepositoryTestSuite) TestDelete_NotFound() {
	ctx := context.Background()
	err := suite.repository.Delete(ctx, "non-existent-id")

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "category not found", err.Error())
}

// TestDeleteAll tests the DeleteAll method
func (suite *CategoryRepositoryTestSuite) TestDeleteAll() {
	ctx := context.Background()
	cat1 := suite.createTestCategory("cat-1", "tech", "Technology")
	cat2 := suite.createTestCategory("cat-2", "sports", "Sports")
	cat3 := suite.createTestCategory("cat-3", "health", "Health")

	ids := []string{cat1.ID, cat2.ID}
	err := suite.repository.DeleteAll(ctx, ids)

	assert.NoError(suite.T(), err)

	categories, _ := suite.repository.GetAll(ctx, 10, 0)
	assert.Len(suite.T(), categories, 1)
	assert.Equal(suite.T(), cat3.ID, categories[0].ID)
}

// TestDeleteAll_NoCategoriesFound tests DeleteAll with non-existent IDs
func (suite *CategoryRepositoryTestSuite) TestDeleteAll_NoCategoriesFound() {
	ctx := context.Background()
	ids := []string{"non-existent-1", "non-existent-2"}
	err := suite.repository.DeleteAll(ctx, ids)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "no categories found with the provided IDs", err.Error())
}

// TestCount tests the Count method
func (suite *CategoryRepositoryTestSuite) TestCount() {
	ctx := context.Background()
	suite.createTestCategory("cat-1", "tech", "Technology")
	suite.createTestCategory("cat-2", "sports", "Sports")
	suite.createTestCategory("cat-3", "health", "Health")

	count, err := suite.repository.Count(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)
}

// TestCount_Empty tests Count when no categories exist
func (suite *CategoryRepositoryTestSuite) TestCount_Empty() {
	ctx := context.Background()
	count, err := suite.repository.Count(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

// TestSuite entry point
func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

func (suite *CategoryRepositoryTestSuite) TestGetAllByQueryParameterIds() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Category{
			ID:    "cat-" + strconv.Itoa(i),
			Slug:  "slug-" + strconv.Itoa(i),
			Title: "Category " + strconv.Itoa(i),
		})
	}
	queryParameter := &utils.QueryParams{
		Limit:  10,
		Offset: 0,
		Sort:   []string{"name", "asc"},
		Filter: map[string]interface{}{
			"ids": []string{"cat-1", "cat-2"},
		},
	}

	count, err := suite.repository.CountByQuery(ctx, queryParameter)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count)
}
