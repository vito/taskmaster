package taskmaster

import (
	"fmt"
)

type Session struct {
	Container Container
	Port      MappedPort
}

func (s *Session) InsertPubkey(publicKey string) error {
	_, err := s.Container.Run(
		fmt.Sprintf(`
      mkdir -p .ssh
      echo '%s' > .ssh/authorized_keys
    `, publicKey),
	)

	return err
}

func (s *Session) StartSSHServer() error {
	_, err := s.Container.Spawn(
		fmt.Sprintf(`
      set -e
      export PATH=/usr/local/bin:$PATH
      dropbearkey -t rsa -f dropbear_rsa_host_key
      dropbear -r dropbear_rsa_host_key -F -E -p 0.0.0.0:%d
    `, s.Port),
	)

	return err
}

func (s *Session) Terminate() error {
	return s.Container.Destroy()
}
