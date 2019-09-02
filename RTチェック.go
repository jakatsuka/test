package main

import (
  "fmt"
  "os"
  "github.com/dghubble/go-twitter/twitter"
  "github.com/dghubble/oauth1"
)

const (
    targetTweetStatus = `1167950194933059587`
    consumerKey = `NSuQ7TQTp2pa0ic8TpKYFWTGT`
    consumerSecret = `w52h9P4LCJR9lSa1dix5BvOj4BAMTZ9P4g2UQMtQNoH3dJsODo`
    accessToken = `232812048-naspXMpuraS5fJl5rzMSNKYXlh890kQPXXnJDphB`
    accessTokenSecret = `7l2tILeMFXZ8WwV4Ilo2ItlupTQozGdCgx1hWusfDztiD`
    until = "2019-08-01"
    saveFilePath = "./retweeted_users.csv"
)

func main() {
    config := oauth1.NewConfig(consumerKey, consumerSecret)
    token := oauth1.NewToken(accessToken, accessTokenSecret)
    httpClient := config.Client(oauth1.NoContext, token)
    twitterClient := twitter.NewClient(httpClient)

    savefile, err := os.Create(saveFilePath)
    if err != nil { panic(err) }
    defer savefile.Close()

    savefile.Write(([]byte)("Link, UserName, RetweetedAt, RetweetID\n"))

    targetTweet, _, _ := twitterClient.Statuses.Show(targetTweetStatus, nil)

    // ツイートに画像が含まれている場合は、本文の最後に短縮URLが付与されるので
    // 正規表現なり、文字をスライスして抜き出す。
    tweetText := string([]rune(targetTweet.Text)[:104])

    var maxID int64 = 1111111111111111111

    for ;; {
        search, _, _ := twitterClient.Search.Tweets(&twitter.SearchTweetParams{
            Count: 100,
            MaxID: maxID,
            Query: tweetText,
            SinceID: targetTweetStatus,
            Until: until,
        })

        if len(search.Statuses) == 0 {
            break
        }

        tmpLines := ""
        for _, status := range search.Statuses {
            link := "https://twitter.com/" + status.User.ScreenName
            tmpLines += fmt.Sprintf("%s, %s, %s, %d\n",
                link,
                status.User.Name,
                status.CreatedAt,
                status.ID,
            )
            maxID = status.ID
        }
        savefile.Write(([]byte)(tmpLines))

    }

}
