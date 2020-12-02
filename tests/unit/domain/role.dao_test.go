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

type UserRoleTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository domain.UserRoleRepoInterface
	userRole   *domain.UserRole
	dsnCount   int64
}

func (s *UserRoleTestSuite) BeforeTest(_, _ string) {
	var (
		err error
	)
	s.dsnCount++
	dsn := fmt.Sprintf("sqlmock_db_userRole_%d", s.dsnCount)
	_, s.mock, err = sqlmock.NewWithDSN(dsn)
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("sqlmock", dsn)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = domain.NewUserRoleRepository(s.DB)
}

func (s *UserRoleTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestUserRolesTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleTestSuite))
}

func (s *UserRoleTestSuite) TestUserRoleRepo_GetAll_Empty() {

	//Test for getting all userRoles from empty table.
	const sqlSelectAll = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sqlSelectAll).
		WillReturnRows(sqlmock.NewRows(nil))

	data, err := s.repository.GetAll()
	if err != nil {
		assert.Fail(s.T(), "An error occurred during repo.GetAll")
	}

	assert.Equal(s.T(), []domain.UserRole{}, data) //result should be an empty slice of userRoles.
}

func (s *UserRoleTestSuite) TestUserRoleRepo_GetAll_NotEmpty() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin").
		AddRow(2, 2, "User")
	s.mock.ExpectQuery(`SELECT (.+) FROM "user_roles"`).
		WillReturnRows(rows)

	data, err := s.repository.GetAll()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	expected := []domain.UserRole{
		{
			ID:     1,
			UserID: 1,
			Name:   "Admin",
		},
		{
			ID:     2,
			UserID: 2,
			Name:   "User",
		},
	}

	assert.Equal(s.T(), expected[0], data[0])
	assert.Equal(s.T(), expected[1], data[1])
}

//Test for getting a single userRole from empty table.
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByID_Empty() {
	rows := sqlmock.NewRows(nil)
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WithArgs(1).WillReturnRows(rows)

	userRole, err := s.repository.GetByID(1)
	require.True(s.T(), userRole == nil)
	require.True(s.T(), err != nil)
}

//Test for getting a single userRole from table with one matching row
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByID_OneValidRow() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin")
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WillReturnRows(rows)

	userRole, err := s.repository.GetByID(1)
	require.True(s.T(), err == nil)

	expected := &domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "Admin",
	}

	assert.Equal(s.T(), expected, userRole)
}

//Test for getting a single userRole from table with no matching rows
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByID_OneInvalidRow() {
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WillReturnError(gorm.Errors{})

	userRole, err := s.repository.GetByID(1)

	require.NotNil(s.T(), err)
	require.Nil(s.T(), userRole)
}

//Test for getting a single userRole from empty table.
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByRoleName_Empty() {
	rows := sqlmock.NewRows(nil)
	const sql = `SELECT (.+) FROM "user_roles" WHERE (.+) `
	s.mock.ExpectQuery(sql).WithArgs("Admin").WillReturnRows(rows)

	userRole, err := s.repository.GetByRole("Admin")
	require.Equal(s.T(), 0, len(userRole))
	require.Nil(s.T(), err)
}

//Test for getting a single userRole from table with one matching row
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByRoleName_OneValidRow() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin")
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WillReturnRows(rows)

	userRole, err := s.repository.GetByID(1)
	require.True(s.T(), err == nil)

	expected := &domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "Admin",
	}

	assert.Equal(s.T(), expected, userRole)
}

//Test for getting a single userRole from table with no matching rows
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByRoleName_OneInvalidRow() {
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WillReturnError(gorm.Errors{})

	userRole, err := s.repository.GetByID(1)

	require.True(s.T(), err != nil)
	require.True(s.T(), userRole == nil)
}

//Test for getting a single userRole from empty table.
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByUserID_Empty() {
	rows := sqlmock.NewRows(nil)
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WithArgs(1).WillReturnRows(rows)

	userRole, err := s.repository.GetByID(1)
	require.True(s.T(), userRole == nil)
	require.True(s.T(), err != nil)
}

//Test for getting a single userRole from table with one matching row
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByUserID_OneValidRow() {
	rows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin")
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WillReturnRows(rows)

	userRole, err := s.repository.GetByID(1)
	require.True(s.T(), err == nil)

	expected := &domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "Admin",
	}

	assert.Equal(s.T(), expected, userRole)
}

//Test for getting a single userRole from table with no matching rows
func (s *UserRoleTestSuite) TestUserRoleRepo_GetByUserID_OneInvalidRow() {
	const sql = `SELECT (.+) FROM "user_roles"`
	s.mock.ExpectQuery(sql).WillReturnError(gorm.Errors{})

	userRole, err := s.repository.GetByID(1)

	require.True(s.T(), err != nil)
	require.True(s.T(), userRole == nil)
}

//Test for updating a non-existing userRole
func (s *UserRoleTestSuite) TestUserRoleRepo_Update_NotExist() {
	const sql = `SELECT`
	s.mock.ExpectQuery(sql).WillReturnError(
		errorUtils.NewNotFoundError("not_found"))

	userRole := &domain.UserRole{}
	var err errorUtils.EntityError = nil
	userRole, err = s.repository.Update(userRole)
	require.True(s.T(), err != nil)
	assert.Equal(s.T(), errorUtils.NewNotFoundError("not_found"), err)
	require.True(s.T(), userRole == nil)
}

//Test for updating an existing userRole
func (s *UserRoleTestSuite) TestUserRoleRepo_Update_Exists() {
	selectRows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin")
	const sqlSelect = `SELECT`
	s.mock.ExpectQuery(sqlSelect).WillReturnRows(selectRows)
	const sqlUpdate = `UPDATE`
	s.mock.ExpectBegin()
	s.mock.ExpectExec(sqlUpdate).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	expected := &domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "Admin",
	}
	userRole, err := s.repository.Update(expected)
	require.True(s.T(), err == nil)
	assert.Equal(s.T(), expected, userRole)
}

//Test for inserting a new userRole
func (s *UserRoleTestSuite) TestUserRoleRepo_Insert_Succeeds() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT`).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	expected := &domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "Admin",
	}
	created, err := s.repository.Create(expected)

	require.True(s.T(), err == nil)
	assert.EqualValues(s.T(), expected.ID, created.ID)
	assert.EqualValues(s.T(), expected.UserID, created.UserID)
	assert.EqualValues(s.T(), expected.Name, created.Name)
}

//Test for when inserting fails (constraints are violated or else)
func (s *UserRoleTestSuite) TestUserRoleRepo_Insert_Fails() {
	expectedErr := errorUtils.NewInternalServerError("server_error")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT`).WillReturnError(expectedErr)
	s.mock.ExpectCommit()

	input := &domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "Admin",
	}
	created, err := s.repository.Create(input)

	require.True(s.T(), created == nil)
	assert.Equal(s.T(), expectedErr, err)
}

//Test for deleting an existing userRole
func (s *UserRoleTestSuite) TestUserRoleRepo_Delete_Succeeds() {
	selectRows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin")
	s.mock.ExpectQuery(`SELECT`).WillReturnRows(selectRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE (.+) SET "deleted_at"=`).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.Delete(1)
	require.True(s.T(), err == nil)
}

//Test for deleting a non-existing userRole
func (s *UserRoleTestSuite) TestUserRoleRepo_Delete_NotExist() {
	expectedErr := errorUtils.NewNotFoundError("not_found")
	s.mock.ExpectQuery(`SELECT`).WillReturnError(expectedErr)
	err := s.repository.Delete(1)

	assert.Equal(s.T(), expectedErr, err)
}

//Test for deleting a non-existing userRole
func (s *UserRoleTestSuite) TestUserRoleRepo_Delete_Fails() {
	expectedErr := errorUtils.NewEntityError(errorUtils.NewError("delete_failed"))
	selectRows := sqlmock.NewRows([]string{"id", "user_id", "roleName"}).
		AddRow(1, 1, "Admin")
	s.mock.ExpectQuery(`SELECT`).WillReturnRows(selectRows)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE (.+) SET "deleted_at"=`).WillReturnError(expectedErr)
	s.mock.ExpectCommit()

	err := s.repository.Delete(1)
	assert.Equal(s.T(), expectedErr, err)
}
