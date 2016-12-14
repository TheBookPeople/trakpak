package trakpak

// Client for P2P Mailing Ltd. Trak-pak service.
// See http://p2pmailing.co.uk/services/trak-pak/

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURL = "http://trakpak.co.uk/API/" // FIXME: Non-SSL???

// NewClient - returns a new TrakPak API client
func NewClient(licenseNo, userID, password string) *Client {
	return &Client{
		AccessRequest: &AccessRequest{
			AccessLicenseNumber: licenseNo,
			UserID:              userID,
			Password:            password,
		},
	}
}

// Client - a TrakPak API client
type Client struct {
	*AccessRequest
	TestMode bool
}

func parseShipmentResponse(r io.Reader) (*ShipmentResponse, error) {

	var ack ShipmentBookingAck

	if err := xml.NewDecoder(r).Decode(&ack); err != nil {
		return nil, err
	}

	if ack.ShipmentResponse == nil {
		return nil, fmt.Errorf(ack.ErrorData)
	}
	return ack.ShipmentResponse, nil
}

// BookShipment - Books a shipment. PDF Label is included in return.
func (c *Client) BookShipment(s *Shipment) (*ShipmentResponse, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	err := enc.Encode(c.AccessRequest)
	if err != nil {
		log.Println("Problem encoding access request: ", c.AccessRequest)
		return nil, err
	}
	err = enc.Encode(s)
	if err != nil {
		log.Println("Problem encoding shipment: ", s)
		return nil, err
	}
	log.Println(buf.String())
	url := baseURL + "?command=create"
	if c.TestMode {
		url += "&testMode=1"
	}
	log.Println(url)
	resp, err := http.Post(url, "text/xml", &buf)
	if err != nil {
		log.Println("Problem posting to:", url)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	//	log.Println(string(b))
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Shipment booking request returned: %v", resp.Status)
	}
	shipmentResponse, err := parseShipmentResponse(bytes.NewReader(b))
	return shipmentResponse, err
}

// VoidShipment - TODO
func (c *Client) VoidShipment(trackingNo string) error {
	return nil
}

/*// ConfirmShipment - TODO
func (c *Client) ConfirmShipment(trackingNo string) (*ShipmentConfirmResponse, error) {
	return nil, nil
}*/
