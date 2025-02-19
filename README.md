# Canonical Kubernetes Snap
[![End to End Tests](https://github.com/canonical/k8s-snap/actions/workflows/e2e.yaml/badge.svg)](https://github.com/canonical/k8s-snap/actions/workflows/e2e.yaml)

**Canonical Kubernetes** is the fastest, easiest way to deploy a fully-conformant Kubernetes cluster. Harnessing pure upstream Kubernetes, this distribution adds the missing pieces (e.g. ingress, dns, networking) for a zero-ops experience.

For more information and instructions, please see the official documentation at: https://ubuntu.com/kubernetes

## Quickstart

Install Canonical Kubernetes and initialise the cluster with:

```bash
sudo snap install k8s --edge --classic
sudo k8s bootstrap
```

Confirm the installation was successful:

```bash
sudo k8s status
```

Use `kubectl` to interact with k8s:

```bash
sudo k8s kubectl get pods -A
```

Remove the snap with:

```bash
sudo snap remove k8s --purge
```


## Build the project from source

To build the Kubernetes snap on an Ubuntu machine you need Snapcraft.

```bash
sudo snap install snapcraft --classic
```

Building the project by running `snapcraft` in the root of this repository. Snapcraft spawns a VM managed by Multipass and builds the snap inside it. If you don’t have Multipass installed, snapcraft will first prompt for its automatic installation.

After snapcraft completes, you can install the newly compiled snap:

```bash
sudo snap install k8s_*.snap --classic --dangerous
```
