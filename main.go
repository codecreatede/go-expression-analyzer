package main

/*

Author Gaurav Sablok
Universitat Potsdam
Date 2024-10-14


Golang estimation of the expression analyzer using go routines.
An implementation acros Golang using the pure structs and it express
normalizes for both the RNA-seq and the metagenomics-RNAseq.

*/


import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"log"
	"strings"
	"strconv"
)


func main () {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}

var samalignment string
var gtffile string

var rootCmd = &cobra.Command{
	Use: "analyze",
	Long: "Analyze the expression from the gene expression alignment",
	Run: analyzeFunc,
}

func init() {
	rootCmd.Flags().StringVarP(&samalignment, "samalignment", "s", "samalignment", "alignment file to be analyzed")
	rootCmd.Flags().StringVarP(&gtffile, "gtf file", "G", "gtf file for the analysis", "gtf analysis")
}


type geneStruct struct {
	gene string
	start int
	end int
}

type geneAlign struct {
	gene string
	position int
}


func analyzeFunc (cmd * cobra.Command, args[] string) {
	fOpen, err := os.Open(samalignment)
	if err != nil {
		log.Fatal(err)
	}


	// sorting the sam alignment
	geneAlignCapture := []geneAlign{}
	genomeAnnotationCapture := []geneStruct{}

	fRead := bufio.NewScanner(fOpen)
	for FRead.Scan() {
		line := fRead.Text()
		if strings.HasPrefix(string(line), "@") {
			continue
		}
		if ! strings.HasPrefix(string(line), "@") {
			geneAlignCapture = append(geneAlignCapture, geneAlign{
				gene: strings.Split(string(line), " ")[2],
				position: strings.Split(string(line), " ")[3],
			})
		}
	}


	// sorting annotations

	gOpen, err := os.Open(gtffile)
	if err != nil {
		log.Fatal(err)
	}
	gRead := bufio.NewScanner(gOpen)
	for gRead.Scan() {
		line := gRead.Text()
		if strings.HasPrefix(string(line), "#") {
			continue
		}
		if ! strings.Hasprefix(string(line), "#")  && strings.Split(string(line))[2] == "gene" {
			genomeAnnotationCapture = append(genomeAnnotationCapture, geneStruct{
        gene: strings.Split(string(line), " ")[0],
				start:strings.Split(string(line), " ")[3],
				end: strings.Split(string(line), " ")[4],
			 })
		}
	}

	// struct iterations and gene length and unqiue gene length estimation
	genenames := []string{}
	genelength := []int{}

	for i := range genomeAnnotationCapture {
		genenames = append(genenames, genomeAnnotationCapture[i].gene)
		start, _ := strconv.Atoi(geneAnnotationCapture[i].start)
		end, _ := strconv.Atoi(genomeAnnotationCapture[i].end)
		genelength = append(genelength, end-start)
	}

	unqiuegenes := make(map[string]string)

	for fRead.Scan() {
		line := fRead.Text()
			uniqueline, additionaline := strings.Split(string(line), " ")[2], strings.Split(string(line), " ")[4]
			uniquegenes[uniqueline] = additionaline
	}

	type geneExpression struct {
		gene string
		express int
	}

	// counting the specific struct for the gene gene aligned from the sam file and then appending the
	// number of the read alignments.

	geneExp := []geneExpression{}

	countAlignment := 0
	for i := range geneAlignCapture {
		for j := range uniquegenes {
			if j == geneAlignCapture[i].position {
				countAlignment++
        geneExp = append(geneExp, geneExpression{
					gene: j,
          express : countAlignment,
				})
			}
		}
	}

	// estimating the values

	geneRPK := []int{}
	for i := range geneExp {
		for j := range geneLength {
			if geneExp[i].gene == genenames[i] {
				intermediate := geneExp[i].express/(geneLength[i]/1000)
				geneRPK = append(geneRPK,intermediate)
			}
		}
	}

	var scale int

	func sum (arr []int) int {
		counter := 0
		for i := range indexarray {
			counter += indexarray[i]
		}
		return counter
	}

	scale := sum(geneRPK)/1000000

	type expressCount struct {
		gene string
		expressCount int
	}

	expression := []expressCount{}

	for i := range genenames {
		expression = append(expression, expressionCount{
			gene: genenames[i],
      expressCount: geneRPK[i]/scale
		})
	}

	file, err := os.Create("expression-metagenome.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for i := range expression {
		_, err := file.WriteString(expression[i].gene + expression[i].count)
	}

}
