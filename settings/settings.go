package settings

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	settingsInstance *AppSettings
)

const (
	SETTINGS_FILENAME      = "settings.json"
	TITLE_JSON_FILENAME    = "titles.json"
	VERSIONS_JSON_FILENAME = "versions.json"
	SLM_VERSION            = "1.4.0"
	SLM_WEB_VERSION        = "1.0.12"
	TITLES_JSON_URL        = "https://tinfoil.io/repo/db/titles.json"
	//VERSIONS_JSON_URL    = "https://tinfoil.io/repo/db/versions.json"
	VERSIONS_JSON_URL      = "https://raw.githubusercontent.com/blawar/titledb/master/versions.json"
)

const (
	TEMPLATE_TITLE_ID    = "TITLE_ID"
	TEMPLATE_TITLE_NAME  = "TITLE_NAME"
	TEMPLATE_DLC_NAME    = "DLC_NAME"
	TEMPLATE_VERSION     = "VERSION"
	TEMPLATE_REGION      = "REGION"
	TEMPLATE_VERSION_TXT = "VERSION_TXT"
	TEMPLATE_TYPE        = "TYPE"
)

type OrganizeOptions struct {
	CreateFolderPerGame  bool   `json:"create_folder_per_game"`
	RenameFiles          bool   `json:"rename_files"`
	DeleteEmptyFolders   bool   `json:"delete_empty_folders"`
	DeleteOldUpdateFiles bool   `json:"delete_old_update_files"`
	FolderNameTemplate   string `json:"folder_name_template"`
	SwitchSafeFileNames  bool   `json:"switch_safe_file_names"`
	FileNameTemplate     string `json:"file_name_template"`
}

type AppSettings struct {
	VersionsEtag           string          `json:"versions_etag"`
	TitlesEtag             string          `json:"titles_etag"`
	Prodkeys               string          `json:"prod_keys"`
	Folder                 string          `json:"folder"`
	ScanFolders            []string        `json:"scan_folders"`
	Port                   int             `json:"port"`
	Debug                  bool            `json:"debug"`
	OrganizeOptions        OrganizeOptions `json:"organize_options"`
	IgnoreDLCTitleIds      []string        `json:"ignore_dlc_title_ids"`
}

func ReadSettingsAsJSON(dataFolder string) string {
	if _, err := os.Stat(filepath.Join(dataFolder, SETTINGS_FILENAME)); err != nil {
		saveDefaultSettings(dataFolder)
	}
	file, _ := os.Open(filepath.Join(dataFolder, SETTINGS_FILENAME))
	bytes, _ := ioutil.ReadAll(file)
	return string(bytes)
}

func ReadSettings(dataFolder string) *AppSettings {
	if settingsInstance != nil {
		return settingsInstance
	}
	settingsInstance = &AppSettings{Debug: false, ScanFolders: []string{},
		OrganizeOptions: OrganizeOptions{SwitchSafeFileNames: true}, Prodkeys: "", IgnoreDLCTitleIds: []string{"01007F600B135007"}}
	if _, err := os.Stat(filepath.Join(dataFolder, SETTINGS_FILENAME)); err == nil {
		file, err := os.Open(filepath.Join(dataFolder, SETTINGS_FILENAME))
		if err != nil {
			zap.S().Warnf("Missing or corrupted config file, creating a new one")
			return saveDefaultSettings(dataFolder)
		} else {
			_ = json.NewDecoder(file).Decode(&settingsInstance)
			return settingsInstance
		}
	} else {
		return saveDefaultSettings(dataFolder)
	}
}

func saveDefaultSettings(dataFolder string) *AppSettings {
	settingsInstance = &AppSettings{
		TitlesEtag:             "W/\"a5b02845cf6bd61:0\"",
		VersionsEtag:           "W/\"2ef50d1cb6bd61:0\"",
		Prodkeys:               dataFolder,
		Folder:                 "/mnt/roms",
		ScanFolders:            []string{},
		IgnoreDLCTitleIds:      []string{},
		Port:                   3000,
		Debug:                  false,
		OrganizeOptions: OrganizeOptions{
			RenameFiles:         false,
			CreateFolderPerGame: false,
			FolderNameTemplate:  fmt.Sprintf("{%v}", TEMPLATE_TITLE_NAME),
			FileNameTemplate: fmt.Sprintf("{%v} ({%v})[{%v}][v{%v}]", TEMPLATE_TITLE_NAME, TEMPLATE_DLC_NAME,
				TEMPLATE_TITLE_ID, TEMPLATE_VERSION),
			DeleteEmptyFolders:   false,
			SwitchSafeFileNames:  true,
			DeleteOldUpdateFiles: false,
		},
	}
	return SaveSettings(settingsInstance, dataFolder)
}

func SaveSettings(settings *AppSettings, dataFolder string) *AppSettings {
	file, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile(filepath.Join(dataFolder, SETTINGS_FILENAME), file, 0644)
	settingsInstance = settings
	return settings
}
