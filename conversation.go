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
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var reg *regexp.Regexp

func (c GoogleGroupConversation) GetAllMessages(client http.Client) []GoogleGroupMessage {
	var ret []GoogleGroupMessage
	res, err := client.Get("https://groups.google.com/g/" + c.GroupName + "/c/" + c.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".wqmMgb").Remove()
	doc.Find(".gmail_quote").Remove()
	doc.Find(".gmail_attr").Remove()

	doc.Find(".BkrUxb").Each(func(i int, s *goquery.Selection) {
		author := s.Find(".s1f8Zd").Text()
		content, _ := s.Find(".ptW7te").Html()
		time := s.Find(".zX2W9c").Text()
		ret = append(ret, GoogleGroupMessage{
			Author: author,
			Content: content,
			Time: time,
		})
	})
	return ret
}

func (c *GoogleGroupConversation) GetAuthorNameToEmailMapping(client http.Client, cookies ...string) {
	var cookie string
	if len(cookies) != 0 {
		cookie = cookies[0]
	}
	targetUrl := "https://groups.google.com/g/" + c.GroupName + "/c/" + c.Id

	req, _ := http.NewRequest("GET", targetUrl, nil)
	req.Header.Set("cookie", cookie)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	htmlStr := doc.Text()

	start := strings.LastIndex(htmlStr, "AF_initDataCallback({")
	end := strings.LastIndex(htmlStr, ", sideChannel: {}}")
	if start >= end {
		return
	}
	data := htmlStr[start:end]
	startIdx := strings.Index(data, "[")
	if startIdx == -1 {
		return
	}
	data = data[startIdx:]

	if reg == nil {
		reg, _ = regexp.Compile("\\[\\[\"([A-Za-z0-9]|[ ])*\",((null)|([\\S]*)),\"[\\S\"]*@[\\S\"]*.[\\S\"]*\"")
	}
	buf := reg.FindAllString(data, -1)
	for i, b := range buf {
		buf[i] = b[strings.LastIndex(b, "[[") + 2:]
	}
	NameToEmail := make(map[string]string)
	for _, b := range buf {
		s := strings.Split(b, ",")
		if len(s) < 3 {
			continue
		}
		NameToEmail[s[0][1:len(s[0])-1]] = s[2][1:len(s[2])-1]
	}
	c.AuthorNameToEmail = NameToEmail
}
