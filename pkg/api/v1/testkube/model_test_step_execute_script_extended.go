/*
 * TestKube API
 *
 * TestKube provides a Kubernetes-native framework for test definition, execution and results
 *
 * API version: 1.0.0
 * Contact: testkube@kubeshop.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package testkube

import (
	"fmt"
)

func (s TestStepExecuteScript) FullName() string {
	return fmt.Sprintf("Script %s.%s", s.Namespace, s.Name)
}

func (s TestStepExecuteScript) Type() TestStepType {
	return EXECUTE_SCRIPT_TestStepType
}

func (s TestStepExecuteScript) StopOnFailure() bool {
	return s.StopTestOnFailure
}
