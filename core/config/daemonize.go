package config

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/takama/daemon"
)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// DaemonRoutine ...
type DaemonRoutine func()

var stdlog, errlog *log.Logger

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage(callback DaemonRoutine) (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	go callback()

	for {
		select {
		case killSignal := <-interrupt:
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

// BackGroundService run in background 
func BackGroundService(name, description string, dependencies []string, callback DaemonRoutine) string {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}

	service := &Service{srv}
	status, err := service.Manage(callback)
	if err != nil {

		if len(os.Args) > 1 {
			command := os.Args[1]
			switch command {

			case "remove":
				if err == daemon.ErrNotInstalled {
					errlog.Println(err)
					os.Exit(0)
				}
			case "stop":
				if err == daemon.ErrAlreadyStopped {
					errlog.Println(err)
					os.Exit(0)
				}
			}
		}

		errlog.Println("Error: ", err)
		os.Exit(1)
	}

	return status
}
