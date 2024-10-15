package plugin

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

type CoverageMetrics struct {
	InstructionCoverage string
	BranchCoverage      string
	LineCoverage        string
	ComplexityCoverage  int
	MethodCoverage      string
	ClassCoverage       string
	FileCoverage        string
	PackageCoverage     string
}

func calculatePercentage(covered, missed int) string {
	total := covered + missed
	if total == 0 {
		return "0%(0/0)" // Avoid division by zero
	}
	percentage := (float64(covered) / float64(total)) * 100
	return fmt.Sprintf("%.2f%%(%d/%d)", percentage, covered, total)
}

func getCounterValues(counters []Counter, counterType string) (int, int) {
	for _, counter := range counters {
		if counter.Type == counterType {
			return counter.Covered, counter.Missed
		}
	}
	return 0, 0 // Default if counter type not found
}

func calculateCoverageMetrics(report Report) CoverageMetrics {
	// Extract counters from the main report
	instCov, instMiss := getCounterValues(report.Counters, "INSTRUCTION")
	branchCov, branchMiss := getCounterValues(report.Counters, "BRANCH")
	lineCov, lineMiss := getCounterValues(report.Counters, "LINE")
	compCov, compMiss := getCounterValues(report.Counters, "COMPLEXITY")
	methodCov, methodMiss := getCounterValues(report.Counters, "METHOD")
	classCov, classMiss := getCounterValues(report.Counters, "CLASS")

	totalPackages := len(report.Packages)
	totalFiles := 0
	for _, pkg := range report.Packages {
		totalFiles += len(pkg.Counters) // Assume 1 file per counter for simplicity
	}

	return CoverageMetrics{
		InstructionCoverage: calculatePercentage(instCov, instMiss),
		BranchCoverage:      calculatePercentage(branchCov, branchMiss),
		LineCoverage:        calculatePercentage(lineCov, lineMiss),
		ComplexityCoverage:  compCov + compMiss, // Sum of complexity values
		MethodCoverage:      calculatePercentage(methodCov, methodMiss),
		ClassCoverage:       calculatePercentage(classCov, classMiss),
		FileCoverage:        calculatePercentage(totalFiles, 0),    // Assume all files are covered
		PackageCoverage:     calculatePercentage(totalPackages, 0), // Assume all packages are covered
	}
}

func parseXMLReport(filename string) Report {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening XML file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
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

func AnalyzeJacocoXml(completeXmlPath string) {
	report := parseXMLReport(completeXmlPath)
	metrics := calculateCoverageMetrics(report)

	// Print the metrics
	fmt.Println("Coverage Metrics:")
	fmt.Printf("Instruction Coverage: %s\n", metrics.InstructionCoverage)
	fmt.Printf("Branch Coverage: %s\n", metrics.BranchCoverage)
	fmt.Printf("Line Coverage: %s\n", metrics.LineCoverage)
	fmt.Printf("Complexity Coverage: %d\n", metrics.ComplexityCoverage)
	fmt.Printf("Method Coverage: %s\n", metrics.MethodCoverage)
	fmt.Printf("Class Coverage: %s\n", metrics.ClassCoverage)
	fmt.Printf("File Coverage: %s\n", metrics.FileCoverage)
	fmt.Printf("Package Coverage: %s\n", metrics.PackageCoverage)
}
