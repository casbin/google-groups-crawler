# Google Groups Crawler

## Installation

``` go
import crawler "github.com/Kininaru/google-groups-crawler"
```

## Usage

We must get an instance of `GoogleGroup` first:

``` go
group := crawler.NewGoogleGroup("group name")
```

### Get all conversations of the group

- For some special reasons, you cannot access Google Groups in some area. You can set up a http proxy, and fill the parameter `http.Client` with it. If you can access Google Groups directly, then you can just fill the parameter like the example code.
- this function returns an array of `GoogleGroupConversation`

``` go
conversations := group.GetConversations(http.Client{})
```

### Get all messages of the conversation

- `conversation` is an instance of `GoogleGroupConversation`
- parameter `http.Client` is the same effect as above
- this function returns an array of `GoogleGroupMessage`

```go
messages := conversation.GetAllMessages(http.Client{})
```



## Data Structure

```
type GoogleGroup struct {
   GroupName string
}

type GoogleGroupConversation struct {
   Author string
   Title string
   Id string
   GroupName string
}

type GoogleGroupMessage struct {
   Author string
   Content string
   Time string
}
```