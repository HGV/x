package pgxx

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testSchema creates a schema with test types and returns a function that connects with it on the
// search path.
func testSchema(t *testing.T) (schema string, connect func() *pgx.Conn) {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	schema = fmt.Sprintf("pgxx_test_%d", time.Now().UnixNano())

	admin, err := pgx.Connect(ctx, url)
	require.NoError(t, err)
	_, err = admin.Exec(ctx, fmt.Sprintf(`
		CREATE SCHEMA %[1]s;
		CREATE TYPE %[1]s.mood AS ENUM ('sad', 'ok', 'happy');
		CREATE TYPE %[1]s.smallrange AS RANGE (subtype = int2);
		CREATE TYPE %[1]s.point2 AS (x int4, y int4);
		CREATE TYPE %[1]s.moodrange AS RANGE (subtype = %[1]s.mood);
		CREATE DOMAIN %[1]s.moods AS %[1]s.mood[];
		CREATE TYPE %[1]s._tag AS ENUM ('a');
		CREATE TYPE %[1]s.tag AS ENUM ('b');`, schema))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
		admin.Close(ctx)
	})

	connect = func() *pgx.Conn {
		config, err := pgx.ParseConfig(url)
		require.NoError(t, err)
		config.RuntimeParams["search_path"] = schema + ", public"
		c, err := pgx.ConnectConfig(ctx, config)
		require.NoError(t, err)
		t.Cleanup(func() { c.Close(ctx) })
		return c
	}
	return schema, connect
}

func TestRegisterTypesRegistersWhatLoadTypesRegisters(t *testing.T) {
	ctx := context.Background()
	schema, connect := testSchema(t)
	names := []string{"mood", "_mood", "smallrange", "_smallrange"}

	// Array and qualified names register the same types as plain ones.
	ours := connect()
	require.NoError(t, RegisterTypes(ctx, ours, "mood[]", schema+".smallrange"))
	reference := connect()
	loaded, err := reference.LoadTypes(ctx, names)
	require.NoError(t, err)
	reference.TypeMap().RegisterTypes(loaded)

	for _, name := range append(names, schema+".mood", schema+"._mood", schema+".smallrange", schema+"._smallrange") {
		want, ok := reference.TypeMap().TypeForName(name)
		require.True(t, ok, name)
		got, ok := ours.TypeMap().TypeForName(name)
		require.True(t, ok, name)

		assert.Equal(t, want.OID, got.OID, name)
		assert.Equal(t, reflect.TypeOf(want.Codec), reflect.TypeOf(got.Codec), name)
		if wantArray, ok := want.Codec.(*pgtype.ArrayCodec); ok {
			assert.Equal(t, wantArray.ElementType.OID, got.Codec.(*pgtype.ArrayCodec).ElementType.OID, name)
		}
		if wantRange, ok := want.Codec.(*pgtype.RangeCodec); ok {
			assert.Equal(t, wantRange.ElementType.OID, got.Codec.(*pgtype.RangeCodec).ElementType.OID, name)
		}
	}
}

func TestRegisteredEnumsAndRangesRoundTrip(t *testing.T) {
	ctx := context.Background()
	_, connect := testSchema(t)
	c := connect()
	require.NoError(t, RegisterTypes(ctx, c, "mood", "smallrange"))

	var mood string
	require.NoError(t, c.QueryRow(ctx, "SELECT $1::mood", "happy").Scan(&mood))
	assert.Equal(t, "happy", mood)

	var moods []string
	require.NoError(t, c.QueryRow(ctx, "SELECT $1::mood[]", []string{"sad", "ok"}).Scan(&moods))
	assert.Equal(t, []string{"sad", "ok"}, moods)

	var r pgtype.Range[pgtype.Int2]
	require.NoError(t, c.QueryRow(ctx, "SELECT '[2,10)'::smallrange").Scan(&r))
	assert.Equal(t, int16(2), r.Lower.Int16)
	assert.Equal(t, int16(10), r.Upper.Int16)

	var contains bool
	require.NoError(t, c.QueryRow(ctx, "SELECT $1::smallrange @> 9::int2", r).Scan(&contains))
	assert.True(t, contains)

	var ranges []pgtype.Range[pgtype.Int2]
	require.NoError(t, c.QueryRow(ctx, "SELECT ARRAY['[0,3)'::smallrange, '[5,)'::smallrange]").Scan(&ranges))
	require.Len(t, ranges, 2)
	assert.Equal(t, pgtype.Unbounded, ranges[1].UpperType)
}

func TestRegisterTypesHandlesDependentAndUnusualTypes(t *testing.T) {
	ctx := context.Background()
	_, connect := testSchema(t)
	c := connect()

	// A range listed before its enum, an array of a composite, a domain over an enum array and an
	// enum whose array is not called _tag, because another type has that name.
	require.NoError(t, RegisterTypes(ctx, c, "moodrange", "mood", "point2[]", "moods", "tag"))

	var tagArrayOID uint32
	var tagArrayName string
	require.NoError(t, c.QueryRow(ctx, "SELECT a.oid, a.typname::text FROM pg_type t JOIN pg_type a ON a.oid = t.typarray WHERE t.oid = 'tag'::regtype").
		Scan(&tagArrayOID, &tagArrayName))
	require.NotEqual(t, "_tag", tagArrayName)

	for name, codec := range map[string]pgtype.Codec{
		"moodrange":  &pgtype.RangeCodec{},
		"_moodrange": &pgtype.ArrayCodec{},
		"point2":     &pgtype.CompositeCodec{},
		"_point2":    &pgtype.ArrayCodec{},
		"moods":      &pgtype.ArrayCodec{},
		tagArrayName: &pgtype.ArrayCodec{},
	} {
		typ, ok := c.TypeMap().TypeForName(name)
		require.True(t, ok, name)
		assert.IsType(t, codec, typ.Codec, name)
	}
	array, _ := c.TypeMap().TypeForName(tagArrayName)
	assert.Equal(t, tagArrayOID, array.OID)
	_, ok := c.TypeMap().TypeForName("_tag")
	assert.False(t, ok, "_tag belongs to another type")

	var moods []string
	require.NoError(t, c.QueryRow(ctx, "SELECT ARRAY['sad', 'ok']::moods").Scan(&moods))
	assert.Equal(t, []string{"sad", "ok"}, moods)
}

func TestRegisterTypesReportsAMissingType(t *testing.T) {
	_, connect := testSchema(t)

	err := RegisterTypes(context.Background(), connect(), "mood", "no_such_type")

	assert.ErrorContains(t, err, "type no_such_type does not exist")
}
