# GPU operator templater

GPU operator templater allows creating code template that implements a GPU operator code for a specific vendor/device. The created code base will require additional updates for the code, but they should be minimal. For a generic GPU operator the created code base should provide the immidiate ability to compile, create image and run on a OCP cluster

## Code components

### KMM module
This module is configured to load the GPU kernel module. Currently the created code provides the ability to load pre-compiled kernel module and to deploy the predefined device-plugin image. In case additional functionlity needs to be supported (in-cluster builds, in-cluster signing, multiple kernel modules loading) the code in internal/kmmmodule/kmmmodule.go needs to be updated

### Node Labeller 
Node labeller is an image that is running on every node that the kernel module is deployed on. It is supposed to label the nodes with the user's specified labels. The image for node labeller is provided by the user

### Node metrics
Node metrix deployes a user's metrix pods on the kernel module's nodes. In the feature we should probably provide additional rules and port mappings in order for it actually to work

## Running templater
There are 2 steps that needs to be execcuted in order to create the initial code: creating the code and adjusting the go.mod

### Creating the code
1. Create a github repo and clone it into your server. Note: the cloned repo must be empty
2. Create templater executable and prepare the configuration file (see sections "Templater executable" and "Templater configuration")
3. Switch into the cloned directory
4. Run the templater with configuration file: <templater path> -f <configuration file path>

### Adjusting go.mod
Adjusting the go.mod needs to be manually, since there are a lot of variables that the go.mod is dependent upon: Golang version, KMM version etc'
1. Run: go mod tidy - this will prepare the go.mod and go.sum files. If needed adjust the go.mod accordingly
2. Run: make manifest - this will prepare the CRDs yamls, rbac and everything needed for deployment in the config directory
3. Run: make generate - this will generate the mocks for unit test

## Templater configuration
Templater configuration is supplied to the templater as an input during templater execution. 
It allows to configure the various components of the operator code and environment.
The configuration file is an YML file that is divided into 4 subcomponents:
API - contains configuration needed for CRD definitions and initialiation of the operator-sdk
KMM - contains confgiuration used by KMM Module
NodeLabeller - contains confgiuration used by NodeLabeller component (optional)
NodeMetrics - contains confgiuration used by NodeMetrics component (optional)

### API component
1. vendor - string that represent the vendor name. For example: test, nvidia, amd, intel, ibm, etc'. Used for configuring labels and imports re-naming
2. codeRepo - github repository that is used for the operator code: For example: github.com/yevgeny-shnaidman/test-gpu-operator
3. version - the API version of the CRD API that the gpu operator will be using: For example: v1alpha1, v1beta1...
4. group - the group part of the CRD api. For example: compgpu
5. domain - the domain part of CRD api. For example: sigs.x-k8s.io

### KMM Component
1. pciVendorID - the PCI vendor id of the GPU device. Will be used as a part of the node selector field of the KMM Module, to schedule worker pod and device plugin pods only on the nodes containing that HW
2. kernelModuleName -  the name of the kernel module to load, as it would be passed to modprobe command
3. enableDevicePlugin - should be set to true, if the device plugin part of KMM should be used
4. devicePluginImage - the image to be used for the DevicePlugin pods.Applicable only if enableDevicePlugin is true
5. enableFirmware - should be set to true, if the kernel module requires the firmware configuration
6. imageFirmwarePath - the path of the firmware directory inside the kernel module image. Applicable only if  imageFirmwarePath is true 
7. driverVersion - the default version of the GPU driver to be used. Will be used a a tag to the module container image created by KMM during in-cluster build

### NodeLabeller component
1. enable - should be set to true, if operator needs to deploy node labeller component.
2. image - the image of the node labeller, that will be deployed by a dedicated DaemonSet

## NodeMetrics
1. enable - should be set to true, if operator needs to deploy node metrics component.
2. image - the image of the node metrics, that will be deployed by a dedicated DaemonSet

### Templater executable
Currently templater executable needs to be built by the user. The following will build the executable
make templater
