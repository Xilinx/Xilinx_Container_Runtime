Command Line Tool
-----------------

Xilinx Container Runtime provides command line tool to show detailed Xilinx devices info on the host.
For some Xilinx FPGA cards, like U30, there is more than one device. Details are as below.

List Card(s)
............

.. code-block:: bash

    xilinx-container-runtime lscard
    CardNum 	    SerialNum       DeviceBDF           UserPF                  MgmtPF                  ShellVersion
    0               XFL1YV0M20E0    0000:00:1e.0        /dev/dri/renderD128                             xilinx_u30_gen3x4_base_1
    0               XFL1YV0M20E0    0000:00:1f.0        /dev/dri/renderD129                             xilinx_u30_gen3x4_base_1


List Device(s)
..............

.. code-block:: bash

    xilinx-container-runtime lsdevice
    DeviceNum 	    SerialNum       DeviceBDF           UserPF                  MgmtPF                  ShellVersion
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
