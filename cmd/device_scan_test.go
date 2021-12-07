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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllDevices(t *testing.T) {
	devices, err := getAllXilinxDevices()
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(devices), 1)
}

func TestGetAllCards(t *testing.T) {
	cards, err := getAllXilinxCards()
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(cards), 1)
}

func TestGetDevicesByDeviceEnv(t *testing.T) {
	testCases := []string{
		"0",
		"12809621T037",
		"0x5001",
	}

	for _, testCase := range testCases {
		visibleDevices, err := getXilinxDevicesByDeviceEnv(testCase)

		require.Nil(t, err)
		require.GreaterOrEqual(t, len(visibleDevices), 1)
	}
}

func TestGetDeviceMajorMinor(t *testing.T) {
	path := "/dev/dri/renderD128"
	major, minor, err := getDeviceMajorMinor(path)
	require.Nil(t, err)
	require.Equal(t, major, int64(226))
	require.Equal(t, minor, int64(128))
}
