package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Registrar json structure
type Registrar struct {
	Writes struct {
		Fail    float64 `json:"fail"`
		Success float64 `json:"success"`
		Total   float64 `json:"total"`
	} `json:"writes"`
	States struct {
		Cleanup float64 `json:"cleanup"`
		Current float64 `json:"current"`
		Update  float64 `json:"update"`
	} `json:"states"`
	Progress []struct {
		Source string  `json:"source"`
		Size   float64 `json:"size"`
		Offset float64 `json:"offset"`
	} `json:"progress"`
}

type registrarCollector struct {
	beatInfo *BeatInfo
	stats    *Stats
	metrics  exportedMetrics
}

// NewRegistrarCollector constructor
func NewRegistrarCollector(beatInfo *BeatInfo, stats *Stats) prometheus.Collector {
	return &registrarCollector{
		beatInfo: beatInfo,
		stats:    stats,
		metrics: exportedMetrics{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "registrar", "writes"),
					"registrar.writes",
					nil, prometheus.Labels{"writes": "fail"},
				),
				eval:    func(stats *Stats) float64 { return stats.Registrar.Writes.Fail },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "registrar", "writes"),
					"registrar.writes",
					nil, prometheus.Labels{"writes": "success"},
				),
				eval:    func(stats *Stats) float64 { return stats.Registrar.Writes.Success },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "registrar", "writes"),
					"registrar.writes",
					nil, prometheus.Labels{"writes": "total"},
				),
				eval:    func(stats *Stats) float64 { return stats.Registrar.Writes.Total },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "registrar", "states"),
					"registrar.states",
					nil, prometheus.Labels{"state": "cleanup"},
				),
				eval:    func(stats *Stats) float64 { return stats.Registrar.States.Cleanup },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "registrar", "states"),
					"registrar.states",
					nil, prometheus.Labels{"state": "current"},
				),
				eval:    func(stats *Stats) float64 { return stats.Registrar.States.Current },
				valType: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "registrar", "states"),
					"registrar.states",
					nil, prometheus.Labels{"state": "update"},
				),
				eval:    func(stats *Stats) float64 { return stats.Registrar.States.Update },
				valType: prometheus.GaugeValue,
			},
		},
	}
}

// Describe returns all descriptions of the collector.
func (c *registrarCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, metric := range c.metrics {
		ch <- metric.desc
	}

}

// Collect returns the current state of all metrics of the collector.
func (c *registrarCollector) Collect(ch chan<- prometheus.Metric) {

	for _, i := range c.metrics {
		ch <- prometheus.MustNewConstMetric(i.desc, i.valType, i.eval(c.stats))
	}

	for _, progress := range c.stats.Registrar.Progress {
		progressSizeDesc := prometheus.NewDesc(
			prometheus.BuildFQName(c.beatInfo.Beat, "registrar", "progress_size"),
			"registrar.progress",
			nil, prometheus.Labels{"progress_size": progress.Source},
		)
		progressOffsetDesc := prometheus.NewDesc(
			prometheus.BuildFQName(c.beatInfo.Beat, "registrar", "progress_offset"),
			"registrar.progress",
			nil, prometheus.Labels{"progress_offset": progress.Source},
		)
		ch <- prometheus.MustNewConstMetric(progressSizeDesc, prometheus.GaugeValue, progress.Size)
		ch <- prometheus.MustNewConstMetric(progressOffsetDesc, prometheus.GaugeValue, progress.Offset)
	}

}
