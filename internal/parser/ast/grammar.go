
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
&zeroOrMoreExpr{
	pos: position{line: 17, col: 10, offset: 127},
	expr: &choiceExpr{
	pos: position{line: 17, col: 11, offset: 128},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 11, offset: 128},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 17, col: 16, offset: 133},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 17, col: 24, offset: 141},
	name: "COMMENT",
},
	},
},
},
&labeledExpr{
	pos: position{line: 17, col: 34, offset: 151},
	label: "us",
	expr: &zeroOrMoreExpr{
	pos: position{line: 17, col: 37, offset: 154},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 38, offset: 155},
	name: "USE",
},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 44, offset: 161},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 47, offset: 164},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 47, offset: 164},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 51, offset: 168},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 17, col: 54, offset: 171},
	label: "firstBlock",
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 65, offset: 182},
	name: "BLOCK",
},
},
&labeledExpr{
	pos: position{line: 17, col: 71, offset: 188},
	label: "otherBlocks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 17, col: 83, offset: 200},
	expr: &seqExpr{
	pos: position{line: 17, col: 84, offset: 201},
	exprs: []interface{}{
&oneOrMoreExpr{
	pos: position{line: 17, col: 84, offset: 201},
	expr: &seqExpr{
	pos: position{line: 17, col: 85, offset: 202},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 85, offset: 202},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 88, offset: 205},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 91, offset: 208},
	name: "WS",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 96, offset: 213},
	name: "BLOCK",
},
	},
},
},
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 104, offset: 221},
	expr: &choiceExpr{
	pos: position{line: 17, col: 105, offset: 222},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 105, offset: 222},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 17, col: 110, offset: 227},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 17, col: 118, offset: 235},
	name: "COMMENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 128, offset: 245},
	name: "EOF",
},
	},
},
},
},
{
	name: "USE",
	pos: position{line: 21, col: 1, offset: 300},
	expr: &actionExpr{
	pos: position{line: 21, col: 8, offset: 307},
	run: (*parser).callonUSE1,
	expr: &seqExpr{
	pos: position{line: 21, col: 8, offset: 307},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 21, col: 8, offset: 307},
	val: "use",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 21, col: 14, offset: 313},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 21, col: 22, offset: 321},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 25, offset: 324},
	name: "USE_ACTION",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 37, offset: 336},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 21, col: 40, offset: 339},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 43, offset: 342},
	name: "USE_VALUE",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 54, offset: 353},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 21, col: 57, offset: 356},
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 57, offset: 356},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 61, offset: 360},
	name: "WS",
},
	},
},
},
},
{
	name: "USE_ACTION",
	pos: position{line: 25, col: 1, offset: 389},
	expr: &choiceExpr{
	pos: position{line: 25, col: 15, offset: 403},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 25, col: 15, offset: 403},
	run: (*parser).callonUSE_ACTION2,
	expr: &choiceExpr{
	pos: position{line: 25, col: 16, offset: 404},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 16, offset: 404},
	val: "timeout",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 28, offset: 416},
	val: "max-age",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 40, offset: 428},
	val: "s-max-age",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 27, col: 5, offset: 473},
	run: (*parser).callonUSE_ACTION7,
	expr: &litMatcher{
	pos: position{line: 27, col: 5, offset: 473},
	val: "cache-control",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "USE_VALUE",
	pos: position{line: 31, col: 1, offset: 556},
	expr: &actionExpr{
	pos: position{line: 31, col: 14, offset: 569},
	run: (*parser).callonUSE_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 31, col: 14, offset: 569},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 31, col: 17, offset: 572},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 31, col: 17, offset: 572},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 31, col: 26, offset: 581},
	name: "Integer",
},
	},
},
},
},
},
{
	name: "BLOCK",
	pos: position{line: 35, col: 1, offset: 618},
	expr: &actionExpr{
	pos: position{line: 35, col: 10, offset: 627},
	run: (*parser).callonBLOCK1,
	expr: &seqExpr{
	pos: position{line: 35, col: 10, offset: 627},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 35, col: 10, offset: 627},
	label: "action",
	expr: &ruleRefExpr{
	pos: position{line: 35, col: 18, offset: 635},
	name: "ACTION_RULE",
},
},
&labeledExpr{
	pos: position{line: 35, col: 31, offset: 648},
	label: "m",
	expr: &zeroOrOneExpr{
	pos: position{line: 35, col: 34, offset: 651},
	expr: &ruleRefExpr{
	pos: position{line: 35, col: 34, offset: 651},
	name: "MODIFIER_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 35, col: 50, offset: 667},
	label: "w",
	expr: &zeroOrOneExpr{
	pos: position{line: 35, col: 53, offset: 670},
	expr: &ruleRefExpr{
	pos: position{line: 35, col: 53, offset: 670},
	name: "WITH_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 35, col: 65, offset: 682},
	label: "f",
	expr: &zeroOrOneExpr{
	pos: position{line: 35, col: 67, offset: 684},
	expr: &choiceExpr{
	pos: position{line: 35, col: 68, offset: 685},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 35, col: 68, offset: 685},
	name: "HIDDEN_RULE",
},
&ruleRefExpr{
	pos: position{line: 35, col: 82, offset: 699},
	name: "ONLY_RULE",
},
	},
},
},
},
&labeledExpr{
	pos: position{line: 35, col: 94, offset: 711},
	label: "fl",
	expr: &zeroOrOneExpr{
	pos: position{line: 35, col: 98, offset: 715},
	expr: &ruleRefExpr{
	pos: position{line: 35, col: 98, offset: 715},
	name: "FLAGS_RULE",
},
},
},
&ruleRefExpr{
	pos: position{line: 35, col: 111, offset: 728},
	name: "WS",
},
	},
},
},
},
{
	name: "ACTION_RULE",
	pos: position{line: 39, col: 1, offset: 774},
	expr: &actionExpr{
	pos: position{line: 39, col: 16, offset: 789},
	run: (*parser).callonACTION_RULE1,
	expr: &seqExpr{
	pos: position{line: 39, col: 16, offset: 789},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 39, col: 16, offset: 789},
	label: "m",
	expr: &ruleRefExpr{
	pos: position{line: 39, col: 19, offset: 792},
	name: "METHOD",
},
},
&ruleRefExpr{
	pos: position{line: 39, col: 27, offset: 800},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 39, col: 35, offset: 808},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 39, col: 38, offset: 811},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 39, col: 45, offset: 818},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 39, col: 48, offset: 821},
	expr: &ruleRefExpr{
	pos: position{line: 39, col: 48, offset: 821},
	name: "ALIAS",
},
},
},
&labeledExpr{
	pos: position{line: 39, col: 56, offset: 829},
	label: "i",
	expr: &zeroOrOneExpr{
	pos: position{line: 39, col: 59, offset: 832},
	expr: &ruleRefExpr{
	pos: position{line: 39, col: 59, offset: 832},
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
	pos: position{line: 43, col: 1, offset: 876},
	expr: &actionExpr{
	pos: position{line: 43, col: 11, offset: 886},
	run: (*parser).callonMETHOD1,
	expr: &choiceExpr{
	pos: position{line: 43, col: 12, offset: 887},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 43, col: 12, offset: 887},
	val: "from",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 43, col: 21, offset: 896},
	val: "to",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 43, col: 28, offset: 903},
	val: "into",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 43, col: 36, offset: 911},
	val: "update",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 43, col: 47, offset: 922},
	val: "delete",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "ALIAS",
	pos: position{line: 47, col: 1, offset: 963},
	expr: &actionExpr{
	pos: position{line: 47, col: 10, offset: 972},
	run: (*parser).callonALIAS1,
	expr: &seqExpr{
	pos: position{line: 47, col: 10, offset: 972},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 47, col: 10, offset: 972},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 47, col: 18, offset: 980},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 47, col: 23, offset: 985},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 47, col: 31, offset: 993},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 47, col: 34, offset: 996},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IN",
	pos: position{line: 51, col: 1, offset: 1023},
	expr: &actionExpr{
	pos: position{line: 51, col: 7, offset: 1029},
	run: (*parser).callonIN1,
	expr: &seqExpr{
	pos: position{line: 51, col: 7, offset: 1029},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 51, col: 7, offset: 1029},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 51, col: 15, offset: 1037},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 51, col: 20, offset: 1042},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 51, col: 28, offset: 1050},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 51, col: 31, offset: 1053},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "MODIFIER_RULE",
	pos: position{line: 55, col: 1, offset: 1091},
	expr: &actionExpr{
	pos: position{line: 55, col: 18, offset: 1108},
	run: (*parser).callonMODIFIER_RULE1,
	expr: &labeledExpr{
	pos: position{line: 55, col: 18, offset: 1108},
	label: "m",
	expr: &oneOrMoreExpr{
	pos: position{line: 55, col: 20, offset: 1110},
	expr: &choiceExpr{
	pos: position{line: 55, col: 21, offset: 1111},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 55, col: 21, offset: 1111},
	name: "HEADERS",
},
&ruleRefExpr{
	pos: position{line: 55, col: 31, offset: 1121},
	name: "TIMEOUT",
},
&ruleRefExpr{
	pos: position{line: 55, col: 41, offset: 1131},
	name: "MAX_AGE",
},
&ruleRefExpr{
	pos: position{line: 55, col: 51, offset: 1141},
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
	pos: position{line: 59, col: 1, offset: 1173},
	expr: &actionExpr{
	pos: position{line: 59, col: 14, offset: 1186},
	run: (*parser).callonWITH_RULE1,
	expr: &seqExpr{
	pos: position{line: 59, col: 14, offset: 1186},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 59, col: 14, offset: 1186},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 59, col: 22, offset: 1194},
	val: "with",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 59, col: 29, offset: 1201},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 59, col: 37, offset: 1209},
	label: "pb",
	expr: &zeroOrOneExpr{
	pos: position{line: 59, col: 40, offset: 1212},
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 40, offset: 1212},
	name: "PARAMETER_BODY",
},
},
},
&labeledExpr{
	pos: position{line: 59, col: 56, offset: 1228},
	label: "kvs",
	expr: &zeroOrOneExpr{
	pos: position{line: 59, col: 60, offset: 1232},
	expr: &ruleRefExpr{
	pos: position{line: 59, col: 60, offset: 1232},
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
	pos: position{line: 63, col: 1, offset: 1278},
	expr: &actionExpr{
	pos: position{line: 63, col: 19, offset: 1296},
	run: (*parser).callonPARAMETER_BODY1,
	expr: &seqExpr{
	pos: position{line: 63, col: 19, offset: 1296},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 63, col: 19, offset: 1296},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 63, col: 23, offset: 1300},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 63, col: 26, offset: 1303},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 63, col: 33, offset: 1310},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 63, col: 37, offset: 1314},
	expr: &ruleRefExpr{
	pos: position{line: 63, col: 37, offset: 1314},
	name: "APPLY_FN",
},
},
},
&ruleRefExpr{
	pos: position{line: 63, col: 48, offset: 1325},
	name: "WS",
},
&zeroOrOneExpr{
	pos: position{line: 63, col: 51, offset: 1328},
	expr: &ruleRefExpr{
	pos: position{line: 63, col: 51, offset: 1328},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 63, col: 55, offset: 1332},
	name: "WS",
},
	},
},
},
},
{
	name: "KEY_VALUE_LIST",
	pos: position{line: 67, col: 1, offset: 1372},
	expr: &actionExpr{
	pos: position{line: 67, col: 19, offset: 1390},
	run: (*parser).callonKEY_VALUE_LIST1,
	expr: &seqExpr{
	pos: position{line: 67, col: 19, offset: 1390},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 67, col: 19, offset: 1390},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 67, col: 25, offset: 1396},
	name: "KEY_VALUE",
},
},
&labeledExpr{
	pos: position{line: 67, col: 35, offset: 1406},
	label: "others",
	expr: &zeroOrMoreExpr{
	pos: position{line: 67, col: 42, offset: 1413},
	expr: &seqExpr{
	pos: position{line: 67, col: 43, offset: 1414},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 67, col: 43, offset: 1414},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 67, col: 46, offset: 1417},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 67, col: 49, offset: 1420},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 67, col: 52, offset: 1423},
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
	pos: position{line: 71, col: 1, offset: 1479},
	expr: &actionExpr{
	pos: position{line: 71, col: 14, offset: 1492},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 71, col: 14, offset: 1492},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 71, col: 14, offset: 1492},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 71, col: 17, offset: 1495},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 71, col: 33, offset: 1511},
	name: "WS",
},
&litMatcher{
	pos: position{line: 71, col: 36, offset: 1514},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 71, col: 40, offset: 1518},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 71, col: 43, offset: 1521},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 71, col: 46, offset: 1524},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 71, col: 53, offset: 1531},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 71, col: 57, offset: 1535},
	expr: &ruleRefExpr{
	pos: position{line: 71, col: 57, offset: 1535},
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
	pos: position{line: 75, col: 1, offset: 1581},
	expr: &actionExpr{
	pos: position{line: 75, col: 13, offset: 1593},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 75, col: 13, offset: 1593},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 75, col: 13, offset: 1593},
	name: "WS",
},
&litMatcher{
	pos: position{line: 75, col: 16, offset: 1596},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 75, col: 21, offset: 1601},
	expr: &ruleRefExpr{
	pos: position{line: 75, col: 21, offset: 1601},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 75, col: 25, offset: 1605},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 75, col: 29, offset: 1609},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 79, col: 1, offset: 1640},
	expr: &choiceExpr{
	pos: position{line: 79, col: 13, offset: 1652},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 79, col: 13, offset: 1652},
	run: (*parser).callonFUNCTION2,
	expr: &choiceExpr{
	pos: position{line: 79, col: 14, offset: 1653},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 79, col: 14, offset: 1653},
	val: "flatten",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 79, col: 26, offset: 1665},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 79, col: 37, offset: 1676},
	val: "json",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 81, col: 5, offset: 1716},
	run: (*parser).callonFUNCTION7,
	expr: &litMatcher{
	pos: position{line: 81, col: 5, offset: 1716},
	val: "contract",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "VALUE",
	pos: position{line: 85, col: 1, offset: 1791},
	expr: &actionExpr{
	pos: position{line: 85, col: 10, offset: 1800},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 85, col: 10, offset: 1800},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 85, col: 13, offset: 1803},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 85, col: 13, offset: 1803},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 85, col: 20, offset: 1810},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 85, col: 29, offset: 1819},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 85, col: 40, offset: 1830},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 89, col: 1, offset: 1866},
	expr: &actionExpr{
	pos: position{line: 89, col: 9, offset: 1874},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 89, col: 9, offset: 1874},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 89, col: 12, offset: 1877},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 89, col: 12, offset: 1877},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 89, col: 25, offset: 1890},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 93, col: 1, offset: 1926},
	expr: &actionExpr{
	pos: position{line: 93, col: 15, offset: 1940},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 93, col: 15, offset: 1940},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 93, col: 15, offset: 1940},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 93, col: 19, offset: 1944},
	name: "WS",
},
&litMatcher{
	pos: position{line: 93, col: 22, offset: 1947},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 97, col: 1, offset: 1979},
	expr: &actionExpr{
	pos: position{line: 97, col: 19, offset: 1997},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 97, col: 19, offset: 1997},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 97, col: 19, offset: 1997},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 97, col: 23, offset: 2001},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 97, col: 26, offset: 2004},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 97, col: 28, offset: 2006},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 97, col: 34, offset: 2012},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 97, col: 37, offset: 2015},
	expr: &seqExpr{
	pos: position{line: 97, col: 38, offset: 2016},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 97, col: 38, offset: 2016},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 97, col: 41, offset: 2019},
	expr: &ruleRefExpr{
	pos: position{line: 97, col: 41, offset: 2019},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 97, col: 45, offset: 2023},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 97, col: 48, offset: 2026},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 97, col: 56, offset: 2034},
	name: "WS",
},
&litMatcher{
	pos: position{line: 97, col: 59, offset: 2037},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 101, col: 1, offset: 2069},
	expr: &actionExpr{
	pos: position{line: 101, col: 11, offset: 2079},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 101, col: 11, offset: 2079},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 101, col: 14, offset: 2082},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 101, col: 14, offset: 2082},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 101, col: 26, offset: 2094},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 105, col: 1, offset: 2129},
	expr: &actionExpr{
	pos: position{line: 105, col: 14, offset: 2142},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 105, col: 14, offset: 2142},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 105, col: 14, offset: 2142},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 18, offset: 2146},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 21, offset: 2149},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 21, offset: 2149},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 25, offset: 2153},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 28, offset: 2156},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 109, col: 1, offset: 2190},
	expr: &actionExpr{
	pos: position{line: 109, col: 18, offset: 2207},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 109, col: 18, offset: 2207},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 109, col: 18, offset: 2207},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 109, col: 22, offset: 2211},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 109, col: 25, offset: 2214},
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 25, offset: 2214},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 29, offset: 2218},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 109, col: 32, offset: 2221},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 36, offset: 2225},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 109, col: 47, offset: 2236},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 109, col: 51, offset: 2240},
	expr: &seqExpr{
	pos: position{line: 109, col: 52, offset: 2241},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 109, col: 52, offset: 2241},
	name: "WS",
},
&litMatcher{
	pos: position{line: 109, col: 55, offset: 2244},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 109, col: 59, offset: 2248},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 109, col: 62, offset: 2251},
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 62, offset: 2251},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 66, offset: 2255},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 109, col: 69, offset: 2258},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 81, offset: 2270},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 109, col: 84, offset: 2273},
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 84, offset: 2273},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 88, offset: 2277},
	name: "WS",
},
&litMatcher{
	pos: position{line: 109, col: 91, offset: 2280},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 113, col: 1, offset: 2325},
	expr: &actionExpr{
	pos: position{line: 113, col: 14, offset: 2338},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 113, col: 14, offset: 2338},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 113, col: 14, offset: 2338},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 113, col: 17, offset: 2341},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2341},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 113, col: 26, offset: 2350},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 113, col: 33, offset: 2357},
	name: "WS",
},
&litMatcher{
	pos: position{line: 113, col: 36, offset: 2360},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 113, col: 40, offset: 2364},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 113, col: 43, offset: 2367},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 113, col: 46, offset: 2370},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 117, col: 1, offset: 2411},
	expr: &actionExpr{
	pos: position{line: 117, col: 14, offset: 2424},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 117, col: 14, offset: 2424},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 117, col: 17, offset: 2427},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 117, col: 17, offset: 2427},
	name: "Null",
},
&ruleRefExpr{
	pos: position{line: 117, col: 24, offset: 2434},
	name: "Boolean",
},
&ruleRefExpr{
	pos: position{line: 117, col: 34, offset: 2444},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 117, col: 43, offset: 2453},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 117, col: 51, offset: 2461},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 117, col: 61, offset: 2471},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "CHAIN",
	pos: position{line: 121, col: 1, offset: 2507},
	expr: &actionExpr{
	pos: position{line: 121, col: 10, offset: 2516},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 121, col: 10, offset: 2516},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 121, col: 10, offset: 2516},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 121, col: 13, offset: 2519},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 121, col: 27, offset: 2533},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 121, col: 30, offset: 2536},
	expr: &seqExpr{
	pos: position{line: 121, col: 31, offset: 2537},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 121, col: 31, offset: 2537},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 121, col: 35, offset: 2541},
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
	pos: position{line: 125, col: 1, offset: 2585},
	expr: &actionExpr{
	pos: position{line: 125, col: 17, offset: 2601},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 125, col: 17, offset: 2601},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 125, col: 21, offset: 2605},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 21, offset: 2605},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 125, col: 32, offset: 2616},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 129, col: 1, offset: 2651},
	expr: &actionExpr{
	pos: position{line: 129, col: 14, offset: 2664},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 129, col: 14, offset: 2664},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 14, offset: 2664},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 129, col: 22, offset: 2672},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 129, col: 29, offset: 2679},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 129, col: 37, offset: 2687},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 40, offset: 2690},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 129, col: 48, offset: 2698},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 129, col: 51, offset: 2701},
	expr: &seqExpr{
	pos: position{line: 129, col: 52, offset: 2702},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 52, offset: 2702},
	name: "WS",
},
&notExpr{
	pos: position{line: 129, col: 55, offset: 2705},
	expr: &choiceExpr{
	pos: position{line: 129, col: 57, offset: 2707},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 57, offset: 2707},
	name: "FLAGS_RULE",
},
&seqExpr{
	pos: position{line: 129, col: 70, offset: 2720},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 129, col: 70, offset: 2720},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 129, col: 73, offset: 2723},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 129, col: 76, offset: 2726},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 129, col: 79, offset: 2729},
	name: "BLOCK",
},
	},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 129, col: 86, offset: 2736},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 129, col: 89, offset: 2739},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 129, col: 92, offset: 2742},
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
	pos: position{line: 133, col: 1, offset: 2779},
	expr: &actionExpr{
	pos: position{line: 133, col: 11, offset: 2789},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 133, col: 11, offset: 2789},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 133, col: 11, offset: 2789},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 133, col: 14, offset: 2792},
	name: "FILTER_VALUE",
},
},
&labeledExpr{
	pos: position{line: 133, col: 28, offset: 2806},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 133, col: 32, offset: 2810},
	expr: &ruleRefExpr{
	pos: position{line: 133, col: 32, offset: 2810},
	name: "MATCHES_FN",
},
},
},
	},
},
},
},
{
	name: "FILTER_VALUE",
	pos: position{line: 137, col: 1, offset: 2853},
	expr: &actionExpr{
	pos: position{line: 137, col: 17, offset: 2869},
	run: (*parser).callonFILTER_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 137, col: 17, offset: 2869},
	label: "fv",
	expr: &choiceExpr{
	pos: position{line: 137, col: 21, offset: 2873},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 21, offset: 2873},
	name: "IDENT_WITH_DOT",
},
&litMatcher{
	pos: position{line: 137, col: 38, offset: 2890},
	val: "*",
	ignoreCase: false,
},
	},
},
},
},
},
{
	name: "MATCHES_FN",
	pos: position{line: 141, col: 1, offset: 2927},
	expr: &actionExpr{
	pos: position{line: 141, col: 15, offset: 2941},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 141, col: 15, offset: 2941},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 141, col: 15, offset: 2941},
	name: "WS",
},
&litMatcher{
	pos: position{line: 141, col: 18, offset: 2944},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 141, col: 23, offset: 2949},
	name: "WS",
},
&litMatcher{
	pos: position{line: 141, col: 26, offset: 2952},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 141, col: 36, offset: 2962},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 141, col: 40, offset: 2966},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 141, col: 45, offset: 2971},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 141, col: 53, offset: 2979},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 145, col: 1, offset: 3020},
	expr: &actionExpr{
	pos: position{line: 145, col: 12, offset: 3031},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 145, col: 12, offset: 3031},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 145, col: 12, offset: 3031},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 145, col: 20, offset: 3039},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 145, col: 30, offset: 3049},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 145, col: 38, offset: 3057},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 145, col: 41, offset: 3060},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 145, col: 49, offset: 3068},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 145, col: 52, offset: 3071},
	expr: &seqExpr{
	pos: position{line: 145, col: 53, offset: 3072},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 145, col: 53, offset: 3072},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 145, col: 56, offset: 3075},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 145, col: 59, offset: 3078},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 145, col: 62, offset: 3081},
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
	pos: position{line: 149, col: 1, offset: 3121},
	expr: &actionExpr{
	pos: position{line: 149, col: 11, offset: 3131},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 149, col: 11, offset: 3131},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 149, col: 11, offset: 3131},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 149, col: 14, offset: 3134},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 149, col: 21, offset: 3141},
	name: "WS",
},
&litMatcher{
	pos: position{line: 149, col: 24, offset: 3144},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 149, col: 28, offset: 3148},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 149, col: 31, offset: 3151},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 149, col: 34, offset: 3154},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 34, offset: 3154},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 149, col: 45, offset: 3165},
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
	pos: position{line: 153, col: 1, offset: 3202},
	expr: &actionExpr{
	pos: position{line: 153, col: 16, offset: 3217},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 153, col: 16, offset: 3217},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 16, offset: 3217},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 153, col: 24, offset: 3225},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 157, col: 1, offset: 3259},
	expr: &actionExpr{
	pos: position{line: 157, col: 12, offset: 3270},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 157, col: 12, offset: 3270},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 12, offset: 3270},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 157, col: 20, offset: 3278},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 30, offset: 3288},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 157, col: 38, offset: 3296},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 157, col: 41, offset: 3299},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 41, offset: 3299},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 157, col: 52, offset: 3310},
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
	pos: position{line: 161, col: 1, offset: 3346},
	expr: &actionExpr{
	pos: position{line: 161, col: 12, offset: 3357},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 161, col: 12, offset: 3357},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 12, offset: 3357},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 161, col: 20, offset: 3365},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 161, col: 30, offset: 3375},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 161, col: 38, offset: 3383},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 161, col: 41, offset: 3386},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 41, offset: 3386},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 161, col: 52, offset: 3397},
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
	pos: position{line: 165, col: 1, offset: 3432},
	expr: &actionExpr{
	pos: position{line: 165, col: 14, offset: 3445},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 165, col: 14, offset: 3445},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 165, col: 14, offset: 3445},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 165, col: 22, offset: 3453},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 165, col: 34, offset: 3465},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 165, col: 42, offset: 3473},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 165, col: 45, offset: 3476},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 165, col: 45, offset: 3476},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 165, col: 56, offset: 3487},
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
	name: "FLAGS_RULE",
	pos: position{line: 169, col: 1, offset: 3523},
	expr: &actionExpr{
	pos: position{line: 169, col: 15, offset: 3537},
	run: (*parser).callonFLAGS_RULE1,
	expr: &seqExpr{
	pos: position{line: 169, col: 15, offset: 3537},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 169, col: 15, offset: 3537},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 169, col: 23, offset: 3545},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 169, col: 25, offset: 3547},
	name: "IGNORE_FLAG",
},
},
&labeledExpr{
	pos: position{line: 169, col: 37, offset: 3559},
	label: "is",
	expr: &zeroOrMoreExpr{
	pos: position{line: 169, col: 40, offset: 3562},
	expr: &seqExpr{
	pos: position{line: 169, col: 41, offset: 3563},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 169, col: 41, offset: 3563},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 169, col: 44, offset: 3566},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 169, col: 47, offset: 3569},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 169, col: 50, offset: 3572},
	name: "IGNORE_FLAG",
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
	name: "IGNORE_FLAG",
	pos: position{line: 173, col: 1, offset: 3615},
	expr: &actionExpr{
	pos: position{line: 173, col: 16, offset: 3630},
	run: (*parser).callonIGNORE_FLAG1,
	expr: &litMatcher{
	pos: position{line: 173, col: 16, offset: 3630},
	val: "ignore-errors",
	ignoreCase: false,
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 177, col: 1, offset: 3677},
	expr: &actionExpr{
	pos: position{line: 177, col: 13, offset: 3689},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 177, col: 13, offset: 3689},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 177, col: 13, offset: 3689},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 177, col: 17, offset: 3693},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 177, col: 20, offset: 3696},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 181, col: 1, offset: 3731},
	expr: &actionExpr{
	pos: position{line: 181, col: 10, offset: 3740},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 181, col: 10, offset: 3740},
	expr: &charClassMatcher{
	pos: position{line: 181, col: 10, offset: 3740},
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
	pos: position{line: 185, col: 1, offset: 3786},
	expr: &actionExpr{
	pos: position{line: 185, col: 19, offset: 3804},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 185, col: 19, offset: 3804},
	expr: &charClassMatcher{
	pos: position{line: 185, col: 19, offset: 3804},
	val: "[a-zA-Z0-9-_.]",
	chars: []rune{'-','_','.',},
	ranges: []rune{'a','z','A','Z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
},
{
	name: "Null",
	pos: position{line: 189, col: 1, offset: 3851},
	expr: &actionExpr{
	pos: position{line: 189, col: 9, offset: 3859},
	run: (*parser).callonNull1,
	expr: &litMatcher{
	pos: position{line: 189, col: 9, offset: 3859},
	val: "null",
	ignoreCase: false,
},
},
},
{
	name: "Boolean",
	pos: position{line: 193, col: 1, offset: 3889},
	expr: &actionExpr{
	pos: position{line: 193, col: 12, offset: 3900},
	run: (*parser).callonBoolean1,
	expr: &choiceExpr{
	pos: position{line: 193, col: 13, offset: 3901},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 193, col: 13, offset: 3901},
	val: "true",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 193, col: 22, offset: 3910},
	val: "false",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "String",
	pos: position{line: 197, col: 1, offset: 3951},
	expr: &actionExpr{
	pos: position{line: 197, col: 11, offset: 3961},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 197, col: 11, offset: 3961},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 197, col: 11, offset: 3961},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 197, col: 15, offset: 3965},
	expr: &seqExpr{
	pos: position{line: 197, col: 17, offset: 3967},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 197, col: 17, offset: 3967},
	expr: &litMatcher{
	pos: position{line: 197, col: 18, offset: 3968},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 197, col: 22, offset: 3972,
},
	},
},
},
&litMatcher{
	pos: position{line: 197, col: 27, offset: 3977},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 201, col: 1, offset: 4012},
	expr: &actionExpr{
	pos: position{line: 201, col: 10, offset: 4021},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 201, col: 10, offset: 4021},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 201, col: 10, offset: 4021},
	expr: &choiceExpr{
	pos: position{line: 201, col: 11, offset: 4022},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 201, col: 11, offset: 4022},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 201, col: 17, offset: 4028},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 201, col: 23, offset: 4034},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 201, col: 31, offset: 4042},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 201, col: 35, offset: 4046},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 205, col: 1, offset: 4084},
	expr: &actionExpr{
	pos: position{line: 205, col: 12, offset: 4095},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 205, col: 12, offset: 4095},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 205, col: 12, offset: 4095},
	expr: &choiceExpr{
	pos: position{line: 205, col: 13, offset: 4096},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 205, col: 13, offset: 4096},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 205, col: 19, offset: 4102},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 205, col: 25, offset: 4108},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 209, col: 1, offset: 4148},
	expr: &choiceExpr{
	pos: position{line: 209, col: 11, offset: 4160},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 209, col: 11, offset: 4160},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 209, col: 17, offset: 4166},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 209, col: 17, offset: 4166},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 209, col: 37, offset: 4186},
	expr: &ruleRefExpr{
	pos: position{line: 209, col: 37, offset: 4186},
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
	pos: position{line: 211, col: 1, offset: 4201},
	expr: &charClassMatcher{
	pos: position{line: 211, col: 16, offset: 4218},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 212, col: 1, offset: 4224},
	expr: &charClassMatcher{
	pos: position{line: 212, col: 23, offset: 4248},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 214, col: 1, offset: 4255},
	expr: &charClassMatcher{
	pos: position{line: 214, col: 10, offset: 4264},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 215, col: 1, offset: 4270},
	expr: &oneOrMoreExpr{
	pos: position{line: 215, col: 35, offset: 4304},
	expr: &choiceExpr{
	pos: position{line: 215, col: 36, offset: 4305},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 215, col: 36, offset: 4305},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 215, col: 44, offset: 4313},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 215, col: 54, offset: 4323},
	name: "NL",
},
	},
},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 216, col: 1, offset: 4328},
	expr: &zeroOrMoreExpr{
	pos: position{line: 216, col: 20, offset: 4347},
	expr: &choiceExpr{
	pos: position{line: 216, col: 21, offset: 4348},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 216, col: 21, offset: 4348},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 216, col: 29, offset: 4356},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 217, col: 1, offset: 4366},
	expr: &choiceExpr{
	pos: position{line: 217, col: 25, offset: 4390},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 25, offset: 4390},
	name: "NL",
},
&litMatcher{
	pos: position{line: 217, col: 30, offset: 4395},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 217, col: 36, offset: 4401},
	name: "COMMENT",
},
	},
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 218, col: 1, offset: 4410},
	expr: &litMatcher{
	pos: position{line: 218, col: 18, offset: 4427},
	val: "\n",
	ignoreCase: false,
},
},
{
	name: "COMMENT",
	pos: position{line: 220, col: 1, offset: 4433},
	expr: &seqExpr{
	pos: position{line: 220, col: 12, offset: 4444},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 220, col: 12, offset: 4444},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 220, col: 17, offset: 4449},
	expr: &seqExpr{
	pos: position{line: 220, col: 19, offset: 4451},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 220, col: 19, offset: 4451},
	expr: &litMatcher{
	pos: position{line: 220, col: 20, offset: 4452},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 220, col: 25, offset: 4457,
},
	},
},
},
&choiceExpr{
	pos: position{line: 220, col: 31, offset: 4463},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 220, col: 31, offset: 4463},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 220, col: 38, offset: 4470},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 222, col: 1, offset: 4476},
	expr: &notExpr{
	pos: position{line: 222, col: 8, offset: 4483},
	expr: &anyMatcher{
	line: 222, col: 9, offset: 4484,
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

func (c *current) onUSE_ACTION2() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonUSE_ACTION2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE_ACTION2()
}

func (c *current) onUSE_ACTION7() (interface{}, error) {
	return nil, errors.New("cache-control action is deprecated")
}

func (p *parser) callonUSE_ACTION7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE_ACTION7()
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

func (c *current) onMETHOD1() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonMETHOD1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMETHOD1()
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

func (c *current) onFUNCTION2() (interface{}, error) {
	return stringify(c.text)
}

func (p *parser) callonFUNCTION2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFUNCTION2()
}

func (c *current) onFUNCTION7() (interface{}, error) {
	return nil, errors.New("contract function is deprecated")
}

func (p *parser) callonFUNCTION7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFUNCTION7()
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

func (c *current) onFILTER_VALUE1(fv interface{}) (interface{}, error) {
	return newFilterValue(fv)
}

func (p *parser) callonFILTER_VALUE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFILTER_VALUE1(stack["fv"])
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

func (c *current) onFLAGS_RULE1(i, is interface{}) (interface{}, error) {
	return newFlags(i, is)
}

func (p *parser) callonFLAGS_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFLAGS_RULE1(stack["i"], stack["is"])
}

func (c *current) onIGNORE_FLAG1() (interface{}, error) {
	return newIgnoreErrors()
}

func (p *parser) callonIGNORE_FLAG1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIGNORE_FLAG1()
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

