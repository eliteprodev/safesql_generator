package dinosql

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const missingVersion = `{
}`

const missingPackages = `{
  "version": "1"
}`

const unknownVersion = `{
  "version": "foo"
}`

const unknownFields = `{
  "foo": "bar"
}`

func TestBadConfigs(t *testing.T) {
	for _, test := range []struct {
		name string
		err  string
		json string
	}{
		{
			"missing version",
			"no version number",
			missingVersion,
		},
		{
			"missing packages",
			"no packages",
			missingPackages,
		},
		{
			"unknown version",
			"invalid version number",
			unknownVersion,
		},
		{
			"unknown fields",
			"json: unknown field \"foo\"",
			unknownFields,
		},
	} {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfig(strings.NewReader(tt.json))
			if err == nil {
				t.Fatalf("expected err; got nil")
			}
			if diff := cmp.Diff(err.Error(), tt.err); diff != "" {
				t.Errorf("differed (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTypeOverrides(t *testing.T) {
	for _, test := range []struct {
		override Override
		pkg      string
		typeName string
		basic    bool
	}{
		{
			Override{
				PostgresType: "uuid",
				GoType:       "github.com/segmentio/ksuid.KSUID",
			},
			"github.com/segmentio/ksuid",
			"ksuid.KSUID",
			false,
		},
		// TODO: Add test for struct pointers
		//
		// {
		// 	Override{
		// 		PostgresType: "uuid",
		// 		GoType:       "github.com/segmentio/*ksuid.KSUID",
		// 	},
		// 	"github.com/segmentio/ksuid",
		// 	"*ksuid.KSUID",
		// 	false,
		// },
		{
			Override{
				PostgresType: "citext",
				GoType:       "string",
			},
			"",
			"string",
			true,
		},
	} {
		tt := test
		t.Run(tt.override.GoType, func(t *testing.T) {
			if err := tt.override.Parse(); err != nil {
				t.Fatalf("override parsing failed; %s", err)
			}
			if diff := cmp.Diff(tt.typeName, tt.override.goTypeName); diff != "" {
				t.Errorf("type name mismatch;\n%s", diff)
			}
			if diff := cmp.Diff(tt.pkg, tt.override.goPackage); diff != "" {
				t.Errorf("package mismatch;\n%s", diff)
			}
			if diff := cmp.Diff(tt.basic, tt.override.goBasicType); diff != "" {
				t.Errorf("basic mismatch;\n%s", diff)
			}
		})
	}
	for _, test := range []struct {
		override Override
		err      string
	}{
		{
			Override{
				PostgresType: "uuid",
				GoType:       "Pointer",
			},
			"Package override `go_type` specifier \"Pointer\" is not a Go basic type e.g. 'string'",
		},
		{
			Override{
				PostgresType: "uuid",
				GoType:       "untyped rune",
			},
			"Package override `go_type` specifier \"untyped rune\" is not a Go basic type e.g. 'string'",
		},
	} {
		tt := test
		t.Run(tt.override.GoType, func(t *testing.T) {
			err := tt.override.Parse()
			if err == nil {
				t.Fatalf("expected pars to fail; got nil")
			}
			if diff := cmp.Diff(tt.err, err.Error()); diff != "" {
				t.Errorf("error mismatch;\n%s", diff)
			}
		})
	}
}
