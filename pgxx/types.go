package pgxx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type typeInfo struct {
	Requested  string  `db:"requested"`
	OID        *uint32 `db:"oid"`
	ElementOID uint32  `db:"element_oid"`
	Schema     string  `db:"schema"`
	Name       string  `db:"name"`
	Kind       string  `db:"kind"`
	ArrayOID   uint32  `db:"array_oid"`
	ArrayName  string  `db:"array_name"`
	SubtypeOID uint32  `db:"subtype_oid"`
}

// RegisterTypes registers custom types and their array types on a connection, typically in
// pgxpool.Config.AfterConnect. Names resolve like in SQL, so unquoted names fold to lower case, and
// an array name such as "_mood" registers its element type. Like LoadTypes, it registers each type
// under its plain and its schema-qualified name.
//
// Enums and ranges are looked up in pg_type in about a millisecond, while LoadTypes can take
// hundreds of them. Other kinds are handed to LoadTypes. Nothing is cached across connections, so no
// connection shares an EnumCodec, which would be unsafe, and new connections pick up a recreated
// type while existing ones keep the old one until the pool replaces them.
//
// A type that does not exist is an error, so new connections fail if the code names a type that was
// renamed or dropped. Rename or drop a type only after the code that uses it is gone.
func RegisterTypes(ctx context.Context, c *pgx.Conn, names ...string) error {
	if len(names) == 0 {
		return nil
	}
	rows, err := c.Query(ctx, `
		SELECT n.requested, t.oid, COALESCE(e.oid, 0) AS element_oid,
			COALESCE(ns.nspname::text, '') AS schema, COALESCE(e.typname::text, '') AS name,
			COALESCE(e.typtype::text, '') AS kind, COALESCE(e.typarray, 0) AS array_oid,
			COALESCE(a.typname::text, '') AS array_name, COALESCE(r.rngsubtype, 0) AS subtype_oid
		FROM unnest($1::text[]) WITH ORDINALITY AS n(requested, ord)
		LEFT JOIN pg_type t ON t.oid = to_regtype(n.requested)
		LEFT JOIN pg_type e ON e.oid = CASE
			WHEN t.typtype = 'b' AND t.typcategory = 'A' AND t.typelem <> 0 THEN t.typelem
			ELSE t.oid
		END
		LEFT JOIN pg_namespace ns ON ns.oid = e.typnamespace
		LEFT JOIN pg_type a ON a.oid = e.typarray
		LEFT JOIN pg_range r ON r.rngtypid = e.oid
		ORDER BY n.ord`, names)
	if err != nil {
		return fmt.Errorf("look up types %v: %w", names, err)
	}
	infos, err := pgx.CollectRows(rows, pgx.RowToStructByName[typeInfo])
	if err != nil {
		return fmt.Errorf("look up types %v: %w", names, err)
	}
	for _, info := range infos {
		if info.OID == nil {
			return fmt.Errorf("type %s does not exist", info.Requested)
		}
	}

	// Enums go first, so a range over one of them finds its subtype whatever the order of names.
	m := c.TypeMap()
	for _, info := range infos {
		if info.Kind == "e" {
			register(m, info, &pgtype.EnumCodec{})
		}
	}
	var others []string
	for _, info := range infos {
		if info.Kind == "e" {
			continue
		}
		if info.Kind == "r" {
			if subtype, ok := m.TypeForOID(info.SubtypeOID); ok {
				register(m, info, &pgtype.RangeCodec{ElementType: subtype})
				continue
			}
		}
		// Built-in types are registered already, and LoadTypes would take them for arrays.
		if _, ok := m.TypeForOID(info.ElementOID); ok {
			continue
		}
		// LoadTypes matches names literally, so it gets the resolved ones.
		others = append(others, info.Schema+"."+info.Name)
		if info.ArrayOID != 0 {
			others = append(others, info.Schema+"."+info.ArrayName)
		}
	}

	if len(others) > 0 {
		types, err := c.LoadTypes(ctx, others)
		if err != nil {
			return fmt.Errorf("load types %v: %w", others, err)
		}
		m.RegisterTypes(types)
	}
	for _, info := range infos {
		if _, ok := m.TypeForOID(*info.OID); !ok {
			return fmt.Errorf("type %s could not be registered", info.Requested)
		}
	}
	return nil
}

// register adds a type and its array under their plain and schema-qualified names. The plain names
// go last, so lookups by OID return them.
func register(m *pgtype.Map, info typeInfo, codec pgtype.Codec) {
	element := &pgtype.Type{Name: info.Name, OID: info.ElementOID, Codec: codec}
	m.RegisterType(&pgtype.Type{Name: info.Schema + "." + info.Name, OID: info.ElementOID, Codec: codec})
	m.RegisterType(element)
	if info.ArrayOID != 0 {
		arrayCodec := &pgtype.ArrayCodec{ElementType: element}
		m.RegisterType(&pgtype.Type{Name: info.Schema + "." + info.ArrayName, OID: info.ArrayOID, Codec: arrayCodec})
		m.RegisterType(&pgtype.Type{Name: info.ArrayName, OID: info.ArrayOID, Codec: arrayCodec})
	}
}
