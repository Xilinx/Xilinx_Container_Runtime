Integration with Podman
-----------------------

Podman is an alternative to docker, with which xilixn container runtime is able to integrate.

Run with Absolute File Path
...........................

Podman is allowed to specify runtime by absolute file path.

.. code-block:: bash

   sudo podman run -it --rm --runtime=/usr/bin/xilinx-container-runtime -e XILINX_VISIBLE_CARDS=0 docker.io/xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash


Configure Runtime for Podman
............................

Optionally, you can update /usr/share/containers/containers.conf, in [engine.runtimes] part, adding below snippet.

.. code-block:: bash

   xilinx = [
      "/usr/bin/xilinx-container-runtime"
   ]

Then, you can specify the runtime as xilinx.

.. code-block:: bash

   sudo podman run -it --rm --runtime=xilinx -e XILINX_VISIBLE_CARDS=0 docker.io/xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04 /bin/bash
