package gen_scaffold

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	autoTableSQL "github.com/hourglasshoro/auto-table/pkg/sql"
	rocheAst "github.com/riita10069/roche/pkg/rochectl/ast"
	"github.com/riita10069/roche/pkg/util"
	"go/ast"
)

func GenerateRepository(name string, targetStruct *ast.StructType, sqlMap map[string]*autoTableSQL.SQL) *File {
	properties, propertiesType := rocheAst.GetPropertyByStructAst(targetStruct)
	createArgument := rocheAst.GetPostArgument(properties, propertiesType)
	updateArgument := append(createArgument, Id("id").Int64())
	dict := GenDict(properties, properties)
	scanArgument := propertyToScan(properties)
	toExecForCreate := propertyToExecForCreate(properties)
	toExecForUpdate := propertyToExecForUpdate(properties)

	infraRepoFile := NewFile("repository")

	infraRepoFile.HeaderComment("Code generated by roche")

	infraRepoFile.Type().Id(name).Struct(
		Id("DB").Id("*sql.DB"),
	)

	// NewStructNameUsecase Constructor
	infraRepoFile.Func().Id("New" + name + "Repository").Params(Id("db").Id("*sql.DB")).Id("repository.I" + name).Block(
		Return(Op("&").Id(name + "Usecase").Values(Dict{
			Id("DB"): Id("db"),
		})),
	)

	// GetList
	infraRepoFile.Func().Params(Id("r").Id(name)).Id("FindAll").Params().Params(Index().Id("*entity."+name), Error()).Block(
		Var().Id("entity").Id("entity."+name),
		Var().Id("entities").Id("[]entity."+name),
		List(Id("rows"), Err()).Op(":=").
			Id("r").Dot("db").Dot("Query").Call(Id(FindAllSQL(name, sqlMap))),
		Defer().Id("rows").Dot("Close").Call(),

		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),

		For(
			Id("rows").Dot("Next").Call(),
		).Block(
			Var().Params(Id(getVarArgumentForScan(properties, propertiesType)+"id int64")),
			Err().Op(":=").Id("rows").Dot("Scan").Call(scanArgument...),
			If(
				Err().Op("!=").Nil(),
			).Block(
				Return(Id("nil"), Err()),
			),
			Id("entity").Op(":=").Op("&").Id("entity."+name).Values(dict),
			Id("entities").Op("=").Id("append").Call(Id("entities"), Id("entity")),
		),
		Err().Op("=").Id("rows").Dot("Err").Call(),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),

		Return(Id("entity"), Err()),
	)

	// GetByID
	infraRepoFile.Func().Params(Id("r").Id(name)).Id("Find").Params(Id("id").Int64()).Params(Id("*entity."+name), Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("db").Dot("Prepare").Call(Id(FindSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		Var().Params(Id(getVarArgumentForScan(properties, propertiesType)+"id int64")),
		List(Id("rows"), Id("err")).Op(":=").Id("prep").Dot("QueryRow").Call(Id("1")).Dot("Scan").Call(scanArgument...),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("rows").Dot("Close").Call(),

		Id("entity").Op(":=").Op("&").Id("entity."+name).Values(dict),

		Return(Id("entity"), Err()),
	)

	// Create
	infraRepoFile.Func().Params(Id("u").Id(name)).Id("Create").Params(createArgument...).Params(Id("*entity."+name), Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("db").Dot("Prepare").Call(Id(CreateSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		List(Id("result"), Id("err")).Op(":=").Id("stmt").Dot("Exec").Call(toExecForCreate...),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),

		//List(Id("id"), Id("err")).Op(":=").Id("result").Dot("LastInsertId").Call(),
		//If(
		//	Err().Op("!=").Nil(),
		//).Block(
		//	Return(Id("nil"), Err()),
		//),

		Id("entity").Op(":=").Op("&").Id("entity."+name).Values(dict),

		Return(Id("entity"), Err()),
	)

	// Update
	infraRepoFile.Func().Params(Id("u").Id(name)).Id("Update").Params(updateArgument...).Params(Id("*entity."+name), Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("db").Dot("Prepare").Call(Id(UpdateSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		List(Id("result"), Id("err")).Op(":=").Id("stmt").Dot("Exec").Call(toExecForUpdate...),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),

		//List(Id("affectedRows"), Id("err")).Op(":=").Id("result").Dot("RowsAffected").Call(),
		//If(
		//	Err().Op("!=").Nil(),
		//).Block(
		//	Return(Id("nil"), Err()),
		//),

		Id("entity").Op(":=").Op("&").Id("entity."+name).Values(dict),

		Return(Id("entity"), Err()),
	)

	// Delete
	infraRepoFile.Func().Params(Id("u").Id(name)).Id("Delete").Params(Id("id").Int64()).Params(Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("db").Dot("Prepare").Call(Id(DeleteSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		List(Id("result"), Id("err")).Op(":=").Id("stmt").Dot("Exec").Call(Id("id")),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),

		Return(Err()),
	)

	return infraRepoFile

}

func propertyToExecForUpdate(properties []string) []Code {
	var execForUpdate []Code
	execForUpdate = append(execForUpdate, Id("id"))
	for _, v := range properties {
		execForUpdate = append(execForUpdate, Id(util.SnakeToLowerCamel(util.CamelToSnake(v))))
	}
	return execForUpdate
}

func propertyToExecForCreate(properties []string) []Code {
	var execForCreate []Code
	for _, v := range properties {
		execForCreate = append(execForCreate, Id(util.SnakeToLowerCamel(util.CamelToSnake(v))))
	}
	return execForCreate
}

func propertyToScan(properties []string) []Code {
	var scanSignature []Code
	scanSignature = append(scanSignature, Id("&id"))
	for _, v := range properties {
		scanSignature = append(scanSignature, Id("&"+util.SnakeToLowerCamel(util.CamelToSnake(v))))
	}
	return scanSignature
}

func getVarArgumentForScan(property []string, propertyType []string) string {
	var postSignature string
	for i := range property {
		postSignature += util.CamelToSnake(util.SnakeToLowerCamel(util.CamelToSnake(property[i]))) + " " + (propertyType[i])
		postSignature += ";"
	}
	return postSignature
}

func FindAllSQL(name string, sqlMap map[string]*autoTableSQL.SQL) string {
	return fmt.Sprintf("\"%s\"", sqlMap[util.CamelToSnake(name)].Record.FindAll)
}

func FindSQL(name string, sqlMap map[string]*autoTableSQL.SQL) string {
	return fmt.Sprintf("\"%s\"", sqlMap[util.CamelToSnake(name)].Record.Find)
}

func CreateSQL(name string, sqlMap map[string]*autoTableSQL.SQL) string {
	return fmt.Sprintf("\"%s\"", sqlMap[util.CamelToSnake(name)].Record.Create)
}

func UpdateSQL(name string, sqlMap map[string]*autoTableSQL.SQL) string {
	return fmt.Sprintf("\"%s\"", sqlMap[util.CamelToSnake(name)].Record.Update)
}

func DeleteSQL(name string, sqlMap map[string]*autoTableSQL.SQL) string {
	return fmt.Sprintf("\"%s\"", sqlMap[util.CamelToSnake(name)].Record.Delete)
}
