package service

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// fakeUserRepoWithSearch extends fakeUserRepo to return meaningful Search results.
type fakeUserRepoWithSearch struct {
	fakeUserRepo
	searchResult []model.User
}

func (r *fakeUserRepoWithSearch) Search(query string) ([]model.User, error) {
	return r.searchResult, nil
}

func newUserSvcWithSearch(repo *fakeUserRepoWithSearch) UserService {
	return NewUserService(repo, newFakeFriendRepo(), &fakeChatRepo{})
}

// --- GetProfile tests ---

func TestGetProfile_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newUserSvc(repo)

	authSvc := newAuthSvc(repo)
	_, _, err := authSvc.Register("profile@example.com", "password123", "profilenick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("profile@example.com")
	require.NoError(t, err)

	got, err := svc.GetProfile(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, "profile@example.com", got.Email)
}

func TestGetProfile_NotFound(t *testing.T) {
	svc := newUserSvc(newFakeUserRepo())

	_, err := svc.GetProfile("nonexistent-id")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

// --- UpdateProfile tests ---

func TestUpdateProfile_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newUserSvc(repo)

	authSvc := newAuthSvc(repo)
	_, _, err := authSvc.Register("update@example.com", "password123", "oldnick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("update@example.com")
	require.NoError(t, err)

	newNick := "newnick"
	avatarURL := "https://example.com/avatar.png"
	statusMsg := "hello world"

	updated, err := svc.UpdateProfile(user.ID, &newNick, &avatarURL, &statusMsg)
	require.NoError(t, err)
	assert.Equal(t, "newnick", updated.Nickname)
	assert.Equal(t, "https://example.com/avatar.png", updated.AvatarURL)
	assert.Equal(t, "hello world", updated.StatusMessage)
}

func TestUpdateProfile_PartialUpdate(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newUserSvc(repo)

	authSvc := newAuthSvc(repo)
	_, _, err := authSvc.Register("partial@example.com", "password123", "originalnick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("partial@example.com")
	require.NoError(t, err)

	newNick := "updatednick"
	// only update nickname, leave avatar and status nil
	updated, err := svc.UpdateProfile(user.ID, &newNick, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "updatednick", updated.Nickname)
	assert.Empty(t, updated.AvatarURL)
}

func TestUpdateProfile_NotFound(t *testing.T) {
	svc := newUserSvc(newFakeUserRepo())

	nick := "nick"
	_, err := svc.UpdateProfile("nonexistent-id", &nick, nil, nil)
	assert.Error(t, err)
}

// --- SearchUsers tests ---

func TestSearchUsers_Success(t *testing.T) {
	base := newFakeUserRepo()
	repo := &fakeUserRepoWithSearch{
		fakeUserRepo: *base,
		searchResult: []model.User{
			{ID: "u1", Nickname: "alice"},
			{ID: "u2", Nickname: "alicia"},
		},
	}
	svc := newUserSvcWithSearch(repo)

	results, err := svc.SearchUsers("ali")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestSearchUsers_Empty(t *testing.T) {
	base := newFakeUserRepo()
	repo := &fakeUserRepoWithSearch{
		fakeUserRepo: *base,
		searchResult: []model.User{},
	}
	svc := newUserSvcWithSearch(repo)

	results, err := svc.SearchUsers("nobody")
	require.NoError(t, err)
	assert.Empty(t, results)
}
