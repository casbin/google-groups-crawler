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
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (g GoogleGroup) GetConversations(client http.Client) []GoogleGroupConversation {
	var ret []GoogleGroupConversation
	res, err := client.Get("https://groups.google.com/g/" + g.GroupName)
	if err != nil {
		return ret
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s\n", res.StatusCode, res.Status)
		return ret
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ret
	}

	doc.Find(".yhgbKd").Each(func(i int, s *goquery.Selection) {
		author := s.Find(".z0zUgf").Text()
		title := s.Find(".iBQX0d").Find(".o1DPKc").Text()
		time := s.Find(".kOkyJc").Find(".tRlaM").Text()
		href, _ := s.Find(".Dysyo").Attr("href")
		hrefs := strings.Split(href, "/")
		id := hrefs[len(hrefs) - 1]
		newConversation := GoogleGroupConversation{
			Author: author,
			Title: title,
			Id: id,
			GroupName: g.GroupName,
			Time: time,
		}
		ret = append(ret, newConversation)
	})
	return ret
}
