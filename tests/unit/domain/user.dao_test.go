package domain

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite

	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository domain.UserRepoInterface
	user       *domain.User
	dsnCount   int64 //dsnCount setup inspired from https://leethax.org/2019/08/28/golang-sqlmock-gorm.html
}

//before each test, let's init the mocked db
func (s *UserTestSuite) BeforeTest(_, _ string) {
	var (
		err error
	)
	s.dsnCount++
	dsn := fmt.Sprintf("sqlmock_db_user_%d", s.dsnCount)
	_, s.mock, err = sqlmock.NewWithDSN(dsn)
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("sqlmock", dsn)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = domain.NewUserRepository(s.DB)

}

func (s *UserTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

//Test for getting all users from empty table.
func (s *UserTestSuite) TestUserRepo_GetAll_Empty() {
	const sqlSelectAll = `SELECT (.+) FROM "users"`
	s.mock.ExpectQuery(sqlSelectAll).
		WillReturnRows(sqlmock.NewRows(nil))

	data, err := s.repository.GetAll()
	if err != nil {
		assert.Fail(s.T(), "An error occurred during repo.GetAll")
	}

	assert.Equal(s.T(), []domain.User{}, data) //result should be an empty slice of users.
}

//Test for getting all users from table with two rows
func (s *UserTestSuite) TestUserRepo_GetAll_NotEmpty() {
	rows := sqlmock.NewRows([]string{"email", "name"}).
		AddRow("devgolang@test.com", "devgolang").
		AddRow("golangdev@test.com", "golangdev")
	s.mock.ExpectQuery(`SELECT (.+) FROM "users"`).
		WillReturnRows(rows)

	data, err := s.repository.GetAll()
	if err != nil {
		assert.Fail(s.T(), "An error occurred during repo.GetAll")
	}

	expected := []domain.User{
		{
			Email: "devgolang@test.com",
			Name:  "devgolang",
		}, {
			Email: "golangdev@test.com",
			Name:  "golangdev",
		},
	}

	assert.Equal(s.T(), expected, data)
}

//Test for getting a single user from empty table.
func (s *UserTestSuite) TestUserRepo_Get_Empty() {
	rows := sqlmock.NewRows(nil)
	const sql = `SELECT (.+) FROM "users"`
	s.mock.ExpectQuery(sql).WithArgs(0).WillReturnRows(rows)

	user, err := s.repository.Get(0)
	require.True(s.T(), user == nil)
	require.True(s.T(), err != nil)
}

//Test for getting a single user from table with one matching row
func (s *UserTestSuite) TestUserRepo_Get_OneValidRow() {
	rows := sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(1, "devgolang@test.com", "devgolang")
	const sql = `SELECT (.+) FROM "users"`
	s.mock.ExpectQuery(sql).WillReturnRows(rows)

	user, err := s.repository.Get(1)
	require.True(s.T(), err == nil)

	expected := &domain.User{
		ID:    1,
		Name:  "devgolang",
		Email: "devgolang@test.com",
	}

	assert.Equal(s.T(), expected, user)
}

//Test for getting a single user from table with no matching rows
func (s *UserTestSuite) TestUserRepo_Get_OneInvalidRow() {
	const sql = `SELECT (.+) FROM "users"`
	s.mock.ExpectQuery(sql).WillReturnError(gorm.Errors{})

	user, err := s.repository.Get(1)

	require.True(s.T(), err != nil)
	require.True(s.T(), user == nil)
}

//Test for updating a non-existing user
func (s *UserTestSuite) TestUserRepo_Update_NotExist() {
	const sql = `SELECT`
	s.mock.ExpectQuery(sql).WillReturnError(
		errorUtils.NewNotFoundError("not_found"))

	user := &domain.User{}
	var err errorUtils.EntityError = nil
	user, err = s.repository.Update(user)
	require.True(s.T(), err != nil)
	assert.Equal(s.T(), errorUtils.NewNotFoundError("not_found"), err)
	require.True(s.T(), user == nil)
}

//Test for updating an existing user
func (s *UserTestSuite) TestUserRepo_Update_Exists() {
	selectRows := sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(1, "devgolang@test.com", "devgolang")
	const sqlSelect = `SELECT`
	s.mock.ExpectQuery(sqlSelect).WillReturnRows(selectRows)
	const sqlUpdate = `UPDATE`
	s.mock.ExpectBegin()
	s.mock.ExpectExec(sqlUpdate).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	expected := &domain.User{
		ID:    1,
		Name:  "dev",
		Email: "dev@test.com",
	}
	user, err := s.repository.Update(expected)
	require.True(s.T(), err == nil)
	assert.Equal(s.T(), expected, user)
}

//Test for inserting a new user
func (s *UserTestSuite) TestUserRepo_Insert_Succeeds() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT`).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	expected := &domain.User{
		ID:    1,
		Name:  "dev",
		Email: "dev@test.com",
	}
	created, err := s.repository.Create(expected)

	require.True(s.T(), err == nil)
	assert.Equal(s.T(), expected, created)
}

//Test for when inserting fails (constraints are violated or else)
func (s *UserTestSuite) TestUserRepo_Insert_Fails() {
	expectedErr := errorUtils.NewInternalServerError("server_error")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT`).WillReturnError(expectedErr)
	s.mock.ExpectCommit()

	input := &domain.User{
		ID:    0,
		Name:  "dev",
		Email: "dev@test.com",
	}
	created, err := s.repository.Create(input)

	require.True(s.T(), created == nil)
	assert.Equal(s.T(), expectedErr, err)
}

//Test for deleting an existing user
func (s *UserTestSuite) TestUserRepo_Delete_Succeeds() {
	selectRows := sqlmock.NewRows([]string{"id", "email", "name"}).
		AddRow(1, "devgolang@test.com", "devgolang")
	s.mock.ExpectQuery(`SELECT`).WillReturnRows(selectRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE (.+) SET "deleted_at"=`).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.Delete(1)
	require.True(s.T(), err == nil)
}

//Test for deleting a non-existing user
func (s *UserTestSuite) TestUserRepo_Delete_NotExist() {
	expectedErr := errorUtils.NewNotFoundError("not_found")
	s.mock.ExpectQuery(`SELECT`).WillReturnError(expectedErr)
	err := s.repository.Delete(1)

	assert.Equal(s.T(), expectedErr, err)
}

//Test for deleting a non-existing user
func (s *UserTestSuite) TestUserRepo_Delete_Fails() {
	expectedErr := errorUtils.NewEntityError(errorUtils.NewError("delete_failed"))
	selectRows := sqlmock.NewRows([]string{"id", "email", "name"}).
		AddRow(1, "devgolang@test.com", "devgolang")
	s.mock.ExpectQuery(`SELECT`).WillReturnRows(selectRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE (.+) SET "deleted_at"=`).WillReturnError(expectedErr)
	s.mock.ExpectCommit()

	err := s.repository.Delete(1)
	assert.Equal(s.T(), expectedErr, err)
}
