package writer

import (
	"log/slog"
	"testing"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/model"
	"github.com/stretchr/testify/require"
)

var input = &UMLTree{
	Tables: []*model.Table{
		{
			Name: "public.posts",
			Pk:   nil,
			Fks:  nil,
			Fields: []*model.Field{
				{
					Name:         "id",
					Type:         "uuid ",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
				{
					Name:         "name",
					Type:         "VARCHAR[34]",
					IsPrimaryKey: false,
					IsNullable:   false,
				},
				{
					Name:         "description",
					Type:         "VARCHAR[514]",
					IsPrimaryKey: false,
					IsNullable:   true,
				},
				{
					Name:         "created_at",
					Type:         "timestamp without time zone ",
					IsPrimaryKey: false,
					IsNullable:   false,
				},
			},
			FieldsByName: map[string]*model.Field{
				"id": {
					Name:         "id",
					Type:         "uuid ",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
				"name": {
					Name:         "name",
					Type:         "VARCHAR[34]",
					IsPrimaryKey: false,
					IsNullable:   false,
				},
				"description": {
					Name:         "description",
					Type:         "VARCHAR[514]",
					IsPrimaryKey: false,
					IsNullable:   true,
				},
				"created_at": {
					Name:         "created_at",
					Type:         "timestamp without time zone ",
					IsPrimaryKey: false,
					IsNullable:   false,
				},
			},
		},
		{
			Name: "public.users",
			Pk:   nil,
			Fks:  nil,
			Fields: []*model.Field{
				{
					Name:         "name",
					Type:         "VARCHAR[34]",
					IsPrimaryKey: false,
					IsNullable:   false,
				},
				{
					Name:         "id",
					Type:         "BIGINT  AUTO_INCREMENT",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
			},
			FieldsByName: map[string]*model.Field{
				"name": {
					Name:         "name",
					Type:         "VARCHAR[34]",
					IsPrimaryKey: false,
					IsNullable:   false,
				},
				"id": {
					Name:         "id",
					Type:         "BIGINT  AUTO_INCREMENT",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
			},
		},
		{
			Name: "public.comments",
			Pk:   nil,
			Fks:  nil,
			Fields: []*model.Field{
				{
					Name:         "user_id",
					Type:         "BIGINT ",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
				{
					Name:         "content",
					Type:         "VARCHAR[514]",
					IsPrimaryKey: false,
					IsNullable:   true,
				},
				{
					Name:         "post_id",
					Type:         "BIGINT  AUTO_INCREMENT",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
			},
			FieldsByName: map[string]*model.Field{
				"user_id": {
					Name:         "user_id",
					Type:         "BIGINT ",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
				"content": {
					Name:         "content",
					Type:         "VARCHAR[514]",
					IsPrimaryKey: false,
					IsNullable:   true,
				},
				"post_id": {
					Name:         "post_id",
					Type:         "BIGINT  AUTO_INCREMENT",
					IsPrimaryKey: true,
					IsNullable:   false,
				},
			},
		},
	},
	Fks: []*model.ForeignKey{
		{
			SourceTable:       "public.posts",
			SourceFields:      []string{"user_id"},
			DestinationTable:  "public.users",
			DestinationFields: []string{"id"},
		},
		{
			SourceTable:       "public.comments",
			SourceFields:      []string{"post_id"},
			DestinationTable:  "public.posts",
			DestinationFields: []string{"id"},
		},
		{
			SourceTable:       "public.comments",
			SourceFields:      []string{"user_id"},
			DestinationTable:  "public.users",
			DestinationFields: []string{"id"},
		},
	},
	Links: map[string]*model.EntityLink{
		"public.posts_public.users": {
			Left: &model.Link{
				SourceName:      "public.posts",
				DestinationName: "public.users",
				IsNullable:      false,
			},
			Right: nil,
		},
		"public.comments_public.posts": {
			Left: &model.Link{
				SourceName:      "public.comments",
				DestinationName: "public.posts",
				IsNullable:      false,
			},
			Right: nil,
		},
		"public.comments_public.users": {
			Left: &model.Link{
				SourceName:      "public.comments",
				DestinationName: "public.users",
				IsNullable:      false,
			},
			Right: nil,
		},
	},
}

const postsTable = `entity public.posts {
	* id, PK, uuid 
--
	* created_at, timestamp without time zone 
	  description, VARCHAR[514]
	* name, VARCHAR[34]
}`

const usersTable = `entity public.users {
	* id, PK, BIGINT  AUTO_INCREMENT
--
	* name, VARCHAR[34]
}`
const commentsTable = `entity public.comments {
	* post_id, PK, BIGINT  AUTO_INCREMENT
	* user_id, PK, BIGINT 
--
	  content, VARCHAR[514]
}`
const linkPostUsers = `public.posts }|--|| public.users`
const linkCommentsPosts = `public.comments }|--|| public.posts`
const linkCommentsUsers = `public.comments }|--|| public.users`

func Test_Build(t *testing.T) {
	global.SetLogger(slog.LevelError)

	res := input.Build()

	require.Contains(t, res, postsTable)
	require.Contains(t, res, usersTable)
	require.Contains(t, res, commentsTable)
	require.Contains(t, res, linkPostUsers)
	require.Contains(t, res, linkCommentsPosts)
	require.Contains(t, res, linkCommentsUsers)
}
