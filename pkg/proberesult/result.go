package proberesult

type ProbeResult struct {
	Duration       float64
	Success        bool
	Message        string
	StatusCode     int
	ContentLength  int64
	TLSVersion     string
	CertExpiryDays int
}

func New(duration float64) *ProbeResult {
	return &ProbeResult{
		Duration: duration,
		Success:  false,
	}
}

func (r *ProbeResult) SetSuccess(success bool) {
	r.Success = success
}

func (r *ProbeResult) SetMessage(message string) {
	r.Message = message
}

func (r *ProbeResult) SetStatusCode(statusCode int) {
	r.StatusCode = statusCode
}

func (r *ProbeResult) SetContentLength(contentLength int64) {
	r.ContentLength = contentLength
}

func (r *ProbeResult) SetTLSVersion(tlsVersion string) {
	r.TLSVersion = tlsVersion
}

func (r *ProbeResult) SetCertExpiryDays(days int) {
	r.CertExpiryDays = days
}
