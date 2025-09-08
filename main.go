package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func JoinQuoted(elems []string, sep string) string {
	var b strings.Builder
	b.WriteByte('"')
	b.WriteString(elems[0])
	b.WriteByte('"')
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteByte('"')
		b.WriteString(s)
		b.WriteByte('"')
	}
	return b.String()
}

func formatSkill(skill string) string {
	l := len(skill)
	if l == 0 {
		return skill
	}
	l -= 1
	if skill[l] == '!' {
		return fmt.Sprintf(`strong("%s")`, skill[:l])
	} else {
		return fmt.Sprintf(`"%s"`, skill)
	}
}

func JoinSkills(skills []string) string {
	var b strings.Builder
	b.WriteString(formatSkill(skills[0]))
	for _, s := range skills[1:] {
		b.WriteByte(',')
		b.WriteString(formatSkill(s))
	}
	return b.String()
}

var funcs = template.FuncMap{"joinQuoted": JoinQuoted, "joinSkills": JoinSkills}

func IfErrExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	cv      = flag.String("cv", "locales/it.yaml", "Cv yaml file")
	out     = flag.String("out", "it.typ", "Output typst file")
	imports = flag.String("imports", "requirements.typ", "Input typst file defining required dependencies")
)

func main() {
	flag.Parse()
	fOut, err := os.OpenFile(*out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	IfErrExit(err)
	defer fOut.Close()
	writer := bufio.NewWriter(fOut)
	templates := os.DirFS("templates")
	cvTempl, err := template.New("cv.tmpl").Funcs(funcs).ParseFS(templates, "*.tmpl")
	IfErrExit(err)
	typstImports, err := os.ReadFile(*imports)
	IfErrExit(err)
	_, err = writer.Write(typstImports)
	IfErrExit(err)
	fields, err := LoadContent(*cv)
	IfErrExit(err)
	for _, field := range fields {
		IfErrExit(cvTempl.ExecuteTemplate(writer, field.Type, field.Content))
	}
	writer.Flush()
}
