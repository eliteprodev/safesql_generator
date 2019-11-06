package dinosql

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/kyleconroy/sqlc/internal/pg"
)

type GenerateSettings struct {
	Version   string            `json:"version"`
	Packages  []PackageSettings `json:"packages"`
	Overrides []Override        `json:"overrides,omitempty"`
	Rename    map[string]string `json:"rename,omitempty"`
}

type PackageSettings struct {
	Name                string     `json:"name"`
	Path                string     `json:"path"`
	Schema              string     `json:"schema"`
	Queries             string     `json:"queries"`
	EmitPreparedQueries bool       `json:"emit_prepared_queries"`
	EmitJSONTags        bool       `json:"emit_json_tags"`
	Overrides           []Override `json:"overrides"`
}

type Override struct {
	// name of the golang type to use, e.g. `github.com/segmentio/ksuid.KSUID`
	GoType string `json:"go_type"`

	// fully qualified name of the Go type, e.g. `github.com/segmentio/ksuid.KSUID`
	PostgresType string `json:"postgres_type"`

	// True if the GoType should override if the maching postgres type is nullable
	Null bool `json:"null"`

	// fully qualified name of the column, e.g. `accounts.id`
	Column string `json:"column"`

	columnName string
	table      pg.FQN
	goTypeName string
	goPackage  string
}

func (o *Override) Parse() error {
	// validate option combinations
	switch {
	case o.Column != "" && o.PostgresType != "":
		return fmt.Errorf("Override specifying both `column` (%q) and `postgres_type` (%q) is not valid.", o.Column, o.PostgresType)
	case o.Column == "" && o.PostgresType == "":
		return fmt.Errorf("Override must specify one of either `column` or `postgres_type`")
	}

	// validate Column
	if o.Column != "" {
		colParts := strings.Split(o.Column, ".")
		switch len(colParts) {
		case 2:
			o.columnName = colParts[1]
			o.table = pg.FQN{Schema: "public", Rel: colParts[0]}
		case 3:
			o.columnName = colParts[2]
			o.table = pg.FQN{Schema: colParts[0], Rel: colParts[1]}
		case 4:
			o.columnName = colParts[3]
			o.table = pg.FQN{Catalog: colParts[0], Schema: colParts[1], Rel: colParts[2]}
		default:
			return fmt.Errorf("Override `column` specifier %q is not the proper format, expected '[catalog.][schema.]colname.tablename'", o.Column)
		}
	}

	// validate GoType
	lastDot := strings.LastIndex(o.GoType, ".")
	if lastDot == -1 {
		return fmt.Errorf("Package override `go_type` specificier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", o.GoType)
	}
	lastSlash := strings.LastIndex(o.GoType, "/")
	if lastSlash == -1 {
		return fmt.Errorf("Package override `go_type` specificier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", o.GoType)
	}
	o.goTypeName = o.GoType[lastSlash+1:]
	o.goPackage = o.GoType[:lastDot]
	isPointer := o.GoType[0] == '*'
	if isPointer {
		o.goPackage = o.goPackage[1:]
		o.goTypeName = "*" + o.goTypeName
	}

	return nil
}

var ErrMissingVersion = errors.New("no version number")
var ErrUnknownVersion = errors.New("invalid version number")
var ErrNoPackages = errors.New("no packages")

func ParseConfig(rd io.Reader) (GenerateSettings, error) {
	dec := json.NewDecoder(rd)
	dec.DisallowUnknownFields()
	var config GenerateSettings
	if err := dec.Decode(&config); err != nil {
		return config, err
	}
	if config.Version == "" {
		return config, ErrMissingVersion
	}
	if config.Version != "1" {
		return config, ErrUnknownVersion
	}
	if len(config.Packages) == 0 {
		return config, ErrNoPackages
	}
	for i := range config.Overrides {
		if err := config.Overrides[i].Parse(); err != nil {
			return config, err
		}
	}
	for j := range config.Packages {
		for i := range config.Packages[j].Overrides {
			if err := config.Packages[j].Overrides[i].Parse(); err != nil {
				return config, err
			}
		}
	}
	return config, nil
}
