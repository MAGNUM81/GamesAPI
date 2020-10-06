package domain

import (
	"GamesAPI/src/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.Errorf(t, err, "An unexpected error occurred when opening a stub database connection.")
	}
	defer db.Close()

	gdb, err := gorm.Open("mssql", db)
	if err != nil {
		assert.Errorf(t, err, "An unexpected error occurred when opening a stub database connection.")
	}
	repo := domain.NewGameRepository(gdb)

	//Test for getting all games from empty table.
	const sqlSelectAll = `SELECT * from games`
	mock.ExpectQuery(sqlSelectAll).
		WillReturnRows(sqlmock.NewRows(nil))

	data, err := repo.GetAll()
	if err != nil {
		assert.Fail(t, "An error occurred during repo.GetAll")
	}

	assert.Equal(t, []domain.Game{}, data) //result should be an empty slice.
	//END: Test for getting all games from empty table
}