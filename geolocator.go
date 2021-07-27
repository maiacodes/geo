package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
)

type geoDB struct {
	cityDB *geoip2.Reader
	asnDB  *geoip2.Reader
}

func (g *geoDB) Setup() {
	c, err := geoip2.Open(os.Getenv("GEOLITE_CITY"))
	if err != nil {
		log.Fatal(err)
	}

	a, err := geoip2.Open(os.Getenv("GEOLITE_ASN"))
	if err != nil {
		log.Fatal(err)
	}

	g.cityDB = c
	g.asnDB = a
}

type geoInfo struct {
	City    string
	Country string
	ISP     string
	ASN     string
	Lat     string
	Long    string
}

func (g *geoDB) Query(ip net.IP) (info geoInfo) {
	city, err := g.cityDB.City(ip)
	if err != nil {
		fmt.Println(err)
		info.City = "Unknown"
		info.Country = "Unknown"
	} else {
		info.City = city.City.Names["en"]
		info.Country = city.Country.Names["en"]
		info.Lat = fmt.Sprint(city.Location.Latitude)
		info.Long = fmt.Sprint(city.Location.Longitude)
	}

	asn, err := g.asnDB.ASN(ip)
	if err != nil {
		fmt.Println(err)
		info.ISP = "Unknown"
		info.ASN = "Unknown"
	} else {
		info.ISP = asn.AutonomousSystemOrganization
		info.ASN = fmt.Sprintf("AS%v", asn.AutonomousSystemNumber)
	}

	return
}
