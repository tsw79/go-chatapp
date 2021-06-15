package util

import (
	"bytes"
	vo "chatapp/internal/app/valobjects"
	"chatapp/internal/lib/config"
	"fmt"
	"net/smtp"
	"strconv"
	"text/template"
)

type EmailUtil struct {
	message  *EmailMessage
	settings *SmtpSettings
}

type EmailMessage struct {
	From    string
	Subject string
	Body    string
	To      []string
}

type SmtpSettings struct {
	Username string
	Password string
	Server   string
	Port     int
}

var t *template.Template

// Set up the template
func (this *EmailUtil) initTemplate(name string) {
	var err error
	var templatePath = config.GetInstance().Url.Templates + "/email/" + name
	t, err = template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
}

// Creates a new instance
func NewEmailUtil(templName string) *EmailUtil {
	this := EmailUtil{}
	this.initTemplate(templName)
	return &this
}

// Setter for SmtpSettings
func (this *EmailUtil) Settings(settings *SmtpSettings) *EmailUtil {
	this.settings = settings
	return this
}

// Setter for EmailUtil
func (this *EmailUtil) Message(message *EmailMessage) *EmailUtil {
	this.message = message
	return this
}

// Sends the email
func (this *EmailUtil) Send(recipient vo.Name) {
	// body holds the entire data of the email
	var body bytes.Buffer
	// Set the MIME version and content type for proper HTML rendering
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", this.message.Subject, headers)))
	// Parse the template and the data struct (object)
	err := t.Execute(&body, struct {
		Name string
	}{
		Name: recipient.String(),
	})
	if err != nil {
		panic(err)
	}
	// Smtp authentication
	auth := smtp.PlainAuth("",
		this.settings.Username,
		this.settings.Password,
		this.settings.Server,
	)
	addr := this.settings.Server + ":" + strconv.Itoa(this.settings.Port)
	// Send mail
	smtp.SendMail(addr, auth, this.message.From, this.message.To, body.Bytes())
}
