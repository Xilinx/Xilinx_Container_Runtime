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

.. _build.rst:

Install Xilinx Container Runtime
--------------------------------

Build and Install from Source Code
..................................

Get Source Code

::

   git clone https://github.com/Xilinx/Xilinx_Container_Runtime.git


Build and install xilinx container runtime requires golang 1.17.1+.

::

    cd Xilinx_Container_Runtime

    # Check required dependencies
    ./configure
    
    # Build binary
    make
    
    # Xilinx Container Runtime will be installed into /usr/bin by default
    sudo make install


To install Xilinx Contaienr Runtime in a different location, the destination directory can be specified while installing.

::

    sudo make install DESTDIR=/opt/xilinx/xcr/bin
    sudo export PATH=$PATH:/opt/xilinx/xcr/bin


Test
....

After successful installation, simply run xiliinx-container-runtime to get some usage information.

::

    xilinx-container-runtime
