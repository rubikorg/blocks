# 📨 mailer rubik block

This is a wrapper block implementation of 
[jordan-wright/email](https://github.com/jordan-wright/email).

### Config

```toml
[mailer]
auth = true
host = rubik.ashishshekar.com
username = someuser@email.com
password = some_secret_password
```

### Import

```go
import (
    _ "github.com/rubikorg/blocks/mailer"
)
```

### Send Email

```go
import (
    "github.com/rubikorg/blocks/mailer"
    r "github.com/rubikorg/rubik"
)

func ctl(en interface{}) r.ByteResponse {
    d := mailer.Details{
        Subject: "Hello",
        Body: "World",
    }
    err := mailer.Send(d)
    if err != nil {
        // do something
    }
}
```

### Credits:

The entire implementation of core code: [jordan-wright](https://github.com/jordan-wright)

