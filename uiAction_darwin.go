package qframelesswindow

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func (f *QFramelessWindow) SetTitleBarActions() {
	t := f.TitleBar

	// TitleBar Actions
	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Raise()
		f.IsTitleBarPressed = true
		f.TitleBarMousePos = e.GlobalPos()
		f.Pos = f.Window.Pos()
	})

	t.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	t.ConnectMouseMoveEvent(func(e *gui.QMouseEvent) {
		if !f.IsTitleBarPressed {
			return
		}
		x := f.Pos.X() + e.GlobalPos().X() - f.TitleBarMousePos.X()
		y := f.Pos.Y() + e.GlobalPos().Y() - f.TitleBarMousePos.Y()
		newPos := core.NewQPoint2(x, y)
		f.Window.Move(newPos)
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		if f.BtnMaximize.IsVisible() {
			f.windowMaximize()
		} else {
			f.windowRestore()
		}
	})

	// Button Actions
	f.BtnMinimize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnMaximize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnRestore.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnClose.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnMinimize.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.Window.SetWindowState(core.Qt__WindowMinimized)
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnMaximize.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.windowMaximize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnRestore.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.windowRestore()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnClose.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.Window.Close()
	})
}

func (f *QFramelessWindow) windowMaximize() {
	f.BtnMaximize.SetVisible(false)
	f.BtnRestore.SetVisible(true)
	f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.Window.SetWindowState(core.Qt__WindowMaximized)
}

func (f *QFramelessWindow) windowRestore() {
	f.BtnMaximize.SetVisible(true)
	f.BtnRestore.SetVisible(false)
	f.Layout.SetContentsMargins(f.shadowMargin, f.shadowMargin, f.shadowMargin, f.shadowMargin)
	f.Window.SetWindowState(core.Qt__WindowNoState)
}
