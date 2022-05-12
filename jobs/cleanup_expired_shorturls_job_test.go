package jobs

import (
	"regexp"
	"testing"
	"time"
	"url-shortener/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCleanupExpiredShortUrlsReturnsNumberOfDeleteRowsOnSuccess(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM \"short_urls\"")).
		WithArgs(testClock{}.Now()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gormDB, err := db.ConnectDatabaseWithoutMigrating(sqlDB)

	if err != nil {
		t.Fatal(err)
	}

	rowsDeleted, err := CleanupExpiredShortUrls(gormDB, testClock{})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(1), rowsDeleted)
}

type testClock struct{}

func (testClock) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 0, 0, 0, time.UTC)
}
