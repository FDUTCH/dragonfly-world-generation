package resourcecompiler


var definitionPath, structurePath string
var resourcesLoaded bool
var resourcesCompiled bool

func SetPaths(DefinitionPath, StructurePath string){
	resourcesLoaded = true

	definitionPath, structurePath = DefinitionPath, StructurePath
	_ = structurePath
}

func CompileAll(){
	if resourcesCompiled{
		return
	}
	
	resourcesCompiled = true
	CompileBiomes()
	CompileFeatures()
	CompileFeatureRules()

	// structures, err := os.ReadDir(structurePath)
	// if err != nil{
	// 	panic(err)
	// }

	// _ = structures
}

func CheckForResources(){
	if !resourcesLoaded{
		panic("Load resources before using the generator\nhttps://github.com/Ikarolyi/dragonfly-world-generation?tab=readme-ov-file#Setup")
	}
}