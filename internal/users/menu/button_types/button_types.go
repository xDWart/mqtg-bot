package button_types

type ButtonType int

const (
	SYSTEM ButtonType = iota - 1
	FOLDER
	SINGLE_VALUE
	TOGGLE
	MULTI_VALUE
	PRINT_LAST_SUB_VALUE
	DRAW_CHART
	COUNT_BUTTON_TYPES
)

var button_types_strs = [...]string{
	FOLDER:               "📂 Folder",
	SINGLE_VALUE:         "🔘 Single-value button",
	TOGGLE:               "🔄 Toggle button",
	MULTI_VALUE:          "🚦 Multi-value button",
	PRINT_LAST_SUB_VALUE: "📤 Print last sub value",
	DRAW_CHART:           "📈 Draw chart",
}

var _ = button_types_strs[COUNT_BUTTON_TYPES-1]

func (bt ButtonType) String() string {
	return button_types_strs[bt]
}

func (bt ButtonType) TypeString() string {
	return "Type: " + button_types_strs[bt]
}

func (bt *ButtonType) NextType(skipFolder bool) ButtonType {
	*bt = (*bt + 1) % COUNT_BUTTON_TYPES
	if skipFolder && *bt == FOLDER {
		bt.NextType(skipFolder)
	}
	return *bt
}
