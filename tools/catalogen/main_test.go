package main

import "testing"

const sample = "# Colombia\n\n" +
	"## Unidades de medida\n\nA continuación las unidades.\n\n" +
	"| Key     | Nombre   |\n| ------- | -------- |\n| unit    | Unidad   |\n| service | Servicio |\n\n" +
	"## Tipos de identificación\n\n" +
	"| Id      | Descripción |\n| ------- | ----------- |\n| NIT\\_X  | NIT escaped |\n| CC      | Cédula      |\n\n" +
	"## Ubicaciones\n\n" +
	"| Prov | Ciudades                          |\n| ---- | --------------------------------- |\n| A    | uno <br /> dos <br /> tres        |\n"

func TestParse(t *testing.T) {
	cats := parse(sample)
	// The geo "Ubicaciones" section's only row is a <br />-joined list, which is
	// filtered out, leaving the category empty and therefore dropped.
	if len(cats) != 2 {
		t.Fatalf("want 2 categories, got %d: %+v", len(cats), cats)
	}

	if cats[0].Key != "unidades-de-medida" || len(cats[0].Entries) != 2 {
		t.Fatalf("units category wrong: %+v", cats[0])
	}
	if cats[0].Entries[0] != (entry{Code: "unit", Name: "Unidad"}) {
		t.Fatalf("first unit wrong: %+v", cats[0].Entries[0])
	}

	if cats[1].Key != "tipos-de-identificacion" {
		t.Fatalf("second category key wrong: %q", cats[1].Key)
	}
	// Markdown escaping is stripped.
	if cats[1].Entries[0].Code != "NIT_X" {
		t.Fatalf("escape not stripped: %q", cats[1].Entries[0].Code)
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Unidades de medida":                  "unidades-de-medida",
		"Tipos de identificación":             "tipos-de-identificacion",
		"Estado de emisión (emission_status)": "estado-de-emision",
		"Regímenes":                           "regimenes",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}
