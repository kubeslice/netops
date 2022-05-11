# netops

The netop pods enforce the QoS Profile for a Slice. It uses Linux TC (Traffic Control) for Slice traffic classification.

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
```

### Running locally on Kind

You can run the operator on your Kind cluster with the below command

```bash
kind load docker-image my-custom-image:unique-tag --name clustername
```

### Verification
You can view the NetOp Pods by using the command below:

```bash
kubectl get pods -n kubeslice-system | grep netop
```

## License
Apache 2.0 License.
