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

Release Note
============

Version 1.0.0
-------------

Xilinx-container-runtime is an extension of runc, with modification to add xilinx devices before running containers. Xilinx container runtime provides command line tool to list Xilinx device(s) or card(s) on host machine, and interface to map Xilinx devices into containers per card or device level. Also, it is easy to integrate Xilinx container runtime with Docker and Podman.