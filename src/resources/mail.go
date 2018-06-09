package resources

import (
	"nomad/api/src/sendmail/themes"
	"github.com/matcornic/hermes"
	"net/smtp"
)

type Mail struct {
	SmtpLogin		string
	SmtpPassword	string
	SmtpHost		string
	SmtpPort		string
	From			string
	Theme			hermes.Theme
	Resources		*Resources
}

func (m Mail) Send(to string, subject string, body hermes.Email) error {
	h := hermes.Hermes{
		Theme: m.Theme,
		Product: hermes.Product{
			Name: "Nomad Space",
			Link: "https://nomad.space/",
			Logo: "http://nomad.space/wp-content/uploads/2018/02/logo-nomad.svg",
			Copyright: "Copyright © 2018 Nomad Space. All rights reserved.",
		},
	}

	emailBody, err := h.GenerateHTML(body)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.SmtpLogin, m.SmtpPassword, m.SmtpHost)

	msg := []byte(
		"MIME-version: 1.0;\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"To: "+to+"\r\n" +
			"From: "+m.From+"\r\n" +
			"Subject: "+subject+"\r\n" +
			"\r\n" +
			emailBody)
	err = smtp.SendMail(m.SmtpHost+":"+m.SmtpPort, auth, m.From, []string{to}, msg)

	return err
}

func (r *Resources) initMail() error {

	r.Mail = Mail{
		SmtpLogin:		r.Config.SmtpLogin,
		SmtpPassword:	r.Config.SmtpPassword,
		SmtpHost:		r.Config.SmtpHost,
		SmtpPort:		r.Config.SmtpPort,
		From:			r.Config.SendmailFrom,
		Theme:			new(themes.Flat),
		Resources:		r,
	}
	return nil
}