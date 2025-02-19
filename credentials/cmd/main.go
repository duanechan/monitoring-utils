package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	credentials "github.com/duanechan/monitoring-utils/credentials/internal"
	"github.com/fatih/color"
)

func main() {
	// Read -path flag
	path := flag.String("path", "", "the file path to the list of recipients")
	flag.Parse()

	for {
		ClearTerminal()
		Header()

		// If path flag is empty, specify filepath on runtime
		if *path == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Filepath: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("error: failed to read CSV file: %s", err)
			}

			*path = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(input, "\n"), "\"", ""), "\r", ""))
		}

		// Parse recipient data
		recipients, err := GetRecipients(*path)
		if err != nil {
			log.Fatalf("error: failed to read recipient data: %s", err)
		}

		// Send emails to each recipient
		sent := 0
		var wg sync.WaitGroup

		for i, r := range recipients {
			wg.Add(1)
			go func(r credentials.User) {
				done := make(chan bool)

				defer func() {
					close(done)
					wg.Done()
				}()

				go ShowLoadingBar(done)

				if err := credentials.SendEmail(r); err != nil {
					color.HiRed("\r✖ Sending credentials email to record (row) %d failed: %s\n", i+1, err)
					done <- false
					return
				}

				done <- true
				sent++

				color.HiGreen("\r✔ Credentials email successfully sent to %s\n", r)

				time.Sleep(100 * time.Millisecond)

			}(r)
		}

		wg.Wait()

		// Print report and ask the user to send another batch or not
		if choice := GenerateReport(sent, len(recipients)); choice == "y" || choice == "Y" {
			*path = ""
			continue
		} else if choice == "n" || choice == "N" {
			os.Exit(0)
		}
	}
}

// Prints a short report of the number of successful and failed sent emails.
func GenerateReport(sentEmails, recipientsLen int) string {
	reader := bufio.NewReader(os.Stdin)

	color.RGB(103, 150, 191).Println(strings.Repeat("_", 75))
	fmt.Println()
	color.New(color.FgHiGreen).Printf("SUCCESS: %d     ", sentEmails)
	color.New(color.FgHiRed).Println("FAILED: ", recipientsLen-sentEmails)
	color.RGB(103, 150, 191).Println(strings.Repeat("_", 75))
	fmt.Println()

	fmt.Printf("Would you like to send another batch? [Y/n]: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error: unable to process input: %s", err)
	}

	return strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(input, "\n"), "\r", ""))
}

// Header text of the Credentials Helper program.
func Header() {
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

// Clears the terminal window.
func ClearTerminal() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Parses raw (CSV file) data and returns a slice of recipients.
func GetRecipients(filepath string) ([]credentials.User, error) {
	if !strings.HasSuffix(filepath, ".csv") {
		return []credentials.User{}, fmt.Errorf("file is not CSV")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return []credentials.User{}, err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return []credentials.User{}, err
	}

	recipients := []credentials.User{}
	invalidEmails := map[int]string{}

	for i, r := range records {
		name := strings.TrimSpace(strings.ReplaceAll(r[0], "\r", ""))
		email := strings.TrimSpace(strings.ReplaceAll(r[1], "\r", ""))

		if !credentials.IsValidEmail(email) {
			invalidEmails[i+1] = email
		}

		recipients = append(recipients, credentials.User{
			Name:  name,
			Email: email,
		})
	}

	if len(invalidEmails) > 0 {
		fmt.Println()
		color.HiYellow("There is/are %d invalid email/s in the file:\n", len(invalidEmails))
		for k, v := range invalidEmails {
			fmt.Printf("-> Record (row) %d: %s\n", k, v)
		}

		fmt.Printf(
			"\n%s Press Enter to %s or CTRL+C to %s.",
			color.New(color.FgHiYellow).Sprintf("Are you sure you want to continue?"),
			color.New(color.FgHiYellow).Sprintf("continue"),
			color.New(color.FgHiGreen).Sprintf("cancel"),
		)
		fmt.Scanln()
		fmt.Println()
	}

	return recipients, nil
}

// Displays a loading animation when sending emails.
func ShowLoadingBar(done chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r")
			return
		case <-ticker.C:
			color.New(color.FgHiYellow).Printf("\rSending email... %s", frames[i])
			i = (i + 1) % len(frames)
		}
	}
}
