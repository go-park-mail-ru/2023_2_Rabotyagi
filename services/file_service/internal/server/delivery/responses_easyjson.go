// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package delivery

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

func easyjson559270aeDecodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery(in *jlexer.Lexer, out *ResponseURLs) {
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
		case "status":
			out.Status = int(in.Int())
		case "body":
			easyjson559270aeDecodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery1(in, &out.Body)
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
func easyjson559270aeEncodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery(out *jwriter.Writer, in ResponseURLs) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Status))
	}
	{
		const prefix string = ",\"body\":"
		out.RawString(prefix)
		easyjson559270aeEncodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery1(out, in.Body)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseURLs) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseURLs) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseURLs) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseURLs) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery(l, v)
}
func easyjson559270aeDecodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery1(in *jlexer.Lexer, out *ResponseURLBody) {
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
		case "urls":
			if in.IsNull() {
				in.Skip()
				out.SlURL = nil
			} else {
				in.Delim('[')
				if out.SlURL == nil {
					if !in.IsDelim(']') {
						out.SlURL = make([]string, 0, 4)
					} else {
						out.SlURL = []string{}
					}
				} else {
					out.SlURL = (out.SlURL)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.SlURL = append(out.SlURL, v1)
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
func easyjson559270aeEncodeGithubComGoParkMailRu20232RabotyagiServicesFileServiceInternalServerDelivery1(out *jwriter.Writer, in ResponseURLBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"urls\":"
		out.RawString(prefix[1:])
		if in.SlURL == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.SlURL {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
