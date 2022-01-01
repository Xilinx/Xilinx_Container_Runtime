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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Xilinx/xilinx-container-runtime/pkg/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
)

const (
	envXLNXVisibleDevices = "XILINX_VISIBLE_DEVICES"
	envXLNXVisibleCards   = "XILINX_VISIBLE_CARDS"
	envXLNXDeviceExlusive = "XILINX_DEVICE_EXLUSIVE"
)

// xilinxContainerRuntime wraps specified runtime, conditionally modifying OCI spec before invoking the spcified runtime
type xilinxContainerRuntime struct {
	logger  *log.Logger
	cfg     *config
	runtime oci.Runtime
	ocispec oci.Spec
}

var _ oci.Runtime = (*xilinxContainerRuntime)(nil)

// Constructor for xilinx container runtime
func newXilinxContainerRuntimeWithLogger(logger *log.Logger, cfg *config, runtime oci.Runtime, ociSpec oci.Spec) (oci.Runtime, error) {
	r := xilinxContainerRuntime{
		logger:  logger,
		cfg:     cfg,
		runtime: runtime,
		ocispec: ociSpec,
	}

	return &r, nil
}

func (r xilinxContainerRuntime) modificationRequired(args []string) bool {
	var previousWasBundle bool
	for _, a := range args {
		// We check for '--bundle|-b create|run|modify' explicitly to ensure
		// that we don't inadvertently trigger a modification if the bundle
		// directory is specified as 'create', 'run' or 'modify'
		if !previousWasBundle && isBundleFlag(a) {
			previousWasBundle = true
			continue
		}

		if !previousWasBundle && (a == "create" || a == "run" || a == "modify") {
			r.logger.Infof("'create', 'run' or 'modify' command detected, modification required")
			return true
		}

		previousWasBundle = false
	}

	r.logger.Infof("No modification required")
	return false
}

func (r xilinxContainerRuntime) forwardingRequired(args []string) bool {
	var previousWasBundle bool
	for _, a := range args {
		// We check for '--bundle|-b modify' explicitly to ensure
		// that we don't inadvertently forward the commands if the
		// bundle directory is specified as 'modify'
		if !previousWasBundle && isBundleFlag(a) {
			previousWasBundle = true
			continue
		}

		if !previousWasBundle && (a == "modify") {
			r.logger.Infof("'modify' command detected, no forwarding required")
			return false
		}

		previousWasBundle = false
	}
	return true
}

func (r xilinxContainerRuntime) deviceExlusionEnabled(spec *specs.Spec) bool {
	deviceExclusiveEnv := ""
	if spec.Process != nil && spec.Process.Env != nil {
		// Check environment variable from OCI Spec file
		for _, str := range spec.Process.Env {
			parts := strings.SplitN(str, "=", 2)

			if len(parts) != 2 {
				continue
			}

			if parts[0] == envXLNXDeviceExlusive {
				deviceExclusiveEnv = parts[1]
			}
		}
	}

	r.logger.Printf("Environment Variable %s=%s", envXLNXDeviceExlusive, deviceExclusiveEnv)
	if deviceExclusiveEnv != "" {
		exlusive, err := strconv.ParseBool(deviceExclusiveEnv)
		if err == nil {
			return exlusive
		} else {
			r.logger.Printf("error getting device exclusive enable status %v", err)
			return r.cfg.exclusive
		}
	}

	return r.cfg.exclusive
}

func (r xilinxContainerRuntime) getDeviceExlusions() (map[string]bool, error) {
	file, err := os.Open(r.cfg.exclusionFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening device exlusion file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var exclusions map[string]bool
	err = decoder.Decode(&exclusions)
	if err != nil {
		return nil, fmt.Errorf("error reading device exlusions from file: %v", err)
	}

	return exclusions, nil
}

func (r xilinxContainerRuntime) setDeviceExlusions(exclusions map[string]bool) error {
	file, err := os.Create(r.cfg.exclusionFilePath)
	if err != nil {
		return fmt.Errorf("error opening device exlusion file: %v", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(exclusions)
	if err != nil {
		return fmt.Errorf("error writing device exlusions to file: %v", err)
	}
	return nil
}

func (r xilinxContainerRuntime) modifyOCISpec() error {
	err := r.ocispec.Load()
	if err != nil {
		return fmt.Errorf("error loading OCI specification for modification: %v", err)
	}

	err = r.ocispec.Modify(r.addXilinxDevices)

	if err != nil {
		return fmt.Errorf("error adding Xilinx devices in OCI Spec: %v", err)
	}

	err = r.ocispec.Flush()
	if err != nil {
		return fmt.Errorf("error writing modified OCI specification: %v", err)
	}
	return nil
}

func (r xilinxContainerRuntime) addXilinxDevices(spec *specs.Spec) error {
	visibleDevicesEnv := ""
	visibleCardsEnv := ""
	var visibleXilinxDevices []xilinxDevice

	if spec.Process != nil && spec.Process.Env != nil {
		// Check environment variable from OCI Spec file
		for _, str := range spec.Process.Env {
			parts := strings.SplitN(str, "=", 2)

			if len(parts) != 2 {
				continue
			}

			if parts[0] == envXLNXVisibleDevices {
				visibleDevicesEnv = parts[1]
			} else if parts[0] == envXLNXVisibleCards {
				visibleCardsEnv = parts[1]
			}
		}
	}

	// Check OS environment variables
	if visibleDevicesEnv == "" {
		visibleDevicesEnv = os.Getenv(envXLNXVisibleDevices)
	}
	if visibleCardsEnv == "" {
		visibleCardsEnv = os.Getenv(envXLNXVisibleCards)
	}

	if visibleDevicesEnv == "" && visibleCardsEnv == "" {
		// Do nothing since no envs specified
		logger.Infof("Environment variable %s and %s is not specified", envXLNXVisibleDevices, envXLNXVisibleCards)
		return nil
	}

	if visibleDevicesEnv != "" {
		devices, err := getXilinxDevicesByDeviceEnv(visibleDevicesEnv)
		if err != nil {
			return fmt.Errorf("error getting xilinx devices: %v", err)
		}
		visibleXilinxDevices = append(visibleXilinxDevices, devices...)
	} else {
		devices, err := getXilinxDevicesByCardEnv(visibleCardsEnv)
		if err != nil {
			return fmt.Errorf("error getting xilinx devices: %v", err)
		}
		visibleXilinxDevices = append(visibleXilinxDevices, devices...)
	}

	r.logger.Printf("There is %d device(s) to be mounted", len(visibleXilinxDevices))

	// check whether it is in device exclusive mode
	if r.deviceExlusionEnabled(spec) {
		// get current device exclusion status from file
		deviceExlusions, err := r.getDeviceExlusions()
		if err != nil {
			return err
		}

		// check whether requested device(s) are excluded by another container
		// set requested device(s) as excluded
		for _, device := range visibleXilinxDevices {
			if deviceExlusions[device.DBDF] {
				r.logger.Printf("Device %s is being used exlusively by another container", device.DBDF)
				return fmt.Errorf("Device %s is being used exlusively by another container", device.DBDF)
			} else {
				r.logger.Printf("Device %s will be used exlusively to this container", device.DBDF)
				deviceExlusions[device.DBDF] = true
			}
		}

		// flush the device exclusion status into file
		err = r.setDeviceExlusions(deviceExlusions)
		if err != nil {
			return err
		}

	}

	for _, device := range visibleXilinxDevices {
		// Check whether the device is in the mount config already
		userMounted, mgmtMounted := false, false
		for _, mount := range spec.Mounts {
			if userMounted && mgmtMounted {
				break
			}
			if device.Nodes.User == mount.Source {
				userMounted = true
			}
			if device.Nodes.Mgmt == mount.Source {
				mgmtMounted = true
			}
		}

		if !userMounted && len(strings.TrimSpace(device.Nodes.User)) != 0 {
			// Mount user node
			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: device.Nodes.User,
				Type:        "none",
				Source:      device.Nodes.User,
				Options:     []string{"nosuid", "noexec", "bind"},
			})
		}

		if !mgmtMounted && len(strings.TrimSpace(device.Nodes.Mgmt)) != 0 {
			// Mount mgmt node
			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: device.Nodes.Mgmt,
				Type:        "none",
				Source:      device.Nodes.Mgmt,
				Options:     []string{"nosuid", "noexec", "bind"},
			})
		}

		// Check whether user device is mapped in Linux Devices config
		deviceMapped := false
		major, minor, err := getDeviceMajorMinor(device.Nodes.User)
		for _, device := range spec.Linux.Resources.Devices {
			if deviceMapped {
				break
			}
			if device.Major == nil || device.Minor == nil {
				continue
			}
			if *(device.Major) == major && *(device.Minor) == minor {
				deviceMapped = true
				break
			}
		}

		if !deviceMapped {
			if err != nil {
				return fmt.Errorf("error getting device major and minor numbers: %v", err)
			}
			spec.Linux.Resources.Devices = append(spec.Linux.Resources.Devices, specs.LinuxDeviceCgroup{
				Allow:  true,
				Type:   "c",
				Major:  &major,
				Minor:  &minor,
				Access: "rw",
			})
		}
	}
	return nil
}

func (r xilinxContainerRuntime) Exec(args []string) error {
	// Check wehther it is rquired to modify the OCI spec
	if r.modificationRequired(args) {
		err := r.modifyOCISpec()
		if err != nil {
			return fmt.Errorf("Fail to modify OCI spec: %v", err)
		}
	}

	if r.forwardingRequired(args) {
		r.logger.Println("Forwarding command to underlying runtime")
		return r.runtime.Exec(args)
	} else {
		r.logger.Println("No forwarding required")
		return nil
	}
}
