/*
 * Copyright (C) 2022, Xilinx Inc - All rights reserved
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
	"sync"
	"time"

	"github.com/Xilinx/xilinx-container-runtime/src/pkg/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
)

const (
	envXLNXVisibleDevices  = "XILINX_VISIBLE_DEVICES"
	envXLNXVisibleCards    = "XILINX_VISIBLE_CARDS"
	envXLNXDeviceExclusive = "XILINX_DEVICE_EXCLUSIVE"
)

// xilinxContainerRuntime wraps specified runtime, conditionally modifying OCI spec before invoking the spcified runtime
type xilinxContainerRuntime struct {
	logger  *log.Logger
	cfg     *config
	runtime oci.Runtime
	ocispec oci.Spec
	mutex   *sync.Mutex
}

type xilinxDeviceExclusions struct {
	Notice  string         `json:"notice"`
	Devices map[string]int `json:"devices"`
}

var _ oci.Runtime = (*xilinxContainerRuntime)(nil)

// Constructor for xilinx container runtime
func newXilinxContainerRuntimeWithLogger(logger *log.Logger, cfg *config, runtime oci.Runtime, ociSpec oci.Spec) (oci.Runtime, error) {
	r := xilinxContainerRuntime{
		logger:  logger,
		cfg:     cfg,
		runtime: runtime,
		ocispec: ociSpec,
		mutex:   new(sync.Mutex),
	}

	return &r, nil
}

// check if modification is required for current command
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

// check if it is 'create' command
func (r xilinxContainerRuntime) addDeviceExclusionsRequired(args []string) bool {
	var previousWasBundle bool
	for _, a := range args {
		// We check for '--bundle|-b create' explicitly to ensure
		// that we don't inadvertently trigger a modification if the bundle
		// directory is specified as 'create'
		if !previousWasBundle && isBundleFlag(a) {
			previousWasBundle = true
			continue
		}

		if !previousWasBundle && (a == "create") {
			r.logger.Infof("'create' command detected, device status update required")
			return true
		}

		previousWasBundle = false
	}

	r.logger.Infof("No device status update required")
	return false
}

// check if it is 'delete' command
func (r xilinxContainerRuntime) deleteDeviceExclusionsRequired(args []string) bool {
	var previousWasBundle bool
	for _, a := range args {
		// We check for '--bundle|-b delete' explicitly to ensure
		// that we don't inadvertently trigger a modification if the bundle
		// directory is specified as 'delete'
		if !previousWasBundle && isBundleFlag(a) {
			previousWasBundle = true
			continue
		}

		if !previousWasBundle && (a == "delete") {
			r.logger.Infof("'delete' command detected, device status update required")
			return true
		}

		previousWasBundle = false
	}

	r.logger.Infof("Not 'delete' commnad, not updating device status")
	return false
}

// check if it is required to forward command to inner runtime (runC)
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

// get visible devices list based on environment variables
func (r xilinxContainerRuntime) getVisibleDevices(spec *specs.Spec) ([]xilinxDevice, error) {
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
		return nil, nil
	}

	if visibleDevicesEnv != "" {
		devices, err := getXilinxDevicesByDeviceEnv(visibleDevicesEnv)
		if err != nil {
			return nil, fmt.Errorf("error getting xilinx devices: %v", err)
		}
		visibleXilinxDevices = append(visibleXilinxDevices, devices...)
	} else {
		devices, err := getXilinxDevicesByCardEnv(visibleCardsEnv)
		if err != nil {
			return nil, fmt.Errorf("error getting xilinx devices: %v", err)
		}
		visibleXilinxDevices = append(visibleXilinxDevices, devices...)
	}

	return visibleXilinxDevices, nil
}

// check if device exclusive is enabled for this container
func (r xilinxContainerRuntime) deviceExclusiveEnabled(spec *specs.Spec) bool {
	deviceExclusiveEnv := ""
	if spec.Process != nil && spec.Process.Env != nil {
		// Check environment variable from OCI Spec file
		for _, str := range spec.Process.Env {
			parts := strings.SplitN(str, "=", 2)

			if len(parts) != 2 {
				continue
			}

			if parts[0] == envXLNXDeviceExclusive {
				deviceExclusiveEnv = parts[1]
			}
		}
	}

	if deviceExclusiveEnv != "" {
		exclusive, err := strconv.ParseBool(strings.ToLower(deviceExclusiveEnv))
		if err == nil {
			return exclusive
		} else {
			r.logger.Printf("error getting device exclusive enable status %v", err)
			return r.cfg.deviceExclusive
		}
	}

	return r.cfg.deviceExclusive
}

// get current device exclusion stats from file
// -1 meaning device is being used by a container exclusively
// non-negtive integers meaning the number of containers are sharing the device
func (r xilinxContainerRuntime) getDeviceExclusions() (map[string]int, error) {
	exclusions := xilinxDeviceExclusions{
		Notice:  "",
		Devices: make(map[string]int),
	}
	if _, err := os.Stat(r.cfg.exclusionFilePath); os.IsNotExist(err) {
		return exclusions.Devices, nil
	}

	file, err := os.Open(r.cfg.exclusionFilePath)

	if err != nil {
		return nil, fmt.Errorf("error opening device exclusion file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&exclusions)
	if err != nil {
		return nil, fmt.Errorf("error reading device exclusions from file: %v", err)
	}

	return exclusions.Devices, nil
}

// save device exclusion stats into file
func (r xilinxContainerRuntime) setDeviceExclusions(devices map[string]int) error {
	file, err := os.Create(r.cfg.exclusionFilePath)
	if err != nil {
		return fmt.Errorf("error opening device exclusion file: %v", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	currentTime := time.Now().Format("2006-01-02 3:4:5 pm")
	exclusions := xilinxDeviceExclusions{
		Notice: fmt.Sprintf(
			"This file stores the status of xilinx devices usage, which was saved on %s. '-1' means the device is being used exclusively. 0 or positive integer is the number of containers currently using respective device.",
			currentTime),
		Devices: devices,
	}

	err = encoder.Encode(exclusions)
	if err != nil {
		return fmt.Errorf("error writing device exclusions to file: %v", err)
	}
	return nil
}

// modify OCI spec to add xilinx devices
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

// check and add device exclusions while creating the container
func (r xilinxContainerRuntime) addDeviceExclusions(spec *specs.Spec) error {
	visibleXilinxDevices, err := r.getVisibleDevices(spec)
	if err != nil {
		return err
	} else if visibleXilinxDevices == nil || len(visibleXilinxDevices) == 0 {
		return nil
	} else {
		r.logger.Infof("Updating device exclusions status for %d device(s)", len(visibleXilinxDevices))
	}

	// get current device exclusion status from file
	deviceExclusions, err := r.getDeviceExclusions()
	if err != nil {
		return err
	}

	// check whether it is in device exclusive mode
	if r.deviceExclusiveEnabled(spec) {
		// In device exclucsive mode, assign device to this
		// container only if the current device exclusion value is 0
		for _, device := range visibleXilinxDevices {
			if deviceExclusions[device.DBDF] != 0 {
				r.logger.Printf("Device %s is being used by another container", device.DBDF)
				return fmt.Errorf("Device %s is being used by another container", device.DBDF)
			} else {
				r.logger.Printf("Device %s will be used exclusively by this container", device.DBDF)
				deviceExclusions[device.DBDF] = -1
			}
		}
	} else {
		// Not in device exclusive mode, assign device to this
		// container if current device exclusion value is not -1
		for _, device := range visibleXilinxDevices {
			if deviceExclusions[device.DBDF] == -1 {
				r.logger.Printf("Device %s is being used exclusively by another container", device.DBDF)
				return fmt.Errorf("Device %s is being used exclusively by another container", device.DBDF)
			} else {
				r.logger.Printf("Device %s will be used by this container", device.DBDF)
				deviceExclusions[device.DBDF] = deviceExclusions[device.DBDF] + 1
			}
		}
	}

	// flush the device exclusion status into file
	r.logger.Printf("Trying to updated device exclusion status to file.")
	err = r.setDeviceExclusions(deviceExclusions)
	if err != nil {
		return err
	}
	return nil
}

// delete device exclusions while deleting the container
func (r xilinxContainerRuntime) deleteDeviceExclusions(spec *specs.Spec) error {

	visibleXilinxDevices, err := r.getVisibleDevices(spec)
	if err != nil {
		return err
	} else if visibleXilinxDevices == nil || len(visibleXilinxDevices) == 0 {
		r.logger.Infof("There is no device used in this container")
		return nil
	} else {
		r.logger.Infof("There is %d device(s) used in this container", len(visibleXilinxDevices))
	}

	// get current device exclusion status from file
	deviceExclusions, err := r.getDeviceExclusions()
	if err != nil {
		return err
	}

	isExclusiveMode := r.deviceExclusiveEnabled(spec)
	for _, device := range visibleXilinxDevices {
		if isExclusiveMode {
			// set the exclusion value 0 from -1
			deviceExclusions[device.DBDF] = 0
		} else {
			// do a decrement from current value
			deviceExclusions[device.DBDF] = deviceExclusions[device.DBDF] - 1
		}
	}

	// flush the device exclusion status into file
	r.logger.Printf("Trying to updated device exclusion status to file.")
	err = r.setDeviceExclusions(deviceExclusions)
	if err != nil {
		return err
	}
	return nil
}

// add xilinx devices in OCI Spec
func (r xilinxContainerRuntime) addXilinxDevices(spec *specs.Spec) error {
	visibleXilinxDevices, err := r.getVisibleDevices(spec)
	if err != nil {
		return err
	} else if visibleXilinxDevices == nil || len(visibleXilinxDevices) == 0 {
		r.logger.Infof("There is no device to be mounted")
		return nil
	} else {
		r.logger.Infof("There is %d device(s) to be mounted", len(visibleXilinxDevices))
	}

	for _, device := range visibleXilinxDevices {
		// Check whether the device is in the mount config already
		userMounted, mgmtMounted := false, false
		for _, mount := range spec.Mounts {
			if userMounted && mgmtMounted {
				break
			}
			if device.Pair.User == mount.Source {
				userMounted = true
			}
			if device.Pair.Mgmt == mount.Source {
				mgmtMounted = true
			}
		}

		if !userMounted && len(strings.TrimSpace(device.Pair.User)) != 0 {
			// Mount user node
			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: device.Pair.User,
				Type:        "none",
				Source:      device.Pair.User,
				Options:     []string{"nosuid", "noexec", "bind"},
			})
		}

		if !mgmtMounted && len(strings.TrimSpace(device.Pair.Mgmt)) != 0 {
			// Mount mgmt node
			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: device.Pair.Mgmt,
				Type:        "none",
				Source:      device.Pair.Mgmt,
				Options:     []string{"nosuid", "noexec", "bind"},
			})
		}

		// Check whether user device is mapped in Linux Devices config
		deviceMapped := false
		major, minor, err := getDeviceMajorMinor(device.Pair.User)
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

// method to be called from main method
func (r xilinxContainerRuntime) Exec(args []string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update device exclusion status if required
	if r.addDeviceExclusionsRequired(args) {
		err := r.ocispec.Load()
		if err != nil {
			return fmt.Errorf("error loading OCI specification for modification: %v", err)
		}
		err = r.ocispec.Modify(r.addDeviceExclusions)
		if err != nil {
			return fmt.Errorf("Fail to update device exclusion status: %v. Please refer to file %s for details",
				err, r.cfg.exclusionFilePath)
		}
	}

	// Add xilinx devices in OCI Spec if required
	if r.modificationRequired(args) {
		err := r.modifyOCISpec()
		if err != nil {
			return fmt.Errorf("Fail to modify OCI spec: %v", err)
		}
	}

	// delete device exclusion status if required
	if r.deleteDeviceExclusionsRequired(args) {
		err := r.ocispec.Load()
		if err != nil {
			return fmt.Errorf("error loading OCI specification for modification: %v", err)
		}
		err = r.ocispec.Modify(r.deleteDeviceExclusions)
		if err != nil {
			return fmt.Errorf("Fail to delete device exclusion status: %v", err)
		}
	}

	// Forward command to inner runtime(runC) if required
	if r.forwardingRequired(args) {
		r.logger.Println("Forwarding command to underlying runtime")
		return r.runtime.Exec(args)
	} else {
		r.logger.Println("No forwarding required")
		return nil
	}
}
