package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sebps/huggingface-client/client"
	"github.com/sebps/huggingface-client/utils"
	"github.com/spf13/cobra"
)

var (
	host                  string
	token                 string
	namespace             string
	inferenceName         string
	inferenceRepository   string
	inferenceFramework    string
	inferenceImage        string
	inferenceImageUrl     string
	inferenceModelPath    string
	inferencePort         int
	inferenceUsername     string
	inferencePassword     string
	inferenceTask         string
	inferenceAccelerator  string
	inferenceVendor       string
	inferenceRegion       string
	inferenceType         string
	inferenceInstanceSize string
	inferenceInstanceType string
	inferenceMinReplica   int
	inferenceMaxReplica   int
	logsReplicaID         string
	startTime             string
	stopTime              string
	metricStep            string
)

func init() {
	endpointCmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Manage endpoints",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			endpoints, err := c.ListEndpoints(namespace, nil)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			printJSON(endpoints)

			return nil
		},
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			inferenceImage, err := utils.BuildInferenceImage(
				inferenceImage,
				inferenceImageUrl,
				inferencePort,
				inferenceModelPath,
				inferenceUsername,
				inferencePassword,
			)
			if err != nil {
				return err
			}

			endpoint := client.Endpoint{
				Name: inferenceName,
				Type: client.EndpointType(inferenceType),
				Provider: client.EndpointProvider{
					Vendor: inferenceVendor,
					Region: inferenceRegion,
				},
				Compute: client.EndpointCompute{
					Accelerator:  client.AcceleratorType(inferenceAccelerator),
					InstanceType: inferenceInstanceType,
					InstanceSize: inferenceInstanceSize,
					Scaling: client.EndpointScaling{
						MinReplica: inferenceMinReplica,
						MaxReplica: inferenceMaxReplica,
					},
				},
				Model: client.EndpointModel{
					Repository: inferenceRepository,
					Framework:  client.EndpointFramework(inferenceFramework),
					Image:      *inferenceImage,
					Task:       client.EndpointTask(inferenceTask),
				},
			}

			createdEndpoint, err := c.CreateEndpoint(namespace, endpoint)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			printJSON(createdEndpoint)

			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Get an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			endpoint, err := c.GetEndpoint(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			printJSON(endpoint)

			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update [name]",
		Short: "Update an existing endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			endpointUpdate := client.EndpointUpdate{}

			// Update Type if provided
			if cmd.Flags().Changed("type") {
				t := client.EndpointType(inferenceType)
				endpointUpdate.Type = &t
			}

			// Update Compute
			compute := &client.EndpointComputeUpdate{}
			anyComputeField := false

			if cmd.Flags().Changed("accelerator") {
				accel := client.AcceleratorType(inferenceAccelerator)
				compute.Accelerator = &accel
				anyComputeField = true
			}
			if cmd.Flags().Changed("instance-type") {
				compute.InstanceType = &inferenceInstanceType
				anyComputeField = true
			}
			if cmd.Flags().Changed("instance-size") {
				compute.InstanceSize = &inferenceInstanceSize
				anyComputeField = true
			}
			if cmd.Flags().Changed("min-replica") || cmd.Flags().Changed("max-replica") {
				scaling := &client.EndpointScalingUpdate{}
				if cmd.Flags().Changed("min-replica") {
					scaling.MinReplica = &inferenceMinReplica
				}
				if cmd.Flags().Changed("max-replica") {
					scaling.MaxReplica = &inferenceMaxReplica
				}
				compute.Scaling = scaling
				anyComputeField = true
			}
			if anyComputeField {
				endpointUpdate.Compute = compute
			}

			// Update Model
			model := &client.EndpointModelUpdate{}
			anyModelField := false

			if cmd.Flags().Changed("repository") {
				model.Repository = &inferenceRepository
				anyModelField = true
			}
			if cmd.Flags().Changed("framework") {
				framework := client.EndpointFramework(inferenceFramework)
				model.Framework = &framework
				anyModelField = true
			}
			if cmd.Flags().Changed("task") {
				task := client.EndpointTask(inferenceTask)
				model.Task = &task
				anyModelField = true
			}
			if cmd.Flags().Changed("image") || cmd.Flags().Changed("url") || cmd.Flags().Changed("path") {
				inferenceImageStruct, err := utils.BuildInferenceImage(
					inferenceImage,
					inferenceImageUrl,
					inferencePort,
					inferenceModelPath,
					inferenceUsername,
					inferencePassword,
				)
				if err != nil {
					return err
				}
				model.Image = inferenceImageStruct
				anyModelField = true
			}

			if anyModelField {
				endpointUpdate.Model = model
			}

			// Nothing was updated
			if endpointUpdate.Compute == nil && endpointUpdate.Model == nil && endpointUpdate.Type == nil {
				return fmt.Errorf("no update flags were provided, nothing to update")
			}

			updatedEndpoint, err := c.UpdateEndpoint(namespace, args[0], endpointUpdate)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			printJSON(updatedEndpoint)

			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			err = c.DeleteEndpoint(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Endpoint deleted successfully.")

			return nil
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [name]",
		Short: "Get logs from an endpoint (optionally filtered by replica)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			var replicaParam *string
			if logsReplicaID != "" {
				replicaParam = &logsReplicaID
			}

			logs, err := c.GetEndpointLogs(namespace, args[0], replicaParam)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println(string(logs))

			return nil
		},
	}

	logsStreamCmd := &cobra.Command{
		Use:   "logs-stream [name]",
		Short: "Stream live logs from an endpoint (optionally filtered by replica)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			var replicaParam *string
			if logsReplicaID != "" {
				replicaParam = &logsReplicaID
			}

			stream, err := c.StreamEndpointLogs(namespace, args[0], replicaParam)
			if err != nil {
				fmt.Println(err)
				stream.Close()
				return nil
			}
			defer stream.Close()

			buf := make([]byte, 4096)
			for {
				n, err := stream.Read(buf)
				if err != nil {
					break
				}
				fmt.Print(string(buf[:n]))
			}

			return nil
		},
	}

	getMetricsCmd := &cobra.Command{
		Use:   "metrics [name]",
		Short: "List endpoint metrics (hardwareUsage or pendingRequests)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if startTime == "" || stopTime == "" {
				return fmt.Errorf("both --start and --stop must be specified")
			}

			startTime, err := utils.ParseTime(startTime)
			if err != nil {
				return fmt.Errorf("invalid --start time format : %w", err)
			}

			stopTime, err := utils.ParseTime(stopTime)
			if err != nil {
				return fmt.Errorf("invalid --stop time format : %w", err)
			}

			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			payload := client.MetricsRequest{
				Start: startTime,
				Stop:  stopTime,
			}

			metricsData, err := c.GetEndpointMetrics(namespace, args[0], payload)
			if err != nil {
				return err
			}

			fmt.Println(string(metricsData))

			printJSON(metricsData)

			return nil
		},
	}

	getMetricCmd := &cobra.Command{
		Use:   "metric [name] [metric]",
		Short: "Get endpoint specific metric, metric argument can be one of the following values : `pending-requests`, `request-count`, `median-latency`, `p95-latency`, `success-throughput`, `bad-request-throughput`, `server-error-throughput`, `cpu-usage`, `memory-usage`, `gpu-usage`, `gpu-memory-usage`, `neuron-usage`, `neuron-memory-usage`, `ready-replicas`, `running-replicas`, `target-replicas`, `average-latency`, `success-rate`, `bad-request-rate`, `server-error-rate`",
		RunE: func(cmd *cobra.Command, args []string) error {
			if startTime == "" || stopTime == "" {
				return fmt.Errorf("both --start and --stop must be specified")
			}

			startTime, err := utils.ParseTime(startTime)
			if err != nil {
				return fmt.Errorf("invalid --start time format : %w", err)
			}

			stopTime, err := utils.ParseTime(stopTime)
			if err != nil {
				return fmt.Errorf("invalid --stop time format : %w", err)
			}

			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			payload := client.MetricRequest{
				From: uint32(startTime.Unix()),
				To:   uint32(stopTime.Unix()),
			}

			if metricStep != "" {
				payload.Step = &metricStep
			}

			if !utils.IsMetricValid(args[1]) {
				return errors.New("invalid metric. metric needs to be one of the following values : `pending-requests`, `request-count`, `median-latency`, `p95-latency`, `success-throughput`, `bad-request-throughput`, `server-error-throughput`, `cpu-usage`, `memory-usage`, `gpu-usage`, `gpu-memory-usage`, `neuron-usage`, `neuron-memory-usage`, `ready-replicas`, `running-replicas`, `target-replicas`, `average-latency`, `success-rate`, `bad-request-rate`, `server-error-rate`")
			}

			metricData, err := c.GetEndpointMetric(namespace, args[0], args[1], payload)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println(string(metricData))

			return nil
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause [name]",
		Short: "Pause an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			err = c.PauseEndpoint(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Endpoint paused successfully.")

			return nil
		},
	}

	getReplicasStatusesCmd := &cobra.Command{
		Use:   "replica [name]",
		Short: "Get endpoint replica statuses",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			replicas, err := c.GetEndpointReplicasStatuses(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println(string(replicas))

			return nil
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume [name]",
		Short: "Resume an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			err = c.ResumeEndpoint(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Endpoint resumed successfully.")

			return nil
		},
	}

	scaleToZeroCmd := &cobra.Command{
		Use:   "scale-to-zero [name]",
		Short: "Scale an endpoint to zero",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			err = c.ScaleEndpointToZero(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Endpoint scaled to zero successfully.")

			return nil
		},
	}

	sseCmd := &cobra.Command{
		Use:   "sse [name]",
		Short: "Stream SSE info from an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewClient(&host, &token)
			if err != nil {
				return err
			}

			stream, err := c.GetEndpointSSE(namespace, args[0])
			if err != nil {
				fmt.Println(err)
				stream.Close()
				return nil
			}
			defer stream.Close()

			buf := make([]byte, 4096)
			for {
				n, err := stream.Read(buf)
				if err != nil {
					break
				}
				fmt.Print(string(buf[:n]))
			}

			return nil
		},
	}

	createCmd.Flags().StringVar(&inferenceName, "name", "", "Endpoint name (required)")
	createCmd.Flags().StringVar(&inferenceRepository, "repository", "", "Model repository (required)")
	createCmd.Flags().StringVar(&inferenceFramework, "framework", "", "Model framework (pytorch, custom, llamacpp)")
	createCmd.Flags().StringVar(&inferenceImage, "image", "huggingface", "Model image (huggingface, huggingfaceNeuron, tgi, tgiNeuron, tei, llamacpp, custom)")
	createCmd.Flags().StringVar(&inferenceImageUrl, "url", "", "Model image url (required for tgi, tgiNeuron, tei, llamacpp, custom) using format https://host/image:tag")
	createCmd.Flags().StringVar(&inferenceTask, "task", "", "Model task (text-generation, classification, etc.)")
	createCmd.Flags().IntVar(&inferencePort, "port", 80, "Endpoint API port")
	createCmd.Flags().StringVar(&inferenceUsername, "username", "", "Endpoint API username ( basic auth )")
	createCmd.Flags().StringVar(&inferencePassword, "password", "", "Endpoint API password ( basic auth )")
	createCmd.Flags().StringVar(&inferenceModelPath, "path", "", "Path to the .gguf file to be loaded")
	createCmd.Flags().StringVar(&inferenceAccelerator, "accelerator", "cpu", "Accelerator type (cpu, gpu, neuron)")
	createCmd.Flags().StringVar(&inferenceVendor, "vendor", "", "Cloud vendor (aws, azure, gcp)")
	createCmd.Flags().StringVar(&inferenceRegion, "region", "", "Cloud region")
	createCmd.Flags().StringVar(&inferenceType, "type", "protected", "Endpoint type (public, protected, private)")
	createCmd.Flags().StringVar(&inferenceInstanceSize, "instance-size", "small", "Instance size")
	createCmd.Flags().StringVar(&inferenceInstanceType, "instance-type", "default", "Instance type")
	createCmd.Flags().IntVar(&inferenceMinReplica, "min-replica", 0, "Endpoint minimum replica")
	createCmd.Flags().IntVar(&inferenceMaxReplica, "max-replica", 0, "Endpoint maximum replica")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("repository")
	createCmd.MarkFlagRequired("framework")
	createCmd.MarkFlagRequired("task")
	createCmd.MarkFlagRequired("vendor")
	createCmd.MarkFlagRequired("region")
	createCmd.MarkFlagRequired("image")

	updateCmd.Flags().StringVar(&inferenceRepository, "repository", "", "Model repository")
	updateCmd.Flags().StringVar(&inferenceFramework, "framework", "", "Model framework (pytorch, custom, llamacpp)")
	updateCmd.Flags().StringVar(&inferenceImage, "image", "huggingface", "Model image (huggingface, huggingfaceNeuron, tgi, tgiNeuron, tei, llamacpp, custom)")
	updateCmd.Flags().StringVar(&inferenceImageUrl, "url", "", "Model image url (for tgi, tgiNeuron, tei, llamacpp, custom) format https://host/image:tag")
	updateCmd.Flags().StringVar(&inferenceTask, "task", "", "Model task (text-generation, classification, etc.)")
	updateCmd.Flags().IntVar(&inferencePort, "port", 80, "Endpoint API port")
	updateCmd.Flags().StringVar(&inferenceUsername, "username", "", "Endpoint API username (basic auth)")
	updateCmd.Flags().StringVar(&inferencePassword, "password", "", "Endpoint API password (basic auth)")
	updateCmd.Flags().StringVar(&inferenceModelPath, "path", "", "Path to the .gguf file to be loaded")
	updateCmd.Flags().StringVar(&inferenceAccelerator, "accelerator", "cpu", "Accelerator type (cpu, gpu, neuron)")
	updateCmd.Flags().StringVar(&inferenceVendor, "vendor", "", "Cloud vendor (aws, azure, gcp)")
	updateCmd.Flags().StringVar(&inferenceRegion, "region", "", "Cloud region")
	updateCmd.Flags().StringVar(&inferenceType, "type", "protected", "Endpoint type (public, protected, private)")
	updateCmd.Flags().StringVar(&inferenceInstanceSize, "instance-size", "small", "Instance size")
	updateCmd.Flags().StringVar(&inferenceInstanceType, "instance-type", "default", "Instance type")
	updateCmd.Flags().IntVar(&inferenceMinReplica, "min-replica", 0, "Endpoint minimum replica")
	updateCmd.Flags().IntVar(&inferenceMaxReplica, "max-replica", 0, "Endpoint maximum replica")

	logsCmd.Flags().StringVar(&logsReplicaID, "replica", "", "Replica ID to filter logs (optional)")

	logsStreamCmd.Flags().StringVar(&logsReplicaID, "replica", "", "Replica ID to filter logs stream (optional)")

	getMetricsCmd.Flags().StringVar(&startTime, "start", "", "Metrics measurement start")
	getMetricsCmd.Flags().StringVar(&stopTime, "stop", "", "Metrics measurement stop")

	getMetricsCmd.MarkFlagRequired("start")
	getMetricsCmd.MarkFlagRequired("stop")

	getMetricCmd.Flags().StringVar(&startTime, "start", "", "Metric measurement start")
	getMetricCmd.Flags().StringVar(&stopTime, "stop", "", "Metric measurement stop")
	getMetricCmd.Flags().StringVar(&metricStep, "step", "", "Duration step ( '1m','5m', etc... ) )")

	getMetricCmd.MarkFlagRequired("start")
	getMetricCmd.MarkFlagRequired("stop")

	endpointCmd.PersistentFlags().StringVar(&token, "token", "", "Authorization Bearer token (required)")
	endpointCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace (organization or user) (required)")
	endpointCmd.PersistentFlags().StringVar(&host, "host", "", "API host URL (optional)")

	endpointCmd.MarkPersistentFlagRequired("token")
	endpointCmd.MarkPersistentFlagRequired("namespace")

	endpointCmd.AddCommand(listCmd)
	endpointCmd.AddCommand(createCmd)
	endpointCmd.AddCommand(getCmd)
	endpointCmd.AddCommand(updateCmd)
	endpointCmd.AddCommand(deleteCmd)
	endpointCmd.AddCommand(logsCmd)
	endpointCmd.AddCommand(logsStreamCmd)
	endpointCmd.AddCommand(getMetricsCmd)
	endpointCmd.AddCommand(getMetricCmd)
	endpointCmd.AddCommand(pauseCmd)
	endpointCmd.AddCommand(getReplicasStatusesCmd)
	endpointCmd.AddCommand(resumeCmd)

	endpointCmd.AddCommand(scaleToZeroCmd)
	endpointCmd.AddCommand(sseCmd)

	rootCmd.AddCommand(endpointCmd)
}

func printJSON(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(string(dataBytes))

	return nil
}
