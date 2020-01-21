---
title: Configuration Reference for devspace.yaml
sidebar_label: Config Reference
---

## `version`
```yaml
version: v1beta6                   # string   | Version of the config
```

<details>
<summary>
### List of supported versions
</summary>
- v1beta6   ***latest***
- v1beta5
- v1beta4
- v1beta3
- v1beta2
- v1beta1
- v1alpha4
- v1alpha3
- v1alpha2 
- v1alpha1
</details>

---
## `images`
```yaml
images:                             # map[string]struct | Images to be built and pushed
  image1:                           # string   | Name of the image
    image: dscr.io/username/image   # string   | Image repository and name 
    tag: v0.0.1                     # string   | Image tag
    dockerfile: ./Dockerfile        # string   | Relative path to the Dockerfile used for building (Default: ./Dockerfile)
    context: ./                     # string   | Relative path to the context used for building (Default: ./)
    entrypoint: []                  # string[] | Override ENTRYPOINT defined in Dockerfile
    cmd: []                         # string[] | Override CMD defined in Dockerfile
    createPullSecret: true          # bool     | Create a pull secret containing your Docker credentials (Default: false)
    build: ...                      # struct   | Build options for this image
  image2: ...
```
[Learn more about building images with DevSpace.](../../cli/image-building/workflow-basics)

### `images[*].build`
```yaml
build:                              # struct   | Build configuration for an image
  docker: ...                       # struct   | Build image with docker and set options for docker
  kaniko: ...                       # struct   | Build image with kaniko and set options for kaniko
  custom: ...                       # struct   | Build image using a custom build script
  disabled: false                   # bool     | Disable image building (Default: false)
```
Notice:
- Setting `docker`, `kaniko` or `custom` will define the build tool for this image.
- You **cannot** use `docker`, `kaniko` and `custom` in combination. 
- If neither `docker`, `kaniko` nor `custom` is specified, `docker` will be used by default.
- By default `docker` will use `kaniko` as fallback when DevSpace is unable to reach the Docker host.

### `images[*].build.docker`
```yaml
docker:                             # struct   | Options for building images with Docker
  preferMinikube: true              # bool     | If available, use minikube's in-built docker daemon instaed of local docker daemon (default: true)
  skipPush: false                   # bool     | Skip pushing image to registry, recommended for minikube (Default: false)
  disableFallback: false            # bool     | Disable using kaniko as fallback when Docker is not installed (Default: false)
  useBuildKit: true                 # bool     | Enable BuildKit for Docker build process (default: false)
  options: ...                      # struct   | Set build general build options
```

### `images[*].build.kaniko`
```yaml
kaniko:                             # struct   | Options for building images with kaniko
  cache: true                       # bool     | Use caching for kaniko build process
  snapshotMode: "time"              # string   | Type of snapshotMode for kaniko build process (compresses layers)
  flags: []                         # string[] | Array of flags for kaniko build command
  namespace: ""                     # string   | Kubernetes namespace to run kaniko build pod in (Default: "" = deployment namespace)
  insecure: false                   # bool     | Allow working with an insecure registry by not validating the SSL certificate (Default: false)
  pullSecret: ""                    # string   | Mount this Kubernetes secret instead of creating one to authenticate to the registry (default: "")
  options: ...                      # struct   | Set build general build options
```

### `images[*].build.custom`
```yaml
custom:                             # struct   | Options for building images with a custom build script
  command: "./scripts/builder"      # string   | Command to be executed for building (e.g. path to build script or executable)
  args: []                          # string[] | Array of arguments for the custom build command
  imageFlag: string                 # string   | Name of the flag that DevSpace uses to pass the image name + tag to the build script
  onChange: []                      # string[] | Array of paths (glob format) to check for file changes to see if image needs to be rebuild
```

### `images[*].build.*.options`
```yaml
options:                            # struct   | Options for building images
  target: ""                        # string   | Target used for multi-stage builds
  network: ""                       # string   | Network mode used for building the image
  buildArgs: {}                     # map[string]string | Key-value map specifying build arguments that will be passed to the build tool (e.g. docker)
```


---
## `deployments`
```yaml
deployments:                        # struct[] | Array of deployments
- name: my-deployment               # string   | Name of the deployment
  namespace: ""                     # string   | Namespace to deploy to (Default: "" = namespace of the active namespace/Space)
  helm: ...                         # struct   | Use Helm as deployment tool and set options for Helm
  kubectl: ...                      # struct   | Use "kubectl apply" as deployment tool and set options for kubectl
```
Notice:
- Setting `component`, `helm` or `kubectl` will define the type of deployment and the deployment tool to be used.
- You **cannot** use `component`, `helm` and `kubectl` in combination.

### `deployments[*].component`
Instead of using `component` directly, the default way of deploying components is to use `helm` as deployment method with `componentChart: true` and with the component config options under `helm.values`. Both methods do exactly the same.

```yaml
helm:
  componentChart: true              # bool     | Use Component chart
  values:                           # struct   | Options for deploying a component = Deployment/StatefulSet
    initContainers: ...             # struct[] | Init Containers of this Deployment/StatefulSet
    containers: ...                 # struct[] | Containers of this Deployment/StatefulSet
    labels: {}                      # map[string]string | Map of Kubernetes labels for labeling the pods of this component
    annotations: {}                 # map[string]string | Map of Kubernetes annotations for annotating the pods of this component
    volumes: ...                    # struct   | Component volumes
    service: ...                    # struct   | Component service
    serviceName: my-service         # string   | Service name for headless service (for StatefulSets)
    ingress: ...                    # struct   | Component ingress
    replicas: 1                     # int      | Number of replicas (Default: 1)
    autoScaling: ...                # struct   | AutoScaling configuration
    rollingUpdate: ...              # struct   | RollingUpdate configuration
    pullSecrets: ...                # string[] | Array of PullSecret names
    podManagementPolicy: OrderedReady # enum     | "OrderedReady" or "Parallel" (for StatefulSets)
    options: ...                    # struct   | Options for deploying this component with helm
```

[Learn more about configuring component deployments.](../../cli/deployment/components/what-are-components)

### `deployments[*].component.containers/initContainers`
```yaml
containers:                         # struct   | Options for deploying a DevSpace component
- name: my-container                # string   | Container name (optional)
  image: dscr.io/username/image     # string   | Image name (optionally with registry URL)
  command:                          # string[] | ENTRYPOINT override
  - sleep
  args:                             # string[] | ARGS override
  - 99999
  stdin: true                       # bool     | Enable stdin (Default: false)
  tty: true                         # bool     | Enable tty (Default: false)
  env:                              # map[interface]interface | Kubernetes env definition for containers
  - name: MY_ENV_VAR
    value: "my-value"
  volumeMounts: ...                 # struct   | VolumeMount Configuration
  resources: ...                    # struct   | Kubernestes resource limits and requests
  livenessProbe: ...                # struct   | Kubernestes livenessProbe
  readinessProbe: ...               # struct   | Kubernestes readinessProbe
```

### `deployments[*].component.containers[*].volumeMounts`
```yaml
volumeMounts: 
  containerPath: /my/path           # string   | Mount path within the container
  volume:                           # struct   | Volume to mount
    name: my-volume                 # string   | Name of the volume to be mounted
    subPath: /in/my/volume          # string   | Path inside to volume to be mounted to the containerPath
    readOnly: false                 # bool     | Mount volume as read-only (Default: false)
```

### `deployments[*].component.volumes`
```yaml
volumes: 	                          # struct   | Array of volumes to be created
- name: my-volume                   # string   | Volume name
  labels: {}                        # map[string]string | Map of Kubernetes labels for labeling the pods of this component
  annotations: {}                   # map[string]string | Map of Kubernetes annotations for annotating the pods of this component
  size: 10Gi                        # string   | Size of the volume in Gi (Gigabytes)
  configMap: ...                    # struct   | Kubernetes ConfigMapVolumeSource
  secret: ...                       # struct   | Kubernetes SecretVolumeSource
```

### `deployments[*].component.service`
```yaml
service: 	                          # struct   | Component service configuration
  name: my-service                  # string   | Name of the service
  labels: {}                        # map[string]string | Map of Kubernetes labels for labeling the pods of this component
  annotations: {}                   # map[string]string | Map of Kubernetes annotations for annotating the pods of this component
  type: NodePort                    # string   | Type of the service (default: ClusterIP)
  ports:                            # array    | Array of service ports
  - port: 80                        # int      | Port exposed by the service
    containerPort: 3000             # int      | Port of the container/pod to redirect traffic to
    protocol: tcp                   # string   | Traffic protocol (tcp, udp)
  externalIPs:                      # array    | Array of externalIPs for the service (discouraged)
  - 123.45.67.890                   # string   | ExternalIP to expose the service on
```

### `deployments[*].component.ingress`
```yaml
ingress: 	                          # struct   | Component service configuration
  name: my-service                  # string   | Name of the ingress (optional)
  labels: {}                        # map[string]string | Map of Kubernetes labels for labeling the ingress
  annotations: {}                   # map[string]string | Map of Kubernetes annotations for annotating the ingress
  tls: false                        # string   | Enable/Disable tls for all hosts OR provide a secret name (default: false)
  rules:                            # struct   | Array of ingress rules
  - host: my-domain.tld             # string   | Hostname for the ingress rule
    tls: false                      # string   | Enable/Disable tls for this host OR provide a secret name (default: false)
    path: /                         # string   | Path for routing the traffic (default: /)
    servicePort: 8080               # int      | Service port to forward traffic to (default: first port of component.service)
    serviceName: ""                 # string   | Name of the Kubernetes service to route the traffic to (default: service defined in component.service)
```

### `deployments[*].component.autoScaling`
```yaml
autoScaling: 	                      # struct   | Auto-Scaling configuration
  horizontal:                       # struct   | Configuration for horizontal auto-scaling
    maxReplicas: 5                  # int      | Max replicas to deploy
    averageCPU: 800m                # string   | Target value for CPU usage
    averageRelativeCPU: 50          # int      | Target percentage of CPU relative to CPU limit configured for the container (e.g. 50%)
    averageMemory: 1Gi              # string   | Target value for memory (RAM) usage
    averageRelativeMemory: 70       # int      | Target percentage of memory relative to memory limit configured for the container
```

### `deployments[*].component.rollingUpdate`
```yaml
rollingUpdate: 	                    # struct   | Rolling-Update configuration
  enabled: false                    # bool     | Enable/Disable rolling update (Default: disabled)
  maxSurge: "25%"                   # string   | Max number of pods to be created above the pod replica limit
  maxUnavailable: "50%"             # string   | Max number of pods unavailable during update process
  partition: 1                      # int      | For partitioned updates of StatefulSets
```

### `deployments[*].helm`
```yaml
helm:                               # struct   | Options for deploying with Helm
  chart: ...                        # struct   | Relative path 
  componentChart: false             # bool     | Use the DevSpace component chart instead of a custom `chart` = deployment is a component (Default: false)
  values: {}                        # struct   | Any object with Helm values to override values.yaml during deployment
  valuesFiles:                      # string[] | Array of paths to values files
  - ./chart/my-values.yaml          # string   | Path to a file to override values.yaml with
  replaceImageTags: true            # bool     | Enable automated tag replacement (Default: true)
  wait: false                       # bool     | Wait for pods to start after deployment (Default: false)
  timeout: 180                      # int      | Timeout to wait for pods to start after deployment (Default: 180)
  force: false                      # bool     | Force deleting and re-creating Kubernetes resources during deployment (Default: false)
  atomic: false                     # bool     | Rollback deployment if it fails (Default: false)
  cleanupOnFail: false              # bool     | Delete resources if rollback fails (Default: false)
  recreate: false                   # bool     | Recreate pods for applicable resources, e.g. deployments (Default: false)
  disableHooks: false               # bool     | Disable hooks (Default: false)
  driver: secrets|configmaps|memory # enum     | Driver used by Helm v3 to store release configuration (Default: secrets)
  v2: false                         # bool     | Use legacy Helm v2 (Default: false)
  tillerNamespace: ""               # string   | Kubernetes namespace to run Tiller in when using Helm v2 (Default: "" = same a deployment namespace)
```
[Learn more about configuring deployments with Helm.](../../cli/deployment/helm-charts/what-are-helm-charts)

### `deployments[*].helm.chart`
```yaml
chart:                              # struct   | Chart to deploy
  name: my-chart                    # string   | Path to local chart on filesystem OR chart name for remote chart in helm chart repository
  version: v1.0.1                   # string   | Chart version
  repo: "https://my-repo.tld/"      # string   | Helm chart repository
  username: "my-username"           # string   | Username for Helm chart repository
  password: "my-password"           # string   | Password for Helm chart repository
```

### `deployments[*].kubectl`
```yaml
kubectl:                            # struct   | Options for deploying with "kubectl apply"
  manifests: []                     # string[] | Array containing glob patterns for the Kubernetes manifests to deploy using "kubectl apply" (e.g. kube or manifests/service.yaml)
  kustomize: false                  # bool     | Use kustomize when deploying manifests via "kubectl apply" (Default: false)
  replaceImageTags: true            # bool     | Enable automated tag replacement (Default: true)
  flags: []                         # string[] | Array of flags for the "kubectl apply" command
  cmdPath: ""                       # string   | Path to the kubectl binary (Default: "" = detect automatically)
```
[Learn more about configuring deployments with Kubectl.](../../cli/deployment/kubernetes-manifests/what-are-manifests)


---
## `dev`
```yaml
dev:                                # struct   | Options for "devspace dev"
  ports: []                         # struct[] | Array of port-forwarding settings for selected pods
  open: []                          # struct[] | Array of auto-open settings
  sync: []                          # struct[] | Array of file sync settings for selected pods
  logs: ...                         # struct   | Options for configuring multi-container log streaming
  autoReload: ...                   # struct   | Options for auto-reloading (i.e. re-deploying deployments and re-building images)
  interactive: ...                  # struct   | Options for configuring the interactive mode
```
[Learn more about development with DevSpace.](../../cli/development/workflow-basics)

### `dev.ports`
```yaml
ports:                              # struct[] | Array of port forwarding settings for selected pods
- imageName: someImage              # string   | Name of an image defined in `images` to select pods with
  labelSelector: ...                # struct   | Key Value map of labels and values to select pods with
  namespace: ""                     # string   | Kubernetes namespace to select pods in
  forward:                          # struct[] | Array of ports to be forwarded
  - port: 8080                      # int      | Forward this port on your local computer
    remotePort: 3000                # int      | Forward traffic to this port exposed by the pod/container selected
    bindAddress: ""                 # string   | Address used for binding / use 0.0.0.0 to bind on all interfaces (Default: "localhost" = 127.0.0.1)
```
[Learn more about port forwarding.](../../cli/development/configuration/port-forwarding)

### `dev.open`
```yaml
open:                               # struct[] | Array of auto-open settings
- url: "https://localhost:3000/"    # string   | URL to open after application has started
```
[Learn more about configuring auto-opening.](../../cli/development/configuration/auto-open)

### `dev.sync`
```yaml
sync:                               # struct[] | Array of file sync settings for selected pods
- imageName: someImage              # string   | Name of an image defined in `images` to select pods with
  labelSelector: ...                # struct   | Key Value map of labels and values to select pods with
  containerName: ""                 # string   | Container name to use after selecting a pod
  namespace: ""                     # string   | Kubernetes namespace to select pods in
  localSubPath: ./                  # string   | Relative path to a local folder that should be synchronized (Default: "./" = entire project)
  containerPath: /app               # string   | Path in the container that should be synchronized with localSubPath (Default is working directory of container ("."))
  excludePaths: []                  # string[] | Paths to exclude files/folders from sync in .gitignore syntax
  downloadExcludePaths: []          # string[] | Paths to exclude files/folders from download in .gitignore syntax
  uploadExcludePaths: []            # string[] | Paths to exclude files/folders from upload in .gitignore syntax
  downloadOnInitialSync: false      # bool     | Download files that exist inside the container but not on the local filesystem during initial sync (Default: false)
  waitInitialSync: false            # bool     | Wait until initial sync is completed before continuing (Default: false)
  bandwidthLimits:                  # struct   | Bandwidth limits for the synchronization algorithm
    download: 0                     # int64    | Max file download speed in kilobytes / second (e.g. 100 means 100 KB/s)
    upload: 0                       # int64    | Max file upload speed in kilobytes / second (e.g. 100 means 100 KB/s)
```
[Learn more about confguring the code synchronization.](../../cli/development/configuration/file-synchronization)

### `dev.logs`
```yaml
logs:                               # struct   | Options for multi-container log streaming in development mode
  disabled: false                   # bool     | Disable log streaming in development mode (Default: false)
  showLast: 200                     # int      | Number of last log lines to show before starting stream (Default: 50)
  images: []                        # string[] | Array of image names referencing images defined in `images` for selecting containers for log streaming
```
[Learn more about configuring multi-container log streaming.](../../cli/development/configuration/logs-streaming)

### `dev.autoReload`
```yaml
autoReload:                         # struct   | Options for auto-reloading (i.e. re-deploying deployments and re-building images)
  paths: []                         # string[] | Array containing glob patterns of files that are watched for auto-reloading (i.e. reload when a file matching any of the patterns changes)
  deployments: []                   # string[] | Array containing names of deployments to watch for auto-reloading (i.e. reload when kubectl manifests or files within the Helm chart change)
  images: []                        # string[] | Array containing names of images to watch for auto-reloading (i.e. reload when the Dockerfile changes)
```

### `dev.interactive`
```yaml
interactive:                        # struct   | Options for interactive mode
  defaultEnabled: false             # bool     | Start interactive mode instead of log streaming by default, even without -i / --interactive flag (Default: false)
  images:                           # struct[] | Array of image override configurations for interactive mode
  - name: default                   # string   | Name of the image to apply this override rule to (key in `images`)
    entrypoint: []                  # string[] | Array defining with the ENTRYPOINT that should be used instead of the ENTRYPOINT defined in the Dockerfile
    cmd: []                         # string[] | Array defining with the CMD that should be used instead of the CMD defined in the Dockerfile
  terminal:                         # struct   | Options for the terminal proxy
    imageName: someImage            # string   | Name of an image defined in `images` to select pods with
    labelSelector: ...              # struct   | Key Value map of labels and values to select pods with
    containerName: ""               # string   | Container name to use after selecting a pod
    namespace: ""                   # string   | Kubernetes namespace to select pods in
    command: []                     # string[] | Array defining the shell command to start the terminal with (Default: ["sh", "-c", "command -v bash >/dev/null 2>&1 && exec bash || exec sh"])
```
[Learn more about interactive mode.](../../cli/development/configuration/interactive-mode)


---
## `dependencies`
```yaml
dependencies:                       # struct[]  | Array of dependencies (other projects containing a devspace.yaml or devspace-configs.yaml) that need to be deployed before this project
- source:                           # struct    | Defines where to find the dependency (exactly one source is allowed)
    git: https://github.com/my-repo # string    | HTTP(S) URL of the git repository (recommended method for referencing dependencies, must have the format of the git remote repo as usually checked out via git clone)
    subPath: repo/sub/path          # string    | Path within the git repo where devspace.yaml can be found
    branch: master                  # string    | Git branch to checkout
    tag: v1.2.3                     # string    | Git tag to checkout
    revision: ac66e49               # string    | Git revision (commit has) to checkout
    path: ../../my-projects/repo    # string    | Path to a project on your local computer (not recommended, instead of using git-related options)
  profile: default                  # string    | Name of the profile used to deploy this dependency (when multiple prpfiles are defined in the devspace.yaml of the dependency)
  skipBuild: false                  # bool      | Do not build images of this dependency (= only start deployments)
  ignoreDependencies: false         # bool      | Do not build and deploy dependencies of this dependency
  namespace: ""                     # string    | Kubernetes namespace to deploy dependency to (Default: default namespace of current kube-context)
```
> You **cannot** use `source.git` and `source.path` in combination. You **must** exactly use one of the two.


---
## `hooks`
```yaml
hooks:                              # struct[]  | Array of hooks to be executed
- command: "./scripts/my-hook"      # string    | Command to be executed when this hook is triggered
  args: []                          # string[]  | Array of arguments for the command of this hook
  when:                             # struct    | Trigger for executing this hook 
    before:                         # struct    | Run hook before a certain execution step
      images: "all"                 # string    | Name of the image you want to run this hook before building OR "all" for running hook before building the first image
      deployments: "all"            # string    | Name of the deployment you want to run this hook before deploying OR "all" for running hook before deploying the first deployment
    after:                          # struct    | Run hook after a certain execution step
      images: "all"                 # string    | Name of the image you want to run this hook after building OR "all" for running hook after building the last image
      deployments: "all"            # string    | Name of the deployment you want to run this hook after deploying OR "all" for running hook after deploying the last deployment
```


---
## `commands`
```yaml
commands:                             # struct[]  | Array of custom commands
- name: "debug-backend"               # string    | Name of the command to run via `devspace run debug-backend`
  command: "devspace dev -i backend"  # string    | Command to be executed when running `devspace run [name]`
```


---
## `vars`
```yaml
vars:                               # struct[]  | Array of config variables
- name: CONFIG_VAR                  # string    | Name of the config variable
  question: "What is CONFIG_VAR?"   # string    | Question to ask the user if no value is found for variable
  options: []                       # string[]  | Options for picker (selector) to show to user (to choose a value for variable)
  password: false                   # bool      | Hide user input when providing value via command-line, i.e. replaces input with `*****` (Default: false)
  validationPattern: ""             # string    | Regexp to validate user input
  validationMessage: ""             # string    | Message to show to user for input validation
  default: ""                       # string    | Default value for variable
  source: "all"                     # enum      | Source for variable (all = default, env, input)
```


---
## `profiles`
```yaml
profiles:                           # struct[]  | Array of config profiles
- name: profile-name                # string    | Name of the profile
  patches:                          # struct[]  | Array of config patches
  - op: "replace"                   # enum      | Patch operation (replace, add, remove)
    path: "images.backend.cmd"      # string    | Jsonpath or xpath to config option that should be patched
    value: ""                       # arbitrary | Value to use for patch operation
    from: ""                        # string    | Jsonpath or xpath to config option which should be used as value for operation
  replace:                          # struct    | Array of replacements for entire config sections
    images: {}                      # struct    | Replacement for entire `images` section
    deployments: {}                 # struct    | Replacement for entire `deployments` section
    dev: {}                         # struct    | Replacement for entire `dev` section
    dependencies: {}                # struct    | Replacement for entire `dependencies` section
    hooks: {}                       # struct    | Replacement for entire `hooks` section
```
