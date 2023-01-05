# NetOps


The netop pods enforce the QoS Profile for a slice. It uses Linux TC (Traffic Control) for the slice traffic classification.

## Getting Started

Please refer to our documentation on:
- [Installing KubeSlice on cloud clusters](https://kubeslice.io/documentation/open-source/0.5.0/getting-started-with-cloud-clusters/installing-kubeslice/installing-the-kubeslice-controller).
- [Installing KubeSlice on kind clusters](https://kubeslice.io/documentation/open-source/0.5.0/tutorials/kind-install-kubeslice-controller).

### Prerequisites

Before you begin, make sure the following prerequisites are met:
* Docker is installed and running on your local machine.
* A running [`kind`](https://kind.sigs.k8s.io/) cluster.
* [`kubectl`](https://kubernetes.io/docs/tasks/tools/) is installed and configured.
* You have prepared the environment to install [`kubeslice-controller`](https://github.com/kubeslice/kubeslice-controller) on the controller cluster
 and [`worker-operator`](https://github.com/kubeslice/worker-operator) on the worker cluster. For more information, see [Prerequisites](https://kubeslice.io/documentation/open-source/0.5.0/getting-started-with-cloud-clusters/prerequisites/).
 
### Local Build and Update 

#### Latest Docker Hub Image

```console
docker pull aveshasystems/netops:latest
```

### Setting up Your Helm Repo

If you have not added avesha helm repo yet, add it.

```console
helm repo add avesha https://kubeslice.github.io/charts/
```

Upgrade the avesha helm repo.

```console
helm repo update
```

### Build docker images

1. Clone the latest version of netops from  the `master` branch. 

```bash
git clone https://github.com/kubeslice/netops.git
cd netops
```

2. Adjust the `VERSION` variable in the Makefile to change the docker tag to be built.
Image is set as `docker.io/aveshasystems/netops:$(VERSION)` in the Makefile. Change this if required.

```bash
make docker-build
```

### Running Locally on Kind Clusters

1. You can load the netops image on your kind cluster using the following command:

```bash
kind load docker-image my-custom-image:unique-tag --name clustername
```

Example

```console
kind load docker-image aveshasystems/netops:1.2.1 --name kind
```

2. Check the loaded image in the cluster. Modify node name if required.

```console
docker exec -it <node-name> crictl images
```

Example

```console
docker exec -it kind-control-plane crictl images
```


### Deploy in a Cluster

Update the chart values file called `yourvaluesfile.yaml` that you have previously created.
Refer to [values.yaml](https://github.com/kubeslice/charts/blob/master/charts/kubeslice-worker/values.yaml) to create `yourvaluesfiel.yaml` and update the netop image subsection to use the local image.

From the sample: 

```
netop:
  image: docker.io/aveshasystems/netops
  tag: 0.1.0
```

Change it to:

```
netop:
  image: <my-custom-image>
  tag: <unique-tag>
```

Deploy the updated chart.

```console
make chart-deploy VALUESFILE=yourvaluesfile.yaml
```

### Verify the Installation

Verify the installation of NetOps by checking the status of pods belonging to the `kubeslice-system` namespace.

```bash
kubectl get pods -n kubeslice-system | grep netop
```
Example output

```
avesha-netop-pnbbr                         1/1     Running   0          4d23h
```

## License
Apache 2.0 License.
