// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF9d30717DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(in *jlexer.Lexer, out *categoryJson) {
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
		case "name":
			out.Name = string(in.String())
		case "parent_id":
			if in.IsNull() {
				in.Skip()
				out.ParentID = nil
			} else {
				if out.ParentID == nil {
					out.ParentID = new(uint64)
				}
				*out.ParentID = uint64(in.Uint64())
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
func easyjsonF9d30717EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(out *jwriter.Writer, in categoryJson) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"parent_id\":"
		out.RawString(prefix)
		if in.ParentID == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.ParentID))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v categoryJson) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF9d30717EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v categoryJson) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF9d30717EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *categoryJson) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF9d30717DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *categoryJson) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF9d30717DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(l, v)
}