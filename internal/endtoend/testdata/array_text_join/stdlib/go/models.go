// Code generated by sqlc. DO NOT EDIT.

package querytest

import ()

type Bar struct {
	ID   string
	Info []string
}

type Foo struct {
	ID  string
	Bar string
}
