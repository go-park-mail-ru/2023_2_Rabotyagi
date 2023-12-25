// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package responses

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

func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses(in *jlexer.Lexer, out *ResponseSuccessful) {
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
			(out.Body).UnmarshalEasyJSON(in)
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses(out *jwriter.Writer, in ResponseSuccessful) {
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
		(in.Body).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseSuccessful) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseSuccessful) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseSuccessful) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseSuccessful) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses1(in *jlexer.Lexer, out *ResponseRedirect) {
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
			(out.Body).UnmarshalEasyJSON(in)
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses1(out *jwriter.Writer, in ResponseRedirect) {
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
		(in.Body).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseRedirect) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseRedirect) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseRedirect) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseRedirect) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses1(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses2(in *jlexer.Lexer, out *ResponseID) {
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
			(out.Body).UnmarshalEasyJSON(in)
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses2(out *jwriter.Writer, in ResponseID) {
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
		(in.Body).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseID) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseID) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseID) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseID) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses2(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses3(in *jlexer.Lexer, out *ResponseBodyRedirect) {
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
		case "redirect_url":
			out.RedirectURL = string(in.String())
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses3(out *jwriter.Writer, in ResponseBodyRedirect) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"redirect_url\":"
		out.RawString(prefix[1:])
		out.String(string(in.RedirectURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseBodyRedirect) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseBodyRedirect) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseBodyRedirect) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseBodyRedirect) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses3(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses4(in *jlexer.Lexer, out *ResponseBodyID) {
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses4(out *jwriter.Writer, in ResponseBodyID) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseBodyID) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseBodyID) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseBodyID) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseBodyID) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses4(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses5(in *jlexer.Lexer, out *ResponseBodyError) {
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
		case "error":
			out.Error = string(in.String())
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses5(out *jwriter.Writer, in ResponseBodyError) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix[1:])
		out.String(string(in.Error))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseBodyError) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseBodyError) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseBodyError) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseBodyError) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses5(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses6(in *jlexer.Lexer, out *ResponseBody) {
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
		case "message":
			out.Message = string(in.String())
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses6(out *jwriter.Writer, in ResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses6(l, v)
}
func easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses7(in *jlexer.Lexer, out *ErrorResponse) {
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
			(out.Body).UnmarshalEasyJSON(in)
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
func easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses7(out *jwriter.Writer, in ErrorResponse) {
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
		(in.Body).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ErrorResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrorResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson576be6ddEncodeGithubComGoParkMailRu20232RabotyagiPkgResponses7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrorResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrorResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson576be6ddDecodeGithubComGoParkMailRu20232RabotyagiPkgResponses7(l, v)
}
