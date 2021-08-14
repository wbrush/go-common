package oauth2

import (
	"fmt"
	"github.com/wbrush/go-common/config"
	"strings"
)

type UserTokenData struct {
	UserId                        int64
	FirstName, LastName, UserName string
}

func GetUserInfoFromToken(conf *config.ServiceParams, token string) (d *UserTokenData, err error) {
	d = &UserTokenData{
		UserId: -1,
	}
	if conf == nil {
		return d, fmt.Errorf("nil config provided")
	}

	if strings.TrimSpace(token) != "" {
		claims, err := TokenDecode(token[6:])
		if err != nil {
			return d, fmt.Errorf("can't decode token: %s", err.Error())
		}
		d.FirstName = claims.UserFirstName
		d.LastName = claims.UserLastName
		d.UserName = claims.Sub
		if conf.Environment == config.EnvironmentTypeProd {
			d.UserId = claims.ProdUserId
		} else if conf.Environment == config.EnvironmentTypeTest {
			d.UserId = claims.TestUserID
		} else {
			d.UserId = claims.DevUserID
		}
	}
	return
}
