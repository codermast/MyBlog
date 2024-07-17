package daos

import (
	"codermast.com/airbbs/models"
	"errors"
	"math"
)

// CreateArticle 文章发布
func CreateArticle(article *models.Article) error {
	result := DB.Create(article)

	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}

// GetArticle 获取所有文章
func GetArticle(status int) (*[]models.Article, error) {
	var articles []models.Article

	if status == -1 {
		result := DB.Table("articles").Order("id desc").Find(&articles)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := DB.Table("articles").Order("id desc").Where("status = ?", status).Find(&articles)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// 根据作者id获取作者名称
	getAuthorNameListByIds(articles)

	return &articles, nil
}

// DeleteArticleByID 根据 ID 删除指定文章
func DeleteArticleByID(articleId string) error {
	result := DB.Table("articles").Where("id = ?", articleId).Delete(nil)

	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}

// GetArticleByID 查询指定 ID 文章
func GetArticleByID(articleID string) (*models.Article, error) {
	var article models.Article
	result := DB.Table("articles").Where("id = ?", articleID).First(&article)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	// 根据作者id获取作者名称
	getAuthorNameById(&article)
	return &article, nil
}

// UpdateArticleByID 根据 ID 更新指定文章
func UpdateArticleByID(article *models.Article) (*models.Article, error) {
	// 1. 首先根据 ID 查文章是否存在
	articleByID, err := GetArticleByID(article.ID)

	// 2. 查询时异常，即查询失败，即文章不存在
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 维护好作者信息
	article.Author = articleByID.Author

	// 此时文章存在，才进行更新
	result := DB.Save(article)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return article, nil
}

// UpdateArticleListStatusById 根据 ID 批量修改文章状态
func UpdateArticleListStatusById(ids []string, status int) error {
	for _, id := range ids {
		result := DB.Table("articles").Where("id = ?", id).Update("status", status)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func GetArticleListPage(articleListPageRequest *models.ArticleListPageRequest) (models.ArticleListPage, error) {
	// 页号
	pageNumber := articleListPageRequest.PageNumber

	// 页面大小
	pageSize := articleListPageRequest.PageSize

	// 偏移量
	offset := (pageNumber - 1) * pageSize

	var articleListPage models.ArticleListPage

	var totalCount int64

	DB.Table("articles").Where("status = ?", 1).Count(&totalCount)
	// 分页查询
	result := DB.Table("articles").Where("status = ?", 1).Limit(pageSize).Offset(offset).Find(&articleListPage.Articles)
	if result.Error != nil {
		return models.ArticleListPage{}, result.Error
	}

	articleListPage.PageNumber = pageNumber
	articleListPage.PageSize = pageSize
	articleListPage.TotalCount = int(totalCount)
	articleListPage.PageCount = int(math.Ceil(float64(int(totalCount) / pageSize)))

	return articleListPage, nil
}

// 根据 作者 id 获取名称
func getAuthorNameListByIds(articles []models.Article) {
	for i := range articles {
		var author models.Author
		DB.Table("users").Where("id = ?", articles[i].Author).First(&author)

		if author.Nickname != "" {
			articles[i].Author = author.Nickname
		} else {
			articles[i].Author = author.Username
		}
	}
}

// 根据 作者 id 获取名称
func getAuthorNameById(article *models.Article) {
	var author models.Author
	DB.Table("users").Where("id = ?", article.Author).First(&author)

	if author.Nickname != "" {
		article.Author = author.Nickname
	} else {
		article.Author = author.Username
	}
}
