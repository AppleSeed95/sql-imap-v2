package testsuite

import (
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-imap"
	appendlimit "github.com/emersion/go-imap-appendlimit"
	"github.com/emersion/go-imap/backend"
	imapsql "github.com/foxcpp/go-imap-sql"
	"gotest.tools/assert"
)

func Backend_AppendLimit(t *testing.T, newBack NewBackFunc, closeBack CloseBackFunc) {
	b := newBack()
	defer closeBack(b)

	bAL, ok := b.(imapsql.AppendLimitBackend)
	if !ok {
		t.Skip("APPENDLIMIT extension is not implemented (need imapsql.AppendLimitBackend interface)")
		t.SkipNow()
	}

	u := getUser(t, b)
	defer assert.NilError(t, u.Logout())

	t.Run("No Limit", func(t *testing.T) {
		assert.NilError(t, bAL.SetMessageLimit(nil))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 300)))
		assert.NilError(t, err)
	})
	t.Run("Under Limit", func(t *testing.T) {
		lim := uint32(500)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 300)))
		assert.NilError(t, err)
	})
	t.Run("Over Limit", func(t *testing.T) {
		lim := uint32(500)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
}

func User_AppendLimit(t *testing.T, newBack NewBackFunc, closeBack CloseBackFunc) {
	b := newBack()
	defer closeBack(b)
	u := getUser(t, b)
	defer assert.NilError(t, u.Logout())

	bAL, ok := b.(imapsql.AppendLimitBackend)
	if !ok {
		t.Skip("APPENDLIMIT extension is not implemented (need imapsql.AppendLimitBackend interface)")
		t.SkipNow()
	}
	uAL, ok := u.(imapsql.AppendLimitUser)
	if !ok {
		t.Skip("APPENDLIMIT extension is not implemented (need imapsql.AppendLimitUser interface)")
		t.SkipNow()
	}

	t.Run("No Limit", func(t *testing.T) {
		assert.NilError(t, uAL.SetMessageLimit(nil))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 300)))
		assert.NilError(t, err)
	})
	t.Run("Under Limit", func(t *testing.T) {
		lim := uint32(500)
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 300)))
		assert.NilError(t, err)
	})
	t.Run("Over Limit", func(t *testing.T) {
		lim := uint32(500)
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
	t.Run("Override backend - Under Limit", func(t *testing.T) {
		lim := uint32(100)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		lim = 500
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 400)))
		assert.NilError(t, err)
	})
	t.Run("Override backend - Over Limit", func(t *testing.T) {
		lim := uint32(1000)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		lim = 500
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
}

func Mailbox_AppendLimit(t *testing.T, newBack NewBackFunc, closeBack CloseBackFunc) {
	b := newBack()
	defer closeBack(b)
	u := getUser(t, b)
	defer assert.NilError(t, u.Logout())

	bAL, ok := b.(imapsql.AppendLimitBackend)
	if !ok {
		t.Skip("APPENDLIMIT extension is not implemented (need imapsql.AppendLimitBackend interface)")
		t.SkipNow()
	}
	uAL, ok := u.(imapsql.AppendLimitUser)
	if !ok {
		t.Skip("APPENDLIMIT extension is not implemented (need imapsql.AppendLimitUser interface)")
		t.SkipNow()
	}

	setMboxLim := func(t *testing.T, mbox backend.Mailbox, val uint32) {
		mAL, ok := mbox.(imapsql.AppendLimitMbox)
		if !ok {
			t.Skip("APPENDLIMIT extension is not implemented (need imapsql.AppendLimitMbox inteface)")
			t.SkipNow()
		}
		assert.NilError(t, mAL.SetMessageLimit(&val))
	}

	t.Run("No Limit", func(t *testing.T) {
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 300)))
		assert.NilError(t, err)
	})
	t.Run("Under Limit", func(t *testing.T) {
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 300)))
		assert.NilError(t, err)
	})
	t.Run("Over Limit", func(t *testing.T) {
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
	t.Run("Override backend - Under Limit", func(t *testing.T) {
		lim := uint32(100)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 400)))
		assert.NilError(t, err)
	})
	t.Run("Override backend - Over Limit", func(t *testing.T) {
		lim := uint32(1000)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		lim = 500
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
	t.Run("Override user - Under Limit", func(t *testing.T) {
		lim := uint32(100)
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 400)))
		assert.NilError(t, err)
	})
	t.Run("Override user - Over Limit", func(t *testing.T) {
		lim := uint32(1000)
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		lim = 500
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
	t.Run("Override backend & user - Under Limit", func(t *testing.T) {
		lim := uint32(200)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		lim = 1000
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		lim = 100
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 400)))
		assert.NilError(t, err)
	})
	t.Run("Override backend & user - Over Limit", func(t *testing.T) {
		lim := uint32(2000)
		assert.NilError(t, bAL.SetMessageLimit(&lim))
		lim = 1000
		assert.NilError(t, uAL.SetMessageLimit(&lim))
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		err := mbox.CreateMessage([]string{}, time.Now(), strings.NewReader(strings.Repeat("A", 700)))
		assert.Error(t, err, appendlimit.ErrTooBig.Error())
	})
	t.Run("Status - No Limit", func(t *testing.T) {
		mbox := getMbox(t, u)

		status, err := mbox.Status([]imap.StatusItem{appendlimit.StatusAppendLimit})
		assert.NilError(t, err)

		assert.Equal(t, appendlimit.MailboxStatusAppendLimit(status), (*uint32)(nil), "Non-nil value for limit")
	})
	t.Run("Status - Limit Present", func(t *testing.T) {
		mbox := getMbox(t, u)
		setMboxLim(t, mbox, 500)

		status, err := mbox.Status([]imap.StatusItem{appendlimit.StatusAppendLimit})
		assert.NilError(t, err)

		assert.Assert(t, appendlimit.MailboxStatusAppendLimit(status) != nil, "Nil value for limit item")
		val := *appendlimit.MailboxStatusAppendLimit(status)
		assert.Equal(t, val, uint32(500), "Wrong value for status item")
	})
}
