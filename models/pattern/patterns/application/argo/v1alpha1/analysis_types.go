package v1alpha1

import (
	"time"

	intstrutil "k8s.io/apimachinery/pkg/util/intstr"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterAnalysisTemplate holds the template for performing canary analysis
type ClusterAnalysisTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec AnalysisTemplateSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// AnalysisTemplateList is a list of AnalysisTemplate resources
type ClusterAnalysisTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterAnalysisTemplate `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AnalysisTemplate holds the template for performing canary analysis
type AnalysisTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec AnalysisTemplateSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// AnalysisTemplateList is a list of AnalysisTemplate resources
type AnalysisTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []AnalysisTemplate `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AnalysisTemplateSpec is the specification for a AnalysisTemplate resource
type AnalysisTemplateSpec struct {
	// Metrics contains the list of metrics to query as part of an analysis run
	Metrics []Metric `json:"metrics" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=metrics"`
	// Args are the list of arguments to the template
	Args []Argument `json:"args,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=args"`
}

// DurationString is a string representing a duration (e.g. 30s, 5m, 1h)
type DurationString string

// Duration converts DurationString into a time.Duration
func (d DurationString) Duration() (time.Duration, error) {
	return time.ParseDuration(string(d))
}

// Metric defines a metric in which to perform analysis
type Metric struct {
	// Name is the name of the metric
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Interval defines an interval string (e.g. 30s, 5m, 1h) between each measurement.
	// If omitted, will perform a single measurement
	Interval DurationString `json:"interval,omitempty" protobuf:"bytes,2,opt,name=interval,casttype=DurationString"`
	// InitialDelay how long the AnalysisRun should wait before starting this metric
	InitialDelay DurationString `json:"initialDelay,omitempty" protobuf:"bytes,3,opt,name=initialDelay,casttype=DurationString"`
	// Count is the number of times to run the measurement. If both interval and count are omitted,
	// the effective count is 1. If only interval is specified, metric runs indefinitely.
	// If count > 1, interval must be specified.
	Count *intstrutil.IntOrString `json:"count,omitempty" protobuf:"bytes,4,opt,name=count"`
	// SuccessCondition is an expression which determines if a measurement is considered successful
	// Expression is a goevaluate expression. The keyword `result` is a variable reference to the
	// value of measurement. Results can be both structured data or primitive.
	// Examples:
	//   result > 10
	//   (result.requests_made * result.requests_succeeded / 100) >= 90
	SuccessCondition string `json:"successCondition,omitempty" protobuf:"bytes,5,opt,name=successCondition"`
	// FailureCondition is an expression which determines if a measurement is considered failed
	// If both success and failure conditions are specified, and the measurement does not fall into
	// either condition, the measurement is considered Inconclusive
	FailureCondition string `json:"failureCondition,omitempty" protobuf:"bytes,6,opt,name=failureCondition"`
	// FailureLimit is the maximum number of times the measurement is allowed to fail, before the
	// entire metric is considered Failed (default: 0)
	FailureLimit *intstrutil.IntOrString `json:"failureLimit,omitempty" protobuf:"bytes,7,opt,name=failureLimit"`
	// InconclusiveLimit is the maximum number of times the measurement is allowed to measure
	// Inconclusive, before the entire metric is considered Inconclusive (default: 0)
	InconclusiveLimit *intstrutil.IntOrString `json:"inconclusiveLimit,omitempty" protobuf:"bytes,8,opt,name=inconclusiveLimit"`
	// ConsecutiveErrorLimit is the maximum number of times the measurement is allowed to error in
	// succession, before the metric is considered error (default: 4)
	ConsecutiveErrorLimit *intstrutil.IntOrString `json:"consecutiveErrorLimit,omitempty" protobuf:"bytes,9,opt,name=consecutiveErrorLimit"`
	// Provider configuration to the external system to use to verify the analysis
	Provider MetricProvider `json:"provider" protobuf:"bytes,10,opt,name=provider"`
}

// EffectiveCount is the effective count based on whether or not count/interval is specified
// If neither count or interval is specified, the effective count is 1
// If only interval is specified, metric runs indefinitely and there is no effective count (nil)
// Otherwise, it is the user specified value
func (m *Metric) EffectiveCount() *intstrutil.IntOrString {
	// Need to check if type is String
	if m.Count == nil || m.Count.IntValue() == 0 {
		if m.Interval == "" {
			one := intstrutil.FromInt(1)
			return &one
		}
		return nil
	}
	return m.Count
}

// MetricProvider which external system to use to verify the analysis
// Only one of the fields in this struct should be non-nil
type MetricProvider struct {
	// Prometheus specifies the prometheus metric to query
	Prometheus *PrometheusMetric `json:"prometheus,omitempty" protobuf:"bytes,1,opt,name=prometheus"`
	// Kayenta specifies a Kayenta metric
	Kayenta *KayentaMetric `json:"kayenta,omitempty" protobuf:"bytes,2,opt,name=kayenta"`
	// Web specifies a generic HTTP web metric
	Web *WebMetric `json:"web,omitempty" protobuf:"bytes,3,opt,name=web"`
	// Datadog specifies a datadog metric to query
	Datadog *DatadogMetric `json:"datadog,omitempty" protobuf:"bytes,4,opt,name=datadog"`
	// Wavefront specifies the wavefront metric to query
	Wavefront *WavefrontMetric `json:"wavefront,omitempty" protobuf:"bytes,5,opt,name=wavefront"`
	// NewRelic specifies the newrelic metric to query
	NewRelic *NewRelicMetric `json:"newRelic,omitempty" protobuf:"bytes,6,opt,name=newRelic"`
	// Job specifies the job metric run
	Job *JobMetric `json:"job,omitempty" protobuf:"bytes,7,opt,name=job"`
}

// AnalysisPhase is the overall phase of an AnalysisRun, MetricResult, or Measurement
type AnalysisPhase string

// Possible AnalysisPhase values
const (
	AnalysisPhasePending      AnalysisPhase = "Pending"
	AnalysisPhaseRunning      AnalysisPhase = "Running"
	AnalysisPhaseSuccessful   AnalysisPhase = "Successful"
	AnalysisPhaseFailed       AnalysisPhase = "Failed"
	AnalysisPhaseError        AnalysisPhase = "Error"
	AnalysisPhaseInconclusive AnalysisPhase = "Inconclusive"
)

// Completed returns whether or not the analysis status is considered completed
func (as AnalysisPhase) Completed() bool {
	switch as {
	case AnalysisPhaseSuccessful, AnalysisPhaseFailed, AnalysisPhaseError, AnalysisPhaseInconclusive:
		return true
	}
	return false
}

// PrometheusMetric defines the prometheus query to perform canary analysis
type PrometheusMetric struct {
	// Address is the HTTP address and port of the prometheus server
	Address string `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`
	// Query is a raw prometheus query to perform
	Query string `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
}

// WavefrontMetric defines the wavefront query to perform canary analysis
type WavefrontMetric struct {
	// Address is the HTTP address and port of the wavefront server
	Address string `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`
	// Query is a raw wavefront query to perform
	Query string `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
}

// NewRelicMetric defines the newrelic query to perform canary analysis
type NewRelicMetric struct {
	// Profile is the name of the secret holding NR account configuration
	Profile string `json:"profile,omitempty" protobuf:"bytes,1,opt,name=profile"`
	// Query is a raw newrelic NRQL query to perform
	Query string `json:"query" protobuf:"bytes,2,opt,name=query"`
}

// JobMetric defines a job to run which acts as a metric
type JobMetric struct {
	Metadata metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec     batchv1.JobSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// AnalysisRun is an instantiation of an AnalysisTemplate
type AnalysisRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              AnalysisRunSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            AnalysisRunStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// AnalysisRunList is a list of AnalysisTemplate resources
type AnalysisRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []AnalysisRun `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AnalysisRunSpec is the spec for a AnalysisRun resource
type AnalysisRunSpec struct {
	// Metrics contains the list of metrics to query as part of an analysis run
	Metrics []Metric `json:"metrics" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=metrics"`
	// Args are the list of arguments used in this run
	Args []Argument `json:"args,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=args"`
	// Terminate is used to prematurely stop the run (e.g. rollout completed and analysis is no longer desired)
	Terminate bool `json:"terminate,omitempty" protobuf:"varint,3,opt,name=terminate"`
}

// Argument is an argument to an AnalysisRun
type Argument struct {
	// Name is the name of the argument
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Value is the value of the argument
	Value *string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// ValueFrom is a reference to where a secret is stored. This field is one of the fields with valueFrom
	ValueFrom *ValueFrom `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

type ValueFrom struct {
	// Secret is a reference to where a secret is stored. This field is one of the fields with valueFrom
	SecretKeyRef *SecretKeyRef `json:"secretKeyRef,omitempty" protobuf:"bytes,1,opt,name=secretKeyRef"`
	//FieldRef is a reference to the fields in metadata which we are referencing. This field is one of the fields with
	//valueFrom
	FieldRef *FieldRef `json:"fieldRef,omitempty" protobuf:"bytes,2,opt,name=fieldRef"`
}

type SecretKeyRef struct {
	// Name is the name of the secret
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Key is the key of the secret to select from.
	Key string `json:"key" protobuf:"bytes,2,opt,name=key"`
}

// AnalysisRunStatus is the status for a AnalysisRun resource
type AnalysisRunStatus struct {
	// Phase is the status of the analysis run
	Phase AnalysisPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=AnalysisPhase"`
	// Message is a message explaining current status
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// MetricResults contains the metrics collected during the run
	MetricResults []MetricResult `json:"metricResults,omitempty" protobuf:"bytes,3,rep,name=metricResults"`
	// StartedAt indicates when the analysisRun first started
	StartedAt *metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,4,opt,name=startedAt"`
}

// MetricResult contain a list of the most recent measurements for a single metric along with
// counters on how often the measurement
type MetricResult struct {
	// Name is the name of the metric
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Phase is the overall aggregate status of the metric
	Phase AnalysisPhase `json:"phase" protobuf:"bytes,2,opt,name=phase,casttype=AnalysisPhase"`
	// Measurements holds the most recent measurements collected for the metric
	Measurements []Measurement `json:"measurements,omitempty" protobuf:"bytes,3,rep,name=measurements"`
	// Message contains a message describing current condition (e.g. error messages)
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// Count is the number of times the metric was measured without Error
	// This is equal to the sum of Successful, Failed, Inconclusive
	Count int32 `json:"count,omitempty" protobuf:"varint,5,opt,name=count"`
	// Successful is the number of times the metric was measured Successful
	Successful int32 `json:"successful,omitempty" protobuf:"varint,6,opt,name=successful"`
	// Failed is the number of times the metric was measured Failed
	Failed int32 `json:"failed,omitempty" protobuf:"varint,7,opt,name=failed"`
	// Inconclusive is the number of times the metric was measured Inconclusive
	Inconclusive int32 `json:"inconclusive,omitempty" protobuf:"varint,8,opt,name=inconclusive"`
	// Error is the number of times an error was encountered during measurement
	Error int32 `json:"error,omitempty" protobuf:"varint,9,opt,name=error"`
	// ConsecutiveError is the number of times an error was encountered during measurement in succession
	// Resets to zero when non-errors are encountered
	ConsecutiveError int32 `json:"consecutiveError,omitempty" protobuf:"varint,10,opt,name=consecutiveError"`
}

// Measurement is a point in time result value of a single metric, and the time it was measured
type Measurement struct {
	// Phase is the status of this single measurement
	Phase AnalysisPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=AnalysisPhase"`
	// Message contains a message describing current condition (e.g. error messages)
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// StartedAt is the timestamp in which this measurement started to be measured
	StartedAt *metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,3,opt,name=startedAt"`
	// FinishedAt is the timestamp in which this measurement completed and value was collected
	FinishedAt *metav1.Time `json:"finishedAt,omitempty" protobuf:"bytes,4,opt,name=finishedAt"`
	// Value is the measured value of the metric
	Value string `json:"value,omitempty" protobuf:"bytes,5,opt,name=value"`
	// Metadata stores additional metadata about this metric result, used by the different providers
	// (e.g. kayenta run ID, job name)
	Metadata map[string]string `json:"metadata,omitempty" protobuf:"bytes,6,rep,name=metadata"`
	// ResumeAt is the  timestamp when the analysisRun should try to resume the measurement
	ResumeAt *metav1.Time `json:"resumeAt,omitempty" protobuf:"bytes,7,opt,name=resumeAt"`
}

type KayentaMetric struct {
	Address string `json:"address" protobuf:"bytes,1,opt,name=address"`

	Application string `json:"application" protobuf:"bytes,2,opt,name=application"`

	CanaryConfigName string `json:"canaryConfigName" protobuf:"bytes,3,opt,name=canaryConfigName"`

	MetricsAccountName       string `json:"metricsAccountName" protobuf:"bytes,4,opt,name=metricsAccountName"`
	ConfigurationAccountName string `json:"configurationAccountName" protobuf:"bytes,5,opt,name=configurationAccountName"`
	StorageAccountName       string `json:"storageAccountName" protobuf:"bytes,6,opt,name=storageAccountName"`

	Threshold KayentaThreshold `json:"threshold" protobuf:"bytes,7,opt,name=threshold"`

	Scopes []KayentaScope `json:"scopes" protobuf:"bytes,8,rep,name=scopes"`
}

type KayentaThreshold struct {
	Pass     int64 `json:"pass" protobuf:"varint,1,opt,name=pass"`
	Marginal int64 `json:"marginal" protobuf:"varint,2,opt,name=marginal"`
}

type KayentaScope struct {
	Name            string      `json:"name" protobuf:"bytes,1,opt,name=name"`
	ControlScope    ScopeDetail `json:"controlScope" protobuf:"bytes,2,opt,name=controlScope"`
	ExperimentScope ScopeDetail `json:"experimentScope" protobuf:"bytes,3,opt,name=experimentScope"`
}

type ScopeDetail struct {
	Scope  string `json:"scope" protobuf:"bytes,1,opt,name=scope"`
	Region string `json:"region" protobuf:"bytes,2,opt,name=region"`
	Step   int64  `json:"step" protobuf:"varint,3,opt,name=step"`
	Start  string `json:"start" protobuf:"bytes,4,opt,name=start"`
	End    string `json:"end" protobuf:"bytes,5,opt,name=end"`
}

type WebMetric struct {
	// URL is the address of the web metric
	URL string `json:"url" protobuf:"bytes,1,opt,name=url"`
	// Headers are optional HTTP headers to use in the request
	Headers []WebMetricHeader `json:"headers,omitempty" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,2,rep,name=headers"`
	// TimeoutSeconds is the timeout for the request in seconds (default: 10)
	TimeoutSeconds int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,3,opt,name=timeoutSeconds"`
	// JSONPath is a JSON Path to use as the result variable (default: "{$}")
	JSONPath string `json:"jsonPath,omitempty" protobuf:"bytes,4,opt,name=jsonPath"`
	// Insecure skips host TLS verification
	Insecure bool `json:"insecure,omitempty" protobuf:"varint,5,opt,name=insecure"`
}

type WebMetricHeader struct {
	Key   string `json:"key" protobuf:"bytes,1,opt,name=key"`
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

type DatadogMetric struct {
	Interval DurationString `json:"interval,omitempty" protobuf:"bytes,1,opt,name=interval,casttype=DurationString"`
	Query    string         `json:"query" protobuf:"bytes,2,opt,name=query"`
}
