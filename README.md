# **Credentials Helper**

A command-line utility program, written in Go, that sends emails containing OfficeTimer credentials to new interns. It simplifies email creation and sending by allowing users to input a CSV file and use a default email layout.
## **Demo**

![Credentials Helper Demo](https://i.imgur.com/colTSvk.gif)

![Email Template](https://i.imgur.com/1fQbCEM.png)


## **Environment Variables**

To run this project, you will need to add the following environment variables to your .env file:

### SMTP Auth

`SMTP_USER` The username for the SMTP server, typically your email address.

`SMTP_PASS` The app password for your SMTP account. Since Google no longer supports less secure apps, you need to generate an app password from [Google App Passwords](https://myaccount.google.com/apppasswords) if you're using Gmail.

### Sender/From

`SENDER_NAME` The name that will appear as the sender of the email.

`SENDER_EMAIL` The email address from which emails will be sent. This should match `SMTP_USER` if using a personal email service like Gmail.

### CC (Optional)

`CC_NAME` The name of the recipient who will be CC’d on the email. This person will receive a copy of the email but will not be the primary recipient.

`CC_EMAIL` The email address of the CC recipient. If left blank, no CC email will be sent.
## **Installation & Setup**

You only have to set up the environment variables, or your `.env` file.

***Note:** Remember to create your app password.* 

![.env Setup](https://i.imgur.com/YsSwhCa.png)

Once finished, open your terminal/command prompt and navigate to the directory where `credentials.exe` and `.env` are located.

![Get directory](https://i.imgur.com/riMF1is.png)

![Navigate with terminal](https://i.imgur.com/z9Vrq5M.png)

Now, you're ready to use the program!
## **Usage**

### Startup

Start the program in command prompt or Powershell by entering:

```cmd
credentials
```

or specify the path directly with...

```cmd
credentials -path "C:/path/to/recipients.csv"
```

### CSV Format

The CSV file containing the list of recipients has only two columns: name and email, as shown below:

    John Doe, johndoe62@gmail.com
    Mary Jane, maryjane242@gmail.com

***Note:** It should contain no headers, otherwise it will read it as an invalid email.*

### Email Report

Once the CSV file is read, it sends the credentials email to the records with a valid email address. Here is a CSV file with two valid emails and one invalid:

![Email Report 1](https://i.imgur.com/g5qzePA.png)

If an email is invalid, the program will display the **line numbers of that invalid email**.

![Email Report 2](https://i.imgur.com/ljCpKGP.png)

The purpose of this feature is once we know where the bad records are, we can extract and put them in another CSV file where we validate and repeat the process.
