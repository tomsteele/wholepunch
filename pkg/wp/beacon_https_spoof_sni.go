package wp

import (
	"fmt"

	"github.com/tomsteele/wholepunch/pkg/http"
)

// BeaconHTTPSGetSpoofSNI sends an HTTPS GET request with a spoofed SNI in the TLS ClientHello.
type BeaconHTTPSGetSpoofSNI struct {
	ServerName string
	UserAgent  string
	ServerAddr string
	Path       string
	ServerURL  string
}

// Name returns the name of the module.
func (b *BeaconHTTPSGetSpoofSNI) Name() string {
	return "https-get-spoof-sni"
}

// Destination returns the server that was connected to
func (b *BeaconHTTPSGetSpoofSNI) Destination() string {
	return b.ServerURL
}

// Success returns a formatted string indicating a successfull connection.
func (b *BeaconHTTPSGetSpoofSNI) Success() string {
	return fmt.Sprintf("The agent was allowed to communicate with %s over HTTPS using the SNI name %s to bypass egress controls.", b.ServerAddr, b.ServerName)
}

// Setup is used to initilize instance variables from BeaconOptions.
func (b *BeaconHTTPSGetSpoofSNI) Setup(o *BeaconOptions) error {
	b.ServerAddr = o.DestinationServerAddress
	b.ServerURL = fmt.Sprintf("https://%s%s", b.ServerAddr, b.Path)
	return nil
}

// Send initiates the HTTPS request with spoofed SNI.
func (b *BeaconHTTPSGetSpoofSNI) Send() (bool, error) {
	return http.TLSGetSpoofSNI(b.ServerURL, b.ServerName, b.UserAgent)
}
