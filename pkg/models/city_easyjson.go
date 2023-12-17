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

func easyjson66d84ff1DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(in *jlexer.Lexer, out *City) {
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
func easyjson66d84ff1EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(out *jwriter.Writer, in City) {
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v City) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66d84ff1EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v City) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66d84ff1EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *City) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66d84ff1DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *City) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66d84ff1DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(l, v)
}