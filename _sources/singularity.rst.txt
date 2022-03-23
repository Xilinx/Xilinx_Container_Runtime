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

Singularity OCI Support
-----------------------

Singularity provides OCI runtime support to create OCI compilant container instance. Before starting a container, xilinx container runtime is able to modify the OCI specs to add xilinx devices. Here is a example.

Mount an OCI Filesystem Bundle
..............................

.. code-block:: bash

   singularity pull docker://xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04
   sudo singularity oci mount ./xilinx_runtime_base_alveo-2021.1-ubuntu-20.04.sif /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04


Modify OCI Specs
................

.. code-block:: bash

   sudo XILINX_VISIBLE_CARDS=0 xilinx-container-runtime modify -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04


Create OCI Compliant Container
..............................

.. code-block:: bash

   sudo singularity oci run -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04 xrt_base
