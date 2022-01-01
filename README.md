<!--
 Copyright (C) 2021, Xilinx Inc - All rights reserved
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
# Xinlinx-Container-Runtime

Xilinx-container-runtime is an extension of runc, with modification to add xilinx devices before running containers.

## Build
    
    make

## Install
    
    sudo make install

## Usage:
    
    xilinx-container-runtime spec
    mkdir rootfs
    docker export $(docker create xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04) | tar -C rootfs -xvf -
    XILINX_VISIBLE_DEVICES=all xilinx-container-runtime run xrt_base
    XILINX_VISIBLE_CARDS=0 xilinx-container-runtime run xrt_base

## Integrate with docker:

Update or create /etc/docker/daemon.json as below:

    {
        "runtimes": {
            "xilinx": {
                "path": "/usr/bin/xilinx-container-runtime",
                "runtimeArgs": []
            }
        }
    }

After updating /etc/docker/daemon.json, it's required to restart docker service to register xilinx-container-runtime.

    sudo systemctl restart docker

    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=0 xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=all 900497858702.dkr.ecr.us-west-2.amazonaws.com/faas-test:GA1.5-al2-rc5-public /bin/bash
    sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 900497858702.dkr.ecr.us-west-2.amazonaws.com/faas-test:GA1.5-al2-rc5-public /bin/bash
    docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all -e XILINX_DEVICE_EXLUSIVE=true xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
    
## Integrate with podman:

Specify the runtime by absolute file path:

    sudo podman run -it --rm --runtime=/usr/bin/xilinx-container-runtime -e XILINX_VISIBLE_CARDS=0 docker.io/xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash

Or, to make it simpler, update /usr/share/containers/containers.conf, in part [engine.runtimes], adding below part:

    xilinx = [
        "/usr/bin/xilinx-container-runtime"
    ]

Then you can specify the runtime as xilinx.
    
    sudo podman run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 docker.io/xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash


## Singularity

    singularity pull docker://xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04
    sudo singularity oci mount ./xilinx_runtime_base_alveo-2021.1-ubuntu-20.04.sif /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04
    sudo XILINX_VISIBLE_CARDS=0 xilinx-container-runtime modify -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04
    sudo singularity oci run -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04 xrt_base
    