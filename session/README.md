gorails/session
===============

[![Build Status](https://travis-ci.org/adjust/gorails.png)](https://travis-ci.org/adjust/gorails)

## Installation

With Go and git installed:

```
go get -u code.google.com/p/go.crypto/pbkdf2
go get -u github.com/adjust/gorails/session
```

Or you can use [Goem](http://big-elephants.com/2013-09/goem-the-missing-go-extension-manager/).

## Usage

```go
import "github.com/adjust/gorails/session"

// session_cookie - raw _<your app name>_session cookie
func getRailsSessionData(session_cookie string) (decrypted_cookie_data []byte, err error) {
  decrypted_cookie_data, err = session.DecryptSignedCookie(session_cookie, secret_key_base, salt)

  return
}

const (
  secret_key_base = "..." // can be found in config/initializers/secret_token.rb
  salt = "encrypted cookie" // default value for Rails 4 app
)
```

After you decrypted session data you might like to deserialize it using [gorails/marshal](https://github.com/adjust/gorails/tree/master/marshal)
