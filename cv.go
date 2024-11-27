package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	EntryType string = "entry"
	SkillType string = "skill"
)

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

func LoadAuthor(path string) (*Author, error) {
	fIn, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fIn.Close()
	fileContent, err := io.ReadAll(fIn)
	if err != nil {
		return nil, err
	}
	author := Author{}
	if err = yaml.Unmarshal(fileContent, &author); err != nil {
		return nil, err
	}
	return &author, nil
}

func LoadContent(path string) (*[]GenericField, error) {
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

func unmarshalEntries(node *yaml.Node, name string) (EntryField, error) {
	var entries EntryField
	err := node.Decode(&entries)
	entries.Title = name
	return entries, err
}
func unmarshalSkills(node *yaml.Node, name string) (SkillField, error) {
	var skills SkillField
	err := node.Decode(&skills)
	skills.Title = name
	return skills, err
}

func unmarshalFields(root *[]yaml.Node) (*[]GenericField, error) {
	// this is not a general implementation, just solving my problem
	var out []GenericField
	var commonInfo struct {
		Type  string
		Title string
	}
	for _, sectionNode := range *root {
		sectionNode.Decode(&commonInfo)
		foundType := commonInfo.Type
		switch foundType {
		case EntryType:
			entry, err := unmarshalEntries(&sectionNode, commonInfo.Title)
			if err != nil {
				return nil, err
			}
			out = append(out, GenericField{Type: commonInfo.Type, Content: entry})
			break
		case SkillType:
			skill, err := unmarshalSkills(&sectionNode, commonInfo.Title)
			if err != nil {
				return nil, err
			}
			out = append(out, GenericField{Type: commonInfo.Type, Content: skill})
		default:
			panic(fmt.Sprintf("Unkown type: %s", commonInfo.Type))
		}
	}
	return &out, nil
}
