package services

import (
	"regexp"
	"testing"
	"url-shortener/db"
	"url-shortener/enums"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetClicks(t *testing.T) {
	type test struct {
		mockCountResult interface{}
		expectedResult  GetClicksResult
	}

	tests := []test{
		{
			mockCountResult: 12,
			expectedResult: GetClicksResult{
				Error:  nil,
				Status: enums.GetClicksResultSuccessful,
				Count:  12,
			},
		},
		{
			mockCountResult: "boom",
			expectedResult: GetClicksResult{
				Error:  nil,
				Status: enums.GetClicksResultUnknownError,
			},
		},
	}

	sqlDB, mock, err := sqlmock.New()

	assert.Nil(t, err)
	defer sqlDB.Close()

	for _, tc := range tests {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM short_urls")).
			WithArgs("slug").
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(tc.mockCountResult))

		gormDB, err := db.ConnectDatabaseWithoutMigrating(sqlDB)

		assert.Nil(t, err)

		subject := GetClicksService{DB: gormDB}
		result := subject.GetClicks("slug")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t, tc.expectedResult.Status, result.Status)
	}
}

func TestGetClicksWithValidSlug(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM short_urls")).
		WithArgs("slug").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(12))

	gormDB, err := db.ConnectDatabaseWithoutMigrating(sqlDB)

	subject := GetClicksService{DB: gormDB}
	result := subject.GetClicks("slug")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, enums.GetClicksResultSuccessful, result.Status)
	assert.Equal(t, int64(12), result.Count)
}

func TestGetClicksWithInvalidSlug(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM short_urls")).
		WithArgs("invalid").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}))

	gormDB, err := db.ConnectDatabaseWithoutMigrating(sqlDB)

	subject := GetClicksService{DB: gormDB}
	result := subject.GetClicks("invalid")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, enums.GetClicksResultNotFound, result.Status)
}

func TestGetClicksWithError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM short_urls")).
		WithArgs("error").
		// Purposely returning a string to trigger error branch
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow("boom"))

	gormDB, err := db.ConnectDatabaseWithoutMigrating(sqlDB)

	subject := GetClicksService{DB: gormDB}
	result := subject.GetClicks("error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, enums.GetClicksResultUnknownError, result.Status)
}
