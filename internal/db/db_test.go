package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	// This test is difficult to run without a real database
	// We'll skip it for now and rely on the postgres_test.go file
	// which tests the actual implementation
	t.Skip("Skipping test that requires a real database")
}

// TestConnection_Close tests the Close method
func TestConnection_Close(t *testing.T) {
	// Create a mock DB that we can control
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test successful close
	mock.ExpectClose()
	if err := conn.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Test error on close
	expectedErr := errors.New("close error")
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	conn = &Connection{db: db}

	mock.ExpectClose().WillReturnError(expectedErr)
	if err := conn.Close(); err != expectedErr {
		t.Errorf("Close() error = %v, want %v", err, expectedErr)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestConnection_DB tests the DB method
func TestConnection_DB(t *testing.T) {
	// Create a mock DB
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test
	result := conn.DB()

	// Assert
	if result != db {
		t.Errorf("DB() = %v, want %v", result, db)
	}
}

// TestConnection_Ping tests the Ping method
func TestConnection_Ping(t *testing.T) {
	// Create a mock DB for successful ping
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test successful ping
	mock.ExpectPing()
	if err := conn.Ping(); err != nil {
		t.Errorf("Ping() error = %v, want nil", err)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Create a mock DB for ping error
	db, mock, err = sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn = &Connection{db: db}

	// Test ping error
	expectedErr := errors.New("ping error")
	mock.ExpectPing().WillReturnError(expectedErr)
	if err := conn.Ping(); err != expectedErr {
		t.Errorf("Ping() error = %v, want %v", err, expectedErr)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestConnection_Begin tests the Begin method
func TestConnection_Begin(t *testing.T) {
	// Create a mock DB for successful begin
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test successful begin
	mock.ExpectBegin()
	tx, err := conn.Begin()
	if err != nil {
		t.Errorf("Begin() error = %v, want nil", err)
	}
	if tx == nil {
		t.Error("Begin() returned nil transaction on success")
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Create a mock DB for begin error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn = &Connection{db: db}

	// Test begin error
	expectedErr := errors.New("begin error")
	mock.ExpectBegin().WillReturnError(expectedErr)
	tx, err = conn.Begin()
	if err != expectedErr {
		t.Errorf("Begin() error = %v, want %v", err, expectedErr)
	}
	if tx != nil {
		t.Error("Begin() returned non-nil transaction on error")
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestConnection_Exec tests the Exec method
func TestConnection_Exec(t *testing.T) {
	// Create a mock DB for successful exec
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test successful exec
	expectedResult := sqlmock.NewResult(1, 2)
	mock.ExpectExec("SELECT 1").WillReturnResult(expectedResult)
	result, err := conn.Exec("SELECT 1")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	
	// Check that result is not nil
	if result == nil {
		t.Error("Exec() returned nil result")
	}
	
	// Check the result values
	lastID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	expectedLastID, _ := expectedResult.LastInsertId()
	expectedRowsAffected, _ := expectedResult.RowsAffected()
	
	if lastID != expectedLastID {
		t.Errorf("LastInsertId = %d, want %d", lastID, expectedLastID)
	}
	if rowsAffected != expectedRowsAffected {
		t.Errorf("RowsAffected = %d, want %d", rowsAffected, expectedRowsAffected)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Create a mock DB for exec error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn = &Connection{db: db}

	// Test exec error
	expectedErr := errors.New("exec error")
	mock.ExpectExec("SELECT 1").WillReturnError(expectedErr)
	result, err = conn.Exec("SELECT 1")
	if err != expectedErr {
		t.Errorf("Exec() error = %v, want %v", err, expectedErr)
	}
	if result != nil {
		t.Error("Exec() returned non-nil result on error")
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestConnection_Query tests the Query method
func TestConnection_Query(t *testing.T) {
	// Create a mock DB for successful query
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test successful query
	expectedRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT 1").WillReturnRows(expectedRows)
	rows, err := conn.Query("SELECT 1")
	if err != nil {
		t.Errorf("Query() error = %v, want nil", err)
	}
	if rows == nil {
		t.Error("Query() returned nil rows on success")
	}
	defer rows.Close()

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Create a mock DB for query error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn = &Connection{db: db}

	// Test query error
	expectedErr := errors.New("query error")
	mock.ExpectQuery("SELECT 1").WillReturnError(expectedErr)
	rows, err = conn.Query("SELECT 1")
	if err != expectedErr {
		t.Errorf("Query() error = %v, want %v", err, expectedErr)
	}
	if rows != nil {
		t.Error("Query() returned non-nil rows on error")
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestConnection_QueryRow tests the QueryRow method
func TestConnection_QueryRow(t *testing.T) {
	// Create a mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create a connection with our mock DB
	conn := &Connection{db: db}

	// Test QueryRow
	expectedRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT 1").WillReturnRows(expectedRows)
	row := conn.QueryRow("SELECT 1")
	
	// We can't directly compare sql.Row objects, but we can verify it's not nil
	if row == nil {
		t.Error("QueryRow() returned nil")
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
