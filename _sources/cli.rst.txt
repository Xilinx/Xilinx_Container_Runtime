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

Command Line Tool
-----------------

Xilinx Container Runtime provides command line tool to show detailed Xilinx devices info on the host.
For some Xilinx FPGA cards, like U30, there is more than one device.
As shown in the following example, these two devices belong to the same Xilinx U30 card, with the card index 0. And you will be getting the indices of these two devices by 'lsdevice' command.

List Card(s)
............

.. code-block:: bash

    xilinx-container-runtime lscard
    CardIndex 	    SerialNum       DeviceBDF           UserPF                  MgmtPF                  ShellVersion
    0               XFL1YV0M20E0    0000:00:1e.0        /dev/dri/renderD128                             xilinx_u30_gen3x4_base_1
    0               XFL1YV0M20E0    0000:00:1f.0        /dev/dri/renderD129                             xilinx_u30_gen3x4_base_1


List Device(s)
..............

.. code-block:: bash

    xilinx-container-runtime lsdevice
    DeviceIndex     SerialNum       DeviceBDF           UserPF                  MgmtPF                  ShellVersion
    0               XFL1YV0M20E0    0000:00:1e.0        /dev/dri/renderD128                             xilinx_u30_gen3x4_base_1
    1               XFL1YV0M20E0    0000:00:1f.0        /dev/dri/renderD129                             xilinx_u30_gen3x4_base_1


Start a Container
.................

Based on previous information, environment variables can be set at the container starting process, so that the corresponding devices will be injected into the container.

Either 'XILINX_VISIBLE_DEVICES' or 'XILINX_VISIBLE_CARDS' can be passed, and acceptable values include 'all' and comma separated integers, like '0,1'.

.. code-block:: bash

   xilinx-container-runtime spec
   mkdir rootfs
   docker export $(docker create xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04) | tar -C rootfs -xvf -
   XILINX_VISIBLE_DEVICES=all xilinx-container-runtime run xrt_base
