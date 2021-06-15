package handlers

import (
	service "chatapp/internal/app/services"
	vo "chatapp/internal/app/valobjects"
	"chatapp/internal/lib/util"
	"log"
	"net/http"
	"os"
)

// Handles new user Registration
func RegistrationHandler(rw http.ResponseWriter, req *http.Request) {

	// pg := &Page{
	// 	Title: "An Example",
	// }

	if req.Method == "POST" {
		// Read the form data
		req.ParseForm()
		name, _ := vo.NewName(req.Form.Get("name"))
		username, _ := vo.NewEmail(req.Form.Get("uname"))
		password := req.Form.Get("passwd")
		// create the DTO
		details := &service.RegistrationDetailsDTO{
			Name:     name,
			Username: username,
			Password: password,
		}
		// RegistrationService will handle the authentication
		registrationService := service.NewRegistrationService(details)
		err := registrationService.Execute()
		if err != nil {
			log.Println(err)
		}
		// Load email settings (Gmail)
		sendVerificationEmail(details)

		// This point is a success. Redirect to welcome screen
		rw.Header()["Location"] = []string{"/verify"}
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func sendVerificationEmail(dto *service.RegistrationDetailsDTO) {
	// Load SMTP settings
	smtpSettings := &util.SmtpSettings{
		Username: os.Getenv("SMTP_GMAIL_USERNAME"),
		Password: os.Getenv("SMTP_GMAIL_PASSWORD"),
		Server:   os.Getenv("SMTP_GMAIL_SERVER"),
		Port:     util.Environment().GetEnv("SMTP_PORT_TLS", "0").AsInt(),
	}
	// Compose the email
	message := &util.EmailMessage{
		From:    "support@chatapp.com",
		To:      []string{dto.Username.String()},
		Subject: "Email Verification",
	}
	// Send the email
	emailer := util.NewEmailUtil("verify_email.html")
	emailer.Settings(smtpSettings)
	emailer.Message(message)
	emailer.Send(dto.Name)
}
