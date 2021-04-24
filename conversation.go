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

	"github.com/PuerkitoBio/goquery"
)

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
