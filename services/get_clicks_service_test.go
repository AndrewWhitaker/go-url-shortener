package services

import (
	"regexp"
	"testing"
	"url-shortener/db"
	"url-shortener/enums"

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

		subject := GetClicksService{DB: gormDB}
		result := subject.GetClicks("slug")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		td.Cmp(t, result, td.Struct(GetClicksResult{}, tc.expectedResult))
	}
}
