package trakpak

import (
	"encoding/base64"
	"encoding/xml"
	"time"
)

// Time - Marshalls time.Time to expected format.
type Time time.Time

var timeFormat = "2006-01-02 15:04:05"

//	"2016-10-21 18:53:53" as "2006-01-02T15:04:05Z07:00"

// MarshalXML - Marshals time to correct format.
func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(time.Time(t).Format(timeFormat), start)
	return nil
}

// Unmarshals time from correct format.
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parsedTime, err := time.Parse(timeFormat, v)
	if err != nil {
		return err
	}
	*t = Time(parsedTime)
	return nil

}

// AccessRequest - TODO
type AccessRequest struct {
	AccessLicenseNumber string `valid:"required,length(1|12)"`
	UserID              string `xml:"UserId" valid:"required,length(1|6)"`
	Password            string `valid:"required,length(1|20)"`
}

// Shipment - TODO
type Shipment struct {
	ShipmentDate        Time // <ShipmentDate>2013-06-24 10:32:44</ShipmentDate>
	Shipper             *Shipper
	Destination         *Destination
	ShipmentInformation *ShipmentInformation
}

// Shipper - TODO
type Shipper struct {
	ShipperCompany     string `valid:"required,length(1|35)"`
	ShipperAddress1    string `valid:"required,length(1|35)"`
	ShipperAddress2    string `valid:"required,length(1|35)"`
	ShipperAddress3    string `valid:"required,length(1|35)"`
	ShipperCity        string `valid:"required,length(1|20)"`
	ShipperCounty      string `valid:"length(1|35)"`
	ShipperCountryCode string `valid:"required,length(2|2)"`
	ShipperPostcode    string `valid:"required,length(1|10)"`
	ShipperContact     string `valid:"length(0|40)"`
	ShipperPhone       string `valid:"required,length(1|10)"`
	ShipperVat         string `valid:"required,length(0|17)"`
	ShipperEmail       string `valid:"email,length(1|40)"`
	ShipperReference   string `valid:"length(1|20)"`
	ShipperDept        string `valid:"length(1|17)"`
}

// Destination - TODO
type Destination struct {
	DestinationCompany     string `valid:"required,length(1|35)"`
	DestinationAddress1    string `valid:"required,length(1|35)"`
	DestinationAddress2    string `valid:"required,length(1|35)"`
	DestinationAddress3    string `valid:"length(10|35)"`
	DestinationCity        string `valid:"required,length(1|20)"`
	DestinationCounty      string `valid:"length(1|35)"`
	DestinationCountryCode string `valid:"required,length(2|2)"`
	DestinationPostCode    string `valid:"length(0|10)"`
	DestinationContact     string `valid:"required,length(1|40)"`
	DestinationPhone       string `valid:"length(0|20)"`
	DestinationVat         string `valid:"length(0|17)"`
	DestinationEmail       string `valid:"required,email,length(1|40)"`
}

// ShipmentInformation - TODO
type ShipmentInformation struct {
	Service            string `valid:"required,length(1|4)"`
	TotalPieces        int
	TotalWeight        float64
	WeightID           string `xml:"WeightId" valid:"length(1|1)"`
	Length             int
	Width              int
	Height             int
	Product            string `valid:"required,length(3|3)"`
	DescriptionOfGoods string `valid:"required,length(1|70)"`
	Value              float64
	ValueCurrency      string `valid:"required,length(3|3)"`
	Terms              string `valid:"length(3|3)"`
	LabelImageFormat   string `valid:"length(3|3)"`
	ItemInformation    []ItemInformation
}

// ItemInformation - TODO
type ItemInformation struct {
	ItemDescription string `valid:"required,length(1|255)"`
	ItemHsCode      string `valid:"length(1|10)"`
	ItemQuantity    int
	ItemValue       float64
	ItemSkuCode     string `valid:"length(1|10)"`
	ItemCOO         string `valid:"length(2|2)"`
}

/*
type Ack struct {
	xml.Name      `xml:"ack"`
	DateTimeStamp Time
}
*/

// ShipmentResponse - the wrapped response to a ShipmentBookingRequest
type ShipmentResponse struct {
	Hawb               string `valid:"required,length(1|12)`
	TrackingNumber     string `valid:"required,length(1|50)`
	TrackingUrl        string `valid:"required"`
	QuickTrackURL      string
	CarrierTrackingUrl string `valid:"required"`
	LabelImage         Label  `valid:"required"`
	LabelImageFormat   string `valid:"required,length(3|3)"`
}

type Label string

func (l *Label) Decode() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(*l))
}

// ShipmentBookingAck - wrapper for the the ShipmentBookingResponse
type ShipmentBookingAck struct {
	xml.Name         `xml:"ack"`
	DateTimeStamp    Time
	ShipmentResponse *ShipmentResponse
	ErrorData        string
}

type ShipmentVoidRequest struct {
	TrackingNumber string `valid:"required,length(1|50)"`
}

type ShipmentConfirmRequest struct {
	TrackingNumber string `valid:"required,length(1|50)"`
}
