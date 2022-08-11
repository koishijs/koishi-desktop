package manage

type KoiManager struct {
	exe string
}

func Manage(exe string) (manager *KoiManager) {
	manager = &KoiManager{
		exe: exe,
	}
	return
}
