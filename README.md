# DevOpsWithKubernetes exercise solutions

My solutions to the [DevOpsWithKubernetes](https://devopswithkubernetes.com/) course exercise.

## Main Projects

### [DummySite CustomResourceDefinition & controller](dummy-site/README.md)

Custom Kubernetes resource allowing to download an arbitrary webpage and host in on the Kubernetes cluster by specifying a simple declarative DummySite resource in yaml.

The controller, once it detects a change in the cluster state (an added DummySite) takes care of downloading the webpage and creating the necessary Kubernetes deployment, service and ingress. It also clears those resources up once it detects a DummySite is deleted from the cluster.

### Todo application

Todo application with a [Go](https://go.dev/) frontend and a [React](https://reactjs.org/) backend. Uses [NATS Messaging](https://nats.io/) to send notification to a broadcaster service (also written in Go), which then sends the notification to a chat in [Telegram](https://telegram.org/) using the [Telegram API](https://core.telegram.org/api).
