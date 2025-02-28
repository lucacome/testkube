// Copyright 2024 Testkube.
//
// Licensed as a Testkube Pro file under the Testkube Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//	https://github.com/kubeshop/testkube/blob/main/licenses/TCL.txt

package data

import "github.com/kubeshop/testkube/pkg/tcl/expressionstcl"

var baseTestWorkflowMachine = expressionstcl.CombinedMachines(EnvMachine, StateMachine, FileMachine)

func GetBaseTestWorkflowMachine() expressionstcl.Machine {
	LoadState()
	return baseTestWorkflowMachine
}
