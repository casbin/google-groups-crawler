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

func (g GoogleGroup) GetAllConversations(client http.Client) []GoogleGroupConversation {
	targetUrl := fmt.Sprintf("https://groups.google.com/g/%s", g.GroupName)
	var ret []GoogleGroupConversation

	req, _ := http.NewRequest("GET", targetUrl, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
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

	var dataArray []interface{}
	err := json.Unmarshal([]byte(body), &dataArray)
	if err != nil {
		return ret
	}
	if len(dataArray) < 3 {
		return ret
	}
	conversationArray, ok := dataArray[2].([]interface{})
	if !ok {
		return ret
	}
	for _, c := range conversationArray {
		cArray, ok := c.([]interface{})
		if !ok {
			continue
		}
		if len(cArray) < 1 {
			continue
		}
		cArray, ok = cArray[0].([]interface{})
		if !ok {
			continue
		}
		if len(cArray) < 6 {
			continue
		}
		id, ok := cArray[1].(string)
		if !ok {
			continue
		}
		title, ok := cArray[2].(string)
		if !ok {
			continue
		}
		cArray, ok = cArray[5].([]interface{})
		if !ok {
			continue
		}
		if len(cArray) < 1 {
			continue
		}
		time, ok := cArray[0].(float64)
		if !ok {
			continue
		}
		ret = append(ret, GoogleGroupConversation{
			GroupName: g.GroupName,
			Id: id,
			Title: title,
			Time: time,
			Cookie: g.Cookie,
		})
	}
	return ret
}
