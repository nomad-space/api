package templates

import (
	"github.com/matcornic/hermes"
	"nomad/api/src/models"
)

func Processed(booking models.Booking) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: booking.FirstName + " " + booking.LastName,
			Intros: []string{
				"Your order " + booking.Id.Hex() + "has been processed successfully.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Firstname", Value: booking.FirstName},
				{Key: "Lastname", Value: booking.LastName},
				{Key: "Hotel", Value: booking.Hotel[0].Name},
				{Key: "Address", Value: booking.Hotel[0].Address},
				{Key: "Location", Value: booking.Location[0].Country + ", " + booking.Location[0].City},
				{Key: "CheckIn", Value: booking.CheckIn.Format("2006-01-02")},
				{Key: "CheckOut", Value: booking.CheckOut.Format("2006-01-02")},
			},
			Actions: []hermes.Action{
				{
					Instructions: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua:",
					Button: hermes.Button{
						Text: "Go to Site",
						Link: "https://mvp.nomad.space",
					},
				},
			},
		},
	}
}