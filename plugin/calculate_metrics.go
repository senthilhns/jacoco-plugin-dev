package plugin

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Define Go structs to match the XML structure
type Report struct {
	XMLName     xml.Name      `xml:"report"`
	Name        string        `xml:"name,attr"`
	SessionInfo []SessionInfo `xml:"sessioninfo"`
	Packages    []Package     `xml:"package"`
	Counters    []Counter     `xml:"counter"`
}

type SessionInfo struct {
	ID    string `xml:"id,attr"`
	Start int64  `xml:"start,attr"`
	Dump  int64  `xml:"dump,attr"`
}

type Package struct {
	Name     string    `xml:"name,attr"`
	Classes  []Class   `xml:"class"`
	Counters []Counter `xml:"counter"`
}

type Class struct {
	Name           string    `xml:"name,attr"`
	SourceFileName string    `xml:"sourcefilename,attr"`
	Methods        []Method  `xml:"method"`
	Counters       []Counter `xml:"counter"`
}

type Method struct {
	Name     string    `xml:"name,attr"`
	Desc     string    `xml:"desc,attr"`
	Line     int       `xml:"line,attr"`
	Counters []Counter `xml:"counter"`
}

type Counter struct {
	Type    string `xml:"type,attr"`
	Missed  int    `xml:"missed,attr"`
	Covered int    `xml:"covered,attr"`
}

func AnalyzeJacocoXml(completeXmlPath string) {
	file, err := os.Open(completeXmlPath)
	if err != nil {
		log.Fatalf("Error opening XML file: %v", err)
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	// Unmarshal the XML data into the Report struct
	var report Report
	err = xml.Unmarshal(data, &report)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	// Print the parsed data
	fmt.Printf("Report Name: %s\n", report.Name)
	for _, session := range report.SessionInfo {
		fmt.Printf("Session ID: %s, Start: %d, Dump: %d\n", session.ID, session.Start, session.Dump)
	}

	for _, pkg := range report.Packages {
		fmt.Printf("Package: %s\n", pkg.Name)
		for _, class := range pkg.Classes {
			fmt.Printf("  Class: %s, Source: %s\n", class.Name, class.SourceFileName)
			for _, method := range class.Methods {
				fmt.Printf("    Method: %s, Description: %s, Line: %d\n", method.Name, method.Desc, method.Line)
				for _, counter := range method.Counters {
					fmt.Printf("      Counter Type: %s, Missed: %d, Covered: %d\n", counter.Type, counter.Missed, counter.Covered)
				}
			}
		}
	}
}
