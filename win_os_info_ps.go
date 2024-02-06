package win_os_info

import "strings"
import "os/exec"
import "bytes"
import "log"
import "strconv"
import(
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//************* PowerShell ****************
type PowerShell struct {
	ps_path string
}

func NewPowerShell() *PowerShell {
	ps, err := exec.LookPath("powershell.exe")
	if err != nil {
		log.Fatal("Cannot find powershell.exe!")
	}
	return &PowerShell{ps_path: ps}
}

func (p *PowerShell) Execute(args ...string) (string, error) {
	args = append([]string{"-NonInteractive", "-NoProfile"}, args...)
	cmd := exec.Command(p.ps_path, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return stderr.String(), err
	}
	return stdout.String(), nil
}

var powershell = NewPowerShell()

//************** WinOsVersion ***************
type WinOsVersion struct {
	major uint64
	minor uint64
	build uint64
	caption string
}

func GetWinOsVersion() (WinOsVersion, error) {
	var win WinOsVersion
	var err error

	version, _ := powershell.Execute("(gwmi win32_operatingsystem).version")
	version_list := strings.Split(version, ".")

	for i, item := range(version_list) {
		version_list[i] = strings.TrimSpace(item)
	}

	if win.major, err = strconv.ParseUint(version_list[0], 10, 64); err != nil {
		return win, err
	}
	if win.minor, err = strconv.ParseUint(version_list[1], 10, 64); err != nil {
		return win, err
	}
	if win.build, err = strconv.ParseUint(version_list[2], 10, 64); err != nil {
		return win, err
	}
	if win.caption, err = powershell.Execute("(gwmi win32_operatingsystem).caption"); err != nil {
		return win, err
	}
	win.caption = strings.TrimSpace(win.caption)

	return win, nil
}

func (w *WinOsVersion) AddFields(acc telegraf.Accumulator,
				 measurement string,
				 tags map[string]string) {
	fields := make(map[string]interface{})
	fields["Major"] = w.major
	fields["Minor"] = w.minor
	fields["Build"] = w.build
	fields["Caption"] = w.caption

	acc.AddFields(measurement, fields, tags)
}


//************** WinOsResults *****************
type WinOsResults struct {
	LastUpdateDate string
	LastSearchDate string
}

func GetWinOsResults() (WinOsResults, error) {
	var res WinOsResults
	var err error

	if res.LastUpdateDate, err = powershell.Execute("(New-Object -ComObject " +
	"Microsoft.Update.AutoUpdate).Results.LastInstallationSuccessDate"); err != nil {
		return res, err
	}
	if res.LastSearchDate, err = powershell.Execute("(New-Object -ComObject" +
	" Microsoft.Update.AutoUpdate).Results.LastLastSearchDate"); err != nil {
		return res, err
	}

	res.LastUpdateDate = strings.TrimSpace(res.LastUpdateDate)
	res.LastSearchDate = strings.TrimSpace(res.LastSearchDate)
	return res, nil
}

func (w *WinOsResults) AddFields(acc telegraf.Accumulator,
				 measurement string,
				 tags map[string]string) {
	fields := make(map[string]interface{})
	fields["LastUpdateDate"] = w.LastUpdateDate
	fields["LastSearchDate"] = w.LastSearchDate

	acc.AddFields(measurement, fields, tags)
}
//*************** WinOsSettings ***************
type WinOsSettings struct {
	settings map[string]interface{}
}

func GetWinOsSettings() (WinOsSettings, error) {
	set, _ := powershell.Execute("(New-Object -ComObject Microsoft.Update.AutoUpdate).Settings")
	set = strings.TrimSpace(set)
	set_lines := strings.Split(set, "\n")
	settings := make(map[string]interface{})
	for _, line := range set_lines {
		pair := strings.Split(line, ":")
		pair[0] = strings.TrimSpace(pair[0])
		pair[1] = strings.TrimSpace(pair[1])
		switch pair[0] {
		case "NotificationLevel",
		     "ScheduledInstallationDay",
		     "ScheduledInstallationTime",
		     "IncludeRecommendedUpdates":
		     value, _ := strconv.Atoi(pair[1])
		     settings[pair[0]] = value
		case "ReadOnly", "Required",
		     "NonAdministratorElevated",
		     "FeatureUpdatesEnabled":
		     value, _ := strconv.ParseBool(pair[1])
		     settings[pair[0]] = value
		default:
			settings[pair[0]] = pair[1]
		}
	}
	return WinOsSettings{settings: settings}, nil
}

func (w *WinOsSettings) AddFields(acc telegraf.Accumulator,
				  measurement string,
				  tags map[string]string) {
	acc.AddFields(measurement, w.settings, tags)
}

//**************** WinOsInfo ****************
type WinOsInfo struct {
	OsVersion bool
	Results bool
	Settings bool
}

var WinOsInfoConfig = `
	## Plugin adds version and caption info
	OsVersion = true

	## Plugin adds update history information such as
	# LastInstallationSuccessDate
	# LastSearchSuccessDate
	Results = true

	## Plugin adds update settings such as:
	# NotificationLevel
	# ReadOnly
	# Required
	# ScheduledInstallationDay
	# ScheduledInstallationTime
	# IncludeRecommendedUpdates
	# NonAdministratorsElevated
	# FeatureUpdatesEnabled
	Settings = true 

	# Configure interval to collect once an hour
	interval = "20s"
`

func (w *WinOsInfo) SampleConfig() string {
	return WinOsInfoConfig
}

func (w *WinOsInfo) Description() string {
	return `##Sends information about Windows system:
		# Windows version and caption
		# Last successful update date
		# Last successful update search date
		# Windows update settings`
}

func (w *WinOsInfo) Gather(acc telegraf.Accumulator) error {
	tags := make(map[string]string)
	if w.OsVersion {
		version, _ := GetWinOsVersion()
		version.AddFields(acc, "win_os_version", tags)
	}
	if w.Results {
		results, _ := GetWinOsResults()
		results.AddFields(acc, "win_os_results", tags)
	}
	if w.Settings {
		settings, _ := GetWinOsSettings()
		settings.AddFields(acc, "win_os_settings", tags)
	}
	return nil
}

func init() {
	inputs.Add("win_os_info", func() telegraf.Input{ return &WinOsInfo{} })
}
