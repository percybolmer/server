package server

import (
	"testing"
)

func TestLoadCredentials(t *testing.T) {
	type credentialTests struct {
		d   []byte
		err error
		n   string
	}

	cases := []credentialTests{
		{[]byte(``), ErrFailedToUnmarshalCredentials, "Unmarshal error"},
		{[]byte(`{ "basicauth": { "username": "percy"} }`), ErrNeedUserAndPassword, "basicAuthValidation"},
		{[]byte(`{ "basicauth": { "username": "percy", "base64": "222=="} }`), nil, "basicAuthBase64Check"},
		{[]byte(`{ "basicauth": { "username": "percy", "password": "test"} }`), nil, "basicAuthWorking"},
		{[]byte(`{ "db": { "username": "test", "password": "test"} }`), ErrInvalidHost, "db_invalidHost"},
		{[]byte(`{ "db": { "host": "123", "username": "", "password": "test"} }`), ErrNeedUserAndPassword, "db_InvalidUserNamePassword"},
		{[]byte(`{ "db": { "host": "123", "username": "percy", "password": "test", "databasename": ""} }`), ErrInvalidDatabaseName, "db_InvalidDataBaseName"},
		{[]byte(`{ "db": { "host": "123", "username": "percy", "password": "test", "databasename": "work"} }`), ErrInvalidDatabasePort, "db_InvalidDataBasePort"},
		{[]byte(`{ "db": { "host": "123", "username": "oera", "password": "test", "databasename": "test", "port": 80} }`), nil, "db_shouldWork"},
	}

	for _, tc := range cases {
		c := NewCredentials()
		err := c.loadCredentials(tc.d)
		if err != tc.err {
			t.Errorf("Failed test named %s, given result was: %v", tc.n, err)
		}
	}

	// Repeat to test LoadNewCredentials at the same time
	for _, tc := range cases {
		c, err := LoadNewCredentials(tc.d)
		if err != tc.err {
			t.Errorf("Failed test named %s, given result was: %v", tc.n, err)
		} else if c == nil {
			t.Errorf("Credentials should not be nil")
		}
	}

}
