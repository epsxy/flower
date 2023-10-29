package writer

import (
	"log/slog"
	"testing"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/model"
	"github.com/stretchr/testify/require"
)

var tablePosts = model.Table{
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
}

var tableUsers = model.Table{
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
}

var tableComments = model.Table{
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
}

var input = &model.UMLTree{
	Options: &model.UMLTreeOptions{
		SplitUnconnected: false,
		SplitDistance:    false,
	},
	Tables: []*model.Table{
		&tablePosts,
		&tableUsers,
		&tableComments,
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
	TablesByName: map[string]*model.Table{
		tablePosts.Name:    &tablePosts,
		tableUsers.Name:    &tableUsers,
		tableComments.Name: &tableComments,
	},
	LinksByTableName: map[string][]*model.EntityLink{
		tablePosts.Name: {
			{
				Left: &model.Link{
					SourceName:      "public.posts",
					DestinationName: "public.users",
					IsNullable:      false,
				},
				Right: nil,
			},
			{
				Left: &model.Link{
					SourceName:      "public.comments",
					DestinationName: "public.posts",
					IsNullable:      false,
				},
				Right: nil,
			},
		},
		tableUsers.Name: {
			{
				Left: &model.Link{
					SourceName:      "public.posts",
					DestinationName: "public.users",
					IsNullable:      false,
				},
				Right: nil,
			},
			{
				Left: &model.Link{
					SourceName:      "public.comments",
					DestinationName: "public.users",
					IsNullable:      false,
				},
				Right: nil,
			},
		},
		tableComments.Name: {
			{
				Left: &model.Link{
					SourceName:      "public.comments",
					DestinationName: "public.posts",
					IsNullable:      false,
				},
				Right: nil,
			},
			{
				Left: &model.Link{
					SourceName:      "public.comments",
					DestinationName: "public.users",
					IsNullable:      false,
				},
				Right: nil,
			},
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

	res := Build(input)

	require.Equal(t, len(res), 1)
	require.Contains(t, res[0], postsTable)
	require.Contains(t, res[0], usersTable)
	require.Contains(t, res[0], commentsTable)
	require.Contains(t, res[0], linkPostUsers)
	require.Contains(t, res[0], linkCommentsPosts)
	require.Contains(t, res[0], linkCommentsUsers)
}

func Test_generateDocumentName(t *testing.T) {
	cases := map[string]struct {
		vertexes []string
		expected string
	}{
		"single table": {
			vertexes: []string{"test_table_name"},
			expected: "test_table_name",
		},
		"multiple tables": {
			vertexes: []string{"table_a", "table_b"},
			expected: "tables",
		},
	}
	for name, c := range cases {
		require.Equal(t, _generateDocumentName(c.vertexes), c.expected, name)
	}
}

func Test_WriteLink(t *testing.T) {
	cases := map[string]struct {
		link     *model.EntityLink
		expected string
	}{
		"link: 1 <-> 1": {
			link: &model.EntityLink{
				Left: &model.Link{
					SourceName:      "table_a",
					DestinationName: "table_b",
					IsNullable:      false,
				},
				Right: &model.Link{
					SourceName:      "table_b",
					DestinationName: "table_a",
					IsNullable:      false,
				},
			},
			expected: "table_a ||--|| table_b\n",
		},
		"link: 0,1 <-> 1": {
			link: &model.EntityLink{
				Left: &model.Link{
					SourceName:      "table_a",
					DestinationName: "table_b",
					IsNullable:      false,
				},
				Right: &model.Link{
					SourceName:      "table_b",
					DestinationName: "table_a",
					IsNullable:      true,
				},
			},
			expected: "table_a |o--|| table_b\n",
		},
		"link: 1 <-> 0,1": {
			link: &model.EntityLink{
				Left: &model.Link{
					SourceName:      "table_a",
					DestinationName: "table_b",
					IsNullable:      true,
				},
				Right: &model.Link{
					SourceName:      "table_b",
					DestinationName: "table_a",
					IsNullable:      false,
				},
			},
			expected: "table_a ||--o| table_b\n",
		},
		"link: 0,1 <-> 0,1": {
			link: &model.EntityLink{
				Left: &model.Link{
					SourceName:      "table_a",
					DestinationName: "table_b",
					IsNullable:      true,
				},
				Right: &model.Link{
					SourceName:      "table_b",
					DestinationName: "table_a",
					IsNullable:      true,
				},
			},
			expected: "table_a |o--o| table_b\n",
		},
		"link: 1 -> N": {
			link: &model.EntityLink{
				Left: nil,
				Right: &model.Link{
					SourceName:      "table_b",
					DestinationName: "table_a",
					IsNullable:      false,
				},
			},
			expected: "table_a ||--|{ table_b\n",
		},
		"link: 0,1 -> N": {
			link: &model.EntityLink{
				Left: nil,
				Right: &model.Link{
					SourceName:      "table_b",
					DestinationName: "table_a",
					IsNullable:      true,
				},
			},
			expected: "table_a |o--|{ table_b\n",
		},
		"link: N <- 1": {
			link: &model.EntityLink{
				Left: &model.Link{
					SourceName:      "table_a",
					DestinationName: "table_b",
					IsNullable:      false,
				},
				Right: nil,
			},
			expected: "table_a }|--|| table_b\n",
		},
		"link: N <- 0, 1": {
			link: &model.EntityLink{
				Left: &model.Link{
					SourceName:      "table_a",
					DestinationName: "table_b",
					IsNullable:      true,
				},
				Right: nil,
			},
			expected: "table_a }|--o| table_b\n",
		},
		"empty link": {
			link:     nil,
			expected: "",
		},
	}
	for name, c := range cases {
		require.Equal(t, WriteLink(c.link), c.expected, name)
	}
}
