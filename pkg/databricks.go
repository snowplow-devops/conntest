package pkg

import "github.com/xo/dburl"

// Functions added below to support Databricks which is not currently supported in xo/dburl
// and xo/usql - the Databricks driver is still in beta
func RegisterDatabricks() {
	dburl.Register(dburl.Scheme{"databricks", GenDatabricks, 0, false, []string{"databricks"}, ""})
}

// GenDatabricks generates a databricks DSN from the passed URL.
// Format is here https://github.com/databricks/databricks-sql-go#usage
//
// databricks://:[your token]@[Workspace hostname][Endpoint HTTP Path]
func GenDatabricks(u *dburl.URL) (string, error) {
	host := u.Hostname()
	if host == "" {
		return "", dburl.ErrMissingHost
	}

	// add auth token, which is passed as the password
	user := ""
	if pass, _ := u.User.Password(); pass != "" {
		user += ":" + pass
	}

	return "databricks://" + user + "@" + host + u.Path, nil
}

func GetProtocols(name string) []string {
	return dburl.Protocols(name)
}

func SchemeDriverAndAliases(name string) (string, []string) {
	return dburl.SchemeDriverAndAliases(name)
}
