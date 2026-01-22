package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomsteele/wholepunch/pkg/wp"
)

var (
	flHTTPSGetSpoofSNIServerPort string
	flHTTPSGetSpoofSNIServerName string
	flHTTPSGetSpoofSNIUserAgent  string
	flHTTPSGetSpoofSNIURLPath    string
)

func init() {
	beaconCmd.AddCommand(httpsGetSpoofSNICmd)
	httpsGetSpoofSNICmd.PersistentFlags().StringVar(&flHTTPSGetSpoofSNIServerPort, "server-port", "443", "HTTPS port to connect to.")
	httpsGetSpoofSNICmd.PersistentFlags().StringVar(&flHTTPSGetSpoofSNIServerName, "server-name", "www.microsoft.com", "Server Name to use in the SNI client hello.")
	httpsGetSpoofSNICmd.PersistentFlags().StringVar(&flHTTPSGetSpoofSNIUserAgent, "user-agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko", "User-Agent to use during HTTPS request.")
	httpsGetSpoofSNICmd.PersistentFlags().StringVar(&flHTTPSGetSpoofSNIURLPath, "path", "/", "URL path to use during HTTPS request.")
}

func httpsGetSpoofSNI(cmd *cobra.Command, args []string) {
	opts := wp.BeaconOptions{
		DestinationServerAddress: fmt.Sprintf("%s:%s", flBeaconServerAddr, flHTTPSGetSpoofSNIServerPort),
	}
	b := wp.BeaconHTTPSGetSpoofSNI{
		ServerName: flHTTPSGetSpoofSNIServerName,
		UserAgent:  flHTTPSGetSpoofSNIUserAgent,
		Path:       flHTTPSGetSpoofSNIURLPath,
	}
	ok, err := wp.RunBeacon(&b, &opts)
	result := wp.MakeBeaconResult(ok, err, &b)
	wp.WriteTableBeaconResults(os.Stdout, []wp.BeaconResult{result}, flBeaconFilterFalse)
}

var httpsGetSpoofSNICmd = &cobra.Command{
	Use:   "https-get-spoof-sni",
	Short: "Send an HTTPS GET request using a spoofed SNI in the TLS ClientHello.",
	Run:   httpsGetSpoofSNI,
}
