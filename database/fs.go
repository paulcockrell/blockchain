package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const databaseFolderName = "database"
const genesisFileName = "genesis.json"
const blockFileName = "block.db"

func InitDataDirIfNotExists(dataDir string, genesis []byte) error {
	if fileExist(getGenesisJSONFilePath(dataDir)) {
		return nil
	}

	if err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm); err != nil {
		return err
	}

	if err := writeGenesisToDisk(getGenesisJSONFilePath(dataDir), genesis); err != nil {
		return err
	}

	if err := writeEmptyBlocksDBToDisk(getBlocksDBFilePath(dataDir)); err != nil {
		return err
	}

	return nil
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, databaseFolderName)
}

func getGenesisJSONFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), genesisFileName)
}

func getBlocksDBFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), blockFileName)
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func writeEmptyBlocksDBToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}
