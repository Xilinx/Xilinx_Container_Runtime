.. 
   Copyright (C) 2022, Xilinx Inc - All rights reserved
  
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
  
       http://www.apache.org/licenses/LICENSE-2.0
  
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

Integration with Docker
-----------------------

Xilinx container runtime is designed to integrate with docker easily.

Add Xilinx Container Runtime into Docker
........................................

Docker is allowed to be configured by /etc/docker/daemon.json, the following code snippet adds xilinx contianer runtime. 

.. code-block:: json

    {
        "runtimes": {
            "xilinx": {
                "path": "/usr/bin/xilinx-container-runtime",
                "runtimeArgs": []
            }
        }
    }

Optionally, you can set xilinx as the default container runtime for docker.

.. code-block:: json

    {
        "default-runtime": "xilinx",
        "runtimes": {
            "xilinx": {
                "path": "/usr/bin/xilinx-container-runtime",
                "runtimeArgs": []
            }
        }
    }

Restart Docker Service
......................

After updating /etc/docker/daemon.json, it's required to restart docker service for the registration of xilinx-container-runtime being effective.

.. code-block:: bash

    sudo systemctl restart docker


Start a Container
.................

Xilinx container runtime can be specified using --runtime flag. If the default runtime was set in /etc/docker/daemon.json, --runtime flag can be omitted.

For environment variables XILINX_VISIBLE_DEVICES and XILINX_VISIBLE_CARDS, the acceptable values include 'all' and comma separated integers of card or device index which can be got from previous command line tools.

**Note: For docker usage, we are using device exclusive mode by default, which assigns some device only to one container exclusively. In this mode, a device will be locked to the specific container from the time of container being created, instead of the time of container being started, till the container is removed. The device would still be locked, if the container is only stopped, but not removed.**

.. code-block:: bash

   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=0,1 xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=all xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash

Disable Device Exclusive Mode
.............................

It is easy to disable the device exclusive mode by setting the environment variable 'XILINX_DEVICE_EXCLUSIVE' to 'false'.

.. code-block:: bash

   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all -e XILINX_DEVICE_EXCLUSIVE=false xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash