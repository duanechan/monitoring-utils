package credentials

import (
	"fmt"
	"log"

	"github.com/wneessen/go-mail"
)

func SendEmail(recipient User) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("error: failed to load email config: %s", err)
	}

	message := mail.NewMsg()

	if config.CC.Exists() {
		if err := message.AddCcFormat(config.CC.Name, config.CC.Email); err != nil {
			log.Fatalf("error: failed to set CC address: %s", err)
		}
	}

	if err := message.AddToFormat(recipient.Name, recipient.Email); err != nil {
		log.Fatalf("error: failed to set To address: %s", err)
	}

	if err := message.FromFormat(config.From.Name, config.From.Email); err != nil {
		log.Fatalf("error: failed to set From address: %s", err)
	}

	message.Subject("OfficeTimer Credentials for the Internship in Knowles Training Institute")
	message.SetImportance(mail.ImportanceUrgent)
	message.SetBodyString(
		mail.TypeTextHTML,
		fmt.Sprintf(
			Template,
			recipient.Name,
			recipient.Email,
			recipient.Email,
			"welcome1#",
			config.From.Name,
			config.From.Email,
			config.From.Email,
		),
	)

	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.SMTPUser),
		mail.WithPassword(config.SMTPPass))

	if err != nil {
		log.Fatalf("error: failed to create mail client: %s", err)
	}

	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}

	fmt.Printf("\rCredentials email successfully sent to %s\n", recipient)

	if config.CC.Exists() {
		fmt.Printf("With CC to %s <%s>", config.CC.Name, config.CC.Email)
	}
}
