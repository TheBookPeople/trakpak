package trakpak

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var shipper = Shipper{
	ShipperCompany:     "BOOK PEOPLE",
	ShipperAddress1:    "PARC MENAI",
	ShipperCity:        "BANGOR",
	ShipperCounty:      "GWYNEDD",
	ShipperCountryCode: "GB",
	ShipperPostcode:    "LL57 4FB",
}

func newClient() *Client {
	licenseNo := os.Getenv("TRAKPAK_LICENSE")
	userID := os.Getenv("TRAKPAK_USER")
	password := os.Getenv("TRAKPAK_PASSWORD")
	if licenseNo == "" {
		panic("TRAKPAK_LICENSE not set")
	}
	if userID == "" {
		panic("TRAKPAK_USER not set")
	}
	if licenseNo == "" {
		panic("TRAKPAK_PASSWORD not set")
	}

	c := NewClient(licenseNo, userID, password)
	c.TestMode = true

	return c
}

func TestLabels(t *testing.T) {

	c := newClient()

	tests := []struct {
		code                 string
		price, weight        float64
		depth, width, height int
		destination          Destination
	}{
		{
			"RSLW", 11.99, 1.606, 32, 24, 19,
			Destination{"", "c/o Sue Steed", "Carrer de la Valleta, 6.", "Ginestar", "Tarragona", "Catalunya", "ES", "43748", "Ben West", "", "", ""},
		},
		{
			"LPSV", 4.99, 1.115, 32, 24, 19,
			Destination{"", "c/o Dan Badge Manufacturing sdn.bhd", "26-S, Jalan Bidara 2/5,", "Taman Bidara", "Selayang, Batu Caves", "Selangor", "MY", "68100", "Sakila bin Mohd Zain", "", "", ""},
		},
		{
			"LPSV", 4.99, 1.115, 32, 24, 19,
			Destination{"", "5485 Byscane Lane", "", "", "Minnetonka", "Minnesota", "US", "55345-5603", "Mrs Cheryl Kerber", "", "", ""},
		},
		{
			"RSLW", 11.99, 1.606, 32, 24, 19,
			Destination{"", "Ruiterslaan 6", "", "", "Wijnegem", "", "BE", "2110", "Frederic Verhulst", "", "", ""},
		},
		{
			"LPSV", 4.99, 1.115, 32, 24, 19,
			Destination{"", "11 Pisani Court", "Golden Grove", "", "Adelaide", "South Australia", "AU", "5125", "Mrs D Richards", "", "", ""},
		},
	}

	for _, test := range tests {
		s := &Shipment{
			ShipmentDate: Time(time.Now()),
			Shipper:      &shipper,
			Destination:  &test.destination,
			ShipmentInformation: &ShipmentInformation{
				Service:            "TPWW",
				TotalPieces:        1,
				TotalWeight:        test.weight,
				WeightID:           "K",
				Width:              test.width,
				Height:             test.height,
				Length:             test.depth,
				DescriptionOfGoods: "Books",
				Value:              4.99,
				ValueCurrency:      "GBP",
				Terms:              "DDU", // FIXME: What is this???
				LabelImageFormat:   "PDF",
				ItemInformation: []ItemInformation{
					{"Book", "", 1, test.price, test.code, "GB"},
				},
			},
		}
		filename := "/tmp/" + test.destination.DestinationCountryCode + ".pdf"
		defer os.Remove(filename)
		if err := ship(c, s, filename); err != nil {
			t.Errorf("Failed on %s: %v", test.destination.DestinationCountryCode, err)
		}
	}
}

func ship(c *Client, s *Shipment, filename string) error {
	resp, err := c.BookShipment(s)
	if err != nil {
		return err
	}
	label, err := resp.LabelImage.Decode()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, label, 0666)
	if err != nil {
		return err
	}
	return nil
}
