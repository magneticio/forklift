# Vamp Forklift command line client

Vamp Forklift is a command line client written in golang and allows to easily set up Organizations and Environments in Vamp.

Forklift requires running and reachable instances of MySql and Vault tied to a Vamp installation.
Please check How to Setup Vamp at the following link https://vamp.io/documentation/installation/v1.0.0/overview/

## Table of Contents
================
- [Vamp Forklift command line client](#vamp-forklift-command-line-client)
- [## Table of Contents](#h2-id%22table-of-contents-16%22table-of-contentsh2)
  - [Development](#development)
  - [Build](#build)
  - [Installation](#installation)
    - [Verifying installation](#verifying-installation)
  - [Usage](#usage)
    - [Organizations](#organizations)
    - [Users](#users)
    - [Environments](#environments)
    - [Artifacts](#artifacts)
    - [Release policy](#release-policy)
    - [Release plan](#release-plan)

## Development

if you have golang installed, it is recommended to git clone Forklift to $GOPATH/src/github.com/magneticio/forklift
This is a requirement for docker builder to work.

It is also recommended to read and follow golang setup for a development environment setup: https://golang.org/doc/install

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

If you have downloaded the binary directly, Just copy the binary for you platform to the user binaries folder for general usage, for MacOS:

```shell
cp forklift-darwin-amd64 /usr/local/bin/forklift
chmod +x /usr/local/bin/forklift
```

If you don't have anything yet and automatically download an install, then follow the commands for your platform:

keep in mind that this installation may not work since this is a private repository.
Manual installation is recommended.

Easy install for MacOS or Linux:
```shell
version=$(curl -s https://api.github.com/repos/magneticio/forklift/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') &&
  base=https://github.com/magneticio/forklift/releases/download/$version &&
  curl -sL $base/forklift-$(uname -s)-$(uname -m) >/usr/local/bin/forklift &&
  chmod +x /usr/local/bin/forklift
```

For general users it is recommended to download the binary for your platform.
Latest release can be found here:
https://github.com/magneticio/forklift/releases/latest

Run get version so see if it is installed correctly:
```
forklift version
```

Now make sure to have a "config.yaml" configuration file in your home under ".forklift" folder, like the one shown below, but with the correct parameters to connect to the database and the key-value store.

```
namespace: vampio
database-enabled: true
database-type: mysql
database-name: vamp-${namespace}
database-url: jdbc:mysql://mysql.default.svc.cluster.local:3306/vamp-${namespace}
database-user: root
database-table: ${namespace}
database-password: secret
key-value-store-url: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_URL}
key-value-store-token: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_TOKEN}
key-value-store-base-path: /secret/vamp/${namespace}
key-value-store-type: vault

```

Mind the fact that setting database-enabled to false (which is also the default value), will disable the database regardless of the values specified in other database related fields.
The configuration path can be changed during the execution of any command by specifying the extra parameter

```shell
--config config-path
```

Where config-path is the path of the configuration file to be used.



Environment variables can be used in combination with the config.
Environment variables overrides the configuration file!

Environment variables:
```shell
  VAMP_FORKLIFT_VAULT_ADDR
    #  Vault address. Example: http://vault.default.svc.cluster.local:8200
  VAMP_FORKLIFT_VAULT_TOKEN
    # Vault token
  VAMP_FORKLIFT_VAULT_CACERT
    # Path of the CA Certificate.
  VAMP_FORKLIFT_VAULT_CLIENT_CERT
    # Path of the Client Certificate for TLS
  VAMP_FORKLIFT_VAULT_CLIENT_KEY
    # Path of the Client Certificate for TLS
  VAMP_FORKLIFT_MYSQL_HOST
    # MySql host address. Example mysql.default.svc.cluster.local:3306
  VAMP_FORKLIFT_MYSQL_CONNECTION_PROPS
    Parameters to use in combination with MySql Url. Example: useSSL=false
  VAMP_FORKLIFT_MYSQL_USER
    # MySql username
  VAMP_FORKLIFT_MYSQL_PASSWORD
    # MySql password
```

Use export to setup environment variables (be careful about empty spaces) :
```shell
export VAMP_FORKLIFT_VAULT_ADDR="http://vault.default.svc.cluster.local:8200"
```

### Verifying installation

To verify the installation you can run the following command, which will return the version of the client's and vamp's versions.

```shell
forklift version
```

It is possible to get all commands and flags by running help:
```shell
forklift help
```

## Usage

Notes: 
  * Organization and environment names should be lowercase alphanumeric, please remove "-" while running examples and use a name proper for you.
  * Organization, Environment, User and Arfifact operations require SQL to be enabled.

### Organizations

Forklift allows for the creation of a new Organization by running:

```shell
forklift create organization organization-name --file ./resources/organization-config.yaml
```

Where organization-config.yaml is the Organization configuration which should correspond to the following Template.

```
vamp:
  persistence:
    database:
      sql:
        database: vamp-${namespace}
        url: jdbc:mysql://mysql.default.svc.cluster.local:3306/vamp-${namespace}?useSSL=false
        database-server-url: jdbc:mysql://mysql.default.svc.cluster.local:3306?useSSL=false
        user: root
        table: ${namespace}
        password: secret
      type: mysql
    key-value-store:
      vault:
        url: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_URL}
        token: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_TOKEN}
      base-path: /secret/vamp/${namespace}
      type: vault
    transformers:
      classes: []
  model:
    resolvers:
      namespace:
      - io.vamp.ee.model.NamespaceValueResolver
  security:
    lookup-hash-salt: b9a277bb-59a5-43d1-9c27-8a72e7e27685
    lookup-hash-algorithm: SHA-1
    session-id-length: 24
    password-hash-algorithm: SHA-512
    password-hash-salt: d4f22852-e281-428f-8968-1265b1c5a1b0
    token-value-length: 24
  pulse:
    elasticsearch:
      index:
        name: vamp-pulse-${namespace}
      url: http://elasticsearch.default.svc.cluster.local:9200
    type: elasticsearch
  metadata:
    namespace:
      title: organization
```

The above configuration can also be provided in JSON format.
Once created, you can list Organization by running

```shell
forklift list organizations
```

update them with

```shell
forklift update organization organization-name --file ./resources/organization-config.yaml
```

and show current configuration with

```shell
forklift show organization organization-name
```

delete them with

```shell
forklift delete organization organization-name
```

### Users

Through Forklift it is also possible to create Users for each Organization.
Users can be created interactively by running the following command whcih specifies the user name, role and organization of belonging.

```shell
forklift create user user-name --role admin --organization organization-name
```

Upon running the command you will be asked to input a new password twice, taking care to use at least six characters, before the user will be created.
It is also possible to create users not interactively by running:

```shell
forklift add user --organization organization-name --file ./resources/user-configuration.json
```

Where user-configuration.json is a file specifying the user configuration and should look like this:

```
{
  "name": "user-name",
  "password":"user-password",
  "kind":"users",
  "roles":["user-role"]
}
```

Mind the fact that add will insert the user or replace it if it already exists.
Once created, you can users Users by running


```shell
forklift update user user-name --role role-name organization organization-name
```

which will require you to specify the password interactively just like with the create.

You can then delete users with

```shell
forklift delete user user-name --organization organization-name
```

list them with

```shell
forklift list users --organization organization-name
```

and show a specific user with

```shell
forklift show user user-name --organization organization-name
```

### Environments

Environments can be created with Forklift by running:

```shell
forklift create environment environment-name --organization organization-name --file ./resources/environment-configuration.yaml --artifacts ./resources/artifacts
```

Where enviroment-configuration.yaml (or json) follows the template below:

```
vamp:
  persistence:
    database:
      sql:
        database: vamp-${parent}
        url: jdbc:mysql://mysql.default.svc.cluster.local:3306/vamp-${parent}?useSSL=false
        database-server-url: jdbc:mysql://mysql.default.svc.cluster.local:3306?useSSL=false
        user: root
        table: ${namespace}
        password: secret
      type: mysql
    key-value-store:
      vault:
        url: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_URL}
        token: ${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_TOKEN}
      base-path: /secret/vamp/${namespace}
      type: vault
    transformers:
      classes: []
  container-driver:
    type: kubernetes
    kubernetes:
      url: https://kubernetes
      vamp-gateway-agent-id: vamp-gateway-agent
      tls-check: false
  lifter:
    artifacts:
    - /usr/local/vamp/artifacts/breeds/quantification.yml
    - /usr/local/vamp/artifacts/workflows/quantification.yml
    - /usr/local/vamp/artifacts/breeds/vamp-workflow-javascript.yml
  gateway-driver:
    marshallers:
    - type: haproxy
      name: '1.8'
      template:
        resource: /io/vamp/gateway_driver/haproxy/template.twig
  model:
    resolvers:
      deployment:
      - io.vamp.ee.model.DisabledConfigurationValueResolver
      namespace:
      - io.vamp.ee.model.NamespaceValueResolver
      workflow:
      - io.vamp.ee.model.WorkflowValueResolver
      - io.vamp.pulse.ElasticsearchValueResolver
  workflow-driver:
    workflow:
      vamp-key-value-store-type: vault
      deployables:
      - type: application/javascript
        breed: vamp-workflow-javascript
      scale:
        cpu: 0.1
        instances: 1
        memory: 128MB
      vamp-key-value-store-connection: ${env://VAMP_WORKFLOW_DRIVER_WORKFLOW_VAMP_KEY_VALUE_STORE_CONNECTION}
      vamp-workflow-execution-period: 60
      vamp-key-value-store-token: ${env://VAMP_WORKFLOW_DRIVER_WORKFLOW_VAMP_KEY_VALUE_STORE_TOKEN}
      vamp-workflow-execution-timeout: 60
      vamp-elasticsearch-url: http://elasticsearch.default.svc.cluster.local:9200
      vamp-key-value-store-path: /secret/vamp/${namespace}/workflows/${workflow}
      vamp-url: http://vamp.default.svc.cluster.local:8080
    type: kubernetes
  pulse:
    elasticsearch:
      index:
        name: vamp-pulse-${namespace}
      url: http://elasticsearch.default.svc.cluster.local:9200
    type: elasticsearch
  operation:
    synchronization:
      period: 3 seconds
      check:
        health-checks: true
        deployable: true
        instances: true
        ports: true
        cpu: true
        environment-variables: true
        memory: true
      deployment:
        refetch-breed-on-update: true
    deployment:
      scale:
        instances: 1
        memory: 256MB
        cpu: 0.2
      arguments: []
    gateway:
      virtual-hosts:
        enabled: false
      selector: namespace(${namespace})
  metadata:
    namespace:
      title: environment
```

The --artifacts flag, on the other hand, provides a path where the specifications for workflows and breeds is provided in yaml format as shown below:

```
name: quantification
kind: workflows
breed: quantification
schedule: daemon
scale:
  cpu: 0.1
  memory: 256MB
  instances: 1
dialects:
  kubernetes:
    imagePullSecrets:
      - name: regsecret
```      

Once created, you can list Environment by running

```shell
forklift list environments --organization organization-name
```

update them with

```shell
forklift update environment environment-name --organization organization-name --file ./resources/environment-configuration.yaml --artifacts ./resources/artifacts
```

show current configuration with

```shell
forklift show environment environment-name --organization organization-name
```

and delete them with

```shell
forklift delete environment environment-name --organization organization-name
```

### Artifacts

Artifacts are breeds and workflows belonging to an environment.

Artifacts can be created or replaced with the following command:

```shell
forklift add artifact artifact-name --organization organization-name --environment environment-name --file ./resources/artifact.yaml
```

where artifact.yaml contains the artifact specification in this form:

```
name: test
kind: breeds
deployable:
  definition: magneticio/vamp-ee-workflows:1.0.4-quantification
ports:
  webport: 8080/http
environment_variables:
  VAMP_URL                            : ${config://vamp.workflow-driver.workflow.vamp-url}
  VAMP_API_TOKEN                      : ${vamp://token}
  VAMP_NAMESPACE                      : ${config://vamp.namespace}
  VAMP_WORKFLOW_EXECUTION_TIMEOUT     : ${config://vamp.workflow-driver.workflow.vamp-workflow-execution-timeout}
  VAMP_KEY_VALUE_STORE_CONNECTION     : ${config://vamp.workflow-driver.workflow.vamp-key-value-store-connection}
  VAMP_KEY_VALUE_STORE_TOKEN          : ${config://vamp.workflow-driver.workflow.vamp-key-value-store-token}
  VAMP_KEY_VALUE_STORE_PATH           : ${config://vamp.workflow-driver.workflow.vamp-key-value-store-path}
  VAMP_WORKFLOW_EXECUTION_PERIOD      : ${config://vamp.workflow-driver.workflow.vamp-workflow-execution-period}
  VAMP_KEY_VALUE_STORE_TYPE           : ${config://vamp.workflow-driver.workflow.vamp-key-value-store-type}
  VAMP_PULSE_ELASTICSEARCH_URL        : ${config://vamp.pulse.elasticsearch.url}
  VAMP_HEALTH                         : true
  VAMP_ELASTICSEARCH_HEALTH_INDEX     : ${es://health}
  VAMP_HEALTH_TIME_WINDOW             : 500
  VAMP_METRICS                        : true
  VAMP_ELASTICSEARCH_METRICS_INDEX    : ${es://metrics}
  VAMP_METRICS_TIME_WINDOW            : 500
  VAMP_CAPACITY                       : true
  VAMP_ELASTICSEARCH_CAPACITY_INDEX   : ${es://capacity}
  VAMP_ALLOCATION                     : true
  VAMP_ELASTICSEARCH_ALLOCATION_INDEX : ${es://allocation}
  VAMP_GATEWAY_DRIVER_ELASTICSEARCH_METRICS_TYPE : log
  VAMP_GATEWAY_DRIVER_ELASTICSEARCH_METRICS_INDEX: vamp-vga-${config://vamp.namespace}-*
```

Just like other resources, artifacts can be listed with

```shell
forklift list artifacts --kind artifact-kind --organization organization-name --environment environment-name
```

Where kind is the kind of the artifact (breeds or workflows).
Artifacts can also be deleted  with

```shell
forklift delete artifact artifact-name --kind artifact-kind --organization organization-name --environment environment-name
```

and shown with

```shell
forklift show artifact artifact-name --kind artifact-kind --organization organization-name --environment environment-name
```

### Release policy

Release policies can be created with the following command:

```shell
forklift add releasepolicy name --organization org --environment env --file ./releasepolicydefinition.json -i json
```

Example release policy:

```json
{
  "maxStartRetries": 10,
  "steps": [
    {
      "duration": "5m",
      "source": { "weight": 100 },
      "target": {
        "weight": 0,
        "condition": "user-agent = iPhone",
        "conditionStrength": 100
      }
    },
    {
      "duration": "5m",
      "source": { "weight": 50 },
      "target": { "weight": 50 }
    },
    {
      "source": { "weight": 0 },
      "target": { "weight": 100 }
    }
  ]
}
```

Release policies can also be deleted with

```shell
forklift delete releasepolicy name --organization org --environment env
```

### Release plan

Release plans can be created with the following command:

```shell
forklift add releaseplan name --file ./releaseplandefinition.json -i json
```

Release plan can also be deleted with

```shell
forklift delete releaseplan name
```