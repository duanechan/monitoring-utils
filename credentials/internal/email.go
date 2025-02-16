package credentials

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

func SendEmail(recipient User) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	message := mail.NewMsg()

	if config.CC.Exists() {
		if err := message.AddCcFormat(config.CC.Name, config.CC.Email); err != nil {
			return err
		}
	}

	if err := message.AddToFormat(recipient.Name, recipient.Email); err != nil {
		return err
	}

	if err := message.FromFormat(config.From.Name, config.From.Email); err != nil {
		return err
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
		return err
	}

	if err := client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

const Template string = `<!DOCTYPE html>
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
              <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="background-color: #f8fafc; border-radius: 8px; border: 1px solid #e2e8f0; width: 90%%;">
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
