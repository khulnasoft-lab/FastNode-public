package clienttelemetry

var (
	// InstallIDs sends to Mixpanel & Customer.io under the install ID
	InstallIDs = Options{mp: InstallID, cio: InstallID}
	// FastnodeOnly sends to Fastnode under the metrics ID
	FastnodeOnly = Options{fastnode: MetricsID}
)

// EventWithFastnodeTelemetry sends events to Mixpanel, CIO, and t.khulnasoft.com
func EventWithFastnodeTelemetry(name string, props map[string]interface{}) {
	Default.Fastnode(MetricsID).Event(name, props)
}

// FastnodeTelemetry aliases FastnodeOnly.Event
func FastnodeTelemetry(name string, props map[string]interface{}) {
	FastnodeOnly.Event(name, props)
}

// Event aliases Default.Event
func Event(name string, props map[string]interface{}) {
	Default.Event(name, props)
}

// Update aliases Default.Update
func Update(props map[string]interface{}) {
	Default.Update(props)
}
