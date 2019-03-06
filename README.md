# Vamp Forklift command line client

Vamp Forklift is a command line client written in golang and allows to easily set up Organizations and Environments in Vamp.

Forklift requires running and reachable instances of MySql and Vault tied to a Vamp installation.
Please check How to Setup Vamp at the following link https://vamp.io/documentation/installation/v1.0.0/overview/

## development

if you have golang installed, it is recommended to git clone Forklift to $GOPATH/src/github.com/magneticio/forklift
This is a requirement for docker builder to work.

It is also recommended to read and follow golang setup for a development environment setup: https://golang.org/doc/install

## build

If you get errors about missing libraries while building, run:
```shell
go get
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

## installation
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

Easy install for MacOS:
```shell
base=base=https://github.com/magneticio/forklift/releases/download/0.1.0 &&
  curl -L $base/forklift-$(uname -s)-$(uname -m) >/usr/local/bin/forklift &&
  chmod +x /usr/local/bin/forklift
```
TODO: add installation for other platforms

For general users it is recommended to download the binary for your platform.
Latest release can be found here:
https://github.com/magneticio/forklift/releases/latest

Now make sure to have a ".forklift.yaml" configuration file in your home, like the one shown below, but with the correct parameters to connect to the database and the key-value store.

```
namespace: vampio
forklift:
  persistence:
    database:
      sql:
        database: vamp-${namespace}
        url: mysql://mysql.default.svc.cluster.local:3306/vamp-${namespace}?useSSL=false
        database-server-url: mysql://mysql.default.svc.cluster.local:3306?useSSL=false
        user: root
        table: ${namespace}
        password: secret
      type: mysql
    key-value-store:
      vault:
        url: http://vault.default.svc.cluster.local:8200
        token: vamp
      base-path: /secret/vamp/${namespace}
      type: vault
  metadata:
    namespace:
      title: organization

```

### Verifying installation

To verify the installation you can run the following command, whcih will return the version of the client's and vamp's versions.

```shell
forklift version
```

It is possible to get all commands and flags by running help:
```shell
forklift help
```

### Organizations

Forklift allows for the creation of a new Organization by running:

```shell
forklift create organization organization-name --configuration ./resources/organization-config.yaml
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

The above configuration can also be provided in Json format.
Once created, you can list Organization by running

```shell
forklift list organization
```

update them with

```shell
forklift update organization organization-name --configuration ./resources/organization-config.yaml
```

and delete them with

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
forklift add user --organization organization-name --configuration ./user-configuration.json
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

Once created, you can delete Users by running

```shell
forklift delete user user-name --organization organization-name
```

### Environments

Environments can be created with Forklift by running:

```shell
forklift create environment environment-name --organization organization-name --configuration ./resources/environment-configuration.yaml --artifacts ./resources/artifacts
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
forklift update environment environment-name --organization organization-name --configuration ./resources/environment-configuration.yaml --artifacts ./resources/artifacts
```

and delete them with

```shell
forklift delete environment environment-name --organization organization-name
```
