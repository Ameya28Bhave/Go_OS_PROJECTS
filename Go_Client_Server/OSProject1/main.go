package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// **********************************

// Create the message , coordinator and worker struct

type worker struct {
	id     int
	job    chan []byte
	result chan []byte
}

type message struct {
	Fname string
	Start int
	End   int
}

type response struct {
	Value  uint64
	Prefix string
	Suffix string
	Error  string
}

// **********************************

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func cleanChunk(chunk string) string {
	//Check for trailing whitespace and Replace all whitespace by space
	ts := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	chunk_ := ts.ReplaceAllString(chunk, " ")

	//Check for trailing whitespace and Replace all whitespace by space
	ds := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	chunk_ = ds.ReplaceAllString(chunk_, " ")

	return chunk_
}

// **********************************

func coordinator(fname string, M float64) {
	var resp = make([][]byte, int(M))
	var msg *message
	var frg map[string]interface{}
	var temp = ""
	var continue_ bool = false
	//var count uint64

	wg := new(sync.WaitGroup)

	job_frag := make(chan []byte, int(M)*10)
	result_frag := make(chan []byte, int(M)*10)

	sum := uint64(0)

	// Read the file and create M chunks of the file

	dat, err := os.ReadFile(fname)
	check(err)
	fmt.Print(string(dat))

	f, err := os.Open(fname)
	check(err)
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n Splitting to %f pieces.\n", M)

	// Divide into M chunks
	frag := math.Floor(float64(fi.Size()-1) / M)

	//	Get the last extra chunks to add to the last thread
	last := int64(fi.Size()) - (int64(frag) * int64(M))

	// Call the thread creator and pass the filename with starting and ending byte
	fsize := 0

	fmt.Println("Total File Size: ", fi.Size())
	for i := 0; i < int(M); i++ {

		if i == int(M)-1 {
			msg = &message{
				fname,
				fsize,
				fsize + int(frag) + int(last) - 1,
			}
		} else {
			msg = &message{
				fname,
				fsize,
				fsize + int(frag) - 1,
			}
		}

		// Create an object to pass the data frag as JSON using Marshal

		FragFile, _ := json.Marshal(msg)
		fmt.Printf("Following JSON sent to Worker %d: %s\n", i, FragFile)
		job_frag <- FragFile
		wrkr := &worker{i, job_frag, result_frag}

		go wrkr.compute(wg)

		wg.Add(1)

		resp[i] = <-result_frag
		fsize += int(frag)
	}

	for index, frag := range resp {
		err := json.Unmarshal(frag, &frg)
		check(err)
		fmt.Printf("Worker %d responded with : Value(%v) , Prefix(%v), Suffix(%v), Error(%v)\n", index, frg["Value"], frg["Prefix"], frg["Suffix"], frg["Error"])
		if frg["Error"].(string) == "" {
			if frg["Prefix"].(string) != "" && frg["Suffix"].(string) != "" {
				// CASE 1
				//count_pre_suf = 3
				temp = fmt.Sprintf("%s%s", temp, frg["Prefix"].(string))
				newValue, _ := strconv.ParseInt(temp, 10, 64)
				sum = (sum + uint64(newValue))
				temp = frg["Suffix"].(string)
				continue_ = true

			} else if frg["Suffix"].(string) != "" {
				// CASE 2
				//count_pre_suf = 2
				if !continue_ {
					temp = frg["Suffix"].(string)
					continue_ = true
				} else {
					a := fmt.Sprintf("%s%s", frg["Prefix"].(string), temp)
					res, _ := strconv.ParseUint(a, 10, 64)
					sum = (sum + res)
					temp = frg["Suffix"].(string)
				}
			} else if frg["Prefix"].(string) != "" {
				// CASE 3
				//count_pre_suf = 2
				if !continue_ {
					newValue, _ := strconv.ParseInt(frg["Prefix"].(string), 10, 64)
					sum = sum + uint64(newValue)
				} else {
					temp = fmt.Sprintf("%s%s", temp, frg["Prefix"].(string))
					newValue, _ := strconv.ParseInt(temp, 10, 64)
					sum = (sum + uint64(newValue))
					temp = ""
					continue_ = false
				}
			} else {
				continue_ = false
			}
		} else {
			if !continue_ {
				temp = frg["Error"].(string)
				continue_ = true
			} else {
				temp = fmt.Sprintf("%s%s", temp, frg["Error"].(string))
				//fmt.Println("Chunk Merged", temp)
			}
		}

		//Compute Sum
		if temp != "" && !continue_ {
			newValue, _ := strconv.ParseInt(temp, 10, 64)
			fmt.Printf("When Temp-> sum + value + newvalue: %d + %d + %d\n", sum, uint64(frg["Value"].(float64)), uint64(newValue))
			sum = (sum + uint64(frg["Value"].(float64)) + uint64(newValue))
			temp = ""
			fmt.Printf(" sum + value: %d \n", sum)
		} else {
			fmt.Printf("When no Temp-> sum + value: %d + %d\n", sum, uint64(frg["Value"].(float64)))
			sum = (sum + uint64(frg["Value"].(float64)))
			fmt.Printf(" sum + value: %d \n", sum)
		}

	} // End of loop
	er := json.Unmarshal(resp[len(resp)-1], &frg)
	check(er)

	//lastIndex := len(resp) - 1
	temp = frg["Suffix"].(string)
	fmt.Println("Last Suffix", temp)
	newValue, _ := strconv.ParseUint(temp, 10, 64)
	sum = (sum + newValue)
	fmt.Printf("Total Sum = %d ", (sum))
	wg.Wait()
}

func (wrkr *worker) compute(wg *sync.WaitGroup) {
	var msgFrag map[string]interface{}
	value := uint64(0)
	error := ""
	prefix := ""
	suffix := ""
	var count_par int

	//	get the chunk
	err := json.Unmarshal(<-wrkr.job, &msgFrag)
	check(err)

	//get the individual attributes from the fragment
	fname := msgFrag["Fname"].(string)
	start := msgFrag["Start"].(float64)
	end := msgFrag["End"].(float64)

	//create chunk
	read, err := os.Open(fname)
	check(err)

	defer read.Close()

	_, e_ := read.Seek(int64(start), 0)
	check(e_)

	chunk := make([]byte, byte(end-start+1))

	_, e2 := read.Read(chunk)
	check(e2)

	chunk_ := string(chunk)
	fmt.Printf("[%s]\n", chunk_)
	chunkStr := cleanChunk(chunk_)

	if strings.Contains(chunkStr, " ") {
		psum := uint64(0)
		chunks := strings.Split(string(chunkStr), " ")
		// if no starting or ending space send it as preffix or suffix
		prefix = chunks[0]
		suffix = chunks[len(chunks)-1]

		// else calculate sum
		for _, element := range chunks[1 : len(chunks)-1] {

			val, _ := strconv.ParseInt(element, 10, 64)
			count_par++
			psum += uint64(val)
		}

		psum = (psum / uint64(count_par))

		// get average sum
		value = psum

	} else {
		//	-- If no average sum return in exception of json file
		error = chunkStr
	}

	res := &response{
		value,
		prefix,
		suffix,
		error,
	}

	pckRes, _ := json.Marshal(res)

	wg.Done()

	wrkr.result <- pckRes
}

func main() {

	//Get the inputs
	M := os.Args[1]
	fname := os.Args[2]

	//Convert string to float
	frag, _ := strconv.ParseFloat(M, 64)

	//Pass the values to cordinator
	coordinator(fname, frag)
}
