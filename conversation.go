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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (c GoogleGroupConversation) GetAllMessages(client http.Client, removeGmailQuote bool) []GoogleGroupMessage {
	targetUrl := fmt.Sprintf("https://groups.google.com/g/%s/c/%s", c.GroupName, c.Id)
	var ret []GoogleGroupMessage

	req, _ := http.NewRequest("GET", targetUrl, nil)
	req.Header.Set("cookie", c.Cookie)
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("Google Groups Crawler: http GET request status code: %d\n", res.StatusCode)
		return ret
	}
	resp, _ := ioutil.ReadAll(res.Body)
	body := string(resp)

	start := strings.LastIndex(body, "AF_initDataCallback({key: 'ds")
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

	var dataArray []interface{}
	err := json.Unmarshal([]byte(body), &dataArray)
	if err != nil {
		return ret
	}
	if len(dataArray) < 3 {
		return ret
	}
	msgArray, ok := dataArray[2].([]interface{})
	if !ok {
		return ret
	}

	for _, msg := range msgArray {
		singleMsgArray, ok := msg.([]interface{})
		if !ok || len(singleMsgArray) < 1 {
			continue
		}
		singleMsgArray, ok = singleMsgArray[0].([]interface{})
		if !ok || len(singleMsgArray) < 2 {
			continue
		}
		singleMsgArray0, ok := singleMsgArray[0].([]interface{})
		if !ok || len(singleMsgArray0) < 9 {
			continue
		}
		authorEmailArray, ok := singleMsgArray0[2].([]interface{})
		if !ok || len(authorEmailArray) < 1 {
			continue
		}
		authorEmailArray, ok = authorEmailArray[0].([]interface{})
		if !ok || len(authorEmailArray) < 3 {
			continue
		}
		author, ok := authorEmailArray[0].(string)
		if !ok {
			continue
		}
		email, ok := authorEmailArray[2].(string)
		if !ok {
			continue
		}
		singleMsgArray0, ok = singleMsgArray0[8].([]interface{})
		if !ok || len(singleMsgArray0) < 1 {
			continue
		}
		time, ok := singleMsgArray0[0].(float64)
		if !ok {
			continue
		}
		singleMsgArray1, ok := singleMsgArray[1].([]interface{})
		if !ok || len(singleMsgArray1) < 2 {
			continue
		}
		singleMsgArray1, ok = singleMsgArray1[1].([]interface{})
		if !ok || len(singleMsgArray1) < 1 {
			continue
		}
		singleMsgArray1, ok = singleMsgArray1[0].([]interface{})
		if !ok || len(singleMsgArray1) < 2 {
			continue
		}
		singleMsgArray1, ok = singleMsgArray1[1].([]interface{})
		if !ok || len(singleMsgArray1) < 2 {
			continue
		}
		content, ok := singleMsgArray1[1].(string)
		if !ok {
			continue
		}

		if removeGmailQuote {
			content = strings.Split(content, "<div class=\"gmail_quote\">")[0]
		}

		ret = append(ret, GoogleGroupMessage{
			Author: author,
			AuthorEmail: email,
			Content: content,
			Time: time,
		})
	}
	return ret
}

