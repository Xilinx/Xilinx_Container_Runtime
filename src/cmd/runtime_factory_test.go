/*
 * Copyright (C) 2021, Xilinx Inc - All rights reserved
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

package main

import (
	"strings"
	"testing"

	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	shim, err := newRuntime([]string{}, nil)
	require.NoError(t, err)
	require.NotNil(t, shim)
}

func TestFindRunc(t *testing.T) {
	testLogger, _ := testlog.NewNullLogger()
	logger.Logger = testLogger

	runcPath, err := findRunc()
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(runcPath, "runc"))
}
