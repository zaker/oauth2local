package register

import (
	"fmt"
)

var _ = `
on open location myUrl

	do shell script "oauth2local callback " & "\"" & myUrl & "\""

end open location
`

func RegMe(urlScheme, locauthPath string) error {

	return fmt.Errorf("not implemented")
}
