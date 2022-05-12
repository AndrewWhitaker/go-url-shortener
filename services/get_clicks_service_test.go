package services

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"
	"url-shortener/db"
	"url-shortener/enums"
	"url-shortener/models"
	"url-shortener/test/helpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

func TestGetClicks(t *testing.T) {
	type test struct {
		mockCountResult interface{}
		expectedResult  td.StructFields
	}

	tests := []test{
		{
			mockCountResult: 12,
			expectedResult: td.StructFields{
				"Error":  nil,
				"Status": enums.GetClicksResultSuccessful,
				"Count":  int64(12),
			},
		},
		{
			mockCountResult: "boom",
			expectedResult: td.StructFields{
				"Error":  td.HasPrefix("sql: Scan error"),
				"Status": enums.GetClicksResultUnknownError,
				"Count":  int64(0),
			},
		},
		{
			mockCountResult: nil,
			expectedResult: td.StructFields{
				"Error":  nil,
				"Status": enums.GetClicksResultNotFound,
				"Count":  int64(0),
			},
		},
	}

	sqlDB, mock, err := sqlmock.New()

	assert.Nil(t, err)
	defer sqlDB.Close()

	for _, tc := range tests {
		rows := sqlmock.NewRows([]string{"count(*)"})

		if tc.mockCountResult != nil {
			rows.AddRow(tc.mockCountResult)
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM short_urls")).
			WithArgs("slug").
			WillReturnRows(rows)

		gormDB, err := db.ConnectDatabaseWithoutMigrating(sqlDB)

		assert.Nil(t, err)

		subject := GetClicksService{
			DB: gormDB, Clock: SystemClock{},
		}

		result := subject.GetClicks("slug", enums.GetClicksTimePeriodAllTime)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		td.Cmp(t, result, td.Struct(GetClicksResult{}, tc.expectedResult))
	}
}

func TestGetClicksFunctional(t *testing.T) {
	// Spins up an actual DB, inserts some data, and tests the query
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	container, sqlDB, err := helpers.CreateTestContainer(ctx, "clicksdb")
	if err != nil {
		t.Fatal(err)
	}

	defer container.Terminate(ctx)
	defer sqlDB.Close()

	gormDB, err := db.ConnectDatabase(sqlDB)
	assert.Nil(t, err)

	clock := TestClock{}

	fmt.Printf("clock = %v", clock)

	subject := GetClicksService{
		DB:    gormDB,
		Clock: clock,
	}

	times := []time.Time{
		time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 17, 16, 45, 0, 0, time.UTC),
		time.Date(2022, 4, 30, 11, 30, 0, 0, time.UTC),
		time.Date(2022, 5, 3, 20, 30, 0, 0, time.UTC),
		time.Date(2022, 5, 8, 11, 30, 0, 0, time.UTC),
		time.Date(2022, 5, 10, 13, 45, 0, 0, time.UTC),
	}

	shortUrl := models.ShortUrl{
		Slug:    "slug",
		LongUrl: "https://www.cloudflare.com",
	}

	err = gormDB.Create(&shortUrl).Error

	assert.Nil(t, err)

	for _, time := range times {
		err = gormDB.Create(&models.Click{
			CreatedAt:  time,
			ShortUrlId: shortUrl.Id,
		}).Error

		assert.Nil(t, err)
	}

	type test struct {
		timePeriod    enums.GetClicksTimePeriod
		expectedCount int64
	}

	tests := []test{
		{timePeriod: enums.GetClicksTimePeriodAllTime, expectedCount: 7},
		{timePeriod: enums.GetClicksTimePeriodPastWeek, expectedCount: 3},
		{timePeriod: enums.GetClicksTimePeriod24Hours, expectedCount: 1},
	}

	for _, tc := range tests {
		actual := subject.GetClicks("slug", tc.timePeriod)
		assert.Equal(t, tc.expectedCount, actual.Count)
	}
}

type TestClock struct{}

func (TestClock) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 0, 0, 0, time.UTC)
}
