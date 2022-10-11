package test_utils

import (
  "unicode/utf8"
)

func toValidUTF8Char( v uint32 ) string {
  if v < 0x20 {
    v = 0x20
  }
  if v > 0x7F && v < 0xA0 {
    v = 0x20
  }
  c := rune( v );
  temp := make([]byte,4)
  size := utf8.EncodeRune( temp, c )
  return string(temp[:size])
}

func ToValidUTF8String( r []byte, l int ) string {
  if len( r ) > l * 3 {
    return ToValidUTF8String( r[:l*3], l )
  }
  if len( r ) < 1 {
    return ""
  }
  if len( r ) < 2 {
    v := uint32( r[ 0 ] )
    return toValidUTF8Char( v )
  }
  if len( r ) < 3 {
    v := ( uint32( r[ 0 ] ) << 8 ) | uint32( r[ 1 ] ) 
    return toValidUTF8Char( v )
  }
  if len( r ) < 4 {
    v := ( uint32( r[ 0 ] ) << 16 ) | ( uint32( r[ 1 ] ) << 8 ) | ( uint32( r[ 2 ] ) ) 
    return toValidUTF8Char( v )
  }
  v := ( uint32( r[ 0 ] ) << 16 ) | ( uint32( r[ 1 ] ) << 8 ) | ( uint32( r[ 2 ] ) ) 
  return toValidUTF8Char( v ) + ToValidUTF8String( r[3:], l )
}

func ToValidUTF8StringBiased( r []byte, l int ) string {
  if l < 1 {
    return ""
  }
  length := len( r )
  if length < 1 {
    return ""
  }
  if r[ 0 ] < 0x20 {
    return ToValidUTF8StringBiased( []byte{ 0x20 }, 1 ) + ToValidUTF8StringBiased( r[1:], l - 1 )
  }
  if r[ 0 ] < 0x80 {
    c := rune( r[ 0 ] );
    temp := make([]byte,4)
    size := utf8.EncodeRune( temp, c )
    return string(temp[:size]) + ToValidUTF8StringBiased( r[1:], l - 1 )
  }
  if length < 2 {
    return ToValidUTF8StringBiased( []byte{ r[ 0 ] & 0x7F }, 1 ) + ToValidUTF8StringBiased( r[1:], l - 1 )
  }
  if r[ 0 ] < 0xE0 {
    c := rune( ( uint32( r[ 0 ] & 0x1F ) << 6 ) | uint32( r[ 1 ] & 0x3F ) );
    temp := make([]byte,4)
    size := utf8.EncodeRune( temp, c )
    return string(temp[:size]) + ToValidUTF8StringBiased( r[2:], l - 1 )
  }
  if length < 3 {
    return ToValidUTF8StringBiased( []byte{ r[ 0 ] & 0x1F | 0xC0, r[ 1 ] }, 2 ) + ToValidUTF8StringBiased( r[2:], l - 1 )
  }
  if r[ 0 ] < 0xF0 {
    c := rune( ( uint32( r[ 0 ] & 0x0F ) << 12 ) | ( uint32( r[ 1 ] & 0x3F ) << 6 ) | uint32( r[ 2 ] & 0x3F ) );
    temp := make([]byte,4)
    size := utf8.EncodeRune( temp, c )
    return string(temp[:size]) + ToValidUTF8StringBiased( r[3:], l - 1 )
  }
  if length < 4 {
    return ToValidUTF8StringBiased( []byte{ r[ 0 ] & 0x0F | 0xE0, r[ 1 ], r[ 2 ] }, 3 ) + ToValidUTF8StringBiased( r[3:], l - 1 )
  }
  c := rune( ( uint32( r[ 0 ] & 0x07 ) << 18 ) | ( uint32( r[ 1 ] & 0x3F ) << 12 ) | ( uint32( r[ 2 ] & 0x3F ) << 6 ) | uint32( r[ 3 ] & 0x3F ) );
  temp := make([]byte,4)
  size := utf8.EncodeRune( temp, c )
  return string(temp[:size]) + ToValidUTF8StringBiased( r[4:], l - 1 )
}

