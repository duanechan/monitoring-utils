package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wneessen/go-mail"
)

type User struct {
	name  string
	email string
}

func (u User) String() string {
	return fmt.Sprintf("%s <%s>", u.name, u.email)
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

func parseFlags() User {
	recipientName := flag.String("name", "", "the name of the recipient")
	recipientEmail := flag.String("email", "", "the email of the recipient")

	flag.Parse()

	*recipientName = strings.TrimSpace(strings.ReplaceAll(*recipientName, "\r", ""))
	*recipientEmail = strings.TrimSpace(strings.ReplaceAll(*recipientEmail, "\r", ""))

	return User{name: *recipientName, email: *recipientEmail}
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("error: failed to load .env file: %s", err)
	// }
	message := mail.NewMsg()

	recipient := parseFlags()

	// ccName := "Samantha Grace Aniversario"
	// ccEmail := "samantha_aniversario@dlsu.edu.ph"
	// if err := message.AddCcFormat(ccName, ccEmail); err != nil {
	// 	log.Fatalf("error: failed to set To address: %s", err)
	// }

	if err := message.AddToFormat(recipient.name, recipient.email); err != nil {
		log.Fatalf("error: failed to set To address: %s", err)
	}

	if err := message.FromFormat("Duane Chan", "chan.duanematthew@gmail.com"); err != nil {
		log.Fatalf("error: failed to set From address: %s", err)
	}

	message.Subject("OfficeTimer Credentials for the Internship in Knowles Training Institute")
	message.SetImportance(mail.ImportanceUrgent)
	message.SetBodyString(
		mail.TypeTextHTML,
		fmt.Sprintf(
			template,
			recipient.name,
			recipient.email,
			recipient.email,
			"welcome1#",
			"Duane Matthew P. Chan",
			"chan.duanematthew@gmail.com",
			"chan.duanematthew@gmail.com",
		),
		// fmt.Sprintf(`
		// 	<b>Here is your <a href="">OfficeTimer</a> account.</b><br>
		// 	<br>
		// 	<b>Username:</b> <a href="mailto:%s">%s</a><br>
		// 	<b>Password: welcome1#</b><br>
		// 	<br>
		// 	--<br>
		// 	<em>Duane Matthew P. Chan</em><br>
		// 	Knowles IT Monitoring Team<br>
		// 	Email: <a href="mailto:chan.duanematthew@gmail.com">chan.duanematthew@gmail.com</a><br>
		// 	Visit us: <a href="https://www.knowlesti.ph">https://www.knowlesti.ph</a>
		// 	`,
		// 	*recipientEmail,
		// 	*recipientEmail,
		// ),
	)

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		log.Fatalf("error: SMTP_USER not set")
	}

	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		log.Fatalf("error: SMTP_PASS not set")
	}

	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtpUser),
		mail.WithPassword(smtpPass))

	if err != nil {
		log.Fatalf("error: failed to create mail client: %s", err)
	}

	done := make(chan bool)
	go showLoadingBar(done)

	if err := client.DialAndSend(message); err != nil {
		done <- true
		log.Fatalf("failed to send mail: %s", err)
	}

	done <- true

	fmt.Printf("\rCredentials email successfully sent to %s\n", recipient)
	// fmt.Printf("With CC to %s <%s>", ccName, ccEmail)
}

const template string = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>OfficeTimer Credentials for the Internship in Knowles Training Institute</title>
</head>
<body style="margin: 0; padding: 15px; background-color: #e9f1f7; font-family: Arial, sans-serif;">
  <table role="presentation" width="100%%" height="100%%" cellspacing="0" cellpadding="0" border="0">
    <tr>
      <td align="center" valign="middle">
        <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="background-color: white; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 10px 15px rgba(0, 0, 0, 0.05);">
          <!-- Logo Section -->
          <tr>
            <td align="center" valign="middle" style="padding: 20px 20px;">
              <a href="https://www.knowlesti.ph" target="_blank" style="display: inline-block;">
                <img src="https://knowlesti.ph/wp-content/uploads/2021/02/ph-LOGO-1.png" width="180" alt="Knowles Training Institute" style="border: 0; display: block;">
              </a>
            </td>
          </tr>

          <!-- Divider -->
          <tr>
            <td align="center" style="padding: 0px 40px;">
              <div style="height: 1px; background-color: #edf2f7;"></div>
            </td>
          </tr>

          <!-- Title Section -->
          <tr>
            <td align="center" style="padding: 20px 40px 0;">
              <h1 style="color: #1a365d; font-size: 24px; margin: 0; font-family: Arial, sans-serif;">Welcome, %s!</h1>
              <p style="color: #4a5568; font-size: 16px; padding: 10px; margin: 0; font-family: Arial, sans-serif;">Here is your OfficeTimer account.</p>
            </td>
          </tr>

          <!-- Credentials Section -->
          <tr>
            <td align="center" style="padding: 20px 40px;">
              <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="background-color: #f8fafc; border-radius: 8px; border: 1px solid #e2e8f0; width: 75%%;">
                <tr>
                  <td style="padding: 30px;">
                    <p style="margin: 0; color: #333; font-size: 16px; line-height: 1.6;">
                      <span style="color: #4a5568;">Username:</span> 
                      <strong><a href="mailto:%s" style="color: #2b6cb0; text-decoration: none;">%s</a></strong>
                    </p>
                    <p style="margin: 15px 0 0 0; color: #333; font-size: 16px; line-height: 1.6;">
                      <span style="color: #4a5568;">Password:</span> 
                      <strong style="color: #2d3748;">%s</strong>
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Login Button -->
          <tr>
            <td align="center" style="padding: 30px 40px;">
              <a href="https://www.officetimer.com/login/" style="background-color: #2b6cb0; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; font-weight: bold; display: inline-block;">Access OfficeTimer</a>
            </td>
          </tr>

          <!-- Divider -->
          <tr>
            <td align="center" style="padding: 0 40px;">
              <div style="height: 1px; background-color: #edf2f7;"></div>
            </td>
          </tr>

          <!-- Footer Section -->
          <tr>
            <td align="center" style="padding: 30px 40px;">
              <p style="margin: 0; color: #4a5568; font-size: 14px; line-height: 1.6;">
                <em style="color: #2d3748;">%s</em><br>
                <span style="color: #4a5568;">Knowles IT Monitoring Team</span><br>
                Email: <a href="mailto:%s" style="color: #2b6cb0; text-decoration: none;">%s</a><br>
                Visit us: <a href="https://www.knowlesti.ph" style="color: #2b6cb0; text-decoration: none;">knowlesti.ph</a>
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`
