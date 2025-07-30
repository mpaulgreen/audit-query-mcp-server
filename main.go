package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"audit-query-mcp-server/server"
)

func main() {
	server := server.NewAuditQueryMCPServer()
	server.GetLogger().Info("OpenShift Audit Query MCP Server started")

	// Run setup if requested
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		runSetup()
		return
	}

	// Run tests if this is a test run
	if len(os.Args) > 1 && os.Args[1] == "test" {
		RunAllTests()
		return
	}

	// Run HTTP server for testing if requested
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		runHTTPServer(server)
		return
	}

	// Show usage information
	showUsage()
}

func showUsage() {
	fmt.Println("🚀 OpenShift Audit Query MCP Server")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ./audit-query-mcp-server setup   - Run environment setup and validation")
	fmt.Println("  ./audit-query-mcp-server test    - Run comprehensive tests")
	fmt.Println("  ./audit-query-mcp-server serve   - Start HTTP server for testing")
	fmt.Println("  ./audit-query-mcp-server         - Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Run setup and validation")
	fmt.Println("  ./audit-query-mcp-server setup")
	fmt.Println()
	fmt.Println("  # Run all tests")
	fmt.Println("  ./audit-query-mcp-server test")
	fmt.Println()
	fmt.Println("  # Start HTTP server for testing")
	fmt.Println("  ./audit-query-mcp-server serve")
	fmt.Println()
	fmt.Println("For production use, integrate this server with the MCP protocol.")
	fmt.Println("See README.md for detailed usage instructions.")
}

func runHTTPServer(srv *server.AuditQueryMCPServer) {
	port := ":3000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// Create a simple HTTP server for testing
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"OpenShift Audit Query MCP Server"}`))
	})

	http.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tools := srv.GetTools()
		response := map[string]interface{}{
			"tools": tools,
		}
		jsonResponse, _ := json.MarshalIndent(response, "", "  ")
		w.Write(jsonResponse)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>OpenShift Audit Query MCP Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .method { color: #0066cc; font-weight: bold; }
    </style>
</head>
<body>
    <h1>🚀 OpenShift Audit Query MCP Server</h1>
    <p>This server provides safe, structured access to OpenShift audit logs.</p>
    
    <h2>Available Endpoints:</h2>
    <div class="endpoint">
        <span class="method">GET</span> <code>/health</code> - Health check endpoint
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <code>/tools</code> - List available MCP tools
    </div>
    
    <h2>Usage:</h2>
    <ul>
        <li><code>curl http://localhost:3000/health</code> - Check server health</li>
        <li><code>curl http://localhost:3000/tools</code> - List available tools</li>
    </ul>
    
    <h2>For Production:</h2>
    <p>This HTTP server is for testing only. For production use, integrate with the MCP protocol.</p>
    
    <p><a href="/health">Health Check</a> | <a href="/tools">View Tools</a></p>
</body>
</html>`
		w.Write([]byte(html))
	})

	srv.GetLogger().Infof("Starting HTTP server on port %s", port)
	srv.GetLogger().Info("Visit http://localhost" + port + " for testing interface")
	srv.GetLogger().Info("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(port, nil); err != nil {
		srv.GetLogger().Errorf("HTTP server failed: %v", err)
		fmt.Printf("❌ Failed to start HTTP server: %v\n", err)
		fmt.Printf("💡 Try using a different port: PORT=8081 ./audit-query-mcp-server serve\n")
		os.Exit(1)
	}
}

func runSetup() {
	fmt.Println("🔍 Testing Audit Query MCP Server Setup")
	fmt.Println("======================================")

	// Check Go version
	fmt.Println("Checking Go version...")
	if goVersion, err := exec.Command("go", "version").Output(); err == nil {
		version := strings.TrimSpace(string(goVersion))
		fmt.Printf("✅ Go version: %s\n", version)
	} else {
		fmt.Println("❌ Go is not installed")
		os.Exit(1)
	}

	// Check if we're in the right directory
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Println("❌ go.mod not found. Please run this script from the audit-query-mcp-server directory")
		os.Exit(1)
	}

	// Check dependencies
	fmt.Println("Checking Go dependencies...")
	if cmd := exec.Command("go", "mod", "tidy"); cmd.Run() == nil {
		fmt.Println("✅ Dependencies are up to date")
	} else {
		fmt.Println("❌ Failed to update dependencies")
		os.Exit(1)
	}

	// Check if .env file exists
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println("✅ .env file found")
	} else {
		fmt.Println("⚠️  .env file not found. Creating from template...")
		if _, err := os.Stat("env.example"); err == nil {
			if input, err := os.ReadFile("env.example"); err == nil {
				if err := os.WriteFile(".env", input, 0644); err == nil {
					fmt.Println("✅ Created .env file from template")
					fmt.Println("⚠️  Please edit .env and add your OpenAI API key")
				} else {
					fmt.Println("❌ Failed to create .env file")
					os.Exit(1)
				}
			} else {
				fmt.Println("❌ Failed to read env.example")
				os.Exit(1)
			}
		} else {
			fmt.Println("❌ env.example not found")
			os.Exit(1)
		}
	}

	// Check OpenShift CLI
	fmt.Println("Checking OpenShift CLI...")
	if ocVersion, err := exec.Command("oc", "version", "--client").Output(); err == nil {
		lines := strings.Split(string(ocVersion), "\n")
		if len(lines) > 0 {
			fmt.Printf("✅ OpenShift CLI: %s\n", strings.TrimSpace(lines[0]))
		}
	} else {
		fmt.Println("❌ OpenShift CLI (oc) is not installed")
		fmt.Println("   Install from: https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/getting-started-cli.html")
		os.Exit(1)
	}

	// Test OpenShift connection
	fmt.Println("Testing OpenShift connection...")
	if whoami, err := exec.Command("oc", "whoami").Output(); err == nil {
		user := strings.TrimSpace(string(whoami))
		fmt.Printf("✅ Connected to OpenShift as: %s\n", user)
	} else {
		fmt.Println("❌ Not connected to OpenShift cluster")
		fmt.Println("   Run: oc login <cluster-url>")
		os.Exit(1)
	}

	// Test basic OpenShift access
	fmt.Println("Testing OpenShift access...")
	if cmd := exec.Command("oc", "get", "nodes"); cmd.Run() == nil {
		fmt.Println("✅ Can access OpenShift cluster nodes")
	} else {
		fmt.Println("❌ Cannot access OpenShift cluster nodes")
		fmt.Println("   Check your permissions and cluster status")
		os.Exit(1)
	}

	// Test audit log access
	fmt.Println("Testing audit log access...")
	auditCmd := exec.Command("oc", "adm", "node-logs", "--role=master", "--path=kube-apiserver/audit.log", "--since=1h")
	auditCmd.Stdout = nil
	auditCmd.Stderr = nil
	if auditCmd.Run() == nil {
		fmt.Println("✅ Can access audit logs")
	} else {
		fmt.Println("⚠️  Cannot access audit logs (this might be normal depending on permissions)")
		fmt.Println("   You may need cluster-admin privileges for audit log access")
	}

	fmt.Println("")
	fmt.Println("🎉 Setup test completed!")
	fmt.Println("")
	fmt.Println("To run the MCP server tests:")
	fmt.Println("  ./audit-query-mcp-server test")
	fmt.Println("")
	fmt.Println("To start the server:")
	fmt.Println("  ./audit-query-mcp-server serve")
}
