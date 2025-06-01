package main

import (
	"fmt"
	"os"

	tvmCmd "rayyanriaz/tool-version-manager/cmd/tvm"
)

func main() {
	if err := tvmCmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
