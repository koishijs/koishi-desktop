package manage

type KoiManager struct {
	exe string
}

func Manage(exe string) (manager *KoiManager, err error) {
	manager = &KoiManager{
		exe: exe,
	}
	return
}
