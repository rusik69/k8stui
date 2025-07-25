package app

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ShowContainerLogs displays logs for a container
func (a *App) ShowContainerLogs(containerName string) error {
	if a.KubeClient == nil || a.SelectedPod == "" || containerName == "" {
		return fmt.Errorf("kubernetes client, pod, or container not selected")
	}

	// Clear the logs view and show loading message
	a.LogsView.Clear()
	a.LogsView.SetText("[yellow]Loading logs...")
	
	// Stop any existing log stream
	if a.logStream != nil {
		a.logStream.Close()
	}
	if a.logStopChan != nil {
		close(a.logStopChan)
	}
	a.logStopChan = make(chan struct{})

	// Create log stream request
	podLogOpts := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
		TailLines: func() *int64 { i := int64(100); return &i }(), // Show last 100 lines
	}

	req := a.KubeClient.CoreV1().Pods(a.CurrentNs).GetLogs(a.SelectedPod, podLogOpts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		a.LogsView.SetText(fmt.Sprintf("[red]Error opening log stream: %v", err))
		return fmt.Errorf("error opening log stream: %v", err)
	}

	a.logStream = stream

	// Start a goroutine to read logs
	go func() {
		defer stream.Close()
		buf := make([]byte, 4096)
		firstRead := true
		for {
			select {
			case <-a.logStopChan:
				return
			default:
				n, err := stream.Read(buf)
				if n > 0 {
					a.App.QueueUpdateDraw(func() {
						if firstRead {
							// Clear loading message on first log data
							a.LogsView.Clear()
							firstRead = false
						}
						fmt.Fprintf(a.LogsView, "%s", string(buf[:n]))
						a.LogsView.ScrollToEnd()
					})
				}
				if err != nil {
					if err.Error() != "EOF" {
						a.App.QueueUpdateDraw(func() {
							a.LogsView.SetText(fmt.Sprintf("[red]Error reading logs: %v", err))
						})
					}
					return
				}
			}
		}
	}()

	a.SelectedCont = containerName
	a.LogsView.SetTitle(fmt.Sprintf(" Logs: %s ", containerName))

	return nil
}


