package main

import (
	"bytes"
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
    "github.com/dyatlov/go-readability"
	"golang.org/x/net/html"
	"html/template"
	"log"
    "path/filepath"
    "io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func get_oembed(embed_url string) interface{} {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get("https://noembed.com/embed?maxwidth=854&maxheight=640&url=" + url.QueryEscape(embed_url))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
    target := struct {
        Error string
        Html template.HTML
    }{}
	json.NewDecoder(r.Body).Decode(target)
    return target
}

func get_readability(embed_url string) string {
    client := &http.Client{Timeout: 10 * time.Second}
    r, err := client.Get(embed_url)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Body.Close()
    html_body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    doc, err := readability.NewDocument(string(html_body))
    if err != nil {
        log.Fatal(err)
    }
    doc.WhitelistTags = []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6", "h7", "h8", "a", "img", "audio", "video", "source"}
    doc.WhitelistAttrs = map[string][]string{
        "img": {"src"},
        "a": {"href"},
        "audio": {"src"},
        "video": {"src"},
    }
    doc.CleanConditionally = false
    html_node, err := html.Parse(strings.NewReader(doc.Content()))
    if err != nil {
        log.Fatal(err)
    }

    base_url, err := url.Parse(embed_url)
    if err != nil {
        log.Fatal(err)
    }
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode {
            for i, a := range n.Attr {
                if a.Key == "href" || a.Key == "src" {
                    new_url, err := url.Parse(a.Val)
                    if err != nil {
                        log.Fatal(err)
                    }
                    n.Attr[i].Val = base_url.ResolveReference(new_url).String()
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(html_node)
    var b bytes.Buffer
    html.Render(&b, html_node)
    return b.String()
}

func get_file_embed(embed_url string) string {
    switch filepath.Ext(strings.ToLower(embed_url)) {
    case ".jpg", ".jpeg", ".png", ".gif", ".svg":
        return "<img src=\"" + string(template.HTMLAttr(embed_url)) + "\">"
    case ".opus", ".weba", ".wav", ".mp3", ".flac":
        return "<audio controls src=\"" + string(template.HTMLAttr(embed_url)) + "\">"
    case ".ogg", ".webm", ".mp4":
        return "<video controls src=\"" + string(template.HTMLAttr(embed_url)) + "\">"
    default:
        return ""
    }
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
    json := get_oembed(string(link)).(struct {
        Error string
        Html template.HTML
    })
    embed := get_file_embed(string(link))
    r_embed := get_readability(string(link))
    log.Println(r_embed)
	if string(content) == "Embed" {
        if embed != "" {
            out.WriteString(embed)
        } else if json.Error == "" && json.Html != "" {
		    out.WriteString(string(json.Html))
        } else {
            out.WriteString(r_embed)
        }
	} else if string(content) == "Image" {
        out.WriteString(embed);
	}
}
