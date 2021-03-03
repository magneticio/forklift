# Vamp Forklift command line client

Vamp Forklift is a command line client written in Go and allows to easily set up Clusters, Applications, Services and
Policies in Vamp.

## Table of Contents

================

- [Vamp Forklift command line client](#vamp-forklift-command-line-client)
    - [Table of Contents](#table-of-contents)
    - [Development](#development)
    - [Build](#build)
    - [Installation](#installation)
        - [Verifying installation](#verifying-installation)
    - [Usage](#usage)
        - [Clusters](#clusters)
        - [Applications](#applications)
        - [Services](#services)
        - [Policies](#policies)
        - [Release plans](#release-plans)

## Development

if you have golang installed, it is recommended to git clone Forklift to \$GOPATH/src/github.com/magneticio/forklift
This is a requirement for docker builder to work.

It is also recommended to read and follow golang setup for a development environment
setup: https://golang.org/doc/install

## Build

If you get errors about missing libraries while building, run:

```shell
GOPRIVATE=github.com/magneticio go get
```

for docker build:

```shell
./build.sh
```

for local build:

```shell
./build.sh local
```

binaries will be under bin directory

## Installation

If you have binaries built locally:
For mac run:

```shell
./bin/forklift-darwin-amd64 --help
```

If you have downloaded the binary directly, Just copy the binary for you platform to the user binaries folder for
general usage, for MacOS:

```shell
cp forklift-darwin-amd64 /usr/local/bin/forklift
chmod +x /usr/local/bin/forklift
```

If you don't have anything yet and automatically download an install, then follow the commands for your platform:

keep in mind that this installation may not work since this is a private repository. Manual installation is recommended.

Easy install for MacOS or Linux:

```shell
version=$(curl -s https://api.github.com/repos/magneticio/forklift/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') &&
  base=https://github.com/magneticio/forklift/releases/download/$version &&
  curl -sL $base/forklift-$(uname -s)-$(uname -m) >/usr/local/bin/forklift &&
  chmod +x /usr/local/bin/forklift
```

For general users it is recommended to download the binary for your platform. Latest release can be found here:
https://github.com/magneticio/forklift/releases/latest

Run get version so see if it is installed correctly:

```
forklift version
```

Now make sure to have a "config.yaml" configuration file in your home under ".forklift" folder, like the one shown
below, but with the correct parameters to connect to the key-value store.

```
key-value-store-url: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_URL}
key-value-store-token: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_TOKEN}
key-value-store-base-path: /secret/vamp/

```

The configuration path can be changed during the execution of any command by specifying the extra parameter

```shell
--config config-path
```

Where config-path is the path of the configuration file to be used.

Environment variables can be used in combination with the config. Environment variables overrides the configuration
file!

Environment variables:

```shell
  VAMP_FORKLIFT_PROJECT
    # Vamp Project ID
  VAMP_FORKLIFT_CLUSTER
    # Vamp Cluster ID
  VAMP_FORKLIFT_VAULT_ADDR
    #  Vault address. Example: http://vault.default.svc.cluster.local:8200
  VAMP_FORKLIFT_VAULT_TOKEN
    # Vault token
  VAMP_FORKLIFT_VAULT_BASE_PATH
    # Vault base path. Example: /secret/vamp/
  VAMP_FORKLIFT_VAULT_CACERT
    # Path of the CA Certificate.
  VAMP_FORKLIFT_VAULT_CLIENT_CERT
    # Path of the Client Certificate for TLS
  VAMP_FORKLIFT_VAULT_CLIENT_KEY
    # Path of the Client Certificate Key for TLS
```

Use export to setup environment variables (be careful about empty spaces) :

```shell
export VAMP_FORKLIFT_VAULT_ADDR="http://vault.default.svc.cluster.local:8200"
```

### Verifying installation

To verify the installation you can run the following command, which will return the version of the client's and vamp's
versions.

```shell
forklift version
```

It is possible to get all commands and flags by running help:

```shell
forklift help
```

## Usage

### Clusters

Forklift allows for the creation and update of clusters by running:

```shell
forklift put cluster 10 --name cluster-name --nats-channel-name name --optimiser-nats-channel-name optimiser-name --nats-token token
```

delete them with

```shell
forklift delete cluster 10
```

### Applications

Forklift allows for the creation and update of applications by running:

```shell
forklift put application 10 --namespace kubernetesNamespace --cluster 8
```

delete them with

```shell
forklift delete application 10 --cluster 8
```

### Services

Forklift allows for the creation and update of services by running:

```shell
forklift put service 10 --cluster 7 --file ./serviceconfig.json`
```

Example of service config:

```json
{
  "application_id": 5,
  "service_id": 10,
  "k8s_namespace": "test",
  "k8s_labels": {
    "app": "nginx-test"
  },
  "version_selector": "version",
  "default_policy_id": 1,
  "ingress_rules": [
    {
      "domain": "test.local",
      "path": "/",
      "port": 8081
    }
  ]
}
```

For headless services service config needs to have a `headless` flag set to true, for example:

```json
{
  "application_id": 5,
  "service_id": 10,
  "k8s_namespace": "test",
  "k8s_labels": {
    "app": "headless-test"
  },
  "version_selector": "version",
  "default_policy_id": 1,
  "ingress_rules": [],
  "headless": true
}
```

delete them with

```shell
forklift delete service 10 --cluster 7 --application 5
```

### Policies

Forklift allows for the creation and update of policies by running:

```shell
forklift put policy 10 --file ./policydefinition.json
```

Example of policy definition:

```json
{
  "type": "release",
  "id": 456,
  "version": 12,
  "steps": [
    {
      "source": {
        "weight": 100
      },
      "target": {
        "weight": 0,
        "condition": "cookie: vamp exists",
        "conditionStrength": 5
      },
      "endAfter": {
        "maxDuration": "1m0s"
      },
      "conditions": [
        {
          "metric": "Health",
          "budget": 70,
          "interval": {
            "type": "rolling",
            "duration": "5m0s"
          },
          "gracePeriod": "1m0s",
          "threshold": 1,
          "operator": "ge"
        },
        {
          "metric": "restarts",
          "budget": 5,
          "interval": {
            "type": "rolling",
            "duration": "5m0s"
          },
          "gracePeriod": "1m0s"
        },
        {
          "metric": "available replicas",
          "budget": 90,
          "interval": {
            "type": "rolling",
            "duration": "5m0s"
          },
          "gracePeriod": "1m0s",
          "threshold": 80,
          "operator": "ge"
        }
      ]
    },
    {
      "source": {
        "weight": 100
      },
      "target": {
        "weight": 0,
        "condition": "cookie: vamp exists",
        "conditionStrength": 10
      },
      "endAfter": {
        "maxDuration": "1m0s"
      },
      "conditions": [
        {
          "metric": "Health",
          "budget": 70,
          "interval": {
            "type": "rolling",
            "duration": "5m0s"
          },
          "gracePeriod": "1m0s",
          "threshold": 1,
          "operator": "ge"
        },
        {
          "metric": "restarts",
          "budget": 5,
          "interval": {
            "type": "rolling",
            "duration": "5m0s"
          },
          "gracePeriod": "1m0s"
        }
      ]
    }
  ],
  "metrics": [
    {
      "name": "restarts",
      "type": "EventBased",
      "value": {
        "source": "k8s-deployment-health",
        "type": "restarts"
      }
    },
    {
      "name": "available replicas",
      "type": "TimeBased",
      "value": {
        "source": "k8s-deployment-health",
        "type": "availablereplicas"
      }
    },
    {
      "name": "Health",
      "type": "TimeBased",
      "value": {
        "source": "k8s-deployment-health"
      }
    }
  ],
  "onSuccess": [
    {
      "type": "http",
      "value": {
        "url": "http://test.local",
        "httpRequest": "POST",
        "headers": [
          "authorization: Bearer xxxyyy",
          "content-type: text/plain"
        ]
      }
    }
  ]
}
```

delete them with

```shell
forklift delete policy 10
```

### Release plans

Release plans can be created with the following command:

```shell
forklift put releaseplan 1.0.1 --cluster 7 --application 6 --service 5 --file ./releaseplandefinition.json
```

Example of release plan:

```json
{
  "status": "not started",
  "service": {
    "name": "sava",
    "version": "1.0.5"
  },
  "releaseGroups": [
    {
      "group": 1,
      "name": "group-1",
      "status": "not started",
      "environments": [
        {
          "id": "test",
          "name": "test",
          "status": "not started"
        }
      ]
    }
  ]
}
```

Release plan can also be deleted with

```shell
forklift delete releaseplan 1.0.1 --cluster 7 --application 6 --service 5
```
