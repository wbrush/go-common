package oauth2

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "errors"
    "strings"
)

func TokenDecode(token string) (claims ClaimSet, err error) {

    s := strings.Split(token, ".")
    if len(s) < 2 {
        err = errors.New("FastToken.decode(): invalid token received")
        return
    }

    // add back missing padding
    t := s[1]
    switch len(t) % 4 {
    case 1:
        t += "==="
    case 2:
        t += "=="
    case 3:
        t += "="
    }

    decoded, err := base64.URLEncoding.DecodeString(t)
    if (err != nil) {
        return
    }

    err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(&claims)
    if (err != nil) {
        //  not a valid token
        //logger.Debug().Println("DBG-> not a valid token; err: ", err.Error())
        //logger.Debug().Println("DBG-> tkn: ", token)
        return
    }

    return
}
