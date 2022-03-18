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

Integration with Podman
-----------------------

Podman is an alternative to docker, with which xilinx container runtime is able to integrate.

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
