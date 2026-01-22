package cmd

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/cobra"
	"github.com/tomsteele/wholepunch/pkg/wp"
	"github.com/tomsteele/xplode"
)

var (
	flSSHServerPort string
	flSSHTimeout    int
	flSSHWorkers    int
)

func init() {
	beaconCmd.AddCommand(sshCmd)
	sshCmd.PersistentFlags().StringVar(&flSSHServerPort, "p", "22", "NMap style port string.")
	sshCmd.PersistentFlags().IntVar(&flSSHTimeout, "timeout", 5000, "Timeout in milliseconds.")
	sshCmd.PersistentFlags().IntVar(&flSSHWorkers, "c", 50, "Max number of concurrent requests.")
}

func sshBeacon(cmd *cobra.Command, args []string) {
	ports, err := xplode.Parse(flSSHServerPort)
	if err != nil {
		fmt.Println("There was an error parsing the port string.")
		fmt.Println(err)
		os.Exit(1)
	}
	results := []wp.BeaconResult{}

	mutex := sync.Mutex{}
	portChan := make(chan int)
	doneChan := make(chan bool)

	w := wow.New(os.Stdout, spin.Get(spin.Pipe), "Working")
	w.Start()
	for i := 0; i < flSSHWorkers; i++ {
		go func() {
			for p := range portChan {
				b := wp.BeaconSSH{
					Timeout: flSSHTimeout,
				}
				opts := wp.BeaconOptions{
					DestinationServerAddress: fmt.Sprintf("%s:%d", flBeaconServerAddr, p),
				}
				ok, err := wp.RunBeacon(&b, &opts)
				result := wp.MakeBeaconResult(ok, err, &b)
				mutex.Lock()
				results = append(results, result)
				mutex.Unlock()
				doneChan <- true
			}
		}()
	}
	go func() {
		for _, p := range ports {
			portChan <- p
		}
	}()
	for i := 0; i < len(ports); i++ {
		<-doneChan
	}
	close(portChan)
	close(doneChan)
	w.Stop()
	fmt.Println()

	sort.Slice(results, func(i, j int) bool {
		iurl, _ := url.Parse(results[i].Destination)
		jurl, _ := url.Parse(results[j].Destination)
		i, _ = strconv.Atoi(iurl.Port())
		j, _ = strconv.Atoi(jurl.Port())
		return i < j
	})
	wp.WriteTableBeaconResults(os.Stdout, results, flBeaconFilterFalse)
}

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Initiates SSH connections to test if the SSH protocol is allowed outbound.",
	Run:   sshBeacon,
}
