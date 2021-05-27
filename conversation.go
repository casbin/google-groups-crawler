// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google_groups_crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SeparateArray(input string) []string {
	var ret []string
	needSkip := false
	level := 0
	var start int
	for i, c := range input {
		if needSkip {
			needSkip = false
			continue
		}
		if c == '\\' {
			needSkip = true
			continue
		}
		if c == '[' {
			if level == 0 {
				start = i
			}
			level++
		}
		if c == ']' {
			level--
			if level == 0 {
				ret = append(ret, input[start:i+1])
			}
		}
	}
	return ret
}

func (c GoogleGroupConversation) GetAllMessages(client http.Client, removeGmailQuote bool, cookies ...string) []GoogleGroupMessage {
	cookie := ""
	if len(cookies) > 0 {
		cookie = cookies[0]
	}
	targetUrl := fmt.Sprintf("https://groups.google.com/g/%s/c/%s", c.GroupName, c.Id)

	var ret []GoogleGroupMessage

	req, _ := http.NewRequest("GET", targetUrl, nil)
	req.Header.Set("cookie", cookie)
	res, _ := client.Do(req)
	if res.StatusCode != 200 {
		fmt.Printf("Google Groups Crawler: http GET request status code: %d\n", res.StatusCode)
		return ret
	}
	resp, _ := ioutil.ReadAll(res.Body)
	body := string(resp)

	start := strings.LastIndex(body, "AF_initDataCallback({key: 'ds:")
	end := strings.LastIndex(body, ", sideChannel: {}});")
	if start >= end {
		return ret
	}
	body = body[start:end]
	start = strings.Index(body, "data:")
	if start < 0 {
		return ret
	}
	body = body[start+5:]
	start = strings.Index(body, "[[[[")
	if start < 0 {
		return ret
	}
	body = body[start+1:len(body)-4]

	conversations := SeparateArray(body)
	for _, c := range conversations {
		c = c[1:len(c)-1]
		sep := SeparateArray(c)
		if len(sep) != 2 {
			continue
		}
		subStr := sep[0]
		subStr = subStr[1:len(subStr)-1]
		sep = SeparateArray(subStr)
		if len(sep) != 3 {
			continue
		}
		emailStart := strings.Index(sep[0], "[[")
		emailEnd := strings.Index(sep[0], "]")
		if emailStart >= emailEnd {
			continue
		}
		emailStr := sep[0][emailStart+2:emailEnd]
		emailSep := strings.Split(emailStr, ",")
		if len(emailSep) < 3 {
			continue
		}
		author := emailSep[0]
		email := emailSep[2]
		if len(sep[1]) < 3 {
			continue
		}
		msgSep := SeparateArray(sep[1][1:len(sep[1])-1])
		if len(msgSep) < 1 || len(msgSep[0]) < 3 {
			continue
		}
		msgSep = SeparateArray(msgSep[0][1:len(msgSep[0])-1])
		if len(msgSep) < 1 || len(msgSep[0]) < 3 {
			continue
		}
		msgSep = SeparateArray(msgSep[0][1:len(msgSep[0])-1])
		if len(msgSep) < 1 {
			continue
		}
		msg := msgSep[0]
		if len(msg) < 3 {
			continue
		}
		msg = msg[1:len(msg)-1]
		msgIndex := strings.Index(msg, ",")
		if msgIndex < 0 {
			continue
		}
		msg = msg[msgIndex+1:]
		sep1 := SeparateArray(sep[0][1:len(sep[0])-1])
		if len(sep1) == 0 {
			continue
		}
		unixTime := sep1[len(sep1)-1]
		if len(unixTime) < 3 {
			continue
		}
		unixTime = unixTime[1:len(unixTime)-1]

		if len(author) < 3 || len(email) < 3 || len(msg) < 3 {
			continue
		}
		author = author[1:len(author)-1]
		email = email[1:len(email)-1]
		msg = msg[1:len(msg)-1]

		if removeGmailQuote {
			msg = strings.Split(msg, "\\u003cdiv class\\u003d\\\"gmail_quote\\\"\\u003e")[0]
		}

		ret = append(ret, GoogleGroupMessage{Author: author, AuthorEmail: email, Content: msg, Time: unixTime})
	}
	return ret
}
