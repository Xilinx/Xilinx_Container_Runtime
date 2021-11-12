.. _build.rst:

Build and Install Xilinx Container Runtime from Source
------------------------------------------------------

Xilinx container runtime requires golang 17+ to build from the source code.

Get Source Code
...............

::

   git clone https://gitenterprise.xilinx.com/FaaSApps/Xilinx_Container_Runtime.git


Build and Install
.................

::

    cd Xilinx_Container_Runtime
    make
    sudo make install


Test
....

::

    xilinx-container-runtime lscard
