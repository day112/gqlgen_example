package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/lk/graphql_demo/dao"
	"github.com/lk/graphql_demo/graph/generated"
	"github.com/lk/graphql_demo/graph/model"
	"log"
)

func (r *mutationResolver) UpdateTeacherInfo(ctx context.Context, req model.UpdateInfo) (*model.TeacherInfo, error) {
	dao.DB.Table("teacher_infos").Where("id=?", req.ID).Updates(model.TeacherInfo{Info: req.Info, Name: req.Name, Avatar: req.Avatar})
	var teacherInfo model.TeacherInfo
	dao.DB.Preload("TeacherScore").Where("name=?", req.Name).Find(&teacherInfo)

	return &teacherInfo, nil
}

func (r *mutationResolver) InsertTeacherInfo(ctx context.Context, req *model.InsertInfo) ([]*model.TeacherInfo, error) {
	var insert model.TeacherInfo
	insert.Name = req.Name
	insert.Info = req.Info
	dao.DB.Model(model.TeacherInfo{}).Create(&insert)
	var teacherInfo []*model.TeacherInfo
	dao.DB.Preload("TeacherScore").Find(&teacherInfo)

	err := dao.ES.BatchAdd(ctx, teacherInfo)
	if err != nil {
		log.Panicf("批量添加失败：%s", err)
	}

	fmt.Println("插入ID：", insert.ID)
	return teacherInfo, nil
}

func (r *queryResolver) TeacherInfo(ctx context.Context) ([]*model.TeacherInfo, error) {
	var teacherInfos []*model.TeacherInfo
	// mysql 查询
	//dao.DB.Table("teacher_infos").Find(&teacherInfos)

	// es 查询
	ids := []uint64{1, 2, 5, 6}
	teacherInfos, _ = dao.ES.MGet(ctx, ids)

	return teacherInfos, nil
}

func (r *queryResolver) TeacherScore(ctx context.Context, id int) (*model.TeacherInfo, error) {
	var teacherInfo model.TeacherInfo
	dao.DB.Debug().Preload("TeacherScore").First(&teacherInfo)

	return &teacherInfo, nil
}

func (r *queryResolver) SearchTeacherInfos(ctx context.Context, req model.Parms) ([]*model.TeacherInfo, error) {
	var teacherInfos []*model.TeacherInfo
	dao.DB.Preload("TeacherScore").Where("name=? and role=?", req.Name, req.Role).Find(&teacherInfos)
	return teacherInfos, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
