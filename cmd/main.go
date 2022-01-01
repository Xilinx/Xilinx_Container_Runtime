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
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/pborman/getopt"
	"github.com/pelletier/go-toml"
)

type config struct {
	debugFilePath     string
	exclusive         bool
	exclusionFilePath string
}

const (
	configOverride = "XCRT_CONFIG_HOME"
	configFilePath = "xilinx-container-runtime/config.toml"
)

var (
	configDir  = "/etc/"
	Version    = ""
	BuildTime  = ""
	GoVersion  = ""
	optVersion = getopt.BoolLong("version", 'v', "print the version")
	optHelp    = getopt.BoolLong("help", 'h', "show help")
)

var logger = NewLogger()

func usage() {
	fmt.Fprint(os.Stderr, "NAME:\n   xilinx-container-runtime\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nUSAGE:\n   xilinx-container-runtime [global options] command [command options] [arguments...]\n")
	fmt.Fprintf(os.Stderr, "\n")
	printVersion()
	fmt.Fprintf(os.Stderr, "\nCOMMANDS:\n")
	fmt.Fprintf(os.Stderr, "   checkpoint\tcheckpoint a running container\n")
	fmt.Fprintf(os.Stderr, "   create\tcreate a container\n")
	fmt.Fprintf(os.Stderr, "   delete\tdelete any resources held by the container often used with detached container\n")
	fmt.Fprintf(os.Stderr, "   events\tdisplay container events such as OOM notifications, cpu, memory, and IO usage statistics\n")
	fmt.Fprintf(os.Stderr, "   exec\t\texecute new process inside the container\n")
	fmt.Fprintf(os.Stderr, "   init\t\tinitialize the namespaces and launch the process\n")
	fmt.Fprintf(os.Stderr, "   kill\t\tkill sends the specified signal (default: SIGTERM) to the container's init process\n")
	fmt.Fprintf(os.Stderr, "   list\t\tlists containers started by runc with the given root\n")
	fmt.Fprintf(os.Stderr, "   lscard\t\tlists xilinx cards in the host\n")
	fmt.Fprintf(os.Stderr, "   lsdevice\t\tlists xilinx devices in the host\n")
	fmt.Fprintf(os.Stderr, "   pause\tpause suspends all processes inside the container\n")
	fmt.Fprintf(os.Stderr, "   ps\t\tps displays the processes running inside a container\n")
	fmt.Fprintf(os.Stderr, "   restore\trestore a container from a previous checkpoint\n")
	fmt.Fprintf(os.Stderr, "   resume\tresumes all processes that have been previously paused\n")
	fmt.Fprintf(os.Stderr, "   run\t\tcreate and run a container\n")
	fmt.Fprintf(os.Stderr, "   spec\t\tcreate a new specification file\n")
	fmt.Fprintf(os.Stderr, "   start\texecutes the user defined process in a created container\n")
	fmt.Fprintf(os.Stderr, "   state\toutput the state of a container\n")
	fmt.Fprintf(os.Stderr, "   update\tupdate container resource constraints\n")
	fmt.Fprintf(os.Stderr, "   help, h\tShows a list of commands or help for one command\n")
	fmt.Fprintf(os.Stderr, "\nGLOBAL OPTIONS:\n")
	fmt.Fprintf(os.Stderr, "   --debug\t\tenable debug output for logging\n")
	fmt.Fprintf(os.Stderr, "   --log value\t\tset the log file path where internal debug information is written\n")
	fmt.Fprintf(os.Stderr, "   --log-format value\tset the format used by logs ('text' (default), or 'json') (default: \"text\")\n")
	fmt.Fprintf(os.Stderr, "   --root value\t\troot directory for storage of container state (this should be located in tmpfs) (default: \"/run/runc\")\n")
	fmt.Fprintf(os.Stderr, "   --criu value\t\tpath to the criu binary used for checkpoint and restore (default: \"criu\")\n")
	fmt.Fprintf(os.Stderr, "   --systemd-cgroup\tenable debug output for logging\n")
	fmt.Fprintf(os.Stderr, "   --rootless value\tenable systemd cgroup support, expects cgroupsPath to be of form \"slice:prefix:name\" for e.g. \"system.slice:runc:434234\"\n")
	fmt.Fprintf(os.Stderr, "   --debug\t\tignore cgroup permission errors ('true', 'false', or 'auto') (default: \"auto\")\n")
	fmt.Fprintf(os.Stderr, "   --help, -h\t\tshow help\n")
	fmt.Fprintf(os.Stderr, "   --version, -v\tprint the version\n")

}

func printVersion() {
	fmt.Fprintf(os.Stderr, "VERSION:\n")
	fmt.Fprintf(os.Stderr, "   xilinx-container-rumtime version %s\n", Version)
	fmt.Fprintf(os.Stderr, "   go:\t%s\n", GoVersion)
	fmt.Fprintf(os.Stderr, "   build time:\t%s\n", BuildTime)
}

func printDevices() {
	xilinxDevices, err := getAllXilinxDevices()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	fmt.Fprintf(os.Stderr, "DeviceNum\tSerialNum\tDeviceBDF\tUserPF\t\t\tMgmtPF\t\t\tShellVersion\n")
	for _, xilinxDevice := range xilinxDevices {
		fmt.Fprintf(os.Stderr, "%-16s%-16s%-16s%-24s%-24s%s\n",
			xilinxDevice.index, xilinxDevice.SN, xilinxDevice.DBDF,
			xilinxDevice.Nodes.User, xilinxDevice.Nodes.Mgmt, xilinxDevice.shellVer)
	}
}

func printCards() {
	xilinxCards, err := getAllXilinxCards()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	fmt.Fprintf(os.Stderr, "CardNum\t\tSerialNum\tDeviceBDF\tUserPF\t\t\tMgmtPF\t\t\tShellVersion\n")
	for _, xilinxCard := range xilinxCards {
		for _, xilinxDevice := range xilinxCard.devices {
			fmt.Fprintf(os.Stderr, "%-16d%-16s%-16s%-24s%-24s%s\n",
				xilinxCard.index, xilinxCard.devices[0].SN, xilinxDevice.DBDF,
				xilinxDevice.Nodes.User, xilinxDevice.Nodes.Mgmt, xilinxCard.devices[0].shellVer)
		}
	}
}

func run(argv []string, cfg *config) (err error) {

	r, err := newRuntime(argv, cfg)
	if err != nil {
		return fmt.Errorf("error creating runtime: %v", err)
	}

	return r.Exec(argv)
}

// Read config values from a toml file or set via environment
func getConfig() (*config, error) {
	cfg := &config{}

	if XDGConfigDir := os.Getenv(configOverride); len(XDGConfigDir) != 0 {
		configDir = XDGConfigDir
	}

	configFilePath := path.Join(configDir, configFilePath)

	tomlContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	toml, err := toml.Load(string(tomlContent))
	if err != nil {
		return nil, err
	}

	cfg.debugFilePath = toml.GetDefault("xilinx-container-runtime.debug", "/dev/null").(string)
	cfg.exclusive = toml.GetDefault("device-exclusion.enabled", false).(bool)
	cfg.exclusionFilePath = toml.GetDefault("device-exclusion.exclusionFilePath", "/etc/xilinx-container-runtime/device-exclusion.json").(string)

	return cfg, nil
}

func main() {
	flag.Usage = usage
	cfg, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
	}

	err = logger.LogToFile(cfg.debugFilePath)
	if err != nil {
		logger.LogToFile("/dev/null")
	}

	logger.Printf("Running %v", os.Args)

	getopt.Getopt(nil)
	args := getopt.Args()

	if len(args) == 0 {
		if *optVersion {
			printVersion()
		} else {
			flag.Usage()
		}
	} else if len(args) == 1 {
		if args[0] == "help" || args[0] == "h" {
			flag.Usage()
		} else if args[0] == "lsdevice" {
			printDevices()
		} else if args[0] == "lscard" {
			printCards()
		} else {
			err := run(os.Args, cfg)
			if err != nil {
				logger.Errorf("Error running %v: %v", os.Args, err)
				fmt.Fprintf(os.Stderr, "Error running %v: %v\n", os.Args, err)
				os.Exit(1)
			}
		}
	} else {
		err := run(os.Args, cfg)
		if err != nil {
			logger.Errorf("Error running %v: %v", os.Args, err)
			fmt.Fprintf(os.Stderr, "Error running %v: %v\n", os.Args, err)
			os.Exit(1)
		}
	}
}
