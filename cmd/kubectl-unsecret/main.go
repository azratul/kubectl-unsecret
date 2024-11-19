package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

const VERSION = "1.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	// Parse arguments
	secretName, namespace, outputFormat, err := parseArgs(args)
	if err != nil {
		logError(err.Error())
		os.Exit(1)
	}

	// Fetch secret
	secretData, err := getSecret(secretName, namespace)
	if err != nil {
		logError(err.Error())
		os.Exit(1)
	}

	// Handle output based on the specified format
	if err := handleOutput(secretData, outputFormat); err != nil {
		logError(err.Error())
		os.Exit(1)
	}
}

func parseArgs(args []string) (string, string, string, error) {
	var secretName string
	namespace := "default"
	outputFormat := "text"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n", "--namespace":
			if i+1 < len(args) {
				namespace = args[i+1]
				i++
			} else {
				return "", "", "", fmt.Errorf("Error: namespace flag requires a value")
			}
		case "-o", "--output":
			if i+1 < len(args) {
				outputFormat = args[i+1]
				i++
			} else {
				return "", "", "", fmt.Errorf("Error: output flag requires a value (text|json|yaml)")
			}
		case "--help":
			printUsage()
			os.Exit(0)
		case "--version":
			fmt.Printf("kubectl-unsecret version %s\n", VERSION)
			os.Exit(0)
		default:
			if strings.HasPrefix(args[i], "--namespace=") {
				namespace = strings.SplitN(args[i], "=", 2)[1]
			} else if !strings.HasPrefix(args[i], "-") {
				secretName = args[i]
			} else {
				return "", "", "", fmt.Errorf("Error: unknown flag %s", args[i])
			}
		}
	}

	if secretName == "" {
		return "", "", "", fmt.Errorf("Error: secret name is required")
	}

	return secretName, namespace, outputFormat, nil
}

func getSecret(secretName, namespace string) (map[string]interface{}, error) {
	cmd := exec.Command("kubectl", "get", "secret", secretName, "-n", namespace, "-o", "jsonpath={.data}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("kubectl error: %s\n%s", err.Error(), string(output))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("Error parsing kubectl output: %s", err.Error())
	}

	return result, nil
}

func handleOutput(secretData map[string]interface{}, outputFormat string) error {
	switch outputFormat {
	case "text":
		return printText(secretData)
	case "json":
		return printJSON(secretData)
	case "yaml":
		return printYAML(secretData)
	default:
		return fmt.Errorf("Error: unsupported output format %s", outputFormat)
	}
}

func printText(secretData map[string]interface{}) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "KEY\tVALUE")
	fmt.Fprintln(w, "----\t-----")

	for key, encodedValue := range secretData {
		encodedStr, ok := encodedValue.(string)
		if !ok {
			logError(fmt.Sprintf("Value for key %s is not a string", key))
			continue
		}

		decodedValue, err := base64.StdEncoding.DecodeString(encodedStr)
		if err != nil {
			logError(fmt.Sprintf("Error decoding value for key %s: %s", key, err.Error()))
			continue
		}

		lines := strings.Split(string(decodedValue), "\n")

		if len(lines) > 0 {
			fmt.Fprintf(w, "%s\t%s\n", key, lines[0])
		}

		for i := 1; i < len(lines); i++ {
			fmt.Fprintf(w, "\t%s\n", lines[i])
		}
	}

	w.Flush()
	return nil
}

func printYAML(secretData map[string]interface{}) error {
	decodedSecret := make(map[string]string)
	for key, encodedValue := range secretData {
		encodedStr, ok := encodedValue.(string)
		if !ok {
			return fmt.Errorf("Error: value for key %s is not a string", key)
		}

		decodedValue, err := base64.StdEncoding.DecodeString(encodedStr)
		if err != nil {
			return fmt.Errorf("Error decoding value for key %s: %s", key, err.Error())
		}

		decodedSecret[key] = string(decodedValue)
	}

	yamlData, err := yaml.Marshal(decodedSecret)
	if err != nil {
		return fmt.Errorf("Error encoding secret to YAML: %s", err.Error())
	}

	yamlString := string(yamlData)
	yamlString = strings.TrimSuffix(yamlString, "\n")

	fmt.Println(yamlString)
	return nil
}

func printJSON(secretData map[string]interface{}) error {
	decodedSecret := make(map[string]string)
	for key, encodedValue := range secretData {
		encodedStr, ok := encodedValue.(string)
		if !ok {
			return fmt.Errorf("Error: value for key %s is not a string", key)
		}

		decodedValue, err := base64.StdEncoding.DecodeString(encodedStr)
		if err != nil {
			return fmt.Errorf("Error decoding value for key %s: %s", key, err.Error())
		}

		decodedSecret[key] = string(decodedValue)
	}

	jsonData, err := json.MarshalIndent(decodedSecret, "", "  ")
	if err != nil {
		return fmt.Errorf("Error encoding secret to JSON: %s", err.Error())
	}
	fmt.Println(string(jsonData))
	return nil
}

func logError(message string) {
	fmt.Fprintf(os.Stderr, "\033[31mError: %s\033[0m\n", message) // Red color for errors
}

func printUsage() {
	fmt.Println(`Usage: kubectl unsecret <secret_name> [-n <namespace> | --namespace <namespace>] [-o <output_format>]
Options:
  -n, --namespace <namespace>   Specify the namespace (default is "default")
  -o, --output <output_format>  Specify the output format (text|json|yaml) (default is "text")
  --help                        Show this help message
  --version                     Show version information`)
}
