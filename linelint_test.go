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
	got, _ := lintFile(strings.NewReader(textWithSingleNewLine))

	if got != true {
		t.Errorf("lintFile(textWithSingleNewLine):\n\tExpected %v, got %v", true, got)
	}
}

func TestLint_ShortTextWithSingleNewLine(t *testing.T) {
	got, _ := lintFile(strings.NewReader(shortTextWithSingleNewLine))

	if got != true {
		t.Errorf("lintFile(textWithSingleNewLine):\n\tExpected %v, got %v", true, got)
	}
}

func TestLint_TextWithTwoNewLines(t *testing.T) {
	got, _ := lintFile(strings.NewReader(textWithTwoNewLines))

	if got != false {
		t.Errorf("lintFile(textWithTwoNewLines):\n\tExpected %v, got %v", false, got)
	}
}

func TestLint_TextWithoutNewLine(t *testing.T) {
	got, _ := lintFile(strings.NewReader(textWithoutNewLine))

	if got != false {
		t.Errorf("lintFile(textWithoutNewLine):\n\tExpected %v, got %v", false, got)
	}
}
