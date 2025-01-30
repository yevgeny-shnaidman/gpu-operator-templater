# GPU operator templater

GPU operator templater allows creating code template that implements a GPU operator code for a specific vendor/device. The created code base will require additional updates for the code, but they should be minimal. For a generic GPU operator the created code base should provide the immidiate ability to compile, create image and run on a OCP cluster

## Code components

### KMM module
This module is configured to load the GPU kernel module. Currently the created code provides the ability to load pre-compiled kernel module and to deploy the predefined device-plugin image. In case additional functionlity needs to be supported (in-cluster builds, in-cluster signing, multiple kernel modules loading) the code in internal/kmmmodule/kmmmodule.go needs to be updated

### Node Labeller 
Node labeller is an image that is running on every node that the kernel module is deployed on. It is supposed to label the nodes with the user's specified labels. The image for node labeller is provided by the user

### Node mertrics
Node metrix deployes a user's metrix pods on the kernel module's nodes. In the feature we should probably provide additional rules and port mappings in order for it actually to work

## Creating operator code
Follow the follwing steps to create the GPU operator code:
1. Create a github repo and clone it into your server. Note: the cloned repo must be empty
2. Create templater executable and prepare the configuration file (see sections "Templater executable" and "Templater configuration")
3. Switch into the cloned directory
4. Run the templater with configuration file: <templater path> -f <configuration file path>
5. Run: go mod tidy - this will prepare the go.mod and go.sum files
6. Run: make manifest - this will prepare the CRDs yamls, rbac and everything needed for deployment in the config directory
7. Run: make generate - this will generate the mocks for unit test

## Templater configuration
1. vendor - string that represent the vendor name. For example: test, nvidia, amd, intel, ibm, etc'
2. codeRepo - github repository that is used for the operator code: For example: github.com/yevgeny-shnaidman/test-gpu-operator
3. apiVersion - the version of the CRD API that the gpu operator will be using: For example: v1alpha1, v1beta1...
4. group - the group part of the CRD api
5. domain - the domain part of CRD api
6. pciVendorID - the PCI vendor id of the GPU device. Will be used as a part of the node selector field of the KMM Module, to schedule worker pod and device plugin pods only on the nodes containing that HW
7. kernelModuleName -  the name of the kernel module to load, as it would be passed to modprobe command
8. defaultDevicePluginImage - the image to be used for the DevicePlugin pods
9. imageFirmwarePath - the path of the firmware directory inside the kernel module image. Will be used by KMM Module
10. defaultDriverVersion - the default version of the GPU driver to be used. Will be used a a tag to the module container image created by KMM during in-cluster build
11. defaultNodeLabellerImage - the image that will be used in the NodeLabeller daemonset
12. operatorImage - the operator image that will be created by the Makefile command
