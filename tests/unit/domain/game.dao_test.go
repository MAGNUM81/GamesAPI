package domain

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils"
	"GamesAPI/src/utils/errorUtils"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)


type GameTestSuite struct {
	suite.Suite
	DB   	   *gorm.DB
	mock 	   sqlmock.Sqlmock

	repository domain.GameRepoInterface
	game       *domain.Game
	dsnCount   int64
}

func (s *GameTestSuite) BeforeTest(_, _ string) {
	var (
		err error
	)
	s.dsnCount++
	dsn := fmt.Sprintf("sqlmock_db_game_%d", s.dsnCount)
	_, s.mock, err = sqlmock.NewWithDSN(dsn)
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("sqlmock", dsn)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = domain.NewGameRepository(s.DB)
}

func(s *GameTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestGamesTestSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}

func (s *GameTestSuite) TestGameRepo_GetAll_Empty() {

	//Test for getting all games from empty table.
	const sqlSelectAll = `SELECT (.+) FROM "games"`
	s.mock.ExpectQuery(sqlSelectAll).
		WillReturnRows(sqlmock.NewRows(nil))

	data, err := s.repository.GetAll()
	if err != nil {
		assert.Fail(s.T(), "An error occurred during repo.GetAll")
	}

	assert.Equal(s.T(), []domain.Game{}, data) //result should be an empty slice of games.
}

func (s *GameTestSuite) TestGameRepo_GetAll_NotEmpty() {
	rows := sqlmock.NewRows([]string{"id", "title", "developer", "publisher", "releaseDate"}).
		AddRow(1, "Rocket League", "Psyonix", "Psyonix", utils.GetDate("2015-07-07")).
		AddRow(2, "The Witcher 3: Wild Hunt", "CD PROJEKT RED", "CD PROJEKT RED", utils.GetDate("2015-05-18"))
	s.mock.ExpectQuery(`SELECT (.+) FROM "games"`).
		WillReturnRows(rows)

	data, err := s.repository.GetAll()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	expected := []domain.Game{
		{
			ID: 		 1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
			ReleaseDate: utils.GetDate("2015-07-07"),
		},
		{
			ID: 		 2,
			Title:       "The Witcher 3: Wild Hunt",
			Developer:   "CD PROJEKT RED",
			Publisher:   "CD PROJEKT RED",
			ReleaseDate: utils.GetDate("2015-05-18"),
		},
	}

	assert.Equal(s.T(), expected[0], data[0])
	assert.Equal(s.T(), expected[1], data[1])
}

//Test for getting a single game from empty table.
func (s *GameTestSuite) TestGameRepo_Get_Empty() {
	rows := sqlmock.NewRows(nil)
	const sql = `SELECT (.+) FROM "games"`
	s.mock.ExpectQuery(sql).WithArgs(0).WillReturnRows(rows)

	game, err := s.repository.Get(0)
	require.True(s.T(), game == nil)
	require.True(s.T(), err != nil)
}

//Test for getting a single game from table with one matching row
func (s *GameTestSuite) TestGameRepo_Get_OneValidRow() {
	rows := sqlmock.NewRows([]string{"id", "title", "developer", "publisher", "releaseDate"}).
		AddRow(1, "Rocket League", "Psyonix", "Psyonix", utils.GetDate("2015-07-07"))
	const sql = `SELECT (.+) FROM "games"`
	s.mock.ExpectQuery(sql).WillReturnRows(rows)

	game, err := s.repository.Get(1)
	require.True(s.T(), err == nil)

	expected := &domain.Game{
		ID:    1,
		Title:  "Rocket League",
		Publisher: "Psyonix",
		Developer: "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
	}

	assert.Equal(s.T(), expected, game)
}

//Test for getting a single game from table with no matching rows
func (s *GameTestSuite) TestGameRepo_Get_OneInvalidRow() {
	const sql = `SELECT (.+) FROM "games"`
	s.mock.ExpectQuery(sql).WillReturnError(gorm.Errors{})

	game, err := s.repository.Get(1)

	require.True(s.T(), err != nil)
	require.True(s.T(), game == nil)
}

//Test for updating a non-existing game
func (s *GameTestSuite) TestGameRepo_Update_NotExist() {
	const sql = `SELECT`
	s.mock.ExpectQuery(sql).WillReturnError(
		errorUtils.NewNotFoundError("not_found"))

	game := &domain.Game{}
	var err errorUtils.EntityError = nil
	game, err = s.repository.Update(game)
	require.True(s.T(), err != nil)
	assert.Equal(s.T(), errorUtils.NewNotFoundError("not_found"), err)
	require.True(s.T(), game == nil)
}

//Test for updating an existing game
func (s *GameTestSuite) TestGameRepo_Update_Exists() {
	selectRows := sqlmock.NewRows([]string{"id", "title", "developer", "publisher", "releaseDate"}).
		AddRow(1, "Rocket League", "Psyonix", "Psyonix", utils.GetDate("2015-07-07"))
	const sqlSelect = `SELECT`
	s.mock.ExpectQuery(sqlSelect).WillReturnRows(selectRows)
	const sqlUpdate = `UPDATE`
	s.mock.ExpectBegin()
	s.mock.ExpectExec(sqlUpdate).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	expected := &domain.Game{
		ID:    1,
		Title:  "Rocket League",
		Publisher: "Psyonix",
		Developer: "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
	}
	game, err := s.repository.Update(expected)
	require.True(s.T(), err == nil)
	assert.Equal(s.T(), expected, game)
}

//Test for inserting a new game
func (s *GameTestSuite) TestGameRepo_Insert_Succeeds() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT`).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	expected := &domain.Game{
		Title:  "Rocket League",
		Publisher: "Psyonix",
		Developer: "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
	}
	created, err := s.repository.Create(expected)

	require.True(s.T(), err == nil)
	assert.EqualValues(s.T(), expected.ID, created.ID)
	assert.EqualValues(s.T(), expected.Title, created.Title)
	assert.EqualValues(s.T(), expected.Publisher, created.Publisher)
	assert.EqualValues(s.T(), expected.Developer, created.Developer)
	assert.EqualValues(s.T(), expected.ReleaseDate, created.ReleaseDate)
}

//Test for when inserting fails (constraints are violated or else)
func (s *GameTestSuite) TestGameRepo_Insert_Fails() {
	expectedErr := errorUtils.NewInternalServerError("server_error")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT`).WillReturnError(expectedErr)
	s.mock.ExpectCommit()

	input := &domain.Game{
		Title:  "dev",
		Publisher: "dev@test.com",
		Developer: "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
	}
	created, err := s.repository.Create(input)

	require.True(s.T(), created == nil)
	assert.Equal(s.T(), expectedErr, err)
}

//Test for deleting an existing game
func (s *GameTestSuite) TestGameRepo_Delete_Succeeds() {
	selectRows := sqlmock.NewRows([]string{"id", "title", "developer", "publisher", "releaseDate"}).
		AddRow(1, "Rocket League", "Psyonix", "Psyonix", utils.GetDate("2015-07-07"))
	s.mock.ExpectQuery(`SELECT`).WillReturnRows(selectRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE (.+) SET "deleted_at"=`).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.Delete(1)
	require.True(s.T(), err == nil)
}

//Test for deleting a non-existing game
func (s *GameTestSuite) TestGameRepo_Delete_NotExist() {
	expectedErr := errorUtils.NewNotFoundError("not_found")
	s.mock.ExpectQuery(`SELECT`).WillReturnError(expectedErr)
	err := s.repository.Delete(1)

	assert.Equal(s.T(), expectedErr, err)
}

//Test for deleting a non-existing game
func (s *GameTestSuite) TestGameRepo_Delete_Fails() {
	expectedErr := errorUtils.NewEntityError(errorUtils.NewError("delete_failed"))
	selectRows := sqlmock.NewRows([]string{"id", "title", "developer", "publisher", "releaseDate"}).
		AddRow(1, "Rocket League", "Psyonix", "Psyonix", utils.GetDate("2015-07-07"))
	s.mock.ExpectQuery(`SELECT`).WillReturnRows(selectRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE (.+) SET "deleted_at"=`).WillReturnError(expectedErr)
	s.mock.ExpectCommit()

	err := s.repository.Delete(1)
	assert.Equal(s.T(), expectedErr, err)
}