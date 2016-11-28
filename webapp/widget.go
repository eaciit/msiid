package main

type Card struct {
	ID      string `bson:"_id" gridcolumn:"1"`
	Title   string `gridcolumn:"2"`
	Widgets []IWidget
}

type IWidget interface {
	Run(parm interface{}) error
	ID() string
	SetID(id string)
	Widgets() []IWidget
	UpdateWidget(IWidget) error
	DeleteWidget(widgetID string) error
}

type WidgetBase struct {
	_id      string
	_widgets []IWidget
}

func (wb *WidgetBase) ID() string {
	return wb._id
}

func (wb *WidgetBase) SetID(id string) {
	wb._id = id
}

func (wb *WidgetBase) Run(ps interface{}) error {
	return nil
}

func (wb *WidgetBase) Widgets() []IWidget {
	return wb._widgets
}

func (wb *WidgetBase) UpdateWidget(w IWidget) error {
	return nil
}

func (wb *WidgetBase) DeleteWidget(id string) error {
	return nil
}
