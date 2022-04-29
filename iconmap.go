package icons

type IconSet map[string]string

var IconMap = map[string]IconSet{
	"heroicons-outline": IconSetHeroiconsOutline,
	"heroicons-solid":   IconSetHeroiconsSolid,
}
