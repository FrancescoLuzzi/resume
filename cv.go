package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	TypeEntry  string = "entry"
	TypeSkill  string = "skill"
	TypeAuthor string = "author"
)

type AuthorField struct {
	Firstname string
	Lastname  string
	Email     string
	Phone     string
	Birth     string
	Address   string
	Language  string
	Github    string
	Linkedin  string
	Positions []string
}

type EntryField struct {
	Title   string
	Entries []struct {
		Title       string
		Date        string
		Description string
		Location    string
		Link        string
		Infos       []string
	}
}

type SkillField struct {
	Title  string
	Skills []struct {
		Name   string
		Values []string
	}
}

type GenericField struct {
	Type    string
	Content any
}

type Author struct {
	Firstname, Lastname, Email, Github, Linkedin, Birth, Address, Language, Phone string
	Positions                                                                     []string
}

func LoadContent(path string) ([]GenericField, error) {
	fIn, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fIn.Close()
	fileContent, err := io.ReadAll(fIn)
	if err != nil {
		return nil, err
	}
	var content []yaml.Node
	if err = yaml.Unmarshal(fileContent, &content); err != nil {
		return nil, err
	}
	return unmarshalFields(&content)
}

func unmarshalField[T any](node *yaml.Node) (T, error) {
	var entry T
	err := node.Decode(&entry)
	return entry, err
}

func unmarshalEntry(entryType string, node yaml.Node) (GenericField, error) {
	var result GenericField
	var err error
	result.Type = entryType
	switch entryType {
	case TypeEntry:
		result.Content, err = unmarshalField[EntryField](&node)
		break
	case TypeSkill:
		result.Content, err = unmarshalField[SkillField](&node)
		break
	case TypeAuthor:
		result.Content, err = unmarshalField[AuthorField](&node)
		break
	default:
		err = fmt.Errorf("Unkown type: %s", entryType)
	}
	return result, err
}

func unmarshalFields(root *[]yaml.Node) ([]GenericField, error) {
	// this is not a general implementation, just solving my problem
	var out []GenericField
	var commonInfo struct {
		Type string
	}
	for _, sectionNode := range *root {
		sectionNode.Decode(&commonInfo)
		field, err := unmarshalEntry(commonInfo.Type, sectionNode)
		if err != nil {
			return nil, err
		}
		if field.Type == TypeAuthor {
			// prepend so that the author entry is always at the top
			out = append([]GenericField{field}, out...)
		} else {
			out = append(out, field)
		}
	}
	return out, nil
}
