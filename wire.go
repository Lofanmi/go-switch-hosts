//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func NewApplication() (*Application, func(), error) {
	panic(wire.Build(Sets))
}
