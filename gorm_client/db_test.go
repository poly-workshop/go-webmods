package gorm_client

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenSqlite_DirectoryCreation(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()
	
	tests := []struct {
		name           string
		dbPath         string
		expectSuccess  bool
		shouldCreateDir bool
	}{
		{
			name:           "simple filename",
			dbPath:         "test.db",
			expectSuccess:  true,
			shouldCreateDir: false,
		},
		{
			name:           "relative path with existing directory",
			dbPath:         filepath.Join(tempDir, "test.db"),
			expectSuccess:  true,
			shouldCreateDir: false,
		},
		{
			name:           "relative path with non-existing directory",
			dbPath:         filepath.Join(tempDir, "subdir", "test.db"),
			expectSuccess:  true,
			shouldCreateDir: true,
		},
		{
			name:           "nested path with non-existing directories",
			dbPath:         filepath.Join(tempDir, "level1", "level2", "test.db"),
			expectSuccess:  true,
			shouldCreateDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			if tt.shouldCreateDir {
				dirPath := filepath.Dir(tt.dbPath)
				os.RemoveAll(dirPath)
				
				// Verify directory doesn't exist
				if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
					t.Fatalf("Directory should not exist before test: %s", dirPath)
				}
			}

			cfg := Config{
				Driver: "sqlite",
				Name:   tt.dbPath,
			}

			db, err := openSqlite(cfg)
			
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("openSqlite() failed: %v", err)
					return
				}
				if db == nil {
					t.Error("openSqlite() returned nil db")
					return
				}
				
				// Close the database
				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
				
				// Verify directory was created if needed
				if tt.shouldCreateDir {
					dirPath := filepath.Dir(tt.dbPath)
					if _, err := os.Stat(dirPath); os.IsNotExist(err) {
						t.Errorf("Directory was not created: %s", dirPath)
					}
				}
			} else {
				if err == nil {
					t.Error("openSqlite() should have failed")
				}
			}
			
			// Clean up after test
			if tt.dbPath != "test.db" {
				os.Remove(tt.dbPath)
				if tt.shouldCreateDir {
					os.RemoveAll(filepath.Dir(tt.dbPath))
				}
			} else {
				os.Remove("test.db")
			}
		})
	}
}

func TestNewDB_SqliteDirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "subdir", "app.db")
	
	// Verify directory doesn't exist
	dirPath := filepath.Dir(dbPath)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		t.Fatalf("Directory should not exist before test: %s", dirPath)
	}

	cfg := Config{
		Driver: "sqlite",
		Name:   dbPath,
	}

	db := NewDB(cfg)
	if db == nil {
		t.Error("NewDB() returned nil")
		return
	}
	
	// Close the database
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
	
	// Verify directory was created
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("Directory was not created: %s", dirPath)
	}
	
	// Clean up
	os.Remove(dbPath)
	os.RemoveAll(dirPath)
}

func TestOpenSqlite_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		dbPath         string
		expectSuccess  bool
		description    string
	}{
		{
			name:           "memory database",
			dbPath:         ":memory:",
			expectSuccess:  true,
			description:    "in-memory database should work",
		},
		{
			name:           "absolute path",
			dbPath:         filepath.Join(t.TempDir(), "abs", "test.db"),
			expectSuccess:  true,
			description:    "absolute path should work",
		},
		{
			name:           "current directory",
			dbPath:         "current.db",
			expectSuccess:  true,
			description:    "current directory should work",
		},
		{
			name:           "dot slash prefix",
			dbPath:         "./dotslash.db",
			expectSuccess:  true,
			description:    "dot slash prefix should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Driver: "sqlite",
				Name:   tt.dbPath,
			}

			db, err := openSqlite(cfg)
			
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("openSqlite() failed for %s: %v", tt.description, err)
					return
				}
				if db == nil {
					t.Errorf("openSqlite() returned nil db for %s", tt.description)
					return
				}
				
				// Close the database
				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
			} else {
				if err == nil {
					t.Errorf("openSqlite() should have failed for %s", tt.description)
				}
			}
			
			// Clean up - be careful with special paths
			if tt.dbPath != ":memory:" && tt.dbPath != "current.db" && tt.dbPath != "./dotslash.db" {
				os.Remove(tt.dbPath)
				dir := filepath.Dir(tt.dbPath)
				if dir != "." && dir != "/" {
					os.RemoveAll(dir)
				}
			} else if tt.dbPath == "current.db" || tt.dbPath == "./dotslash.db" {
				os.Remove(tt.dbPath)
			}
		})
	}
}