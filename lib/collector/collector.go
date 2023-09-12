package collector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/sonarcloud-exporter/internal"
	"github.com/whyeasy/sonarcloud-exporter/lib/client"
)

const (
	METRIC_PREFIX = "sonarcloud_"
)

// Collector struct for holding Prometheus Desc and Exporter Client
type Collector struct {
	up     *prometheus.Desc
	client *client.ExporterClient

	projectInfo *prometheus.Desc

	metrics map[string]*prometheus.Desc

	linesOfCode     *prometheus.Desc
	codeCoverage    *prometheus.Desc
	vulnerabilities *prometheus.Desc
	bugs            *prometheus.Desc
	codeSmells      *prometheus.Desc
}

// New creates a new Collecotor with Prometheus descriptors
func New(c *client.ExporterClient, cfg internal.Config) *Collector {
	log.Info("Creating collector")

	metrics := make(map[string]*prometheus.Desc)

	for _, metric := range strings.Split(cfg.Metrics, ",") {
		metrics[metric] = prometheus.NewDesc(METRIC_PREFIX+metric, fmt.Sprintf("%s within a project in SonarCloud", metric), []string{"project_key"}, nil)
	}

	return &Collector{
		up:          prometheus.NewDesc("sonarcloud_up", "Whether Sonarcloud scrape was successfull", nil, nil),
		client:      c,
		projectInfo: prometheus.NewDesc(METRIC_PREFIX+"project_info", "General information about projects", []string{"project_name", "project_qualifier", "project_key", "project_organization"}, nil),
		metrics:     metrics,
	}
}

// Describe the metrics that are collected
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	ch <- c.projectInfo

	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	log.Info("Running scrape")

	if stats, err := c.client.GetStats(); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectProjectInfo(c, ch, stats)

		collectMeasurements(c, ch, stats)

		log.Info("Scrape Complete")
	}
}

func collectProjectInfo(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, project := range *stats.Projects {
		value := 0.0
		if project.LastAnalysis != nil {
			value = float64(project.LastAnalysis.Unix())
		}

		ch <- prometheus.MustNewConstMetric(c.projectInfo, prometheus.GaugeValue, value, project.Name, project.Qualifier, project.Key, project.Organization)
	}
}

func collectMeasurements(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, measurement := range *stats.Measurements {
		value, err := strconv.ParseFloat(measurement.Value, 64)
		if err != nil {
			log.Error(err)
		}
		metricDesc := c.metrics[measurement.Metric]

		ch <- prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value, measurement.Key)
	}
}
