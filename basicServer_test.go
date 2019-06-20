package server

import (
	"testing"
)

var (
	workingDB = []byte(`{ "db": { "host": "localhost", "username": "postgres", "password": "509710", "databasename": "basicserver", "port": 5432} }`)
)

func getWorkingBasicServer() *basicServer {
	b := NewBasicServer("t")

	cred, err := LoadNewCredentials(workingDB)
	if err != nil {
		panic(err)
	}
	b.credentials = cred
	return b
}
func TestNewBasicServer(t *testing.T) {
	b := NewBasicServer("t")
	if b == nil {
		t.Fatal("This should not be able to fail")
	}
}
func TestConnect(t *testing.T) {
	b := getWorkingBasicServer()

	type connectionTests struct {
		ip  string
		err error
		n   string
	}

	cases := []connectionTests{
		{ip: "", err: ErrFailedToConnect, n: "Should fail to connect with a Bad IP"},
		{ip: "localhost", err: nil, n: "Should Work"},
	}

	for _, tc := range cases {
		b.credentials.DB.Host = tc.ip
		err := b.Connect(b.credentials)
		if err != tc.err {
			t.Errorf("Expected: %s  but got: %s", tc.err, err)
		}
	}
	/* Also try speicla case where NIL credentials shuold fail */
	b.credentials.DB = nil
	err := b.Connect(b.credentials)
	if err != ErrMissingCredentials {
		t.Error("Expected: ", ErrMissingCredentials, " But got:", err)
	}

}
func TestDisconnect(t *testing.T) {
	b := getWorkingBasicServer()

	/* Test closing Nil Database */
	err := b.Disconnect()
	if err != ErrCannotCloseNilDatabase {
		t.Error(err)
	}

	err = b.Connect(b.credentials)
	if err != nil {
		t.Error(err)
	}
	/* Now test to disconnect without errors */
	err = b.Disconnect()
	if err != nil {
		t.Error(err)
	}
}

func TestInterfaceImplementation(t *testing.T) {
	var _ Server = (*basicServer)(nil)
}
