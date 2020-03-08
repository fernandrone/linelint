package rules

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
	got, fix, _ := NewEndOfFileRule().Lint(strings.NewReader(textWithSingleNewLine))

	if fix != nil {
		t.Errorf("NewEndOfFileRule().Lint(textWithSingleNewLine):\n\tExpected nil, got:\n%v", string(fix))
	}

	if got != true {
		t.Errorf("NewEndOfFileRule().Lint(textWithSingleNewLine):\n\tExpected %v, got %v", true, got)
	}
}

func TestLint_ShortTextWithSingleNewLine(t *testing.T) {
	got, fix, _ := NewEndOfFileRule().Lint(strings.NewReader(shortTextWithSingleNewLine))

	if fix != nil {
		t.Errorf("NewEndOfFileRule().Lint(shortTextWithSingleNewLine):\n\tExpected nil, got:\n%v", string(fix))
	}

	if got != true {
		t.Errorf("NewEndOfFileRule().Lint(shortTextWithSingleNewLine):\n\tExpected %v, got %v", true, got)
	}
}

func TestLint_TextWithTwoNewLines(t *testing.T) {
	got, fixed, _ := NewEndOfFileRule().Lint(strings.NewReader(textWithTwoNewLines))

	if got != false {
		t.Errorf("NewEndOfFileRule().Lint(textWithTwoNewLines):\n\tExpected %v, got %v", false, got)
	}

	if string(fixed) != textWithSingleNewLine {
		t.Errorf("NewEndOfFileRule().Lint(textWithTwoNewLines): autofix did not work\n\tExpected:\n%q\n\tGot:\n%q", textWithSingleNewLine, string(fixed))
	}
}

func TestLint_TextWithoutNewLine(t *testing.T) {
	got, fixed, _ := NewEndOfFileRule().Lint(strings.NewReader(textWithoutNewLine))

	if string(fixed) != textWithSingleNewLine {
		t.Errorf("NewEndOfFileRule().Lint(textWithoutNewLine): autofix did not work\n\tExpected:\n%q\n\tGot:\n%q", textWithSingleNewLine, string(fixed))
	}

	if got != false {
		t.Errorf("NewEndOfFileRule().Lint(textWithoutNewLine):\n\tExpected %v, got %v", false, got)
	}
}

func TestLint_NotTextFile(t *testing.T) {
	// the 0xFFFD UTF-8 control character should make the 'IsText' check fail
	got, _, err := NewEndOfFileRule().Lint(strings.NewReader(string([]rune{0xFFFD, '👋'})))

	if err == nil {
		t.Errorf("NewEndOfFileRule().Lint(textNotText):\n\tExpected err, got nil")
	}

	if got != false {
		t.Errorf("NewEndOfFileRule().Lint(textNotText):\n\tExpected %v, got %v", false, got)
	}
}
