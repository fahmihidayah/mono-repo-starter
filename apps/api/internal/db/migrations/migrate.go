package migrations

import (
	"log"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/security"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

// AutoMigrate runs automatic migrations for all domain models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto migrations...")

	// Migrate all models
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Account{},
		&domain.UserSession{},
		&domain.TokenBlacklist{},
		&domain.Post{},
		&domain.Category{},
		&domain.Media{},
		&domain.TokenBlacklist{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("✓ Migrations completed successfully")
	return nil
}

// DropAllTables drops all tables (use with caution!)
func DropAllTables(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	// Drop join table first (many-to-many relationship)
	err := db.Migrator().DropTable("post_categories")
	if err != nil {
		log.Printf("Failed to drop join table: %v", err)
		// Continue anyway - table might not exist
	}

	// Drop main tables
	err = db.Migrator().DropTable(
		&domain.User{},
		&domain.Account{},
		&domain.UserSession{},
		&domain.TokenBlacklist{},
		&domain.Post{},
		&domain.Category{},
		&domain.Media{},
		&domain.TokenBlacklist{},
	)

	if err != nil {
		log.Printf("Drop tables failed: %v", err)
		return err
	}

	log.Println("✓ All tables dropped successfully")
	return nil
}

// CreateIndexes creates additional indexes for better performance
func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating additional indexes...")

	// Add any custom indexes here if needed
	// Example:
	// db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")

	log.Println("✓ Indexes created successfully")
	return nil
}

// SeedData seeds initial data into the database
func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// Check if admin user already exists
	var existingUser domain.User
	result := db.Where("email = ?", "admin@example.com").First(&existingUser)

	if result.Error == nil {
		log.Println("Admin user already exists, skipping seed")
		return nil
	}

	// Create admin user
	hashedPassword, err := security.HashPassword("admin123")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	adminUser := &domain.User{
		ID:             utils.GenerateUUID(),
		Name:           "Admin User",
		Email:          "admin@example.com",
		HashedPassword: hashedPassword,
		IsVerified:     true,
		IsSuperUser:    true,
		// Note: CreatedAt and UpdatedAt will be set automatically by GORM autoCreateTime/autoUpdateTime tags
	}

	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return err
	}

	log.Println("✓ Admin user created: admin@example.com / admin123")
	log.Println("✓ Seeding completed successfully")
	return nil
}
