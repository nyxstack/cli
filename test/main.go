package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nyxstack/cli"
)

var (
	// Global flags
	verbose bool
	config  string
	dryRun  bool
	timeout time.Duration
)

func main() {
	root := cli.Root("cloudctl").
		Description("Cloud infrastructure management CLI").
		Flag(&verbose, "verbose", "v", false, "Enable verbose output").
		Flag(&config, "config", "c", "~/.cloudctl/config.yaml", "Configuration file path").
		Flag(&timeout, "timeout", "t", 5*time.Minute, "Global operation timeout").
		Action(func(ctx context.Context, cmd *cli.Command) error {
			cmd.ShowHelp()
			return nil
		}).
		PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
			if verbose {
				fmt.Printf("[DEBUG] Executing command: %s\n", cmd.GetName())
				fmt.Printf("[DEBUG] Config: %s\n", config)
			}
			return nil
		})

	// Add all command groups
	root.AddCommand(buildDeployCommands())
	root.AddCommand(buildServerCommands())
	root.AddCommand(buildDatabaseCommands())
	root.AddCommand(buildStorageCommands())
	root.AddCommand(buildMonitorCommands())
	root.AddCommand(buildNetworkCommands())

	// Add completion
	cli.AddCompletion(root)

	// Execute
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// buildDeployCommands creates the deploy command tree
func buildDeployCommands() *cli.Command {
	var (
		replicas    int
		tag         string
		force       bool
		rollback    bool
		canary      int
		healthCheck bool
	)

	deploy := cli.Cmd("deploy").
		Description("Deployment management").
		Flag(&dryRun, "dry-run", "", false, "Simulate deployment without making changes")

	// deploy app
	deployApp := cli.Cmd("app").
		Description("Deploy an application").
		Arg("name", "Application name", true).
		Arg("environment", "Target environment (dev/staging/prod)", true).
		Flag(&replicas, "replicas", "r", 3, "Number of replicas").
		Flag(&tag, "tag", "", "latest", "Docker image tag").
		Flag(&force, "force", "f", false, "Force deployment even if health checks fail").
		Flag(&canary, "canary", "", 0, "Canary deployment percentage (0-100)").
		PreRun(func(ctx context.Context, cmd *cli.Command) error {
			if verbose {
				fmt.Println("[VALIDATE] Checking deployment prerequisites...")
			}
			time.Sleep(100 * time.Millisecond)
			return nil
		}).
		Action(func(ctx context.Context, cmd *cli.Command, name, env string) error {
			fmt.Printf("🚀 Deploying %s to %s environment\n", name, env)
			if dryRun {
				fmt.Println("   [DRY RUN] No changes will be made")
			}
			fmt.Printf("   Image: %s:%s\n", name, tag)
			fmt.Printf("   Replicas: %d\n", replicas)

			if canary > 0 {
				fmt.Printf("   Canary: %d%% traffic\n", canary)
			}

			// Simulate deployment steps
			steps := []string{
				"Pulling image",
				"Creating deployment manifest",
				"Applying configuration",
				"Scaling replicas",
				"Waiting for pods to be ready",
			}

			for i, step := range steps {
				fmt.Printf("   [%d/%d] %s...", i+1, len(steps), step)
				time.Sleep(200 * time.Millisecond)
				fmt.Println(" ✓")
			}

			fmt.Printf("✅ Successfully deployed %s to %s\n", name, env)
			return nil
		})

	// deploy rollback
	deployRollback := cli.Cmd("rollback").
		Description("Rollback a deployment").
		Arg("name", "Application name", true).
		Arg("environment", "Target environment", true).
		Flag(&rollback, "revision", "", false, "Rollback to specific revision").
		Action(func(ctx context.Context, cmd *cli.Command, name, env string) error {
			fmt.Printf("⏮️  Rolling back %s in %s environment\n", name, env)
			time.Sleep(500 * time.Millisecond)
			fmt.Println("   Finding previous revision...")
			time.Sleep(300 * time.Millisecond)
			fmt.Println("   Applying rollback...")
			time.Sleep(400 * time.Millisecond)
			fmt.Printf("✅ Successfully rolled back %s\n", name)
			return nil
		})

	// deploy status
	deployStatus := cli.Cmd("status").
		Description("Check deployment status").
		Arg("name", "Application name", true).
		Flag(&healthCheck, "health", "", true, "Include health check results").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("📊 Deployment status for: %s\n", name)
			fmt.Println("   Environment: production")
			fmt.Println("   Status: Running")
			fmt.Println("   Replicas: 3/3 ready")
			fmt.Println("   Version: v1.2.3")
			fmt.Println("   Uptime: 5d 12h 34m")
			if healthCheck {
				fmt.Println("\n   Health Checks:")
				fmt.Println("   ✓ Liveness: Passing")
				fmt.Println("   ✓ Readiness: Passing")
				fmt.Println("   ✓ Startup: Passing")
			}
			return nil
		})

	deploy.AddCommand(deployApp)
	deploy.AddCommand(deployRollback)
	deploy.AddCommand(deployStatus)

	return deploy
}

// buildServerCommands creates the server command tree
func buildServerCommands() *cli.Command {
	var (
		instanceType string
		region       string
		count        int
		diskSize     int
		autoScale    bool
	)

	server := cli.Cmd("server").
		Description("Server management")

	// server create
	serverCreate := cli.Cmd("create").
		Description("Create new server instances").
		Arg("name", "Server name", true).
		Flag(&instanceType, "type", "", "t3.medium", "Instance type").
		Flag(&region, "region", "", "us-east-1", "AWS region").
		Flag(&count, "count", "n", 1, "Number of instances").
		Flag(&diskSize, "disk", "d", 50, "Disk size in GB").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🖥️  Creating %d server(s): %s\n", count, name)
			fmt.Printf("   Type: %s\n", instanceType)
			fmt.Printf("   Region: %s\n", region)
			fmt.Printf("   Disk: %dGB\n", diskSize)

			for i := 1; i <= count; i++ {
				fmt.Printf("   [%d/%d] Provisioning instance...", i, count)
				time.Sleep(300 * time.Millisecond)
				fmt.Printf(" ✓ (ID: srv-%s-%03d)\n", name, i)
			}

			fmt.Printf("✅ Successfully created %d server(s)\n", count)
			return nil
		})

	// server list
	serverList := cli.Cmd("list").
		Description("List all servers").
		Flag(&region, "region", "", "", "Filter by region").
		Action(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("📋 Server Instances:")
			servers := []struct {
				name   string
				status string
				ip     string
				region string
			}{
				{"web-01", "running", "10.0.1.15", "us-east-1"},
				{"web-02", "running", "10.0.1.16", "us-east-1"},
				{"db-01", "running", "10.0.2.20", "us-west-2"},
				{"cache-01", "stopped", "10.0.3.10", "eu-west-1"},
			}

			fmt.Printf("%-15s %-10s %-15s %-12s\n", "NAME", "STATUS", "IP", "REGION")
			fmt.Println("-----------------------------------------------------------")
			for _, s := range servers {
				if region == "" || s.region == region {
					statusIcon := "●"
					if s.status == "stopped" {
						statusIcon = "○"
					}
					fmt.Printf("%-15s %s %-8s %-15s %-12s\n", s.name, statusIcon, s.status, s.ip, s.region)
				}
			}
			return nil
		})

	// server stop
	serverStop := cli.Cmd("stop").
		Description("Stop server instances").
		Arg("name", "Server name or ID", true).
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("⏸️  Stopping server: %s\n", name)
			time.Sleep(400 * time.Millisecond)
			fmt.Println("   Draining connections...")
			time.Sleep(300 * time.Millisecond)
			fmt.Println("   Shutting down gracefully...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("✅ Server %s stopped\n", name)
			return nil
		})

	// server start
	serverStart := cli.Cmd("start").
		Description("Start server instances").
		Arg("name", "Server name or ID", true).
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("▶️  Starting server: %s\n", name)
			time.Sleep(500 * time.Millisecond)
			fmt.Println("   Booting system...")
			time.Sleep(300 * time.Millisecond)
			fmt.Println("   Initializing services...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("✅ Server %s started\n", name)
			return nil
		})

	// server scale
	serverScale := cli.Cmd("scale").
		Description("Scale server group").
		Arg("group", "Server group name", true).
		Arg("replicas", "Target replica count", true).
		Flag(&autoScale, "auto", "a", false, "Enable auto-scaling").
		Action(func(ctx context.Context, cmd *cli.Command, group string, replicas int) error {
			fmt.Printf("📈 Scaling %s to %d replicas\n", group, replicas)
			if autoScale {
				fmt.Println("   Auto-scaling enabled")
			}
			time.Sleep(600 * time.Millisecond)
			fmt.Println("   Adjusting capacity...")
			time.Sleep(400 * time.Millisecond)
			fmt.Printf("✅ Scaled %s successfully\n", group)
			return nil
		})

	server.AddCommand(serverCreate)
	server.AddCommand(serverList)
	server.AddCommand(serverStop)
	server.AddCommand(serverStart)
	server.AddCommand(serverScale)

	return server
}

// buildDatabaseCommands creates the database command tree
func buildDatabaseCommands() *cli.Command {
	var (
		engine      string
		version     string
		storage     int
		backup      bool
		replication bool
	)

	db := cli.Cmd("database").
		Description("Database management")

	// db create
	dbCreate := cli.Cmd("create").
		Description("Create a new database").
		Arg("name", "Database name", true).
		Flag(&engine, "engine", "e", "postgres", "Database engine (postgres/mysql/redis)").
		Flag(&version, "version", "", "14.0", "Engine version").
		Flag(&storage, "storage", "s", 100, "Storage size in GB").
		Flag(&replication, "replicas", "r", false, "Enable replication").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🗄️  Creating database: %s\n", name)
			fmt.Printf("   Engine: %s %s\n", engine, version)
			fmt.Printf("   Storage: %dGB\n", storage)
			if replication {
				fmt.Println("   Replication: Enabled (1 primary + 2 replicas)")
			}

			time.Sleep(800 * time.Millisecond)
			fmt.Println("   Provisioning storage...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("   Initializing database...")
			time.Sleep(400 * time.Millisecond)
			fmt.Println("   Configuring security...")
			time.Sleep(300 * time.Millisecond)

			fmt.Printf("✅ Database %s created\n", name)
			fmt.Printf("   Endpoint: %s.db.internal:5432\n", name)
			return nil
		})

	// db backup
	dbBackup := cli.Cmd("backup").
		Description("Backup a database").
		Arg("name", "Database name", true).
		Flag(&backup, "incremental", "i", false, "Incremental backup").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			backupType := "Full"
			if backup {
				backupType = "Incremental"
			}
			fmt.Printf("💾 Creating %s backup for: %s\n", backupType, name)
			time.Sleep(1 * time.Second)
			fmt.Println("   Backing up data...")
			time.Sleep(700 * time.Millisecond)
			fmt.Println("   Compressing...")
			time.Sleep(400 * time.Millisecond)
			fmt.Printf("✅ Backup completed: backup-%s-%s.sql.gz\n", name, time.Now().Format("20060102-150405"))
			return nil
		})

	// db restore
	dbRestore := cli.Cmd("restore").
		Description("Restore database from backup").
		Arg("name", "Database name", true).
		Arg("backup-file", "Backup file path", true).
		Action(func(ctx context.Context, cmd *cli.Command, name, backupFile string) error {
			fmt.Printf("♻️  Restoring %s from %s\n", name, backupFile)
			time.Sleep(500 * time.Millisecond)
			fmt.Println("   Validating backup...")
			time.Sleep(600 * time.Millisecond)
			fmt.Println("   Restoring data...")
			time.Sleep(1 * time.Second)
			fmt.Println("   Rebuilding indexes...")
			time.Sleep(400 * time.Millisecond)
			fmt.Printf("✅ Database %s restored\n", name)
			return nil
		})

	// db migrate
	dbMigrate := cli.Cmd("migrate").
		Description("Run database migrations").
		Arg("name", "Database name", true).
		Arg("direction", "Migration direction (up/down)", false).
		Action(func(ctx context.Context, cmd *cli.Command, name, direction string) error {
			if direction == "" {
				direction = "up"
			}
			fmt.Printf("🔄 Running migrations %s for: %s\n", direction, name)
			migrations := []string{"001_create_users", "002_add_indexes", "003_alter_schema"}
			for i, m := range migrations {
				fmt.Printf("   [%d/%d] %s...", i+1, len(migrations), m)
				time.Sleep(300 * time.Millisecond)
				fmt.Println(" ✓")
			}
			fmt.Println("✅ Migrations completed")
			return nil
		})

	db.AddCommand(dbCreate)
	db.AddCommand(dbBackup)
	db.AddCommand(dbRestore)
	db.AddCommand(dbMigrate)

	return db
}

// buildStorageCommands creates the storage command tree
func buildStorageCommands() *cli.Command {
	var (
		bucket     string
		recursive  bool
		acl        string
		versioning bool
		encryption bool
	)

	storage := cli.Cmd("storage").
		Description("Object storage management")

	// storage upload
	storageUpload := cli.Cmd("upload").
		Description("Upload files to storage").
		Arg("source", "Source file or directory", true).
		Arg("destination", "Destination path", true).
		Flag(&bucket, "bucket", "b", "default", "Storage bucket").
		Flag(&recursive, "recursive", "r", false, "Upload directory recursively").
		Flag(&encryption, "encrypt", "e", true, "Enable encryption").
		Action(func(ctx context.Context, cmd *cli.Command, source, dest string) error {
			fmt.Printf("⬆️  Uploading %s to %s/%s\n", source, bucket, dest)
			if recursive {
				fmt.Println("   Mode: Recursive")
			}
			if encryption {
				fmt.Println("   Encryption: AES-256")
			}

			files := []string{"file1.txt", "file2.json", "file3.csv"}
			for i, f := range files {
				fmt.Printf("   [%d/%d] %s...", i+1, len(files), f)
				time.Sleep(200 * time.Millisecond)
				fmt.Println(" ✓")
			}
			fmt.Println("✅ Upload completed")
			return nil
		})

	// storage download
	storageDownload := cli.Cmd("download").
		Description("Download files from storage").
		Arg("source", "Source path", true).
		Arg("destination", "Local destination", true).
		Flag(&bucket, "bucket", "b", "default", "Storage bucket").
		Action(func(ctx context.Context, cmd *cli.Command, source, dest string) error {
			fmt.Printf("⬇️  Downloading %s/%s to %s\n", bucket, source, dest)
			time.Sleep(500 * time.Millisecond)
			fmt.Println("   Transferring...")
			time.Sleep(400 * time.Millisecond)
			fmt.Println("✅ Download completed")
			return nil
		})

	// storage list
	storageList := cli.Cmd("list").
		Description("List storage objects").
		Arg("path", "Path to list", false).
		Flag(&bucket, "bucket", "b", "default", "Storage bucket").
		Action(func(ctx context.Context, cmd *cli.Command, path string) error {
			if path == "" {
				path = "/"
			}
			fmt.Printf("📂 Contents of %s:%s\n\n", bucket, path)
			fmt.Printf("%-30s %-12s %-20s\n", "NAME", "SIZE", "MODIFIED")
			fmt.Println("----------------------------------------------------------------")
			objects := []struct {
				name     string
				size     string
				modified string
			}{
				{"logs/app.log", "2.5 MB", "2025-11-18 10:30"},
				{"backups/db.sql.gz", "150 MB", "2025-11-17 23:00"},
				{"assets/logo.png", "45 KB", "2025-11-15 14:20"},
			}
			for _, obj := range objects {
				fmt.Printf("%-30s %-12s %-20s\n", obj.name, obj.size, obj.modified)
			}
			return nil
		})

	// storage bucket
	storageBucket := cli.Cmd("bucket").
		Description("Bucket management")

	bucketCreate := cli.Cmd("create").
		Description("Create a new bucket").
		Arg("name", "Bucket name", true).
		Flag(&versioning, "versioning", "v", false, "Enable versioning").
		Flag(&acl, "acl", "", "private", "Access control (private/public-read)").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🪣 Creating bucket: %s\n", name)
			fmt.Printf("   ACL: %s\n", acl)
			if versioning {
				fmt.Println("   Versioning: Enabled")
			}
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("✅ Bucket %s created\n", name)
			return nil
		})

	bucketDelete := cli.Cmd("delete").
		Description("Delete a bucket").
		Arg("name", "Bucket name", true).
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🗑️  Deleting bucket: %s\n", name)
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("✅ Bucket %s deleted\n", name)
			return nil
		})

	storageBucket.AddCommand(bucketCreate)
	storageBucket.AddCommand(bucketDelete)

	storage.AddCommand(storageUpload)
	storage.AddCommand(storageDownload)
	storage.AddCommand(storageList)
	storage.AddCommand(storageBucket)

	return storage
}

// buildMonitorCommands creates the monitoring command tree
func buildMonitorCommands() *cli.Command {
	var (
		duration  string
		metric    string
		threshold float64
	)

	monitor := cli.Cmd("monitor").
		Description("Monitoring and metrics")

	// monitor metrics
	monitorMetrics := cli.Cmd("metrics").
		Description("View metrics").
		Arg("resource", "Resource to monitor", true).
		Flag(&duration, "duration", "d", "1h", "Time duration (1h/24h/7d)").
		Flag(&metric, "metric", "m", "cpu", "Metric type (cpu/memory/disk/network)").
		Action(func(ctx context.Context, cmd *cli.Command, resource string) error {
			fmt.Printf("📊 Metrics for %s (%s, last %s)\n\n", resource, metric, duration)
			time.Sleep(300 * time.Millisecond)

			fmt.Println("CPU Usage:")
			fmt.Println("  ▓▓▓▓▓▓▓░░░  65%")
			fmt.Println("\nMemory Usage:")
			fmt.Println("  ▓▓▓▓▓▓▓▓░░  78%")
			fmt.Println("\nDisk I/O:")
			fmt.Println("  ▓▓▓░░░░░░░  32%")
			fmt.Println("\nNetwork:")
			fmt.Println("  ▓▓▓▓▓░░░░░  51%")

			return nil
		})

	// monitor logs
	monitorLogs := cli.Cmd("logs").
		Description("View logs").
		Arg("resource", "Resource name", true).
		Flag(&duration, "since", "s", "10m", "Show logs since (10m/1h/24h)").
		Action(func(ctx context.Context, cmd *cli.Command, resource string) error {
			fmt.Printf("📜 Logs for %s (last %s)\n\n", resource, duration)
			time.Sleep(200 * time.Millisecond)

			logs := []string{
				"[INFO] Application started",
				"[INFO] Connected to database",
				"[WARN] High memory usage detected",
				"[INFO] Request processed: 200 OK",
				"[ERROR] Connection timeout to service-X",
				"[INFO] Retrying connection...",
				"[INFO] Connection successful",
			}

			for _, log := range logs {
				fmt.Printf("%s %s\n", time.Now().Format("15:04:05"), log)
				time.Sleep(100 * time.Millisecond)
			}

			return nil
		})

	// monitor alerts
	monitorAlerts := cli.Cmd("alerts").
		Description("Manage alerts")

	alertList := cli.Cmd("list").
		Description("List active alerts").
		Action(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("🚨 Active Alerts:")
			fmt.Println()
			alerts := []struct {
				severity string
				resource string
				message  string
			}{
				{"CRITICAL", "db-prod-01", "High CPU usage: 95%"},
				{"WARNING", "web-server-03", "Memory usage above threshold"},
				{"INFO", "cache-01", "Connection pool at 80%"},
			}

			for _, a := range alerts {
				icon := "⚠️"
				if a.severity == "CRITICAL" {
					icon = "🔴"
				} else if a.severity == "INFO" {
					icon = "🔵"
				}
				fmt.Printf("%s [%s] %s: %s\n", icon, a.severity, a.resource, a.message)
			}
			return nil
		})

	alertCreate := cli.Cmd("create").
		Description("Create an alert rule").
		Arg("name", "Alert name", true).
		Flag(&metric, "metric", "m", "cpu", "Metric to monitor").
		Flag(&threshold, "threshold", "t", 80.0, "Alert threshold").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🔔 Creating alert: %s\n", name)
			fmt.Printf("   Metric: %s\n", metric)
			fmt.Printf("   Threshold: %.1f%%\n", threshold)
			time.Sleep(200 * time.Millisecond)
			fmt.Println("✅ Alert created")
			return nil
		})

	monitorAlerts.AddCommand(alertList)
	monitorAlerts.AddCommand(alertCreate)

	monitor.AddCommand(monitorMetrics)
	monitor.AddCommand(monitorLogs)
	monitor.AddCommand(monitorAlerts)

	return monitor
}

// buildNetworkCommands creates the network command tree
func buildNetworkCommands() *cli.Command {
	var (
		cidr       string
		vpcID      string
		protocol   string
		port       int
		sourceIP   string
		targetPort int
	)

	network := cli.Cmd("network").
		Description("Network management")

	// network vpc
	networkVPC := cli.Cmd("vpc").
		Description("VPC management")

	vpcCreate := cli.Cmd("create").
		Description("Create a VPC").
		Arg("name", "VPC name", true).
		Flag(&cidr, "cidr", "", "10.0.0.0/16", "CIDR block").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🌐 Creating VPC: %s\n", name)
			fmt.Printf("   CIDR: %s\n", cidr)
			time.Sleep(400 * time.Millisecond)
			fmt.Println("   Creating subnets...")
			time.Sleep(300 * time.Millisecond)
			fmt.Println("   Configuring routing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("✅ VPC %s created (vpc-12345678)\n", name)
			return nil
		})

	vpcList := cli.Cmd("list").
		Description("List VPCs").
		Action(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("🌐 Virtual Private Clouds:")
			fmt.Println()
			fmt.Printf("%-20s %-15s %-10s\n", "NAME", "CIDR", "VPC-ID")
			fmt.Println("--------------------------------------------------")
			vpcs := []struct {
				name string
				cidr string
				id   string
			}{
				{"prod-vpc", "10.0.0.0/16", "vpc-12345678"},
				{"staging-vpc", "10.1.0.0/16", "vpc-87654321"},
				{"dev-vpc", "10.2.0.0/16", "vpc-11223344"},
			}
			for _, v := range vpcs {
				fmt.Printf("%-20s %-15s %-10s\n", v.name, v.cidr, v.id)
			}
			return nil
		})

	networkVPC.AddCommand(vpcCreate)
	networkVPC.AddCommand(vpcList)

	// network firewall
	networkFirewall := cli.Cmd("firewall").
		Description("Firewall rules management")

	fwAddRule := cli.Cmd("add-rule").
		Description("Add firewall rule").
		Arg("name", "Rule name", true).
		Flag(&protocol, "protocol", "p", "tcp", "Protocol (tcp/udp/icmp)").
		Flag(&port, "port", "", 80, "Port number").
		Flag(&sourceIP, "source", "s", "0.0.0.0/0", "Source IP/CIDR").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("🔥 Adding firewall rule: %s\n", name)
			fmt.Printf("   Protocol: %s\n", protocol)
			fmt.Printf("   Port: %d\n", port)
			fmt.Printf("   Source: %s\n", sourceIP)
			time.Sleep(200 * time.Millisecond)
			fmt.Println("✅ Firewall rule added")
			return nil
		})

	fwList := cli.Cmd("list").
		Description("List firewall rules").
		Action(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("🔥 Firewall Rules:")
			fmt.Println()
			fmt.Printf("%-20s %-10s %-8s %-20s\n", "NAME", "PROTOCOL", "PORT", "SOURCE")
			fmt.Println("----------------------------------------------------------------")
			rules := []struct {
				name     string
				protocol string
				port     string
				source   string
			}{
				{"allow-http", "tcp", "80", "0.0.0.0/0"},
				{"allow-https", "tcp", "443", "0.0.0.0/0"},
				{"allow-ssh", "tcp", "22", "10.0.0.0/8"},
			}
			for _, r := range rules {
				fmt.Printf("%-20s %-10s %-8s %-20s\n", r.name, r.protocol, r.port, r.source)
			}
			return nil
		})

	networkFirewall.AddCommand(fwAddRule)
	networkFirewall.AddCommand(fwList)

	// network lb (load balancer)
	networkLB := cli.Cmd("loadbalancer").
		Description("Load balancer management")

	lbCreate := cli.Cmd("create").
		Description("Create load balancer").
		Arg("name", "Load balancer name", true).
		Flag(&vpcID, "vpc", "", "", "VPC ID").
		Flag(&targetPort, "target-port", "", 80, "Target port").
		Action(func(ctx context.Context, cmd *cli.Command, name string) error {
			fmt.Printf("⚖️  Creating load balancer: %s\n", name)
			if vpcID != "" {
				fmt.Printf("   VPC: %s\n", vpcID)
			}
			fmt.Printf("   Target Port: %d\n", targetPort)
			time.Sleep(500 * time.Millisecond)
			fmt.Println("   Provisioning...")
			time.Sleep(400 * time.Millisecond)
			fmt.Println("   Configuring health checks...")
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("✅ Load balancer %s created\n", name)
			fmt.Printf("   DNS: %s-123456.elb.amazonaws.com\n", name)
			return nil
		})

	lbList := cli.Cmd("list").
		Description("List load balancers").
		Action(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("⚖️  Load Balancers:")
			fmt.Println()
			fmt.Printf("%-20s %-10s %-15s\n", "NAME", "STATUS", "TARGETS")
			fmt.Println("--------------------------------------------------")
			lbs := []struct {
				name    string
				status  string
				targets string
			}{
				{"prod-lb", "active", "3/3 healthy"},
				{"staging-lb", "active", "2/2 healthy"},
				{"api-lb", "active", "5/6 healthy"},
			}
			for _, lb := range lbs {
				fmt.Printf("%-20s %-10s %-15s\n", lb.name, lb.status, lb.targets)
			}
			return nil
		})

	networkLB.AddCommand(lbCreate)
	networkLB.AddCommand(lbList)

	network.AddCommand(networkVPC)
	network.AddCommand(networkFirewall)
	network.AddCommand(networkLB)

	return network
}
