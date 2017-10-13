gorails/session
===============

[![Build Status](https://travis-ci.org/adjust/gorails.png)](https://travis-ci.org/adjust/gorails)

## Installation

With Go and git installed:

```
go get -u github.com/adjust/gorails/session
```

## Usage

```go
import "github.com/adjust/gorails/session"

// sessionCookie - raw _<your app name>_session cookie
func getRailsSessionData(sessionCookie string) (decryptedCookieData []byte, err error) {
  decryptedCookieData, err = session.DecryptSignedCookie(sessionCookie, secretKeyBase, salt, signSalt)

  return
}

const (
  secretKeyBase = "..."                      // can be found in config/initializers/secret_token.rb
  salt          = "encrypted cookie"         // default value for Rails 4 app
  signSalt      = "signed encrypted cookie"  // default value for Rails 4 app
)
```

After you decrypted session data you might like to deserialize it using [gorails/marshal](https://github.com/adjust/gorails/tree/master/marshal) if your Rails version is less than v4.1 and you use the default serializer config.

Rails use JSON as its default serializer from v4.1, so you can deserialize the decrypted session data as a common JSON data as what [test](https://github.com/adjust/gorails/blob/master/session/session_test.go) does.
