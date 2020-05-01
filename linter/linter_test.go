package linter

import (
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

func TestLint_IsNotText(t *testing.T) {
	// the 0xFFFD UTF-8 control character should make the 'IsText' check fail
	got := IsText([]byte(string([]rune{0xFFFD, '👋'})))

	if got != false {
		t.Errorf("NewEndOfFileRule(defaultTestConf).Lint(textNotText):\n\tExpected %v, got %v", false, got)
	}
}
