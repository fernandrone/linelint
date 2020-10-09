package main

import (
	"strings"

	"github.com/fernandrone/linelint/linter"
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

func Example_NoFix() {
	input := Input{
		Paths:  []string{"-"},
		Stdin:  strings.NewReader(textWithSingleNewLine),
		Config: linter.NewConfig(),
	}

	if err := run(input); err != nil {
		panic(err)
	}

	// Output:
}

func Example_Fix() {
	input := Input{
		Paths:  []string{"-"},
		Stdin:  strings.NewReader(textWithTwoNewLines),
		Config: linter.NewConfig(),
	}

	if err := run(input); err != nil {
		panic(err)
	}

	// Output:
	// As armas e os Barões assinalados
	// Que da Ocidental praia Lusitana,
	// Por mares nunca de antes navegados
	// Passaram ainda além da Taprobana,
	// Em perigos e guerras esforçados,
	// Mais do que prometia a força humana,
	// E entre gente remota edificaram
	// Novo reino, que tanto sublimaram;
}
