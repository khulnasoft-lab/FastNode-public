package community

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestSignupManager() *SignupManager {
	db, err := gorm.Open("postgres", "postgres://XXXXXXX:XXXXXXX@localhost/XXXXXXX?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.SingularTable(true)
	db.DropTableIfExists(&Signup{})
	db.DropTableIfExists(&Download{})
	manager := NewSignupManager(db)
	if err := manager.Migrate(); err != nil {
		log.Fatalln(err)
	}
	return manager
}

func requireCleanupSignupManager(t *testing.T, m *SignupManager) {
	require.NoError(t, m.db.Close())
}

func buildServerClient() (*httptest.Server, *Server, *http.Client, *App) {
	ts, server, app := makeTestServer()

	client := makeTestClient()

	return ts, server, client, app
}

// --

func TestSignup(t *testing.T) {
	metadata := `{"a": "hello", "b": 5}`
	updatedMetadata := `{"a": "hello_updated", "b": 7}`

	manager := setupTestSignupManager()
	defer requireCleanupSignupManager(t, manager)

	_, err := manager.CreateOrUpdateSignup("test@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")
	signup, err := manager.Get("test@khulnasoft.com")
	require.Nil(t, err, "signup failed")
	require.Equal(t, "test@khulnasoft.com", signup.Email, "emails different")
	require.Equal(t, metadata, signup.Metadata, "metadata different")
	require.NotEqual(t, 0, signup.ID, "ID should be greater than zero")

	_, err = manager.CreateOrUpdateSignup("test@khulnasoft.com", updatedMetadata, "")
	require.Nil(t, err, "signup failed")
	updated, err := manager.Get("test@khulnasoft.com")
	require.Nil(t, err, "signup failed")
	require.NotEqual(t, 0, updated.ID, "ID should be greater than zero")
	require.Equal(t, "test@khulnasoft.com", updated.Email, "emails different")
	require.Equal(t, updatedMetadata, updated.Metadata, "metadata different")

	_, err = manager.CreateOrUpdateSignup("test2@khulnasoft.com", "blah blah", "")
	require.NotNil(t, err, "metadata is invalid, should have returned error")

	_, err = manager.CreateOrUpdateSignup("test2@khulnasoft.com", "{}", "192.168")
	require.NotNil(t, err, "ip is invalid, should have returned error")
}

func TestInvite(t *testing.T) {
	manager := setupTestSignupManager()
	defer requireCleanupSignupManager(t, manager)

	metadata := `{"a": "hello", "b": 5}`
	_, err := manager.CreateOrUpdateSignup("test@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")
	_, err = manager.Invite("test@khulnasoft.com", "testhost")
	require.Nil(t, err, "invite failed")
	signup, err := manager.Get("test@khulnasoft.com")
	require.NotEmpty(t, signup.InviteCode, "invite code is empty")
	require.NotEmpty(t, signup.InvitedTimestamp, "invite timestamp is empty")
	log.Println(signup.InviteCode)

	// lazy signup on invite
	_, err = manager.Invite("test2@khulnasoft.com", "testhost")
	require.Nil(t, err, "invite failed")
	signup, err = manager.Get("test2@khulnasoft.com")
	require.Nil(t, err, "lazy signup failed")
	require.Equal(t, "{}", signup.Metadata, "metadata should be empty")
	require.Empty(t, signup.ClientIP, "client ip should be empty")
	require.NotEmpty(t, signup.InviteCode, "invite code is empty")
	require.NotEmpty(t, signup.InvitedTimestamp, "invite timestamp is empty")
	log.Println(signup.InviteCode)
}

func TestHandleSignup(t *testing.T) {
	ts, _, client, app := buildServerClient()
	defer requireCleanupApp(t, app)

	url := makeTestURL(ts.URL, "/api/signups")

	metadata := `{"a": "hello", "b": 5}`
	resp, err := client.Post(url, "application/json", marshal(&Signup{
		Email:    "test@khulnasoft.com",
		Metadata: metadata,
		Secret:   defaultSecret,
	}))
	require.Nil(t, err, "POST new signup failed")
	require.Equal(t, 200, resp.StatusCode, "response status code not 200")

	signup, err := app.Signups.Get("test@khulnasoft.com")
	require.Equal(t, "test@khulnasoft.com", signup.Email, "emails different")
	require.Equal(t, metadata, signup.Metadata, "metadata different")
	require.NotEqual(t, 0, signup.ID, "ID should be greater than zero")
}

func TestHandleInviteAndValidate(t *testing.T) {
	ts, _, client, app := buildServerClient()
	defer requireCleanupApp(t, app)

	url := makeTestURL(ts.URL, "/api/invite")

	metadata := `{"a": "hello", "b": 5}`
	_, err := app.Signups.CreateOrUpdateSignup("test1@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")

	_, err = app.Signups.CreateOrUpdateSignup("test2@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")

	resp, err := client.Post(url, "application/json", marshal(&inviteData{
		Emails: []string{"test1@khulnasoft.com", "test2@khulnasoft.com"},
		Secret: defaultSecret,
	}))
	require.Nil(t, err, "POST new invites failed")
	require.Equal(t, 200, resp.StatusCode, "response status code not 200")

	_, err = app.Signups.CreateOrUpdateSignup("test3@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")

	resp, err = client.Post(url, "application/json", marshal(&inviteData{
		Emails: []string{"test1@khulnasoft.com", "test2@khulnasoft.com", "test3@khulnasoft.com"},
		Secret: defaultSecret,
	}))
	require.Nil(t, err, "POST new invites failed")
	require.Equal(t, 200, resp.StatusCode, "response status code not 200")

	signup, err := app.Signups.Get("test1@khulnasoft.com")
	require.Equal(t, "test1@khulnasoft.com", signup.Email, "emails different")
	require.NotEmpty(t, signup.InviteCode, "invite code is empty")
	require.NotEmpty(t, signup.InvitedTimestamp, "invite timestamp is empty")
	require.NotEqual(t, 0, signup.ID, "ID should be greater than zero")

	_, err = app.Signups.Validate(signup.InviteCode)
	require.Nil(t, err, "validation of invite code failed")

	signup, err = app.Signups.Get("test2@khulnasoft.com")
	require.Equal(t, "test2@khulnasoft.com", signup.Email, "emails different")
	require.NotEmpty(t, signup.InviteCode, "invite code is empty")
	require.NotEmpty(t, signup.InvitedTimestamp, "invite timestamp is empty")
	require.NotEqual(t, 0, signup.ID, "ID should be greater than zero")

	_, err = app.Signups.Validate(signup.InviteCode)
	require.Nil(t, err, "validation of invite code failed")

	signup, err = app.Signups.Get("test3@khulnasoft.com")
	require.Equal(t, "test3@khulnasoft.com", signup.Email, "emails different")
	require.NotEmpty(t, signup.InviteCode, "invite code is empty")
	require.NotEmpty(t, signup.InvitedTimestamp, "invite timestamp is empty")
	require.NotEqual(t, 0, signup.ID, "ID should be greater than zero")

	_, err = app.Signups.Validate(signup.InviteCode)
	require.Nil(t, err, "validation of invite code failed")

	_, err = app.Signups.Validate("random invite code")
	require.NotNil(t, err, "validation of random invite code succeeded")
}

func TestAll(t *testing.T) {
	manager := setupTestSignupManager()
	defer requireCleanupSignupManager(t, manager)

	metadata := `{"a": "hello", "b": 5}`

	for i := 0; i < 5; i++ {
		_, err := manager.CreateOrUpdateSignup(fmt.Sprintf("test%d@khulnasoft.com", i), metadata, "")
		require.Nil(t, err, "signup failed")
	}

	_, err := manager.Invite("test3@khulnasoft.com", "testhost")
	require.Nil(t, err, "invite failed")

	all, err := manager.All()
	require.Nil(t, err, "getting all signups failed")

	require.Equal(t, 5, len(all), "should be 5 total")
}

func TestDeduceClientIP(t *testing.T) {
	req, err := http.NewRequest("POST", "blah.khulnasoft.com/api/signups", strings.NewReader(""))
	require.NoError(t, err, "error creating new test request")
	req.Header.Add("X-Forwarded-For", "192.168.30.10, 192.168.30.20")

	ip := deduceClientIP(req)
	assert.Equal(t, "192.168.30.10", ip, "incorrect deduction of client IP")

	req, err = http.NewRequest("POST", "blah.khulnasoft.com/api/signups", strings.NewReader(""))
	assert.NoError(t, err, "error creating new test request")
	req.RemoteAddr = "192.168.30.50"

	ip = deduceClientIP(req)
	assert.Equal(t, "192.168.30.50", ip, "incorrect deduction of client IP")
}

func TestUnsubscribe(t *testing.T) {
	manager := setupTestSignupManager()
	defer requireCleanupSignupManager(t, manager)

	metadata := `{"a": "hello", "b": 5}`

	signup, err := manager.CreateOrUpdateSignup("test@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")
	assert.Equal(t, signup.Unsubscribed, false, "unsubscribed should default to false")

	signup, err = manager.Unsubscribe("test@khulnasoft.com")
	require.Nil(t, err, "unsubscribe failed")
	assert.Equal(t, signup.Unsubscribed, true, "unsubscribed should be true")

	signup, err = manager.Unsubscribe("test@khulnasoft.com")
	require.Nil(t, err, "unsubscribe failed")
	assert.Equal(t, signup.Unsubscribed, true, "unsubscribed should be true")
}

func TestSubscribe(t *testing.T) {
	manager := setupTestSignupManager()
	defer requireCleanupSignupManager(t, manager)

	metadata := `{"a": "hello", "b": 5}`

	signup, err := manager.CreateOrUpdateSignup("test@khulnasoft.com", metadata, "")
	require.Nil(t, err, "signup failed")
	assert.Equal(t, signup.Unsubscribed, false, "unsubscribed should default to false")

	signup, err = manager.Unsubscribe("test@khulnasoft.com")
	require.Nil(t, err, "unsubscribe failed")
	assert.Equal(t, signup.Unsubscribed, true, "unsubscribed should be true")

	signup, err = manager.Subscribe("test@khulnasoft.com")
	require.Nil(t, err, "subscribe failed")
	assert.Equal(t, signup.Unsubscribed, false, "unsubscribed should be false")

	signup, err = manager.Subscribe("test@khulnasoft.com")
	require.Nil(t, err, "subscribe failed")
	assert.Equal(t, signup.Unsubscribed, false, "unsubscribed should be false")
}

// --

func marshal(contents interface{}) io.Reader {
	marshalled, err := json.Marshal(contents)
	if err != nil {
		log.Fatalf("error marshalling contents: %v", err)
	}
	return bytes.NewBuffer(marshalled)
}
