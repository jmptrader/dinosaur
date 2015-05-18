package main

import (
	"fmt"

	"github.com/FcoManueel/Dinosaur/dino"
	ui "github.com/gizak/termui"
)

func main() {
	fmt.Println("Hello dinosaur! Enjoy your evolution. ")

	d := dino.New(20)

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	scheme := ui.Theme()
	scheme.BorderFg = ui.ColorMagenta
	scheme.BorderLabelTextFg = ui.ColorYellow

	//	bgColor := ui.ColorMagenta
	bgColor := ui.ColorCyan
	scheme.ParTextBg = bgColor
	scheme.BlockBg = bgColor
	scheme.BodyBg = bgColor
	scheme.BorderBg = bgColor
	scheme.BorderLabelTextBg = ui.ColorMagenta
	scheme.ListItemBg = ui.ColorWhite
	scheme.ListItemFg = ui.ColorBlack
	scheme.GaugeBar = ui.ColorGreen

	scheme.ParTextFg = ui.ColorMagenta
	scheme.ParTextBg = ui.ColorYellow
	//	scheme := ui.ColorScheme{
	//		BorderFg: ui.ColorCyan,
	//		ParTextFg: ui.ColorCyan,
	//		BarChartBar: ui.ColorRed,
	//		BorderLabelTextFg: ui.ColorMagenta,
	//		BorderLabelTextBg: ui.ColorCyan,
	//		BodyBg: ui.ColorWhite,
	//		BlockBg: ui.ColorWhite,
	//
	//	}
	//	ui.UseTheme("helloworld")
	ui.SetTheme(scheme)
	p := ui.NewPar("Welcome to dinosaur! A Operating System simulator written \nin Go, with memory management and process scheduling\n\n:Press Enter to evolve\t\t\t:Press q to quit")
	p.Height = 6
	p.Width = 60
	p.TextFgColor = ui.ColorMagenta
	p.TextBgColor = ui.ColorYellow
	p.Border.Label = "Dinosaur"

	halfWidth := 29

	cpuExec := ui.NewPar("")
	//cpuExec.Border.FgColor = scheme.BorderFg
	cpuExec.Width = halfWidth
	cpuExec.Height = 3
	cpuExec.Border.Label = "CPU"
	cpuExec.Y = 6

	ioExec := ui.NewPar("")
	//ioExec.Border.FgColor = scheme.BorderFg
	ioExec.Width = halfWidth
	ioExec.Height = 3
	ioExec.Border.Label = "IO"
	ioExec.X = 31
	ioExec.Y = 6

	strsNew := []string{}
	newProcs := ui.NewList()
	newProcs.Items = strsNew
	newProcs.Border.Label = "New"
	newProcs.Height = 12
	newProcs.Width = halfWidth
	newProcs.Y = 9

	strsReady := []string{""}
	readyProcs := ui.NewList()
	readyProcs.Items = strsReady
	readyProcs.Border.Label = "Ready"
	readyProcs.Height = 2 + len(strsReady)
	readyProcs.Width = halfWidth
	readyProcs.X = 2 + newProcs.Width
	readyProcs.Y = 9

	mem := ui.NewGauge()
	mem.Percent = 0
	mem.Width = 43
	mem.Height = 3
	mem.Y = 21
	mem.Border.Label = "Occupied Memory"
	mem.Border.LabelFgColor = scheme.BorderLabelTextFg

	frag := ui.NewPar("")
	//frag.Border.FgColor = scheme.BorderFg
	frag.Width = 15
	frag.Height = 3
	frag.Border.Label = "Fragmented"
	frag.X = 45
	frag.Y = 21
	frag.PaddingLeft = 4

	memLayout := ui.NewPar("")
	//memLayout.Border.FgColor = scheme.BorderFg
	memLayout.Width = 17
	memLayout.Height = 20
	memLayout.Border.Label = "Memory"
	memLayout.X = 60
	memLayout.Y = 0

	draw := func(state *dino.DinoState, d *dino.Dino) {
		mem.Percent = 100 - int(100*float32(state.FreeMemory)/float32(d.MemorySize()))
		newProcs.Items = state.NewQ
		readyProcs.Items = state.InteractiveQ
		readyProcs.Height = 2 + len(readyProcs.Items)
		if state.ExecutedByCPU != nil {
			cpuExec.PaddingLeft = 6
			cpuExec.Text = "Executed: " + state.ExecutedByCPU.Name
		} else {
			cpuExec.PaddingLeft = 8
			cpuExec.Text = "Not executed"
		}
		if state.ExecutedByIO != nil {
			ioExec.PaddingLeft = 6
			ioExec.Text = "Executed: " + state.ExecutedByIO.Name
		} else {
			ioExec.PaddingLeft = 8
			ioExec.Text = "Not executed"
		}
		if state.ExtFragmentation {
			ioExec.PaddingLeft = 1
			frag.Text = "Yes " + state.FragmentationProcess.Name
		} else {
			ioExec.PaddingLeft = 6
			frag.Text = "No!"
		}
		memString := ""
		for i, _ := range d.Memory {
			if i%10 == 0 {
				memString += "\n"
			}

			mark := "o"
			if d.Memory[i] != nil {
				mark = "X"
			} else {
				mark = "-"
			}
			memString += mark
		}

		memLayout.Text = memString
		ui.Render(p, newProcs, readyProcs, mem, cpuExec, ioExec, frag, memLayout)
	}

	evt := ui.EventCh()

	i := 0
	welcomeMessage := ui.NewPar("Dinosaur")
	welcomeMessage.HasBorder = false
	welcomeMessage.Width = 25
	welcomeMessage.Height = 3
	welcomeMessage.X = 35
	welcomeMessage.Y = 11
	ui.Render(welcomeMessage)

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			} else if e.Type == ui.EventKey && e.Key == ui.KeyEnter {
				state, err := d.Step()
				if err != nil {
					panic("Error while calculating step")
				}

				draw(state, d)
				i++
			}
		}
	}
	//d.Run(100)
	//  // I think that this way will make easier the communication with the front end
	//    i := 0
	//    for {
	//        state := d.Step()
	//        fmt.Printf(`\n----------------------------\n
	//                     Step %d: %+v\n`, i, state)
	//        SendToBrowser(state) // update clients
	//        i++
	//    }
}
