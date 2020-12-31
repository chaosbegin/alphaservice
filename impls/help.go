package impls

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	goMdHtml "github.com/gomarkdown/markdown/html"
	goMdParser "github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

var HelpSrv Help

type Help struct {
	DocMap       sync.Map
	extensions   goMdParser.Extensions
	rendererOpts goMdHtml.RendererOptions
	policy       *bluemonday.Policy
	imgDir       string
	toc          *HelpToc
	md           goldmark.Markdown
	id           int32
}

type HelpToc struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Key      string     `json:"key"`
	Children []*HelpToc `json:"children,omitempty"`
}

func (this *Help) getId() int {
	if atomic.LoadInt32(&this.id) > int32(1<<30) {
		atomic.StoreInt32(&this.id, 1)
		return int(atomic.LoadInt32(&this.id))
	} else {
		return int(atomic.AddInt32(&this.id, 1))
	}
}

func (this *Help) walk(path string, info os.FileInfo, node *HelpToc) {
	var err error
	dirName := this.FixName(filepath.Base(path))
	files := this.listFiles(path)
	for _, filename := range files {
		fpath := filepath.Join(path, filename)
		fio, _ := os.Lstat(fpath)

		if fio.IsDir() {
			if this.IsHelpName(filename) {
				//logs.Trace("filename:", filename, " fixname:", this.FixName(filename))
				name := ""
				key := ""
				tpath := filepath.Join(fpath, this.FixName(filename)+".md")
				_, err = os.Lstat(tpath)
				if err == nil {
					name, key, err = this.readOneMd(tpath)
					if err != nil {
						logs.Error(err.Error())
					}
				}

				if len(name) < 1 {
					name = this.FixName(filename)
				}

				child := HelpToc{this.getId(), name, key, []*HelpToc{}}
				node.Children = append(node.Children, &child)
				this.walk(fpath, fio, &child)
			}

		} else {
			if (dirName + ".md") == filename {
				//logs.Trace("same dir name, ignore...")
				continue
			}

			if strings.ToLower(filepath.Ext(filename)) == ".md" {
				name, key, err := this.readOneMd(fpath)
				if err != nil {
					logs.Error(err.Error())
					continue
				}

				child := HelpToc{Id: this.getId(), Name: name, Key: key}
				node.Children = append(node.Children, &child)
			}
		}
	}

	return
}

func (this *Help) listFiles(dirname string) []string {
	f, _ := os.Open(dirname)
	names, _ := f.Readdirnames(-1)
	f.Close()
	sort.Strings(names)
	return names
}

func (this *Help) readOneMd(path string) (string, string, error) {
	fExt := filepath.Ext(path)
	fName := filepath.Base(path)
	fName = strings.TrimSuffix(fName, fExt)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", errors.New("read file:" + path + " failed, " + err.Error())
	}

	opts := this.rendererOpts
	doc := goMdParser.NewWithExtensions(this.extensions).Parse(f)
	title := this.firstHeaderText(doc)
	if title == "" {
		title = fName
	}

	body := markdown.Render(doc, goMdHtml.NewRenderer(opts))
	body = this.policy.SanitizeBytes(body)

	//mdBuf := bytes.NewBuffer(f[:0])
	//if err := this.md.Convert(f, mdBuf); err != nil {
	//	return "","",err
	//}

	//body := this.policy.SanitizeBytes(mdBuf.Bytes())

	//logs.Trace("title:",title)

	page := struct {
		Title     string
		StyleHref string
		Style     template.CSS
		Body      template.HTML
		WithHL    bool
	}{
		Title:  title,
		Style:  style,
		Body:   template.HTML(body),
		WithHL: true,
	}

	buf := bytes.NewBuffer(f[:0]) // reuse b to reduce allocations
	if err := contentTemplate.Execute(buf, page); err != nil {
		return "", "", errors.New("file:" + path + " html template execute failed, " + err.Error())
	}

	htmlBuf := buf.Bytes()
	htmlBuf = bytes.Replace(htmlBuf, []byte("/Users/alphawolf/Personal/code/go/src/alphawolf.com/alphaservice/docs/img"), []byte(this.imgDir), -1)

	fName = this.FixName(fName)
	_, ok := this.DocMap.Load(fName)
	if ok {
		return "", "", errors.New("duplicate help doc name: " + fName)
	} else {
		this.DocMap.Store(fName, htmlBuf)
	}

	return title, fName, nil
}

func (this *Help) Initialize() {
	exePath, _ := os.Executable()
	pwd := filepath.Dir(exePath)
	docPath := filepath.Join(pwd, "docs")

	this.md = goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.TaskList),
		goldmark.WithParserOptions(),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	this.extensions = goMdParser.CommonExtensions | goMdParser.Attributes | goMdParser.OrderedListStart | goMdParser.SuperSubscript | goMdParser.Mmark | goMdParser.HardLineBreak | goMdParser.EmptyLinesBreakList
	this.rendererOpts = goMdHtml.RendererOptions{Flags: goMdHtml.CommonFlags | goMdHtml.TOC}
	this.policy = bluemonday.UGCPolicy().AllowAttrs("class").OnElements("code")
	this.imgDir = beego.AppConfig.String("help::img_dir")

	this.toc = &HelpToc{this.getId(), "docs", "", []*HelpToc{}}
	fileInfo, _ := os.Lstat(docPath)

	this.walk(docPath, fileInfo, this.toc)

	//tocBytes, _ := util.JsonIter.Marshal(this.toc)
	//tocString := string(tocBytes)
	//logs.Trace(tocString)
	return

}

func (this *Help) IsHelpName(name string) bool {
	if len(name) > 0 && name[0] == '_' {
		return true
	} else {
		return false
	}
}

func (this *Help) FixName(name string) string {
	ns := strings.Split(name, "_")
	nsLen := len(ns)
	if nsLen > 0 {
		return ns[nsLen-1]
	} else {
		return ""
	}
}

func (this *Help) GetToc() []*HelpToc {
	if this.toc != nil && this.toc.Children != nil {
		return this.toc.Children
	} else {
		toc := make([]*HelpToc, 0)
		return toc
	}
}

func (this *Help) Refresh() {
	this.DocMap.Range(func(key, value interface{}) bool {
		this.DocMap.Delete(key)
		return true
	})

	this.Initialize()
}

func (this *Help) firstHeaderText(doc ast.Node) string {
	var title string
	walkFn := func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}
		switch n := node.(type) {
		case *ast.Heading:
			if n.Level != 1 {
				return ast.GoToNext
			}
			title = string(this.childLiterals(n))
			return ast.Terminate
		case *ast.Code, *ast.CodeBlock, *ast.BlockQuote:
			return ast.SkipChildren
		}
		return ast.GoToNext
	}
	_ = ast.Walk(doc, ast.NodeVisitorFunc(walkFn))
	return title
}

func (this *Help) childLiterals(node ast.Node) []byte {
	if l := node.AsLeaf(); l != nil {
		return l.Literal
	}
	var out [][]byte
	for _, n := range node.GetChildren() {
		if lit := this.childLiterals(n); lit != nil {
			out = append(out, lit)
		}
	}
	if out == nil {
		return nil
	}
	return bytes.Join(out, nil)
}

const style = `body {
	font-family: Charter, Constantia, serif;
	font-size: 1rem;
	line-height: 170%;
	max-width: 45em;
	margin: auto;
	padding-right: 1em;
	padding-left: 1em;
	color: #333;
	background: white;
	text-rendering: optimizeLegibility;
}

@media only screen and (max-width: 480px) {
	body {
		font-size: 125%;
		text-rendering: auto;
	}
}

a {color: #a08941; text-decoration: none;}
a:hover {color: #c6b754; text-decoration: underline;}

h1 a, h2 a, h3 a, h4 a, h5 a {
	text-decoration: none;
	color: gray;
	break-after: avoid;
}
h1 a:hover, h2 a:hover, h3 a:hover, h4 a:hover, h5 a:hover {
	text-decoration: none;
	color: gray;
}
h1, h2, h3, h4, h5 {
	font-weight: bold;
	color: gray;
}

h1 {
	font-size: 150%;
}

h2 {
	font-size: 130%;
}

h3 {
	font-size: 110%;
}

h4, h5 {
	font-size: 100%;
	font-style: italic;
}

pre {
	background-color: rgb(240,240,240);
    color: #77b9fb;
    background-color: #15151b;
	padding: 0.5em;
	overflow: auto;
}
code, pre {
	font-family: Consolas, "PT Mono", monospace;
}
pre { font-size: 90%; }

hr { border:none; text-align:center; color:gray; }
hr:after {
	content:"\2766";
	display:inline-block;
	font-size:1.5em;
}

dt code {
	font-weight: bold;
}
dd p {
	margin-top: 0;
}

blockquote {
	border-left:thick solid lightgrey;
	color: #111111;
	padding: 0 0.5em;
}

img {display:block;margin:auto;max-width:100%}

table, td, th {
	border:thin solid lightgrey;
	border-collapse:collapse;
	vertical-align:middle;
}
td, th {padding:0.2em 0.5em}
tr:nth-child(even) {background-color: rgba(200,200,200,0.2)}

nav#toc {margin:1em 0 1em 0}
nav#toc summary {font-weight:bold; color:gray}
nav#toc ul:after {
	content:"\2042";
	text-align:center;
	display:block;
	color:gray;
}
nav#toc ul {margin:0; list-style:none; padding-left:0}
nav#toc ul li.h2 {padding-left:1em}
nav#toc ul li.h3 {padding-left:2em}
nav#toc ul li.h4 {padding-left:3em}
nav#toc ul li.h5 {padding-left:4em}
nav#toc ul li.h6 {padding-left:5em}

nav#site {
	font-size:90%;
	text-align:right;
	padding:.5em;
	border-bottom: 1px solid gray;
}
nav#site a:before {content:"\2767\0020"}

footer summary {font-weight:bold; color:gray}

summary {cursor:pointer; outline:none}
summary:only-child {display:none}

@media print {
	nav {display: none}
	pre {overflow-wrap:break-word; white-space:pre-wrap}
}`

var indexTemplate = template.Must(template.New("index").Parse(indexTpl))
var pageTemplate = template.Must(template.New("page").Parse(pageTpl))
var contentTemplate = template.Must(template.New("content").Parse(contentTpl))

const indexTpl = `<!doctype html><head><meta charset="utf-8"><title>{{.Title}}</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
{{if .StyleHref}}<link rel="stylesheet" href="{{.StyleHref}}">{{end -}}
{{if .Style}}<style>{{.Style}}</style>{{end}}</head><body id="mdserver-autoindex">{{if .WithSearch}}<form method="get">
<input type="search" name="q" minlength="3" placeholder="Substring search" autofocus required>
<input type="submit"></form>{{end}}
<h1>{{.Title}}</h1><ul>{{$prev := "."}}
{{range .Index}}{{if ne .Subdir $prev}}{{$prev = .Subdir}}</ul><h2>{{.Subdir}}</h2><ul>{{end}}<li><a href="{{.File}}">{{.Title}}</a></li>
{{end}}</ul></body>
`

const pageTpl = `<!doctype html><head><meta charset="utf-8"><title>{{.Title}}</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
{{if .StyleHref}}<link rel="stylesheet" href="{{.StyleHref}}">{{end -}}
{{if .Style}}<style>{{.Style}}</style>{{end}}<script>
document.addEventListener('DOMContentLoaded', function() {
	htmlTableOfContents();
} );
function htmlTableOfContents( documentRef ) {
	var documentRef = documentRef || document;
	var headings = [].slice.call(documentRef.body.querySelectorAll('article h1, article h2, article h3, article h4, article h5, article h6'));
	if (headings.length < 2) { return };
	var toc = documentRef.querySelector("nav#toc details");
	var ul = documentRef.createElement( "ul" );
	headings.forEach(function (heading, index) {
		var ref = heading.getAttribute( "id" );
		var link = documentRef.createElement( "a" );
		link.setAttribute( "href", "#"+ ref );
		link.textContent = heading.textContent;
		var li = documentRef.createElement( "li" );
		li.setAttribute( "class", heading.tagName.toLowerCase() );
		li.appendChild( link );
		ul.appendChild( li );
	});
	toc.appendChild( ul );
}
</script>{{if .WithHL}}
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/styles/default.min.css" integrity="sha256-zcunqSn1llgADaIPFyzrQ8USIjX2VpuxHzUwYisOwo8=" crossorigin="anonymous" referrerpolicy="no-referrer">
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/highlight.min.js" integrity="sha256-aYTdUrn6Ow1DDgh5JTc3aDGnnju48y/1c8s1dgkYPQ8=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script>
document.addEventListener('DOMContentLoaded', (event) => {
	document.querySelectorAll('pre code[class^="language-"]').forEach((block) => {
		hljs.highlightBlock(block);
	});
});
</script>{{end}}
</head><body><nav id="site"><a href="/?index">目录</a></nav>
<nav id="toc"><details open><summary>内容</summary></details></nav>
<ul id="toc"></ul>
<article>
{{.Body}}
</article></body>
`

//const contentTpl =`<nav id="toc"><details open><summary>内容</summary></details></nav>
//<ul id="toc"></ul>
//<article>
//{{.Body}}
//</article>`

const contentTpl = `{{.Body}}`
