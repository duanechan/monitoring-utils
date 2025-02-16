package main

import (
	"bufio"
	credentials "credentials/internal"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func main() {
	path := flag.String("path", "", "the file path to the list of recipients")
	flag.Parse()

	for {
		clearTerminal()
		header()

		if *path == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Filepath: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("error: failed to read CSV file: %s", err)
			}

			*path = strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(input, "\n"), "\r", ""))
		}

		recipients, err := getRecipients(*path)
		if err != nil {
			log.Fatalf("error: failed to read recipient data: %s", err)
		}

		sent := 0
		var wg sync.WaitGroup

		for _, r := range recipients {
			wg.Add(1)
			go func(r credentials.User) {
				done := make(chan bool)

				defer func() {
					close(done)
					wg.Done()
				}()

				go showLoadingBar(done)

				if err := credentials.SendEmail(r); err != nil {
					color.Red("\r✖ Failed to send credentials email: %s\n", err)
					done <- false
					return
				}

				done <- true
				sent++

				color.Green("\r✔ Credentials email successfully sent to %s\n", r)

				time.Sleep(100 * time.Millisecond)

			}(r)
		}

		wg.Wait()

		if choice := generateReport(sent, len(recipients)); choice == "y" || choice == "Y" {
			*path = ""
			continue
		} else if choice == "n" || choice == "N" {
			os.Exit(0)
		}
	}
}

func generateReport(sentEmails, recipientsLen int) string {
	reader := bufio.NewReader(os.Stdin)

	color.RGB(103, 150, 191).Println(strings.Repeat("_", 75))
	fmt.Println()
	color.New(color.FgGreen).Printf("SUCCESS: %d     ", sentEmails)
	color.New(color.FgRed).Println("FAILED: ", recipientsLen-sentEmails)
	color.RGB(103, 150, 191).Println(strings.Repeat("_", 75))
	fmt.Println()

	fmt.Printf("Would you like to send another batch? [Y/n]: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error: unable to process input: %s", err)
	}

	return strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(input, "\n"), "\r", ""))
}

func header() {
	color.RGB(103, 150, 191).Println(` _____              _            _   _       _       _   _      _                 
/  __ \            | |          | | (_)     | |     | | | |    | |                
| /  \/_ __ ___  __| | ___ _ __ | |_ _  __ _| |___  | |_| | ___| |_ __   ___ _ __ 
| |   | '__/ _ \/ _` + "`" + ` |/ _ | '_ \| __| |/ _` + "`" + ` | / __| |  _  |/ _ | | '_ \ / _ | '__|
| \__/| | |  __| (_| |  __| | | | |_| | (_| | \__ \ | | | |  __| | |_) |  __| |   
 \____|_|  \___|\__,_|\___|_| |_|\__|_|\__,_|_|___/ \_| |_/\___|_| .__/ \___|_|   
                                                                 | |              
                                                                 |_|     `)
	fmt.Println()
	fmt.Println()
}

func clearTerminal() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func getRecipients(filepath string) ([]credentials.User, error) {
	if !strings.HasSuffix(filepath, ".csv") {
		return []credentials.User{}, fmt.Errorf("file is not CSV")
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return []credentials.User{}, err
	}

	recipients := []credentials.User{}

	recipientsData := strings.Split(string(bytes), "\n")
	for _, d := range recipientsData {
		fields := strings.Split(d, ",")
		recipients = append(recipients, credentials.User{
			Name:  strings.TrimSpace(strings.ReplaceAll(fields[0], "\r", "")),
			Email: strings.TrimSpace(strings.ReplaceAll(fields[1], "\r", "")),
		})
	}

	return recipients, nil
}

func showLoadingBar(done chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r") // Clear the line
			return
		case <-ticker.C:
			color.New(color.FgYellow).Printf("\rSending email... %s", frames[i])
			// fmt.Printf("\rSending email... %s", frames[i])
			i = (i + 1) % len(frames)
		}
	}
}
