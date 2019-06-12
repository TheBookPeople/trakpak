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

	"github.com/pkg/errors"
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
	Verbose  bool
}

func parseShipmentResponse(r io.Reader) (*ShipmentResponse, error) {

	var ack ShipmentBookingAck

	if err := xml.NewDecoder(r).Decode(&ack); err != nil {
		return nil, errors.Wrap(err, "Problem decoding TrakPak shipment booking acknowledgement")
	}

	if ack.ShipmentResponse == nil {
		return nil, errors.Errorf("Error: Trakpak shipment reponse was nil: %s", ack.ErrorData)
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
		return nil, errors.Wrapf(err, "Problem encoding TrakPak access request: %v", c.AccessRequest)
	}
	err = enc.Encode(s)
	if err != nil {
		return nil, errors.Wrapf(err, "Problem encoding TrakPak shipment: %v", s)
	}
	log.Println(buf.String())
	url := baseURL + "?command=create"
	if c.TestMode {
		url += "&testMode=1"
	}
	log.Println(url)
	resp, err := http.Post(url, "text/xml", &buf)
	if err != nil {
		return nil, errors.Wrapf(err, "Probem booking TrakPak shipment; failed to POST to %s", url)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if c.Verbose {
		log.Println(string(b))
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Shipment booking request returned: %v", resp.Status)
	}
	return parseShipmentResponse(bytes.NewReader(b))
}

// VoidShipment - TODO
func (c *Client) VoidShipment(trackingNo string) error {
	return nil
}

/*// ConfirmShipment - TODO
func (c *Client) ConfirmShipment(trackingNo string) (*ShipmentConfirmResponse, error) {
	return nil, nil
}*/
