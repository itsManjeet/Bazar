package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func refreshData(refbtn *gtk.Button) {
	updcont := getWidget("headerContainer").(*gtk.Box)

	glist := updcont.GetChildren()
	glist.Foreach(func(item interface{}) {
		gtkWid := item.(*gtk.Widget)
		gtkWid.Destroy()
	})

	cmd := exec.Command("app", "refresh")
	s := true
	glib.IdleAdd(refProgress.SetVisible, true)
	go func() {
		for s {
			glib.IdleAdd(refProgress.Pulse, "refreshing database")
			time.Sleep(time.Second)
		}

	}()
	if err := cmd.Start(); err != nil {
		glib.IdleAdd(showError, err.Error())
		glib.IdleAdd(refbtn.SetSensitive, true)
		glib.IdleAdd(refProgress.SetVisible, false)
		s = false
		return
	}

	if err := cmd.Wait(); err != nil {
		glib.IdleAdd(showError, err.Error())
		glib.IdleAdd(refbtn.SetSensitive, true)
		glib.IdleAdd(refProgress.SetVisible, false)
		s = false
		return
	}

	cmdout, _ := exec.Command("app", "dry-check-update").Output()
	apd := make([]appData, 0)
	if len(cmdout) != 0 {
		a := strings.Split(string(cmdout), "\n")
		upda := ""
		for _, x := range a {
			if len(x) == 0 {
				continue
			}
			if x[0] == '\033' {
				continue
			}
			upda += " " + x
		}

		applist = listapps()

		for _, x := range strings.Split(upda, " ") {
			ap, err := getFromAppList(x)
			if err != nil {
				continue
			}
			apd = append(apd, ap)
		}
	}

	glib.IdleAdd(refProgress.SetVisible, false)
	glib.IdleAdd(refbtn.SetSensitive, true)
	s = false

	fmt.Println("Update applist", apd)

	if len(apd) > 0 {
		fmt.Println("update found", len(apd))
		updbtn, err := gtk.ButtonNewWithLabel(fmt.Sprintf("(%d) Update !", len(apd)))
		if err != nil {
			glib.IdleAdd(showError, err.Error())
			return
		}

		glib.IdleAdd(updcont.Add, updbtn)
		glib.IdleAdd(updbtn.Show)

		updbtn.Connect("clicked", onUpdateButtonClick, apd)
	}

	return

}

func onUpdateButtonClick(btn *gtk.Button, apd []appData) {
	go doUpdate(apd)
}
