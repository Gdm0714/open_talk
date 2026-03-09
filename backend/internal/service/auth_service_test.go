package service

import (
	"errors"
	"testing"
	"time"

	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- in-process fake repository ---

type fakeUserRepo struct {
	users map[string]*model.User // keyed by email
	byID  map[string]*model.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users: make(map[string]*model.User),
		byID:  make(map[string]*model.User),
	}
}

func (r *fakeUserRepo) Create(user *model.User) error {
	if user.ID == "" {
		user.ID = "generated-id-" + user.Email
	}
	r.users[user.Email] = user
	r.byID[user.ID] = user
	return nil
}

func (r *fakeUserRepo) FindByEmail(email string) (*model.User, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (r *fakeUserRepo) FindByID(id string) (*model.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (r *fakeUserRepo) Update(user *model.User) error {
	r.users[user.Email] = user
	r.byID[user.ID] = user
	return nil
}

func (r *fakeUserRepo) Delete(id string) error {
	u, ok := r.byID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	delete(r.users, u.Email)
	delete(r.byID, id)
	return nil
}

func (r *fakeUserRepo) Search(query string) ([]model.User, error) {
	return nil, nil
}

// --- helpers ---

func testConfig() *config.Config {
	return &config.Config{JWTSecret: "test-secret-key-for-unit-tests"}
}

func newAuthSvc(repo *fakeUserRepo) AuthService {
	return NewAuthService(repo, testConfig())
}

// --- Register tests ---

func TestAuthService_Register_SuccessReturnsUserAndToken(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	user, token, err := svc.Register("valid@example.com", "password123", "validnick")

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, "valid@example.com", user.Email)
	assert.Equal(t, "validnick", user.Nickname)
}

func TestAuthService_Register_PasswordIsHashedInStoredUser(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	user, _, err := svc.Register("hash@example.com", "password123", "hashnick")

	require.NoError(t, err)
	assert.NotEqual(t, "password123", user.Password, "password should be stored hashed, not plaintext")
}

func TestAuthService_Register_ReturnsErrorOnInvalidEmail(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	_, _, err := svc.Register("not-an-email", "password123", "validnick")

	assert.ErrorIs(t, err, ErrInvalidEmail)
}

func TestAuthService_Register_ReturnsErrorOnShortPassword(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	_, _, err := svc.Register("ok@example.com", "short", "validnick")

	assert.ErrorIs(t, err, ErrInvalidPassword)
}

func TestAuthService_Register_ReturnsErrorOnInvalidNickname(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	_, _, err := svc.Register("ok@example.com", "password123", "x")

	assert.ErrorIs(t, err, ErrInvalidNickname)
}

func TestAuthService_Register_ReturnsErrorOnDuplicateEmail(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newAuthSvc(repo)

	_, _, err := svc.Register("dup@example.com", "password123", "firstnick")
	require.NoError(t, err)

	_, _, err = svc.Register("dup@example.com", "password123", "secondnick")
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

// --- Login tests ---

func TestAuthService_Login_SuccessReturnsUserAndToken(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())
	_, _, err := svc.Register("login@example.com", "password123", "loginnick")
	require.NoError(t, err)

	user, token, err := svc.Login("login@example.com", "password123")

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_ReturnsErrorOnWrongPassword(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())
	_, _, err := svc.Register("wp@example.com", "correctpass", "wpnick")
	require.NoError(t, err)

	_, _, err = svc.Login("wp@example.com", "wrongpass")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_ReturnsErrorWhenUserNotFound(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	_, _, err := svc.Login("ghost@example.com", "password123")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

// --- JWT generation and validation tests ---

func TestAuthService_GeneratedToken_ContainsUserIDAsSubject(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())
	user, token, err := svc.Register("jwt@example.com", "password123", "jwtnick")
	require.NoError(t, err)

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-unit-tests"), nil
	})

	require.NoError(t, err)
	assert.True(t, parsed.Valid)
	assert.Equal(t, user.ID, claims.Subject)
}

func TestAuthService_GeneratedToken_ExpiresInFuture(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())
	_, token, err := svc.Register("exp@example.com", "password123", "expnick")
	require.NoError(t, err)

	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-unit-tests"), nil
	})
	require.NoError(t, err)

	assert.True(t, claims.ExpiresAt.Time.After(time.Now()), "token should expire in the future")
}

func TestAuthService_RefreshToken_SuccessReturnsNewToken(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())
	_, original, err := svc.Register("refresh@example.com", "password123", "refreshnick")
	require.NoError(t, err)

	newToken, err := svc.RefreshToken(original)

	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
}

func TestAuthService_RefreshToken_ReturnsErrorOnInvalidToken(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	_, err := svc.RefreshToken("this.is.not.a.valid.token")

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestAuthService_RefreshToken_ReturnsErrorOnTokenSignedWithWrongSecret(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	wrongSecretClaims := jwt.RegisteredClaims{
		Subject:   "some-user-id",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, wrongSecretClaims)
	badToken, err := tok.SignedString([]byte("wrong-secret"))
	require.NoError(t, err)

	_, err = svc.RefreshToken(badToken)
	assert.True(t, errors.Is(err, ErrInvalidToken))
}

// --- ChangePassword tests ---

func TestChangePassword_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newAuthSvc(repo)

	_, _, err := svc.Register("cp@example.com", "oldpassword1", "cpnick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("cp@example.com")
	require.NoError(t, err)

	err = svc.ChangePassword(user.ID, "oldpassword1", "newpassword1")
	assert.NoError(t, err)
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newAuthSvc(repo)

	_, _, err := svc.Register("cpwrong@example.com", "correctpass1", "cpwrongnick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("cpwrong@example.com")
	require.NoError(t, err)

	err = svc.ChangePassword(user.ID, "wrongpass1", "newpassword1")
	assert.ErrorIs(t, err, ErrWrongPassword)
}

func TestChangePassword_InvalidNewPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newAuthSvc(repo)

	_, _, err := svc.Register("cpshort@example.com", "oldpassword1", "cpshortnick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("cpshort@example.com")
	require.NoError(t, err)

	err = svc.ChangePassword(user.ID, "oldpassword1", "short")
	assert.ErrorIs(t, err, ErrInvalidPassword)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	svc := newAuthSvc(newFakeUserRepo())

	err := svc.ChangePassword("nonexistent-id", "oldpassword1", "newpassword1")
	assert.Error(t, err)
}
