package main

import (
	"bytes"
	"encoding/json"
	"github.com/dyatlov/go-readability"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func get_oembed(embed_url string) string {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get("https://noembed.com/embed?maxwidth=854&maxheight=640&url=" + url.QueryEscape(embed_url))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	target := struct {
		Error string
		Html  string
	}{}
	json.NewDecoder(r.Body).Decode(&target)
	if target.Error != "" {
		return ""
	}
	return target.Html
}

func get_rembed(embed_url string) string {
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
		"img":   {"src"},
		"a":     {"href"},
		"audio": {"src"},
		"video": {"src"},
	}
	doc.CleanConditionally = false
	html_node, err := html.Parse(strings.NewReader(doc.Content()))
	if err != nil {
		log.Fatal(err)
	}

	//translate the relative urls in the document to absolute urls
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

func get_fembed(embed_url string) string {
	switch filepath.Ext(strings.ToLower(embed_url)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".svg":
		return "<img src=\"" + string(template.HTMLAttr(embed_url)) + "\">"
	case ".opus", ".weba", ".wav", ".mp3", ".flac":
		return "<audio controls src=\"" + string(template.HTMLAttr(embed_url)) + "\">"
	case ".ogg", ".ogv", ".webm", ".mp4":
		return "<video controls src=\"" + string(template.HTMLAttr(embed_url)) + "\">"
	default:
		return ""
	}
}

func render_text(text_raw string, process_title bool) (template.HTML, string) {
	var html_raw bytes.Buffer
	text_san := string(html.EscapeString(strings.Replace(text_raw, "\r\n", "\n", -1)))
	re_str := ""
	if process_title {
		re_str = re_str + `(?:\A\s*(?P<title>(?U).*)[\t\f\r ]*\n)|`
	}
	re_str = re_str +
		`(?sU:` + "```" + `(?P<pre>.*)` + "```" + `\n??)|` +
		`(?P<url>https?://[^\s]+)|` +
		`(?P<hashtag>#[_\pLetter\pNumber]+)|` +
		`(?P<break>\n)|` +
		`(?P<normal>.)`
	re := regexp.MustCompile(re_str)
	matches := re.FindAllStringSubmatch(text_san, -1)
	title := ""
	html_raw.WriteString("<p>")
	for _, match := range matches {
		for group_idx, group := range match {
			if group == "" {
				continue
			}
			name := re.SubexpNames()[group_idx]
			switch name {
			case "title":
				html_raw.WriteString("<h2>" + group + "</h2>")
				title = group
			case "pre":
				html_raw.WriteString("<pre>" + group + "</pre>")
			case "url":
				if fembed := get_fembed(string(group)); fembed != "" {
					html_raw.WriteString(fembed)
				} else if oembed := get_oembed(string(group)); oembed != "" {
					html_raw.WriteString(oembed)
				} else if rembed := get_rembed(string(group)); rembed != "" {
					html_raw.WriteString(rembed)
				} else {
					html_raw.WriteString("<a href=\"" + group + "\">" + group + "</a>")
				}
			case "hashtag":
				html_raw.WriteString("<a href=\"/search/?query=" + group + "\">" + group + "</a>")
			case "break":
				html_raw.WriteString("</p><p>")
			case "normal":
				html_raw.WriteString(group)
			}
		}
	}
	html_raw.WriteString("</p>")
	policy := bluemonday.UGCPolicy()
	policy.AllowElements("iframe").AllowAttrs("src", "width", "height", "frameBorder").OnElements("iframe")
	policy.AllowElements("video").AllowAttrs("src", "type").OnElements("video")
	policy.AllowElements("audio").AllowAttrs("src", "type").OnElements("audio")
	policy.AllowElements("image").AllowAttrs("src", "type").OnElements("image")
	policy.AllowElements("source").AllowAttrs("src", "type").OnElements("source")
	html_san := template.HTML(policy.Sanitize(html_raw.String()))
	return html_san, title
}
