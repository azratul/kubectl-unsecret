package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: decode-secret <secret_name> [-n <namespace> | --namespace <namespace>]")
		os.Exit(1)
	}

	var secretName string
	var namespace string = "default"

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-n" || os.Args[i] == "--namespace" {
			if i+1 < len(os.Args) {
				namespace = os.Args[i+1]
				i++
			} else {
				fmt.Println("Error: namespace flag requires a value")
				os.Exit(1)
			}
		} else if strings.HasPrefix(os.Args[i], "--namespace=") {
			namespace = strings.SplitN(os.Args[i], "=", 2)[1]
		} else if !strings.HasPrefix(os.Args[i], "-") {
			secretName = os.Args[i]
		}
	}

	if secretName == "" {
		fmt.Println("Error: secret name is required")
		os.Exit(1)
	}

	cmd := exec.Command("kubectl", "get", "secret", secretName, "-n", namespace, "-o", "jsonpath={.data}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", output)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		fmt.Printf("Error parsing kubectl output: %s\n", err.Error())
		os.Exit(1)
	}

	for key, encodedValue := range result {
		encodedStr, ok := encodedValue.(string)
		if !ok {
			fmt.Printf("Error: value for key %s is not a string\n", key)
			continue
		}

		decodedValue, err := base64.StdEncoding.DecodeString(encodedStr)
		if err != nil {
			fmt.Printf("Error decoding value for key %s: %s\n", key, err.Error())
			continue
		}
		fmt.Printf("%s: %s\n", key, decodedValue)
	}
}
