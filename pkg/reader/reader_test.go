package reader

import (
	"log/slog"
	"testing"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/model"
	"github.com/stretchr/testify/require"
)

const testData = `
-- posts table
CREATE TABLE public.posts (
  id uuid NOT NULL,
  name VARCHAR(34) NOT NULL,
  description VARCHAR(514),
  created_at timestamp without time zone NOT NULL,
);

ALTER TABLE ONLY public.posts
  ADD CONSTRAINT posts_pkey PRIMARY KEY (id);

-- users table
CREATE TABLE public.users (
  name VARCHAR(34) NOT NULL,
  id BIGINT NOT NULL AUTO_INCREMENT,
);

ALTER TABLE ONLY public.users
  ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- comments table
CREATE TABLE public.comments (
  user_id BIGINT NOT NULL,
  content VARCHAR(514),
  post_id BIGINT NOT NULL AUTO_INCREMENT,
);

ALTER TABLE ONLY public.comments
  ADD CONSTRAINT comments_pkey PRIMARY KEY (user_id, post_id);

ALTER TABLE public.posts ADD CONSTRAINT
  FOREIGN KEY (user_id)
  REFERENCES public.users(id);

ALTER TABLE public.comments ADD CONSTRAINT
  FOREIGN KEY (post_id)
  REFERENCES public.posts(id);

ALTER TABLE public.comments ADD CONSTRAINT
  FOREIGN KEY (user_id)
  REFERENCES public.users(id);
`

// FIXME: trailing spaces around types
func Test_Read(t *testing.T) {
	global.SetLogger(slog.LevelError)
	res := Read(testData)

	expectedFks := []*model.ForeignKey{
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
	}
	expectedLinks := map[string]*model.EntityLink{
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
	}
	expectedTables := []*model.Table{
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
	}

	require.Equal(t, expectedFks, res.Fks)
	require.Equal(t, expectedLinks, res.Links)
	require.Equal(t, expectedTables, res.Tables)
}
