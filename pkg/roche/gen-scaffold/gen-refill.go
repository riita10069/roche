package gen_scaffold

import . "github.com/dave/jennifer/jen"

func GenRefill(before string, after string, propertyList []string) []*Statement {
	if after == "" {
		return []*Statement{}
	}
	var refilled []*Statement
	if before == "" {
		for _, p := range propertyList {
			refilled = append(refilled, Id(after).Dot(p).Op("=").Id(p))
		}
	} else {
		for _, p := range propertyList {
			refilled = append(refilled, Id(after).Dot(p).Op("=").Id(before).Dot(p))
		}
	}
	return refilled
}

func GenDict(properties []string, values []string) Dict {
	d := Dict{}
	for i, _ := range properties {
		d[Id(properties[i])] = Id(values[i])
	}
	return d
}

