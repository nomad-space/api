package templates

import "github.com/matcornic/hermes"

func Welcome(username string, link string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"Welcome to Nomad Space! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Nomad Space, please click here:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Confirm your account",
						Link:  link,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}