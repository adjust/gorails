gorails/marshal
===============

[![Build Status](https://travis-ci.org/adjust/gorails.png)](https://travis-ci.org/adjust/gorails)

This is a parser/unmarshaller for objects serialized via [Ruby's
Marshal#dump][marshal-dump]. The binary format is [described in some
depth][marshal-format], if there's need for future extension.

## Installation

With Go and git installed:

```
go get -u github.com/adjust/gorails/marshal
```

## Usage

```go
import (
  "errors"
  "github.com/adjust/gorails/marshal"
)

func getAuthUserId(decrypted_session_data []byte) (user_id int64, err error) {
  unauthorized_user := errors.New("Unauthorized user")
  invalid_auth_data := errors.New("Invalid auth data")

  session_data, err := marshal.CreateMarshalledObject(decrypted_session_data).GetAsMap()
  if err != nil {
    return
  }

  warden_data, ok := session_data["warden.user.user.key"]
  if !ok {
    return 0, unauthorized_user
  }

  warden_user_key, err := warden_data.GetAsArray()
  if err != nil {
    return
  }
  if len(warden_user_key) < 1 {
    return 0, invalid_auth_data
  }

  user_data, err := warden_user_key[0].GetAsArray()
  if err != nil {
    return
  }
  if len(user_data) < 1 {
    return 0, invalid_auth_data
  }

  user_id, err = user_data[0].GetAsInteger()

  return
}
```

[marshal-dump]: http://ruby-doc.org/core-2.5.0/Marshal.html#method-c-dump
[marshal-format]: https://docs.ruby-lang.org/en/2.1.0/marshal_rdoc.html
