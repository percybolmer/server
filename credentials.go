package server

import (
	b64 "encoding/base64"
	"encoding/json"
	"reflect"
)

// Credentials is a container for different Credentials that can be used to connect to a Server
// Add new Credential types here
type Credentials struct {
	BasicAuth *basicAuthorization `json:"basicauth"`
	Ssh       *ssh                `json:"ssh"`
	DB        *db                 `json:"db"`
}

type Validator interface {
	validate() error
}

// NewCredentials is used to create a new Credentials object, and return a pointer to it
// @See LoadNewCredentials if you want to Load and Create Credentials at the same time
func NewCredentials() *Credentials {
	return &Credentials{}
}

// LoadNewCredentials will return a new Pointer to Credentials but Also load them with data
func LoadNewCredentials(data []byte) (*Credentials, error) {
	c := &Credentials{}
	return c, c.loadCredentials(data)
}

// loadCredentials takes a byte array containing JSON data and Unmarshalls it
// will return err if unmarshal fails, will also return error on validation failures
func (c *Credentials) loadCredentials(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		/* Replace this when xerrors is part of std lib and instead Wrap the unmarshal error with our custom Unmarshal err */
		return ErrFailedToUnmarshalCredentials
	}
	/* Validate all non Nil Credentials */
	v := reflect.ValueOf(c)
	/* c is a pointer so we need to indirect it */
	nopointer := reflect.Indirect(v)
	for i := 0; i < nopointer.NumField(); i++ {
		// See if field is nil
		if !nopointer.Field(i).IsNil() {
			/* Only validate if its part of the validator interface */
			interfacetype := nopointer.Field(i).Interface()

			originalField, ok := interfacetype.(Validator)
			if ok {
				err := originalField.validate()
				if err != nil {
					return err
				}
			}

		}
	}
	return nil

}

// basicAuthorization is a struct containing all neccessary data and functions to handle basic Auth
type basicAuthorization struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// Base64 is the base64 encoded version of the Username + Password
	Base64 string `json:"base64"`
}

// validate is used to enable Validator interface
func (b *basicAuthorization) validate() error {
	if (b.Username == "" || b.Password == "") && b.Base64 == "" {
		return ErrNeedUserAndPassword
	} else if b.Base64 == "" {
		/* Encode the b64 string if its not set in the Config */
		b.Base64 = b64.StdEncoding.EncodeToString([]byte(b.Username + ":" + b.Password))
	}
	return nil
}

// ssh contains neccessary data and methods to handle all SSH connection
type ssh struct {
	Publickey string `json:"publickey"`
}

// validate is used to enable validator interface
func (s *ssh) validate() error {
	if s.Publickey == "" {
		return ErrSshInvalidKey
	}
	return nil
}

type db struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"databasename"`
}

// validate makes db part of the Validator interface
func (d *db) validate() error {
	if d.Host == "" {
		return ErrInvalidHost
	} else if d.Username == "" || d.Password == "" {
		return ErrNeedUserAndPassword
	} else if d.DatabaseName == "" {
		return ErrInvalidDatabaseName
	} else if d.Port == 0 {
		return ErrInvalidDatabasePort
	}

	return nil
}
