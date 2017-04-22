package main

import (
    "net/url"
    "encoding/json"
	"bytes"
    "strings"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"html/template"
	"log"
	"net/http"
	"time"
)

func get_oembed(embed_url string, target interface{}) {
    client := &http.Client{Timeout: 10 * time.Second}
    r, err := client.Get("https://noembed.com/embed?maxwidth=854&maxheight=640&url=" + url.QueryEscape(embed_url))
    if err != nil {
        log.Fatal(err)
    }
    defer r.Body.Close()
    json.NewDecoder(r.Body).Decode(target)
}
func render_markdown(markdown_raw string) string {
    renderer := &renderer{Html: blackfriday.HtmlRenderer(0, "", "").(*blackfriday.Html)}
    markdown_san := string(html.EscapeString(strings.Replace(markdown_raw, "\r", "\n", -1)))
    html_raw := string(blackfriday.Markdown([]byte(markdown_san), renderer, 0))
    policy := bluemonday.UGCPolicy()
    policy.AllowElements("iframe").AllowAttrs("src", "width", "height", "frameBorder").OnElements("iframe")
    html_san := policy.Sanitize(html_raw)
    html_san = html_raw
    return html_san
}
type renderer struct {
    *blackfriday.Html
}
func (r *renderer) Image(out *bytes.Buffer, link []byte, title []byte, content []byte) {
    type oembed_return struct {
        Html template.HTML
    }
    json := &oembed_return{}
    get_oembed(string(link), json)
    if (string(content) == "Embed") {
        out.WriteString(string(json.Html))
    } else if string(content) == "Image" {
        out.WriteString("<img src=\"")
        out.WriteString(string(template.HTMLAttr(link)))
        out.WriteString("\">")
    }
}
