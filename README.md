# 🍕 provider-pizza
I was disappointed to learn [Crossplane](https://crossplane.io) doesn't support my favorite flavor of the cloud - pizza. So I had to build a Domino's Pizza provider to make the cloud-native community's universal cloud API complete.

## Usage
To create an order:
```
kubectl apply -f examples/order/order.yaml
```
Which should show you something like:
```
NAME            ORDER STATUS   PRICE   STORE ADDRESS                                 MANAGER   AGE
example-order   Bad            3.29    2320 B North 45th Street Seattle, WA 98103    Sarah     2m48s
```

## Limitations
Credit card payments aren't yet working. I'll be adding support soon. You'll store credit card details as a K8s secret and specify `Type: CreditCard` in the secret. See `examples/secret-example.yaml` for details.

## Developing

Run against a Kubernetes cluster:

```console
make run
```

Install `latest` into Kubernetes cluster where Crossplane is installed:

```console
make install
```

Install local build into [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
cluster where Crossplane is installed:

```console
make install-local
```

Build, push, and install:

```console
make all
```

Build image:

```console
make image
```

Push image:

```console
make push
```

Build binary:

```console
make build
```

## Credit
Having never written Go before, and being very unfamiliar with Kubernetes controllers, I had a ton of inspiration from https://github.com/rudoi/cruster-api. This provider is basically just a wrapper around https://github.com/rudoi/pizza-go.