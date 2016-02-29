package dataaccess

import (
	"fmt"
	"path"
	"strconv"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type dataAccess struct {
	gen *generator.Generator
}

func init() {
	generator.RegisterPlugin(new(dataAccess))
}

func (d *dataAccess) Name() string {
	return "dataaccess"
}

func (d *dataAccess) Init(g *generator.Generator) {
	d.gen = g
}

func (d *dataAccess) Generate(file *generator.FileDescriptor) {
	d.gen.P(`var _ raw.Store`)
	for _, m := range file.Messages() {
		if m.Options == nil {
			continue
		}

		v, err := proto.GetExtension(m.Options, E_Model)
		if err != nil {
			continue
		}

		t, ok := v.(*bool)
		if !ok || !(*t) {
			continue
		}

		d.gen.P()
		d.gen.P(`type `, m.GetName(), `BitMap uint32`)
		d.gen.P()
		d.gen.P(`const (`)
		for _, f := range m.Field {
			d.gen.P("\t", m.GetName(), f.GetName(), "Bit ", m.GetName(), "BitMap = 1 << ", fmt.Sprintf("%d", f.GetNumber()-1))
		}
		d.gen.P(`)`)
		d.gen.P()

		for _, f := range m.Field {
			if f.Options == nil {
				continue
			}

			v, err := proto.GetExtension(f.Options, E_Primarykey)
			if err != nil {
				continue
			}

			t, ok := v.(*bool)
			if !ok || !(*t) {
				continue
			}

			typename, _ := d.gen.GoType(m, f)

			fieldName := f.GetName()
			if gogoproto.IsCustomName(f) {
				fieldName = gogoproto.GetCustomName(f)
			}

			d.gen.P(`func encode`, m.GetName(), f.GetName(), `KeyString(k *`, m.GetName(), `) string {`)
			d.gen.P("\treturn \"/", m.GetName(), "/\" + proto.CompactTextString(k)")
			d.gen.P(`}`)
			d.gen.P()

			d.gen.P(`func encode`, m.GetName(), `Prefix() string {`)
			d.gen.P("\treturn \"/", m.GetName(), "\"")
			d.gen.P(`}`)
			d.gen.P()

			d.gen.P(`func Create`, m.GetName(), `(store raw.Store, `, f.GetName(), ` `, typename, `, o *`, m.GetName(), `) error {`)
			d.gen.P("\tkey := encode", m.GetName(), f.GetName(), `KeyString(&`, m.GetName(), `{`, fieldName, `:`, f.GetName(), `})`)
			d.gen.P("\treturn store.CreateObject(key, o)")
			d.gen.P(`}`)
			d.gen.P()

			d.gen.P(`func Delete`, m.GetName(), `(store raw.Store, `, f.GetName(), ` `, typename, `) error {`)
			d.gen.P("\tkey := encode", m.GetName(), f.GetName(), `KeyString(&`, m.GetName(), `{`, fieldName, `:`, f.GetName(), `})`)
			d.gen.P("\treturn store.DeleteObject(key)")
			d.gen.P(`}`)
			d.gen.P()

			d.gen.P(`func Get`, m.GetName(), `(store raw.Store, `, f.GetName(), ` `, typename, `) *`, m.GetName(), ` {`)
			d.gen.P("\tkey := encode", m.GetName(), f.GetName(), `KeyString(&`, m.GetName(), `{`, fieldName, `:`, f.GetName(), `})`)
			d.gen.P("\treturn store.Object(key).(*", m.GetName(), ")")
			d.gen.P(`}`)
			d.gen.P()

			d.gen.P(`func List`, m.GetName(), `(store raw.Store, `, f.GetName(), ` `, typename, `) []*`, m.GetName(), ` {`)
			d.gen.P("\tkey := encode", m.GetName(), `Prefix()`)
			d.gen.P("\tolist := store.ListObjects(key)")
			d.gen.P("\tlist := make([]*", m.GetName(), ", 0, len(olist))")
			d.gen.P("\tfor _,o := range olist {")
			d.gen.P("\t\tlist = append(list, o.(*", m.GetName(), "))")
			d.gen.P("\t}")
			d.gen.P(`return list`)
			d.gen.P(`}`)
			d.gen.P()

		}
	}
}

func (d *dataAccess) GenerateImports(file *generator.FileDescriptor) {
	d.gen.P("import (")
	d.gen.P(strconv.Quote(path.Join(d.gen.ImportPrefix, "github.com/docker/swarm-v2/state/raw")))
	d.gen.P(")")
	d.gen.P()
}
