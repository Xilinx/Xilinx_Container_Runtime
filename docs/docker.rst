Integration with Docker
-----------------------

Xilinx container runtime was designed to integrate with docker easily.

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

After updating /etc/docker/daemon.json, it's required to restart docker service to register xilinx-container-runtime.

.. code-block:: bash

    sudo systemctl restart docker


Start a Container
.................

Xilinx container runtime can be specified using --runtime flag.

Notice: If the default runtime was set in /etc/docker/daemon.json, --runtime flag can be omitted.

.. code-block:: bash

   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=0 xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=all xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash

Assign Device Exclusively
.........................

For docker usage, we can enable device exclusive mode, which allows devices to be assigned to one container exclusively. It is easy to be enabled by setting the environment variable 'XILINX_DEVICE_EXCLUSIVE' to 'true'.

.. code-block:: bash

   sudo docker run -it --rm --runtime=xilinx -e XILINX_VISIBLE_DEVICES=all -e XILINX_DEVICE_EXLUSIVE=true xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash