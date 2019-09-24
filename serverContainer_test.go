package server

import (
	"errors"
	"testing"
	"time"

	"github.com/percybolmer/credentials"
)

func getWorkingBasicServer(id string) basicserver {
	return basicserver{
		id: id,
	}
}

// basicserver is just a make up struct used to test server functionality
type basicserver struct {
	id string
}

var pingCounter int

func (b basicserver) Connect(c *credentials.Credentials) error {
	return nil
}
func (b basicserver) Reconnect() error {
	return nil
}
func (b basicserver) Disconnect() error     { return nil }
func (b basicserver) TestConnection() error { return nil }
func (b basicserver) Ping() error {
	if pingCounter == 1 {
		return errors.New("Should return an error")
	}
	pingCounter++
	return nil
}

// GetUniqueIdentifier is responsible to return a Unique identifer for each Server object
func (b basicserver) GetUniqueIdentifier() string { return b.id }

func TestAdd(t *testing.T) {
	/* Test if Duplicates are handled correctly */
	sc := NewServerContainer()

	bs := getWorkingBasicServer("test")

	/* Try Adding BS twice to the ServerContainer, second tiem should fail */
	err := sc.Add(bs)
	if err != nil {
		t.Error(err)
	}

	err = sc.Add(bs)
	if err == nil {
		t.Error("Should have failed to add the same server twice")
	}
}

func TestRemove(t *testing.T) {
	/* Test if Duplicates are handled correctly */
	sc := NewServerContainer()

	bs := getWorkingBasicServer("test")

	/* Try Adding BS twice to the ServerContainer, second tiem should fail */
	err := sc.Add(bs)
	if err != nil {
		t.Error(err)
	}

	sc.Remove(bs)
	if len(sc.servers) != 0 {
		t.Error("Failed to remove server")
	}
}

func TestManageConnections(t *testing.T) {
	// No idea how I should test this properly, Here I start a select waiting for cancel after X seconds,

	sc := NewServerContainer()

	bs := getWorkingBasicServer("test")
	bs2 := getWorkingBasicServer("t2")
	testErr(sc.Add(bs), t)
	testErr(sc.Add(bs2), t)

	if len(sc.servers) != 2 {
		t.Fatal("Failed to properly add testing servers")
	}
	sc.SetConnectionInterval(5)
	sc.ManageConnections()
	/* This Go func will Wait a few Seconds and then Disconnect the Server */
	pingCounter = 0
	go func() {
		time.Sleep(7 * time.Second)
		sc.Err <- sc.servers["t2"].Disconnect()
		sc.SetConnectionInterval(10)
	}()
	/* Kill chan is used to cancel for loop */
	killChan := make(chan int)
	go func() {
		time.Sleep(25 * time.Second)
		killChan <- 1
	}()
	for {
		select {
		case err := <-sc.Err:
			if err != nil {
				t.Fatal(err)
			}
		case <-killChan:
			/* The server should now be able to Ping correctly */
			err := sc.servers["t2"].Ping()
			if err != nil {
				t.Fatal(err)
			}
			/* Stop managing and Disconnect a server, Wait and it should NOT be connected again*/
			sc.StopManagingConnections()
			sc.servers["t2"].Disconnect()
			time.Sleep(11 * time.Second)
			err = sc.servers["t2"].Ping()
			if err == nil {
				t.Fatal("This database should have been closed")
			}

			return
		}
	}

}

func testErr(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
