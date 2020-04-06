package main

import "fmt"

type ErrDependencyNotFound struct {
	Dependency string
}

func (e *ErrDependencyNotFound) Error() string {
	return fmt.Sprintf("failed to find %s dependency", e.Dependency)
}

type ErrPluginExisted struct {
	Name string
}

func (e *ErrPluginExisted) Error() string {
	return fmt.Sprintf("%s plugin has been registered", e.Name)
}
