package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProtoConfig struct {
	ServiceName string
	ProtoFiles  []string
	OutputDir   string
}

var services = map[string]ProtoConfig{
	"user": {
		ServiceName: "user",
		ProtoFiles:  []string{"../user/proto/user.proto"},
		OutputDir:   "../user/proto",
	},
	"product": {
		ServiceName: "product",
		ProtoFiles:  []string{"../product/proto/product.proto"},
		OutputDir:   "../product/proto",
	},
	"order": {
		ServiceName: "order",
		ProtoFiles:  []string{"../order/proto/order.proto", "../order/proto/product.proto"},
		OutputDir:   "../order/proto",
	},
	"payment": {
		ServiceName: "payment",
		ProtoFiles:  []string{"../payment/proto/payment.proto", "../payment/proto/order.proto"},
		OutputDir:   "../payment/proto",
	},
	"broker": {
		ServiceName: "broker",
		ProtoFiles: []string{
			"../broker/proto/user.proto",
			"../broker/proto/product.proto",
			"../broker/proto/order.proto",
			"../broker/proto/payment.proto",
		},
		OutputDir: "../broker/proto",
	},
}

func main() {
	var (
		serviceName = flag.String("service", "", "Service name to generate proto for")
		all         = flag.Bool("all", false, "Generate all proto files")
		list        = flag.Bool("list", false, "List available services")
		help        = flag.Bool("help", false, "Show help")
		verbose     = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *list {
		listServices()
		return
	}

	// Check prerequisites
	if !checkPrerequisites(*verbose) {
		log.Fatal("Prerequisites not met. Please install required tools.")
	}

	if *all {
		generateAllServices(*verbose)
	} else if *serviceName != "" {
		generateService(*serviceName, *verbose)
	} else {
		fmt.Println("No action specified. Use -help for usage information.")
		showHelp()
	}
}

func showHelp() {
	fmt.Println("gRPC Proto Generator for Go")
	fmt.Println("===========================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run proto-generator.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -service string    Generate proto for specific service")
	fmt.Println("  -all              Generate all proto files")
	fmt.Println("  -list             List available services")
	fmt.Println("  -verbose          Verbose output")
	fmt.Println("  -help             Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run proto-generator.go -all")
	fmt.Println("  go run proto-generator.go -service user")
	fmt.Println("  go run proto-generator.go -service product")
	fmt.Println("  go run proto-generator.go -list")
	fmt.Println()
	fmt.Println("Prerequisites:")
	fmt.Println("  - protoc (Protocol Compiler)")
	fmt.Println("  - protoc-gen-go")
	fmt.Println("  - protoc-gen-go-grpc")
	fmt.Println()
	fmt.Println("Install prerequisites:")
	fmt.Println("  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	fmt.Println("  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
}

func listServices() {
	fmt.Println("Available services:")
	for name, config := range services {
		fmt.Printf("  %-10s (%d proto files)\n", name, len(config.ProtoFiles))
		if len(config.ProtoFiles) > 0 {
			for _, protoFile := range config.ProtoFiles {
				fmt.Printf("    → %s\n", protoFile)
			}
		}
	}
}

func checkPrerequisites(verbose bool) bool {
	if verbose {
		fmt.Println("Checking prerequisites...")
	}

	// Check protoc
	if err := exec.Command("protoc", "--version").Run(); err != nil {
		fmt.Println("✗ protoc not found. Please install Protocol Compiler.")
		fmt.Println("  Download from: https://github.com/protocolbuffers/protobuf/releases")
		return false
	}
	if verbose {
		fmt.Println("✓ protoc found")
	}

	// Check protoc-gen-go
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		cmd := exec.Command("go", "env", "GOPATH")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("✗ Could not determine GOPATH")
			return false
		}
		gopath = strings.TrimSpace(string(output))
	}

	protocGenGo := filepath.Join(gopath, "bin", "protoc-gen-go")
	if _, err := os.Stat(protocGenGo); os.IsNotExist(err) {
		protocGenGo += ".exe" // Try Windows extension
		if _, err := os.Stat(protocGenGo); os.IsNotExist(err) {
			fmt.Println("✗ protoc-gen-go not found")
			fmt.Println("  Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest")
			return false
		}
	}
	if verbose {
		fmt.Println("✓ protoc-gen-go found")
	}

	// Check protoc-gen-go-grpc
	protocGenGoGrpc := filepath.Join(gopath, "bin", "protoc-gen-go-grpc")
	if _, err := os.Stat(protocGenGoGrpc); os.IsNotExist(err) {
		protocGenGoGrpc += ".exe" // Try Windows extension
		if _, err := os.Stat(protocGenGoGrpc); os.IsNotExist(err) {
			fmt.Println("✗ protoc-gen-go-grpc not found")
			fmt.Println("  Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
			return false
		}
	}
	if verbose {
		fmt.Println("✓ protoc-gen-go-grpc found")
	}

	return true
}

func generateService(serviceName string, verbose bool) {
	config, exists := services[serviceName]
	if !exists {
		fmt.Printf("Unknown service: %s\n", serviceName)
		fmt.Println("Available services:")
		for name := range services {
			fmt.Printf("  %s\n", name)
		}
		return
	}

	fmt.Printf("Generating %s service proto files...\n", serviceName)

	for _, protoFile := range config.ProtoFiles {
		generateProtoFile(protoFile, config.OutputDir, verbose)
	}

	fmt.Printf("✓ %s service proto files generated successfully\n", serviceName)
}

func generateAllServices(verbose bool) {
	fmt.Println("Generating all proto files...")
	fmt.Println()

	for name := range services {
		fmt.Printf("Generating %s service...\n", name)
		generateService(name, verbose)
		fmt.Println()
	}

	fmt.Println("✓ All proto files generated successfully!")
}

func generateProtoFile(protoFile, outputDir string, verbose bool) {
	if _, err := os.Stat(protoFile); os.IsNotExist(err) {
		fmt.Printf("Proto file not found: %s\n", protoFile)
		return
	}

	if verbose {
		fmt.Printf("Generating Go code for: %s\n", protoFile)
	}

	protoDir := filepath.Dir(protoFile)
	fileName := filepath.Base(protoFile)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// Generate Go code
	cmd := exec.Command("protoc",
		fmt.Sprintf("--proto_path=%s", protoDir),
		fmt.Sprintf("--go_out=%s", outputDir),
		fmt.Sprintf("--go-grpc_out=%s", outputDir),
		fileName)

	if verbose {
		fmt.Printf("Running: %s\n", cmd.String())
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("✗ Failed to generate: %s\n", fileName)
		fmt.Printf("Error: %v\n", err)
		if len(output) > 0 {
			fmt.Printf("Output: %s\n", string(output))
		}
		return
	}

	if verbose {
		fmt.Printf("✓ Generated: %s\n", fileName)

		// List generated files
		baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		pbFile := filepath.Join(outputDir, baseName+".pb.go")
		grpcFile := filepath.Join(outputDir, baseName+"_grpc.pb.go")

		if _, err := os.Stat(pbFile); err == nil {
			fmt.Printf("  → %s\n", pbFile)
		}
		if _, err := os.Stat(grpcFile); err == nil {
			fmt.Printf("  → %s\n", grpcFile)
		}
	}
}

func promptUser(message string) string {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
