package gotil

func GetHomeDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

func EtcHostsFilename() string {
	return "C:\\windows\\system32\\drivers\\etc\\hosts"
}
