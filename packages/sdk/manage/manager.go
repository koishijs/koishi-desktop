package manage

type KoiManager struct {
	exe string
}

func NewKoiManager(exe string) (manager *KoiManager) {
	manager = &KoiManager{
		exe: exe,
	}
	return
}
