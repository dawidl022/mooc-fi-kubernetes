# CNCF Cloud Native Landscape technologies used

## App Definition and Development

### Database

- MariaDB used outside the course
- MongoDB used outside the course
- MySQL used outside the course
- PostgreSQL used since part 2 as the database for our applications
- Redis used in an example in part 2 to learn about StatefulSets

### Streaming & Messaging

- NATS used in part 4 to send messages between applications
- RabbitMQ used outside the course

### Application definition and image build

- HELM used since part 2 to install k8s packages (called charts)
- Bitnami used outside the course
- Docker Compose used in DevOpsWithDocker
- Gradle used outside the course
- OpenAPI used in part 5 to define the schema of a CustomResourceDefinition
- Skaffold used outside the course

### Continuous Integration & Delivery

- Argo rollouts used in part 4 to create a canary release
- Flux used in part 4 to implement GitOps
- GitHub actions used since part 1 to deploy Docker images to a container registry
- GitLab used outside the course

## Orchestration and management

### Scheduling & Orchestration

- Kubernetes used throughout the course

### Coordination & service discovery

- etcd indirectly used as a dependency of Kubernetes to persist the cluster state

### Remote Procedure Call

- hRPC used outside the course

### Service Proxy

- Contour used in part 5 as the ingress controller for Knative
- Nginx used since part 1 to serve static files over HTTP
- Traefik used since part 1 as the default ingress controller in k3s

### Service Mesh

- Linkerd used in part 5 as the service mesh of choice

## Runtime

### Container Runtime

- containerd used since part 1 as the container runtime in Kubernetes

### Cloud Native Network

- cilium used outside the course

## Provisioning

### Automation & Configuration

- Terraform used outside the course

### Container Registry

- Google Container Registry used in part 3 to push application docker images to and pull from for GKE

### Security & Compliance

- Tetragon used outside the course as a component of cilium

## Platform

### Certified Kubernetes - Distribution

- k3s used since part 1 as the underlying distribution k3d uses
- Red Hat okd used outside the course

### Certified Kubernetes - Hosted

- Google Kubernetes Engine used in part 3 as the hosted k8s platform

### PaaS/Container Service

- Akka used outside the course
- heroku used outside the course

## Observability and Analysis

### Monitoring

- Prometheus used in part 2 to monitor the cluster
- Hubble used outside the project as a component of cilium
- Grafana used in part 2 to view monitoring data

### Logging

- Grafana Loki used in part 2 to view and query logs

## Serverless

### Installable platform

- Knative used in part 5 as the serverless platform running on top of Kubernetes
