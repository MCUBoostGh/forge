package contents

const CMakePresetsContent = `{
  "version": 3,
  "cmakeMinimumRequired": {
	"major": 3,
	"minor": 20
  },
  "configurePresets": [
	{
	  "name": "default",
	  "hidden": true,
	  "generator": "Ninja",
	  "binaryDir": "${sourceDir}/build"
	}
  ]
}
`