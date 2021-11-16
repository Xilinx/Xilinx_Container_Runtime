Integration with Singularity OCI
--------------------------------

Singularity provides OCI runtime support to create OCI compilant container instance. Xilinx container runtime is able to integrate with OCI command to add xilinx devices.

Mount an OCI Filesystem Bundle
..............................

.. code-block:: bash

   singularity pull docker://xilinx/xilinx_runtime_base:alveo-2021.1-ubuntu-20.04
   sudo singularity oci mount ./xilinx_runtime_base_alveo-2021.1-ubuntu-20.04.sif /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04


Modify OCI Configuration
........................

.. code-block:: bash

   sudo XILINX_VISIBLE_CARDS=0 xilinx-container-runtime modify -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04


Create OCI Compliant Container
..............................

.. code-block:: bash

   sudo singularity oci run -b /var/tmp/xilinx_runtime_base_alveo-2021.1-ubuntu-20.04 xrt_base
