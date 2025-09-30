package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

func parseInterval(s string) (time.Duration, error) {
	re := regexp.MustCompile(`^(\d+)([smh])$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid interval format: %s", s)
	}
	num, _ := strconv.Atoi(matches[1])
	unit := matches[2]

	switch unit {
	case "s":
		return time.Duration(num) * time.Second, nil
	case "m":
		return time.Duration(num) * time.Minute, nil
	case "h":
		return time.Duration(num) * time.Hour, nil
	}
	return 0, fmt.Errorf("unsupported unit: %s", unit)
}

func formatElapsedTime(elapsed time.Duration) string {
    hours := int(elapsed.Hours())
    minutes := int(elapsed.Minutes()) % 60
    seconds := int(elapsed.Seconds()) % 60
    return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func runCommand(cmdStr string) {
    cmd := exec.Command("sh", "-c", cmdStr)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error running command: %v\n", err)
    }
    if len(output) > 0 {
        fmt.Print(string(output))
    }
}


func main() {
	cmdStr := flag.String("c", "", "command to run (required)")
	intervalStr := flag.String("i", "5m", "interval (e.g. 10s, 5m, 1h)")
	flag.Parse()

	if *cmdStr == "" {
		fmt.Println("Usage: rep-cmd -c 'command' -i 5m")
		os.Exit(1)
	}

	interval, err := parseInterval(*intervalStr)
	if err != nil {
		fmt.Println("Error parsing interval:", err)
		os.Exit(1)
	}

	fmt.Printf("Running `%s` every %s. Press Ctrl+C to stop.\n", *cmdStr, interval)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("\nReceived signal: %v, exiting...\n", sig)
		done <- true
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	lastResetTime := time.Now()
	
	timerTicker := time.NewTicker(time.Second)
	defer timerTicker.Stop()

	runCommand(*cmdStr)

	for {
		select {
		case <-timerTicker.C:
			elapsed := time.Since(lastResetTime)
			fmt.Printf("\r\033[K[%s] ", formatElapsedTime(elapsed))
		case <-ticker.C:
			lastResetTime = time.Now()
			fmt.Println()
			runCommand(*cmdStr)
		case <-done:
			fmt.Println()
			return
		}
	}
}
