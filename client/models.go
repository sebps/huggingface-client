package client

import "time"

type Endpoint struct {
	Name                 string                  `json:"name"`
	Type                 EndpointType            `json:"type"`
	Provider             EndpointProvider        `json:"provider"`
	Compute              EndpointCompute         `json:"compute"`
	Model                EndpointModel           `json:"model"`
	Tags                 []string                `json:"tags,omitempty"`
	CacheHttpResponses   *bool                   `json:"cacheHttpResponses,omitempty"`
	ExperimentalFeatures *ExperimentalFeatures   `json:"experimentalFeatures,omitempty"`
	PrivateService       *EndpointPrivateService `json:"privateService,omitempty"`
	Route                *RouteSpec              `json:"route,omitempty"`
}

type EndpointWithStatus struct {
	Name                 string                  `json:"name"`
	Type                 EndpointType            `json:"type"`
	Provider             EndpointProvider        `json:"provider"`
	Compute              EndpointCompute         `json:"compute"`
	Model                EndpointModel           `json:"model"`
	Tags                 []string                `json:"tags"`
	CacheHttpResponses   *bool                   `json:"cacheHttpResponses,omitempty"`
	ExperimentalFeatures *ExperimentalFeatures   `json:"experimentalFeatures"`
	PrivateService       *EndpointPrivateService `json:"privateService,omitempty"`
	Route                *RouteSpec              `json:"route,omitempty"`
	Status               EndpointStatus          `json:"status"`
}

type EndpointUpdate struct {
	Compute              *EndpointComputeUpdate `json:"compute,omitempty"`
	Model                *EndpointModelUpdate   `json:"model,omitempty"`
	ExperimentalFeatures *ExperimentalFeatures  `json:"experimentalFeatures,omitempty"`
	Route                *RouteSpec             `json:"route,omitempty"`
	Tags                 []string               `json:"tags,omitempty"`
	Type                 *EndpointType          `json:"type,omitempty"`
}

type EndpointProvider struct {
	Vendor string `json:"vendor"`
	Region string `json:"region"`
}

type EndpointCompute struct {
	Accelerator  AcceleratorType `json:"accelerator"`
	ID           *string         `json:"id,omitempty"`
	InstanceType string          `json:"instanceType"`
	InstanceSize string          `json:"instanceSize"`
	Scaling      EndpointScaling `json:"scaling"`
}

type EndpointComputeUpdate struct {
	Accelerator  *AcceleratorType       `json:"accelerator,omitempty"`
	InstanceType *string                `json:"instanceType,omitempty"`
	InstanceSize *string                `json:"instanceSize,omitempty"`
	Scaling      *EndpointScalingUpdate `json:"scaling,omitempty"`
}

type EndpointScaling struct {
	MinReplica         int             `json:"minReplica"`
	MaxReplica         int             `json:"maxReplica"`
	Measure            *ScalingMeasure `json:"measure,omitempty"`
	Metric             *ScalingMetric  `json:"metric,omitempty"`
	ScaleToZeroTimeout *int            `json:"scaleToZeroTimeout,omitempty"`
	Threshold          *float64        `json:"threshold,omitempty"`
}

type EndpointScalingUpdate struct {
	MinReplica         *int            `json:"minReplica,omitempty"`
	MaxReplica         *int            `json:"maxReplica,omitempty"`
	Measure            *ScalingMeasure `json:"measure,omitempty"`
	Metric             *ScalingMetric  `json:"metric,omitempty"`
	ScaleToZeroTimeout *int            `json:"scaleToZeroTimeout,omitempty"`
	Threshold          *float64        `json:"threshold,omitempty"`
}

type ScalingMeasure struct {
	HardwareUsage   *float64 `json:"hardwareUsage,omitempty"`
	PendingRequests *float64 `json:"pendingRequests,omitempty"`
}

type ScalingMetric string

const (
	ScalingMetricHardwareUsage   ScalingMetric = "hardwareUsage"
	ScalingMetricPendingRequests ScalingMetric = "pendingRequests"
)

type AcceleratorType string

const (
	AcceleratorCPU    AcceleratorType = "cpu"
	AcceleratorGPU    AcceleratorType = "gpu"
	AcceleratorNeuron AcceleratorType = "neuron"
)

type EndpointModel struct {
	Repository  string             `json:"repository"`
	Framework   EndpointFramework  `json:"framework"`
	Image       EndpointModelImage `json:"image"`
	Args        []string           `json:"args,omitempty"`
	Command     []string           `json:"command,omitempty"`
	Env         map[string]string  `json:"env,omitempty"`
	Secrets     map[string]*string `json:"secrets,omitempty"`
	FromCatalog *bool              `json:"fromCatalog,omitempty"`
	Revision    *string            `json:"revision,omitempty"`
	Task        EndpointTask       `json:"task"`
}

type EndpointModelUpdate struct {
	Repository *string             `json:"repository,omitempty"`
	Framework  *EndpointFramework  `json:"framework,omitempty"`
	Image      *EndpointModelImage `json:"image,omitempty"`
	Args       []string            `json:"args,omitempty"`
	Command    []string            `json:"command,omitempty"`
	Env        map[string]string   `json:"env,omitempty"`
	Secrets    map[string]*string  `json:"secrets,omitempty"`
	Revision   *string             `json:"revision,omitempty"`
	Task       *EndpointTask       `json:"task,omitempty"`
}

type EndpointFramework string

const (
	FrameworkCustom   EndpointFramework = "custom"
	FrameworkPytorch  EndpointFramework = "pytorch"
	FrameworkLlamaCpp EndpointFramework = "llamacpp"
)

type EndpointTask string

type EndpointModelImage struct {
	HuggingFace       *HuggingFaceImage       `json:"huggingface,omitempty"`
	HuggingFaceNeuron *HuggingFaceNeuronImage `json:"huggingfaceNeuron,omitempty"`
	TGI               *TGIImage               `json:"tgi,omitempty"`
	TGINeuron         *TGINeuronImage         `json:"tgiNeuron,omitempty"`
	TEI               *TEIImage               `json:"tei,omitempty"`
	LlamaCpp          *LlamaCppImage          `json:"llamacpp,omitempty"`
	Custom            *CustomImage            `json:"custom,omitempty"`
}

type HuggingFaceImage struct{}

type HuggingFaceNeuronImage struct {
	BatchSize      *int   `json:"batchSize,omitempty"`
	NeuronCache    string `json:"neuronCache"`
	SequenceLength *int   `json:"sequenceLength,omitempty"`
}

type TGIImage struct {
	HealthRoute           *string       `json:"healthRoute,omitempty"`
	Port                  int           `json:"port"`
	URL                   string        `json:"url"`
	MaxBatchPrefillTokens *int          `json:"maxBatchPrefillTokens,omitempty"`
	MaxBatchTotalTokens   *int          `json:"maxBatchTotalTokens,omitempty"`
	MaxInputLength        *int          `json:"maxInputLength,omitempty"`
	MaxTotalTokens        *int          `json:"maxTotalTokens,omitempty"`
	DisableCustomKernels  bool          `json:"disableCustomKernels"`
	Quantize              *QuantizeType `json:"quantize,omitempty"`
}

type TGINeuronImage struct {
	HealthRoute           *string       `json:"healthRoute,omitempty"`
	Port                  int           `json:"port"`
	URL                   string        `json:"url"`
	MaxBatchPrefillTokens *int          `json:"maxBatchPrefillTokens,omitempty"`
	MaxBatchTotalTokens   *int          `json:"maxBatchTotalTokens,omitempty"`
	MaxInputLength        *int          `json:"maxInputLength,omitempty"`
	MaxTotalTokens        *int          `json:"maxTotalTokens,omitempty"`
	HfAutoCastType        *AutoCastType `json:"hfAutoCastType,omitempty"`
	HfNumCores            *int          `json:"hfNumCores,omitempty"`
}

type TEIImage struct {
	HealthRoute           *string      `json:"healthRoute,omitempty"`
	Port                  int          `json:"port"`
	URL                   string       `json:"url"`
	MaxBatchTokens        *int         `json:"maxBatchTokens,omitempty"`
	MaxConcurrentRequests *int         `json:"maxConcurrentRequests,omitempty"`
	Pooling               *PoolingType `json:"pooling,omitempty"`
}

type LlamaCppImage struct {
	HealthRoute *string      `json:"healthRoute,omitempty"`
	Port        int          `json:"port"`
	URL         string       `json:"url"`
	CtxSize     int          `json:"ctxSize"`
	Mode        *ModelMode   `json:"mode,omitempty"`
	ModelPath   string       `json:"modelPath"`
	NGpuLayers  int          `json:"nGpuLayers"`
	NParallel   int          `json:"nParallel"`
	Pooling     *PoolingType `json:"pooling,omitempty"`
	ThreadsHttp *int         `json:"threadsHttp,omitempty"`
	Variant     *string      `json:"variant,omitempty"`
}

type Credentials struct {
	Username string  `json:"username"`
	Password *string `json:"password,omitempty"`
}

type QuantizeType string

const (
	QuantizeAWQ          QuantizeType = "awq"
	QuantizeBitsAndBytes QuantizeType = "bitsandbytes"
	QuantizeEETQ         QuantizeType = "eetq"
	QuantizeGPTQ         QuantizeType = "gptq"
)

type AutoCastType string

const (
	AutoCastBF16 AutoCastType = "bf16"
	AutoCastFP16 AutoCastType = "fp16"
)

type PoolingType string

const (
	PoolingMean PoolingType = "mean"
	PoolingCLS  PoolingType = "cls"
	PoolingLast PoolingType = "last"
	PoolingRank PoolingType = "rank"
)

type ModelMode string

const (
	ModelModeEmbeddings ModelMode = "embeddings"
	ModelModeReranking  ModelMode = "reranking"
)

type CustomImage struct {
	URL         string       `json:"url"`
	HealthRoute *string      `json:"healthRoute,omitempty"`
	Port        int          `json:"port"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

type EndpointPrivateService struct {
	AccountID string `json:"accountId"`
	Shared    bool   `json:"shared"`
}

type RouteSpec struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type EndpointStatus struct {
	CreatedAt     time.Time       `json:"createdAt"`
	CreatedBy     EndpointAccount `json:"createdBy"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	UpdatedBy     EndpointAccount `json:"updatedBy"`
	State         EndpointState   `json:"state"`
	Message       string          `json:"message"`
	ReadyReplica  int             `json:"readyReplica"`
	TargetReplica int             `json:"targetReplica"`
	ErrorMessage  *string         `json:"errorMessage,omitempty"`
	URL           *string         `json:"url,omitempty"`
	Private       *PrivateStatus  `json:"private,omitempty"`
}

type EndpointAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PrivateStatus struct {
	ServiceName *string `json:"serviceName,omitempty"`
}

type EndpointState string

const (
	StatePending      EndpointState = "pending"
	StateInitializing EndpointState = "initializing"
	StateUpdating     EndpointState = "updating"
	StateUpdateFailed EndpointState = "updateFailed"
	StateRunning      EndpointState = "running"
	StatePaused       EndpointState = "paused"
	StateFailed       EndpointState = "failed"
	StateScaledToZero EndpointState = "scaledToZero"
)

type EndpointType string

const (
	TypePublic    EndpointType = "public"
	TypeVerified  EndpointType = "verified"
	TypeProtected EndpointType = "protected"
	TypePrivate   EndpointType = "private"
)

type ExperimentalFeatures struct {
	CacheHttpResponses bool      `json:"cacheHttpResponses"`
	KvRouter           *KvRouter `json:"kvRouter,omitempty"`
}

type KvRouter struct {
	Tag string `json:"tag"`
}
