package colly

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func main() {

}

func CollyInit(username string) string {
	var lock sync.Mutex
	var htmlStr string
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) { // 當Visit訪問網頁後，網頁響應(Response)時候執行的事情
		// fmt.Println(string(r.Body)) // 返回的Response物件r.Body 是[]Byte格式，要再轉成字串
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println(err)
	})

	// 抓類別Class 名稱
	c.OnHTML(".f_4No > section > a > div > img", func(e *colly.HTMLElement) {
		imageHref := e.Attr("srcset")
		url := strings.Split(imageHref, " ")[0]
		resp, err := http.Get(url)
		if err != nil {
			return
		}

		var ext string
		switch resp.Header["Content-Type"][0] {
		case "image/jpeg":
			ext = "jpg"
		// case "image/svg+xml":
		// 	ext = "svg"
		default:
			ext = ""
		}
		if ext != "" {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			base64Encoding := getDataImageBase64(body)
			lock.Lock()
			html := fmt.Sprintf("<img src=\"%s\" />", base64Encoding)
			htmlStr += html
			lock.Unlock()
			// name := fmt.Sprintf("%d.%s", e.Index, ext)
			// os.WriteFile(name, body, 0666) //存下圖檔
		}
	})

	url := "https://www.deviantart.com/" + username
	c.Visit(url)
	return htmlStr
}

func getDataImageBase64(source []byte) string {
	var base64Encoding string

	mimeType := http.DetectContentType(source)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(source)
	return base64Encoding
}
