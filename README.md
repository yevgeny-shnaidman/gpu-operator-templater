# GPU Operator Templater

GPU Operator Templater scaffolds a ready-to-build Kubernetes Operator that enables GPU devices for a specific vendor/device. It bootstraps an Operator repo using Operator SDK and renders a vendor-specific codebase that integrates with Kernel Module Management (KMM), with optional Node Labeller and Node Metrics components. The generated repo should require only minimal code changes before you can build an image and deploy on an OpenShift (OCP) or Kubernetes cluster.

## What it generates
- KMM-based module management to load a GPU kernel module and optionally deploy a device plugin
- A `DeviceConfig` CRD and API for configuring driver and device-plugin deployment
- Optional Node Labeller DaemonSet (labels GPU nodes)
- Optional Node Metrics DaemonSet (deploys user-provided metrics pods)
- Makefile targets to build, generate manifests, and deploy the Operator

## How it works
1. Uses an embedded `operator-sdk` to `init` a new repo and `create api` for `DeviceConfig`.
2. Renders templates into the target repo (adjusted for your vendor, API group/version, etc.).
3. Applies small `go.mod` adjustments automatically. You may still need to run `go mod tidy`.

## Prerequisites
- Linux amd64 environment
- Go 1.21+ installed for local builds
- Docker or compatible container runtime (for building images)
- Git and curl
- Access to a Kubernetes/OCP cluster for deployment (optional until you deploy)

## Build the templater
```bash
make templater
```
This builds the `templater` binary in the repo root.

## Quick start
1. Create an empty GitHub repository for your new Operator and clone it locally. The target directory must be empty.
2. Prepare a configuration file (see examples in `examples/` and the Configuration reference below).
3. From inside the empty target repo directory, run the templater:
   ```bash
   /path/to/gpu-operator-templater/templater -f /path/to/config.yaml
   ```
4. Post-generation steps (run inside your newly generated repo):
   ```bash
   go mod tidy
   make manifests
   make generate
   # Optional: run unit tests
   make unit-test
   ```
5. Build and push the Operator image (replace with your repository):
   ```bash
   make docker-build IMG=<your-registry/your-operator>:<tag>
   make docker-push IMG=<your-registry/your-operator>:<tag>
   ```
6. Install CRDs and deploy the Operator to your cluster:
   ```bash
   make install
   make deploy IMG=<your-registry/your-operator>:<tag>
   ```
7. Create a `DeviceConfig` instance to trigger driver/device-plugin deployment. Example skeleton (adjust to your `group`, `domain`, and `version`):
   ```yaml
   apiVersion: <group>.<domain>/<version>
   kind: DeviceConfig
   metadata:
     name: example
     namespace: <target-namespace>
   spec:
     # Optional overrides
     driversImage: <your-driver-image>
     driversVersion: <driver-version>
     devicePluginImage: <device-plugin-image>
     selector:
       node-role.kubernetes.io/worker: ""
   ```

## Configuration reference
Supply a YAML file via `-f` to the templater. See `examples/` for full samples.

Top-level keys:
- `api`:
  - `vendor` (string): Vendor name, lowercase (e.g., `nvidia`, `amd`, `intel`). Used for labels and import aliases.
  - `codeRepo` (string): Module path to use for the generated repo (e.g., `github.com/your-org/your-gpu-operator`).
  - `version` (string): API version for the CRD (e.g., `v1alpha1`).
  - `apiGroup` (string): Group for the API (e.g., `compgpu`).
  - `domain` (string): Domain for the API (e.g., `sigs.x-k8s.io`).
- `kmm`:
  - `pciVendorID` (string): PCI vendor ID used for node selection (targets nodes with that HW).
  - `kernelModuleName` (string): Name passed to `modprobe` for the GPU kernel module.
  - `enableDevicePlugin` (bool): Whether to deploy a device plugin via KMM.
  - `devicePluginImage` (string): Default device plugin image if not specified in the `DeviceConfig`.
  - `enableFirmware` (bool): Whether firmware files are required by the kernel module.
  - `imageFirmwarePath` (string): Path to firmware directory inside the driver image.
  - `driverVersion` (string): Default driver version, used as image tag for in-cluster builds.
  - `enableInClusterBuild` (bool): Enable KMM in-cluster driver image builds.
- `nodeLabeller` (optional, omit to disable):
  - `image` (string): Default image for node labeller DaemonSet.
- `nodeMetrics` (optional, omit to disable):
  - `image` (string): Default image for node metrics DaemonSet.
- `operatorImage` (string): Base image name (without tag) used by the generated Makefile (e.g., `quay.io/you/your-gpu-operator`).

Notes:
- Optional components are enabled by including their section. Omit `nodeLabeller` or `nodeMetrics` to disable them.
- The presence of `devicePluginImage` is honored only when `enableDevicePlugin` is true.

## Example configs
Minimal (KMM only):
```yaml
api:
  vendor: test
  codeRepo: github.com/yevgeny-shnaidman/test-gpu-operator
  version: v1alpha1
  apiGroup: compgpu
  domain: sigs.x-k8s.io
kmm:
  pciVendorID: "2040"
  kernelModuleName: testgpu
  devicePluginImage: rocm/k8s-device-plugin
  imageFirmwarePath: testFirmwareDir/updates
  driverVersion: el9-6.1.1
operatorImage: quay.io/yshnaidm/test-gpu-operator
```

With Node Labeller:
```yaml
api:
  vendor: test
  codeRepo: github.com/yevgeny-shnaidman/test-gpu-operator
  version: v1alpha1
  apiGroup: compgpu
  domain: sigs.x-k8s.io
kmm:
  pciVendorID: "2040"
  kernelModuleName: testgpu
  devicePluginImage: rocm/k8s-device-plugin
  imageFirmwarePath: testFirmwareDir/updates
  driverVersion: el9-6.1.1
nodeLabeller:
  image: rocm/k8s-device-plugin:labeller-latest
operatorImage: quay.io/yshnaidm/test-gpu-operator
```

With Node Metrics:
```yaml
api:
  vendor: test
  codeRepo: github.com/yevgeny-shnaidman/test-gpu-operator
  version: v1alpha1
  apiGroup: compgpu
  domain: sigs.x-k8s.io
kmm:
  pciVendorID: "2040"
  kernelModuleName: testgpu
  devicePluginImage: rocm/k8s-device-plugin
  imageFirmwarePath: testFirmwareDir/updates
  driverVersion: el9-6.1.1
nodeMetrics:
  image: rocm/k8s-node-metrics
operatorImage: quay.io/yshnaidm/test-gpu-operator
```

You can find these under `examples/` as:
- `config-kmm-only.yaml`
- `config-with-kmm-nodelabeller.yaml`
- `config-with-kmm-nodemetrics.yaml`
- `config.yaml` (full example)

## After generation: common tasks in the new repo
- Generate and install manifests: `make manifests && make install`
- Generate mocks and deep-copies: `make generate`
- Build/push controller image: `make docker-build IMG=... && make docker-push IMG=...`
- Deploy the controller: `make deploy IMG=...`
- Bundle and catalog (optional): `make bundle bundle-build`

## Troubleshooting and tips
- If `go.mod` contains unexpected versions, run `go mod tidy` and adjust as needed.
- Ensure your `IMG`/`operatorImage` point to a registry you can push to.
- `enableInClusterBuild` requires appropriate KMM configuration and build environment in-cluster.
- This tool scaffolds a baseline. Depending on your hardware and driver packaging, you may need to extend the controller logic or templates in the generated repo.
