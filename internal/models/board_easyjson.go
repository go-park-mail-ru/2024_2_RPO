// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson202377feDecodeRPOBackInternalModels(in *jlexer.Lexer, out *MemberWithPermissions) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			if in.IsNull() {
				in.Skip()
				out.User = nil
			} else {
				if out.User == nil {
					out.User = new(UserProfile)
				}
				(*out.User).UnmarshalEasyJSON(in)
			}
		case "role":
			out.Role = string(in.String())
		case "addedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.AddedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
		case "addedBy":
			if in.IsNull() {
				in.Skip()
				out.AddedBy = nil
			} else {
				if out.AddedBy == nil {
					out.AddedBy = new(UserProfile)
				}
				(*out.AddedBy).UnmarshalEasyJSON(in)
			}
		case "updatedBy":
			if in.IsNull() {
				in.Skip()
				out.UpdatedBy = nil
			} else {
				if out.UpdatedBy == nil {
					out.UpdatedBy = new(UserProfile)
				}
				(*out.UpdatedBy).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels(out *jwriter.Writer, in MemberWithPermissions) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		if in.User == nil {
			out.RawString("null")
		} else {
			(*in.User).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"role\":"
		out.RawString(prefix)
		out.String(string(in.Role))
	}
	{
		const prefix string = ",\"addedAt\":"
		out.RawString(prefix)
		out.Raw((in.AddedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"updatedAt\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"addedBy\":"
		out.RawString(prefix)
		if in.AddedBy == nil {
			out.RawString("null")
		} else {
			(*in.AddedBy).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"updatedBy\":"
		out.RawString(prefix)
		if in.UpdatedBy == nil {
			out.RawString("null")
		} else {
			(*in.UpdatedBy).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MemberWithPermissions) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MemberWithPermissions) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MemberWithPermissions) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MemberWithPermissions) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels1(in *jlexer.Lexer, out *InviteLink) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "inviteLinkUuid":
			out.InviteLinkUUID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels1(out *jwriter.Writer, in InviteLink) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"inviteLinkUuid\":"
		out.RawString(prefix[1:])
		out.String(string(in.InviteLinkUUID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InviteLink) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InviteLink) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InviteLink) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InviteLink) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels1(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels2(in *jlexer.Lexer, out *Comment) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "text":
			out.Text = string(in.String())
		case "isEdited":
			out.IsEdited = bool(in.Bool())
		case "createdBy":
			if in.IsNull() {
				in.Skip()
				out.CreatedBy = nil
			} else {
				if out.CreatedBy == nil {
					out.CreatedBy = new(UserProfile)
				}
				(*out.CreatedBy).UnmarshalEasyJSON(in)
			}
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels2(out *jwriter.Writer, in Comment) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsEdited))
	}
	{
		const prefix string = ",\"createdBy\":"
		out.RawString(prefix)
		if in.CreatedBy == nil {
			out.RawString("null")
		} else {
			(*in.CreatedBy).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Comment) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Comment) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Comment) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Comment) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels2(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels3(in *jlexer.Lexer, out *Column) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "title":
			out.Title = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels3(out *jwriter.Writer, in Column) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Column) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Column) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Column) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Column) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels3(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels4(in *jlexer.Lexer, out *CheckListField) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "title":
			out.Title = string(in.String())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "isDone":
			out.IsDone = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels4(out *jwriter.Writer, in CheckListField) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"isDone\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsDone))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CheckListField) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CheckListField) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CheckListField) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CheckListField) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels4(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels5(in *jlexer.Lexer, out *CardDetails) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "card":
			if in.IsNull() {
				in.Skip()
				out.Card = nil
			} else {
				if out.Card == nil {
					out.Card = new(Card)
				}
				(*out.Card).UnmarshalEasyJSON(in)
			}
		case "checkList":
			if in.IsNull() {
				in.Skip()
				out.CheckList = nil
			} else {
				in.Delim('[')
				if out.CheckList == nil {
					if !in.IsDelim(']') {
						out.CheckList = make([]CheckListField, 0, 1)
					} else {
						out.CheckList = []CheckListField{}
					}
				} else {
					out.CheckList = (out.CheckList)[:0]
				}
				for !in.IsDelim(']') {
					var v1 CheckListField
					(v1).UnmarshalEasyJSON(in)
					out.CheckList = append(out.CheckList, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "attachments":
			if in.IsNull() {
				in.Skip()
				out.Attachments = nil
			} else {
				in.Delim('[')
				if out.Attachments == nil {
					if !in.IsDelim(']') {
						out.Attachments = make([]Attachment, 0, 1)
					} else {
						out.Attachments = []Attachment{}
					}
				} else {
					out.Attachments = (out.Attachments)[:0]
				}
				for !in.IsDelim(']') {
					var v2 Attachment
					(v2).UnmarshalEasyJSON(in)
					out.Attachments = append(out.Attachments, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "comments":
			if in.IsNull() {
				in.Skip()
				out.Comments = nil
			} else {
				in.Delim('[')
				if out.Comments == nil {
					if !in.IsDelim(']') {
						out.Comments = make([]Comment, 0, 1)
					} else {
						out.Comments = []Comment{}
					}
				} else {
					out.Comments = (out.Comments)[:0]
				}
				for !in.IsDelim(']') {
					var v3 Comment
					(v3).UnmarshalEasyJSON(in)
					out.Comments = append(out.Comments, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "assignedUsers":
			if in.IsNull() {
				in.Skip()
				out.AssignedUsers = nil
			} else {
				in.Delim('[')
				if out.AssignedUsers == nil {
					if !in.IsDelim(']') {
						out.AssignedUsers = make([]UserProfile, 0, 0)
					} else {
						out.AssignedUsers = []UserProfile{}
					}
				} else {
					out.AssignedUsers = (out.AssignedUsers)[:0]
				}
				for !in.IsDelim(']') {
					var v4 UserProfile
					(v4).UnmarshalEasyJSON(in)
					out.AssignedUsers = append(out.AssignedUsers, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels5(out *jwriter.Writer, in CardDetails) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"card\":"
		out.RawString(prefix[1:])
		if in.Card == nil {
			out.RawString("null")
		} else {
			(*in.Card).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"checkList\":"
		out.RawString(prefix)
		if in.CheckList == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.CheckList {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"attachments\":"
		out.RawString(prefix)
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v7, v8 := range in.Attachments {
				if v7 > 0 {
					out.RawByte(',')
				}
				(v8).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"comments\":"
		out.RawString(prefix)
		if in.Comments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.Comments {
				if v9 > 0 {
					out.RawByte(',')
				}
				(v10).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"assignedUsers\":"
		out.RawString(prefix)
		if in.AssignedUsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.AssignedUsers {
				if v11 > 0 {
					out.RawByte(',')
				}
				(v12).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CardDetails) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CardDetails) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CardDetails) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CardDetails) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels5(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels6(in *jlexer.Lexer, out *Card) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "cardUuid":
			out.UUID = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "coverImageUrl":
			out.CoverImageURL = string(in.String())
		case "columnId":
			out.ColumnID = int64(in.Int64())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
		case "deadline":
			if in.IsNull() {
				in.Skip()
				out.Deadline = nil
			} else {
				if out.Deadline == nil {
					out.Deadline = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Deadline).UnmarshalJSON(data))
				}
			}
		case "isDone":
			out.IsDone = bool(in.Bool())
		case "hasCheckList":
			out.HasCheckList = bool(in.Bool())
		case "hasAttachments":
			out.HasAttachments = bool(in.Bool())
		case "hasAssignedUsers":
			out.HasAssignedUsers = bool(in.Bool())
		case "hasComments":
			out.HasComments = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels6(out *jwriter.Writer, in Card) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"cardUuid\":"
		out.RawString(prefix)
		out.String(string(in.UUID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"coverImageUrl\":"
		out.RawString(prefix)
		out.String(string(in.CoverImageURL))
	}
	{
		const prefix string = ",\"columnId\":"
		out.RawString(prefix)
		out.Int64(int64(in.ColumnID))
	}
	{
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"updatedAt\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	if in.Deadline != nil {
		const prefix string = ",\"deadline\":"
		out.RawString(prefix)
		out.Raw((*in.Deadline).MarshalJSON())
	}
	{
		const prefix string = ",\"isDone\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsDone))
	}
	{
		const prefix string = ",\"hasCheckList\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasCheckList))
	}
	{
		const prefix string = ",\"hasAttachments\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasAttachments))
	}
	{
		const prefix string = ",\"hasAssignedUsers\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasAssignedUsers))
	}
	{
		const prefix string = ",\"hasComments\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasComments))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Card) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Card) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Card) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Card) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels6(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels7(in *jlexer.Lexer, out *BoardContent) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "myRole":
			out.MyRole = string(in.String())
		case "allCards":
			if in.IsNull() {
				in.Skip()
				out.Cards = nil
			} else {
				in.Delim('[')
				if out.Cards == nil {
					if !in.IsDelim(']') {
						out.Cards = make([]Card, 0, 0)
					} else {
						out.Cards = []Card{}
					}
				} else {
					out.Cards = (out.Cards)[:0]
				}
				for !in.IsDelim(']') {
					var v13 Card
					(v13).UnmarshalEasyJSON(in)
					out.Cards = append(out.Cards, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "allColumns":
			if in.IsNull() {
				in.Skip()
				out.Columns = nil
			} else {
				in.Delim('[')
				if out.Columns == nil {
					if !in.IsDelim(']') {
						out.Columns = make([]Column, 0, 2)
					} else {
						out.Columns = []Column{}
					}
				} else {
					out.Columns = (out.Columns)[:0]
				}
				for !in.IsDelim(']') {
					var v14 Column
					(v14).UnmarshalEasyJSON(in)
					out.Columns = append(out.Columns, v14)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "boardInfo":
			if in.IsNull() {
				in.Skip()
				out.BoardInfo = nil
			} else {
				if out.BoardInfo == nil {
					out.BoardInfo = new(Board)
				}
				(*out.BoardInfo).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels7(out *jwriter.Writer, in BoardContent) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"myRole\":"
		out.RawString(prefix[1:])
		out.String(string(in.MyRole))
	}
	{
		const prefix string = ",\"allCards\":"
		out.RawString(prefix)
		if in.Cards == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v15, v16 := range in.Cards {
				if v15 > 0 {
					out.RawByte(',')
				}
				(v16).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"allColumns\":"
		out.RawString(prefix)
		if in.Columns == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.Columns {
				if v17 > 0 {
					out.RawByte(',')
				}
				(v18).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"boardInfo\":"
		out.RawString(prefix)
		if in.BoardInfo == nil {
			out.RawString("null")
		} else {
			(*in.BoardInfo).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BoardContent) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BoardContent) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BoardContent) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BoardContent) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels7(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels8(in *jlexer.Lexer, out *Board) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "backgroundImageUrl":
			out.BackgroundImageURL = string(in.String())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
		case "lastVisitAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.LastVisitAt).UnmarshalJSON(data))
			}
		case "myInviteLinkUuid":
			if in.IsNull() {
				in.Skip()
				out.MyInviteUUID = nil
			} else {
				if out.MyInviteUUID == nil {
					out.MyInviteUUID = new(string)
				}
				*out.MyInviteUUID = string(in.String())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels8(out *jwriter.Writer, in Board) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"backgroundImageUrl\":"
		out.RawString(prefix)
		out.String(string(in.BackgroundImageURL))
	}
	{
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"updatedAt\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"lastVisitAt\":"
		out.RawString(prefix)
		out.Raw((in.LastVisitAt).MarshalJSON())
	}
	{
		const prefix string = ",\"myInviteLinkUuid\":"
		out.RawString(prefix)
		if in.MyInviteUUID == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.MyInviteUUID))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Board) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Board) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Board) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Board) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels8(l, v)
}
func easyjson202377feDecodeRPOBackInternalModels9(in *jlexer.Lexer, out *Attachment) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "originalName":
			out.OriginalName = string(in.String())
		case "fileName":
			out.FileName = string(in.String())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson202377feEncodeRPOBackInternalModels9(out *jwriter.Writer, in Attachment) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"originalName\":"
		out.RawString(prefix)
		out.String(string(in.OriginalName))
	}
	{
		const prefix string = ",\"fileName\":"
		out.RawString(prefix)
		out.String(string(in.FileName))
	}
	{
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Attachment) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson202377feEncodeRPOBackInternalModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Attachment) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson202377feEncodeRPOBackInternalModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Attachment) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson202377feDecodeRPOBackInternalModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Attachment) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson202377feDecodeRPOBackInternalModels9(l, v)
}