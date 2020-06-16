
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
&ruleRefExpr{
	pos: position{line: 17, col: 30, offset: 147},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 33, offset: 150},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 33, offset: 150},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 37, offset: 154},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 17, col: 40, offset: 157},
	label: "firstBlock",
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 51, offset: 168},
	name: "BLOCK",
},
},
&labeledExpr{
	pos: position{line: 17, col: 57, offset: 174},
	label: "otherBlocks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 17, col: 69, offset: 186},
	expr: &seqExpr{
	pos: position{line: 17, col: 70, offset: 187},
	exprs: []interface{}{
&oneOrMoreExpr{
	pos: position{line: 17, col: 70, offset: 187},
	expr: &seqExpr{
	pos: position{line: 17, col: 71, offset: 188},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 71, offset: 188},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 74, offset: 191},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 77, offset: 194},
	name: "WS",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 82, offset: 199},
	name: "BLOCK",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 90, offset: 207},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 93, offset: 210},
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 93, offset: 210},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 97, offset: 214},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 100, offset: 217},
	name: "EOF",
},
	},
},
},
},
{
	name: "USE",
	pos: position{line: 21, col: 1, offset: 272},
	expr: &actionExpr{
	pos: position{line: 21, col: 8, offset: 279},
	run: (*parser).callonUSE1,
	expr: &seqExpr{
	pos: position{line: 21, col: 8, offset: 279},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 21, col: 8, offset: 279},
	val: "use",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 21, col: 14, offset: 285},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 21, col: 22, offset: 293},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 25, offset: 296},
	name: "USE_ACTION",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 37, offset: 308},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 21, col: 40, offset: 311},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 43, offset: 314},
	name: "USE_VALUE",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 54, offset: 325},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 21, col: 57, offset: 328},
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 57, offset: 328},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 61, offset: 332},
	name: "WS",
},
	},
},
},
},
{
	name: "USE_ACTION",
	pos: position{line: 25, col: 1, offset: 361},
	expr: &actionExpr{
	pos: position{line: 25, col: 15, offset: 375},
	run: (*parser).callonUSE_ACTION1,
	expr: &choiceExpr{
	pos: position{line: 25, col: 16, offset: 376},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 16, offset: 376},
	val: "timeout",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 28, offset: 388},
	val: "max-age",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 40, offset: 400},
	val: "s-max-age",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "USE_VALUE",
	pos: position{line: 29, col: 1, offset: 444},
	expr: &actionExpr{
	pos: position{line: 29, col: 14, offset: 457},
	run: (*parser).callonUSE_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 29, col: 14, offset: 457},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 29, col: 17, offset: 460},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 29, col: 17, offset: 460},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 29, col: 26, offset: 469},
	name: "Integer",
},
	},
},
},
},
},
{
	name: "BLOCK",
	pos: position{line: 33, col: 1, offset: 506},
	expr: &actionExpr{
	pos: position{line: 33, col: 10, offset: 515},
	run: (*parser).callonBLOCK1,
	expr: &seqExpr{
	pos: position{line: 33, col: 10, offset: 515},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 33, col: 10, offset: 515},
	label: "action",
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 18, offset: 523},
	name: "ACTION_RULE",
},
},
&labeledExpr{
	pos: position{line: 33, col: 31, offset: 536},
	label: "m",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 34, offset: 539},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 34, offset: 539},
	name: "MODIFIER_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 50, offset: 555},
	label: "w",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 53, offset: 558},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 53, offset: 558},
	name: "WITH_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 65, offset: 570},
	label: "f",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 67, offset: 572},
	expr: &choiceExpr{
	pos: position{line: 33, col: 68, offset: 573},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 33, col: 68, offset: 573},
	name: "HIDDEN_RULE",
},
&ruleRefExpr{
	pos: position{line: 33, col: 82, offset: 587},
	name: "ONLY_RULE",
},
	},
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 94, offset: 599},
	label: "fl",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 98, offset: 603},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 98, offset: 603},
	name: "FLAGS_RULE",
},
},
},
&ruleRefExpr{
	pos: position{line: 33, col: 111, offset: 616},
	name: "WS",
},
	},
},
},
},
{
	name: "ACTION_RULE",
	pos: position{line: 37, col: 1, offset: 662},
	expr: &actionExpr{
	pos: position{line: 37, col: 16, offset: 677},
	run: (*parser).callonACTION_RULE1,
	expr: &seqExpr{
	pos: position{line: 37, col: 16, offset: 677},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 37, col: 16, offset: 677},
	label: "m",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 19, offset: 680},
	name: "METHOD",
},
},
&ruleRefExpr{
	pos: position{line: 37, col: 27, offset: 688},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 37, col: 35, offset: 696},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 38, offset: 699},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 37, col: 45, offset: 706},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 48, offset: 709},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 48, offset: 709},
	name: "ALIAS",
},
},
},
&labeledExpr{
	pos: position{line: 37, col: 56, offset: 717},
	label: "i",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 59, offset: 720},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 59, offset: 720},
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
	pos: position{line: 41, col: 1, offset: 764},
	expr: &actionExpr{
	pos: position{line: 41, col: 11, offset: 774},
	run: (*parser).callonMETHOD1,
	expr: &choiceExpr{
	pos: position{line: 41, col: 12, offset: 775},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 41, col: 12, offset: 775},
	val: "from",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 21, offset: 784},
	val: "to",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 28, offset: 791},
	val: "into",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 36, offset: 799},
	val: "update",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 47, offset: 810},
	val: "delete",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "ALIAS",
	pos: position{line: 45, col: 1, offset: 851},
	expr: &actionExpr{
	pos: position{line: 45, col: 10, offset: 860},
	run: (*parser).callonALIAS1,
	expr: &seqExpr{
	pos: position{line: 45, col: 10, offset: 860},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 45, col: 10, offset: 860},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 45, col: 18, offset: 868},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 45, col: 23, offset: 873},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 45, col: 31, offset: 881},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 45, col: 34, offset: 884},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IN",
	pos: position{line: 49, col: 1, offset: 911},
	expr: &actionExpr{
	pos: position{line: 49, col: 7, offset: 917},
	run: (*parser).callonIN1,
	expr: &seqExpr{
	pos: position{line: 49, col: 7, offset: 917},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 7, offset: 917},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 49, col: 15, offset: 925},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 49, col: 20, offset: 930},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 49, col: 28, offset: 938},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 49, col: 31, offset: 941},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "MODIFIER_RULE",
	pos: position{line: 53, col: 1, offset: 979},
	expr: &actionExpr{
	pos: position{line: 53, col: 18, offset: 996},
	run: (*parser).callonMODIFIER_RULE1,
	expr: &labeledExpr{
	pos: position{line: 53, col: 18, offset: 996},
	label: "m",
	expr: &oneOrMoreExpr{
	pos: position{line: 53, col: 20, offset: 998},
	expr: &choiceExpr{
	pos: position{line: 53, col: 21, offset: 999},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 53, col: 21, offset: 999},
	name: "HEADERS",
},
&ruleRefExpr{
	pos: position{line: 53, col: 31, offset: 1009},
	name: "TIMEOUT",
},
&ruleRefExpr{
	pos: position{line: 53, col: 41, offset: 1019},
	name: "MAX_AGE",
},
&ruleRefExpr{
	pos: position{line: 53, col: 51, offset: 1029},
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
	pos: position{line: 57, col: 1, offset: 1061},
	expr: &actionExpr{
	pos: position{line: 57, col: 14, offset: 1074},
	run: (*parser).callonWITH_RULE1,
	expr: &seqExpr{
	pos: position{line: 57, col: 14, offset: 1074},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 57, col: 14, offset: 1074},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 57, col: 22, offset: 1082},
	val: "with",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 57, col: 29, offset: 1089},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 57, col: 37, offset: 1097},
	label: "pb",
	expr: &zeroOrOneExpr{
	pos: position{line: 57, col: 40, offset: 1100},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 40, offset: 1100},
	name: "PARAMETER_BODY",
},
},
},
&labeledExpr{
	pos: position{line: 57, col: 56, offset: 1116},
	label: "kvs",
	expr: &zeroOrOneExpr{
	pos: position{line: 57, col: 60, offset: 1120},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 60, offset: 1120},
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
	pos: position{line: 61, col: 1, offset: 1166},
	expr: &actionExpr{
	pos: position{line: 61, col: 19, offset: 1184},
	run: (*parser).callonPARAMETER_BODY1,
	expr: &seqExpr{
	pos: position{line: 61, col: 19, offset: 1184},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 61, col: 19, offset: 1184},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 61, col: 23, offset: 1188},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 26, offset: 1191},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 61, col: 33, offset: 1198},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 61, col: 37, offset: 1202},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 37, offset: 1202},
	name: "APPLY_FN",
},
},
},
&ruleRefExpr{
	pos: position{line: 61, col: 48, offset: 1213},
	name: "WS",
},
&zeroOrOneExpr{
	pos: position{line: 61, col: 51, offset: 1216},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 51, offset: 1216},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 61, col: 55, offset: 1220},
	name: "WS",
},
	},
},
},
},
{
	name: "KEY_VALUE_LIST",
	pos: position{line: 65, col: 1, offset: 1260},
	expr: &actionExpr{
	pos: position{line: 65, col: 19, offset: 1278},
	run: (*parser).callonKEY_VALUE_LIST1,
	expr: &seqExpr{
	pos: position{line: 65, col: 19, offset: 1278},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 65, col: 19, offset: 1278},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 25, offset: 1284},
	name: "KEY_VALUE",
},
},
&labeledExpr{
	pos: position{line: 65, col: 35, offset: 1294},
	label: "others",
	expr: &zeroOrMoreExpr{
	pos: position{line: 65, col: 42, offset: 1301},
	expr: &seqExpr{
	pos: position{line: 65, col: 43, offset: 1302},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 43, offset: 1302},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 46, offset: 1305},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 49, offset: 1308},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 52, offset: 1311},
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
	pos: position{line: 69, col: 1, offset: 1367},
	expr: &actionExpr{
	pos: position{line: 69, col: 14, offset: 1380},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 69, col: 14, offset: 1380},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 69, col: 14, offset: 1380},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 17, offset: 1383},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 69, col: 33, offset: 1399},
	name: "WS",
},
&litMatcher{
	pos: position{line: 69, col: 36, offset: 1402},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 69, col: 40, offset: 1406},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 69, col: 43, offset: 1409},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 46, offset: 1412},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 69, col: 53, offset: 1419},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 69, col: 57, offset: 1423},
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 57, offset: 1423},
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
	pos: position{line: 73, col: 1, offset: 1469},
	expr: &actionExpr{
	pos: position{line: 73, col: 13, offset: 1481},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 73, col: 13, offset: 1481},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 73, col: 13, offset: 1481},
	name: "WS",
},
&litMatcher{
	pos: position{line: 73, col: 16, offset: 1484},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 73, col: 21, offset: 1489},
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 21, offset: 1489},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 73, col: 25, offset: 1493},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 29, offset: 1497},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 77, col: 1, offset: 1528},
	expr: &actionExpr{
	pos: position{line: 77, col: 13, offset: 1540},
	run: (*parser).callonFUNCTION1,
	expr: &choiceExpr{
	pos: position{line: 77, col: 14, offset: 1541},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 77, col: 14, offset: 1541},
	val: "flatten",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 26, offset: 1553},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 37, offset: 1564},
	val: "json",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VALUE",
	pos: position{line: 81, col: 1, offset: 1603},
	expr: &actionExpr{
	pos: position{line: 81, col: 10, offset: 1612},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 81, col: 10, offset: 1612},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 81, col: 13, offset: 1615},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 81, col: 13, offset: 1615},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 81, col: 20, offset: 1622},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 81, col: 29, offset: 1631},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 81, col: 40, offset: 1642},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 85, col: 1, offset: 1678},
	expr: &actionExpr{
	pos: position{line: 85, col: 9, offset: 1686},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 85, col: 9, offset: 1686},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 85, col: 12, offset: 1689},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 85, col: 12, offset: 1689},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 85, col: 25, offset: 1702},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 89, col: 1, offset: 1738},
	expr: &actionExpr{
	pos: position{line: 89, col: 15, offset: 1752},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 89, col: 15, offset: 1752},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 15, offset: 1752},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 89, col: 19, offset: 1756},
	name: "WS",
},
&litMatcher{
	pos: position{line: 89, col: 22, offset: 1759},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 93, col: 1, offset: 1791},
	expr: &actionExpr{
	pos: position{line: 93, col: 19, offset: 1809},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 93, col: 19, offset: 1809},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 93, col: 19, offset: 1809},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 93, col: 23, offset: 1813},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 93, col: 26, offset: 1816},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 28, offset: 1818},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 93, col: 34, offset: 1824},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 93, col: 37, offset: 1827},
	expr: &seqExpr{
	pos: position{line: 93, col: 38, offset: 1828},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 93, col: 38, offset: 1828},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 93, col: 41, offset: 1831},
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 41, offset: 1831},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 45, offset: 1835},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 93, col: 48, offset: 1838},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 56, offset: 1846},
	name: "WS",
},
&litMatcher{
	pos: position{line: 93, col: 59, offset: 1849},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 97, col: 1, offset: 1881},
	expr: &actionExpr{
	pos: position{line: 97, col: 11, offset: 1891},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 97, col: 11, offset: 1891},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 97, col: 14, offset: 1894},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 97, col: 14, offset: 1894},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 97, col: 26, offset: 1906},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 101, col: 1, offset: 1941},
	expr: &actionExpr{
	pos: position{line: 101, col: 14, offset: 1954},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 101, col: 14, offset: 1954},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 101, col: 14, offset: 1954},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 101, col: 18, offset: 1958},
	name: "WS",
},
&litMatcher{
	pos: position{line: 101, col: 21, offset: 1961},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 105, col: 1, offset: 1995},
	expr: &actionExpr{
	pos: position{line: 105, col: 18, offset: 2012},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 105, col: 18, offset: 2012},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 105, col: 18, offset: 2012},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 22, offset: 2016},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 105, col: 25, offset: 2019},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 29, offset: 2023},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 105, col: 40, offset: 2034},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 105, col: 44, offset: 2038},
	expr: &seqExpr{
	pos: position{line: 105, col: 45, offset: 2039},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 105, col: 45, offset: 2039},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 48, offset: 2042},
	val: ",",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 52, offset: 2046},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 52, offset: 2046},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 56, offset: 2050},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 105, col: 59, offset: 2053},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 71, offset: 2065},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 74, offset: 2068},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 109, col: 1, offset: 2113},
	expr: &actionExpr{
	pos: position{line: 109, col: 14, offset: 2126},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 109, col: 14, offset: 2126},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 109, col: 14, offset: 2126},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 109, col: 17, offset: 2129},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 109, col: 17, offset: 2129},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 109, col: 26, offset: 2138},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 33, offset: 2145},
	name: "WS",
},
&litMatcher{
	pos: position{line: 109, col: 36, offset: 2148},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 109, col: 40, offset: 2152},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 109, col: 43, offset: 2155},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 46, offset: 2158},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 113, col: 1, offset: 2199},
	expr: &actionExpr{
	pos: position{line: 113, col: 14, offset: 2212},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 113, col: 14, offset: 2212},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 113, col: 17, offset: 2215},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2215},
	name: "Null",
},
&ruleRefExpr{
	pos: position{line: 113, col: 24, offset: 2222},
	name: "Boolean",
},
&ruleRefExpr{
	pos: position{line: 113, col: 34, offset: 2232},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 113, col: 43, offset: 2241},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 113, col: 51, offset: 2249},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 113, col: 61, offset: 2259},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "CHAIN",
	pos: position{line: 117, col: 1, offset: 2295},
	expr: &actionExpr{
	pos: position{line: 117, col: 10, offset: 2304},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 117, col: 10, offset: 2304},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 117, col: 10, offset: 2304},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 117, col: 13, offset: 2307},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 117, col: 27, offset: 2321},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 117, col: 30, offset: 2324},
	expr: &seqExpr{
	pos: position{line: 117, col: 31, offset: 2325},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 117, col: 31, offset: 2325},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 117, col: 35, offset: 2329},
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
	pos: position{line: 121, col: 1, offset: 2373},
	expr: &actionExpr{
	pos: position{line: 121, col: 17, offset: 2389},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 121, col: 17, offset: 2389},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 121, col: 21, offset: 2393},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 121, col: 21, offset: 2393},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 121, col: 32, offset: 2404},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 125, col: 1, offset: 2439},
	expr: &actionExpr{
	pos: position{line: 125, col: 14, offset: 2452},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 125, col: 14, offset: 2452},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 14, offset: 2452},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 125, col: 22, offset: 2460},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 125, col: 29, offset: 2467},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 125, col: 37, offset: 2475},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 125, col: 40, offset: 2478},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 125, col: 48, offset: 2486},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 125, col: 51, offset: 2489},
	expr: &seqExpr{
	pos: position{line: 125, col: 52, offset: 2490},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 52, offset: 2490},
	name: "WS",
},
&notExpr{
	pos: position{line: 125, col: 55, offset: 2493},
	expr: &choiceExpr{
	pos: position{line: 125, col: 57, offset: 2495},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 57, offset: 2495},
	name: "FLAGS_RULE",
},
&seqExpr{
	pos: position{line: 125, col: 70, offset: 2508},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 70, offset: 2508},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 73, offset: 2511},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 76, offset: 2514},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 79, offset: 2517},
	name: "BLOCK",
},
	},
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 125, col: 86, offset: 2524},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 89, offset: 2527},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 92, offset: 2530},
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
	pos: position{line: 129, col: 1, offset: 2567},
	expr: &actionExpr{
	pos: position{line: 129, col: 11, offset: 2577},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 129, col: 11, offset: 2577},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 129, col: 11, offset: 2577},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 14, offset: 2580},
	name: "IDENT_WITH_DOT",
},
},
&labeledExpr{
	pos: position{line: 129, col: 30, offset: 2596},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 129, col: 34, offset: 2600},
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 34, offset: 2600},
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
	pos: position{line: 133, col: 1, offset: 2643},
	expr: &actionExpr{
	pos: position{line: 133, col: 15, offset: 2657},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 133, col: 15, offset: 2657},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 133, col: 15, offset: 2657},
	name: "WS",
},
&litMatcher{
	pos: position{line: 133, col: 18, offset: 2660},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 133, col: 23, offset: 2665},
	name: "WS",
},
&litMatcher{
	pos: position{line: 133, col: 26, offset: 2668},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 133, col: 36, offset: 2678},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 133, col: 40, offset: 2682},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 133, col: 45, offset: 2687},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 133, col: 53, offset: 2695},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 137, col: 1, offset: 2736},
	expr: &actionExpr{
	pos: position{line: 137, col: 12, offset: 2747},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 137, col: 12, offset: 2747},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 12, offset: 2747},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 137, col: 20, offset: 2755},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 137, col: 30, offset: 2765},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 137, col: 38, offset: 2773},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 137, col: 41, offset: 2776},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 137, col: 49, offset: 2784},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 137, col: 52, offset: 2787},
	expr: &seqExpr{
	pos: position{line: 137, col: 53, offset: 2788},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 53, offset: 2788},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 137, col: 56, offset: 2791},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 137, col: 59, offset: 2794},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 137, col: 62, offset: 2797},
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
	pos: position{line: 141, col: 1, offset: 2837},
	expr: &actionExpr{
	pos: position{line: 141, col: 11, offset: 2847},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 141, col: 11, offset: 2847},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 141, col: 11, offset: 2847},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 141, col: 14, offset: 2850},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 141, col: 21, offset: 2857},
	name: "WS",
},
&litMatcher{
	pos: position{line: 141, col: 24, offset: 2860},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 141, col: 28, offset: 2864},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 141, col: 31, offset: 2867},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 141, col: 34, offset: 2870},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 141, col: 34, offset: 2870},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 141, col: 45, offset: 2881},
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
	pos: position{line: 145, col: 1, offset: 2918},
	expr: &actionExpr{
	pos: position{line: 145, col: 16, offset: 2933},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 145, col: 16, offset: 2933},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 145, col: 16, offset: 2933},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 145, col: 24, offset: 2941},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 149, col: 1, offset: 2975},
	expr: &actionExpr{
	pos: position{line: 149, col: 12, offset: 2986},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 149, col: 12, offset: 2986},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 12, offset: 2986},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 149, col: 20, offset: 2994},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 149, col: 30, offset: 3004},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 149, col: 38, offset: 3012},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 149, col: 41, offset: 3015},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 41, offset: 3015},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 149, col: 52, offset: 3026},
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
	pos: position{line: 153, col: 1, offset: 3062},
	expr: &actionExpr{
	pos: position{line: 153, col: 12, offset: 3073},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 153, col: 12, offset: 3073},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 12, offset: 3073},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 153, col: 20, offset: 3081},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 153, col: 30, offset: 3091},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 153, col: 38, offset: 3099},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 153, col: 41, offset: 3102},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 41, offset: 3102},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 153, col: 52, offset: 3113},
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
	pos: position{line: 157, col: 1, offset: 3148},
	expr: &actionExpr{
	pos: position{line: 157, col: 14, offset: 3161},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 157, col: 14, offset: 3161},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 14, offset: 3161},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 157, col: 22, offset: 3169},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 34, offset: 3181},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 157, col: 42, offset: 3189},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 157, col: 45, offset: 3192},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 45, offset: 3192},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 157, col: 56, offset: 3203},
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
	pos: position{line: 161, col: 1, offset: 3239},
	expr: &actionExpr{
	pos: position{line: 161, col: 15, offset: 3253},
	run: (*parser).callonFLAGS_RULE1,
	expr: &seqExpr{
	pos: position{line: 161, col: 15, offset: 3253},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 15, offset: 3253},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 161, col: 23, offset: 3261},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 161, col: 25, offset: 3263},
	name: "IGNORE_FLAG",
},
},
&labeledExpr{
	pos: position{line: 161, col: 37, offset: 3275},
	label: "is",
	expr: &zeroOrMoreExpr{
	pos: position{line: 161, col: 40, offset: 3278},
	expr: &seqExpr{
	pos: position{line: 161, col: 41, offset: 3279},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 41, offset: 3279},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 161, col: 44, offset: 3282},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 161, col: 47, offset: 3285},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 161, col: 50, offset: 3288},
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
	pos: position{line: 165, col: 1, offset: 3331},
	expr: &actionExpr{
	pos: position{line: 165, col: 16, offset: 3346},
	run: (*parser).callonIGNORE_FLAG1,
	expr: &litMatcher{
	pos: position{line: 165, col: 16, offset: 3346},
	val: "ignore-errors",
	ignoreCase: false,
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 169, col: 1, offset: 3393},
	expr: &actionExpr{
	pos: position{line: 169, col: 13, offset: 3405},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 169, col: 13, offset: 3405},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 169, col: 13, offset: 3405},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 169, col: 17, offset: 3409},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 169, col: 20, offset: 3412},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 173, col: 1, offset: 3447},
	expr: &actionExpr{
	pos: position{line: 173, col: 10, offset: 3456},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 173, col: 10, offset: 3456},
	expr: &charClassMatcher{
	pos: position{line: 173, col: 10, offset: 3456},
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
	pos: position{line: 177, col: 1, offset: 3502},
	expr: &actionExpr{
	pos: position{line: 177, col: 19, offset: 3520},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 177, col: 19, offset: 3520},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 19, offset: 3520},
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
	pos: position{line: 181, col: 1, offset: 3567},
	expr: &actionExpr{
	pos: position{line: 181, col: 9, offset: 3575},
	run: (*parser).callonNull1,
	expr: &litMatcher{
	pos: position{line: 181, col: 9, offset: 3575},
	val: "null",
	ignoreCase: false,
},
},
},
{
	name: "Boolean",
	pos: position{line: 185, col: 1, offset: 3605},
	expr: &actionExpr{
	pos: position{line: 185, col: 12, offset: 3616},
	run: (*parser).callonBoolean1,
	expr: &choiceExpr{
	pos: position{line: 185, col: 13, offset: 3617},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 185, col: 13, offset: 3617},
	val: "true",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 185, col: 22, offset: 3626},
	val: "false",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "String",
	pos: position{line: 189, col: 1, offset: 3667},
	expr: &actionExpr{
	pos: position{line: 189, col: 11, offset: 3677},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 189, col: 11, offset: 3677},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 189, col: 11, offset: 3677},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 189, col: 15, offset: 3681},
	expr: &seqExpr{
	pos: position{line: 189, col: 17, offset: 3683},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 189, col: 17, offset: 3683},
	expr: &litMatcher{
	pos: position{line: 189, col: 18, offset: 3684},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 189, col: 22, offset: 3688,
},
	},
},
},
&litMatcher{
	pos: position{line: 189, col: 27, offset: 3693},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 193, col: 1, offset: 3728},
	expr: &actionExpr{
	pos: position{line: 193, col: 10, offset: 3737},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 193, col: 10, offset: 3737},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 193, col: 10, offset: 3737},
	expr: &choiceExpr{
	pos: position{line: 193, col: 11, offset: 3738},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 193, col: 11, offset: 3738},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 193, col: 17, offset: 3744},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 193, col: 23, offset: 3750},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 193, col: 31, offset: 3758},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 193, col: 35, offset: 3762},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 197, col: 1, offset: 3800},
	expr: &actionExpr{
	pos: position{line: 197, col: 12, offset: 3811},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 197, col: 12, offset: 3811},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 197, col: 12, offset: 3811},
	expr: &choiceExpr{
	pos: position{line: 197, col: 13, offset: 3812},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 197, col: 13, offset: 3812},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 197, col: 19, offset: 3818},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 197, col: 25, offset: 3824},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 201, col: 1, offset: 3864},
	expr: &choiceExpr{
	pos: position{line: 201, col: 11, offset: 3876},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 201, col: 11, offset: 3876},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 201, col: 17, offset: 3882},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 201, col: 17, offset: 3882},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 201, col: 37, offset: 3902},
	expr: &ruleRefExpr{
	pos: position{line: 201, col: 37, offset: 3902},
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
	pos: position{line: 203, col: 1, offset: 3917},
	expr: &charClassMatcher{
	pos: position{line: 203, col: 16, offset: 3934},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 204, col: 1, offset: 3940},
	expr: &charClassMatcher{
	pos: position{line: 204, col: 23, offset: 3964},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 206, col: 1, offset: 3971},
	expr: &charClassMatcher{
	pos: position{line: 206, col: 10, offset: 3980},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 207, col: 1, offset: 3986},
	expr: &oneOrMoreExpr{
	pos: position{line: 207, col: 35, offset: 4020},
	expr: &choiceExpr{
	pos: position{line: 207, col: 36, offset: 4021},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 207, col: 36, offset: 4021},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 207, col: 44, offset: 4029},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 207, col: 54, offset: 4039},
	name: "NL",
},
	},
},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 208, col: 1, offset: 4044},
	expr: &zeroOrMoreExpr{
	pos: position{line: 208, col: 20, offset: 4063},
	expr: &choiceExpr{
	pos: position{line: 208, col: 21, offset: 4064},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 208, col: 21, offset: 4064},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 208, col: 29, offset: 4072},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 209, col: 1, offset: 4082},
	expr: &choiceExpr{
	pos: position{line: 209, col: 25, offset: 4106},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 209, col: 25, offset: 4106},
	name: "NL",
},
&litMatcher{
	pos: position{line: 209, col: 30, offset: 4111},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 209, col: 36, offset: 4117},
	name: "COMMENT",
},
	},
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 210, col: 1, offset: 4126},
	expr: &litMatcher{
	pos: position{line: 210, col: 18, offset: 4143},
	val: "\n",
	ignoreCase: false,
},
},
{
	name: "COMMENT",
	pos: position{line: 212, col: 1, offset: 4149},
	expr: &seqExpr{
	pos: position{line: 212, col: 12, offset: 4160},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 212, col: 12, offset: 4160},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 212, col: 17, offset: 4165},
	expr: &seqExpr{
	pos: position{line: 212, col: 19, offset: 4167},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 212, col: 19, offset: 4167},
	expr: &litMatcher{
	pos: position{line: 212, col: 20, offset: 4168},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 212, col: 25, offset: 4173,
},
	},
},
},
&choiceExpr{
	pos: position{line: 212, col: 31, offset: 4179},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 212, col: 31, offset: 4179},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 212, col: 38, offset: 4186},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 214, col: 1, offset: 4192},
	expr: &notExpr{
	pos: position{line: 214, col: 8, offset: 4199},
	expr: &anyMatcher{
	line: 214, col: 9, offset: 4200,
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

