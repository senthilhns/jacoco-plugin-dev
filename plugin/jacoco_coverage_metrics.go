package plugin

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Report struct {
	XMLName  xml.Name  `xml:"report"`
	Counters []Counter `xml:"counter"`
	Packages []Package `xml:"package"`
}

type Counter struct {
	Type    string `xml:"type,attr"`
	Missed  int    `xml:"missed,attr"`
	Covered int    `xml:"covered,attr"`
}

type Package struct {
	Name     string    `xml:"name,attr"`
	Counters []Counter `xml:"counter"`
}

type JacocoCoverageThresholds struct {
	InstructionCoverageThreshold string
	BranchCoverageThreshold      string
	LineCoverageThreshold        string
	ComplexityCoverageThreshold  int
	MethodCoverageThreshold      string
	ClassCoverageThreshold       string
}

func CalculatePercentage(covered, missed int) string {
	total := covered + missed
	if total == 0 {
		return "0%(0/0)"
	}
	percentage := (float64(covered) / float64(total)) * 100
	return fmt.Sprintf("%.2f%%(%d/%d)", percentage, covered, total)
}

func GetCounterValues(counters []Counter, counterType string) (int, int) {
	for _, counter := range counters {
		if counter.Type == counterType {
			return counter.Covered, counter.Missed
		}
	}
	return 0, 0
}

func CalculateCoverageMetrics(report Report) JacocoCoverageThresholds {

	instructionCoverage, instructionMiss := GetCounterValues(report.Counters, "INSTRUCTION")
	branchCoverage, branchMiss := GetCounterValues(report.Counters, "BRANCH")
	lineCoverage, lineMiss := GetCounterValues(report.Counters, "LINE")
	complexityCoverage, complexityMiss := GetCounterValues(report.Counters, "COMPLEXITY")
	methodCoverage, methodMiss := GetCounterValues(report.Counters, "METHOD")
	classCoverage, classMiss := GetCounterValues(report.Counters, "CLASS")

	return JacocoCoverageThresholds{
		InstructionCoverageThreshold: CalculatePercentage(instructionCoverage, instructionMiss),
		BranchCoverageThreshold:      CalculatePercentage(branchCoverage, branchMiss),
		LineCoverageThreshold:        CalculatePercentage(lineCoverage, lineMiss),
		ComplexityCoverageThreshold:  complexityCoverage + complexityMiss,
		MethodCoverageThreshold:      CalculatePercentage(methodCoverage, methodMiss),
		ClassCoverageThreshold:       CalculatePercentage(classCoverage, classMiss),
	}
}

func ParseXMLReport(filename string) Report {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening XML file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	var report Report
	err = xml.Unmarshal(data, &report)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}
	return report
}

func GetJacocoCoverageThresholds(completeXmlPath string) JacocoCoverageThresholds {
	report := ParseXMLReport(completeXmlPath)
	coverageThresholds := CalculateCoverageMetrics(report)

	fmt.Println("Coverage Metrics:")
	fmt.Printf("Instruction Coverage: %s\n", coverageThresholds.InstructionCoverageThreshold)
	fmt.Printf("Branch Coverage: %s\n", coverageThresholds.BranchCoverageThreshold)
	fmt.Printf("Line Coverage: %s\n", coverageThresholds.LineCoverageThreshold)
	fmt.Printf("Complexity Coverage: %d\n", coverageThresholds.ComplexityCoverageThreshold)
	fmt.Printf("Method Coverage: %s\n", coverageThresholds.MethodCoverageThreshold)
	fmt.Printf("Class Coverage: %s\n", coverageThresholds.ClassCoverageThreshold)

	return coverageThresholds
}
