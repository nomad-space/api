package templates

import (
	"github.com/matcornic/hermes"
	"nomad/api/src/models"
)

func Processed(booking models.Booking) hermes.Email {

	var outros []string
	var intros string

	switch booking.Status {
	case models.BOOKING_STATUS_NEW:
		intros = "Your order " + booking.Id.Hex() + "has been successfully processed."
		outros = []string{"Our website does not charge any fees to the travelers and hotels.  We provide a direct link between you and any hotel you choose. In addition, we check the availability of cashback on average 15% depending on the price policy of the hotel.",
					"Our system requires up to 72 hours to check the booking terms and conditions, as well as the availability of the cashback. During this time we will send you a notice with details of your booking or advise you if the booking is not available.",
					"If the booking is available, payment will be made in accordance with the payment rules of the hotel you have chosen. In case of cancellation or if the booking terms and conditions are not suitable for you, we will offer you other closest hotels in the location."}
	case models.BOOKING_STATUS_REFUSAL:
		intros = "Unfortunately, booking of the hotel you have chosen is not available now."
		outros = []string{"We would be happy to offer you other closest hotels in the location. Alternatively, please make a repeated search for other hotels in our system."}
	case models.BOOKING_STATUS_CONSENT:
		intros = "Your order " + booking.Id.Hex() + "has been confirmed."
		outros = []string{"The cashback you can get in case you book this hotel is 15%.",
					"If you cancel your booking later, refund of the payment shall be regulated by the payment policy< of the hotel.",
					"If you accept the proposed terms and conditions, please complete your booking.",
					"In case of cancellation or if the booking terms and conditions are not suitable for you, we would be happy to offer you other closest hotels in the location. Alternatively, please make a repeated search for other hotels in our system."}
	}


	return hermes.Email{
		Body: hermes.Body{
			Name: booking.FirstName + " " + booking.LastName,
			Intros: []string{
				intros,
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
			Outros: outros,
		},
	}
}