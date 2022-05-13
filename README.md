# Netops


The netop pods enforce the QoS Profile for a Slice. It uses Linux TC (Traffic Control) for Slice traffic classification.

## Getting Started

[TBD add link to getting started] 
It is strongly recommended to use a released version.

### Prerequisites

* Docker installed and running in your local machine
* A running [`kind`](https://kind.sigs.k8s.io/)
* [`kubectl`](https://kubernetes.io/docs/tasks/tools/) installed and configured
* Follow the getting started from above, to install [`kubeslice-controller`](https://github.com/kubeslice/kubeslice-controller) and [`worker-operator`](https://github.com/kubeslice/worker-operator)

### Local build and update 

#### Latest docker image
[TBD link to docker hub]

### Setting up your helm repo

If you have not added avesha helm repo yet, add it

```console
helm repo add avesha https://kubeslice.github.io/charts/
```

upgrade the avesha helm repo

```console
helm repo update
```

### Build docker images

```bash
git clone https://github.com/kubeslice/netops.git
cd netops
make docker-build
```

### Running locally on Kind

Load the docker image into kind cluster

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
