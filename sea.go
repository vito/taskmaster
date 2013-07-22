package taskmaster

import (
	"github.com/vito/gordon"
	"errors"
)

type SEA struct {
	WardenSocketPath string
	Registry         *Registry
}

func NewSEA(wardenSocketPath string) *SEA {
	return &SEA{
    WardenSocketPath: wardenSocketPath,
		Registry: NewRegistry(),
	}
}

func (s *SEA) Start(id, publicKey string) error {
	client := warden.NewClient(
		&warden.ConnectionInfo{
			SocketPath: s.WardenSocketPath,
		},
	)

	err := client.Connect()
	if err != nil {
		return err
	}

	container, err := NewWardenContainer(client)
	if err != nil {
		return err
	}

	mappedPort, err := container.NetIn()
	if err != nil {
		return err
	}

	session := &Session{
		Container: container,
		Port:      mappedPort,
	}

	err = session.InsertPubkey(publicKey)
	if err != nil {
		return err
	}

	err = session.StartSSHServer()
	if err != nil {
		return err
	}

	s.Registry.Register(id, session)

	return nil
}

func (s *SEA) Stop(id string) error {
  session := s.Registry.Lookup(id)
  if session == nil {
    return errors.New("unknown session")
  }

  s.Registry.Unregister(id)

  return session.Terminate()
}
