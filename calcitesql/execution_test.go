package calcitesql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExecuteQuery(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name  string
		query string
		mock  func()
	}{
		{
			name:  "simple select query",
			query: "SELECT * FROM test",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test1").
					AddRow(2, "test2")
				mock.ExpectQuery("SELECT \\* FROM test").WillReturnRows(rows)
			},
		},
		{
			name:  "empty result query",
			query: "SELECT * FROM empty",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT \\* FROM empty").WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ExecuteQuery(db, tt.query)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
