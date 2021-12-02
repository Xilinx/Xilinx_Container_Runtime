.. _build.rst:

Install Xilinx Container Runtime
--------------------------------

Install the RPM package
.......................

::

   wget https://.../xilinx-container-runtime_version.rpm
   sudo yum reinstall ./xilinx-container-runtme_version.rpm


Install the DEB package
.......................

::

   wget https://.../xilinx-container-runtime_version.deb
   sudo apt install  --reinstall ./xilinx-container-runtme_version.deb



Build and Install from source code
..................................

Get Source Code

::

   git clone https://gitenterprise.xilinx.com/FaaSApps/Xilinx_Container_Runtime.git


Build and install xilinx container runtime requires golang 17.1+.

::

    cd Xilinx_Container_Runtime
    make
    sudo make install


Test
....

After successful installation, simply run xiliinx-container-runtime to get some usage information.

::

    xilinx-container-runtime
