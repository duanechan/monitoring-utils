package main

import (
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
	clearTerminal()
	header()

	path := flag.String("path", "", "the file path to the list of recipients")
	flag.Parse()

	recipients, err := getRecipients(*path)
	if err != nil {
		log.Fatalf("error: failed to read recipient data: %s", err)
	}

	var wg sync.WaitGroup

	for _, r := range recipients {
		wg.Add(1)
		go func(r credentials.User) {
			defer wg.Done()

			done := make(chan bool)
			go showLoadingBar(done)
			credentials.SendEmail(r)
			done <- true
			time.Sleep(100 * time.Millisecond)

			close(done)
		}(r)
	}

	wg.Wait()

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

// func parseFlags() []credentials.User {
// 	recipients, err := getRecipients(*path)
// 	if err != nil {
// 		log.Fatalf("error: failed to read recipients data: %s", err)
// 	}

// 	recipientName := flag.String("rname", "", "the name of the recipient")
// 	recipientEmail := flag.String("remail", "", "the email of the recipient")

// 	ccName := flag.String("cname", "", "the name of the cc recipient")
// 	ccEmail := flag.String("cemail", "", "the email of the cc recipient")

// 	flag.Parse()

// 	recipient = credentials.User{
// 		Name:  strings.TrimSpace(strings.ReplaceAll(*recipientName, "\r", "")),
// 		Email: strings.TrimSpace(strings.ReplaceAll(*recipientEmail, "\r", "")),
// 	}

// 	if *ccName != "" && *ccEmail != "" {
// 		cc = credentials.User{
// 			Name:  strings.TrimSpace(strings.ReplaceAll(*ccName, "\r", "")),
// 			Email: strings.TrimSpace(strings.ReplaceAll(*ccEmail, "\r", "")),
// 		}
// 	}

// 	return recipients
// }

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
			fmt.Printf("\rSending email... %s", frames[i])
			i = (i + 1) % len(frames)
		}
	}
}
