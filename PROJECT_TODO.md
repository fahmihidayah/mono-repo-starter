# Project TODO - API Enhancement Roadmap

**Project:** Mono Starter API
**Created:** March 6, 2026
**Status:** Planning Phase
**Priority:** High - Critical issues affecting production readiness

---

## Overview

This TODO list addresses the critical issues identified in the Project Structure Analysis. The goal is to improve the API from a **B grade (82/100)** to an **A grade (95+/100)** by addressing inconsistencies in request/response DTOs, package naming, and cross-cutting concerns.

**Estimated Total Time:** 3-4 weeks
**Risk Level:** Low to Medium
**Breaking Changes:** None (backward compatible)

---

## Phase 1: Package Naming Standardization (Priority: CRITICAL)
**Duration:** 2-3 hours
**Risk:** Low
**Dependencies:** None

### Task 1.1: Update Posts Package Name
- [ ] Open `internal/data/request/posts/post_request.go`
- [ ] Change `package request` to `package posts`
- [ ] Verify no other files in posts directory need updates
- [ ] Run `go build ./...` to check compilation
- [ ] Run `go test ./internal/data/request/posts/...` for tests

**Commands:**
```bash
cd apps/api
sed -i 's/^package request$/package posts/' internal/data/request/posts/post_request.go
go build ./...
go test ./internal/data/request/posts/...
```

**Expected Files Changed:** 1
**Breaking Changes:** None (import paths remain the same)

---

### Task 1.2: Update Users Package Name
- [ ] List all files in `internal/data/request/users/` directory
- [ ] Update `user_request.go` - change `package request` to `package users`
- [ ] Update `login_user.go` - change `package request` to `package users`
- [ ] Update `verify_user.go` - change `package request` to `package users`
- [ ] Update `reset_password_user.go` - change `package request` to `package users`
- [ ] Update `filter_user.go` - change `package request` to `package users`
- [ ] Run `go build ./...` to check compilation
- [ ] Run `go test ./internal/data/request/users/...` for tests

**Commands:**
```bash
cd apps/api
find internal/data/request/users -name "*.go" -exec sed -i 's/^package request$/package users/' {} \;
go build ./...
go test ./internal/data/request/users/...
```

**Expected Files Changed:** 5
**Breaking Changes:** None

---

### Task 1.3: Verify Package Naming Consistency
- [ ] Verify posts package: `grep "^package" internal/data/request/posts/post_request.go`
- [ ] Verify users package: `grep "^package" internal/data/request/users/*.go`
- [ ] Verify categories package: `grep "^package" internal/data/request/categories/*.go`
- [ ] Verify media package: `grep "^package" internal/data/request/media/*.go`
- [ ] Ensure all show feature-specific package names

**Success Criteria:**
```
✅ internal/data/request/posts/      → package posts
✅ internal/data/request/users/      → package users
✅ internal/data/request/categories/ → package categories
✅ internal/data/request/media/      → package media
```

---

### Task 1.4: Full Build & Test Verification
- [ ] Run full build: `go build ./...`
- [ ] Run all tests: `go test ./...`
- [ ] Check for import errors in controllers
- [ ] Check for import errors in services
- [ ] Verify Swagger docs still generate correctly
- [ ] Test one endpoint per feature (manual smoke test)

**Commands:**
```bash
cd apps/api
go build ./...
go test ./... -v
go run cmd/api/main.go  # Start server
# Test endpoints manually or with curl
```

---

### Task 1.5: Documentation & Commit
- [ ] Update any internal documentation referencing old package names
- [ ] Update README if it mentions request package structure
- [ ] Create git commit with descriptive message
- [ ] Push to feature branch for review

**Commit Message:**
```
refactor: standardize request package names to feature-specific

- Changed package request → package posts (posts feature)
- Changed package request → package users (users feature)
- Categories and media already correct
- No breaking changes - import paths unchanged
- Improves code clarity and Go convention adherence

Resolves: Package naming inconsistency identified in structure analysis
```

---

## Phase 2: Response DTO Implementation (Priority: CRITICAL)
**Duration:** 1 week
**Risk:** Medium
**Dependencies:** None (can run parallel to Phase 1)

### Task 2.1: Create Response Directory Structure
- [ ] Create `internal/data/response/common/` directory
- [ ] Create `internal/data/response/users/` directory
- [ ] Create `internal/data/response/posts/` directory
- [ ] Create `internal/data/response/categories/` directory
- [ ] Create `internal/data/response/media/` directory

**Commands:**
```bash
cd apps/api/internal/data/response
mkdir -p common users posts categories media
```

---

### Task 2.2: Move Common Response Files
- [ ] Move `error_response.go` to `common/error_response.go`
- [ ] Move `pagination_response.go` to `common/pagination_response.go`
- [ ] Move `response.go` to `common/web_response.go`
- [ ] Update package declarations to `package common`
- [ ] Update all imports across codebase to use new paths

**Commands:**
```bash
cd apps/api/internal/data/response
mv error_response.go common/
mv pagination_response.go common/
mv response.go common/web_response.go

# Update package declarations
sed -i 's/^package response$/package common/' common/*.go

# Find all files importing these
grep -r "internal/data/response\"" internal/
```

**Expected Files to Update:** ~10-15 files (controllers, services)

---

### Task 2.3: Create Posts Response DTOs
- [ ] Create `internal/data/response/posts/post_response.go`
- [ ] Define `PostResponse` struct with proper JSON tags
- [ ] Define `PostListResponse` struct for list operations
- [ ] Add transformation function `ToPostResponse(domain.Post) PostResponse`
- [ ] Add transformation function `ToPostListResponse([]domain.Post) []PostResponse`
- [ ] Include computed fields (e.g., excerpt, category names)

**File:** `internal/data/response/posts/post_response.go`
```go
package posts

import "github.com/fahmihidayah/go-api-orchestrator/internal/domain"

// PostResponse represents the API response for a post
type PostResponse struct {
    ID          string   `json:"id"`
    Slug        string   `json:"slug"`
    Title       string   `json:"title"`
    Content     string   `json:"content"`
    Excerpt     string   `json:"excerpt"`           // Computed field
    UserID      string   `json:"user_id"`
    Categories  []CategorySummary `json:"categories"` // Nested response
    CreatedAt   int64    `json:"created_at"`
    UpdatedAt   int64    `json:"updated_at"`
}

type CategorySummary struct {
    ID    string `json:"id"`
    Title string `json:"title"`
    Slug  string `json:"slug"`
}

// ToPostResponse transforms domain.Post to PostResponse
func ToPostResponse(post *domain.Post) *PostResponse {
    if post == nil {
        return nil
    }

    categories := make([]CategorySummary, 0, len(post.Categories))
    for _, cat := range post.Categories {
        categories = append(categories, CategorySummary{
            ID:    cat.ID,
            Title: cat.Title,
            Slug:  cat.Slug,
        })
    }

    // Generate excerpt (first 200 chars)
    excerpt := post.Content
    if len(excerpt) > 200 {
        excerpt = excerpt[:200] + "..."
    }

    return &PostResponse{
        ID:         post.ID,
        Slug:       post.Slug,
        Title:      post.Title,
        Content:    post.Content,
        Excerpt:    excerpt,
        UserID:     post.UserID,
        Categories: categories,
        CreatedAt:  post.CreatedAt,
        UpdatedAt:  post.UpdatedAt,
    }
}

// ToPostListResponse transforms slice of posts
func ToPostListResponse(posts []domain.Post) []*PostResponse {
    responses := make([]*PostResponse, 0, len(posts))
    for i := range posts {
        responses = append(responses, ToPostResponse(&posts[i]))
    }
    return responses
}
```

---

### Task 2.4: Create Categories Response DTOs
- [ ] Create `internal/data/response/categories/category_response.go`
- [ ] Define `CategoryResponse` struct
- [ ] Define `CategoryListResponse` struct
- [ ] Add transformation functions
- [ ] Include computed fields (e.g., post count)

**File:** `internal/data/response/categories/category_response.go`
```go
package categories

import "github.com/fahmihidayah/go-api-orchestrator/internal/domain"

// CategoryResponse represents the API response for a category
type CategoryResponse struct {
    ID        string `json:"id"`
    Slug      string `json:"slug"`
    Title     string `json:"title"`
    PostCount int    `json:"post_count,omitempty"` // Optional computed field
    CreatedAt int64  `json:"created_at"`
    UpdatedAt int64  `json:"updated_at"`
}

// ToCategoryResponse transforms domain.Category to CategoryResponse
func ToCategoryResponse(category *domain.Category) *CategoryResponse {
    if category == nil {
        return nil
    }

    return &CategoryResponse{
        ID:        category.ID,
        Slug:      category.Slug,
        Title:     category.Title,
        CreatedAt: category.CreatedAt,
        UpdatedAt: category.UpdatedAt,
    }
}

// ToCategoryListResponse transforms slice of categories
func ToCategoryListResponse(categories []domain.Category) []*CategoryResponse {
    responses := make([]*CategoryResponse, 0, len(categories))
    for i := range categories {
        responses = append(responses, ToCategoryResponse(&categories[i]))
    }
    return responses
}
```

---

### Task 2.5: Create Media Response DTOs
- [ ] Create `internal/data/response/media/media_response.go`
- [ ] Define `MediaResponse` struct
- [ ] Define `MediaListResponse` struct
- [ ] Add transformation functions
- [ ] Include computed fields (e.g., full URL, thumbnail URL)

**File:** `internal/data/response/media/media_response.go`
```go
package media

import (
    "fmt"
    "github.com/fahmihidayah/go-api-orchestrator/internal/domain"
)

// MediaResponse represents the API response for media
type MediaResponse struct {
    ID          string `json:"id"`
    FileName    string `json:"file_name"`
    FilePath    string `json:"file_path"`
    FileSize    int64  `json:"file_size"`
    MimeType    string `json:"mime_type"`
    Description string `json:"description"`
    URL         string `json:"url"`           // Computed: Full URL
    ThumbnailURL string `json:"thumbnail_url,omitempty"` // Computed: Thumbnail if image
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}

// ToMediaResponse transforms domain.Media to MediaResponse
func ToMediaResponse(media *domain.Media, baseURL string) *MediaResponse {
    if media == nil {
        return nil
    }

    // Construct full URL
    url := fmt.Sprintf("%s/%s", baseURL, media.FilePath)

    return &MediaResponse{
        ID:          media.ID,
        FileName:    media.FileName,
        FilePath:    media.FilePath,
        FileSize:    media.FileSize,
        MimeType:    media.MimeType,
        Description: media.Description,
        URL:         url,
        CreatedAt:   media.CreatedAt,
        UpdatedAt:   media.UpdatedAt,
    }
}

// ToMediaListResponse transforms slice of media
func ToMediaListResponse(mediaList []domain.Media, baseURL string) []*MediaResponse {
    responses := make([]*MediaResponse, 0, len(mediaList))
    for i := range mediaList {
        responses = append(responses, ToMediaResponse(&mediaList[i], baseURL))
    }
    return responses
}
```

---

### Task 2.6: Move & Enhance User Response DTO
- [ ] Move `user_response.go` to `internal/data/response/users/user_response.go`
- [ ] Update package declaration to `package users`
- [ ] Verify transformation functions exist
- [ ] Add `UserListResponse` if missing
- [ ] Update imports across codebase

---

### Task 2.7: Update Post Controller to Use Response DTOs
- [ ] Import `postResponse "github.com/.../response/posts"`
- [ ] Update `GetList` to transform domain.Post → PostResponse
- [ ] Update `GetOne` to transform domain.Post → PostResponse
- [ ] Update `Create` to transform domain.Post → PostResponse
- [ ] Update `Update` to transform domain.Post → PostResponse
- [ ] Update `Delete` to transform domain.Post → PostResponse
- [ ] Test all endpoints

**Example Change:**
```go
// Before
func (c *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
    post, err := c.postService.GetByID(r.Context(), id)
    if err != nil {
        c.SendNotFound(w, err.Error())
        return
    }
    c.SendOne(w, post)  // Returns domain.Post directly
}

// After
func (c *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
    post, err := c.postService.GetByID(r.Context(), id)
    if err != nil {
        c.SendNotFound(w, err.Error())
        return
    }
    response := postResponse.ToPostResponse(post)
    c.SendOne(w, response)  // Returns PostResponse DTO
}
```

---

### Task 2.8: Update Category Controller to Use Response DTOs
- [ ] Import `categoryResponse "github.com/.../response/categories"`
- [ ] Update all methods to use CategoryResponse
- [ ] Test all endpoints

---

### Task 2.9: Update Media Controller to Use Response DTOs
- [ ] Import `mediaResponse "github.com/.../response/media"`
- [ ] Update all methods to use MediaResponse
- [ ] Pass baseURL from config for URL generation
- [ ] Test all endpoints

---

### Task 2.10: Update Swagger Documentation
- [ ] Update all Swagger annotations to reference response DTOs
- [ ] Example: `@Success 200 {object} posts.PostResponse`
- [ ] Regenerate Swagger docs
- [ ] Verify API documentation is correct

**Commands:**
```bash
cd apps/api
# Regenerate swagger docs (adjust command based on your setup)
swag init -g cmd/api/main.go
```

---

### Task 2.11: Response DTO Testing & Verification
- [ ] Write unit tests for ToPostResponse transformation
- [ ] Write unit tests for ToCategoryResponse transformation
- [ ] Write unit tests for ToMediaResponse transformation
- [ ] Test all API endpoints manually
- [ ] Verify response structure matches expected format
- [ ] Check that computed fields work correctly

**Test File Example:** `internal/data/response/posts/post_response_test.go`
```go
package posts

import (
    "testing"
    "github.com/fahmihidayah/go-api-orchestrator/internal/domain"
)

func TestToPostResponse(t *testing.T) {
    post := &domain.Post{
        ID:      "123",
        Title:   "Test Post",
        Content: "This is test content",
        // ... other fields
    }

    response := ToPostResponse(post)

    if response.ID != post.ID {
        t.Errorf("Expected ID %s, got %s", post.ID, response.ID)
    }
    // ... more assertions
}
```

---

## Phase 3: Media Request DTO Standardization (Priority: HIGH)
**Duration:** 4-6 hours
**Risk:** Medium (involves multipart form handling)
**Dependencies:** Phase 1 complete

### Task 3.1: Create BaseMediaRequest
- [ ] Open `internal/data/request/media/upload_media.go`
- [ ] Define `BaseMediaRequest` struct
- [ ] Extract common fields (Description, etc.)
- [ ] Consider multipart form handling

**Code:**
```go
package media

import "mime/multipart"

// BaseMediaRequest contains common fields for media operations
type BaseMediaRequest struct {
    Description string `json:"description" form:"description" example:"Product image"`
}
```

---

### Task 3.2: Rename UploadMediaRequest to CreateMediaRequest
- [ ] Rename `UploadMediaRequest` to `CreateMediaRequest`
- [ ] Embed `BaseMediaRequest`
- [ ] Update all references in media controller
- [ ] Update all references in media service

**Code:**
```go
// CreateMediaRequest represents the request for creating/uploading media
type CreateMediaRequest struct {
    File *multipart.FileHeader `form:"file" validate:"required" binding:"required"`
    BaseMediaRequest          `form:",inline"`
}
```

---

### Task 3.3: Update UpdateMediaRequest
- [ ] Add ID field to UpdateMediaRequest
- [ ] Embed BaseMediaRequest
- [ ] Follow same pattern as other features

**Code:**
```go
// UpdateMediaRequest represents the request for updating media
type UpdateMediaRequest struct {
    ID string `json:"id" validate:"required,uuid4" binding:"required,uuid4"`
    BaseMediaRequest
}
```

---

### Task 3.4: Update Media Controller
- [ ] Replace `UploadMediaRequest` with `CreateMediaRequest` in Create method
- [ ] Update variable names and comments
- [ ] Test file upload functionality
- [ ] Verify multipart form handling still works

---

### Task 3.5: Update Media Service
- [ ] Update service interface if needed
- [ ] Update implementation to use CreateMediaRequest
- [ ] Test media creation
- [ ] Test media updates

---

### Task 3.6: Update Swagger Docs for Media
- [ ] Update @Param annotations to use CreateMediaRequest
- [ ] Update example values
- [ ] Regenerate Swagger documentation
- [ ] Verify media endpoints show correctly

---

## Phase 4: Service Layer Standardization (Priority: MEDIUM)
**Duration:** 1 week
**Risk:** Medium (changes business logic layer)
**Dependencies:** Phases 1-3 complete

### Task 4.1: Analyze Service Method Signatures
- [ ] Document all Update method signatures across services
- [ ] Document all Delete method signatures across services
- [ ] Identify inconsistencies
- [ ] Choose standard pattern to adopt

**Current State Analysis:**
```go
// UserService - ID in request
Update(ctx, *UpdateUserRequest) (*domain.User, error)

// CategoryService - ID as parameter
Update(ctx, id string, *CreateCategoryRequest) (*domain.Category, error)

// PostService - ID as parameter + userID
Update(ctx, id string, *CreatePostRequest, userID string) (*domain.Post, error)

// MediaService - ID as parameter
Update(ctx, id string, *UpdateMediaRequest) (*domain.Media, error)
```

---

### Task 4.2: Decision - Choose Standard Pattern
- [ ] Evaluate pros/cons of each approach
- [ ] Decide on single pattern for all services
- [ ] Document decision in ARCHITECTURE.md

**Recommended Pattern:**
```go
// Option 1: ID in Request (like UserService)
// Pros: Single parameter, clear ownership
// Cons: Requires setting ID from URL in controller

// Option 2: ID as Parameter (like Category/Post)
// Pros: Explicit, matches REST semantics
// Cons: More parameters, redundant if ID also in request

// Recommendation: Option 2 with UpdateRequest (not CreateRequest)
Update(ctx context.Context, id string, req *UpdateXRequest) (*domain.X, error)
```

---

### Task 4.3: Update Category Service Signature
- [ ] Change `Update(ctx, id, *CreateCategoryRequest)` to `Update(ctx, id, *UpdateCategoryRequest)`
- [ ] Update interface definition
- [ ] Update implementation
- [ ] Update all controller calls
- [ ] Run tests

---

### Task 4.4: Update Post Service Signature
- [ ] Change `Update(ctx, id, *CreatePostRequest, userID)` to `Update(ctx, id, *UpdatePostRequest, userID)`
- [ ] Update interface definition
- [ ] Update implementation
- [ ] Update all controller calls
- [ ] Maintain authorization (userID) parameter
- [ ] Run tests

---

### Task 4.5: Standardize Authorization Pattern
- [ ] Document which resources require user ownership
  - [ ] Posts: User-owned (requires userID)
  - [ ] Media: User-owned (requires userID)
  - [ ] Categories: Global resource (no userID)
  - [ ] Users: Self or admin (special handling)
- [ ] Ensure consistent ownership checks
- [ ] Consider middleware vs. service layer auth
- [ ] Document pattern in ARCHITECTURE.md

---

### Task 4.6: Update DeleteAll Methods
- [ ] Verify all DeleteAll methods use same signature: `(ctx, ids []string, userID?) error`
- [ ] Add authorization where needed
- [ ] Remove unused BulkDeleteRequest types
- [ ] Update tests

---

## Phase 5: Frontend Type Alignment (Priority: CRITICAL)
**Duration:** 3-5 days
**Risk:** Medium (affects all frontend features)
**Dependencies:** Phase 2 complete (response DTOs created)

### Task 5.1: Create Type Transformation Utilities
- [ ] Create `apps/frontend/app/lib/utils/api-transformer.ts`
- [ ] Implement snake_case → camelCase converter
- [ ] Implement Unix timestamp → Date converter
- [ ] Add type guards and validators

**File:** `apps/frontend/app/lib/utils/api-transformer.ts`
```typescript
// Transform API response (snake_case) to frontend types (camelCase)
export function transformApiResponse<T>(data: any): T {
  if (!data) return data;

  if (Array.isArray(data)) {
    return data.map(item => transformApiResponse(item)) as T;
  }

  if (typeof data === 'object') {
    const transformed: any = {};

    for (const [key, value] of Object.entries(data)) {
      // Convert snake_case to camelCase
      const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());

      // Convert timestamps
      if ((key === 'created_at' || key === 'updated_at') && typeof value === 'number') {
        transformed[camelKey] = new Date(value * 1000);
      } else if (typeof value === 'object') {
        transformed[camelKey] = transformApiResponse(value);
      } else {
        transformed[camelKey] = value;
      }
    }

    return transformed as T;
  }

  return data;
}

// Transform frontend data to API format (camelCase → snake_case)
export function transformToApiFormat<T>(data: any): T {
  if (!data) return data;

  if (Array.isArray(data)) {
    return data.map(item => transformToApiFormat(item)) as T;
  }

  if (typeof data === 'object') {
    const transformed: any = {};

    for (const [key, value] of Object.entries(data)) {
      // Convert camelCase to snake_case
      const snakeKey = key.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);

      if (typeof value === 'object' && value instanceof Date) {
        transformed[snakeKey] = Math.floor(value.getTime() / 1000);
      } else if (typeof value === 'object') {
        transformed[snakeKey] = transformToApiFormat(value);
      } else {
        transformed[snakeKey] = value;
      }
    }

    return transformed as T;
  }

  return data;
}
```

---

### Task 5.2: Standardize Post Types
- [ ] Open `apps/frontend/app/features/posts/types/index.ts`
- [ ] Change all fields to camelCase
- [ ] Change timestamp type from `number | null` to `Date`
- [ ] Update Zod schemas
- [ ] Add transformation in API client

**Before:**
```typescript
export type Post = {
  id: string;
  title: string;
  slug: string;
  content: string;
  user_id: string;           // ❌ snake_case
  created_at: Date | null;   // ❌ Could be inconsistent
  updated_at: Date | null;
};
```

**After:**
```typescript
export type Post = {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string;           // ✅ New computed field from backend
  userId: string;            // ✅ camelCase
  categories: CategorySummary[]; // ✅ Nested response
  createdAt: Date;           // ✅ Always Date, not null
  updatedAt: Date;
};

export type CategorySummary = {
  id: string;
  title: string;
  slug: string;
};
```

---

### Task 5.3: Standardize Media Types
- [ ] Open `apps/frontend/app/features/media/types/index.ts`
- [ ] Change all fields to camelCase
- [ ] Change `created_at: number` to `createdAt: Date`
- [ ] Change `updated_at: number` to `updatedAt: Date`
- [ ] Add new computed fields from backend (url, thumbnailUrl)

**Before:**
```typescript
export type Media = {
  id: string;
  file_name: string;        // ❌ snake_case
  created_at: number | null; // ❌ number type
  updated_at: number | null;
};
```

**After:**
```typescript
export type Media = {
  id: string;
  fileName: string;         // ✅ camelCase
  filePath: string;
  fileSize: number;
  mimeType: string;
  description: string;
  url: string;              // ✅ Computed from backend
  thumbnailUrl?: string;    // ✅ Optional computed field
  createdAt: Date;          // ✅ Date type
  updatedAt: Date;
};
```

---

### Task 5.4: Update Category Types (Already Correct)
- [ ] Verify `apps/frontend/app/features/category/types/index.ts`
- [ ] Ensure all fields are camelCase (should already be correct)
- [ ] Verify Date types
- [ ] Add any missing computed fields from backend

---

### Task 5.5: Update All API Clients to Use Transformers
- [ ] Update `features/posts/api/index.ts` to use transformApiResponse
- [ ] Update `features/media/api/index.ts` to use transformApiResponse
- [ ] Update `features/category/api/index.ts` to use transformApiResponse
- [ ] Update `features/users/api/users.ts` to use transformApiResponse

**Example:**
```typescript
// features/posts/api/index.ts
import { transformApiResponse } from "@/lib/utils/api-transformer";

export const postApiClient = new ApiClient<Post, PostFormData>({
  endpoint: "posts",
  transformResponse: (data) => transformApiResponse<Post>(data),
});
```

---

### Task 5.6: Update All Zod Schemas
- [ ] Update post schemas to match new types
- [ ] Update media schemas to match new types
- [ ] Update category schemas (verify)
- [ ] Update user schemas (verify)

---

### Task 5.7: Update All Components
- [ ] Find all components using old field names (snake_case)
- [ ] Update to use new field names (camelCase)
- [ ] Test all features thoroughly

**Commands to find usages:**
```bash
cd apps/frontend
grep -r "created_at" app/features/posts/
grep -r "updated_at" app/features/posts/
grep -r "user_id" app/features/posts/
grep -r "file_name" app/features/media/
```

---

### Task 5.8: Frontend Testing
- [ ] Test all post operations (list, create, update, delete)
- [ ] Test all media operations
- [ ] Test all category operations
- [ ] Test all user operations
- [ ] Verify dates display correctly
- [ ] Verify computed fields show correctly

---

## Phase 6: Missing Features (Priority: MEDIUM)
**Duration:** 3-5 days
**Risk:** Low
**Dependencies:** Phases 1-5 complete

### Task 6.1: Add Missing Frontend Hooks - Posts
- [ ] Create `apps/frontend/app/features/posts/hooks/usePostForm.ts`
- [ ] Create `apps/frontend/app/features/posts/hooks/usePostList.ts`
- [ ] Create `apps/frontend/app/features/posts/hooks/index.ts`
- [ ] Follow pattern from users/hooks
- [ ] Update components to use hooks

---

### Task 6.2: Add Missing Frontend Hooks - Categories
- [ ] Create `apps/frontend/app/features/category/hooks/useCategoryForm.ts`
- [ ] Create `apps/frontend/app/features/category/hooks/useCategoryList.ts`
- [ ] Create `apps/frontend/app/features/category/hooks/index.ts`
- [ ] Update components to use hooks

---

### Task 6.3: Extract Auth Feature
- [ ] Create `apps/frontend/app/features/auth/` directory
- [ ] Move auth-related code from users feature
- [ ] Create auth actions, API, components, types
- [ ] Update imports across codebase
- [ ] Test login/register/logout flows

---

### Task 6.4: Add Error Boundaries
- [ ] Create `apps/frontend/app/components/error/ErrorBoundary.tsx`
- [ ] Create `apps/frontend/app/components/error/ErrorFallback.tsx`
- [ ] Wrap route components with error boundaries
- [ ] Test error handling

---

## Phase 7: Documentation & Testing (Priority: HIGH)
**Duration:** 3-4 days
**Risk:** Low
**Dependencies:** All previous phases

### Task 7.1: Update API Documentation
- [ ] Regenerate Swagger documentation
- [ ] Document request/response DTO patterns
- [ ] Document authentication/authorization model
- [ ] Create API usage examples
- [ ] Document error response format

---

### Task 7.2: Create Architecture Documentation
- [ ] Create `ARCHITECTURE.md` file
- [ ] Document layered architecture
- [ ] Document request/response DTO pattern
- [ ] Document service layer patterns
- [ ] Document authorization strategy
- [ ] Add diagrams if helpful

---

### Task 7.3: Write Integration Tests
- [ ] Test complete request/response flow for posts
- [ ] Test complete request/response flow for categories
- [ ] Test complete request/response flow for media
- [ ] Test error scenarios
- [ ] Test authorization

---

### Task 7.4: Update Project Structure Analysis
- [ ] Re-run structure analysis
- [ ] Update scores in PROJECT_STRUCTURE_ANALYSIS.md
- [ ] Mark completed items
- [ ] Document improvements

---

## Phase 8: Final Verification (Priority: CRITICAL)
**Duration:** 2 days
**Risk:** Low
**Dependencies:** All previous phases complete

### Task 8.1: Full Backend Verification
- [ ] Run `go build ./...` - must succeed
- [ ] Run `go test ./...` - all tests must pass
- [ ] Run `go vet ./...` - no issues
- [ ] Run `golangci-lint run` - if available
- [ ] Check Swagger UI - all endpoints documented
- [ ] Manual test each endpoint

---

### Task 8.2: Full Frontend Verification
- [ ] Run `npm run typecheck` - must succeed
- [ ] Run `npm run build` - must succeed
- [ ] Run `npm run test` - if tests exist
- [ ] Manual test all features
- [ ] Check browser console - no errors
- [ ] Verify responsive design

---

### Task 8.3: Integration Testing
- [ ] Test complete user flow: register → login → create post → edit → delete
- [ ] Test media upload and display
- [ ] Test category management
- [ ] Test error handling
- [ ] Test authorization (try unauthorized actions)
- [ ] Test pagination
- [ ] Test filtering and sorting

---

### Task 8.4: Performance Check
- [ ] Check API response times
- [ ] Check frontend load times
- [ ] Check database query performance
- [ ] Identify any N+1 query issues
- [ ] Check bundle size

---

### Task 8.5: Final Documentation
- [ ] Update README with new patterns
- [ ] Update CHANGELOG
- [ ] Document breaking changes (should be none)
- [ ] Create upgrade guide if needed
- [ ] Update code comments

---

## Success Criteria

### Backend (Target: 95/100)
- [x] All request packages use feature-specific names ✅
- [ ] All features have complete response DTOs ⬜
- [ ] Response/request structure is consistent ⬜
- [ ] Service method signatures are standardized ⬜
- [ ] Authorization pattern is clear and consistent ⬜
- [ ] All tests pass ⬜
- [ ] Swagger documentation is complete ⬜

### Frontend (Target: 95/100)
- [ ] All types use camelCase ⬜
- [ ] All timestamps are Date objects ⬜
- [ ] API transformation layer is implemented ⬜
- [ ] All features have hooks ⬜
- [ ] Auth is extracted to separate feature ⬜
- [ ] Error boundaries are implemented ⬜
- [ ] All tests pass ⬜

### Overall
- [ ] No breaking changes introduced ⬜
- [ ] All existing functionality works ⬜
- [ ] Code is cleaner and more maintainable ⬜
- [ ] Documentation is up to date ⬜
- [ ] Team can understand new patterns ⬜

---

## Risk Management

### Potential Risks

1. **Breaking API Contract**
   - **Mitigation:** Keep transformation layer, maintain backward compatibility
   - **Rollback Plan:** Revert commits, restore old structure

2. **Frontend Type Errors**
   - **Mitigation:** Thorough TypeScript checking, gradual rollout
   - **Rollback Plan:** Feature flags to disable new transformations

3. **Performance Regression**
   - **Mitigation:** Monitor response times, optimize transformations
   - **Rollback Plan:** Remove transformation layer if too slow

4. **Team Confusion**
   - **Mitigation:** Clear documentation, team training session
   - **Rollback Plan:** N/A - address through better docs

---

## Progress Tracking

**Last Updated:** March 6, 2026

### Phase Completion Status

| Phase | Status | Completion % | Notes |
|-------|--------|--------------|-------|
| Phase 1: Package Naming | ⬜ Not Started | 0% | Quick win, start here |
| Phase 2: Response DTOs | ⬜ Not Started | 0% | Critical for API consistency |
| Phase 3: Media Standardization | ⬜ Not Started | 0% | Medium priority |
| Phase 4: Service Layer | ⬜ Not Started | 0% | Can run parallel to Phase 5 |
| Phase 5: Frontend Types | ⬜ Not Started | 0% | Depends on Phase 2 |
| Phase 6: Missing Features | ⬜ Not Started | 0% | Nice to have |
| Phase 7: Documentation | ⬜ Not Started | 0% | Final polish |
| Phase 8: Verification | ⬜ Not Started | 0% | Must complete |

**Overall Progress:** 0% (0/8 phases complete)

---

## Next Steps

1. **Immediate (Today):**
   - [ ] Review this TODO with team
   - [ ] Get approval for Phase 1
   - [ ] Start Phase 1 (2-3 hours)

2. **This Week:**
   - [ ] Complete Phase 1
   - [ ] Start Phase 2 (Response DTOs)
   - [ ] Complete 2-3 features in Phase 2

3. **Next Week:**
   - [ ] Complete Phase 2
   - [ ] Start Phase 3 and 5 in parallel
   - [ ] Begin frontend alignment

4. **Week 3:**
   - [ ] Complete Phases 3, 4, 5
   - [ ] Start Phase 6 (optional features)

5. **Week 4:**
   - [ ] Complete Phase 7 (Documentation)
   - [ ] Complete Phase 8 (Final verification)
   - [ ] Celebrate achieving A grade! 🎉

---

**Note:** This is a living document. Update task status as you progress, add notes, and adjust timelines as needed.
