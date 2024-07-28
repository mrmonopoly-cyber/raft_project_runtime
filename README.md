# Raft Project Runtime

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Starting the Cluster](#starting-the-cluster)
  - [Interacting with the Cluster](#interacting-with-the-cluster)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Overview

Raft Project Runtime implements the Raft consensus algorithm for distributed systems, ensuring consistency across a cluster of nodes. The project includes tools for cluster management, client interaction, environment setup, VM creation, and testing.

## Features

- **Consensus Algorithm**: Implements Raft for distributed consistency.
- **Cluster Management**: Tools to manage and monitor clusters.
- **Client API**: Interfaces for client-cluster interaction.
- **Environment Setup**: Scripts for development and deployment.
- **VM Support**: Automated setup for virtual machines.
- **Testing Framework**: Comprehensive suite for validation.

## Prerequisites

- **Go**: Install the Go programming language.
- **Shell**: Required for setup and management scripts.
- **VM Software**: VirtualBox or similar for VM management.

## Installation

1. **Clone the repository**:
   ```sh
   git clone https://github.com/mrmonopoly-cyber/raft_project_runtime.git

2. **Navigate to the project directory**:
   ```sh
   cd raft_project_runtime
   ```

3. **Set up the environment**:

Right now the automatic setup only works on Arch-Linux

   ```sh
   ./env_setup/setup.sh
   ```

## Usage

### Starting the Cluster

1. **Admin operations**:
   ```sh
   ./cluster_admin/bin/admin
   ```
2. **Client operations**:
   ```sh
   ./cluster_client/bin/client
   ```
## Contributing

We welcome contributions to the Raft Project Runtime. To contribute, please follow these steps:

1. **Fork the repository** on GitHub.
2. **Clone your forked repository** to your local machine:
   ```sh
   git clone https://github.com/your-username/raft_project_runtime.git
   ```
3. **Create a new branch for your feature or bug fix**:
   ```sh
   git checkout -b feature-branch
   ```
4. **Make your changes and commit them with a descriptive message**:
   ```sh
   git commit -am 'Add new feature'
   ```
5. **Push your changes to your forked repository**:
   ```sh
   git push origin feature-branch
   ```
6. **Create a Pull Request on the original repository**.

## License

This project is licensed under the Unlicense.

### Unlicense

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or distribute this software, either in source code form or as a compiled binary, for any purpose, commercial or non-commercial, and by any means.

For more details, see the [LICENSE](LICENSE) file.
