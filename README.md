<!--
 Copyright (C) 2022, Xilinx Inc - All rights reserved
 Xilinx Container Runtime
 
 Licensed under the Apache License, Version 2.0 (the "License"). You may
 not use this file except in compliance with the License. A copy of the
 License is located at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 License for the specific language governing permissions and limitations
 under the License. 
-->
# Xilinx-Container-Runtime

Xilinx-container-runtime is an extension of runc, with modification to add xilinx devices before running containers.

## Build
    
    make

## Install
    
    sudo make install

## Usage
    
By running xilinx-container-runtime, it is required to have runC installed. If docker has been installed, the runC is installed already.

    # Create OCI specs
    xilinx-container-runtime spec

    # Create root file system for container
    mkdir rootfs
    
    # Export a docker image as the root file system
    docker export $(docker create xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04) | tar -C rootfs -xvf -
    
    # Starting a container and inject xilinx devices
    # Accepted value for XILINX_VISIBLE_DEVICES include 0,1,...|all 
    XILINX_VISIBLE_DEVICES=0 xilinx-container-runtime run xrt_base
    # Accepted value for XILINX_VISIBLE_CARDS include 0,1,...|all
    XILINX_VISIBLE_CARDS=0 xilinx-container-runtime run xrt_base

## Integrate with docker

Docker allows users to add custom runtime by editing /etc/docker/daemon.json in "runtimes" part. Below is a code snap to add xilinx container runtime. 

    {
        "runtimes": {
            "xilinx": {
                "path": "/usr/bin/xilinx-container-runtime",
                "runtimeArgs": []
            }
        }
    }

After updating /etc/docker/daemon.json, it is required to restart docker service to register xilinx-container-runtime.

    sudo systemctl restart docker

    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=0 xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=all 900497858702.dkr.ecr.us-west-2.amazonaws.com/faas-test:GA1.5-al2-rc5-public /bin/bash
    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 900497858702.dkr.ecr.us-west-2.amazonaws.com/faas-test:GA1.5-al2-rc5-public /bin/bash

For docker usage, we can enable device exclusive mode, which allows devices to be assigned to one container exclusively. It is easy to be enabled by setting the environment variable 'XILINX_DEVICE_EXCLUSIVE' to 'true'.

    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all -e XILINX_DEVICE_EXLUSIVE=true xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
    
## Integrate with podman

Specify the runtime by absolute file path:

    sudo podman run -it --rm --runtime=/usr/bin/xilinx-container-runtime -e XILINX_VISIBLE_CARDS=0 docker.io/xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash

Or, to make it simpler, update /usr/share/containers/containers.conf, in part [engine.runtimes], adding below part:

    xilinx = [
        "/usr/bin/xilinx-container-runtime"
    ]

Then you can specify the runtime as xilinx.
    
    sudo podman run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 docker.io/xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash


## Singularity OCI Support

Currently, xilinx container runtime partly supports singularity. Before running a singularity OCI container, OCI specs can be modified by 'modify' command.

    singularity pull docker://xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04
    sudo singularity oci mount ./xilinx_runtime_base_alveo-2021.1-ubuntu-20.04.sif /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04
    sudo XILINX_VISIBLE_CARDS=0 xilinx-container-runtime modify -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04
    sudo singularity oci run -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04 xrt_base
    