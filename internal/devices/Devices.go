package devices

import (
	"embed"
	"fmt"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Catalog struct {
	ID            string            `yaml:"id"`
	Vendor        string            `yaml:"vendor"`
	Family        string            `yaml:"family"`
	Series        string            `yaml:"series"`
	Architecture  string            `yaml:"architecture"`
	CPU           string            `yaml:"cpu"`
	FPU           string            `yaml:"fpu"`
	FloatABI      string            `yaml:"float_abi"`
	FlashKB       int               `yaml:"flash_kb"`
	RAMKB         int               `yaml:"ram_kb"`
	Pins          int               `yaml:"pins"`
	Package       string            `yaml:"package"`
	STM32Device   string            `yaml:"stm32_device"`
	LinkerScript  string            `yaml:"linker_script"`
	OpenOCDTarget string            `yaml:"openocd_target"`
	Presets       map[string]Preset `yaml:"presets"`
	Aliases       []string          `yaml:"aliases"`
}

type Preset struct {
	BuildType string `yaml:"build_type"`
}

type Filter struct {
	Vendor string
	Family string
	Series string
}

//go:embed st/*.yaml
var stCatalog embed.FS

var (
	catalog     map[string]Catalog
	aliasIndex  map[string]string
	catalogOnce sync.Once
	catalogErr  error
)

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func ensureCatalog() error {
	catalogOnce.Do(func() {
		catalogErr = loadCatalog()
	})
	return catalogErr
}

func loadCatalog() error {
	catalog = make(map[string]Catalog)
	aliasIndex = make(map[string]string)

	dirEntries, err := stCatalog.ReadDir("st")
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		data, err := stCatalog.ReadFile("st/" + dirEntry.Name())
		if err != nil {
			return err
		}

		var devices map[string]Catalog
		err = yaml.Unmarshal(data, &devices)
		if err != nil {
			return err
		}

		for id, device := range devices {
			if strings.HasPrefix(id, "x-") {
				continue // anchor-only bases
			}
			device.ID = id
			catalog[id] = device

			aliasIndex[normalize(id)] = id
			for _, alias := range device.Aliases {
				aliasIndex[normalize(alias)] = id
			}
		}
	}
	return nil
}

func Lookup(id string) (Catalog, error) {
	return Resolve(id)
}

func List(filter Filter) []Catalog {
	if err := ensureCatalog(); err != nil {
		return nil
	}

	vendor := normalize(filter.Vendor)
	family := normalize(filter.Family)
	series := normalize(filter.Series)

	result := make([]Catalog, 0, len(catalog))
	for _, device := range catalog {
		if vendor != "" && normalize(device.Vendor) != vendor {
			continue
		}
		if family != "" && normalize(device.Family) != family {
			continue
		}
		if series != "" && normalize(device.Series) != series {
			continue
		}
		result = append(result, device)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func Resolve(input string) (Catalog, error) {
	if err := ensureCatalog(); err != nil {
		return Catalog{}, err
	}
	key := normalize(input)
	if key == "" {
		return Catalog{}, fmt.Errorf("device is empty")
	}
	id, ok := aliasIndex[key]
	if !ok {
		return Catalog{}, fmt.Errorf("unknown device %q (not a catalog id or alias)", input)
	}
	return catalog[id], nil
}
