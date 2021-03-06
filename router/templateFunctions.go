package router

import (
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/nicksnyder/go-i18n/i18n"
	"html/template"
	"log"
	"math"
	"net/url"
	"strconv"
)

var FuncMap = template.FuncMap{
	"min": math.Min,
	"genRoute": func(name string, params ...string) string {
		url, err := Router.Get(name).URL(params...)
		if err == nil {
			return url.String()
		}
		return "error"
	},
	"genRouteWithQuery": func(name string, currentUrl *url.URL, params ...string) template.HTML {
		url, err := Router.Get(name).URL(params...)
		if err == nil {
			return template.HTML(template.HTMLEscapeString(url.String() + "?" + currentUrl.RawQuery)) // TODO: Review application of character escaping
		}
		return "error"
	},
	"genNav": func(nav Navigation, currentUrl *url.URL, pagesSelectable int) template.HTML {
		var ret = ""
		if (nav.TotalItem > 0) {
		maxPages := math.Ceil(float64(nav.TotalItem) / float64(nav.MaxItemPerPage))

		if nav.CurrentPage-1 > 0 {
			url, _ := Router.Get(nav.Route).URL("page", "1")
			ret = ret + "<li><a id=\"page-prev\" href=\"" + url.String() + "?" + currentUrl.RawQuery + "\" aria-label=\"Previous\"><span aria-hidden=\"true\">&laquo;</span></a></li>"
		}
		startValue := 1
		if nav.CurrentPage > pagesSelectable/2 {
			startValue = (int(math.Min((float64(nav.CurrentPage)+math.Floor(float64(pagesSelectable)/2)), maxPages)) - pagesSelectable + 1)
		}
		endValue := (startValue + pagesSelectable - 1)
		if endValue > int(maxPages) {
			endValue = int(maxPages)
		}
		log.Println(nav.TotalItem)
		for i := startValue; i <= endValue; i++ {
			pageNum := strconv.Itoa(i)
			url, _ := Router.Get(nav.Route).URL("page", pageNum)
			ret = ret + "<li"
			if i == nav.CurrentPage {
				ret = ret + " class=\"active\""
			}

			ret = ret + "><a href=\"" + url.String() + "?" + currentUrl.RawQuery + "\">" + strconv.Itoa(i) + "</a></li>"
		}
		if nav.CurrentPage < int(maxPages) {
			url, _ := Router.Get(nav.Route).URL("page", strconv.Itoa(nav.CurrentPage+1))
			ret = ret + "<li><a id=\"page-next\" href=\"" + url.String() + "?" + currentUrl.RawQuery + "\" aria-label=\"Next\"><span aria-hidden=\"true\">&raquo;</span></a></li>"
		}
		}
		return template.HTML(ret)
	},
	"T": i18n.IdentityTfunc,
	"Ts": i18n.IdentityTfunc,
	"getAvatar": func(hash string, size int) string {
		return "https://www.gravatar.com/avatar/" + hash + "?s=" + strconv.Itoa(size)
	},
	"CurrentOrAdmin":       userPermission.CurrentOrAdmin,
	"CurrentUserIdentical": userPermission.CurrentUserIdentical,
	"HasAdmin":             userPermission.HasAdmin,
	"GetRole":              userPermission.GetRole,
	"IsFollower":           userPermission.IsFollower,
	"NoEncode": func(str string) template.HTML {
		return template.HTML(str)
	},
	"calcWidthSeed": func(seed uint32, leech uint32) float64 {
		return float64(float64(seed)/(float64(seed)+float64(leech)))*100
	},
	"calcWidthLeech": func(seed uint32, leech uint32) float64 {
		return float64(float64(leech)/(float64(seed)+float64(leech)))*100
	},
}
