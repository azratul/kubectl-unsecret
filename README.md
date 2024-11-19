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
wget https://github.com/azratul/kubectl-unsecret/releases/download/v1.1.0/kubectl-unsecret-linux-amd64
chmod +x kubectl-unsecret-linux-amd64
sudo mv kubectl-unsecret-linux-amd64 /usr/local/bin/kubectl-unsecret
```

#### macOS

```sh
wget https://github.com/azratul/kubectl-unsecret/releases/download/v1.1.0/kubectl-unsecret-darwin-amd64
chmod +x kubectl-unsecret-darwin-amd64
sudo mv kubectl-unsecret-darwin-amd64 /usr/local/bin/kubectl-unsecret
```

#### Windows:

Download the binary from the releases page and add it to your PATH.

## Usage

```sh 
kubectl unsecret <secret_name> [-n <namespace> | --namespace <namespace>] [-o <output_format>]
```

Options:
- `-n, --namespace <namespace>`: Specify the namespace (default is "default")
- `-o, --output <output_format>`: Specify the output format (text|json|yaml) (default is "text")
- `--help`: Show this help message
- `--version`: Show version information

## Example

```sh
kubectl unsecret my-secret --namespace=default -o json
```

## License

This project is licensed under the GNU GENERAL PUBLIC License - see the LICENSE file for details.
