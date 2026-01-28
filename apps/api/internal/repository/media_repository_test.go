package repository

import (
	"context"
	"testing"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MediaRepositoryTestSuite defines the test suite for MediaRepository
type MediaRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository IMediaRepository
}

// SetupSuite runs once before all tests in the suite
func (suite *MediaRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&domain.Media{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = MediaRepositoryProvider(db)
}

// SetupTest runs before each test
func (suite *MediaRepositoryTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM media")
}

// TearDownSuite runs once after all tests in the suite
func (suite *MediaRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Helper function to create a test media
func (suite *MediaRepositoryTestSuite) createTestMedia(id, path, fileName string) *domain.Media {
	media := &domain.Media{
		ID:       id,
		Alt:      "Test image",
		Url:      "http://localhost:8080" + path,
		Path:     path,
		FileName: fileName,
		MimeType: "image/jpeg",
		FileSize: 1024,
		Width:    800,
		Height:   600,
	}
	suite.db.Create(media)
	return media
}

// TestCreate tests the Create method
func (suite *MediaRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	media := &domain.Media{
		ID:       "media-1",
		Alt:      "Test image",
		Url:      "http://localhost:8080/uploads/test.jpg",
		Path:     "/uploads/test.jpg",
		FileName: "test.jpg",
		MimeType: "image/jpeg",
		FileSize: 2048,
		Width:    1024,
		Height:   768,
	}

	err := suite.repository.Create(ctx, media)

	assert.NoError(suite.T(), err)

	var savedMedia domain.Media
	suite.db.First(&savedMedia, "id = ?", media.ID)
	assert.Equal(suite.T(), media.ID, savedMedia.ID)
	assert.Equal(suite.T(), media.Alt, savedMedia.Alt)
	assert.Equal(suite.T(), media.Path, savedMedia.Path)
	assert.Equal(suite.T(), media.FileName, savedMedia.FileName)
	assert.Equal(suite.T(), media.FileSize, savedMedia.FileSize)
}

// TestGetByID tests the GetByID method
func (suite *MediaRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	expectedMedia := suite.createTestMedia("media-1", "/uploads/test.jpg", "test.jpg")

	media, err := suite.repository.GetByID(ctx, expectedMedia.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), media)
	assert.Equal(suite.T(), expectedMedia.ID, media.ID)
	assert.Equal(suite.T(), expectedMedia.Path, media.Path)
	assert.Equal(suite.T(), expectedMedia.FileName, media.FileName)
}

// TestGetByID_NotFound tests getting a non-existent media
func (suite *MediaRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	media, err := suite.repository.GetByID(ctx, "non-existent-id")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), media)
	assert.Equal(suite.T(), "media not found", err.Error())
}

// TestGetByPath tests the GetByPath method
func (suite *MediaRepositoryTestSuite) TestGetByPath() {
	ctx := context.Background()
	expectedMedia := suite.createTestMedia("media-1", "/uploads/test.jpg", "test.jpg")

	media, err := suite.repository.GetByPath(ctx, expectedMedia.Path)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), media)
	assert.Equal(suite.T(), expectedMedia.ID, media.ID)
	assert.Equal(suite.T(), expectedMedia.Path, media.Path)
}

// TestGetByPath_NotFound tests getting media with non-existent path
func (suite *MediaRepositoryTestSuite) TestGetByPath_NotFound() {
	ctx := context.Background()
	media, err := suite.repository.GetByPath(ctx, "/non-existent/path.jpg")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), media)
	assert.Equal(suite.T(), "media not found", err.Error())
}

// TestGetAll tests the GetAll method
func (suite *MediaRepositoryTestSuite) TestGetAll() {
	ctx := context.Background()
	suite.createTestMedia("media-1", "/uploads/test1.jpg", "test1.jpg")
	suite.createTestMedia("media-2", "/uploads/test2.jpg", "test2.jpg")
	suite.createTestMedia("media-3", "/uploads/test3.jpg", "test3.jpg")

	mediaList, err := suite.repository.GetAll(ctx, 10, 0)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), mediaList, 3)
}

// TestGetAll_WithPagination tests GetAll with limit and offset
func (suite *MediaRepositoryTestSuite) TestGetAll_WithPagination() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		suite.db.Create(&domain.Media{
			ID:       "media-" + string(rune(i)),
			Path:     "/uploads/test" + string(rune(i)) + ".jpg",
			FileName: "test" + string(rune(i)) + ".jpg",
			MimeType: "image/jpeg",
			FileSize: 1024,
		})
	}

	mediaList, err := suite.repository.GetAll(ctx, 2, 2)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), mediaList, 2)
}

// TestGetAll_Empty tests GetAll when no media exist
func (suite *MediaRepositoryTestSuite) TestGetAll_Empty() {
	ctx := context.Background()
	mediaList, err := suite.repository.GetAll(ctx, 10, 0)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), mediaList)
}

// TestUpdate tests the Update method
func (suite *MediaRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	media := suite.createTestMedia("media-1", "/uploads/original.jpg", "original.jpg")

	media.Alt = "Updated alt text"
	media.FileName = "updated.jpg"

	err := suite.repository.Update(ctx, media)

	assert.NoError(suite.T(), err)

	updatedMedia, _ := suite.repository.GetByID(ctx, media.ID)
	assert.Equal(suite.T(), "Updated alt text", updatedMedia.Alt)
	assert.Equal(suite.T(), "updated.jpg", updatedMedia.FileName)
}

// TestDelete tests the Delete method
func (suite *MediaRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	media := suite.createTestMedia("media-1", "/uploads/test.jpg", "test.jpg")

	err := suite.repository.Delete(ctx, media.ID)

	assert.NoError(suite.T(), err)

	deletedMedia, err := suite.repository.GetByID(ctx, media.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), deletedMedia)
}

// TestDelete_NotFound tests deleting a non-existent media
func (suite *MediaRepositoryTestSuite) TestDelete_NotFound() {
	ctx := context.Background()
	err := suite.repository.Delete(ctx, "non-existent-id")

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "media not found", err.Error())
}

// TestDeleteAll tests the DeleteAll method
func (suite *MediaRepositoryTestSuite) TestDeleteAll() {
	ctx := context.Background()
	media1 := suite.createTestMedia("media-1", "/uploads/test1.jpg", "test1.jpg")
	media2 := suite.createTestMedia("media-2", "/uploads/test2.jpg", "test2.jpg")
	media3 := suite.createTestMedia("media-3", "/uploads/test3.jpg", "test3.jpg")

	ids := []string{media1.ID, media2.ID}
	err := suite.repository.DeleteAll(ctx, ids)

	assert.NoError(suite.T(), err)

	mediaList, _ := suite.repository.GetAll(ctx, 10, 0)
	assert.Len(suite.T(), mediaList, 1)
	assert.Equal(suite.T(), media3.ID, mediaList[0].ID)
}

// TestDeleteAll_NoMediaFound tests DeleteAll with non-existent IDs
func (suite *MediaRepositoryTestSuite) TestDeleteAll_NoMediaFound() {
	ctx := context.Background()
	ids := []string{"non-existent-1", "non-existent-2"}
	err := suite.repository.DeleteAll(ctx, ids)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "no media found with the provided IDs", err.Error())
}

// TestCount tests the Count method
func (suite *MediaRepositoryTestSuite) TestCount() {
	ctx := context.Background()
	suite.createTestMedia("media-1", "/uploads/test1.jpg", "test1.jpg")
	suite.createTestMedia("media-2", "/uploads/test2.jpg", "test2.jpg")
	suite.createTestMedia("media-3", "/uploads/test3.jpg", "test3.jpg")

	count, err := suite.repository.Count(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count)
}

// TestCount_Empty tests Count when no media exist
func (suite *MediaRepositoryTestSuite) TestCount_Empty() {
	ctx := context.Background()
	count, err := suite.repository.Count(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
}

// TestSuite entry point
func TestMediaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MediaRepositoryTestSuite))
}
