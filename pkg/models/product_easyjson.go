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

func easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels(in *jlexer.Lexer, out *ProductInFeed) {
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
		case "title":
			out.Title = string(in.String())
		case "price":
			out.Price = uint64(in.Uint64())
		case "city_id":
			out.CityID = uint64(in.Uint64())
		case "available_count":
			out.AvailableCount = uint32(in.Uint32())
		case "delivery":
			out.Delivery = bool(in.Bool())
		case "safe_deal":
			out.SafeDeal = bool(in.Bool())
		case "in_favourites":
			out.InFavourites = bool(in.Bool())
		case "is_active":
			out.IsActive = bool(in.Bool())
		case "premium":
			out.Premium = bool(in.Bool())
		case "images":
			if in.IsNull() {
				in.Skip()
				out.Images = nil
			} else {
				in.Delim('[')
				if out.Images == nil {
					if !in.IsDelim(']') {
						out.Images = make([]Image, 0, 4)
					} else {
						out.Images = []Image{}
					}
				} else {
					out.Images = (out.Images)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Image
					easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels1(in, &v1)
					out.Images = append(out.Images, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "favourites":
			out.Favourites = uint64(in.Uint64())
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
func easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels(out *jwriter.Writer, in ProductInFeed) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Price))
	}
	{
		const prefix string = ",\"city_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CityID))
	}
	{
		const prefix string = ",\"available_count\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.AvailableCount))
	}
	{
		const prefix string = ",\"delivery\":"
		out.RawString(prefix)
		out.Bool(bool(in.Delivery))
	}
	{
		const prefix string = ",\"safe_deal\":"
		out.RawString(prefix)
		out.Bool(bool(in.SafeDeal))
	}
	{
		const prefix string = ",\"in_favourites\":"
		out.RawString(prefix)
		out.Bool(bool(in.InFavourites))
	}
	{
		const prefix string = ",\"is_active\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsActive))
	}
	{
		const prefix string = ",\"premium\":"
		out.RawString(prefix)
		out.Bool(bool(in.Premium))
	}
	{
		const prefix string = ",\"images\":"
		out.RawString(prefix)
		if in.Images == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Images {
				if v2 > 0 {
					out.RawByte(',')
				}
				easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels1(out, v3)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"favourites\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Favourites))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProductInFeed) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProductInFeed) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProductInFeed) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProductInFeed) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels(l, v)
}
func easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels1(in *jlexer.Lexer, out *Image) {
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
		case "url":
			out.URL = string(in.String())
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
func easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels1(out *jwriter.Writer, in Image) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	out.RawByte('}')
}
func easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels2(in *jlexer.Lexer, out *ProductID) {
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
		case "product_id":
			out.ProductID = uint64(in.Uint64())
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
func easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels2(out *jwriter.Writer, in ProductID) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"product_id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ProductID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProductID) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProductID) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCf3f67efEncodeGithubComGoParkMailRu20232RabotyagiPkgModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProductID) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProductID) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCf3f67efDecodeGithubComGoParkMailRu20232RabotyagiPkgModels2(l, v)
}
