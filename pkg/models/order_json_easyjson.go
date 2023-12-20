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

func easyjsonFdf1e785DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(in *jlexer.Lexer, out *orderJSON) {
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
			out.ID = uint64(in.Uint64())
		case "owner_id":
			out.OwnerID = uint64(in.Uint64())
		case "product_id":
			out.ProductID = uint64(in.Uint64())
		case "count":
			out.Count = uint32(in.Uint32())
		case "status":
			out.Status = uint8(in.Uint8())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updated_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
		case "closed_at":
			if in.IsNull() {
				in.Skip()
				out.ClosedAt = nil
			} else {
				if out.ClosedAt == nil {
					out.ClosedAt = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ClosedAt).UnmarshalJSON(data))
				}
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
func easyjsonFdf1e785EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(out *jwriter.Writer, in orderJSON) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"owner_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.OwnerID))
	}
	{
		const prefix string = ",\"product_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.ProductID))
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.Count))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.Uint8(uint8(in.Status))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"updated_at\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"closed_at\":"
		out.RawString(prefix)
		if in.ClosedAt == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.ClosedAt).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v orderJSON) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFdf1e785EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v orderJSON) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFdf1e785EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *orderJSON) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFdf1e785DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *orderJSON) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFdf1e785DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(l, v)
}
