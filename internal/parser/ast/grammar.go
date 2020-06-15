
package ast

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar {
	rules: []*rule{
{
	name: "QUERY",
	pos: position{line: 17, col: 1, offset: 118},
	expr: &actionExpr{
	pos: position{line: 17, col: 10, offset: 127},
	run: (*parser).callonQUERY1,
	expr: &seqExpr{
	pos: position{line: 17, col: 10, offset: 127},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 10, offset: 127},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 13, offset: 130},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 13, offset: 130},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 17, offset: 134},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 17, col: 20, offset: 137},
	label: "us",
	expr: &zeroOrMoreExpr{
	pos: position{line: 17, col: 23, offset: 140},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 24, offset: 141},
	name: "USE",
},
},
},
&labeledExpr{
	pos: position{line: 17, col: 30, offset: 147},
	label: "firstBlock",
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 41, offset: 158},
	name: "BLOCK",
},
},
&labeledExpr{
	pos: position{line: 17, col: 47, offset: 164},
	label: "otherBlocks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 17, col: 59, offset: 176},
	expr: &seqExpr{
	pos: position{line: 17, col: 60, offset: 177},
	exprs: []interface{}{
&oneOrMoreExpr{
	pos: position{line: 17, col: 60, offset: 177},
	expr: &seqExpr{
	pos: position{line: 17, col: 61, offset: 178},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 61, offset: 178},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 64, offset: 181},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 67, offset: 184},
	name: "WS",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 72, offset: 189},
	name: "BLOCK",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 80, offset: 197},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 83, offset: 200},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 83, offset: 200},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 87, offset: 204},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 90, offset: 207},
	name: "EOF",
},
	},
},
},
},
{
	name: "USE",
	pos: position{line: 21, col: 1, offset: 262},
	expr: &actionExpr{
	pos: position{line: 21, col: 8, offset: 269},
	run: (*parser).callonUSE1,
	expr: &seqExpr{
	pos: position{line: 21, col: 8, offset: 269},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 21, col: 8, offset: 269},
	val: "use",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 21, col: 14, offset: 275},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 21, col: 22, offset: 283},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 25, offset: 286},
	name: "USE_ACTION",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 37, offset: 298},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 21, col: 40, offset: 301},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 43, offset: 304},
	name: "USE_VALUE",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 54, offset: 315},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 21, col: 57, offset: 318},
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 57, offset: 318},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 61, offset: 322},
	name: "WS",
},
	},
},
},
},
{
	name: "USE_ACTION",
	pos: position{line: 25, col: 1, offset: 351},
	expr: &actionExpr{
	pos: position{line: 25, col: 15, offset: 365},
	run: (*parser).callonUSE_ACTION1,
	expr: &choiceExpr{
	pos: position{line: 25, col: 16, offset: 366},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 16, offset: 366},
	val: "timeout",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 28, offset: 378},
	val: "max-age",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 40, offset: 390},
	val: "s-max-age",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "USE_VALUE",
	pos: position{line: 29, col: 1, offset: 434},
	expr: &actionExpr{
	pos: position{line: 29, col: 14, offset: 447},
	run: (*parser).callonUSE_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 29, col: 14, offset: 447},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 29, col: 17, offset: 450},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 29, col: 17, offset: 450},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 29, col: 26, offset: 459},
	name: "Integer",
},
	},
},
},
},
},
{
	name: "BLOCK",
	pos: position{line: 33, col: 1, offset: 496},
	expr: &actionExpr{
	pos: position{line: 33, col: 10, offset: 505},
	run: (*parser).callonBLOCK1,
	expr: &seqExpr{
	pos: position{line: 33, col: 10, offset: 505},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 33, col: 10, offset: 505},
	label: "action",
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 18, offset: 513},
	name: "ACTION_RULE",
},
},
&labeledExpr{
	pos: position{line: 33, col: 31, offset: 526},
	label: "m",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 34, offset: 529},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 34, offset: 529},
	name: "MODIFIER_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 50, offset: 545},
	label: "w",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 53, offset: 548},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 53, offset: 548},
	name: "WITH_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 65, offset: 560},
	label: "f",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 67, offset: 562},
	expr: &choiceExpr{
	pos: position{line: 33, col: 68, offset: 563},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 33, col: 68, offset: 563},
	name: "ONLY_RULE",
},
&ruleRefExpr{
	pos: position{line: 33, col: 80, offset: 575},
	name: "HIDDEN_RULE",
},
	},
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 94, offset: 589},
	label: "fl",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 98, offset: 593},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 98, offset: 593},
	name: "FLAG_RULE",
},
},
},
&ruleRefExpr{
	pos: position{line: 33, col: 110, offset: 605},
	name: "WS",
},
	},
},
},
},
{
	name: "ACTION_RULE",
	pos: position{line: 37, col: 1, offset: 651},
	expr: &actionExpr{
	pos: position{line: 37, col: 16, offset: 666},
	run: (*parser).callonACTION_RULE1,
	expr: &seqExpr{
	pos: position{line: 37, col: 16, offset: 666},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 37, col: 16, offset: 666},
	label: "m",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 19, offset: 669},
	name: "METHOD",
},
},
&ruleRefExpr{
	pos: position{line: 37, col: 27, offset: 677},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 37, col: 35, offset: 685},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 38, offset: 688},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 37, col: 45, offset: 695},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 48, offset: 698},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 48, offset: 698},
	name: "ALIAS",
},
},
},
&labeledExpr{
	pos: position{line: 37, col: 56, offset: 706},
	label: "i",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 59, offset: 709},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 59, offset: 709},
	name: "IN",
},
},
},
	},
},
},
},
{
	name: "METHOD",
	pos: position{line: 41, col: 1, offset: 753},
	expr: &actionExpr{
	pos: position{line: 41, col: 11, offset: 763},
	run: (*parser).callonMETHOD1,
	expr: &labeledExpr{
	pos: position{line: 41, col: 11, offset: 763},
	label: "m",
	expr: &choiceExpr{
	pos: position{line: 41, col: 14, offset: 766},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 41, col: 14, offset: 766},
	val: "from",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 23, offset: 775},
	val: "to",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 30, offset: 782},
	val: "into",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 38, offset: 790},
	val: "update",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 49, offset: 801},
	val: "delete",
	ignoreCase: false,
},
	},
},
},
},
},
{
	name: "ALIAS",
	pos: position{line: 45, col: 1, offset: 842},
	expr: &actionExpr{
	pos: position{line: 45, col: 10, offset: 851},
	run: (*parser).callonALIAS1,
	expr: &seqExpr{
	pos: position{line: 45, col: 10, offset: 851},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 45, col: 10, offset: 851},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 45, col: 18, offset: 859},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 45, col: 23, offset: 864},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 45, col: 31, offset: 872},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 45, col: 34, offset: 875},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IN",
	pos: position{line: 49, col: 1, offset: 902},
	expr: &actionExpr{
	pos: position{line: 49, col: 7, offset: 908},
	run: (*parser).callonIN1,
	expr: &seqExpr{
	pos: position{line: 49, col: 7, offset: 908},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 7, offset: 908},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 49, col: 15, offset: 916},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 49, col: 20, offset: 921},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 49, col: 28, offset: 929},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 49, col: 31, offset: 932},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "MODIFIER_RULE",
	pos: position{line: 53, col: 1, offset: 970},
	expr: &actionExpr{
	pos: position{line: 53, col: 18, offset: 987},
	run: (*parser).callonMODIFIER_RULE1,
	expr: &labeledExpr{
	pos: position{line: 53, col: 18, offset: 987},
	label: "m",
	expr: &oneOrMoreExpr{
	pos: position{line: 53, col: 20, offset: 989},
	expr: &choiceExpr{
	pos: position{line: 53, col: 21, offset: 990},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 53, col: 21, offset: 990},
	name: "HEADERS",
},
&ruleRefExpr{
	pos: position{line: 53, col: 31, offset: 1000},
	name: "TIMEOUT",
},
&ruleRefExpr{
	pos: position{line: 53, col: 41, offset: 1010},
	name: "MAX_AGE",
},
&ruleRefExpr{
	pos: position{line: 53, col: 51, offset: 1020},
	name: "S_MAX_AGE",
},
	},
},
},
},
},
},
{
	name: "WITH_RULE",
	pos: position{line: 57, col: 1, offset: 1052},
	expr: &actionExpr{
	pos: position{line: 57, col: 14, offset: 1065},
	run: (*parser).callonWITH_RULE1,
	expr: &seqExpr{
	pos: position{line: 57, col: 14, offset: 1065},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 57, col: 14, offset: 1065},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 57, col: 22, offset: 1073},
	val: "with",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 57, col: 29, offset: 1080},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 57, col: 37, offset: 1088},
	label: "pb",
	expr: &zeroOrOneExpr{
	pos: position{line: 57, col: 40, offset: 1091},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 40, offset: 1091},
	name: "PARAMETER_BODY",
},
},
},
&labeledExpr{
	pos: position{line: 57, col: 56, offset: 1107},
	label: "kvs",
	expr: &zeroOrOneExpr{
	pos: position{line: 57, col: 60, offset: 1111},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 60, offset: 1111},
	name: "KEY_VALUE_LIST",
},
},
},
	},
},
},
},
{
	name: "PARAMETER_BODY",
	pos: position{line: 61, col: 1, offset: 1157},
	expr: &actionExpr{
	pos: position{line: 61, col: 19, offset: 1175},
	run: (*parser).callonPARAMETER_BODY1,
	expr: &seqExpr{
	pos: position{line: 61, col: 19, offset: 1175},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 61, col: 19, offset: 1175},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 61, col: 23, offset: 1179},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 26, offset: 1182},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 61, col: 33, offset: 1189},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 61, col: 37, offset: 1193},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 37, offset: 1193},
	name: "APPLY_FN",
},
},
},
&ruleRefExpr{
	pos: position{line: 61, col: 48, offset: 1204},
	name: "WS",
},
&zeroOrOneExpr{
	pos: position{line: 61, col: 51, offset: 1207},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 51, offset: 1207},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 61, col: 55, offset: 1211},
	name: "WS",
},
	},
},
},
},
{
	name: "KEY_VALUE_LIST",
	pos: position{line: 65, col: 1, offset: 1251},
	expr: &actionExpr{
	pos: position{line: 65, col: 19, offset: 1269},
	run: (*parser).callonKEY_VALUE_LIST1,
	expr: &seqExpr{
	pos: position{line: 65, col: 19, offset: 1269},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 65, col: 19, offset: 1269},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 25, offset: 1275},
	name: "KEY_VALUE",
},
},
&labeledExpr{
	pos: position{line: 65, col: 35, offset: 1285},
	label: "others",
	expr: &zeroOrMoreExpr{
	pos: position{line: 65, col: 42, offset: 1292},
	expr: &seqExpr{
	pos: position{line: 65, col: 43, offset: 1293},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 43, offset: 1293},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 46, offset: 1296},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 49, offset: 1299},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 52, offset: 1302},
	name: "KEY_VALUE",
},
	},
},
},
},
	},
},
},
},
{
	name: "KEY_VALUE",
	pos: position{line: 69, col: 1, offset: 1358},
	expr: &actionExpr{
	pos: position{line: 69, col: 14, offset: 1371},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 69, col: 14, offset: 1371},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 69, col: 14, offset: 1371},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 17, offset: 1374},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 69, col: 33, offset: 1390},
	name: "WS",
},
&litMatcher{
	pos: position{line: 69, col: 36, offset: 1393},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 69, col: 40, offset: 1397},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 69, col: 43, offset: 1400},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 46, offset: 1403},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 69, col: 53, offset: 1410},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 69, col: 57, offset: 1414},
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 57, offset: 1414},
	name: "APPLY_FN",
},
},
},
	},
},
},
},
{
	name: "APPLY_FN",
	pos: position{line: 73, col: 1, offset: 1460},
	expr: &actionExpr{
	pos: position{line: 73, col: 13, offset: 1472},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 73, col: 13, offset: 1472},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 73, col: 13, offset: 1472},
	name: "WS",
},
&litMatcher{
	pos: position{line: 73, col: 16, offset: 1475},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 73, col: 21, offset: 1480},
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 21, offset: 1480},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 73, col: 25, offset: 1484},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 29, offset: 1488},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 77, col: 1, offset: 1519},
	expr: &actionExpr{
	pos: position{line: 77, col: 13, offset: 1531},
	run: (*parser).callonFUNCTION1,
	expr: &choiceExpr{
	pos: position{line: 77, col: 14, offset: 1532},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 77, col: 14, offset: 1532},
	val: "flatten",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 26, offset: 1544},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 37, offset: 1555},
	val: "json",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VALUE",
	pos: position{line: 81, col: 1, offset: 1594},
	expr: &actionExpr{
	pos: position{line: 81, col: 10, offset: 1603},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 81, col: 10, offset: 1603},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 81, col: 13, offset: 1606},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 81, col: 13, offset: 1606},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 81, col: 20, offset: 1613},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 81, col: 29, offset: 1622},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 81, col: 40, offset: 1633},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 85, col: 1, offset: 1669},
	expr: &actionExpr{
	pos: position{line: 85, col: 9, offset: 1677},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 85, col: 9, offset: 1677},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 85, col: 12, offset: 1680},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 85, col: 12, offset: 1680},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 85, col: 25, offset: 1693},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 89, col: 1, offset: 1729},
	expr: &actionExpr{
	pos: position{line: 89, col: 15, offset: 1743},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 89, col: 15, offset: 1743},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 15, offset: 1743},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 89, col: 19, offset: 1747},
	name: "WS",
},
&litMatcher{
	pos: position{line: 89, col: 22, offset: 1750},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 93, col: 1, offset: 1782},
	expr: &actionExpr{
	pos: position{line: 93, col: 19, offset: 1800},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 93, col: 19, offset: 1800},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 93, col: 19, offset: 1800},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 93, col: 23, offset: 1804},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 93, col: 26, offset: 1807},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 28, offset: 1809},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 93, col: 34, offset: 1815},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 93, col: 37, offset: 1818},
	expr: &seqExpr{
	pos: position{line: 93, col: 38, offset: 1819},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 93, col: 38, offset: 1819},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 93, col: 41, offset: 1822},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 93, col: 44, offset: 1825},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 93, col: 47, offset: 1828},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 55, offset: 1836},
	name: "WS",
},
&litMatcher{
	pos: position{line: 93, col: 58, offset: 1839},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 97, col: 1, offset: 1871},
	expr: &actionExpr{
	pos: position{line: 97, col: 11, offset: 1881},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 97, col: 11, offset: 1881},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 97, col: 14, offset: 1884},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 97, col: 14, offset: 1884},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 97, col: 26, offset: 1896},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 101, col: 1, offset: 1931},
	expr: &actionExpr{
	pos: position{line: 101, col: 14, offset: 1944},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 101, col: 14, offset: 1944},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 101, col: 14, offset: 1944},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 101, col: 18, offset: 1948},
	name: "WS",
},
&litMatcher{
	pos: position{line: 101, col: 21, offset: 1951},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 105, col: 1, offset: 1985},
	expr: &actionExpr{
	pos: position{line: 105, col: 18, offset: 2002},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 105, col: 18, offset: 2002},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 105, col: 18, offset: 2002},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 22, offset: 2006},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 105, col: 25, offset: 2009},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 29, offset: 2013},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 105, col: 40, offset: 2024},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 105, col: 44, offset: 2028},
	expr: &seqExpr{
	pos: position{line: 105, col: 45, offset: 2029},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 105, col: 45, offset: 2029},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 48, offset: 2032},
	val: ",",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 52, offset: 2036},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 52, offset: 2036},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 56, offset: 2040},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 105, col: 59, offset: 2043},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 71, offset: 2055},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 74, offset: 2058},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 109, col: 1, offset: 2103},
	expr: &actionExpr{
	pos: position{line: 109, col: 14, offset: 2116},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 109, col: 14, offset: 2116},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 109, col: 14, offset: 2116},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 109, col: 17, offset: 2119},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 109, col: 17, offset: 2119},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 109, col: 26, offset: 2128},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 33, offset: 2135},
	name: "WS",
},
&litMatcher{
	pos: position{line: 109, col: 36, offset: 2138},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 109, col: 40, offset: 2142},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 109, col: 43, offset: 2145},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 46, offset: 2148},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 113, col: 1, offset: 2189},
	expr: &actionExpr{
	pos: position{line: 113, col: 14, offset: 2202},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 113, col: 14, offset: 2202},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 113, col: 17, offset: 2205},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2205},
	name: "Null",
},
&ruleRefExpr{
	pos: position{line: 113, col: 24, offset: 2212},
	name: "Boolean",
},
&ruleRefExpr{
	pos: position{line: 113, col: 34, offset: 2222},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 113, col: 43, offset: 2231},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 113, col: 51, offset: 2239},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 113, col: 61, offset: 2249},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "CHAIN",
	pos: position{line: 117, col: 1, offset: 2285},
	expr: &actionExpr{
	pos: position{line: 117, col: 10, offset: 2294},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 117, col: 10, offset: 2294},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 117, col: 10, offset: 2294},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 117, col: 13, offset: 2297},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 117, col: 27, offset: 2311},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 117, col: 30, offset: 2314},
	expr: &seqExpr{
	pos: position{line: 117, col: 31, offset: 2315},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 117, col: 31, offset: 2315},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 117, col: 35, offset: 2319},
	name: "CHAINED_ITEM",
},
	},
},
},
},
	},
},
},
},
{
	name: "CHAINED_ITEM",
	pos: position{line: 121, col: 1, offset: 2363},
	expr: &actionExpr{
	pos: position{line: 121, col: 17, offset: 2379},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 121, col: 17, offset: 2379},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 121, col: 21, offset: 2383},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 121, col: 21, offset: 2383},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 121, col: 32, offset: 2394},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 125, col: 1, offset: 2429},
	expr: &actionExpr{
	pos: position{line: 125, col: 14, offset: 2442},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 125, col: 14, offset: 2442},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 14, offset: 2442},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 125, col: 22, offset: 2450},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 125, col: 29, offset: 2457},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 125, col: 37, offset: 2465},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 125, col: 40, offset: 2468},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 125, col: 48, offset: 2476},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 125, col: 51, offset: 2479},
	expr: &seqExpr{
	pos: position{line: 125, col: 52, offset: 2480},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 52, offset: 2480},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 55, offset: 2483},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 58, offset: 2486},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 61, offset: 2489},
	name: "FILTER",
},
	},
},
},
},
	},
},
},
},
{
	name: "FILTER",
	pos: position{line: 129, col: 1, offset: 2526},
	expr: &actionExpr{
	pos: position{line: 129, col: 11, offset: 2536},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 129, col: 11, offset: 2536},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 129, col: 11, offset: 2536},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 14, offset: 2539},
	name: "IDENT_WITH_DOT",
},
},
&labeledExpr{
	pos: position{line: 129, col: 30, offset: 2555},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 129, col: 34, offset: 2559},
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 34, offset: 2559},
	name: "MATCHES_FN",
},
},
},
	},
},
},
},
{
	name: "MATCHES_FN",
	pos: position{line: 133, col: 1, offset: 2602},
	expr: &actionExpr{
	pos: position{line: 133, col: 15, offset: 2616},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 133, col: 15, offset: 2616},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 133, col: 15, offset: 2616},
	name: "WS",
},
&litMatcher{
	pos: position{line: 133, col: 18, offset: 2619},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 133, col: 23, offset: 2624},
	name: "WS",
},
&litMatcher{
	pos: position{line: 133, col: 26, offset: 2627},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 133, col: 36, offset: 2637},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 133, col: 40, offset: 2641},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 133, col: 45, offset: 2646},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 133, col: 53, offset: 2654},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 137, col: 1, offset: 2695},
	expr: &actionExpr{
	pos: position{line: 137, col: 12, offset: 2706},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 137, col: 12, offset: 2706},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 12, offset: 2706},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 137, col: 20, offset: 2714},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 137, col: 30, offset: 2724},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 137, col: 38, offset: 2732},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 137, col: 41, offset: 2735},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 137, col: 49, offset: 2743},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 137, col: 52, offset: 2746},
	expr: &seqExpr{
	pos: position{line: 137, col: 53, offset: 2747},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 53, offset: 2747},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 137, col: 56, offset: 2750},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 137, col: 59, offset: 2753},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 137, col: 62, offset: 2756},
	name: "HEADER",
},
	},
},
},
},
	},
},
},
},
{
	name: "HEADER",
	pos: position{line: 141, col: 1, offset: 2796},
	expr: &actionExpr{
	pos: position{line: 141, col: 11, offset: 2806},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 141, col: 11, offset: 2806},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 141, col: 11, offset: 2806},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 141, col: 14, offset: 2809},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 141, col: 21, offset: 2816},
	name: "WS",
},
&litMatcher{
	pos: position{line: 141, col: 24, offset: 2819},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 141, col: 28, offset: 2823},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 141, col: 31, offset: 2826},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 141, col: 34, offset: 2829},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 141, col: 34, offset: 2829},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 141, col: 45, offset: 2840},
	name: "String",
},
	},
},
},
	},
},
},
},
{
	name: "HIDDEN_RULE",
	pos: position{line: 145, col: 1, offset: 2877},
	expr: &actionExpr{
	pos: position{line: 145, col: 16, offset: 2892},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 145, col: 16, offset: 2892},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 145, col: 16, offset: 2892},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 145, col: 24, offset: 2900},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 149, col: 1, offset: 2934},
	expr: &actionExpr{
	pos: position{line: 149, col: 12, offset: 2945},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 149, col: 12, offset: 2945},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 12, offset: 2945},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 149, col: 20, offset: 2953},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 149, col: 30, offset: 2963},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 149, col: 38, offset: 2971},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 149, col: 41, offset: 2974},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 41, offset: 2974},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 149, col: 52, offset: 2985},
	name: "Integer",
},
	},
},
},
	},
},
},
},
{
	name: "MAX_AGE",
	pos: position{line: 153, col: 1, offset: 3021},
	expr: &actionExpr{
	pos: position{line: 153, col: 12, offset: 3032},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 153, col: 12, offset: 3032},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 12, offset: 3032},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 153, col: 20, offset: 3040},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 153, col: 30, offset: 3050},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 153, col: 38, offset: 3058},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 153, col: 41, offset: 3061},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 41, offset: 3061},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 153, col: 52, offset: 3072},
	name: "Integer",
},
	},
},
},
	},
},
},
},
{
	name: "S_MAX_AGE",
	pos: position{line: 157, col: 1, offset: 3107},
	expr: &actionExpr{
	pos: position{line: 157, col: 14, offset: 3120},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 157, col: 14, offset: 3120},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 14, offset: 3120},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 157, col: 22, offset: 3128},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 34, offset: 3140},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 157, col: 42, offset: 3148},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 157, col: 45, offset: 3151},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 45, offset: 3151},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 157, col: 56, offset: 3162},
	name: "Integer",
},
	},
},
},
	},
},
},
},
{
	name: "FLAG_RULE",
	pos: position{line: 161, col: 1, offset: 3198},
	expr: &actionExpr{
	pos: position{line: 161, col: 14, offset: 3211},
	run: (*parser).callonFLAG_RULE1,
	expr: &seqExpr{
	pos: position{line: 161, col: 14, offset: 3211},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 14, offset: 3211},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 161, col: 22, offset: 3219},
	val: "ignore-errors",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 165, col: 1, offset: 3266},
	expr: &actionExpr{
	pos: position{line: 165, col: 13, offset: 3278},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 165, col: 13, offset: 3278},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 165, col: 13, offset: 3278},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 165, col: 17, offset: 3282},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 165, col: 20, offset: 3285},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 169, col: 1, offset: 3320},
	expr: &actionExpr{
	pos: position{line: 169, col: 10, offset: 3329},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 169, col: 10, offset: 3329},
	expr: &charClassMatcher{
	pos: position{line: 169, col: 10, offset: 3329},
	val: "[A-Za-z0-9_-]",
	chars: []rune{'_','-',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
{
	name: "IDENT_WITH_DOT",
	pos: position{line: 173, col: 1, offset: 3375},
	expr: &actionExpr{
	pos: position{line: 173, col: 19, offset: 3393},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 173, col: 19, offset: 3393},
	expr: &charClassMatcher{
	pos: position{line: 173, col: 19, offset: 3393},
	val: "[A-Za-z0-9-_.]",
	chars: []rune{'-','_','.',},
	ranges: []rune{'A','Z','a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
{
	name: "Null",
	pos: position{line: 177, col: 1, offset: 3440},
	expr: &actionExpr{
	pos: position{line: 177, col: 9, offset: 3448},
	run: (*parser).callonNull1,
	expr: &litMatcher{
	pos: position{line: 177, col: 9, offset: 3448},
	val: "null",
	ignoreCase: false,
},
},
},
{
	name: "Boolean",
	pos: position{line: 181, col: 1, offset: 3478},
	expr: &actionExpr{
	pos: position{line: 181, col: 12, offset: 3489},
	run: (*parser).callonBoolean1,
	expr: &choiceExpr{
	pos: position{line: 181, col: 13, offset: 3490},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 181, col: 13, offset: 3490},
	val: "true",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 181, col: 22, offset: 3499},
	val: "false",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "String",
	pos: position{line: 185, col: 1, offset: 3540},
	expr: &actionExpr{
	pos: position{line: 185, col: 11, offset: 3550},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 185, col: 11, offset: 3550},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 185, col: 11, offset: 3550},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 185, col: 15, offset: 3554},
	expr: &seqExpr{
	pos: position{line: 185, col: 17, offset: 3556},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 185, col: 17, offset: 3556},
	expr: &litMatcher{
	pos: position{line: 185, col: 18, offset: 3557},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 185, col: 22, offset: 3561,
},
	},
},
},
&litMatcher{
	pos: position{line: 185, col: 27, offset: 3566},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 189, col: 1, offset: 3601},
	expr: &actionExpr{
	pos: position{line: 189, col: 10, offset: 3610},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 189, col: 10, offset: 3610},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 189, col: 10, offset: 3610},
	expr: &choiceExpr{
	pos: position{line: 189, col: 11, offset: 3611},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 189, col: 11, offset: 3611},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 189, col: 17, offset: 3617},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 189, col: 23, offset: 3623},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 189, col: 31, offset: 3631},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 189, col: 35, offset: 3635},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 193, col: 1, offset: 3673},
	expr: &actionExpr{
	pos: position{line: 193, col: 12, offset: 3684},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 193, col: 12, offset: 3684},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 193, col: 12, offset: 3684},
	expr: &choiceExpr{
	pos: position{line: 193, col: 13, offset: 3685},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 193, col: 13, offset: 3685},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 193, col: 19, offset: 3691},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 193, col: 25, offset: 3697},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 197, col: 1, offset: 3737},
	expr: &choiceExpr{
	pos: position{line: 197, col: 11, offset: 3749},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 197, col: 11, offset: 3749},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 197, col: 17, offset: 3755},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 197, col: 17, offset: 3755},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 197, col: 37, offset: 3775},
	expr: &ruleRefExpr{
	pos: position{line: 197, col: 37, offset: 3775},
	name: "DecimalDigit",
},
},
	},
},
	},
},
},
{
	name: "DecimalDigit",
	pos: position{line: 199, col: 1, offset: 3790},
	expr: &charClassMatcher{
	pos: position{line: 199, col: 16, offset: 3807},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 200, col: 1, offset: 3813},
	expr: &charClassMatcher{
	pos: position{line: 200, col: 23, offset: 3837},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 202, col: 1, offset: 3844},
	expr: &charClassMatcher{
	pos: position{line: 202, col: 10, offset: 3853},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 203, col: 1, offset: 3859},
	expr: &charClassMatcher{
	pos: position{line: 203, col: 18, offset: 3876},
	val: "[\\n\\r]",
	chars: []rune{'\n','\r',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 205, col: 1, offset: 3884},
	expr: &choiceExpr{
	pos: position{line: 205, col: 25, offset: 3908},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 205, col: 25, offset: 3908},
	name: "NL",
},
&litMatcher{
	pos: position{line: 205, col: 30, offset: 3913},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 205, col: 36, offset: 3919},
	name: "COMMENT",
},
	},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 206, col: 1, offset: 3928},
	expr: &zeroOrMoreExpr{
	pos: position{line: 206, col: 20, offset: 3947},
	expr: &choiceExpr{
	pos: position{line: 206, col: 21, offset: 3948},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 206, col: 21, offset: 3948},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 206, col: 29, offset: 3956},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 207, col: 1, offset: 3966},
	expr: &oneOrMoreExpr{
	pos: position{line: 207, col: 35, offset: 4000},
	expr: &choiceExpr{
	pos: position{line: 207, col: 36, offset: 4001},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 207, col: 36, offset: 4001},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 207, col: 44, offset: 4009},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 207, col: 54, offset: 4019},
	name: "NL",
},
	},
},
},
},
{
	name: "COMMENT",
	pos: position{line: 209, col: 1, offset: 4025},
	expr: &seqExpr{
	pos: position{line: 209, col: 12, offset: 4036},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 209, col: 12, offset: 4036},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 209, col: 17, offset: 4041},
	expr: &seqExpr{
	pos: position{line: 209, col: 19, offset: 4043},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 209, col: 19, offset: 4043},
	expr: &litMatcher{
	pos: position{line: 209, col: 20, offset: 4044},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 209, col: 25, offset: 4049,
},
	},
},
},
&choiceExpr{
	pos: position{line: 209, col: 31, offset: 4055},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 209, col: 31, offset: 4055},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 209, col: 38, offset: 4062},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 211, col: 1, offset: 4068},
	expr: &notExpr{
	pos: position{line: 211, col: 8, offset: 4075},
	expr: &anyMatcher{
	line: 211, col: 9, offset: 4076,
},
},
},
	},
}
func (c *current) onQUERY1(us, firstBlock, otherBlocks interface{}) (interface{}, error) {
	return newQuery(us, firstBlock, otherBlocks)
}

func (p *parser) callonQUERY1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQUERY1(stack["us"], stack["firstBlock"], stack["otherBlocks"])
}

func (c *current) onUSE1(r, v interface{}) (interface{}, error) {
	return newUse(r, v)
}

func (p *parser) callonUSE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE1(stack["r"], stack["v"])
}

func (c *current) onUSE_ACTION1() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonUSE_ACTION1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE_ACTION1()
}

func (c *current) onUSE_VALUE1(v interface{}) (interface{}, error) {
	return newUseValue(v)
}

func (p *parser) callonUSE_VALUE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE_VALUE1(stack["v"])
}

func (c *current) onBLOCK1(action, m, w, f, fl interface{}) (interface{}, error) {
	return newBlock(action, m, w, f, fl)
}

func (p *parser) callonBLOCK1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBLOCK1(stack["action"], stack["m"], stack["w"], stack["f"], stack["fl"])
}

func (c *current) onACTION_RULE1(m, r, a, i interface{}) (interface{}, error) {
	return newActionRule(m, r, a, i)
}

func (p *parser) callonACTION_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onACTION_RULE1(stack["m"], stack["r"], stack["a"], stack["i"])
}

func (c *current) onMETHOD1(m interface{}) (interface{}, error) {
	return newMethod(m), nil
}

func (p *parser) callonMETHOD1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMETHOD1(stack["m"])
}

func (c *current) onALIAS1(a interface{}) (interface{}, error) {
	return a, nil
}

func (p *parser) callonALIAS1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onALIAS1(stack["a"])
}

func (c *current) onIN1(t interface{}) (interface{}, error) {
	return newIn(t)
}

func (p *parser) callonIN1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIN1(stack["t"])
}

func (c *current) onMODIFIER_RULE1(m interface{}) (interface{}, error) {
	return m, nil
}

func (p *parser) callonMODIFIER_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMODIFIER_RULE1(stack["m"])
}

func (c *current) onWITH_RULE1(pb, kvs interface{}) (interface{}, error) {
	return newWith(pb, kvs)
}

func (p *parser) callonWITH_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWITH_RULE1(stack["pb"], stack["kvs"])
}

func (c *current) onPARAMETER_BODY1(t, fn interface{}) (interface{}, error) {
	return newParameterBody(t, fn)
}

func (p *parser) callonPARAMETER_BODY1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPARAMETER_BODY1(stack["t"], stack["fn"])
}

func (c *current) onKEY_VALUE_LIST1(first, others interface{}) (interface{}, error) {
	return newKeyValueList(first, others)
}

func (p *parser) callonKEY_VALUE_LIST1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKEY_VALUE_LIST1(stack["first"], stack["others"])
}

func (c *current) onKEY_VALUE1(k, v, fn interface{}) (interface{}, error) {
	return newKeyValue(k, v, fn)
}

func (p *parser) callonKEY_VALUE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKEY_VALUE1(stack["k"], stack["v"], stack["fn"])
}

func (c *current) onAPPLY_FN1(fn interface{}) (interface{}, error) {
	return fn, nil
}

func (p *parser) callonAPPLY_FN1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAPPLY_FN1(stack["fn"])
}

func (c *current) onFUNCTION1() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonFUNCTION1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFUNCTION1()
}

func (c *current) onVALUE1(v interface{}) (interface{}, error) {
	return newValue(v)
}

func (p *parser) callonVALUE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVALUE1(stack["v"])
}

func (c *current) onLIST1(l interface{}) (interface{}, error) {
	return l, nil
}

func (p *parser) callonLIST1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLIST1(stack["l"])
}

func (c *current) onEMPTY_LIST1() (interface{}, error) {
	return newEmptyList()
}

func (p *parser) callonEMPTY_LIST1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEMPTY_LIST1()
}

func (c *current) onPOPULATED_LIST1(i, ii interface{}) (interface{}, error) {
	return newList(i, ii)
}

func (p *parser) callonPOPULATED_LIST1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPOPULATED_LIST1(stack["i"], stack["ii"])
}

func (c *current) onOBJECT1(o interface{}) (interface{}, error) {
	return o, nil
}

func (p *parser) callonOBJECT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOBJECT1(stack["o"])
}

func (c *current) onEMPTY_OBJ1() (interface{}, error) {
	return newEmptyObject()
}

func (p *parser) callonEMPTY_OBJ1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEMPTY_OBJ1()
}

func (c *current) onPOPULATED_OBJ1(oe, oes interface{}) (interface{}, error) {
	return newPopulatedObject(oe, oes)
}

func (p *parser) callonPOPULATED_OBJ1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPOPULATED_OBJ1(stack["oe"], stack["oes"])
}

func (c *current) onOBJ_ENTRY1(k, v interface{}) (interface{}, error) {
	return newObjectEntry(k, v)
}

func (p *parser) callonOBJ_ENTRY1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOBJ_ENTRY1(stack["k"], stack["v"])
}

func (c *current) onPRIMITIVE1(p interface{}) (interface{}, error) {
	return newPrimitive(p)
}

func (p *parser) callonPRIMITIVE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPRIMITIVE1(stack["p"])
}

func (c *current) onCHAIN1(i, ii interface{}) (interface{}, error) {
	return newChain(i, ii)
}

func (p *parser) callonCHAIN1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCHAIN1(stack["i"], stack["ii"])
}

func (c *current) onCHAINED_ITEM1(ci interface{}) (interface{}, error) {
	return newChained(ci)
}

func (p *parser) callonCHAINED_ITEM1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCHAINED_ITEM1(stack["ci"])
}

func (c *current) onONLY_RULE1(f, fs interface{}) (interface{}, error) {
	return newOnly(f, fs)
}

func (p *parser) callonONLY_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onONLY_RULE1(stack["f"], stack["fs"])
}

func (c *current) onFILTER1(f, fn interface{}) (interface{}, error) {
	return newFilter(f, fn)
}

func (p *parser) callonFILTER1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFILTER1(stack["f"], stack["fn"])
}

func (c *current) onMATCHES_FN1(arg interface{}) (interface{}, error) {
	return newMatchesFunction(arg)
}

func (p *parser) callonMATCHES_FN1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMATCHES_FN1(stack["arg"])
}

func (c *current) onHEADERS1(h, hs interface{}) (interface{}, error) {
	return newHeaders(h, hs)
}

func (p *parser) callonHEADERS1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHEADERS1(stack["h"], stack["hs"])
}

func (c *current) onHEADER1(n, v interface{}) (interface{}, error) {
	return newHeader(n, v)
}

func (p *parser) callonHEADER1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHEADER1(stack["n"], stack["v"])
}

func (c *current) onHIDDEN_RULE1() (interface{}, error) {
	return newHidden()
}

func (p *parser) callonHIDDEN_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHIDDEN_RULE1()
}

func (c *current) onTIMEOUT1(t interface{}) (interface{}, error) {
	return newTimeout(t)
}

func (p *parser) callonTIMEOUT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTIMEOUT1(stack["t"])
}

func (c *current) onMAX_AGE1(t interface{}) (interface{}, error) {
	return newMaxAge(t)
}

func (p *parser) callonMAX_AGE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMAX_AGE1(stack["t"])
}

func (c *current) onS_MAX_AGE1(t interface{}) (interface{}, error) {
	return newSmaxAge(t)
}

func (p *parser) callonS_MAX_AGE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onS_MAX_AGE1(stack["t"])
}

func (c *current) onFLAG_RULE1() (interface{}, error) {
	return newIgnoreErrors()
}

func (p *parser) callonFLAG_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFLAG_RULE1()
}

func (c *current) onVARIABLE1(v interface{}) (interface{}, error) {
	return newVariable(v)
}

func (p *parser) callonVARIABLE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVARIABLE1(stack["v"])
}

func (c *current) onIDENT1() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonIDENT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIDENT1()
}

func (c *current) onIDENT_WITH_DOT1() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonIDENT_WITH_DOT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIDENT_WITH_DOT1()
}

func (c *current) onNull1() (interface{}, error) {
	return newNull()
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
}

func (c *current) onBoolean1() (interface{}, error) {
	return newBoolean(c.text)
}

func (p *parser) callonBoolean1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolean1()
}

func (c *current) onString1() (interface{}, error) {
	return newString(c.text)
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onFloat1() (interface{}, error) {
	return newFloat(c.text)
}

func (p *parser) callonFloat1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFloat1()
}

func (c *current) onInteger1() (interface{}, error) {
	return newInteger(c.text)
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}


var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule          = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch         = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos    position
	expr   interface{}
	run    func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs: new(errList),
		data: b,
		pt: savepoint{position: position{line: 1}},
		recover: true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v interface{}
	b bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug bool
	depth  int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules  map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth) + ">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth) + "<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}

