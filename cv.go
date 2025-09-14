package main

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

const (
	TypeEntry   string = "entry"
	TypeSkill   string = "skill"
	TypeAuthor  string = "author"
	TypeAboutMe string = "about-me"
)

type GenericField struct {
	Type    string
	Content any
}

type CV struct {
	Author  AuthorField
	AboutMe AboutMeField
	Entries []EntryField
	Skills  []SkillField
}

func (c *CV) AsGenericFields() []GenericField {
	var result []GenericField
	result = append(result, GenericField{TypeAuthor, c.Author})
	result = append(result, GenericField{TypeAboutMe, c.AboutMe})
	for _, entry := range c.Entries {
		result = append(result, GenericField{TypeEntry, entry})
	}
	for _, skill := range c.Skills {
		result = append(result, GenericField{TypeSkill, skill})
	}
	return result
}

type AuthorField struct {
	Address   string
	Birth     string
	Email     string
	Firstname string
	Github    string
	Language  string
	Lastname  string
	Phone     string
	Linkedin  string
	Positions []string
}

type AboutMeField struct {
	Title       string
	Description string
}

type EntryField struct {
	Title   string
	Entries []struct {
		Title       string
		Date        string
		Description string
		Location    string
		Link        string
		Infos       string
	}
}

type SkillField struct {
	Title  string
	Skills []struct {
		Name   string
		Values []string
	}
}

type Author struct {
	Firstname string
	Lastname  string
	Email     string
	Github    string
	Linkedin  string
	Birth     string
	Address   string
	Language  string
	Phone     string
	Positions []string
}

func LoadCvYamlFile(path string) (CV, error) {
	var result CV

	fIn, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer fIn.Close()

	fileContent, err := io.ReadAll(fIn)
	if err != nil {
		return result, err
	}

	// parse to AST
	f, err := parser.ParseBytes(fileContent, 0)
	if err != nil {
		return result, err
	}

	for _, doc := range f.Docs {
		if doc == nil || doc.Body == nil {
			continue
		}
		if seq, ok := doc.Body.(*ast.SequenceNode); ok {
			err = unmarshalCvFromNodes(seq.Values, &result)
		} else {
			err = unmarshalCvFromNodes([]ast.Node{doc.Body}, &result)
		}
		return result, err
	}

	return result, fmt.Errorf("no document found in YAML")
}

func unmarshalCvFromNodes(nodes []ast.Node, cv *CV) error {
	var err error
	var commonInfo struct {
		Type string `yaml:"type"`
	}
	for _, node := range nodes {
		if node == nil {
			continue
		}
		err = yaml.NodeToValue(node, &commonInfo)
		switch commonInfo.Type {
		case TypeEntry:
			var entry EntryField
			err = yaml.NodeToValue(node, &entry)
			cv.Entries = append(cv.Entries, entry)
			break

		case TypeSkill:
			var skill SkillField
			err = yaml.NodeToValue(node, &skill)
			cv.Skills = append(cv.Skills, skill)

		case TypeAuthor:
			err = yaml.NodeToValue(node, &cv.Author)
			break

		case TypeAboutMe:
			err = yaml.NodeToValue(node, &cv.AboutMe)
			break

		default:
			return fmt.Errorf("unknown type: %s", commonInfo.Type)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
