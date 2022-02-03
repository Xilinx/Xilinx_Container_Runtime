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
	"fmt"
	"os"
	"testing"

	"github.com/opencontainers/runtime-spec/specs-go"
	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestArgsGetConfigFilePath(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	testCases := []struct {
		bundleDir   string
		ociSpecPath string
	}{
		{
			ociSpecPath: fmt.Sprintf("%v/config.json", wd),
		},
		{
			bundleDir:   "/foo/bar",
			ociSpecPath: "/foo/bar/config.json",
		},
		{
			bundleDir:   "/foo/bar/",
			ociSpecPath: "/foo/bar/config.json",
		},
	}

	for i, tc := range testCases {
		cp, err := getOCISpecFilePath(tc.bundleDir)

		require.NoErrorf(t, err, "%d: %v", i, tc)
		require.Equalf(t, tc.ociSpecPath, cp, "%d: %v", i, tc)
	}
}

func TestAddXilinxDevices(t *testing.T) {
	logger, logHook := testlog.NewNullLogger()
	shim := xilinxContainerRuntime{
		logger: logger,
	}

	testCases := []struct {
		spec      *specs.Spec
		shouldAdd bool
	}{
		{
			spec:      &specs.Spec{},
			shouldAdd: false,
		},
		{
			spec: &specs.Spec{
				Process: &specs.Process{
					Env: []string{},
				},
			},
			shouldAdd: false,
		},
		{
			spec: &specs.Spec{
				Process: &specs.Process{
					Env: []string{
						"XILINX_VISIBLE_DEVICES=all",
					},
				},
				Linux: &specs.Linux{
					Resources: &specs.LinuxResources{
						Devices: []specs.LinuxDeviceCgroup{},
					},
				},
			},
			shouldAdd: true,
		},
		{
			spec: &specs.Spec{
				Process: &specs.Process{
					Env: []string{
						"XILINX_VISIBLE_DEVICES=0",
					},
				},
				Linux: &specs.Linux{
					Resources: &specs.LinuxResources{
						Devices: []specs.LinuxDeviceCgroup{},
					},
				},
			},
			shouldAdd: true,
		},
		{
			spec: &specs.Spec{
				Process: &specs.Process{
					Env: []string{
						"XILINX_VISIBLE_CARDS=0",
					},
				},
				Linux: &specs.Linux{
					Resources: &specs.LinuxResources{
						Devices: []specs.LinuxDeviceCgroup{},
					},
				},
			},
			shouldAdd: true,
		},
	}

	for i, tc := range testCases {
		logHook.Reset()

		var numMounts int
		if tc.spec.Mounts != nil {
			numMounts = len(tc.spec.Mounts)
		}

		err := shim.addXilinxDevices(tc.spec)
		require.NoErrorf(t, err, "%d: %v", i, tc)
		if tc.shouldAdd {
			require.Greater(t, len(tc.spec.Mounts), numMounts, "%d: %v", i, tc)
		} else {
			if tc.spec.Mounts != nil {
				require.Equal(t, numMounts, len(tc.spec.Mounts), "%d: %v", i, tc)
			}
		}

	}
}

func TestModificationRequired(t *testing.T) {
	logger, logHook := testlog.NewNullLogger()

	testCases := []struct {
		shim         xilinxContainerRuntime
		shouldModify bool
		args         []string
	}{
		{
			shim:         xilinxContainerRuntime{},
			shouldModify: false,
		},
		{
			shim:         xilinxContainerRuntime{},
			args:         []string{"create"},
			shouldModify: true,
		},
		{
			shim:         xilinxContainerRuntime{},
			args:         []string{"run"},
			shouldModify: true,
		},
		{
			shim:         xilinxContainerRuntime{},
			args:         []string{"--bundle=create"},
			shouldModify: false,
		},
		{
			shim:         xilinxContainerRuntime{},
			args:         []string{"--bundle", "create"},
			shouldModify: false,
		},
		{
			shim:         xilinxContainerRuntime{},
			args:         []string{"--bundle=run"},
			shouldModify: false,
		},
		{
			shim:         xilinxContainerRuntime{},
			args:         []string{"--bundle", "run"},
			shouldModify: false,
		},
	}

	for i, tc := range testCases {
		tc.shim.logger = logger
		logHook.Reset()

		require.Equal(t, tc.shouldModify, tc.shim.modificationRequired(tc.args), "%d: %v", i, tc)
	}
}
