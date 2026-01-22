package wp

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// BeaconSSH connects to a server using SSH.
type BeaconSSH struct {
	ServerAddr string
	Timeout    int
}

// Name returns the name of the module.
func (b *BeaconSSH) Name() string {
	return "ssh"
}

// Destination returns the server that was connected to.
func (b *BeaconSSH) Destination() string {
	return fmt.Sprintf("ssh://%s", b.ServerAddr)
}

// Success returns a formatted string indicating a successfull connection.
func (b *BeaconSSH) Success() string {
	return fmt.Sprintf("The agent was allowed to reach %s using SSH.", b.ServerAddr)
}

// Setup is used to initilize instance variables from BeaconOptions.
func (b *BeaconSSH) Setup(o *BeaconOptions) error {
	b.ServerAddr = o.DestinationServerAddress
	return nil
}

// Send initiates the SSH connection.
func (b *BeaconSSH) Send() (bool, error) {
	conn, err := net.DialTimeout("tcp", b.ServerAddr, time.Duration(b.Timeout)*time.Millisecond)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Duration(b.Timeout) * time.Millisecond))

	// Dummy config for egress testing. Auth will fail but that's fine.
	config := &ssh.ClientConfig{
		User: "egress-test",
		Auth: []ssh.AuthMethod{
			ssh.Password("egress-test"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(b.Timeout) * time.Millisecond,
	}

	sshConn, _, _, err := ssh.NewClientConn(conn, b.ServerAddr, config)
	if sshConn != nil {
		sshConn.Close()
	}

	// Auth errors mean SSH was reachable.
	if err != nil {
		if _, ok := err.(*ssh.ServerAuthError); ok {
			return true, nil
		}
		if err.Error() == "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain" {
			return true, nil
		}
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false, err
		}
		// Got a response from an SSH server.
		return true, nil
	}

	return true, nil
}
