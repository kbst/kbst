<p align="center">
 <img src="./assets/favicon.png" alt="Kubestack, The Open Source Gitops Framework" width="25%" height="25%" />
</p>

<h1 align="center">Kubestack CLI</h1>
<h3 align="center">CLI for the Kubestack GitOps Framework</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues/kbst/kbst.svg)](https://github.com/kbst/kbst/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/kbst/kbst.svg)](https://github.com/kbst/kbst/pulls)

</div>

<div align="center">

![GitHub Repo stars](https://img.shields.io/github/stars/kbst/kbst?style=social)
![Twitter Follow](https://img.shields.io/twitter/follow/kubestack?style=social)

</div>


<h3 align="center"><a href="#Contributing">Join Our Contributors!</a></h3>

<div align="center">

<a href="https://github.com/kbst/kbst/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=kbst/kbst&max=36" />
</a>

</div>

## Introduction

`kbst` is an all-in-one CLI to scaffold your Infrastructure as Code repository and deploy your entire platform stack locally for faster iteration.

This CLI is part of the [Kubestack Gitops framework](https://www.kubestack.com) for managing Kubernetes services based on Terraform and Kustomize.

`kbst` improves the GitOps developer experience by making a number of common tasks easier.
> The CLI never makes any changes to any cloud environment. All changes are exclusively to the local environment and working directory.

### Features
* Scaffolds an Infrastructure as Code (IaC) repository ready to deploy to your target cloud environment
* Provisions a local environment mirroring your cloud infrastructure using Kubernetes in Docker (KinD)
* Watches for changes in configuration files and re-deploys parts of infrastructure to immediately reflact updates
* Provides Docker container including all tooling for interacting manually your cloud environment as needed


## Installing kbst

### Linux
```
# Download the latest release
curl -LO "https://github.com/kbst/kbst/releases/download/$(curl -s https://www.kubestack.com/cli-latest.txt)/kbst_linux_amd64.zip"

# Extract the binary into your PATH e.g. /usr/local/bin
sudo unzip -d /usr/local/bin/ kbst_linux_amd64.zip kbst

# Verify the binary works
kbst --version
```

### MacOS
```
# Download the latest release
curl -LO "https://github.com/kbst/kbst/releases/download/$(curl -s https://www.kubestack.com/cli-latest.txt)/kbst_darwin_amd64.zip"

# Extract the binary into your PATH e.g. /usr/local/bin
sudo unzip -d /usr/local/bin/ kbst_darwin_amd64.zip kbst

# Verify the binary works
kbst --version
```

### Windows
> Windows instructions require WSL2 and Docker Desktop for Windows with the WSL2 backend.
```
# Download the latest release
curl -LO "https://github.com/kbst/kbst/releases/download/$(curl -s https://www.kubestack.com/cli-latest.txt)/kbst_linux_amd64.zip"

# Extract the binary into your PATH e.g. /usr/local/bin
sudo unzip -d /usr/local/bin/ kbst_linux_amd64.zip kbst

# Verify the binary works
kbst --version
```

### Build from Source
```
# Clone the repository
git clone https://github.com/cctechwiz-forks/kbst.git

# Enter repository directory
cd kbst

# Install binary (requires you $GOPATH/bin to be on your system $PATH)
make install
```


## Using the kbst CLI

`kbst` has four commands:
* `help` - Help about any command
* `local` - Start a localhost development environment
* `manifest` - Add, update, and remove services from the [catalog](https://github.com/kbst/catalog)
* `repository` - Create and change Kubestack repositories

`kbst local` has two sub-commands:
* `apply` - Watch and apply changes to the localhost development environment
* `destroy` - Destroy the localhost development environment

`kbst manifest` has three sub-commands:
* `install` - Install and vendor a manifest from the catalog
* `remove` - Remove a vendored manifest from all environments
* `update` - Update vendored manifests from the catalog

`kbst repository` only has a single sub-command:
* `init` - Scaffold a new repository


## Getting Started with Kubestack

For the easiest way to get started, [visit the official Kubestack quickstart](https://www.kubestack.com/infrastructure/documentation/quickstart). This tutorial will help you get started with the Kubestack GitOps framework. It is divided into three steps.

1. Develop Locally
    * Scaffold your repository and tweak your config in a local development environment that simulates your actual cloud configuration using Kubernetes in Docker (KinD).
3. Provision Infrastructure
    * Set-up cloud prerequisites and bootstrap Kubestack's environment and clusters on your cloud provider for the first time.
4. Set-up Automation
    * Integrate CI/CD to automate changes following Kubestack's GitOps workflow.


## Getting Help

**Official Documentation**  
Refer to the [official documentation](https://www.kubestack.com/framework/documentation) for a deeper dive into how to use and configure Kubetack.

**Community Help**  
If you have any questions while following the tutorial, join the [#kubestack](https://app.slack.com/client/T09NY5SBT/CMBCT7XRQ) channel on the Kubernetes community. To create an account request an [invitation](https://slack.k8s.io/).

**Professional Services**  
For organizations interested in accelerating their GitOps journey, [professional services](https://www.kubestack.com/lp/professional-services) are available.


## Contributing
Contributions to the Kubestack framework are welcome and encouraged. Before contributing, please read the [Contributing](./CONTRIBUTING.md) and [Code of Conduct](./CODE_OF_CONDUCT.md) Guidelines.

One super simple way to contribute to the success of this project is to give it a star.  

<div align="center">

![GitHub Repo stars](https://img.shields.io/github/stars/kbst/kbst?style=social)

</div>


## Kubestack Repositories
* [kbst/terraform-kubestack](https://github.com/kbst/terraform-kubestack)  
    * Terraform GitOps Framework - Everything you need to build reliable automation for AKS, EKS and GKE Kubernetes clusters in one free and open-source framework.
* [kbst/kbst](https://github.com/kbst/kbst) (this repository)  
    * Kubestack Framework CLI - All-in-one CLI to scaffold your Infrastructure as Code repository and deploy your entire platform stack locally for faster iteration.
* [kbst/terraform-provider-kustomization](https://github.com/kbst/terraform-provider-kustomization)  
    * Kustomize Terraform Provider - A Kubestack maintained Terraform provider for Kustomize, available in the [Terraform registry](https://registry.terraform.io/providers/kbst/kustomization/latest).
* [kbst/catalog](https://github.com/kbst/catalog)  
    * Catalog of cluster services as Kustomize bases - Continuously tested and updated Kubernetes services, installed and customizable using native Terraform syntax.

