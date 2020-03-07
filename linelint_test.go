package main

import (
	"strings"
	"testing"
)

const textWithSingleNewLine = `
As armas e os Barões assinalados
Que da Ocidental praia Lusitana,
Por mares nunca de antes navegados
Passaram ainda além da Taprobana,
Em perigos e guerras esforçados,
Mais do que prometia a força humana,
E entre gente remota edificaram
Novo reino, que tanto sublimaram;
`

const textWithTwoNewLines = `
As armas e os Barões assinalados
Que da Ocidental praia Lusitana,
Por mares nunca de antes navegados
Passaram ainda além da Taprobana,
Em perigos e guerras esforçados,
Mais do que prometia a força humana,
E entre gente remota edificaram
Novo reino, que tanto sublimaram;

`

const textWithoutNewLine = `
As armas e os Barões assinalados
Que da Ocidental praia Lusitana,
Por mares nunca de antes navegados
Passaram ainda além da Taprobana,
Em perigos e guerras esforçados,
Mais do que prometia a força humana,
E entre gente remota edificaram
Novo reino, que tanto sublimaram;`

const shortTextWithSingleNewLine = `#
`

func TestLint_TextWithSingleNewLine(t *testing.T) {
	got, _ := singleNewLineRule.lint(strings.NewReader(textWithSingleNewLine))

	if got != true {
		t.Errorf("singleNewLineRule.lint(textWithSingleNewLine):\n\tExpected %v, got %v", true, got)
	}
}

func TestLint_ShortTextWithSingleNewLine(t *testing.T) {
	got, _ := singleNewLineRule.lint(strings.NewReader(shortTextWithSingleNewLine))

	if got != true {
		t.Errorf("singleNewLineRule.lint(textWithSingleNewLine):\n\tExpected %v, got %v", true, got)
	}
}

func TestLint_TextWithTwoNewLines(t *testing.T) {
	got, _ := singleNewLineRule.lint(strings.NewReader(textWithTwoNewLines))

	if got != false {
		t.Errorf("singleNewLineRule.lint(textWithTwoNewLines):\n\tExpected %v, got %v", false, got)
	}
}

func TestLint_TextWithoutNewLine(t *testing.T) {
	got, _ := singleNewLineRule.lint(strings.NewReader(textWithoutNewLine))

	if got != false {
		t.Errorf("singleNewLineRule.lint(textWithoutNewLine):\n\tExpected %v, got %v", false, got)
	}
}

func TestLint_NotTextFile(t *testing.T) {
	// the 0xFFFD UTF-8 control character should make the 'IsText' check fail
	got, err := singleNewLineRule.lint(strings.NewReader(string([]rune{0xFFFD, '👋'})))

	if err == nil {
		t.Errorf("singleNewLineRule.lint(textNotText):\n\tExpected err, got nil")
	}

	if got != false {
		t.Errorf("singleNewLineRule.lint(textNotText):\n\tExpected %v, got %v", false, got)
	}
}
