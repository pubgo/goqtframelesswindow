package qframelesswindow

import (
	"fmt"
	"runtime"

        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/svg"
        "github.com/therecipe/qt/widgets"
)

type Edge int

const (
        None Edge = 0x0
        Left Edge = 0x1
        Top Edge = 0x2
        Right Edge = 0x4
        Bottom Edge = 0x8
        TopLeft Edge = 0x10
        TopRight Edge = 0x20
        BottomLeft Edge = 0x40
        BottomRight Edge = 0x80
)

type QFramelessWindow struct {
	Window         *widgets.QMainWindow
	Widget         *widgets.QWidget

	windowColor    string

	borderSize     int
	Layout         *widgets.QVBoxLayout
	// Layout         *widgets.QGridLayout

	WindowWidget   *widgets.QFrame
	// WindowWidget   *widgets.QWidget
	WindowVLayout  *widgets.QVBoxLayout

	TitleBar       *widgets.QWidget
	TitleBarLayout *widgets.QHBoxLayout
	TitleLabel     *widgets.QLabel
	TitleBarBtnWidget *widgets.QWidget
	TitleBarBtnLayout *widgets.QHBoxLayout

	// for darwin
	BtnMinimize    *widgets.QToolButton
	BtnMaximize    *widgets.QToolButton
	BtnRestore     *widgets.QToolButton
	BtnClose       *widgets.QToolButton

	// for windows, linux
	IconMinimize    *svg.QSvgWidget
	IconMaximize    *svg.QSvgWidget
	IconRestore     *svg.QSvgWidget
	IconClose       *svg.QSvgWidget

	isCursorChanged bool
	isDragStart     bool
	dragPos         *core.QPoint
	pressedEdge     Edge

	Content        *widgets.QWidget

	Pos            *core.QPoint
	MousePos       *core.QPoint
	IsMousePressed bool
}

func NewQFramelessWindow() *QFramelessWindow {
	f := &QFramelessWindow{}
	f.Window = widgets.NewQMainWindow(nil, 0)
	f.Widget = widgets.NewQWidget(nil, 0)
	f.SetborderSize(4)
	f.Window.SetCentralWidget(f.Widget)
	f.SetupUI(f.Widget)
	f.SetWindowFlags()
	f.SetAttribute()
	f.SetWindowActions()
        f.SetTitleBarActions()

	return f
}

func (f *QFramelessWindow) SetborderSize(size int) {
	f.borderSize = size
}

func (f *QFramelessWindow) SetupUI(widget *widgets.QWidget) {
        //f.Layout = widgets.NewQVBoxLayout2(widget)

	widget.SetObjectName("QFramelessWindow")
	window := widget.Window()
	window.InstallEventFilter(window)

        // f.Layout = widgets.NewQGridLayout(widget)
        f.Layout = widgets.NewQVBoxLayout2(widget)
        f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.Layout.SetSpacing(0)

        f.WindowWidget = widgets.NewQFrame(widget, 0)
	// f.WindowWidget.InstallEventFilter(f.WindowWidget)
	// f.WindowWidget.InstallEventFilter(window)
	// window.InstallEventFilter(f.WindowWidget)
        // f.WindowWidget = widgets.NewQWidget(widget, 0)

        //f.WindowWidget.SetObjectName("QFramelessWidget")
	f.WindowWidget.SetSizePolicy2(widgets.QSizePolicy__Expanding | widgets.QSizePolicy__Maximum , widgets.QSizePolicy__Expanding)

	// windowVLayout is the following structure layout
	// +-----------+
	// |           |
	// +-----------+
	// |           |
	// +-----------+
	// |           |
	// +-----------+
        f.WindowVLayout = widgets.NewQVBoxLayout2(f.WindowWidget)
        f.WindowVLayout.SetContentsMargins(f.borderSize, f.borderSize, f.borderSize, 0)
        f.WindowVLayout.SetContentsMargins(0, 0, 0, 0)
	f.WindowVLayout.SetSpacing(0)
	f.WindowWidget.SetLayout(f.WindowVLayout)

	// create titlebar widget
	f.TitleBar = widgets.NewQWidget(f.WindowWidget, 0)
        f.TitleBar.SetObjectName("titleBar")
	f.TitleBar.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	//f.TitleBar.ConnectEventFilter(f.EventFilter)

	// titleBarLayout is the following structure layout
	// +--+--+--+--+
	// |  |  |  |  |
	// +--+--+--+--+
        f.TitleBarLayout = widgets.NewQHBoxLayout2(f.TitleBar)
        f.TitleBarLayout.SetContentsMargins(0, 0, 0, 0)

        f.TitleLabel = widgets.NewQLabel(nil, 0)
        f.TitleLabel.SetObjectName("TitleLabel")
        f.TitleLabel.SetAlignment(core.Qt__AlignCenter)

	if runtime.GOOS == "darwin" {
		f.SetTitleBarButtonsForDarwin()
	} else {
		f.SetTitleBarButtons()
	}


	// create window content
        f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// Set widget to layout
        f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
        f.WindowVLayout.AddWidget(f.Content, 0, 0)

        f.Layout.AddWidget(f.WindowWidget, 0, 0)
}

func (f *QFramelessWindow) SetTitleBarButtons() {
	iconSize := 14
        f.TitleBarLayout.SetSpacing(iconSize/2)
	f.IconMinimize = svg.NewQSvgWidget(nil)
	f.IconMinimize.SetFixedSize2(iconSize, iconSize)
	f.IconMinimize.SetObjectName("IconMinimize")
	f.IconMaximize = svg.NewQSvgWidget(nil)
	f.IconMaximize.SetFixedSize2(iconSize, iconSize)
	f.IconMaximize.SetObjectName("IconMaximize")
	f.IconRestore = svg.NewQSvgWidget(nil)
	f.IconRestore.SetFixedSize2(iconSize, iconSize)
	f.IconRestore.SetObjectName("IconRestore")
	f.IconClose = svg.NewQSvgWidget(nil)
	f.IconClose.SetFixedSize2(iconSize, iconSize)
	f.IconClose.SetObjectName("IconClose")

	f.IconMinimize.Hide()
	f.IconMaximize.Hide()
	f.IconRestore.Hide()
	f.IconClose.Hide()

        f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignRight)
        f.TitleBarLayout.AddWidget(f.TitleLabel, 0, 0)
        f.TitleBarLayout.AddWidget(f.IconMinimize, 0, 0)
        f.TitleBarLayout.AddWidget(f.IconMaximize, 0, 0)
        f.TitleBarLayout.AddWidget(f.IconRestore, 0, 0)
        f.TitleBarLayout.AddWidget(f.IconClose, 0, 0)
}

func (f *QFramelessWindow) SetTitleBarButtonsForDarwin() {
	btnSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__ToolButton)
	f.BtnMinimize = widgets.NewQToolButton(f.TitleBar)
	f.BtnMinimize.SetObjectName("BtnMinimize")
	f.BtnMinimize.SetSizePolicy(btnSizePolicy)
	
	f.BtnMaximize = widgets.NewQToolButton(f.TitleBar)
	f.BtnMaximize.SetObjectName("BtnMaximize")
	f.BtnMaximize.SetSizePolicy(btnSizePolicy)
	
	f.BtnRestore = widgets.NewQToolButton(f.TitleBar)
	f.BtnRestore.SetObjectName("BtnRestore")
	f.BtnRestore.SetSizePolicy(btnSizePolicy)
	f.BtnRestore.SetVisible(false)
	
	f.BtnClose = widgets.NewQToolButton(f.TitleBar)
	f.BtnClose.SetObjectName("BtnClose")
	f.BtnClose.SetSizePolicy(btnSizePolicy)
	
        f.TitleBarLayout.SetSpacing(0)
	f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignLeft)
	f.TitleBarLayout.AddWidget(f.BtnClose, 0, 0)
	f.TitleBarLayout.AddWidget(f.BtnMinimize, 0, 0)
	f.TitleBarLayout.AddWidget(f.BtnMaximize, 0, 0)
	f.TitleBarLayout.AddWidget(f.BtnRestore, 0, 0)
	f.TitleBarLayout.AddWidget(f.TitleLabel, 0, 0)
}

func (f *QFramelessWindow) SetAttribute() {
	f.Widget.Window().SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.Widget.Window().SetAttribute(core.Qt__WA_NoSystemBackground, true)
	f.Widget.Window().SetAttribute(core.Qt__WA_Hover, true)
	f.Widget.Window().SetMouseTracking(true)
}

func (f *QFramelessWindow) SetWidgetColor(color string) {
	f.windowColor = color
	style := fmt.Sprintf("background-color: %s", color)
	// f.Widget.SetStyleSheet(fmt.Sprintf(" .QFramelessWindow { border: 2px solid #ccc; border-radius: 11px; %s}", style))
	// f.Widget.Window().SetStyleSheet("* { background-color: rgba(0, 0, 0, 0); }")
	f.Widget.SetStyleSheet("* { background-color: rgba(0, 0, 0, 0); }")

	// f.WindowWidget.SetStyleSheet(fmt.Sprintf(" .QWidget { %s}", style))
	// f.WindowWidget.SetStyleSheet(fmt.Sprintf(" .QFrame { border: 2px solid %s; border-radius: 6px; background-color: rgba(0, 0, 0, 0); }", color))
	f.WindowWidget.SetStyleSheet(fmt.Sprintf(" .QFrame { border: 1px solid %s; padding: 6px; border-radius: 6px; %s; }", color, style))
	// f.TitleBar.SetStyleSheet(fmt.Sprintf(" .QWidget { %s; }", style))

	if runtime.GOOS == "darwin" {
		// padding titlebar
		f.TitleLabel.SetStyleSheet(" * {padding-right: 60px}")
		f.SetWindowButtonColorInDarwin()
	} else {
		// padding titlebar
		f.TitleLabel.SetStyleSheet(" * {padding-left: 70px}")
		f.IconMinimize.SetStyleSheet(`
		#IconMinimize { 
			background-color:none;
			border:none;
		}
		#IconMinimize:hover { 
			background-color:none;
			border:none;
		}
		`)
		f.IconMaximize.SetStyleSheet(`
		#IconMaximize { 
			background-color:none;
			border:none;
		}
		#IconMaximize:hover { 
			background-color:none;
			border:none;
		}
		`)

		f.IconRestore.SetStyleSheet(`
		#IconRestore { 
			background-color:none;
			border:none;
		}
		#IconRestore:hover { 
			background-color:none;
			border:none;
		}
		`)

		f.IconClose.SetStyleSheet(`
		#IconClose { 
			background-color:none;
			border:none;
		}
		#IconClose:hover { 
			background-color:none;
			border:none;
		}
		`)
	}
}

func (f *QFramelessWindow) SetWindowButtonColorInDarwin() {
	window := f.Widget.Window()
	var baseStyle, restoreAndMaximizeColor, minimizeColor, closeColor string
	baseStyle = ` #BtnMinimize, #BtnMaximize, #BtnRestore, #BtnClose {
		min-width: 10px;
		min-height: 10px;
		max-width: 10px;
		max-height: 10px;
		border-radius: 6px;
		border-width: 1px;
		border-style: solid;
		margin: 4px;
	}`
	if window.IsActiveWindow() {
		restoreAndMaximizeColor = `
			#BtnRestore, #BtnMaximize {
				background-color: rgb(53, 202, 74);
				border-color: rgb(34, 182, 52);
			}
		`
		minimizeColor = `
			#BtnMinimize {
				background-color: rgb(253, 190, 65);
				border-color: rgb(239, 170, 47);
			}
		`
		closeColor = `
			#BtnClose {
				background-color: rgb(252, 98, 93);
				border-color: rgb(239, 75, 71);
			}
		`
	} else {
		restoreAndMaximizeColor = `
			#BtnRestore, #BtnMaximize {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgb(128, 128, 128, 0.2);
			}
		`
		minimizeColor = `
			#BtnMinimize {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgb(128, 128, 128, 0.2);
			}
		`
		closeColor = `
			#BtnClose {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgb(128, 128, 128, 0.2);
			}
		`
	}
	MaximizeColorHover := `
		#BtnMaximize:hover {
			background-color: rgb(53, 202, 74);
			border-color: rgb(34, 182, 52);
			background-image: url(":/icons/MaximizeHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	RestoreColorHover := `
		#BtnRestore:hover {
			background-color: rgb(53, 202, 74);
			border-color: rgb(34, 182, 52);
			background-image: url(":/icons/RestoreHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	minimizeColorHover := `
		#BtnMinimize:hover {
			background-color: rgb(253, 190, 65);
			border-color: rgb(239, 170, 47);
			background-image: url(":/icons/MinimizeHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	closeColorHover := `
		#BtnClose:hover {
			background-color: rgb(252, 98, 93);
			border-color: rgb(239, 75, 71);
			background-image: url(":/icons/CloseHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	f.BtnMinimize.SetStyleSheet(baseStyle+minimizeColor+minimizeColorHover)
	f.BtnMaximize.SetStyleSheet(baseStyle+restoreAndMaximizeColor+MaximizeColorHover)
	f.BtnRestore.SetStyleSheet(baseStyle+restoreAndMaximizeColor+RestoreColorHover)
	f.BtnClose.SetStyleSheet(baseStyle+closeColor+closeColorHover)
}

func (f *QFramelessWindow) SetWindowFlags() {
	f.Widget.Window().SetWindowFlag(core.Qt__Window, true)
	f.Widget.Window().SetWindowFlag(core.Qt__FramelessWindowHint, true)
	f.Widget.Window().SetWindowFlag(core.Qt__WindowSystemMenuHint, true)
}

func (f *QFramelessWindow) SetTitle(title string) {
	f.TitleLabel.SetText(title)
}

func (f *QFramelessWindow) SetTitleColor(color string) {
	f.TitleLabel.SetStyleSheet(fmt.Sprintf(" *{padding-right: 60px; color: %s; }", color))

	if runtime.GOOS != "darwin" {
		SvgMinimize := fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M20,14H4V10H20" />
		</svg>
		`, color) 
		f.IconMinimize.Load2(core.NewQByteArray2(SvgMinimize, len(SvgMinimize)))

		SvgMaximize := fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M4,4H20V20H4V4M6,8V18H18V8H6Z" />
		</svg>
		`, color) 
		f.IconMaximize.Load2(core.NewQByteArray2(SvgMaximize, len(SvgMaximize)))

		SvgRestore := fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M4,8H8V4H20V16H16V20H4V8M16,8V14H18V6H10V8H16M6,12V18H14V12H6Z" />
		</svg>
		`, color) 
		f.IconRestore.Load2(core.NewQByteArray2(SvgRestore, len(SvgRestore)))

		SvgClose := fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M13.46,12L19,17.54V19H17.54L12,13.46L6.46,19H5V17.54L10.54,12L5,6.46V5H6.46L12,10.54L17.54,5H19V6.46L13.46,12Z" />
		</svg>
		`, color) 
		f.IconClose.Load2(core.NewQByteArray2(SvgClose, len(SvgClose)))

		f.IconMinimize.Show()
		f.IconMaximize.Show()
		f.IconRestore.Show()
		f.IconRestore.SetVisible(false)
		f.IconClose.Show()
	}
}

func (f *QFramelessWindow) SetContent(layout widgets.QLayout_ITF) {
	f.Content.SetLayout(layout)
}

func (f *QFramelessWindow) UpdateWidget() {
	f.Widget.Update()
	f.Widget.Window().Update()
}

func (f *QFramelessWindow) SetWindowActions() {


	// Ref: https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341#37507341
	f.Widget.Window().ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
	// f.WindowWidget.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
		switch event.Type() {
		case core.QEvent__ActivationChange :
			if runtime.GOOS == "darwin" {
				f.SetWindowButtonColorInDarwin()
			}

		case core.QEvent__HoverMove :
	 		f.updateCursorShape(e.GlobalPos())

		case core.QEvent__Leave :
			f.Widget.Window().UnsetCursor()

		case core.QEvent__MouseMove :
			if f.isDragStart {
				startPos := f.Widget.Window().FrameGeometry().TopLeft()
				newX :=startPos.X() + e.Pos().X() - f.dragPos.X()
				newY :=startPos.Y() + e.Pos().Y() - f.dragPos.Y()
				newPoint := core.NewQPoint2(newX, newY)
				f.Widget.Window().Move(newPoint)
			}
			if f.pressedEdge != None {

				left := f.Widget.Window().FrameGeometry().Left()
				top := f.Widget.Window().FrameGeometry().Top()
				right := f.Widget.Window().FrameGeometry().Right()
				bottom := f.Widget.Window().FrameGeometry().Bottom()

				switch f.pressedEdge {
				case Top:
					top = e.GlobalPos().Y()
				case Bottom:
					bottom = e.GlobalPos().Y()
				case Left:
					left = e.GlobalPos().X()
				case Right:
					right = e.GlobalPos().X()
				case TopLeft:
					top = e.GlobalPos().Y()
					left = e.GlobalPos().X()
				case TopRight:
					top = e.GlobalPos().Y()
					right = e.GlobalPos().X()
				case BottomLeft:
					bottom = e.GlobalPos().Y()
					left = e.GlobalPos().X()
				case BottomRight:
					bottom = e.GlobalPos().Y()
					right = e.GlobalPos().X()
				default:
				}

				topLeftPoint := core.NewQPoint2(left, top)
				rightBottomPoint := core.NewQPoint2(right, bottom)
				newRect := core.NewQRect2(topLeftPoint, rightBottomPoint)
				if newRect.Width() < f.Widget.Window().MinimumWidth() {
					left = f.Widget.Window().FrameGeometry().X()
				}
				if newRect.Height() < f.Widget.Window().MinimumHeight() {
					top = f.Widget.Window().FrameGeometry().Y()
				}
				topLeftPoint = core.NewQPoint2(left, top)
				rightBottomPoint = core.NewQPoint2(right, bottom)
				newRect = core.NewQRect2(topLeftPoint, rightBottomPoint)

				f.Widget.Window().SetGeometry(newRect)
			}
		case core.QEvent__MouseButtonPress :
			f.pressedEdge = f.calcCursorPos(e.GlobalPos(), f.Widget.Window().FrameGeometry())
			if f.pressedEdge != None {
				margins := core.NewQMargins2(f.borderSize*2, f.borderSize, f.borderSize*2, f.borderSize*2)
				if f.Widget.Window().Rect().MarginsRemoved(margins).Contains3(e.Pos().X(), e.Pos().Y()) {
					f.isDragStart = true
					f.dragPos = e.Pos()
				}
			}
		case core.QEvent__MouseButtonRelease :
			f.isDragStart = false
			f.pressedEdge  = None

		default:
			//fmt.Println(event.Type())
		}

		return f.Widget.EventFilter(watched, event)
	})
}

func (f *QFramelessWindow) updateCursorShape(pos *core.QPoint) {
	if f.Widget.Window().IsFullScreen() || f.Widget.Window().IsMaximized() {
		if f.isCursorChanged {
			f.Widget.Window().UnsetCursor()
		}
	}
	hoverEdge := f.calcCursorPos(pos, f.Widget.Window().FrameGeometry())
	f.isCursorChanged = true
	cursor := gui.NewQCursor()
	switch hoverEdge {
	case Top, Bottom:
		cursor.SetShape(core.Qt__SizeVerCursor)
		f.Widget.Window().SetCursor(cursor)
	case Left, Right:
		cursor.SetShape(core.Qt__SizeHorCursor)
		f.Widget.Window().SetCursor(cursor)
	case TopLeft, BottomRight:
		cursor.SetShape(core.Qt__SizeFDiagCursor)
		f.Widget.Window().SetCursor(cursor)
	case TopRight, BottomLeft:
		cursor.SetShape(core.Qt__SizeBDiagCursor)
		f.Widget.Window().SetCursor(cursor)
	default:
		f.Widget.Window().UnsetCursor()
		f.isCursorChanged = false
	}
}

func (f *QFramelessWindow) calcCursorPos(pos *core.QPoint, rect *core.QRect) Edge {
	borderSize := f.borderSize + 1
	doubleBorderSize := borderSize * 2
	var onLeft, onRight, onBottom, onTop, onBottomLeft, onBottomRight, onTopRight, onTopLeft bool

	onBottomLeft = (pos.X() <= (rect.X() + doubleBorderSize)) && pos.X() >= rect.X() &&
		       (pos.Y() <= (rect.Y() + rect.Height())) && (pos.Y() >= (rect.Y() + rect.Height() - doubleBorderSize))
	if onBottomLeft {
		return BottomLeft
	}
	
	onBottomRight = (pos.X() >= (rect.X() + rect.Width() - doubleBorderSize)) && (pos.X() <= (rect.X() + rect.Width())) &&
	                (pos.Y() >= (rect.Y() + rect.Height() - doubleBorderSize)) && (pos.Y() <= (rect.Y() + rect.Height()))
	if onBottomRight {
		return BottomRight
	}
	
	onTopRight = (pos.X() >= (rect.X() + rect.Width() - doubleBorderSize)) && (pos.X() <= (rect.X() + rect.Width())) &&
		     (pos.Y() >= rect.Y()) && (pos.Y() <= (rect.Y() + doubleBorderSize))
	if onTopRight {
		return TopRight
	}
	
	onTopLeft = pos.X() >= rect.X() && (pos.X() <= (rect.X() + doubleBorderSize)) &&
		    pos.Y() >= rect.Y() && (pos.Y() <= (rect.Y() + doubleBorderSize))
	if onTopLeft {
		return TopLeft
	}

	onLeft = (pos.X() >= (rect.X() - doubleBorderSize)) && (pos.X() <= (rect.X() + doubleBorderSize)) &&
	         (pos.Y() <= (rect.Y() + rect.Height() - doubleBorderSize)) &&
		 (pos.Y() >= rect.Y() + doubleBorderSize)
	if onLeft {
		return Left
	}

	onRight = (pos.X() >= (rect.X() + rect.Width() - doubleBorderSize)) &&
	          (pos.X() <= (rect.X() + rect.Width())) &&
		  (pos.Y() >= (rect.Y() + doubleBorderSize)) && (pos.Y() <= (rect.Y() + rect.Height() - doubleBorderSize))
	if onRight {
		return Right
	}
	
	onBottom = (pos.X() >= (rect.X() + doubleBorderSize)) && (pos.X() <= (rect.X() + rect.Width() - doubleBorderSize)) &&
	           (pos.Y() >= (rect.Y() + rect.Height() - doubleBorderSize)) && (pos.Y() <= (rect.Y() + rect.Height()))
	if onBottom {
		return Bottom
	}

	onTop = (pos.X() >= (rect.X() + borderSize)) && (pos.X() <= (rect.X() + rect.Width() - borderSize)) &&
		(pos.Y() >= rect.Y()) && (pos.Y() <= (rect.Y() + borderSize))
	if onTop {
		return Top
	}
	
	
	return None
}
