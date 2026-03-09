package repository

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserRepoDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))
	return db
}

func makeUser(email, nickname string) *model.User {
	return &model.User{
		Email:    email,
		Password: "hashed",
		Nickname: nickname,
	}
}

func TestUserRepository_Create_StoresUser(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("create@example.com", "creator")
	err := repo.Create(user)

	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}

func TestUserRepository_Create_AssignsUUID(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("uuid@example.com", "uuiduser")
	require.NoError(t, repo.Create(user))

	assert.Len(t, user.ID, 36)
}

func TestUserRepository_FindByID_ReturnsUserWhenExists(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("find@example.com", "finder")
	require.NoError(t, repo.Create(user))

	found, err := repo.FindByID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_FindByID_ReturnsErrorWhenNotFound(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	_, err := repo.FindByID("nonexistent-id")

	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_FindByEmail_ReturnsUserWhenExists(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("email@example.com", "emailuser")
	require.NoError(t, repo.Create(user))

	found, err := repo.FindByEmail("email@example.com")

	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindByEmail_ReturnsErrorWhenNotFound(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	_, err := repo.FindByEmail("notexist@example.com")

	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserRepository_Update_ChangesNickname(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("update@example.com", "oldnick")
	require.NoError(t, repo.Create(user))

	user.Nickname = "newnick"
	require.NoError(t, repo.Update(user))

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "newnick", found.Nickname)
}

func TestUserRepository_Update_ChangesStatusMessage(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("status@example.com", "statususer")
	require.NoError(t, repo.Create(user))

	user.StatusMessage = "I am here"
	require.NoError(t, repo.Update(user))

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "I am here", found.StatusMessage)
}

func TestUserRepository_Delete_SoftDeletesUser(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	user := makeUser("delete@example.com", "deluser")
	require.NoError(t, repo.Create(user))

	require.NoError(t, repo.Delete(user.ID))

	_, err := repo.FindByID(user.ID)
	assert.Error(t, err, "soft-deleted user should not be found via FindByID")
}

func TestUserRepository_Search_ReturnsMatchingByEmail(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	require.NoError(t, repo.Create(makeUser("alice@example.com", "alice")))
	require.NoError(t, repo.Create(makeUser("bob@example.com", "bob")))

	results, err := repo.Search("alice")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "alice@example.com", results[0].Email)
}

func TestUserRepository_Search_ReturnsMatchingByNickname(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	require.NoError(t, repo.Create(makeUser("c@example.com", "charlie")))
	require.NoError(t, repo.Create(makeUser("d@example.com", "delta")))

	results, err := repo.Search("charlie")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "charlie", results[0].Nickname)
}

func TestUserRepository_Search_ReturnsEmptySliceWhenNoMatch(t *testing.T) {
	repo := NewUserRepository(setupUserRepoDB(t))

	require.NoError(t, repo.Create(makeUser("e@example.com", "echo")))

	results, err := repo.Search("zzznomatch")
	require.NoError(t, err)
	assert.Empty(t, results)
}
