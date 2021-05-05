package dao

import (
	"context"
	"fmt"
	"github.com/lk/graphql_demo/graph/model"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)

const (
	author     = "lk"
	project    = "es"
	mappingTpl = `{
	"mappings":{
		"properties":{
			"id": 			{ "type": "long" },
			"name": 		{ "type": "text" },
			"avatar":		{ "type": "text" },
			"info":			{ "type": "text" },
			"score":		{ "type": "long" },
			"role":			{ "type": "text" },
			}
		}
	}`
	esRetryLimit = 3 //bulk 错误重试机制
)

type UserES struct {
	index   string
	mapping string
	client  *elastic.Client
}

var ES *UserES

func NewEsClient() *elastic.Client {
	url := fmt.Sprintf("http:"+"//%s:%s", os.Getenv("HOST"), os.Getenv("ES_PORT"))
	client, err := elastic.NewClient(
		//elastic 服务地址
		elastic.SetURL(url),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		log.Fatalf("Failed to create elastic client: %s", err)
	}
	return client
}

func NewUserES(client *elastic.Client) {
	index := fmt.Sprintf("%s_%s", author, project)
	ES = &UserES{
		client:  client,
		index:   index,
		mapping: mappingTpl,
	}

	ES.init()
}

func (es *UserES) init() {
	ctx := context.Background()
	exists, err := es.client.IndexExists(es.index).Do(ctx)
	if err != nil {
		fmt.Printf("userEs init exist failed err is %s\n", err)
		return
	}

	if !exists {
		_, err := es.client.CreateIndex(es.index).Body(es.mapping).Do(ctx)
		if err != nil {
			fmt.Printf("userEs init failed err is %s\n", err)
			return
		}
	}
}

func (es *UserES) BatchAdd(ctx context.Context, user []*model.TeacherInfo) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchAdd(ctx, user); err != nil {
			fmt.Println("batch add failed ", err)
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchAdd(ctx context.Context, user []*model.TeacherInfo) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		doc := elastic.NewBulkIndexRequest().Id(strconv.FormatUint(uint64(u.ID), 10)).Doc(u)
		req.Add(doc)
	}
	if req.NumberOfActions() < 0 {
		return nil
	}
	res, err := req.Do(ctx)
	if err != nil {
		return err
	}
	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	// 所以如果没有出错，说明全部成功了，直接返回即可
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}

func (es *UserES) BatchUpdate(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchUpdate(ctx, user); err != nil {
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchUpdate(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		u.UpdateTime = uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
		doc := elastic.NewBulkUpdateRequest().Id(strconv.FormatUint(u.ID, 10)).Doc(u)
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}
	res, err := req.Do(ctx)
	if err != nil {
		return err
	}
	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	// 所以如果没有出错，说明全部成功了，直接返回即可
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}

func (es *UserES) BatchDel(ctx context.Context, user []*model.UserEs) error {
	var err error
	for i := 0; i < esRetryLimit; i++ {
		if err = es.batchDel(ctx, user); err != nil {
			continue
		}
		return err
	}
	return err
}

func (es *UserES) batchDel(ctx context.Context, user []*model.UserEs) error {
	req := es.client.Bulk().Index(es.index)
	for _, u := range user {
		doc := elastic.NewBulkDeleteRequest().Id(strconv.FormatUint(u.ID, 10))
		req.Add(doc)
	}

	if req.NumberOfActions() < 0 {
		return nil
	}

	res, err := req.Do(ctx)
	if err != nil {
		return err
	}
	// 任何子请求失败，该 `errors` 标志被设置为 `true` ，并且在相应的请求报告出错误明细
	// 所以如果没有出错，说明全部成功了，直接返回即可
	if !res.Errors {
		return nil
	}
	for _, it := range res.Failed() {
		if it.Error == nil {
			continue
		}
		return &elastic.Error{
			Status:  it.Status,
			Details: it.Error,
		}
	}
	return nil
}

// 根据id 批量获取

func (es *UserES) MGet(ctx context.Context, IDS []uint64) ([]*model.TeacherInfo, error) {
	userES := make([]*model.TeacherInfo, 0, len(IDS))
	idStr := make([]string, 0, len(IDS))
	for _, id := range IDS {
		idStr = append(idStr, strconv.FormatUint(id, 10))
	}
	resp, err := es.client.Search(es.index).Query(
		elastic.NewIdsQuery().Ids(idStr...)).Size(len(IDS)).Do(ctx)

	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}
	for _, e := range resp.Each(reflect.TypeOf(&model.TeacherInfo{})) {
		us := e.(*model.TeacherInfo)
		userES = append(userES, us)
	}
	return userES, nil
}

func (es *UserES) Search(ctx context.Context, filter *model.EsSearch) ([]*model.TeacherInfo, error) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(filter.MustQuery...)
	boolQuery.MustNot(filter.MustNotQuery...)
	boolQuery.Should(filter.ShouldQuery...)
	boolQuery.Filter(filter.Filters...)

	// 当should不为空时，保证至少匹配should中的一项
	if len(filter.MustQuery) == 0 && len(filter.MustNotQuery) == 0 && len(filter.ShouldQuery) > 0 {
		boolQuery.MinimumShouldMatch("1")
	}

	service := es.client.Search().Index(es.index).Query(boolQuery).SortBy(filter.Sorters...).From(filter.From).Size(filter.Size)
	resp, err := service.Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}
	userES := make([]*model.TeacherInfo, 0)
	for _, e := range resp.Each(reflect.TypeOf(&model.TeacherInfo{})) {
		us := e.(*model.TeacherInfo)
		userES = append(userES, us)
	}
	return userES, nil
}
