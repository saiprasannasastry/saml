package main

import (
	"encoding/xml"
)

//SAML Response with Signed Message & Assertion

type Decoder struct {
	XMLName   xml.Name
	Issuer    Issuer    `xml:"Issuer"`
	Status    Status    `xml:"Status"`
	Signature Signature `xml:"Signature,omitempty"`
	Assertion Assertion `xml:"Assertion"`
}

type Issuer struct {
	Text string `xml:",chardata"`
	//	Saml2  string `xml:"saml2,attr"`
	//	Format string `xml:"Format,attr"`
}

type Signature struct {
	//XMLName xml.Name
	//Id      string `xml:"Id,attr"`
	//SignedInfo     SignedInfo
	SignatureValue string  `xml:"SignatureValue"`
	KeyInfo        KeyInfo `xml:"KeyInfo"`
	//	DS   string   `xml:"xmlns:dsig,attr"`

}

type KeyInfo struct {
	Cert X509Data `xml:"X509Data"`
}
type X509Data struct {
	Text            string `xml:",chardata"`
	X509Certificate string `xml:"X509Certificate"`
}

//We dont need Status as if the user is not authorized the app, okta will throw him out
type Status struct {
	StatusCode StatusCode `xml:"StatusCode"`
}
type StatusCode struct {
	Value string `xml:"Value,attr"`
	Text  string `xml:",chardata"`
}

type Assertion struct {
	Conditions Conditions `xml:"Conditions"`
	Signature  Signature  `xml:"Signature,omitempty"`

	//	AuthnStatement      AuthnStatement      `xml:"AuthnStatement"`
	AttributeStatement AttributeStatement `xml:"AttributeStatement"`
}
type Conditions struct {
	NotBefore string `xml:"NotBefore,attr"`
	//This is what we should use to create session in ACM
	NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
	//this is our URI
	AudienceRestriction AudienceRestriction `xml:"AudienceRestriction"`
}

type AudienceRestriction struct {
	Audience string `xml:"Audience"`
}

type AttributeStatement struct {
	Attribute []Attribute `xml:"Attribute"`
}
type Attribute struct {
	Name           string           `xml:"Name,attr"`
	AttributeValue []AttributeValue `xml:"AttributeValue"`
}

type AttributeValue struct {
	Text string `xml:",chardata"`
}
