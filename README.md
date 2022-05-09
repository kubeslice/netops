# netops

The NetOp Pods enforce the QoS Profile for a Slice. It uses Linux TC (Traffic Control) for Slice traffic classification.

## Getting Started

It is strongly recommended to use a released version.

### Prerequisites

* Docker installed and running in your local machine
* A running [`kind`](https://kind.sigs.k8s.io/) or [`Docker Desktop Kubernetes`](https://docs.docker.com/desktop/kubernetes/)
  cluster 
* [`kubectl`](https://kubernetes.io/docs/tasks/tools/) installed and configured

### Build and push docker images

```bash
git clone https://github.com/kubeslice/netops.git
cd netops
make docker-build
make docker-push
```

### Usages
You can view the NetOp Pods by using the command below:

```bash
kubectl get pods -n avesha-system | grep netops
```

## License
This project is released under the Apache 2.0 License.