Node Metrics Collection
Performance Management

Purpose
Collection, summarization and persistence of node metrics for inclusion in performance analysis. Metrics are persisted for anonymous analysis and sharing with the Cloud Native community at-large. Publicly-shared and project-agnostic analysis serves to increase confidence in service mesh adoption.
Guiding Principles
Meshery may require local storage for temporary use w/o guarantee of data resiliency. 
Meshery can require and install certain temporary infrastructure on the cluster (e.g. InfluxDB).
Goals
Users will want to compare service mesh overhead against application resource consumption.
Users 
Design
Node metrics should not overlap, but complement request metrics gleaned from Meshery’s load generators, if possible. Tail latencies are readily available from the load-generator Fortio. 
Trends, not static interval sampling may prove both easier to store (less data) and potentially more insightful benchmarks.
Reports / Node Metrics
While many metrics may be viewed in Meshery in real-time, only a certain set will be considered for long-term persistence. The following list of per node metrics to collect and store by both control plane and by namespace (application).

There are a few ways we can go about this:
Cluster wide metrics
Node level metrics

Cluster-wide Metrics:
Here is a list of cluster wide metrics we can collect from Prometheus node exporter:

CPU utilization - total
CPU requests commitment - total
CPU limits commitment - total
Memory utilization - total
Memory requests commitment - total
Memory limits commitment - total
CPU usage by namespace - time series per namespace
Memory usage without cache by namespace - time series per namespace

Consumption breakdown by control plane vs data plane vs application would be interesting.

The advantage with cluster wide metrics is that Meshery doesn’t have to include logic to account for individual nodes.
Per Node Metrics:
Here is a list of node level metrics we can collect from Prometheus node exporter:

CPU utilization - total
CPU usage
Memory utilization - total
Memory usage
Used 
Buffers
Cached
Free
System load
Load 1m
Load 5m
Load 15m
Logical cores
CPU Usage by core (time series for each of the cores)
Network (gather latency and throughput through request load generator)
Storage (not gathered for now)

The challenge with collecting these metrics is that they will have to be collected for each of the nodes. Hence, we will first have to query Prometheus to get a list of nodes, and then query again for each of the nodes collect all the above list of metrics.
Mechanics
Glean from Grafana dashboards and service mesh-specific metrics tracked through Prometheus, but gather directly from Prometheus.
Requirements
Prometheus node exporter deployed (i.e. as a daemonset),
Reachable URL to Prometheus endpoint for querying.
Metrics Retrieval
Use https://github.com/prometheus/client_golang and one for creating clients that talk to the Prometheus HTTP API.
 Execution plan for persisting server-side metrics
In the UI, performance page, a UUID is generated when the page loads.
The UUID will be submitted to the server when the load test is initiated.
when the static charts are created, the UUID will be associated with them. When the static charts make calls, they will include the UUID to indicate the calls are indeed from a static chart.
Whenever the prometheus query range handler sees UUID in query params, we will persist the queries (start, end and step can be computed based on the duration of the test) in a global map for easier retrieval.
To support parallel test executions we should probably store the queries in a map[string]map[string]struct{}{} - to store UUID vs all queries that UUID
For horizontal scaling, we can eventually store the map in a dedicated cache layer, like redis, later.
When the load test completes, all the queries in the map will be used to query prometheus with the start and end time of the load test (which can be obtained from fortio response) with a computed step values and the results will now be associated with the query and persisted in SaaS.
Also, do we want to call prometheus for percentile queries? Probably NOT, bcoz we are going to be persisting the time series… may be we should
What if there is a network error or any error during this background process?
Meshery will retry like 10 times with exponential backoff
If failure is permanent, we just persist the board json and let Meshery talk to Prometheus directly?
We also need to send the static board json for persistence
Right after the data is persisted, we should clear the map entry for the UUID.
Right after the completion, a new UUID should be generated in the UI.
What-ifs:
Q: What if the test is run for too short a duration? like 5s. . . in which case the static charts might not all have made their calls?
A: Show the user a message and don’t persist results.
Q: What if the payload generated with all the needed metrics is too large to ship for persistence?
A: Either restrict the maximum time length that a test can be run or show the user a message stating that any results beyond the hour limit will not be persisted.
We might have to send it in chunks of 1 or 2MB to play it safe. Of course, the server should have the capability to put them together as well.


Metrics Summarization (for later)
What level of granularity should be supported? 

A factor of duration (time) and number of samples per minute (count)?

(duration  / # of samples per min) = granularity
	(3600 secs / 20 ) = 3 samples per min

A factor of duration (time) and sample interval (time)?

(duration / interval) = consistent granularity
(3600 secs / 10 sec) = 6 per min

A factor of duration (time) and number of samples per test (count)?

(duration / interval) = granularity
Metrics Storage
??GB in Postgres free tier (output is likely too large for pg)
??GB in S3 free tier (most likely storage)
20GB limit in Dynamodb free tier (too slow and costly)
Implementation
Per node metrics:

The node_boot_time_seconds is invoked at the beginning of the page load to fetch the list of nodes.

http://prometheus-node:32758/api/datasources/proxy/1/api/v1/series?match[]=node_boot_time_seconds%7Bcluster%3D%22%22%2C%20job%3D%22node-exporter%22%7D&start=1568392571&end=1568396171

http://prometheus-node:30234/api/datasources/proxy/1/api/v1/series?match[]=node_boot_time_seconds%7Bcluster%3D%22%22%2C%20job%3D%22node-exporter%22%7D&start=1568392571&end=1568396171

Response:
{"status":"success","data":[{"__name__":"node_boot_time_seconds","endpoint":"https","instance":"10.199.75.57:9100","job":"node-exporter","namespace":"monitoring","pod":"node-exporter-8b6kq","service":"node-exporter"},{"__name__":"node_boot_time_seconds","endpoint":"https","instance":"10.199.75.64:9100","job":"node-exporter","namespace":"monitoring","pod":"node-exporter-9mnpt","service":"node-exporter"}]}

We can make this call, fetch the list of instances and inject it into a templated version of the Grafana board. This way there will be one board for every instance.
Cluster metrics:
http://10.199.75.64:30234/d/efa86fd1d0c121a26444b636a3f509a8/kubernetes-compute-resources-cluster?orgId=1

The advantage with this board is that we are NOT tied to the node but rather it gets us information on the cluster wide cpu and memory usage based on the namespace.

One thing to note: if we want to go this route, support for panels of type “singlestat” will have to be added into Meshery.

Tasks
Explore the value of reporting trends vs. sampling.
Research generally accepted promql for summarizations.
Identify specific cpu and mem metrics to use.


