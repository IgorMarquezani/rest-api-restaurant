package models

func (t Tab) FieldsAdress() []any {
	return []any{&t.Number, &t.RoomId, &t.Table, &t.RoomId, &t.Maded, &t.PayValue}
}
