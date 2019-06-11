//
// graylog-exporter
//
// Prometheus Exporter for Graylog API
//
// Author: Martin Weber <martin.weber@de.clara.net>
// Company: Claranet GmbH
// Version: 0.2.0
//

package main

import (
	"./graylog"

	"flag"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

var defaultMetrics = "jvm.memory.total.max,jvm.memory.total.used,jvm.memory.total.init,org.graylog2.throughput.output,org.graylog2.throughput.input,org.graylog2.throughput.input.1-sec-rate,org.graylog2.throughput.output.1-sec-rate,org.graylog2.journal.entries-uncommitted"

var (
	namespace       = "graylog"
	graylogUrl      = flag.String("graylog.url", "", "Graylog URL to connect to API https://graylog.local.host:9000")
	graylogUser     = flag.String("graylog.username", "", "Graylog API User")
	graylogPassword = flag.String("graylog.password", "", "Graylog API User Password")
	listenAddress   = flag.String("listen-address", ":9404", "The address to lisiten on for HTTP requests.")
	graylogMetrics  = flag.String("graylog.metrics", defaultMetrics, "Graylog metrics to export")
)

var (
	graylogApi *graylog.Graylog
	// Current Session Age
//	graylogSessionAge		int64		= 0
)

type Exporter struct {
	CountNodes                 *prometheus.GaugeVec
	MessageThroughput          *prometheus.GaugeVec
	JournalReadEventsPerSecond *prometheus.GaugeVec
	JournalUncomittedEntries   *prometheus.GaugeVec
	ClusterNodeMetric          *prometheus.GaugeVec
	SystemMetric               *prometheus.GaugeVec
}

func NewExporter() *Exporter {
	defaultLabels := []string{"hostname"}

	return &Exporter{
		CountNodes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "count_nodes",
			Help: "Count Nodes in Graylog Cluster",
		}, []string{}),
		MessageThroughput: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "message_throughput",
			Help: "Message Throughput of current node",
		}, defaultLabels),
		JournalReadEventsPerSecond: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "journal_read_events",
			Help: "Journal Read Events per second",
		}, defaultLabels),
		JournalUncomittedEntries: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "journal_uncommitted_entries",
			Help: "Uncommited entries of Jouranl",
		}, defaultLabels),
		ClusterNodeMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "cluster_node_metric",
			Help: "Cluster Node Metrics",
		}, []string{"hostname", "metric"}),
		SystemMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "system_metric",
			Help: "System Metrics",
		}, []string{"metric"}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.CountNodes.Describe(ch)
	e.MessageThroughput.Describe(ch)
	e.JournalReadEventsPerSecond.Describe(ch)
	e.JournalUncomittedEntries.Describe(ch)
	e.ClusterNodeMetric.Describe(ch)
	e.SystemMetric.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	nodes := graylogApi.GetSystemClusterNodes()
	{
		g := e.CountNodes.WithLabelValues()
		g.Set(float64(nodes.Total))
		g.Collect(ch)
	}

	var nodeMetrics = strings.Split(*graylogMetrics, ",")
	for i := 0; i < len(nodes.Nodes); i++ {
		node := nodes.Nodes[i]

		{
			tput := graylogApi.GetSystemThroughput()
			g := e.MessageThroughput.WithLabelValues(node.Hostname)
			g.Set(float64(tput.Throughput))
			g.Collect(ch)
		}

		{
			journal := graylogApi.GetClusterJournal(node.NodeId)

			g := e.JournalReadEventsPerSecond.WithLabelValues(node.Hostname)
			g.Set(float64(journal.ReadEventsPerSecond))
			g.Collect(ch)

			g = e.JournalUncomittedEntries.WithLabelValues(node.Hostname)
			g.Set(float64(journal.UncommittedJournalEntries))
			g.Collect(ch)

		}

		metricValues := graylogApi.GetClusterNodeMetrics(node.NodeId, nodeMetrics)
		log.Info(metricValues.Total)

		for i := 0; i < metricValues.Total; i++ {
			m := metricValues.Metrics[i]
			g := e.ClusterNodeMetric.WithLabelValues(node.Hostname, m.FullName)

			if m.Type == "counter" {
				g.Set(float64(m.Metric.Count))
			} else {
				g.Set(float64(m.Metric.Value))
			}
			g.Collect(ch)
		}
	}

	systemMetrics := graylogApi.GetSystemMetrics()
	for metric := range systemMetrics.Gauges {
		value := systemMetrics.Gauges[metric].Value
		g := e.SystemMetric.WithLabelValues(metric)
		g.Set(value)
		g.Collect(ch)
	}

}

func main() {
	flag.Parse()

	log.Debug("Create Graylog instance")
	graylogApi = graylog.NewGraylog(*graylogUrl, *graylogUser, *graylogPassword)

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>NameNode Exporter</title></head>
		<body>
		<h1>NameNode Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>`))
	})

	log.Printf("Starting Server: %s", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
