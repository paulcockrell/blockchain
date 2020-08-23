package wallet

import "path/filepath"

const keystoreDirName = "keystore"
const PaulcAccount = "0x2858e7edea0FeD8d3255f1b5b4d9987568095062"
const DavecAccount = "0x21973d33e048f5ce006fd7b41f51725c30e4b76b"
const KimcAccount = "0x84470a31D271ea400f34e7A697F36bE0e866a716"

func GetKeystoreDirPath(dataDir string) string {
	return filepath.Join(dataDir, keystoreDirName)
}
