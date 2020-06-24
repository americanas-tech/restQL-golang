package ast

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func newQuery(uses, firstBlock, otherBlocks interface{}) (Query, error) {
	var q Query

	useList := uses.([]interface{})
	if len(useList) > 0 {
		us := make([]Use, len(useList))
		for i, u := range useList {
			us[i] = u.(Use)
		}

		q.Use = us
	}

	fb := firstBlock.(Block)
	blocks := []interface{}{fb}

	if otherBlocks != nil {
		otherBs := otherBlocks.([]interface{})
		otherBs = flatten(otherBs)

		blocks = append(blocks, otherBs...)
	}

	q.Blocks = newBlockList(blocks)

	return q, nil
}

func newBlockList(blocks []interface{}) []Block {
	var result []Block

	for _, b := range blocks {
		switch b := b.(type) {
		case Block:
			result = append(result, b)
		case []interface{}:
			bs := newBlockList(b)
			result = append(result, bs...)
		default:
			continue
		}
	}

	return result
}

func newUse(rule, value interface{}) (Use, error) {
	r := rule.(string)
	v := value.(UseValue)

	return Use{Key: r, Value: v}, nil
}

func newUseValue(value interface{}) (UseValue, error) {
	vInt, ok := value.(int)
	if ok {
		return UseValue{Int: &vInt}, nil
	}

	sInt, ok := value.(string)
	if ok {
		return UseValue{String: &sInt}, nil
	}

	return UseValue{}, errors.Errorf("unknown use value type : %T", value)
}

func newBlock(action, modifiers, with, filter, ignore interface{}) (Block, error) {
	ac := action.(actionRule)
	block := Block{
		Method:   ac.Method,
		Resource: ac.Resource,
		Alias:    ac.Alias,
		In:       ac.In,
	}

	if modifiers != nil {
		ms := modifiers.([]interface{})

		for _, m := range ms {
			var q Qualifier

			switch m := m.(type) {
			case []HeaderItem:
				q = Qualifier{Headers: m}
			case *TimeoutValue:
				q = Qualifier{Timeout: m}
			case *MaxAgeValue:
				q = Qualifier{MaxAge: m}
			case *SMaxAgeValue:
				q = Qualifier{SMaxAge: m}
			default:
				continue
			}

			block.Qualifiers = append(block.Qualifiers, q)
		}
	}

	if with != nil {
		w := with.(*Parameters)
		q := Qualifier{With: w}

		block.Qualifiers = append(block.Qualifiers, q)
	}

	if filter != nil {
		var q Qualifier

		switch filter := filter.(type) {
		case []Filter:
			q = Qualifier{Only: filter}
		case hidden:
			q = Qualifier{Hidden: true}
		default:
			return Block{}, fmt.Errorf("got an unknown value of type %T", filter)
		}

		block.Qualifiers = append(block.Qualifiers, q)
	}

	if ignore != nil {
		ig := ignore.(ignoreErrors)
		q := Qualifier{IgnoreErrors: bool(ig)}

		block.Qualifiers = append(block.Qualifiers, q)
	}

	return block, nil
}

type actionRule struct {
	Method   string
	Resource string
	Alias    string
	In       []string
}

func newActionRule(method, resource, alias, in interface{}) (actionRule, error) {
	m := method.(string)
	r := resource.(string)

	ar := actionRule{Method: m, Resource: r}

	if alias != nil {
		a := alias.(string)
		ar.Alias = a
	}

	if in != nil {
		i := in.([]string)
		ar.In = i
	}

	return ar, nil
}

func newIn(target interface{}) ([]string, error) {
	t := target.(string)
	path := strings.Split(t, ".")
	return path, nil
}

func newWith(parameterBody, keyValues interface{}) (*Parameters, error) {
	if parameterBody == nil && keyValues == nil {
		return nil, errors.New("empty with clause is not allowed")
	}

	var p Parameters

	if parameterBody != nil {
		pb := parameterBody.(*ParameterBody)
		p.Body = pb
	}

	if keyValues != nil {
		kvs := keyValues.([]KeyValue)
		if len(kvs) > 0 {
			p.KeyValues = kvs
		}
	}

	return &p, nil
}

func newParameterBody(target, functions interface{}) (*ParameterBody, error) {
	t := target.(string)
	pb := ParameterBody{Target: t}

	if functions != nil {
		fns := newFunctionList(functions)
		pb.Functions = fns
	}

	return &pb, nil
}

func newKeyValueList(first, others interface{}) ([]KeyValue, error) {
	kv := first.(KeyValue)
	kvs := []interface{}{kv}

	if others != nil {
		otherKvs := others.([]interface{})
		if len(otherKvs) > 0 {
			otherKvs = flatten(otherKvs)
			kvs = append(kvs, otherKvs...)
		}
	}

	return castKeyValueList(kvs), nil
}

func castKeyValueList(kvs []interface{}) []KeyValue {
	var result []KeyValue

	for _, kv := range kvs {
		switch kv := kv.(type) {
		case KeyValue:
			result = append(result, kv)
		case []interface{}:
			result = append(result, castKeyValueList(kv)...)
		default:
			continue
		}
	}

	return result
}

func newKeyValue(key, value, functions interface{}) (KeyValue, error) {
	k := key.(string)
	v := value.(Value)

	kv := KeyValue{Key: k, Value: v}

	if functions != nil {
		kv.Functions = newFunctionList(functions)
	}

	return kv, nil
}

func newFunctionList(functions interface{}) []string {
	fns := functions.([]interface{})
	var result []string

	for _, fn := range fns {
		if fn, ok := fn.(string); ok {
			result = append(result, fn)
		}
	}

	return result
}

func newValue(value interface{}) (Value, error) {
	switch value := value.(type) {
	case *Primitive:
		return Value{Primitive: value}, nil
	case variable:
		v := string(value)
		return Value{Variable: &v}, nil
	case []Value:
		return Value{List: value}, nil
	case []ObjectEntry:
		return Value{Object: value}, nil
	default:
		return Value{}, fmt.Errorf("got an unknown value of type %T", value)
	}
}

func newEmptyList() ([]Value, error) {
	return []Value{}, nil
}

func newList(first, others interface{}) ([]Value, error) {
	fi := first.(Value)
	list := []Value{fi}

	if others != nil {
		oi := others.([]interface{})
		oi = flatten(oi)

		for _, i := range oi {
			if i, ok := i.(Value); ok {
				list = append(list, i)
			}
		}
	}

	return list, nil
}

func newEmptyObject() ([]ObjectEntry, error) {
	return []ObjectEntry{}, nil
}

func newPopulatedObject(first, others interface{}) ([]ObjectEntry, error) {
	fe := first.(ObjectEntry)
	entries := []ObjectEntry{fe}

	if others != nil {
		oe := others.([]interface{})
		oe = flatten(oe)

		for _, e := range oe {
			if e, ok := e.(ObjectEntry); ok {
				entries = append(entries, e)
			}
		}
	}

	return entries, nil
}

func newObjectEntry(key, value interface{}) (ObjectEntry, error) {
	k := key.(string)
	v := value.(Value)

	return ObjectEntry{Key: k, Value: v}, nil
}

func newChain(first, others interface{}) ([]Chained, error) {
	fc := first.(Chained)
	chain := []Chained{fc}

	if others != nil {
		oc := others.([]interface{})
		if len(oc) > 0 {
			oc = flatten(oc)

			var result []Chained
			for _, c := range oc {
				if c, ok := c.(Chained); ok {
					result = append(result, c)
				}
			}

			chain = append(chain, result...)
		}
	}

	return chain, nil
}

func newChained(chainItem interface{}) (Chained, error) {
	switch chainItem := chainItem.(type) {
	case variable:
		i := string(chainItem)
		return Chained{PathVariable: i}, nil
	case string:
		return Chained{PathItem: chainItem}, nil
	default:
		return Chained{}, fmt.Errorf("got an unknown type : %T", chainItem)
	}
}

type variable string

func newChainPathVariable(pathVariable interface{}) (variable, error) {
	return newVariable(pathVariable)
}

func newVariable(v interface{}) (variable, error) {
	vStr := v.(string)
	return variable(vStr), nil
}

func newPrimitive(value interface{}) (*Primitive, error) {
	var p Primitive

	switch value := value.(type) {
	case int:
		p.Int = &value
	case string:
		p.String = &value
	case float64:
		p.Float = &value
	case bool:
		p.Boolean = &value
	case []Chained:
		p.Chain = value
	case null:
		p.Null = true
	}

	return &p, nil
}

func newOnly(first, others interface{}) ([]Filter, error) {
	ff := first.(Filter)
	filters := []Filter{ff}

	if others != nil {
		fs := others.([]interface{})
		if len(fs) > 0 {
			fs = flatten(fs)

			for _, f := range fs {
				if f, ok := f.(Filter); ok {
					filters = append(filters, f)
				}
			}

		}
	}

	return filters, nil
}

func newFilter(identifier, matchArg interface{}) (Filter, error) {
	ident := identifier.(string)
	fields := strings.Split(ident, ".")
	filter := Filter{Field: fields}

	if matchArg != nil {
		ma := matchArg.(string)
		filter.Match = ma
	}

	return filter, nil
}

func newFilterValue(value interface{}) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil
	case []byte:
		return stringify(value)
	default:
		return "", fmt.Errorf("got an unknown type : %T", value)
	}
}

func newMatchesFunction(arg interface{}) (string, error) {
	a := arg.(string)
	return a, nil
}

type hidden bool

func newHidden() (hidden, error) {
	return true, nil
}

func newHeaders(first, others interface{}) ([]HeaderItem, error) {
	fh := first.(HeaderItem)
	headers := []HeaderItem{fh}

	if others != nil {
		hs := others.([]interface{})
		if len(hs) > 0 {
			hs = flatten(hs)

			for _, h := range hs {
				if h, ok := h.(HeaderItem); ok {
					headers = append(headers, h)
				}
			}

		}
	}

	return headers, nil
}

func newHeader(name, value interface{}) (HeaderItem, error) {
	n := name.(string)
	v, err := newHeaderValue(value)
	if err != nil {
		return HeaderItem{}, err
	}

	return HeaderItem{Key: n, Value: v}, nil
}

func newHeaderValue(value interface{}) (HeaderValue, error) {
	switch value := value.(type) {
	case variable:
		v := string(value)
		return HeaderValue{Variable: &v}, nil
	case string:
		return HeaderValue{String: &value}, nil
	default:
		return HeaderValue{}, fmt.Errorf("got an unknown type : %T", value)
	}
}

func newTimeout(value interface{}) (*TimeoutValue, error) {
	switch value := value.(type) {
	case variable:
		v := string(value)
		return &TimeoutValue{Variable: &v}, nil
	case int:
		return &TimeoutValue{Int: &value}, nil
	default:
		return &TimeoutValue{}, fmt.Errorf("got an unknown type : %T", value)
	}
}

func newMaxAge(value interface{}) (*MaxAgeValue, error) {
	switch value := value.(type) {
	case variable:
		v := string(value)
		return &MaxAgeValue{Variable: &v}, nil
	case int:
		return &MaxAgeValue{Int: &value}, nil
	default:
		return &MaxAgeValue{}, fmt.Errorf("got an unknown type : %T", value)
	}
}

func newSmaxAge(value interface{}) (*SMaxAgeValue, error) {
	switch value := value.(type) {
	case variable:
		v := string(value)
		return &SMaxAgeValue{Variable: &v}, nil
	case int:
		return &SMaxAgeValue{Int: &value}, nil
	default:
		return &SMaxAgeValue{}, fmt.Errorf("got an unknown type : %T", value)
	}
}

type ignoreErrors bool

func newFlags(ignoreFlag, others interface{}) (ignoreErrors, error) {
	i := ignoreFlag.(ignoreErrors)
	return i, nil
}

func newIgnoreErrors() (ignoreErrors, error) {
	return true, nil
}

func newBoolean(boolean []byte) (bool, error) {
	return strconv.ParseBool(string(boolean))
}

func newString(str []byte) (string, error) {
	return strconv.Unquote(string(str))
}

func newFloat(float []byte) (float64, error) {
	return strconv.ParseFloat(string(float), 64)
}

func newInteger(integer []byte) (int, error) {
	i, err := strconv.ParseInt(string(integer), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

type null struct{}

func newNull() (null, error) {
	return struct{}{}, nil
}

func stringify(s []byte) (string, error) {
	return string(s), nil
}

func flatten(ii []interface{}) []interface{} {
	var res []interface{}
	for _, i := range ii {
		if i == nil {
			continue
		}
		switch t := i.(type) {
		case []interface{}:
			res = append(res, flatten(t)...)
		default:
			res = append(res, i)
		}
	}

	return res
}
