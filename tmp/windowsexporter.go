package windows

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TestItem struct {
	Level  int
	TestId string
	Script string
}

type WindowsExporter struct {
	LogDir     string
	ScriptPath string
	Level      int
	TestItems  []*TestItem
}

func NewTestItem(level int, testId, script string) *TestItem {
	testItem := &TestItem{
		Level:  level,
		TestId: testId,
		Script: script,
	}
	return testItem
}

func (w *WindowsExporter) Setup() {
	w.TestItems = []*TestItem{
		NewTestItem(0, "system", `
			Get-WmiObject -Class Win32_ComputerSystem
			`),
		NewTestItem(0, "os_conf", `
			Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion" | 
			Format-List
			`),
		NewTestItem(0, "os", `
			Get-WmiObject Win32_OperatingSystem | 
			Format-List Caption,CSDVersion,ProductType,OSArchitecture 
			`),
		NewTestItem(0, "driver", `
			Get-WmiObject Win32_PnPSignedDriver
			`),
		NewTestItem(0, "fips", `
			Get-Item "HKLM:System\CurrentControlSet\Control\Lsa\FIPSAlgorithmPolicy"
			`),
		NewTestItem(0, "virturalization", `
			Get-WmiObject -Class Win32_ComputerSystem | Select Model | FL
			`),
		NewTestItem(0, "storage_timeout", `
			Get-ItemProperty "HKLM:SYSTEM\CurrentControlSet\Services\disk" 
			`),
		NewTestItem(0, "monitor", `
			Get-WmiObject Win32_DesktopMonitor | FL
			`),
		NewTestItem(0, "ie_version", `
			Get-ItemProperty "HKLM:SOFTWARE\Microsoft\Internet Explorer"
			`),
		NewTestItem(2, "system_log", `
			Get-EventLog system | Where-Object { $_.EntryType -eq "Error" } | FL
			`),
		NewTestItem(2, "apps_log", `
			Get-EventLog application | Where-Object { $_.EntryType -eq "Error" } | FL
			`),
		NewTestItem(0, "ntp", `
			(Get-Item "HKLM:System\CurrentControlSet\Services\W32Time\Parameters").GetValue("NtpServer")
			`),
		NewTestItem(0, "user_account_control", `
			Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System"
			`),
		NewTestItem(0, "remote_desktop", `
			(Get-Item "HKLM:System\CurrentControlSet\Control\Terminal Server").GetValue("fDenyTSConnections")
			`),
		NewTestItem(0, "cpu", `
			Get-WmiObject -Class Win32_Processor | Format-List DeviceID, Name, MaxClockSpeed, SocketDesignation, NumberOfCores, NumberOfLogicalProcessors
			`),
		NewTestItem(0, "memory", `
			Get-WmiObject Win32_OperatingSystem |
			select TotalVirtualMemorySize,TotalVisibleMemorySize,
			FreePhysicalMemory,FreeVirtualMemory,FreeSpaceInPagingFiles
			`),
		NewTestItem(0, "dns", `
			Get-DnsClientServerAddress|FL
			`),
		NewTestItem(0, "nic_teaming_config", `
			Get-NetLbfoTeamNic
			`),
		NewTestItem(0, "tcp", `
			Get-ItemProperty "HKLM:SYSTEM\CurrentControlSet\Services\Tcpip\Parameters" |
			Format-List
			`),
		NewTestItem(0, "network", `
			Get-WmiObject Win32_NetworkAdapterConfiguration |
			Where{$_.IpEnabled -Match "True"} |
			Select ServiceName, MacAddress, IPAddress, DefaultIPGateway, Description, IPSubnet |
			Format-List
			`),
		NewTestItem(0, "network_profile", `
			Get-NetConnectionProfile
			`),
		NewTestItem(0, "net_bind", `
			Get-NetAdapterBinding | FL
			`),
		NewTestItem(0, "net_ip", `
			Get-NetIPInterface | FL
			`),
		NewTestItem(0, "firewall", `
			Get-NetFirewallRule -Direction Inbound -Enabled True
			`),
		NewTestItem(0, "filesystem", `
			Get-WmiObject Win32_LogicalDisk | Format-List *
			`),
		NewTestItem(2, "user", `
			$result = @()
			$accountObjList =  Get-CimInstance -ClassName Win32_Account
			$userObjList = Get-CimInstance -ClassName Win32_UserAccount
			foreach($userObj in $userObjList)
			{
			    $IsLocalAccount = ($userObjList | ?{$_.SID -eq $userObj.SID}).LocalAccount
			    if($IsLocalAccount)
			    {
			        $query = "WinNT://{0}/{1},user" -F $env:COMPUTERNAME,$userObj.Name
			        $dirObj = New-Object -TypeName System.DirectoryServices.DirectoryEntry -ArgumentList $query
			        $UserFlags = $dirObj.InvokeGet("UserFlags")
			        $DontExpirePasswd = [boolean]($UserFlags -band 0x10000)
			        $AccountDisable   = [boolean]($UserFlags -band 0x2)
			        $obj = New-Object -TypeName PsObject
			        Add-Member -InputObject $obj -MemberType NoteProperty -Name "UserName" -Value $userObj.Name
			        Add-Member -InputObject $obj -MemberType NoteProperty -Name "DontExpirePasswd" -Value $DontExpirePasswd
			        Add-Member -InputObject $obj -MemberType NoteProperty -Name "AccountDisable" -Value $AccountDisable
			        Add-Member -InputObject $obj -MemberType NoteProperty -Name "SID" -Value $userObj.SID
			        $result += $obj
			    }
			}
			$result | Format-List
			`),
		NewTestItem(0, "etc_hosts", `
			Get-Content "$($env:windir)\system32\Drivers\etc\hosts"
			`),
		NewTestItem(0, "patch_lists", `
			wmic qfe
			`),
		NewTestItem(1, "service", `
			Get-Service | FL
			`),
		NewTestItem(1, "packages", `
			Get-WmiObject Win32_Product |
			Select-Object Name, Vendor, Version |
			Format-List
			Get-ChildItem -Path(
			'HKLM:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall',
			'HKCU:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall') |
			% { Get-ItemProperty $_.PsPath | Select-Object DisplayName, Publisher, DisplayVersion } |
			Format-List
			`),
		NewTestItem(1, "feature", `
			Get-WindowsFeature | ?{$_.InstallState -eq [Microsoft.Windows.ServerManager.Commands.InstallState]::Installed} | FL
			`),
		NewTestItem(2, "task_scheduler", `
			Get-ScheduledTask |
			? {$_.State -eq "Ready"} |
			Get-ScheduledTaskInfo |
			? {$_.NextRunTime -ne $null}|
			Format-List
			`),
	}
}

func (w *WindowsExporter) CreateScript() error {
	log.Info("create temporary log dir for test ", w.LogDir)
	tmpl, err := template.ParseFiles("template/ps1.tpl")
	if err != nil {
		return errors.Wrap(err, "failed read template")
	}
	outPath := filepath.Join(w.LogDir, "get_windows_inventory.ps1")
	outFile, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed create script")
	}
	defer outFile.Close()
	var testItems []*TestItem
	for _, testItem := range w.TestItems {
		if testItem.Level <= w.Level {
			log.Info("add test item ", testItem.TestId)
			testItems = append(testItems, testItem)
		}
	}
	if err := tmpl.Execute(outFile, testItems); err != nil {
		return errors.Wrap(err, "failed generate script")
	}
	log.Info("created windows inventory script ", outPath)
	w.ScriptPath = outPath
	return nil
}

func (w *WindowsExporter) Run() error {
	if err := w.CreateScript(); err != nil {
		return errors.Wrap(err, "run windows inventory")
	}
	if runtime.GOOS != "windows" {
		return fmt.Errorf("windows powershell environment only")
	}
	cmdPowershell := []string{
		"powershell",
		w.ScriptPath,
		w.LogDir,
	}
	log.Info(cmdPowershell)
	cmd := exec.Command(cmdPowershell[0], cmdPowershell[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed create windows inventory command pipe")
	}
	startTime := time.Now()
	testId := ""
	intervalTime := startTime
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if testId != "" {
			log.Info(testId, ", Elapse: ", time.Since(intervalTime))
		}
		testId = scanner.Text()
		intervalTime = time.Now()
	}
	if testId != "" {
		log.Info(testId, ", Elapse: ", time.Since(intervalTime))
	}
	cmd.Wait()
	// tio := &timeout.Timeout{
	// 	Cmd:       cmd,
	// 	Duration:  defaultTimeoutDuration,
	// 	KillAfter: timeoutKillAfter,
	// }
	// exitstatus, stdout, stderr, err := tio.Run()

	log.Infof("finish windows inventory script, elapse [%s]", time.Since(startTime))
	// if err != nil {
	// 	return fmt.Errorf("test3 %s", err)
	// }
	// enc := japanese.ShiftJIS
	// bb, _, err := transform.Bytes(enc.NewDecoder(), []byte(stderr))

	// log.Info("RC:", exitstatus)
	// log.Info("STDOUT:", stdout)
	// log.Info("STDERR:", string(bb))

	return nil
}
