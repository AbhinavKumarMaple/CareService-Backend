package user

import (
	"regexp"
	"testing"
	"time"

	domainUser "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)
	cleanup := func() { db.Close() }
	return gormDB, mock, cleanup
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestTableName(t *testing.T) {
	u := &User{}
	assert.Equal(t, "users", u.TableName())
}

func TestNewUserRepository(t *testing.T) {
	db, _, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	assert.NotNil(t, repo)
}

func TestToDomainMapper(t *testing.T) {
	u := &User{
		ID:        uuid.New(),
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Status:    true,
		Role:      "caregiver",
		Location:  domainUser.Location{HouseNumber: "1", Street: "Main St"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	d := u.toDomainMapper()
	assert.Equal(t, u.UserName, d.UserName)
	assert.Equal(t, u.Role, d.Role)
	assert.Equal(t, u.Location.Street, d.Location.Street)
}

func TestFromDomainMapper(t *testing.T) {
	d := &domainUser.User{
		ID:        uuid.New(),
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Status:    true,
		Role:      "caregiver",
		Location:  domainUser.Location{HouseNumber: "1", Street: "Main St"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	u := fromDomainMapper(d)
	assert.Equal(t, d.UserName, u.UserName)
	assert.Equal(t, d.Role, u.Role)
	assert.Equal(t, d.Location.Street, u.Location.Street)
}

func TestArrayToDomainMapper(t *testing.T) {
	arr := &[]User{{ID: uuid.New(), UserName: "A"}, {ID: uuid.New(), UserName: "B"}}
	d := arrayToDomainMapper(arr)
	assert.Len(t, *d, 2)
	assert.Equal(t, "A", (*d)[0].UserName)
}

func TestRepository_GetAll(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	rows := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password", "role", "location_house_number", "location_street", "location_city", "location_state", "location_pincode", "location_lat", "location_long"}).
		AddRow(uuid.New(), "user1", "a@a.com", "A", "B", true, "hash1", "caregiver", "1", "Main St", "City", "State", "12345", 1.0, 2.0).
		AddRow(uuid.New(), "user2", "b@b.com", "C", "D", false, "hash2", "client", "2", "Second St", "Town", "State", "67890", 3.0, 4.0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(rows)
	users, err := repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, *users, 2)
}

func TestRepository_GetByID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	id1 := uuid.New()
	id2 := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password", "role", "location_house_number", "location_street", "location_city", "location_state", "location_pincode", "location_lat", "location_long"}).
		AddRow(id1, "user1", "a@a.com", "A", "B", true, "hash1", "caregiver", "1", "Main St", "City", "State", "12345", 1.0, 2.0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id1, 1).WillReturnRows(rows)
	user, err := repo.GetByID(id1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.UserName)
	// Not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id2, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password", "role", "location_house_number", "location_street", "location_city", "location_state", "location_pincode", "location_lat", "location_long"}))
	user, err = repo.GetByID(id2)
	assert.Error(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uuid.Nil, user.ID)
}

func TestRepository_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	domainU := &domainUser.User{
		UserName:     "user1",
		Email:        "a@a.com",
		FirstName:    "A",
		LastName:     "B",
		Status:       true,
		HashPassword: "hash1",
		Role:         "caregiver",
		Location:     domainUser.Location{HouseNumber: "1", Street: "Main St"},
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","user_name","email","first_name","last_name","status","hash_password","role","profile_picture","location_house_number","location_street","location_city","location_state","location_pincode","location_lat","location_long","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`)).
		WithArgs(sqlmock.AnyArg(), domainU.UserName, domainU.Email, domainU.FirstName, domainU.LastName, domainU.Status, domainU.HashPassword, domainU.Role, domainU.ProfilePicture, domainU.Location.HouseNumber, domainU.Location.Street, domainU.Location.City, domainU.Location.State, domainU.Location.Pincode, domainU.Location.Lat, domainU.Location.Long, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	user, err := repo.Create(domainU)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.UserName)
}

func TestRepository_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)
	id1 := uuid.New()
	id2 := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(id1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.Delete(id1)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(id2).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	err = repo.Delete(id2)
	assert.Error(t, err)
}

func TestRepository_GetByEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	logger := setupLogger(t)
	repo := NewUserRepository(db, logger)

	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password", "role", "location_house_number", "location_street", "location_city", "location_state", "location_pincode", "location_lat", "location_long"}).
		AddRow(uuid.New(), "user1", email, "A", "B", true, "hash1", "caregiver", "1", "Main St", "City", "State", "12345", 1.0, 2.0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).WillReturnRows(rows)
	user, err := repo.GetByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	// Not found
	emailNotFound := "notfound@example.com"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(emailNotFound, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "status", "hash_password", "role", "location_house_number", "location_street", "location_city", "location_state", "location_pincode", "location_lat", "location_long"}))
	user, err = repo.GetByEmail(emailNotFound)
	assert.Error(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uuid.Nil, user.ID)
}
