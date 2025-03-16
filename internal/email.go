// Copyright © 2025 Duane Matthew P. Chan

package email

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	"github.com/wneessen/go-mail"
)

type Email struct {
	TemplateType TemplateType
	TemplateData map[string]interface{}
	To           User
	Config       EmailConfig
}

type TemplateType string

const (
	CredentialsTemplate TemplateType = "credentials"
	LateTemplate        TemplateType = "late"
)

// Checks if given email address is valid.
func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}(?:\.[a-zA-Z]{2,})?$`)
	return regex.MatchString(email)
}

// Template definitions
var templates = map[TemplateType]*template.Template{}

// Initialize templates
func init() {
	// Parse and store templates
	templates[CredentialsTemplate] = template.Must(template.New("credentials").Parse(credentialsHTML))
	templates[LateTemplate] = template.Must(template.New("late").Parse(lateHTML))
}

// Sends the email to the given recipient.
func (e Email) Send() error {
	if !IsValidEmail(e.To.Email) {
		return fmt.Errorf("invalid recipient email")
	}

	// Get template
	tmpl, exists := templates[e.TemplateType]
	if !exists {
		return fmt.Errorf("template not found: %s", e.TemplateType)
	}

	// Execute template with data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, e.TemplateData); err != nil {
		return err
	}

	// Create new email message
	message := mail.NewMsg()

	// Subject
	subjects := map[TemplateType]string{
		CredentialsTemplate: "OfficeTimer Credentials for the Internship in Knowles Training Institute",
		LateTemplate:        "Important Reminder for Late Interns",
	}
	message.Subject(subjects[e.TemplateType])

	// Importance
	message.SetImportance(mail.ImportanceUrgent)

	// CC
	if e.Config.CC.Exists() {
		if err := message.AddCcFormat(e.Config.CC.Name, e.Config.CC.Email); err != nil {
			return err
		}
	}

	// Recipient
	if err := message.AddToFormat(e.To.Name, e.To.Email); err != nil {
		return err
	}

	// Sender
	if err := message.FromFormat(e.Config.From.Name, e.Config.From.Email); err != nil {
		return err
	}

	// Email body
	message.SetBodyString(mail.TypeTextHTML, body.String())

	// SMTP server configuration
	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(e.Config.SMTPUser),
		mail.WithPassword(e.Config.SMTPPass))

	if err != nil {
		return err
	}

	// Send email
	if err := client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

// Template HTML as constants
const credentialsHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>OfficeTimer Credentials for the Internship in Knowles Training Institute</title>
</head>
<body style="margin: 0; padding: 15px; background-color: #e9f1f7; font-family: Arial, sans-serif;">
  <table role="presentation" width="100%" height="100%" cellspacing="0" cellpadding="0" border="0">
    <tr>
      <td align="center" valign="middle">
        <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="background-color: white; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 10px 15px rgba(0, 0, 0, 0.05);">
          <!-- Logo Section -->
          <tr>
            <td align="center" valign="middle" style="padding: 20px 20px;">
              <a href="https://www.knowlesti.sg" target="_blank" style="display: inline-block;">
                <img src="https://i.imgur.com/Q9nLEZA.png" width="200" alt="Knowles Training Institute" style="border: 0; display: block;">
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
              <h1 style="color: #1a365d; font-size: 24px; margin: 0; font-family: Arial, sans-serif;">Welcome, {{ .Name }}!</h1>
              <p style="color: #4a5568; font-size: 16px; padding: 10px; margin: 0; font-family: Arial, sans-serif;">Here is your OfficeTimer account.</p>
            </td>
          </tr>

          <!-- Credentials Section -->
          <tr>
            <td align="center" style="padding: 20px 40px;">
              <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="background-color: #f8fafc; border-radius: 8px; border: 1px solid #e2e8f0; width: 90%;">
                <tr>
                  <td style="padding: 30px;">
                    <p style="margin: 0; color: #333; font-size: 16px; line-height: 1.6;">
                      <span style="color: #4a5568;">Username:</span> 
                      <strong><a href="mailto:{{ .Username }}" style="color: #2b6cb0; text-decoration: none;">{{ .Username }}</a></strong>
                    </p>
                    <p style="margin: 15px 0 0 0; color: #333; font-size: 16px; line-height: 1.6;">
                      <span style="color: #4a5568;">Password:</span> 
                      <strong style="color: #2d3748;">{{ .Password }}</strong>
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
                <em style="color: #2d3748;">{{ .SenderName }}</em><br>
                <span style="color: #4a5568;">Knowles IT Monitoring Team</span><br>
                Email: <a href="mailto:{{ .SenderEmail }}" style="color: #2b6cb0; text-decoration: none;">{{ .SenderEmail }}</a><br>
                Visit us: <a href="https://www.philippines.knowlesti.com" style="color: #2b6cb0; text-decoration: none;">philippines.knowlesti.com</a>
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>
`

const lateHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Important Reminder for Late Interns</title>
</head>
<body style="margin: 0; padding: 15px; background-color: #e9f1f7; font-family: Arial, sans-serif;">
  <table role="presentation" width="100%" height="100%" cellspacing="0" cellpadding="0" border="0">
    <tr>
      <td align="center" valign="middle">
        <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="background-color: white; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 10px 15px rgba(0, 0, 0, 0.05);">
          <!-- Logo Section -->
          <tr>
            <td align="center" valign="middle" style="padding: 20px 20px;">
              <a href="https://www.knowlesti.sg" target="_blank" style="display: inline-block;">
                <img src="https://i.imgur.com/Q9nLEZA.png" width="200" alt="Knowles Training Institute" style="border: 0; display: block;">
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
              <h1 style="color: #1a365d; font-size: 24px; margin: 0; font-family: Arial, sans-serif;">Important Reminder for Late Interns</h1>
              <p style="color: #4a5568; font-size: 16px; padding: 10px; margin: 0; font-family: Arial, sans-serif; text-align: justify;">
              <br><br>
              Dear Intern,<br><br>
              We hope this message finds you well. As you know, punctuality is an essential aspect of professionalism and contributes significantly to the success of any workplace. We understand that unforeseen circumstances may sometimes cause delays, but it is crucial to prioritize timeliness in your internship experience. 
              <br><br>
              We kindly remind all interns who have been late to take this matter seriously and make the necessary adjustments to ensure your punctuality moving forward. Remember, being on time not only demonstrates your commitment and respect for your work but also allows you to maximize your learning opportunities and contribute effectively to the team. 
              <br><br>
              To help you improve your punctuality, we suggest the following:
              </p>
              <ul style="color: #4a5568; font-size: 16px; padding: 10px; margin: 0; font-family: Arial, sans-serif; text-align: justify;">
              <li><strong>Plan ahead:</strong> Set your alarm clock early enough to provide ample time for your morning routine and commute. Consider any potential traffic or public transportation delays.</li>
              <br><br>
              <li><strong>Prepare in advance:</strong> Organize your essentials, such as your work bag and necessary documents, the night before to avoid last-minute rushes or forgotten items.</li>
              <br><br>
              <li><strong>Communicate proactively:</strong> If you encounter an unexpected situation that may cause tardiness, immediately notify your supervisor or the appropriate person. Prompt communication demonstrates responsibility and enables your team to plan accordingly.</li>
              <br><br>
              <li><strong>Seek support:</strong> If you struggle with punctuality, don't hesitate to seek guidance from your mentor, supervisor, or colleagues. They can provide valuable advice or resources to help you manage your time effectively.</li>
              </ul>
              <p style="color: #4a5568; font-size: 16px; padding: 10px; margin: 0; font-family: Arial, sans-serif; text-align: justify;">
              Please remember that your time with us is a valuable learning experience, and developing strong professional habits, such as punctuality, will greatly benefit your future career endeavors. 
              <br><br>
              We believe in your potential and are confident that you can make the necessary adjustments to improve your timeliness. If you have any questions or need further assistance, don't hesitate to reach out to your supervisor or the intern coordinator. 
              <br><br>
              Thank you for your attention, and we look forward to your continued growth and success during your internship. 
              <br><br>
              Best regards,<br>
              Monitoring Team<br>
              </p>
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
                <em style="color: #2d3748;">{{ .SenderName }}</em><br>
                <span style="color: #4a5568;">Knowles IT Monitoring Team</span><br>
                Email: <a href="mailto:{{ .SenderEmail }}" style="color: #2b6cb0; text-decoration: none;">{{ .SenderEmail }}</a><br>
                Visit us: <a href="https://www.philippines.knowlesti.com" style="color: #2b6cb0; text-decoration: none;">philippines.knowlesti.com</a>
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>
`
