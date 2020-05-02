package agent

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *Config) InitHome() error {
	createDirs := []string{
		c.OutDir, c.WorkDir, c.NodeDir, c.WorkCommonDir, c.ArchiveDir, c.SslDir,
	}
	for _, createDir := range createDirs {
		if err := os.MkdirAll(createDir, 0777); err != nil {
			return fmt.Errorf("initialize director for agent : %s", err)
		}
	}
	return nil
}

func (c *Config) InitAgent() error {
	if err := c.InitHome(); err != nil {
		return errors.Wrap(err, "initialize home")
	}
	log.Info("write pid ", c.PidFile)
	c.WriteWorkFileNumber(c.PidFile, os.Getpid())
	return nil
}

func (config *Config) CheckExitFile() (string, error) {
	status := ""
	lines, _ := config.ReadWorkFileHead(config.ExitFlag, 1)
	if len(lines) > 0 {
		status = lines[0]
	}
	return status, nil
}

func checkNoSpace(value string) bool {
	for _, v := range value {
		if unicode.IsSpace(v) || v == '/' || v == '\\' {
			return false
		}
	}
	return true
}

func (config *Config) CheckHostname(hostname string) bool {
	c := hostname[0]
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return checkNoSpace(hostname)
	}
	return false
}

func (config *Config) CheckHAStatus() error {
	schedule := config.Schedule
	if schedule == nil {
		return fmt.Errorf("schedule not found for check ha status")
	}
	if schedule.HanodeEnable == false {
		return nil
	}
	log.Info("[0] HA service Check =========================")
	log.Info(schedule.HanodeCmd)
	CmdInfo := &CommandInfo{CmdLine: schedule.HanodeCmd, Timeout: 60}
	CmdInfo.ExecCommandNoRedirect()
	if CmdInfo.Executed == false {
		log.Error("check ha status command not execute")
		return nil
	}
	serviceName := strings.TrimFunc(CmdInfo.OutBuf, unicode.IsSpace)
	if config.CheckHostname(serviceName) == false {
		log.Error("invalid check ha servicename ", serviceName)
		return nil
	}
	log.Info("service name : ", serviceName)
	config.ServiceName = serviceName
	return nil
}

func (config *Config) AuthLicense(expire int) error {
	schedule := config.Schedule
	log.Debug("AuthLicense")
	if schedule == nil {
		return fmt.Errorf("schedule not found for auth license")
	}
	license := schedule.License
	if license == nil {
		return fmt.Errorf("license not found for auth license")
	}
	log.Debug("HOST CHECK", strings.Compare(license.Hostname, config.Host))
	if strings.Compare(license.Hostname, config.Host) != 0 {
		return fmt.Errorf("invalid license host : [%s, %s]", license.Hostname, config.Host)
	}
	currentDate := GetCurrentTime(expire, YYYYMMDD)
	oldFlag := strings.Compare(license.Expired, currentDate)
	log.Debugf("AuthLicense %v", oldFlag)
	log.Debugf("AuthLicense %v %v", currentDate, license.Expired)
	if strings.Compare(license.Expired, currentDate) == -1 {
		return fmt.Errorf("invalid license expired : %s < %s",
			currentDate, license.Expired)
	}
	// オリジナルは MD5 ライブラリで 40桁(320bit)のハッシュを返すが、
	// Go ライブラリの場合、32桁(256bit)のハッシュとなり、ロジックが
	// 異なる。実装は保留し、ハッシュチェックはスキップする
	h := md5.New()
	io.WriteString(h, config.Host)
	io.WriteString(h, license.Expired)
	io.WriteString(h, schedule.SiteKey)
	// checksum := h.Sum(nil)
	//   if license.Code != fmt.Sprintf("%x", checksum) {
	// return fmt.Errorf("invalid license code : %s", license.Code)
	//   }
	return nil
}

func (config *Config) UnzipSSLConf() error {
	// sslPath := config.GetWorkfilePath("sslconf.zip")
	sslPath := filepath.Join(config.WorkCommonDir, "sslconf.zip")
	log.Info("unzip ", sslPath)
	return Unzip(sslPath, config.Home)
}

func (config *Config) PostTask(task *Task) {
	ds := task.Datastore
	id := ds.Path()
	// elapse := task.EndTime.Sub(task.StartTime)
	log.Info("end task [", id, "]")
	if err := config.SaveReport(task, ds); err != nil {
		log.Error("report failed ", err)
	}
	log.Debug("check send enable ", config.Schedule.RemhostEnable)
	if config.Schedule.RemhostEnable == true {
		if err := config.ArchiveData(task, ds); err != nil {
			log.Error("zip failed ", err)
			return
		}
		if err := config.SendCollectorDataAll(ds); err != nil {
			log.Error("send failed ", err)
			return
		}
	}
	if err := config.PurgeData(task, ds); err != nil {
		log.Error("purge failed ", err)
		return
	}
}

func (config *Config) SaveReport(task *Task, ds *Datastore) error {
	// yaml, err := task.MakeReport()
	// fileName := fmt.Sprintf("stat_%s.log", task.Job.StaJob)
	// log.Debug("make report ", ds.Path())
	// filePath := filepath.Join(ds.AbsDir(), fileName)
	// file, err := os.Create(filePath)
	// if err != nil {
	// 	return fmt.Errorf("write report file %s: %s", filePath, err)
	// }
	// defer file.Close()
	// file.Write(([]byte)(yaml))
	return nil
}

// ArchiveDataは、指定したタスクのログに対して zip 圧縮をします。
// * arc_{host}__{stat}_{date}_{time}.zip ファイル名作成を作成します。
// * {stat}/{date}/{time} ディレクトリパスを作成します。
// * ログ保存ディレクトリ(ODir)に移動します。
// * 次のコマンドと同じ zip 圧縮。zip -r arcxxx.zip ODIR
func (config *Config) ArchiveData(task *Task, ds *Datastore) error {
	if ds == nil {
		return fmt.Errorf("output log is nil")
	}
	zipFile := ds.ZipFile(config.GetServiceOrHostName())
	zipPath := config.GetArchivefilePath(zipFile)
	log.Info("archive ", zipFile)
	return Zip(zipPath, config.OutDir, ds.Path())
}

func (config *Config) SendCollectorData(zip string) error {
	cmd := filepath.Join(config.BinDir, "getperfsoap")
	confFile := config.ParameterFile
	cmdLine := fmt.Sprintf("%s --send -c %s %s", cmd, confFile, zip)
	log.Info("send   ", zip)
	log.Debug("COMMAND ", cmdLine)
	CmdInfo := &CommandInfo{
		CmdLine: cmdLine,
		Timeout: 15,
	}
	startTime := time.Now()
	err := CmdInfo.ExecCommandNoRedirect()
	if err != nil {
		log.Errorf("send failed %s [%s] : error : %s, command : %s", zip, time.Since(startTime), err, cmdLine)
		return err
	} else {
		log.Infof("sended %s [%s]", zip, time.Since(startTime))
		return config.RemoveArchiveFile(zip)
	}
}

func (config *Config) SendCollectorDataAll(ds *Datastore) error {
	recoveryHour := config.Schedule.RecoveryHour
	hostName := config.GetServiceOrHostName()
	oldZip := ds.OldZipFile(hostName, recoveryHour)
	zipPrefix := fmt.Sprintf("arc_%s__%s_", hostName, ds.StatName)
	zips, err := ioutil.ReadDir(config.ArchiveDir)
	if err != nil {
		return err
	}
	for _, zip := range zips {
		zipFile := zip.Name()
		if !strings.HasPrefix(zipFile, zipPrefix) || !strings.HasSuffix(zipFile, ".zip") {
			continue
		}
		oldFlag := strings.Compare(zipFile, oldZip)
		if oldFlag < 0 {
			config.RemoveArchiveFile(zipFile)
		} else {
			if err := config.SendCollectorData(zipFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func (config *Config) PurgeData(task *Task, ds *Datastore) error {
	saveHour := config.Schedule.SaveHour
	logDir := config.OutDir
	purgeTime := ds.StartTime.Add(-1 * time.Hour * time.Duration(saveHour))
	dateDirs, err := ioutil.ReadDir(logDir)
	if err != nil {
		return err
	}
	checkDate := purgeTime.Format("20060102")
	for _, dateDir := range dateDirs {
		if dateDir.Name() < checkDate {
			purgeDir := filepath.Join(logDir, dateDir.Name())
			fmt.Printf("remove %s\n", purgeDir)
			if err := os.RemoveAll(purgeDir); err != nil {
				return err
			}
		}
	}
	checkTime := purgeTime.Format("150405")
	checkDatePath := filepath.Join(logDir, checkDate)
	timeDirs, err := ioutil.ReadDir(checkDatePath)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	for _, timeDir := range timeDirs {
		if timeDir.Name() < checkTime {
			purgeDir := filepath.Join(checkDatePath, timeDir.Name())
			fmt.Printf("remove %s\n", purgeDir)
			if err := os.RemoveAll(purgeDir); err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckLicense はライセンスファイルをダウンロードし、ライセンスをチェックします
func (config *Config) CheckLicense(expired int) error {
	// retry := DEFAULT_SOAP_RETRY
	// authOk := false

	// for authOk == false && retry > 0 {
	// 	if err := config.AuthLicense(expired); err == nil {
	// 		authOk = true
	// 		break
	// 	} else {
	// 		log.Info("auth license error : ", err)
	// 	}
	// 	log.Info("auth error resync License.txt ", retry)
	// 	if err := config.DownloadLicense(); err != nil {
	// 		log.Error("remhost license web service failed ", err)
	// 	} else if err := config.UnzipSSLConf(); err != nil {
	// 		log.Error("unzip sslconf.zip failed ", err)
	// 	} else {
	// 		if err := config.LoadLicense(); err != nil {
	// 			log.Error("can't load ssl license ", err)
	// 		} else if err := config.AuthLicense(expired); err == nil {
	// 			authOk = true
	// 			continue
	// 		}
	// 	}
	// 	if retry != DEFAULT_SOAP_RETRY {
	// 		time.Sleep(time.Second * CHECK_LICENSE_INTERVAL)
	// 	}
	// 	retry--
	// }
	// if authOk {
	// 	return nil
	// } else {
	// 	return fmt.Errorf("license check failed")
	// }
	return nil
}

func (config *Config) DownloadLicense() error {
	cmd := filepath.Join(config.BinDir, "getperfsoap")
	confFile := config.ParameterFile
	cmdLine := fmt.Sprintf("%s --get -c %s %s", cmd, confFile, "sslconf.zip")
	log.Info("Get ", cmdLine)
	CmdInfo := &CommandInfo{
		CmdLine: cmdLine,
		Timeout: config.Schedule.SoapTimeout,
	}
	if err := CmdInfo.ExecCommandNoRedirect(); err != nil {
		return fmt.Errorf("failed to get sslconf.zip %s", err)
	}
	if CmdInfo.ExitCode != 0 {
		return fmt.Errorf("failed to get sslconf.zip")
	}
	return nil
}

func (config *Config) CheckDiskUtil() error {
	// ExecSOAPCommandPM("--get", "sslconf.zip")
	disk, err := CheckDiskFree(config.OutDir)
	if err != nil {
		return err
	}
	diskCapacity := config.Schedule.DiskCapacity
	diskUtil := int(100.0 * disk.Free / disk.All)
	if diskCapacity > diskUtil {
		return fmt.Errorf("disk usage %d > %d", diskCapacity, diskUtil)
	}
	return nil
}

func (config *Config) Stop() error {
	log.Info("stop process")
	config.RemoveWorkFile(config.ExitFlag)
	config.RemoveWorkFile("_running_flg")
	os.RemoveAll(config.WorkDir)
	os.Exit(0)
	return nil
}

func (config *Config) PrepareTask(job *Job) error {
	log.Debug("[S] CHECK ==================")
	if err := config.CheckHAStatus(); err != nil {
		return fmt.Errorf("HA service ... NG : %s", err)
	}
	if config.Schedule.RemhostEnable == true && job.Id == 1 {
		if err := config.CheckLicense(0); err != nil {
			return fmt.Errorf("License ... NG : %s", err)
		}
	}
	if err := config.CheckDiskUtil(); err != nil {
		return fmt.Errorf("Check Diskfree ... NG : %s", err)
	}
	return nil
}

func (config *Config) Run() error {
	ctx, cancel := MakeContext(0)
	defer cancel()
	return config.RunWithContext(ctx)
}

func (config *Config) RunCollector(job *Job, count int, tasks chan<- *Task) {
	timeout := time.Duration(job.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan bool, 1)
	ds, err := config.NewDatastore(job.Name, time.Now(), TIME_ON)
	if err != nil {
		log.Error("new out log ", err)
	}
	task := NewTask(job, ds.PathAbs())
	task.Datastore = ds
	id := ds.Path()
	log.Info("run task [", id, "]")
	log.Debugf("run task COLLECTOR %d %s", job.Id, job.Name)
	go func(ctx context.Context) {
		// ctx を第一引数にして、時間のかかる処理を実行。
		// time.Sleep(time.Second * time.Duration(job.Duration))
		if err := config.PrepareTask(job); err != nil {
			log.Error("prepare task ", err)
		} else if err := task.RunWithContext(ctx); err != nil {
			log.Error("task ", err)
		}
		ch <- true
	}(ctx)
	select {
	case <-ch:
		log.Debug("task end ", id)
		task.Status = END
	case <-ctx.Done():
		log.Println(id, ctx.Err())
		task.Status = TIMEOUT
		cancel()
	}
	task.EndTime = time.Now()
	tasks <- task
}

func (config *Config) RunWithContext(ctx context.Context) error {
	config.InitAgent()
	log.Info("run agent with conext")
	if config.Schedule.RemhostEnable == true {
		// config.LoadLicense()
		// if err := config.CheckLicense(0); err != nil {
		// 	return fmt.Errorf("License ... NG : %s", err)
		// }
	}
	config.Mode = RUN
	config.WriteLineWorkFile("_running_flg", "")
	jobs := config.Schedule.Jobs

	config.FileServe()
	tasks := make(chan *Task, len(jobs)*2)
	for statName, job := range jobs {
		job.Name = statName
		if job.Enable == false {
			continue
		}
		go func(job *Job) {
			count := 0
			log.Info("run task ", job)
			go config.RunCollector(job, count, tasks)
			interval := time.Second * time.Duration(job.Interval)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			count++
			for {
				select {
				case <-ticker.C:
					log.Info("run task ", job)
					// go config.RunJob(c, count, tasks)
					count++
				}
			}
		}(job)
	}
	log.Info("start process monitor")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()
	count := 0

	for {
		select {
		case task := <-tasks:
			config.PostTask(task)
		case <-ticker.C:
			status, err := config.CheckExitFile()
			if err != nil {
				log.Error("check exit file ", err)
			} else if status == "STOP" {
				config.Stop()
			}
			// log.Println("CHECK stop flag file", count)
			count++
		case <-quit:
			log.Debug("CATCH SIGNAL")
			config.Stop()
			os.Exit(0)
		}
	}
}
