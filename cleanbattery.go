package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"bytes"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/glib"
	"os"
)

const (
	TIMEOUT = 3000
	NEWLINE = byte(10)
	STATUS_CHARGING = "Charging"
	STATUS_DISCHARGING = "Discharging"
)

var (
	status_icon *gtk.GtkStatusIcon
	full int64
)

func main() {
	gtk.Init(&os.Args)
	glib.SetApplicationName("zzcleanbattery")

	buf, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/energy_full")
	if err != nil { panic(err) }
	str := string(bytes.Split(buf, []byte{NEWLINE})[0])
	full, err = strconv.ParseInt(str, 10, 64)
	if err != nil { panic(err) }

	status_icon = gtk.StatusIcon()
	status_icon.SetTitle("zzcleanbattery")

	update_icon()

	glib.TimeoutAdd(TIMEOUT, update_icon)
	gtk.Main()
}

func update_icon() bool {

	var (
		hours int64
		minutes int64
		seconds int64
		pfull int64
		rate int64
		now int64
		status string
		text string
	)

	buf, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/energy_now")
	if err != nil { panic(err) }
	str := string(bytes.Split(buf, []byte{NEWLINE})[0])
	now, err = strconv.ParseInt(str, 10, 64)
	if err != nil { panic(err) }

	buf, err = ioutil.ReadFile("/sys/class/power_supply/BAT0/power_now")
	if err != nil { panic(err) }
	str = string(bytes.Split(buf, []byte{NEWLINE})[0])
	rate, err = strconv.ParseInt(str, 10, 64)
	if err != nil { panic(err) }

	buf, err = ioutil.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil { panic(err) }
	status = string(bytes.Split(buf, []byte{NEWLINE})[0])

	pfull = now * 100 / full

	if rate > 0 {
		switch status {
		case STATUS_CHARGING:
			seconds = 3600 * (full - now) / rate
		case STATUS_DISCHARGING:
			seconds = 3600 * now / rate
		default:
			seconds = 0
		}
	} else {
		seconds = 0
	}
	hours = seconds / 3600
	seconds -= hours * 3600
	minutes = seconds / 60
	seconds -= minutes * 60
	if seconds == 0 {
		text = fmt.Sprintf("%s, %d%%\n", status, pfull)
	} else {
		text = fmt.Sprintf("%s, %d%%, %d:%d remaining\n",
		           status,
		           pfull,
		           hours,
		           minutes)
	}

	status_icon.SetTooltipText(text)
	status_icon.SetFromIconName(get_icon_name(status, pfull))
	return true
}

func get_icon_name(status string, pfull int64) string {
	if status == STATUS_DISCHARGING {
		switch {
		case pfull < 10:
			return "battery_empty"
		case pfull < 20:
			return "battery_caution"
		case pfull < 40:
			return "battery_low"
		case pfull < 60:
			return "battery_two_thirds"
		case pfull < 75:
			return "battery_third_fourth"
		default:
			return "battery_full"
		}
	} else if status == STATUS_CHARGING {
		return "battery_charged"
	}
	return "battery_plugged"
}