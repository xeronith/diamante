package analytics

type IMeasurementsProvider interface {
	SubmitMeasurement(string, Tags, Fields)
	SubmitMeasurementAsync(string, Tags, Fields)
}
