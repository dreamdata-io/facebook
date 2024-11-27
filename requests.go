package facebook

import (
	"github.com/dreamdata-io/facebook/internal"
	"strings"
)

type Result = internal.Result
type Method = internal.Method
type Params = internal.Params

func FieldsParams(fields ...string) Params {
	return internal.MakeParams(map[string]string{
		"fields": strings.Join(fields, ","),
	})
}
