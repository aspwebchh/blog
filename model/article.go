package model

import (
	"../common"
	"../db_helper"
	"strconv"
)

type Article struct {

}

func (self *Article) GetArticle( id int) (map[string]string)  {
	sql := "select * from article where id = " + strconv.Itoa(id);
	result := DbHelper.GetDataBase().Query(sql);
	if(  len(result) > 0 ) {
		return result[0];
	} else {
		return nil;
	}
}

func (self *Article) GetAll() []map[string]string {
	sql := "select id, title, author ,add_time from article order by id desc";
	return DbHelper.GetDataBase().Query(sql);
}

func (self *Article) GetArticleList( start int, end int, cateId int, kw string) ([]map[string]string, int)  {
	kw = common.Trim(kw);
	limit := strconv.Itoa(start) + "," + strconv.Itoa(end);
	where := " 1=1";
	args := make([]interface{},0);
	if( cateId != 0 ) {
		where += " and id in (select article_id from article_categories where cate_id in (select DISTINCT(id) from article_category where pid = "+ strconv.Itoa(cateId) +"))";
	}
	if( kw != "" ) {
		args = append(args, kw);
		where += " and instr(title,?)";
	}
	countSql := "select count(1) from article where " + where ;
	sql := "select id , title from article where "+ where +" order by id desc limit " + limit;

	dataCount,_:= strconv.Atoi( DbHelper.GetDataBase().GetSingle(countSql,args...) );
	dataList := DbHelper.GetDataBase().Query(sql,args...);
	return dataList,dataCount;
}


func ( self * Article ) GetCates()  []map[string]string {
	sql := "select id, name from article_category where pid = 0 order by num asc";
	result := DbHelper.GetDataBase().Query(sql);
	return result;
}

func ( self * Article ) GetCateName( cateId int ) string {
	sql := "select name from article_category where id = " + strconv.Itoa( cateId );
	result := DbHelper.GetDataBase().GetSingle(sql);
	if( result == "") {
		return "所有文章";
	}
	return result;
}