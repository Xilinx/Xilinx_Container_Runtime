/*
 * Copyright (C) 2019-2021, Xilinx Inc - All rights reserved
 * Xilinx Container Runtime
 *
 * Licensed under the Apache License, Version 2.0 (the "License"). You may
 * not use this file except in compliance with the License. A copy of the
 * License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package oci

// MockExecRuntime wraps a SyscallExecRuntime, intercepting the exec call for testing
type MockExecRuntime struct {
	SyscallExecRuntime
	execMock
}

// WithMockExec wraps a specified SyscallExecRuntime with a mocked exec function for testing
func WithMockExec(e SyscallExecRuntime, execResult error) *MockExecRuntime {
	m := MockExecRuntime{
		SyscallExecRuntime: e,
		execMock:           execMock{result: execResult},
	}
	// overrdie the exec function to the mocked exec function.
	m.SyscallExecRuntime.exec = m.execMock.exec
	return &m
}

type execMock struct {
	argv0  string
	argv   []string
	envv   []string
	result error
}

func (m *execMock) exec(argv0 string, argv []string, envv []string) error {
	m.argv0 = argv0
	m.argv = argv
	m.envv = envv

	return m.result
}
