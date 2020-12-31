package impls

import (
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

var defaultCiphers = []string{
	"chacha20-poly1305@openssh.com",
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	"aes128-gcm@openssh.com",
	"arcfour256", "arcfour128", "arcfour",
	"aes128-cbc",
	"3des-cbc",
}

var defaultHostKeyAlgos = []string{
	"ssh-rsa-cert-v01@openssh.com", "ssh-dss-cert-v01@openssh.com",
	"ecdsa-sha2-nistp256-cert-v01@openssh.com",
	"ecdsa-sha2-nistp384-cert-v01@openssh.com",
	"ecdsa-sha2-nistp521-cert-v01@openssh.com",
	"ssh-ed25519-cert-v01@openssh.com",
	"ssh-rsa",
	"ssh-dss",
	"ecdsa-sha2-nistp256",
	"ecdsa-sha2-nistp384",
	"ecdsa-sha2-nistp521",
	"ssh-ed25519",
}

func DialSSH(server, username string, timeout int, auth ...ssh.AuthMethod) (ssh.Conn, <-chan ssh.NewChannel, <-chan *ssh.Request, error) {
	sshConfig := ssh.Config{}
	sshConfig.SetDefaults()
	sshConfig.Ciphers = defaultCiphers

	config := &ssh.ClientConfig{
		Config: sshConfig,
		User:   username,
		Auth:   auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if strings.Index(server, ":") < 0 {
		server += ":22"
	}

	conn, err := net.DialTimeout("tcp", server, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, nil, nil, err
	}

	return ssh.NewClientConn(conn, server, config)

}
