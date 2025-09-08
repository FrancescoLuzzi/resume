package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
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
	var content []yaml.Node
	if err = yaml.Unmarshal(fileContent, &content); err != nil {
		return result, err
	}
	err = unmarshalCv(&content, &result)
	return result, err
}

func unmarshalCv(root *[]yaml.Node, cv *CV) error {
	// this is not a general implementation, just solving my problem
	var err error
	var commonInfo struct {
		Type string
	}
	for _, sectionNode := range *root {
		sectionNode.Decode(&commonInfo)
		switch commonInfo.Type {
		case TypeEntry:
			var entry EntryField
			err = sectionNode.Decode(&entry)
			cv.Entries = append(cv.Entries, entry)
			break
		case TypeSkill:
			var skill SkillField
			err = sectionNode.Decode(&skill)
			cv.Skills = append(cv.Skills, skill)
			break
		case TypeAuthor:
			err = sectionNode.Decode(&cv.Author)
			break
		case TypeAboutMe:
			err = sectionNode.Decode(&cv.AboutMe)
			break
		default:
			err = fmt.Errorf("Unkown type: %s", commonInfo.Type)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
