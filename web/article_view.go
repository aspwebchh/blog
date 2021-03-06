package web

import(
	"net/http"
	"../common"
	"strconv"
	"strings"
	"regexp"
	"../cache"
)


type ArticleView struct{
	http.Handler
}


func ( self * ArticleView ) ServeHTTP( response http.ResponseWriter, request *http.Request ) {
	handleArticleContent(response,request,"html/article_view.html")
}

func handleArticleContent(  response http.ResponseWriter, request *http.Request ,htmlPageTemplate string )  {
	id, error := strconv.Atoi( request.FormValue("id") )
	if error != nil {
		response.Write([]byte("参数错误"))
		return
	}
	var pageContentOut = make(chan string)
	go func() {
		var pageContent = common.GetFileContent(htmlPageTemplate)
		pageContentOut <- pageContent
		close(pageContentOut)
	}()
	var articleDataOut = make( chan map[string]string )
	go func() {
		//var article = model.Article{}
		//var data = article.GetArticle(id)
		var data = cache.GetArticle(id)
		articleDataOut <- data
		close(articleDataOut)
	}()

	var pageContent = <- pageContentOut
	var data = <- articleDataOut

	if ( len(data) == 0 ) || ( data["disabled"] == "1" && !loginModel.IsLogin(request) ) {
		var notFoundPageContent = common.GetFileContent("html/404.html")
		response.WriteHeader(404)
		response.Write([]byte(notFoundPageContent))
		return
	}

	pageContent = strings.Replace(pageContent, "${title}",data["title"],-1)
	pageContent = strings.Replace(pageContent, "${content}",data["content"],-1)
	pageContent = strings.Replace(pageContent,"${time}",data["add_date"],-1)
	pageContent = strings.Replace(pageContent,"${author}",data["author"],-1)
	pageContent = code(pageContent)
	pageContent = setClientResourceVersion(pageContent)
	response.Write([]byte(pageContent))

	go func() {
		articleVisitInfoModel.Add(id, getIpAddr(request))
	}()
}

func code( content string ) string {
	var regex = regexp.MustCompile(`<pre[^>]*>([\s\S]+?)<\/pre>`)
	return regex.ReplaceAllStringFunc(content, func(str string) string {
		var replaceHtmlRegex = regexp.MustCompile(`<\/?(p|span)(>|[\s|\/][^>]*>)`)
		str = replaceHtmlRegex.ReplaceAllString(str,"")
		start := regexp.MustCompile(`<pre[^>]*>`)
		end := regexp.MustCompile(`<\/pre>`)
		str = start.ReplaceAllStringFunc(str, func(str string) string  {
			return str + "<code>"
		})
		str = end.ReplaceAllStringFunc(str, func(str string) string  {
			return "</code>" + str
		})
		return str
	})
}


type ArticleViewForMobile struct{
	http.Handler
}

func ( self * ArticleViewForMobile ) ServeHTTP( response http.ResponseWriter, request *http.Request ) {
	handleArticleContent(response,request,"html/mob_article_view.html")
}

