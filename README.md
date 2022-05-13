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

#### Latest docker hub image

```console
docker pull aveshasystems/netops:latest
```

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

1. Clone the latest version of netops from  the `master` branch. 

```bash
git clone https://github.com/kubeslice/netops.git
cd netops
```

2. Adjust `VERSION` variable in the Makefile to change the docker tag to be built.
Image is set as `docker.io/aveshasystems/netops:$(VERSION)` in the Makefile. Change this if required

```bash
make docker-build
```

### Running locally on Kind

Load the docker image into kind cluster

```bash
kind load docker-image my-custom-image:unique-tag --name clustername
```

### Deploy in a cluster

Update chart values file `yourvaluesfile.yaml` that you have previously created.
Refer to [values.yaml](https://github.com/kubeslice/charts/blob/master/kubeslice-worker/values.yaml) to create `yourvaluesfiel.yaml` and update the netop image subsection to use the local image.

From the sample , 

```
netop:
  image: docker.io/aveshasystems/netops
  tag: 0.1.0
```

Change it to ,

```
netop:
  image: <my-custom-image>
  tag: <unique-tag>
```

Deploy the updated chart

```console
make chart-deploy VALUESFILE=yourvaluesfile.yaml
```

### Verify the NetOp Pods are running:

```bash
kubectl get pods -n kubeslice-system | grep netop
```

## License
Apache 2.0 License.
