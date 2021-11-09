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

// Runtime is an interface for a runtime shim. The Exec method accepts a list
// of command line arguments, and returns an error / nil.
type Runtime interface {
	Exec([]string) error
}
