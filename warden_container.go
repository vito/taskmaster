package taskmaster

import (
	"github.com/vito/gordon"
)

type WardenContainer struct {
	Handle string
	client *warden.Client
}

func NewWardenContainer(client *warden.Client) (*WardenContainer, error) {
	response, err := client.Create()
	if err != nil {
		return nil, err
	}

	return &WardenContainer{
		Handle: *response.Handle,
		client: client,
	}, nil
}

func (c *WardenContainer) Destroy() error {
	_, err := c.client.Destroy(c.Handle)
	return err
}

func (c *WardenContainer) Spawn(script string) (JobId, error) {
	res, err := c.client.Spawn(c.Handle, script)
	if err != nil {
		return 0, err
	}

	return JobId(res.GetJobId()), nil
}

func (c *WardenContainer) Stream(jobId JobId) (chan *StreamOutput, error) {
	responses, err := c.client.Stream(c.Handle, uint32(jobId))
	if err != nil {
		return nil, err
	}

	outputs := make(chan *StreamOutput)

	go func() {
		for {
			response, ok := <-responses
			if !ok {
				close(outputs)
				break
			}

			outputs <- &StreamOutput{
				Name: response.GetName(),
				Data: response.GetData(),

				Finished:   response.ExitStatus != nil,
				ExitStatus: response.GetExitStatus(),
			}
		}
	}()

	return outputs, nil
}

func (c *WardenContainer) Run(script string) (*JobInfo, error) {
	res, err := c.client.Run(c.Handle, script)
	if err != nil {
		return nil, err
	}

	return &JobInfo{
		ExitStatus: res.GetExitStatus(),
	}, nil
}

func (c *WardenContainer) NetIn() (MappedPort, error) {
	res, err := c.client.NetIn(c.Handle)
	if err != nil {
		return 0, err
	}

	return MappedPort(res.GetHostPort()), nil
}

func (c *WardenContainer) CopyIn(src, dst string) error {
	_, err := c.client.CopyIn(c.Handle, src, dst)
	return err
}
