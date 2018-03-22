

# terraform-provider-vcloud-director

## Overview

This repository implements a Terraform provider to work with [VMware vCloud Director](https://www.vmware.com/products/vcloud-director.html) resources.

***Terraform*** is a tool for building, changing, and versioning infrastructure safely and efficiently. Terraform can manage existing and popular service providers as well as custom in-house solutions.

***Terraform Provider*** is a tool which is based on Terraform and works with configuration files. These configuration files contain all the information which is relevant enough to perform ***create/update/read/delete*** operations on the resources available on vCloud Director.

## Architecture

***Terraform Provider*** has been developed using Python and GO. It uses Client-Server model inside the hood where the client has been written using GO and server has been written using Python language. The core reason to use two different languages is to make a bridge between Terraform and Pyvcloud API. Pyvcloud is the SDK developed by VMware and provides an medium to talk to vCloud Director. Terraform uses GO to communicate where Pyvcloud has been written in Python3.

We are using GRPC Protocol to handle the communication between GO client and Python server.

![alt text](https://raw.githubusercontent.com/vmware/terraform-provider-vcloud-director/master/docs/architecture.jpg)

## Try it out

See the [Setup.md](https://github.com/vmware/terraform-provider-vcloud-director/blob/master/docs/SETUP.md)
See the [Release Docs](https://vmware.github.io/terraform-provider-vcloud-director/)

## Contributing

The terraform-provider-vcloud-director project team welcomes contributions from the community. Before you start working with terraform-provider-vcloud-director, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[BSD-2](LICENSE.txt)
