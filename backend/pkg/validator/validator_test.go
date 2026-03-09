package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ValidateEmail tests

func TestValidateEmail_AcceptsStandardEmail(t *testing.T) {
	assert.True(t, ValidateEmail("user@example.com"))
}

func TestValidateEmail_AcceptsEmailWithSubdomain(t *testing.T) {
	assert.True(t, ValidateEmail("user@mail.example.com"))
}

func TestValidateEmail_AcceptsEmailWithPlusSign(t *testing.T) {
	assert.True(t, ValidateEmail("user+tag@example.com"))
}

func TestValidateEmail_AcceptsEmailWithDots(t *testing.T) {
	assert.True(t, ValidateEmail("first.last@example.org"))
}

func TestValidateEmail_AcceptsEmailWithNumbers(t *testing.T) {
	assert.True(t, ValidateEmail("user123@example123.com"))
}

func TestValidateEmail_RejectsEmailWithoutAtSign(t *testing.T) {
	assert.False(t, ValidateEmail("userexample.com"))
}

func TestValidateEmail_RejectsEmailWithoutDomain(t *testing.T) {
	assert.False(t, ValidateEmail("user@"))
}

func TestValidateEmail_RejectsEmailWithoutLocalPart(t *testing.T) {
	assert.False(t, ValidateEmail("@example.com"))
}

func TestValidateEmail_RejectsEmptyString(t *testing.T) {
	assert.False(t, ValidateEmail(""))
}

func TestValidateEmail_RejectsEmailWithoutTLD(t *testing.T) {
	assert.False(t, ValidateEmail("user@example"))
}

func TestValidateEmail_RejectsEmailWithSpaces(t *testing.T) {
	assert.False(t, ValidateEmail("user @example.com"))
}

// ValidatePassword tests

func TestValidatePassword_AcceptsEightCharacterPassword(t *testing.T) {
	assert.True(t, ValidatePassword("12345678"))
}

func TestValidatePassword_AcceptsPasswordLongerThanEightChars(t *testing.T) {
	assert.True(t, ValidatePassword("supersecretpassword"))
}

func TestValidatePassword_AcceptsPasswordWithUnicodeRunes(t *testing.T) {
	// "가나다라마바사아" is 8 Korean runes — should pass
	assert.True(t, ValidatePassword("가나다라마바사아"))
}

func TestValidatePassword_RejectsSevenCharacterPassword(t *testing.T) {
	assert.False(t, ValidatePassword("1234567"))
}

func TestValidatePassword_RejectsSingleCharacterPassword(t *testing.T) {
	assert.False(t, ValidatePassword("a"))
}

func TestValidatePassword_RejectsEmptyString(t *testing.T) {
	assert.False(t, ValidatePassword(""))
}

// ValidateNickname tests

func TestValidateNickname_AcceptsTwoCharacterNickname(t *testing.T) {
	assert.True(t, ValidateNickname("ab"))
}

func TestValidateNickname_AcceptsTwentyCharacterNickname(t *testing.T) {
	assert.True(t, ValidateNickname("12345678901234567890"))
}

func TestValidateNickname_AcceptsNicknameInValidRange(t *testing.T) {
	assert.True(t, ValidateNickname("JohnDoe"))
}

func TestValidateNickname_AcceptsUnicodeNickname(t *testing.T) {
	// "홍길동" is 3 Korean runes — valid
	assert.True(t, ValidateNickname("홍길동"))
}

func TestValidateNickname_RejectsSingleCharacterNickname(t *testing.T) {
	assert.False(t, ValidateNickname("a"))
}

func TestValidateNickname_RejectsTwentyOneCharacterNickname(t *testing.T) {
	assert.False(t, ValidateNickname("123456789012345678901"))
}

func TestValidateNickname_RejectsEmptyString(t *testing.T) {
	assert.False(t, ValidateNickname(""))
}
