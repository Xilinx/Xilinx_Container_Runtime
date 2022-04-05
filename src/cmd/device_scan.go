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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	SysfsDevices      = "/sys/bus/pci/devices"
	MgmtPrefix        = "/dev/xclmgmt"
	UserPrefix        = "/dev/dri"
	QdmaPrefix        = "/dev/xfpga"
	QDMASTR           = "dma.qdma.u"
	UserPFKeyword     = "drm"
	DRMSTR            = "renderD"
	ROMSTR            = "rom"
	SNSTR             = "xmc.u."
	DSAverFile        = "VBNV"
	DSAtsFile         = "timestamp"
	InstanceFile      = "instance"
	MgmtFile          = "mgmt_pf"
	UserFile          = "user_pf"
	VendorFile        = "vendor"
	DeviceFile        = "device"
	SNFile            = "serial_num"
	XilinxVendorID    = "0x10ee"
	AdvantechVendorID = "0x13fe"
	AWSVendorID       = "0x1d0f"
	AristaVendorID    = "0x3475"
)

// Usually, there is management PF and uer PF within one Xilinx device
type xilinxPair struct {
	Mgmt string // Management function node
	User string // User function node
	Qdma string
}

type xilinxDevice struct {
	index     string      // integer numbered
	shellVer  string      // Xilinx shell version
	timestamp string      // DSA timestamp
	DBDF      string      // this is for user pf
	deviceID  string      // devid of the user pf
	SN        string      // serial number
	Pair      *xilinxPair // pair with UserPF and MgmtPF node
}

// For some Xilinx card, like U30, there is multiple devices
type xilinxCard struct {
	index   int // integer numbered
	devices []xilinxDevice
}

func getInstance(DBDF string) (string, error) {
	strArray := strings.Split(DBDF, ":")
	domain, err := strconv.ParseUint(strArray[0], 16, 16)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[0])
	}
	bus, err := strconv.ParseUint(strArray[1], 16, 8)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[1])
	}
	strArray = strings.Split(strArray[2], ".")
	dev, err := strconv.ParseUint(strArray[0], 16, 8)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[0])
	}
	fc, err := strconv.ParseUint(strArray[1], 16, 8)
	if err != nil {
		return "", fmt.Errorf("strconv failed: %s\n", strArray[1])
	}
	ret := domain*65536 + bus*256 + dev*8 + fc
	return strconv.FormatUint(ret, 10), nil
}

func getFileNameFromPrefix(dir string, prefix string) (string, error) {
	userFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Can't read folder %s", dir)
	}
	for _, userFile := range userFiles {
		fname := userFile.Name()

		if !strings.HasPrefix(fname, prefix) {
			continue
		}
		return fname, nil
	}
	return "", nil
}

func getFileContent(file string) (string, error) {
	if buf, err := ioutil.ReadFile(file); err != nil {
		return "", fmt.Errorf("Can't read file %s", file)
	} else {
		ret := strings.Trim(string(buf), "\n")
		return ret, nil
	}
}

func fileExist(fname string) bool {
	if _, err := os.Stat(fname); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func isMgmtPf(pciID string) bool {
	fname := path.Join(SysfsDevices, pciID, MgmtFile)
	return fileExist(fname)
}

func isUserPf(pciID string) bool {
	fname := path.Join(SysfsDevices, pciID, UserFile)
	return fileExist(fname)
}

// Return a list of all Xilinx devices on host.
func getAllXilinxDevices() ([]xilinxDevice, error) {
	var devices []xilinxDevice
	pairMap := make(map[string]*xilinxPair)
	pciFiles, err := ioutil.ReadDir(SysfsDevices)
	if err != nil {
		return nil, fmt.Errorf("Can't read folder %s", SysfsDevices)
	}

	for _, pciFile := range pciFiles {
		pciID := pciFile.Name()

		fname := path.Join(SysfsDevices, pciID, VendorFile)
		vendorID, err := getFileContent(fname)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(vendorID, XilinxVendorID) != true &&
			strings.EqualFold(vendorID, AristaVendorID) != true &&
			strings.EqualFold(vendorID, AWSVendorID) != true &&
			strings.EqualFold(vendorID, AdvantechVendorID) != true {
			continue
		}

		DBD := pciID[:len(pciID)-2]
		if _, ok := pairMap[DBD]; !ok {
			pairMap[DBD] = &xilinxPair{
				Mgmt: "",
				User: "",
				Qdma: "",
			}
		}

		if isUserPf(pciID) { //user pf
			userDBDF := pciID
			romFolder, err := getFileNameFromPrefix(path.Join(SysfsDevices, pciID), ROMSTR)
			count := 0
			if err != nil {
				return nil, err
			}
			for romFolder == "" {
				if count >= 36 {
					break
				}
				time.Sleep(10 * time.Second)
				romFolder, err = getFileNameFromPrefix(path.Join(SysfsDevices, pciID), ROMSTR)
				if romFolder != "" {
					time.Sleep(20 * time.Second)
					break
				}
				fmt.Println(count, pciID, romFolder, err)
				count += 1
			}
			SNFolder, err := getFileNameFromPrefix(path.Join(SysfsDevices, pciID), SNSTR)
			if err != nil {
				return nil, err
			}
			// get dsa version
			fname = path.Join(SysfsDevices, pciID, romFolder, DSAverFile)
			content, err := getFileContent(fname)
			if err != nil {
				return nil, err
			}
			dsaVer := content
			// get dsa timestamp
			fname = path.Join(SysfsDevices, pciID, romFolder, DSAtsFile)
			content, err = getFileContent(fname)
			if err != nil {
				return nil, err
			}
			dsaTs := content
			// get device id
			fname = path.Join(SysfsDevices, pciID, DeviceFile)
			content, err = getFileContent(fname)
			if err != nil {
				return nil, err
			}
			devid := content
			// get Serial Number
			fname = path.Join(SysfsDevices, pciID, SNFolder, SNFile)
			content, err = getFileContent(fname)
			SN := ""
			if err == nil {
				SN = content
			}
			// get user PF node
			userpf, err := getFileNameFromPrefix(path.Join(SysfsDevices, pciID, UserPFKeyword), DRMSTR)
			if err != nil {
				return nil, err
			}
			userNode := path.Join(UserPrefix, userpf)
			pairMap[DBD].User = userNode

			//get qdma device node if it exists
			instance, err := getInstance(userDBDF)
			if err != nil {
				return nil, err
			}

			qdmaFolder, err := getFileNameFromPrefix(path.Join(SysfsDevices, pciID), QDMASTR)
			if err != nil {
				return nil, err
			}

			if qdmaFolder != "" {
				pairMap[DBD].Qdma = path.Join(QdmaPrefix, QDMASTR+instance+".0")
			}

			devices = append(devices, xilinxDevice{
				index:     strconv.Itoa(len(devices)),
				shellVer:  dsaVer,
				timestamp: dsaTs,
				DBDF:      userDBDF,
				deviceID:  devid,
				SN:        SN,
				Pair:      pairMap[DBD],
			})
		} else if isMgmtPf(pciID) { //mgmt pf
			// get mgmt instance
			fname = path.Join(SysfsDevices, pciID, InstanceFile)
			content, err := getFileContent(fname)
			if err != nil {
				return nil, err
			}
			pairMap[DBD].Mgmt = MgmtPrefix + content
		}
	}
	return devices, nil
}

// Return a list of device based on device environment variable, like '0,1', 'all', etc.
func getXilinxDevicesByDeviceEnv(visibleDevicesEnv string) ([]xilinxDevice, error) {
	allDevices, err := getAllXilinxDevices()
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(visibleDevicesEnv, "ALL") {
		return allDevices, nil
	}

	visibleDevices := []xilinxDevice{}
	parts := strings.Split(visibleDevicesEnv, ",")
	for _, part := range parts {
		for _, device := range allDevices {
			if part == device.index || part == device.deviceID || part == device.SN {
				visibleDevices = append(visibleDevices, device)
			}
		}
	}

	return visibleDevices, nil
}

// Return a list of all Xilinx cards on host
func getAllXilinxCards() ([]xilinxCard, error) {
	allDevices, err := getAllXilinxDevices()
	if err != nil {
		return nil, err
	}

	cards := []xilinxCard{}
	m := make(map[string]int)
	for _, device := range allDevices {
		if strings.TrimSpace(device.SN) == "" {
			// No serial number found, treated as a single device card
			index := len(cards)
			cards = append(cards, xilinxCard{
				index: index,
				devices: []xilinxDevice{
					device,
				},
			})
		} else {
			index, exisited := m[device.SN]
			if !exisited {
				index = len(cards)
				cards = append(cards, xilinxCard{
					index:   index,
					devices: []xilinxDevice{},
				})
				m[device.SN] = index
			}
			cards[index].devices = append(cards[index].devices, device)
		}
	}

	return cards, nil
}

// Return a list of devices baed on card number, like '0', '1', etc.
func getXilinxDevicesByCardNum(num int) ([]xilinxDevice, error) {
	allcards, err := getAllXilinxCards()
	if err != nil {
		return nil, err
	}

	if num >= len(allcards) {
		return nil, fmt.Errorf("card number %d not existed", num)
	}
	return allcards[num].devices, nil
}

// Return a list of devices based on card environment variable, like '0,1' or 'all', etc.
func getXilinxDevicesByCardEnv(visibleCardsEnv string) ([]xilinxDevice, error) {
	if strings.EqualFold(visibleCardsEnv, "all") {
		return getAllXilinxDevices()
	}

	var visibleXilinxDevices []xilinxDevice
	cardNums := strings.Split(visibleCardsEnv, ",")
	for _, cardNum := range cardNums {
		num, err := strconv.Atoi(cardNum)
		if err != nil {
			return nil, fmt.Errorf("only int numbers allowed for env %s", envXLNXVisibleCards)
		}
		card, err := getXilinxDevicesByCardNum(num)
		if err != nil {
			return nil, fmt.Errorf("error getting xilinx card: %v", err)
		}
		visibleXilinxDevices = append(visibleXilinxDevices, card...)
	}

	return visibleXilinxDevices, nil
}

// Return device major and minor numbers based on device path
func getDeviceMajorMinor(devPath string) (int64, int64, error) {
	stat := syscall.Stat_t{}
	err := syscall.Stat(devPath, &stat)
	if err != nil {
		return 0, 0, err
	}
	major := int64(stat.Rdev >> 8)
	minor := int64(stat.Rdev & (256 - 1))
	return major, minor, nil
}
