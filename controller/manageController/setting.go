package manageController

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"irisweb/config"
	"irisweb/model"
	"irisweb/provider"
	"irisweb/request"
	"strings"
)

func SettingSystem(ctx iris.Context) {
	system := config.JsonData.System
	var templateNames []string
	//读取目录
	readerInfos, err := ioutil.ReadDir(fmt.Sprintf("%stemplate", config.ExecPath))
	if err != nil {
		fmt.Println(err)
		//怎么会不存在？
	}
	for _, info := range readerInfos {
		if info.IsDir() {
			templateNames = append(templateNames, info.Name())
		}
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": iris.Map{
			"system":         system,
			"template_names": templateNames,
		},
	})
}

func SettingSystemForm(ctx iris.Context) {
	var req request.SystemConfig
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	if !strings.HasPrefix(req.AdminUri, "/") || len(req.AdminUri) < 2 {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "后台路径需要以/开头",
		})
		return
	}

	req.SiteLogo = strings.Replace(req.SiteLogo, config.JsonData.System.BaseUrl, "", -1)

	//进行一些限制
	if req.TemplateType == config.TemplateTypeSeparate {
		if req.MobileUrl == req.BaseUrl {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  "手机端域名不能和电脑端域名一样",
			})
			return
		} else if req.MobileUrl == "" {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  "你选择了电脑+手机模板类型，请填写手机端域名",
			})
			return
		}
	}

	config.JsonData.System.SiteName = req.SiteName
	config.JsonData.System.SiteLogo = req.SiteLogo
	config.JsonData.System.SiteIcp = req.SiteIcp
	config.JsonData.System.SiteCopyright = req.SiteCopyright
	config.JsonData.System.AdminUri = req.AdminUri
	config.JsonData.System.SiteClose = req.SiteClose
	config.JsonData.System.SiteCloseTips = req.SiteCloseTips
	config.JsonData.System.TemplateName = req.TemplateName
	config.JsonData.System.BaseUrl = req.BaseUrl
	config.JsonData.System.MobileUrl = req.MobileUrl
	config.JsonData.System.TemplateType = req.TemplateType

	err := config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置已更新",
	})
}

func SettingContent(ctx iris.Context) {
	system := config.JsonData.Content

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": system,
	})
}

func SettingContentForm(ctx iris.Context) {
	var req request.ContentConfig
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	req.DefaultThumb = strings.Replace(req.DefaultThumb, config.JsonData.System.BaseUrl, "", -1)

	config.JsonData.Content.RemoteDownload = req.RemoteDownload
	config.JsonData.Content.FilterOutlink = req.FilterOutlink
	config.JsonData.Content.ResizeImage = req.ResizeImage
	config.JsonData.Content.ResizeWidth = req.ResizeWidth
	config.JsonData.Content.ThumbCrop = req.ThumbCrop
	config.JsonData.Content.ThumbWidth = req.ThumbWidth
	config.JsonData.Content.ThumbHeight = req.ThumbHeight
	config.JsonData.Content.DefaultThumb = req.DefaultThumb

	err := config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置已更新",
	})
}

//重建所有的thumb
func SettingThumbRebuild(ctx iris.Context) {
	go provider.ThumbRebuild()
	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "缩略图已更新",
	})
}

func SettingIndex(ctx iris.Context) {
	system := config.JsonData.Index

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": system,
	})
}

func SettingIndexForm(ctx iris.Context) {
	var req request.IndexConfig
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.JsonData.Index.SeoTitle = req.SeoTitle
	config.JsonData.Index.SeoKeywords = req.SeoKeywords
	config.JsonData.Index.SeoDescription = req.SeoDescription

	err := config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置已更新",
	})
}

func SettingNav(ctx iris.Context) {
	navList, _ := provider.GetNavList(false)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": navList,
	})
}

func SettingNavForm(ctx iris.Context) {
	var req request.NavConfig
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	var nav *model.Nav
	var err error
	if req.Id > 0 {
		nav, err = provider.GetNavById(req.Id)
		if err != nil {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  err.Error(),
			})
			return
		}
	} else {
		nav = &model.Nav{
			Status: 1,
		}
	}

	nav.Title = req.Title
	nav.SubTitle = req.SubTitle
	nav.Description = req.Description
	nav.ParentId = req.ParentId
	nav.NavType = req.NavType
	nav.PageId = req.PageId
	nav.Link = req.Link
	nav.Sort = req.Sort
	nav.Status = 1
	if nav.NavType == model.NavTypeSystem {
		//内置菜单
		nav.PageId = req.InnerPageId
	}

	err = nav.Save(config.DB)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置已更新",
	})
}

func SettingNavDelete(ctx iris.Context) {
	var req request.NavConfig
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}
	nav, err := provider.GetNavById(req.Id)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	err = nav.Delete(config.DB)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "导航已删除",
	})
}

func SettingContact(ctx iris.Context) {
	system := config.JsonData.Contact

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": system,
	})
}

func SettingContactForm(ctx iris.Context) {
	var req request.ContactConfig
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	req.Qrcode = strings.Replace(req.Qrcode, config.JsonData.System.BaseUrl, "", -1)

	config.JsonData.Contact.UserName = req.UserName
	config.JsonData.Contact.Cellphone = req.Cellphone
	config.JsonData.Contact.Address = req.Address
	config.JsonData.Contact.Email = req.Email
	config.JsonData.Contact.Wechat = req.Wechat
	config.JsonData.Contact.Qrcode = req.Qrcode

	err := config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置已更新",
	})
}
