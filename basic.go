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

import "strings"

type GoogleGroup struct {
	GroupName string
	Cookie string
}

type GoogleGroupConversation struct {
	Author string
	Title string
	Id string
	GroupName string
	AuthorNameToEmail map[string]string
}

type GoogleGroupMessage struct {
	Author string
	Content string
	Time string
}

func NewGoogleGroup(name string, cookie ...string) GoogleGroup {
	ret := GoogleGroup{GroupName: strings.Split(name, "@")[0]}
	if len(cookie) != 0 {
		ret.Cookie = cookie[0]
	}
	return ret
}
