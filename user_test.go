package imapsql

import (
	"testing"

	"gotest.tools/assert"
)

func TestUserCaseInsensitivity(t *testing.T) {
	b := initTestBackend().(*Backend)
	defer cleanBackend(b)

	assert.NilError(t, b.CreateUser("foXcpp"))
	_, err := b.Login(nil, "foXCpp", "")
	assert.NilError(t, err, "b.Login")
	u1, err := b.GetUser("Foxcpp")
	assert.NilError(t, err, "b.GetUser")
	u2, err := b.GetOrCreateUser("FOXcpp")
	assert.NilError(t, err, "b.GetOrCreateUser")

	assert.NilError(t, u1.CreateMailbox("BOX"))
	_, mbox, err := u2.GetMailbox("BOX", true, &noopConn{})
	assert.NilError(t, err, "u2.GetMailbox")
	defer mbox.Close()
}

func TestInboxCreation(t *testing.T) {
	b := initTestBackend().(*Backend)
	defer cleanBackend(b)

	assert.NilError(t, b.CreateUser("foxcpp"))

	u, err := b.GetUser("foxcpp")
	assert.NilError(t, err, "b.GetUser")

	_, mbox, err := u.GetMailbox("INBOX", true, &noopConn{})
	assert.NilError(t, err, "u.GetMailbox")
	defer mbox.Close()
}
