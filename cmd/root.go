/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
)

var (
	gracefullWaitTime time.Duration
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "portk [flags] <port>",
	Short: "Kill process with port",
	Long: `
The "portk" command is used to kill a process that is using a specific port on your machine. It works by finding the process ID (PID) of the process associated with the specified port and then terminating it.

To use this command, simply run:

	portk <port>

Replace "<port>" with the port number you want to kill. For example, if you want to kill a process using port 8080, you would run:

	portk 8080

The "portk" command also has a flag option that allows you to specify a gracefull wait time before forcefully terminating the process. To use this flag, run:

	portk --wait <time> <port>

Replace "<time>" with the number of seconds you want to wait before forcefully terminating the process, and replace "<port>" with the port number you want to kill. For example, if you want to wait 10 seconds before forcefully terminating a process using port 8080, you would run:

	portk --wait 10s 8080

The "portk" command is useful for quickly killing processes that are consuming ports you no longer need or that may be causing issues on your machine.
`,
	SilenceUsage: true,
	RunE:         run,
}

func init() {
	rootCmd.Flags().DurationVarP(&gracefullWaitTime, "wait", "w", 3*time.Second, "Wait time in seconds before killing the process")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	log.SetFlags(0)

	if len(args) == 0 {
		return fmt.Errorf("Please provide a port number")
	}
	portQuery, err := strconv.Atoi(args[0])
	if err != nil || portQuery == 0 || portQuery > 65535 {
		return fmt.Errorf("Port number should be between 0 and 65535")
	}

	pid, err := GetProcessIDFromPort(int32(portQuery))
	if err != nil {
		return fmt.Errorf("Could not find process for port %d", portQuery)
	}
	err = KillProcess(pid)
	if err != nil {
		return fmt.Errorf("Could not kill process %d", pid)
	}
	log.Println("Done")
	return nil
}

func GetProcessIDFromPort(port int32) (int32, error) {
	var cmd *exec.Cmd
	var output []byte
	var err error

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	} else {
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err = cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	line := lines[1]
	fields := strings.Fields(line)
	if len(fields) > 4 {
		processID, _ := strconv.Atoi(fields[1])
		return int32(processID), nil
	}

	return 0, fmt.Errorf("no process found using port %d", port)
}

func KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}

	if gracefullWaitTime <= 0 {
		log.Println("Killing Process...")
		return p.Kill()
	}

	terminateChan := make(chan string, 2)
	err = p.SendSignal(syscall.SIGINT)
	log.Println("Gracefully terminating...")
	if err != nil {
		log.Println("Killing Process...")
		return p.Kill()
	}
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			exists, err := process.PidExists(pid)
			if !exists || err != nil {
				terminateChan <- "killed"
			}
		}
	}()
	select {
	// Case statement
	case _ = <-terminateChan:
		return nil

	// Still alive
	case <-time.After(gracefullWaitTime):
		log.Println("Killing Process...")
		return p.Kill()
	}
}
