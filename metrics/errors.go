package metrics

import "errors"

var (
	ErrRegistryClosed    = errors.New("metrics: registry is closed")
	ErrMetricNotFound    = errors.New("metrics: metric not found")
	ErrMetricExists      = errors.New("metrics: metric already exists")
	ErrInvalidMetricName = errors.New("metrics: invalid metric name")
	ErrInvalidLabelName  = errors.New("metrics: invalid label name")
	ErrLabelMismatch     = errors.New("metrics: label names do not match")
	ErrExporterFailed    = errors.New("metrics: exporter failed")
)
