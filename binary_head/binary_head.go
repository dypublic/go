package main

import (
	"bufio"
	"fmt"
	"os"
	//	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	//    "reflect"
	"errors"
)

const ( // iota is reset to 0
	SB_start    = iota // c0 == 0
	SB_module   = iota // c1 == 1
	SB_resource = iota // c2 == 2
	SB_index    = iota // c2 == 2
)

type meta struct {
	len  int
	path string
}
type fileMeta map[int]int
type resourceTable map[int]meta

func metaTableCompare(table resourceTable, outer fileMeta) bool {
	if table == nil || outer == nil {
		return false
	}
	if len(table) != len(outer) {
		return false
	}
	for k, v := range table {
		if v_out, ok := outer[k]; !ok {
			return false
		} else {
			if v.len != v_out {
				return false
			}
		}
	}
	return true
}

var RegularResHead *regexp.Regexp = nil
var RegularResLen *regexp.Regexp = nil

func processBlock(scanner *bufio.Scanner, w *bufio.Writer) int {
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 3 {
			if line == "};" {
				//fmt.Println("End of block with }")
				break
			}
			fmt.Println("")
		}
		if line[len(line)-1] == ',' {
			line = line[:len(line)-1]
		}
		words := strings.Split(line, ",")
		for _, word := range words {
			byte_int, err := strconv.ParseUint(word, 0, 8)
			if err != nil {
				log.Panic("Parse bytes error,", err)
			}
			count++
			w.WriteByte(byte(byte_int))
		}
	}
	w.Flush()
	return count
}

func processModule(scanner *bufio.Scanner, wf *os.File) int {
	w := bufio.NewWriter(wf)

	scanner.Scan()
	line := scanner.Text()
	if !strings.HasPrefix(line, "const char sb_model[]") {
		fmt.Println("Error, no SB_model found!")
		return 0
	}

	return processBlock(scanner, w)
}

func compileRegexpHead() *regexp.Regexp {
	if RegularResHead == nil {
		//const unsigned char res_3157[] = {
		RegularResHead = regexp.MustCompile("\\sres_(\\d+)\\[\\]\\s")
	}
	return RegularResHead

	//fmt.Printf("%q\n", re.FindStringSubmatch("-axbyc-"))
	//fmt.Printf("%q\n", re.FindStringSubmatch("-abzc-"))
}

func compileRegexpLen() *regexp.Regexp {
	if RegularResLen == nil {
		//int res_3157_nbytes = 371;
		RegularResLen = regexp.MustCompile("\\bres_(\\d+)_nbytes\\b\\s=\\s\\b(\\d+);")
	}
	return RegularResLen
}

func processResource(scanner *bufio.Scanner, wf *os.File, indexSize fileMeta) int {
	w := bufio.NewWriter(wf)
	scanner.Scan()
	line := scanner.Text()
	if !strings.HasPrefix(line, "const unsigned char res_") {
		fmt.Println("Error, no res block found!,", line)
	}
	re := compileRegexpHead()
	head_words := re.FindStringSubmatch(line)
	//fmt.Println(head_words)
	if head_words == nil {
		fmt.Println("head not found, the line is, ", line)
		return 0
	}
	resIndex, err := strconv.Atoi(head_words[1])
	if err != nil {
		fmt.Println("conver number(", head_words[1], ") to int fail, ", err)
		return 0
	}
	//print(line)
	size := processBlock(scanner, w)

	scanner.Scan()
	endLine := scanner.Text()
	regLen := compileRegexpLen()
	endWords := regLen.FindStringSubmatch(endLine)
	//fmt.Println(end_words)
	nums, err := convertResult(endWords)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	resIndex, resLen := nums[0], nums[1]
	if resIndex != resIndex || size != resLen {
		fmt.Printf("res_index or len error, i_start %d, i_end %d, len %d, len_end %d\n",
            resIndex, resIndex, size, resLen)
		return 0
	}
	fmt.Printf("res:%d, len:%d\n", resIndex, size)
	indexSize[resIndex] = size
	return size
}

func convertResult(words []string) ([]int, error) {
	if words == nil || len(words) != 3 {
		return nil, errors.New("end not found")
	}
	if words[1] == "" || words[2] == "" {
		return nil, errors.New("end number not found")
	}

	index, err := strconv.Atoi(words[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintln(
			"conver number(", words[1], ") to int fail, ", err))
	}

	resLen, err := strconv.Atoi(words[2])
	if err != nil {
		return nil, errors.New(fmt.Sprintln(
			"conver number(", words[2], ") to int fail, ", err))
	}
	nums := []int{index, resLen}
	return nums, nil
}

func convertTableResult(words []string) (string, []int, error) {
	if words == nil || len(words) != 4 {
		return "", nil, errors.New("end not found")
	}
	if words[1] == "" || words[2] == "" || words[3] == "" {
		return "", nil, errors.New("end number not found")
	}
	path := words[1]
	index, err := strconv.Atoi(words[2])
	if err != nil {
		return "", nil, errors.New(fmt.Sprintln(
			"conver number(", words[2], ") to int fail, ", err))
	}

	resLen, err := strconv.Atoi(words[3])
	if err != nil {
		return "", nil, errors.New(fmt.Sprintln(
			"conver number(", words[3], ") to int fail, ", err))
	}
	nums := []int{index, resLen}
	return path, nums, nil
}

func checkingIndex(scanner *bufio.Scanner, indexSizeMap fileMeta) (resourceTable, error) {
	fmt.Println(len(indexSizeMap))
	foundHeadCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "const struct _sb_image_resources sb_model_images") {
			//fmt.Println("Error, no res block found!,", line)
			break
		}
		foundHeadCount++
	}
	checkingMap := make(resourceTable)
	for scanner.Scan() {
		line := scanner.Text()
		//{ "/images/list_btn_background5.png", res_1667, 317 },
		rex_item := regexp.MustCompile(`\{\s"(.+)",\sres_(\d+),\s(\d+)\s\},`)
		words := rex_item.FindStringSubmatch(line)
		path, nums, errConver := convertTableResult(words)
		if errConver != nil {
			if strings.Contains(line, "NULL, NULL, 0") {
				//fmt.Println("found NULL end")
				break
			}
			fmt.Println(errConver, "at ", line)
		}
		index, resLen := nums[0], nums[1]
		fmt.Printf("res:%d, len:%d\n", index, resLen)
		checkingMap[index] = meta{resLen, path}
	}
	//
	if len(checkingMap) != len(indexSizeMap) {
		fmt.Println("len of map is not equal")
	}
	equal := metaTableCompare(checkingMap, indexSizeMap)
	if !equal {
		return nil, errors.New("it is not equal for check the index and size")
	}
	fmt.Println("check success! resource number is: ", len(indexSizeMap))
	return checkingMap, nil
}

func pointerPart(f *os.File, res resourceTable, startPos int) error {
	SBModuleLine := "const char* sb_model = BIN_ADDR;\r\n"
	if _, err := f.WriteString(SBModuleLine); err != nil {
		fmt.Println("fail to write the SB module line")
		return err
	}
	pointerTemplate := "const unsigned char* res_%d = BIN_ADDR + %d;\r\n"
	pos := startPos
	for i, length := 0, len(res); i < length; i++ {
		line := fmt.Sprintf(pointerTemplate, i, pos)
		if _, err := f.WriteString(line); err != nil {
			fmt.Println("fail to write the pointers line:", i)
			return err
		}
		pos += res[i].len
	}
	lenLine := fmt.Sprintf("\r\nconst int total_len = %d;\r\n\r\n", pos)
	if _, err := f.WriteString(lenLine); err != nil {
		fmt.Println("fail to write the length line")
		return err
	}
	return nil
}

func headFileGen(res resourceTable, sbModuleSize int) error {
	firstLine := `const struct _sb_image_resources sb_model_images[] = {` + "\r\n"
	//{ "/images/list_btn_background5.png", res_1667, 317 },
	contentTemplate := `{ "%s", res_%d, %d },` + "\r\n"
	endLine := "{ NULL, NULL, 0 },\r\n};\r\n"
	f, err := os.Create("ArmingStation_s.h")
	if err != nil {
		fmt.Println("create head file fail")
		return err
	}
	defer f.Close()

	if err := pointerPart(f, res, sbModuleSize); err != nil {
		return err
	}

	_, err_f := f.WriteString(firstLine)
	if err_f != nil {
		fmt.Println("write end line fail", err_f)
		return err_f
	}
	for i, metaValue := range res {
		item := fmt.Sprintf(contentTemplate, metaValue.path, i, metaValue.len)
		_, err := f.WriteString(item)
		if err != nil {
			fmt.Println("write content line fail")
			return err
		}
	}
	_, err_e := f.WriteString(endLine)
	if err_e != nil {
		fmt.Println("write end line fail", err_e)
		return err
	}
	println("write all content to file")
	return nil
}

func main() {
	headFile, err := os.Open("ArmingStation.h")
	if err != nil {
		fmt.Println("Open head file fail!")
		os.Exit(1)
	}
	defer headFile.Close()
	//buf_reader := bufio.NewReader(head_file)
	binFile, err := os.Create("binary.bin")
	if err != nil {
		fmt.Println("Open output file fail!")
	}
	defer binFile.Close()
	scanner := bufio.NewScanner(headFile)
	sbModuleLength := processModule(scanner, binFile)
	ongoing := true
	indexSizeMap := make(map[int]int)
	for ongoing {
		if processResource(scanner, binFile, indexSizeMap) <= 0 {
			ongoing = false
		}
	}
	resTable, err := checkingIndex(scanner, indexSizeMap)
	if err != nil {
		fmt.Println(err)
	}
	headFileGen(resTable, sbModuleLength)
}
