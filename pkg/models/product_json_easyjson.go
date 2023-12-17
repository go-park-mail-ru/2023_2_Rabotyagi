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

func easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(in *jlexer.Lexer, out *productJson) {
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
		case "saler_id":
			out.SalerID = uint64(in.Uint64())
		case "category_id":
			out.CategoryID = uint64(in.Uint64())
		case "city_id":
			out.CityID = uint64(in.Uint64())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "price":
			out.Price = uint64(in.Uint64())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "premium_begin":
			if in.IsNull() {
				in.Skip()
				out.PremiumBegin = nil
			} else {
				if out.PremiumBegin == nil {
					out.PremiumBegin = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.PremiumBegin).UnmarshalJSON(data))
				}
			}
		case "premium_expire":
			if in.IsNull() {
				in.Skip()
				out.PremiumExpire = nil
			} else {
				if out.PremiumExpire == nil {
					out.PremiumExpire = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.PremiumExpire).UnmarshalJSON(data))
				}
			}
		case "views":
			out.Views = uint32(in.Uint32())
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
					easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels1(in, &v1)
					out.Images = append(out.Images, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "price_history":
			if in.IsNull() {
				in.Skip()
				out.PriceHistory = nil
			} else {
				in.Delim('[')
				if out.PriceHistory == nil {
					if !in.IsDelim(']') {
						out.PriceHistory = make([]PriceHistoryRecord, 0, 2)
					} else {
						out.PriceHistory = []PriceHistoryRecord{}
					}
				} else {
					out.PriceHistory = (out.PriceHistory)[:0]
				}
				for !in.IsDelim(']') {
					var v2 PriceHistoryRecord
					easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels2(in, &v2)
					out.PriceHistory = append(out.PriceHistory, v2)
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
func easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(out *jwriter.Writer, in productJson) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"saler_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.SalerID))
	}
	{
		const prefix string = ",\"category_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CategoryID))
	}
	{
		const prefix string = ",\"city_id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.CityID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Price))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"premium_begin\":"
		out.RawString(prefix)
		if in.PremiumBegin == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.PremiumBegin).MarshalJSON())
		}
	}
	{
		const prefix string = ",\"premium_expire\":"
		out.RawString(prefix)
		if in.PremiumExpire == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.PremiumExpire).MarshalJSON())
		}
	}
	{
		const prefix string = ",\"views\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.Views))
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
			for v3, v4 := range in.Images {
				if v3 > 0 {
					out.RawByte(',')
				}
				easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels1(out, v4)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"price_history\":"
		out.RawString(prefix)
		if in.PriceHistory == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.PriceHistory {
				if v5 > 0 {
					out.RawByte(',')
				}
				easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels2(out, v6)
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
func (v productJson) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v productJson) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *productJson) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *productJson) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels(l, v)
}
func easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels2(in *jlexer.Lexer, out *PriceHistoryRecord) {
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
		case "price":
			out.Price = uint64(in.Uint64())
		case "created_at":
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
func easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels2(out *jwriter.Writer, in PriceHistoryRecord) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Price))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}
func easyjsonB0091c22DecodeGithubComGoParkMailRu20232RabotyagiPkgModels1(in *jlexer.Lexer, out *Image) {
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
func easyjsonB0091c22EncodeGithubComGoParkMailRu20232RabotyagiPkgModels1(out *jwriter.Writer, in Image) {
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
