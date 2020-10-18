package domain

import (
	"GamesAPI/src/domain"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GameTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository domain.GameRepoInterface
	game     *domain.Game
	dsnCount int64
}

func (s *GameTestSuite) BeforeTest(_, _ string) {
	var (
		err error
	)
	dsn := fmt.Sprintf("sqlmock_db_%d", s.dsnCount)
	_, s.mock, err = sqlmock.NewWithDSN(dsn)
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("sqlmock", dsn)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = domain.NewGameRepository(s.DB)
	s.dsnCount++
}

func(s *GameTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestGamesTestSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}

func (s *GameTestSuite) TestGameRepo_GetAll() {

	//Test for getting all games from empty table.
	const sqlSelectAll = `SELECT (.+) FROM "games"`
	s.mock.ExpectQuery(sqlSelectAll).
		WillReturnRows(sqlmock.NewRows(nil))

	data, err := s.repository.GetAll()
	if err != nil {
		assert.Fail(s.T(), "An error occurred during repo.GetAll")
	}

	assert.Equal(s.T(), []domain.Game{}, data) //result should be an empty slice.
	//END: Test for getting all games from empty table
}