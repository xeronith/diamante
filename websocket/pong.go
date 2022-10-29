package websocket

import "time"

func (c *Conn) WritePong() error {
	var timeout time.Duration = 4000000000
	deadline := time.Now().Add(timeout)

	buf := make([]byte, 0, 139)
	buf = append(buf, 138, 0)

	timer := time.NewTimer(timeout)
	select {
	case <-c.mu:
		timer.Stop()
	case <-timer.C:
		return errWriteTimeout
	}

	defer func() { c.mu <- true }()

	c.writeErrMu.Lock()
	err := c.writeErr
	c.writeErrMu.Unlock()
	if err != nil {
		return err
	}

	c.conn.SetWriteDeadline(deadline)
	_, err = c.conn.Write(buf)
	if err != nil {
		return c.writeFatal(err)
	}

	return err
}
