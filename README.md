# kubectl-unsecret

`kubectl-unsecret` is a plugin for `kubectl` that allows you to decode Kubernetes secrets and display them in a readable format.

## Installation

### Using `go install`

You can install this plugin using `go install`:

```sh
go install github.com/azratul/kubectl-unsecret/cmd/kubectl-unsecret@latest
```

### Manual Installation

Download the pre-compiled binary from the Releases page and place it in your $PATH.

#### Linux

```sh
wget https://github.com/azratul/kubectl-unsecret/releases/download/v1.0.0/kubectl-unsecret-linux-amd64
chmod +x kubectl-unsecret-linux-amd64
sudo mv kubectl-unsecret-linux-amd64 /usr/local/bin/kubectl-unsecret
```

#### macOS

```sh
wget https://github.com/azratul/kubectl-unsecret/releases/download/v1.0.0/kubectl-unsecret-darwin-amd64
chmod +x kubectl-unsecret-darwin-amd64
sudo mv kubectl-unsecret-darwin-amd64 /usr/local/bin/kubectl-unsecret
```

#### Windows:

Download the binary from the releases page and add it to your PATH.

## Usage

```sh 
kubectl unsecret [SECRET_NAME] [--namespace=NAMESPACE]
```

- `SECRET_NAME`: The name of the secret to decode.
- `--namespace=NAMESPACE` (optional): The namespace of the secret. Defaults to the current namespace if not specified.

## Example

```sh
kubectl unsecret my-secret --namespace=default
```

## License

This project is licensed under the GNU GENERAL PUBLIC License - see the LICENSE file for details.
