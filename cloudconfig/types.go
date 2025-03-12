package cloudconfig

type Format string

var (
	JSONFormat Format = "json"
	YAMLFormat Format = "yaml"
)

func (f Format) Valid() bool {
	return f == JSONFormat || f == YAMLFormat
}
