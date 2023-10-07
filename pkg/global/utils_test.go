package global

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_removeConstraints(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected []string
	}{
		"nominal": {
			input: []string{
				"	CREATE INDEX author_pkey_idx ON ...",
				"CREATE TABLE public.posts (",
				"  id uuid NOT NULL, ",
				"  author_uuid uuid NOT NULL, ",
				"  created_at timestamp without time zone,",
				"  CONSTRAINT valid_id CHECK",
				"  CONSTRAINT valid_author_id CHECK",
				");",
				"CREATE INDEX posts_id_idx ON ...",
				"ALTER TABLE ONLY public.posts",
				"	ADD CONSTRAINT posts_pkey PRIMARY KEY(id)",
			},
			expected: []string{
				"CREATE TABLE public.posts (",
				"  id uuid NOT NULL, ",
				"  author_uuid uuid NOT NULL, ",
				"  created_at timestamp without time zone,",
				");",
				"ALTER TABLE ONLY public.posts",
				"	ADD CONSTRAINT posts_pkey PRIMARY KEY(id)",
			},
		},
	}

	for _, c := range cases {
		require.Equal(t, RemoveConstraints(c.input), c.expected)
	}
}

func Test_removeComments(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected []string
	}{
		"nominal": {
			input: []string{
				"-- define public.posts table",
				"CREATE TABLE public.posts (",
				"  id uuid NOT NULL, -- PK",
				"  author_uuid uuid NOT NULL, -- FK to author table",
				"  created_at timestamp without time zone,",
				"  -- test this is a comment",
				");",
			},
			expected: []string{
				"",
				"CREATE TABLE public.posts (",
				"  id uuid NOT NULL, ",
				"  author_uuid uuid NOT NULL, ",
				"  created_at timestamp without time zone,",
				"  ",
				");",
			},
		},
	}

	for _, c := range cases {
		require.Equal(t, RemoveComments(c.input), c.expected)
	}
}

func Test_splitAndGetFields(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected []string
	}{
		"nominal": {
			input:    "'hey', `its`, me, 'mario`",
			expected: []string{"hey", "its", "me", "mario"},
		},
	}

	for _, c := range cases {
		require.Equal(t, SplitAndGetFields(c.input), c.expected)
	}
}

func Test_unquote(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected string
	}{
		"should convert string with single quotes": {
			input:    "'hey'",
			expected: "hey",
		},
		"should convert string with backtick quotes": {
			input:    "`hey`",
			expected: "hey",
		},
		"should convert string with quotes": {
			input:    " 'hey` ",
			expected: "hey",
		},
		"should do nothing": {
			input:    "hey",
			expected: "hey",
		},
	}

	for _, c := range cases {
		require.Equal(t, Unquote(c.input), c.expected)
	}
}

func Test_cleanUpLine(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected []string
	}{
		"nominal": {
			input:    []string{"      hey", "        its", "\nme", "\tmario", "!"},
			expected: []string{"hey", "its", "me", "mario", "!"},
		},
	}

	for _, c := range cases {
		require.Equal(t, CleanUpLine(c.input), c.expected)
	}
}

func Test_contains(t *testing.T) {
	cases := map[string]struct {
		inputArr     []string
		inputElement string
		expected     bool
	}{
		"should contain": {
			inputArr:     []string{"hey", "its", "me", "mario"},
			inputElement: "mario",
			expected:     true,
		},
		"should not contain": {
			inputArr:     []string{"hey", "its", "me", "mario"},
			inputElement: "luigi",
			expected:     false,
		},
	}

	for _, c := range cases {
		require.Equal(t, Contains(c.inputArr, c.inputElement), c.expected)
	}
}
