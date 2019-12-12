package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	// "log"
)

//const filePath string = "./data/users.txt"

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	bufUsers := bytes.Buffer{}

	lines := strings.Split(string(fileContents), "\n")

	users := make([]map[string]interface{}, 0, len(lines))
	for _, line := range lines {
		user := make(map[string]interface{})
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	pattern1 := regexp.MustCompile("Android")
	pattern2 := regexp.MustCompile("MSIE")

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

			//if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
			if ok := pattern1.MatchString(browser); ok && err == nil {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
			if ok := pattern2.MatchString(browser); ok && err == nil {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		//foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
		//bufUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email))
		bufUsers.WriteString("[")
		bufUsers.WriteString(strconv.Itoa(i))
		bufUsers.WriteString("] ")
		bufUsers.WriteString(string(user["name"].(string)))
		bufUsers.WriteString(" <")
		bufUsers.WriteString(email)
		bufUsers.WriteString(">\n")

	}

	fmt.Fprintln(out, "found users:\n"+bufUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
