package windows

// func TestWindowsInventoryInit(t *testing.T) {
// 	if runtime.GOOS != "windows" {
// 		return
// 	}
// 	w := &WindowsInventory{}
// 	w.Setup()
// 	if len(w.TestItems) < 10 {
// 		t.Error("test item is few")
// 	}
// 	if w.TestItems[0].TestId != "system" {
// 		t.Error("invalid first test item : system")
// 	}
// }

// func TestWindowsInventoryCreateScript(t *testing.T) {
// 	if runtime.GOOS != "windows" {
// 		return
// 	}
// 	w := &WindowsInventory{}
// 	w.Setup()
// 	if err := w.CreateScript(); err != nil {
// 		t.Error(err)
// 	}
// 	if w.ScriptPath != "get_windows_inventory.ps1" {
// 		t.Error("not found windows inventory script")
// 	}
// }

// func TestWindowsInventoryRun(t *testing.T) {
// 	if runtime.GOOS != "windows" {
// 		return
// 	}
// 	logDir, _ := ioutil.TempDir("", "log")
// 	w := &WindowsInventory{}
// 	w.LogDir = logDir
// 	w.Setup()
// 	if err := w.Run(); err != nil {
// 		t.Error(err)
// 	}
// }

// func TestSetLogBackground(t *testing.T) {
// 	if err := SetLogBackground(); err != nil {
// 		t.Error("set log foreground", err)
// 	}
// 	log.Info("succeeded")
// 	log.Warn("not correct")
// 	log.Error("something error")
// }
