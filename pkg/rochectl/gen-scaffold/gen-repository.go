package gen_scaffold

import (
	"fmt"
	"go/ast"

	. "github.com/dave/jennifer/jen"
	autoTableSQL "github.com/hourglasshoro/auto-table/pkg/sql"
	rocheAst "github.com/riita10069/roche/pkg/rochectl/ast"
	"github.com/riita10069/roche/pkg/util"
)

func GenerateRepository(name string, targetStruct *ast.StructType, sqlMap map[string]*autoTableSQL.SQL, moduleName string) *File {
	properties, propertiesType := rocheAst.GetPropertyByStructAst(targetStruct)
	dict := GenDict(properties, properties)
	scanArgument := propertyToScan(properties)
	toExecForCreate := propertyToExecForCreate(properties, "e")
	toExecForUpdate := propertyToExecForUpdate(properties, "e")

	infraRepoFile := NewFile("repository")

	infraRepoFile.HeaderComment("Code generated by roche")

	entityPath := moduleName + "/domain/entity"
	domainRepoPath := moduleName + "/domain/repository"
	infraRepoFile.Id("import").Parens(Id("\"" + entityPath + "\"").Id(";").Id("\"" + domainRepoPath + "\"").Id(";").Id("\"database/sql\""))

	var repository = name + "Repository"
	infraRepoFile.Type().Id(name + "Repository").Struct(
		Id("DB").Id("*sql.DB"),
	)

	// NewStructNameUsecase Constructor
	infraRepoFile.Func().Id("New" + name + "Repository").Params(Id("db").Id("*sql.DB")).Id("*" + repository).Block(
		Return(Op("&").Id(repository).Values(Dict{
			Id("DB"): Id("db"),
		})),
	)

	// GetList
	infraRepoFile.Func().Params(Id("r").Id(repository)).Id("GetList").Params().Params(Index().Id("*entity."+name), Error()).Block(
		Var().Id("entities").Id("[]*entity."+name),
		List(Id("rows"), Err()).Op(":=").
			Id("r").Dot("DB").Dot("Query").Call(Id(FindAllSQL(name, sqlMap))),
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
			Err().Op("=").Id("rows").Dot("Scan").Call(scanArgument...),
			If(
				Err().Op("!=").Nil(),
			).Block(
				Return(Id("nil"), Err()),
			),
			Id("e").Op(":=").Op("&").Id("entity."+name).Values(dict),
			Id("entities").Op("=").Id("append").Call(Id("entities"), Id("e")),
		),
		Err().Op("=").Id("rows").Dot("Err").Call(),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),

		Return(Id("entities"), Err()),
	)

	// GetByID
	infraRepoFile.Func().Params(Id("r").Id(repository)).Id("GetByID").Params(Id("id").Int64()).Params(Id("*entity."+name), Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("DB").Dot("Prepare").Call(Id(FindSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		Var().Params(Id(getVarArgumentForScan(properties, propertiesType))),
		Id("err").Op("=").Id("stmt").Dot("QueryRow").Call(Id("1")).Dot("Scan").Call(scanArgument...),
		If(
			Err().Op("!=").Nil(),
		).Block(
			If(
				Id("err").Op("==").Id("sql.ErrNoRows").Block(
					Return(Id("&entity."+name+"{}"), Nil()),
				),
			),
			Return(Id("nil"), Err()),
		),

		Id("e").Op(":=").Op("&").Id("entity."+name).Values(dict),

		Return(Id("e"), Err()),
	)

	// Create
	infraRepoFile.Func().Params(Id("r").Id(repository)).Id("Create").Params(Id("e").Id("*entity."+name)).Params(Id("*entity."+name), Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("DB").Dot("Prepare").Call(Id(CreateSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		List(Id("_"), Id("err")).Op("=").Id("stmt").Dot("Exec").Call(toExecForCreate...),
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

		Return(Id("e"), Err()),
	)

	// Update
	infraRepoFile.Func().Params(Id("r").Id(repository)).Id("Update").Params(Id("id").Int64(), Id("e").Id("*entity."+name)).Params(Id("*entity."+name), Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("DB").Dot("Prepare").Call(Id(UpdateSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Id("nil"), Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		List(Id("_"), Id("err")).Op("=").Id("stmt").Dot("Exec").Call(toExecForUpdate...),
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

		Return(Id("e"), Err()),
	)

	// Delete
	infraRepoFile.Func().Params(Id("r").Id(repository)).Id("Delete").Params(Id("id").Int64()).Params(Error()).Block(
		List(Id("stmt"), Id("err")).Op(":=").Id("r").Dot("DB").Dot("Prepare").Call(Id(DeleteSQL(name, sqlMap))),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),
		Defer().Id("stmt").Dot("Close").Call(),

		List(Id("_"), Id("err")).Op("=").Id("stmt").Dot("Exec").Call(Id("id")),
		If(
			Err().Op("!=").Nil(),
		).Block(
			Return(Err()),
		),

		Return(Err()),
	)

	return infraRepoFile

}

func propertyToExecForUpdate(properties []string, object string) []Code {
	var execForUpdate []Code
	execForUpdate = append(execForUpdate, Id("id"))
	for _, v := range properties {
		execForUpdate = append(execForUpdate, Id(object+"."+util.SnakeToUpperCamel(util.CamelToSnake(v))))
	}
	return execForUpdate
}

func propertyToExecForCreate(properties []string, object string) []Code {
	var execForCreate []Code
	for _, v := range properties {
		execForCreate = append(execForCreate, Id(object+"."+util.SnakeToUpperCamel(util.CamelToSnake(v))))
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
		postSignature += util.SnakeToLowerCamel(util.CamelToSnake(property[i])) + " " + (propertyType[i])
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
