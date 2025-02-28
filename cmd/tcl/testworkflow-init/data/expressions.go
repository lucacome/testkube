// Copyright 2024 Testkube.
//
// Licensed as a Testkube Pro file under the Testkube Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//	https://github.com/kubeshop/testkube/blob/main/licenses/TCL.txt

package data

import (
	"os"
	"strings"

	"github.com/kubeshop/testkube/pkg/tcl/expressionstcl"
	"github.com/kubeshop/testkube/pkg/tcl/expressionstcl/libs"
)

var aliases = map[string]string{
	"always": `true`,
	"never":  `false`,

	"error":   `failed`,
	"success": `passed`,

	"self.error":   `self.failed`,
	"self.success": `self.passed`,

	"passed": `!status`,
	"failed": `bool(status) && status != "skipped"`,

	"self.passed": `!self.status`,
	"self.failed": `bool(self.status) && self.status != "skipped"`,
}

var LocalMachine = expressionstcl.NewMachine().
	Register("status", expressionstcl.MustCompile("self.status"))

var RefMachine = expressionstcl.NewMachine().
	RegisterAccessor(func(name string) (interface{}, bool) {
		if name == "_ref" {
			return Step.Ref, true
		}
		return nil, false
	})

var AliasMachine = expressionstcl.NewMachine().
	RegisterAccessorExt(func(name string) (interface{}, bool, error) {
		alias, ok := aliases[name]
		if !ok {
			return nil, false, nil
		}
		expr, err := expressionstcl.Compile(alias)
		if err != nil {
			return expr, false, err
		}
		expr, err = expr.Resolve(RefMachine)
		return expr, true, err
	})

var StateMachine = expressionstcl.NewMachine().
	RegisterAccessor(func(name string) (interface{}, bool) {
		if name == "status" {
			return State.GetStatus(), true
		} else if name == "self.status" {
			return State.GetSelfStatus(), true
		}
		return nil, false
	}).
	RegisterAccessorExt(func(name string) (interface{}, bool, error) {
		if strings.HasPrefix(name, "output.") {
			return State.GetOutput(name[7:])
		}
		return nil, false, nil
	})

var EnvMachine = expressionstcl.NewMachine().
	RegisterAccessor(func(name string) (interface{}, bool) {
		if strings.HasPrefix(name, "env.") {
			return os.Getenv(name[4:]), true
		}
		return nil, false
	})

var RefSuccessMachine = expressionstcl.NewMachine().
	RegisterAccessor(func(ref string) (interface{}, bool) {
		s := State.GetStep(ref)
		return s.Status == StepStatusPassed || s.Status == StepStatusSkipped, s.HasStatus
	})

var RefStatusMachine = expressionstcl.NewMachine().
	RegisterAccessor(func(ref string) (interface{}, bool) {
		return string(State.GetStep(ref).Status), true
	})

var wd, _ = os.Getwd()
var FileMachine = libs.NewFsMachine(os.DirFS("/"), wd)

func Template(tpl string, m ...expressionstcl.Machine) (string, error) {
	m = append(m, AliasMachine, baseTestWorkflowMachine)
	return expressionstcl.EvalTemplate(tpl, m...)
}

func Expression(expr string, m ...expressionstcl.Machine) (expressionstcl.StaticValue, error) {
	m = append(m, AliasMachine, baseTestWorkflowMachine)
	return expressionstcl.EvalExpression(expr, m...)
}

func RefSuccessExpression(expr string) (expressionstcl.StaticValue, error) {
	return expressionstcl.EvalExpression(expr, RefSuccessMachine)
}

func RefStatusExpression(expr string) (expressionstcl.StaticValue, error) {
	return expressionstcl.EvalExpression(expr, RefStatusMachine)
}
