{{ define "gtime_test.go" }}
package gtime

import (
  "database/sql/driver"
  "database/sql"
)

func init () {
  var _ = []driver.Valuer{
    {{- range .Formats }}
    Time{{.Name}}{},
    {{- end }}
  }

  var _ = []sql.Scanner{
    {{- range .Formats }}
    &Time{{.Name}}{},
    {{- end }}
  }
}
{{ end }}