
package ast

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar {
	rules: []*rule{
{
	name: "QUERY",
	pos: position{line: 18, col: 1, offset: 129},
	expr: &actionExpr{
	pos: position{line: 18, col: 10, offset: 138},
	run: (*parser).callonQUERY1,
	expr: &seqExpr{
	pos: position{line: 18, col: 10, offset: 138},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 18, col: 10, offset: 138},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 18, col: 13, offset: 141},
	expr: &ruleRefExpr{
	pos: position{line: 18, col: 13, offset: 141},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 18, col: 17, offset: 145},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 18, col: 20, offset: 148},
	label: "us",
	expr: &zeroOrMoreExpr{
	pos: position{line: 18, col: 23, offset: 151},
	expr: &ruleRefExpr{
	pos: position{line: 18, col: 24, offset: 152},
	name: "USE",
},
},
},
&labeledExpr{
	pos: position{line: 18, col: 30, offset: 158},
	label: "firstBlock",
	expr: &ruleRefExpr{
	pos: position{line: 18, col: 41, offset: 169},
	name: "BLOCK",
},
},
&labeledExpr{
	pos: position{line: 18, col: 47, offset: 175},
	label: "otherBlocks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 18, col: 59, offset: 187},
	expr: &seqExpr{
	pos: position{line: 18, col: 60, offset: 188},
	exprs: []interface{}{
&oneOrMoreExpr{
	pos: position{line: 18, col: 60, offset: 188},
	expr: &seqExpr{
	pos: position{line: 18, col: 61, offset: 189},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 18, col: 61, offset: 189},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 18, col: 64, offset: 192},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 18, col: 67, offset: 195},
	name: "WS",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 18, col: 72, offset: 200},
	name: "BLOCK",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 18, col: 80, offset: 208},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 18, col: 83, offset: 211},
	expr: &ruleRefExpr{
	pos: position{line: 18, col: 83, offset: 211},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 18, col: 87, offset: 215},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 18, col: 90, offset: 218},
	name: "EOF",
},
	},
},
},
},
{
	name: "USE",
	pos: position{line: 22, col: 1, offset: 273},
	expr: &actionExpr{
	pos: position{line: 22, col: 8, offset: 280},
	run: (*parser).callonUSE1,
	expr: &seqExpr{
	pos: position{line: 22, col: 8, offset: 280},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 22, col: 8, offset: 280},
	val: "use",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 22, col: 14, offset: 286},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 22, col: 22, offset: 294},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 22, col: 25, offset: 297},
	name: "USE_RULE",
},
},
&ruleRefExpr{
	pos: position{line: 22, col: 35, offset: 307},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 22, col: 38, offset: 310},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 22, col: 41, offset: 313},
	name: "USE_VALUE",
},
},
&ruleRefExpr{
	pos: position{line: 22, col: 52, offset: 324},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 22, col: 55, offset: 327},
	expr: &ruleRefExpr{
	pos: position{line: 22, col: 55, offset: 327},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 22, col: 59, offset: 331},
	name: "WS",
},
	},
},
},
},
{
	name: "USE_RULE",
	pos: position{line: 26, col: 1, offset: 360},
	expr: &actionExpr{
	pos: position{line: 26, col: 13, offset: 372},
	run: (*parser).callonUSE_RULE1,
	expr: &choiceExpr{
	pos: position{line: 26, col: 14, offset: 373},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 26, col: 14, offset: 373},
	val: "timeout",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 26, col: 26, offset: 385},
	val: "max-age",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 26, col: 38, offset: 397},
	val: "s-max-age",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "USE_VALUE",
	pos: position{line: 30, col: 1, offset: 443},
	expr: &actionExpr{
	pos: position{line: 30, col: 14, offset: 456},
	run: (*parser).callonUSE_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 30, col: 14, offset: 456},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 30, col: 17, offset: 459},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 30, col: 17, offset: 459},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 30, col: 26, offset: 468},
	name: "Integer",
},
	},
},
},
},
},
{
	name: "BLOCK",
	pos: position{line: 34, col: 1, offset: 505},
	expr: &actionExpr{
	pos: position{line: 34, col: 10, offset: 514},
	run: (*parser).callonBLOCK1,
	expr: &seqExpr{
	pos: position{line: 34, col: 10, offset: 514},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 34, col: 10, offset: 514},
	label: "action",
	expr: &ruleRefExpr{
	pos: position{line: 34, col: 18, offset: 522},
	name: "ACTION_RULE",
},
},
&labeledExpr{
	pos: position{line: 34, col: 31, offset: 535},
	label: "m",
	expr: &zeroOrOneExpr{
	pos: position{line: 34, col: 34, offset: 538},
	expr: &ruleRefExpr{
	pos: position{line: 34, col: 34, offset: 538},
	name: "MODIFIER_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 34, col: 50, offset: 554},
	label: "w",
	expr: &zeroOrOneExpr{
	pos: position{line: 34, col: 53, offset: 557},
	expr: &ruleRefExpr{
	pos: position{line: 34, col: 53, offset: 557},
	name: "WITH_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 34, col: 65, offset: 569},
	label: "f",
	expr: &zeroOrOneExpr{
	pos: position{line: 34, col: 67, offset: 571},
	expr: &choiceExpr{
	pos: position{line: 34, col: 68, offset: 572},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 34, col: 68, offset: 572},
	name: "ONLY_RULE",
},
&ruleRefExpr{
	pos: position{line: 34, col: 80, offset: 584},
	name: "HIDDEN_RULE",
},
	},
},
},
},
&labeledExpr{
	pos: position{line: 34, col: 94, offset: 598},
	label: "fl",
	expr: &zeroOrOneExpr{
	pos: position{line: 34, col: 98, offset: 602},
	expr: &ruleRefExpr{
	pos: position{line: 34, col: 98, offset: 602},
	name: "FLAG_RULE",
},
},
},
&ruleRefExpr{
	pos: position{line: 34, col: 110, offset: 614},
	name: "WS",
},
	},
},
},
},
{
	name: "ACTION_RULE",
	pos: position{line: 38, col: 1, offset: 660},
	expr: &actionExpr{
	pos: position{line: 38, col: 16, offset: 675},
	run: (*parser).callonACTION_RULE1,
	expr: &seqExpr{
	pos: position{line: 38, col: 16, offset: 675},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 38, col: 16, offset: 675},
	label: "m",
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 19, offset: 678},
	name: "METHOD",
},
},
&ruleRefExpr{
	pos: position{line: 38, col: 27, offset: 686},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 38, col: 35, offset: 694},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 38, offset: 697},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 38, col: 45, offset: 704},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 38, col: 48, offset: 707},
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 48, offset: 707},
	name: "ALIAS",
},
},
},
&labeledExpr{
	pos: position{line: 38, col: 56, offset: 715},
	label: "i",
	expr: &zeroOrOneExpr{
	pos: position{line: 38, col: 59, offset: 718},
	expr: &ruleRefExpr{
	pos: position{line: 38, col: 59, offset: 718},
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
	pos: position{line: 42, col: 1, offset: 762},
	expr: &actionExpr{
	pos: position{line: 42, col: 11, offset: 772},
	run: (*parser).callonMETHOD1,
	expr: &labeledExpr{
	pos: position{line: 42, col: 11, offset: 772},
	label: "m",
	expr: &choiceExpr{
	pos: position{line: 42, col: 14, offset: 775},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 42, col: 14, offset: 775},
	val: "from",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 42, col: 23, offset: 784},
	val: "to",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 42, col: 30, offset: 791},
	val: "into",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 42, col: 38, offset: 799},
	val: "update",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 42, col: 49, offset: 810},
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
	pos: position{line: 46, col: 1, offset: 851},
	expr: &actionExpr{
	pos: position{line: 46, col: 10, offset: 860},
	run: (*parser).callonALIAS1,
	expr: &seqExpr{
	pos: position{line: 46, col: 10, offset: 860},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 46, col: 10, offset: 860},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 46, col: 18, offset: 868},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 46, col: 23, offset: 873},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 46, col: 31, offset: 881},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 46, col: 34, offset: 884},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IN",
	pos: position{line: 50, col: 1, offset: 911},
	expr: &actionExpr{
	pos: position{line: 50, col: 7, offset: 917},
	run: (*parser).callonIN1,
	expr: &seqExpr{
	pos: position{line: 50, col: 7, offset: 917},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 50, col: 7, offset: 917},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 50, col: 15, offset: 925},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 50, col: 20, offset: 930},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 50, col: 28, offset: 938},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 50, col: 31, offset: 941},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "MODIFIER_RULE",
	pos: position{line: 54, col: 1, offset: 979},
	expr: &actionExpr{
	pos: position{line: 54, col: 18, offset: 996},
	run: (*parser).callonMODIFIER_RULE1,
	expr: &labeledExpr{
	pos: position{line: 54, col: 18, offset: 996},
	label: "m",
	expr: &oneOrMoreExpr{
	pos: position{line: 54, col: 20, offset: 998},
	expr: &choiceExpr{
	pos: position{line: 54, col: 21, offset: 999},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 54, col: 21, offset: 999},
	name: "HEADERS",
},
&ruleRefExpr{
	pos: position{line: 54, col: 31, offset: 1009},
	name: "TIMEOUT",
},
&ruleRefExpr{
	pos: position{line: 54, col: 41, offset: 1019},
	name: "MAX_AGE",
},
&ruleRefExpr{
	pos: position{line: 54, col: 51, offset: 1029},
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
	pos: position{line: 58, col: 1, offset: 1061},
	expr: &actionExpr{
	pos: position{line: 58, col: 14, offset: 1074},
	run: (*parser).callonWITH_RULE1,
	expr: &seqExpr{
	pos: position{line: 58, col: 14, offset: 1074},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 58, col: 14, offset: 1074},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 58, col: 22, offset: 1082},
	val: "with",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 58, col: 29, offset: 1089},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 58, col: 37, offset: 1097},
	label: "pb",
	expr: &zeroOrOneExpr{
	pos: position{line: 58, col: 40, offset: 1100},
	expr: &ruleRefExpr{
	pos: position{line: 58, col: 40, offset: 1100},
	name: "PARAMETER_BODY",
},
},
},
&labeledExpr{
	pos: position{line: 58, col: 56, offset: 1116},
	label: "kvs",
	expr: &zeroOrOneExpr{
	pos: position{line: 58, col: 60, offset: 1120},
	expr: &ruleRefExpr{
	pos: position{line: 58, col: 60, offset: 1120},
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
	pos: position{line: 62, col: 1, offset: 1166},
	expr: &actionExpr{
	pos: position{line: 62, col: 19, offset: 1184},
	run: (*parser).callonPARAMETER_BODY1,
	expr: &seqExpr{
	pos: position{line: 62, col: 19, offset: 1184},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 62, col: 19, offset: 1184},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 62, col: 23, offset: 1188},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 62, col: 26, offset: 1191},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 62, col: 33, offset: 1198},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 62, col: 37, offset: 1202},
	expr: &ruleRefExpr{
	pos: position{line: 62, col: 37, offset: 1202},
	name: "APPLY_FN",
},
},
},
&ruleRefExpr{
	pos: position{line: 62, col: 48, offset: 1213},
	name: "WS",
},
&zeroOrOneExpr{
	pos: position{line: 62, col: 51, offset: 1216},
	expr: &ruleRefExpr{
	pos: position{line: 62, col: 51, offset: 1216},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 62, col: 55, offset: 1220},
	name: "WS",
},
	},
},
},
},
{
	name: "KEY_VALUE_LIST",
	pos: position{line: 66, col: 1, offset: 1260},
	expr: &actionExpr{
	pos: position{line: 66, col: 19, offset: 1278},
	run: (*parser).callonKEY_VALUE_LIST1,
	expr: &seqExpr{
	pos: position{line: 66, col: 19, offset: 1278},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 66, col: 19, offset: 1278},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 66, col: 25, offset: 1284},
	name: "KEY_VALUE",
},
},
&labeledExpr{
	pos: position{line: 66, col: 35, offset: 1294},
	label: "others",
	expr: &zeroOrMoreExpr{
	pos: position{line: 66, col: 42, offset: 1301},
	expr: &seqExpr{
	pos: position{line: 66, col: 43, offset: 1302},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 66, col: 43, offset: 1302},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 66, col: 46, offset: 1305},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 66, col: 49, offset: 1308},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 66, col: 52, offset: 1311},
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
	pos: position{line: 70, col: 1, offset: 1367},
	expr: &actionExpr{
	pos: position{line: 70, col: 14, offset: 1380},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 70, col: 14, offset: 1380},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 70, col: 14, offset: 1380},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 70, col: 17, offset: 1383},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 70, col: 33, offset: 1399},
	name: "WS",
},
&litMatcher{
	pos: position{line: 70, col: 36, offset: 1402},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 70, col: 40, offset: 1406},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 70, col: 43, offset: 1409},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 70, col: 46, offset: 1412},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 70, col: 53, offset: 1419},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 70, col: 57, offset: 1423},
	expr: &ruleRefExpr{
	pos: position{line: 70, col: 57, offset: 1423},
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
	pos: position{line: 74, col: 1, offset: 1469},
	expr: &actionExpr{
	pos: position{line: 74, col: 13, offset: 1481},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 74, col: 13, offset: 1481},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 74, col: 13, offset: 1481},
	name: "WS",
},
&litMatcher{
	pos: position{line: 74, col: 16, offset: 1484},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 74, col: 21, offset: 1489},
	expr: &ruleRefExpr{
	pos: position{line: 74, col: 21, offset: 1489},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 74, col: 25, offset: 1493},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 74, col: 29, offset: 1497},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 78, col: 1, offset: 1528},
	expr: &actionExpr{
	pos: position{line: 78, col: 13, offset: 1540},
	run: (*parser).callonFUNCTION1,
	expr: &choiceExpr{
	pos: position{line: 78, col: 14, offset: 1541},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 78, col: 14, offset: 1541},
	val: "flatten",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 78, col: 26, offset: 1553},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 78, col: 37, offset: 1564},
	val: "json",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VALUE",
	pos: position{line: 82, col: 1, offset: 1605},
	expr: &actionExpr{
	pos: position{line: 82, col: 10, offset: 1614},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 82, col: 10, offset: 1614},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 82, col: 13, offset: 1617},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 82, col: 13, offset: 1617},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 82, col: 20, offset: 1624},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 82, col: 29, offset: 1633},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 82, col: 40, offset: 1644},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 86, col: 1, offset: 1680},
	expr: &actionExpr{
	pos: position{line: 86, col: 9, offset: 1688},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 86, col: 9, offset: 1688},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 86, col: 12, offset: 1691},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 86, col: 12, offset: 1691},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 86, col: 25, offset: 1704},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 90, col: 1, offset: 1740},
	expr: &actionExpr{
	pos: position{line: 90, col: 15, offset: 1754},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 90, col: 15, offset: 1754},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 90, col: 15, offset: 1754},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 90, col: 19, offset: 1758},
	name: "WS",
},
&litMatcher{
	pos: position{line: 90, col: 22, offset: 1761},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 94, col: 1, offset: 1793},
	expr: &actionExpr{
	pos: position{line: 94, col: 19, offset: 1811},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 94, col: 19, offset: 1811},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 94, col: 19, offset: 1811},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 94, col: 23, offset: 1815},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 94, col: 26, offset: 1818},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 94, col: 28, offset: 1820},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 94, col: 34, offset: 1826},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 94, col: 37, offset: 1829},
	expr: &seqExpr{
	pos: position{line: 94, col: 38, offset: 1830},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 94, col: 38, offset: 1830},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 94, col: 41, offset: 1833},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 94, col: 44, offset: 1836},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 94, col: 47, offset: 1839},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 94, col: 55, offset: 1847},
	name: "WS",
},
&litMatcher{
	pos: position{line: 94, col: 58, offset: 1850},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 98, col: 1, offset: 1882},
	expr: &actionExpr{
	pos: position{line: 98, col: 11, offset: 1892},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 98, col: 11, offset: 1892},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 98, col: 14, offset: 1895},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 98, col: 14, offset: 1895},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 98, col: 26, offset: 1907},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 102, col: 1, offset: 1942},
	expr: &actionExpr{
	pos: position{line: 102, col: 14, offset: 1955},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 102, col: 14, offset: 1955},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 102, col: 14, offset: 1955},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 102, col: 18, offset: 1959},
	name: "WS",
},
&litMatcher{
	pos: position{line: 102, col: 21, offset: 1962},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 106, col: 1, offset: 1996},
	expr: &actionExpr{
	pos: position{line: 106, col: 18, offset: 2013},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 106, col: 18, offset: 2013},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 106, col: 18, offset: 2013},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 106, col: 22, offset: 2017},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 106, col: 25, offset: 2020},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 106, col: 29, offset: 2024},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 106, col: 40, offset: 2035},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 106, col: 44, offset: 2039},
	expr: &seqExpr{
	pos: position{line: 106, col: 45, offset: 2040},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 106, col: 45, offset: 2040},
	name: "WS",
},
&litMatcher{
	pos: position{line: 106, col: 48, offset: 2043},
	val: ",",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 106, col: 52, offset: 2047},
	expr: &ruleRefExpr{
	pos: position{line: 106, col: 52, offset: 2047},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 106, col: 56, offset: 2051},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 106, col: 59, offset: 2054},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 106, col: 71, offset: 2066},
	name: "WS",
},
&litMatcher{
	pos: position{line: 106, col: 74, offset: 2069},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 110, col: 1, offset: 2114},
	expr: &actionExpr{
	pos: position{line: 110, col: 14, offset: 2127},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 110, col: 14, offset: 2127},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 110, col: 14, offset: 2127},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 110, col: 17, offset: 2130},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 110, col: 17, offset: 2130},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 110, col: 26, offset: 2139},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 110, col: 33, offset: 2146},
	name: "WS",
},
&litMatcher{
	pos: position{line: 110, col: 36, offset: 2149},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 110, col: 40, offset: 2153},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 110, col: 43, offset: 2156},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 110, col: 46, offset: 2159},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 114, col: 1, offset: 2200},
	expr: &actionExpr{
	pos: position{line: 114, col: 14, offset: 2213},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 114, col: 14, offset: 2213},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 114, col: 17, offset: 2216},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 114, col: 17, offset: 2216},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 114, col: 26, offset: 2225},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 114, col: 34, offset: 2233},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 114, col: 44, offset: 2243},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "CHAIN",
	pos: position{line: 118, col: 1, offset: 2279},
	expr: &actionExpr{
	pos: position{line: 118, col: 10, offset: 2288},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 118, col: 10, offset: 2288},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 118, col: 10, offset: 2288},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 118, col: 13, offset: 2291},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 118, col: 27, offset: 2305},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 118, col: 30, offset: 2308},
	expr: &seqExpr{
	pos: position{line: 118, col: 31, offset: 2309},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 118, col: 31, offset: 2309},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 118, col: 35, offset: 2313},
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
	pos: position{line: 122, col: 1, offset: 2357},
	expr: &actionExpr{
	pos: position{line: 122, col: 17, offset: 2373},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 122, col: 17, offset: 2373},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 122, col: 21, offset: 2377},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 122, col: 21, offset: 2377},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 122, col: 32, offset: 2388},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 126, col: 1, offset: 2423},
	expr: &actionExpr{
	pos: position{line: 126, col: 14, offset: 2436},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 126, col: 14, offset: 2436},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 126, col: 14, offset: 2436},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 126, col: 22, offset: 2444},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 126, col: 29, offset: 2451},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 126, col: 37, offset: 2459},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 126, col: 40, offset: 2462},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 126, col: 48, offset: 2470},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 126, col: 51, offset: 2473},
	expr: &seqExpr{
	pos: position{line: 126, col: 52, offset: 2474},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 126, col: 52, offset: 2474},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 126, col: 55, offset: 2477},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 126, col: 58, offset: 2480},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 126, col: 61, offset: 2483},
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
	pos: position{line: 130, col: 1, offset: 2520},
	expr: &actionExpr{
	pos: position{line: 130, col: 11, offset: 2530},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 130, col: 11, offset: 2530},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 130, col: 11, offset: 2530},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 130, col: 14, offset: 2533},
	name: "IDENT_WITH_DOT",
},
},
&labeledExpr{
	pos: position{line: 130, col: 30, offset: 2549},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 130, col: 34, offset: 2553},
	expr: &ruleRefExpr{
	pos: position{line: 130, col: 34, offset: 2553},
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
	pos: position{line: 134, col: 1, offset: 2596},
	expr: &actionExpr{
	pos: position{line: 134, col: 15, offset: 2610},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 134, col: 15, offset: 2610},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 134, col: 15, offset: 2610},
	name: "WS",
},
&litMatcher{
	pos: position{line: 134, col: 18, offset: 2613},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 134, col: 23, offset: 2618},
	name: "WS",
},
&litMatcher{
	pos: position{line: 134, col: 26, offset: 2621},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 134, col: 36, offset: 2631},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 134, col: 40, offset: 2635},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 134, col: 45, offset: 2640},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 134, col: 53, offset: 2648},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 138, col: 1, offset: 2689},
	expr: &actionExpr{
	pos: position{line: 138, col: 12, offset: 2700},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 138, col: 12, offset: 2700},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 138, col: 12, offset: 2700},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 138, col: 20, offset: 2708},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 138, col: 30, offset: 2718},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 138, col: 38, offset: 2726},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 138, col: 41, offset: 2729},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 138, col: 49, offset: 2737},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 138, col: 52, offset: 2740},
	expr: &seqExpr{
	pos: position{line: 138, col: 53, offset: 2741},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 138, col: 53, offset: 2741},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 138, col: 56, offset: 2744},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 138, col: 59, offset: 2747},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 138, col: 62, offset: 2750},
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
	pos: position{line: 142, col: 1, offset: 2790},
	expr: &actionExpr{
	pos: position{line: 142, col: 11, offset: 2800},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 142, col: 11, offset: 2800},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 142, col: 11, offset: 2800},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 142, col: 14, offset: 2803},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 142, col: 21, offset: 2810},
	name: "WS",
},
&litMatcher{
	pos: position{line: 142, col: 24, offset: 2813},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 142, col: 28, offset: 2817},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 142, col: 31, offset: 2820},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 142, col: 34, offset: 2823},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 142, col: 34, offset: 2823},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 142, col: 45, offset: 2834},
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
	pos: position{line: 146, col: 1, offset: 2871},
	expr: &actionExpr{
	pos: position{line: 146, col: 16, offset: 2886},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 146, col: 16, offset: 2886},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 146, col: 16, offset: 2886},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 146, col: 24, offset: 2894},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 150, col: 1, offset: 2928},
	expr: &actionExpr{
	pos: position{line: 150, col: 12, offset: 2939},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 150, col: 12, offset: 2939},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 150, col: 12, offset: 2939},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 150, col: 20, offset: 2947},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 150, col: 30, offset: 2957},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 150, col: 38, offset: 2965},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 150, col: 41, offset: 2968},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 150, col: 41, offset: 2968},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 150, col: 52, offset: 2979},
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
	pos: position{line: 154, col: 1, offset: 3015},
	expr: &actionExpr{
	pos: position{line: 154, col: 12, offset: 3026},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 154, col: 12, offset: 3026},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 154, col: 12, offset: 3026},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 154, col: 20, offset: 3034},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 154, col: 30, offset: 3044},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 154, col: 38, offset: 3052},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 154, col: 41, offset: 3055},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 154, col: 41, offset: 3055},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 154, col: 52, offset: 3066},
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
	pos: position{line: 158, col: 1, offset: 3101},
	expr: &actionExpr{
	pos: position{line: 158, col: 14, offset: 3114},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 158, col: 14, offset: 3114},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 158, col: 14, offset: 3114},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 158, col: 22, offset: 3122},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 158, col: 34, offset: 3134},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 158, col: 42, offset: 3142},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 158, col: 45, offset: 3145},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 158, col: 45, offset: 3145},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 158, col: 56, offset: 3156},
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
	pos: position{line: 162, col: 1, offset: 3192},
	expr: &actionExpr{
	pos: position{line: 162, col: 14, offset: 3205},
	run: (*parser).callonFLAG_RULE1,
	expr: &seqExpr{
	pos: position{line: 162, col: 14, offset: 3205},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 162, col: 14, offset: 3205},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 162, col: 22, offset: 3213},
	val: "ignore-errors",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 166, col: 1, offset: 3260},
	expr: &actionExpr{
	pos: position{line: 166, col: 13, offset: 3272},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 166, col: 13, offset: 3272},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 166, col: 13, offset: 3272},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 166, col: 17, offset: 3276},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 166, col: 20, offset: 3279},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 170, col: 1, offset: 3314},
	expr: &actionExpr{
	pos: position{line: 170, col: 10, offset: 3323},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 170, col: 10, offset: 3323},
	expr: &charClassMatcher{
	pos: position{line: 170, col: 10, offset: 3323},
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
	pos: position{line: 174, col: 1, offset: 3371},
	expr: &actionExpr{
	pos: position{line: 174, col: 19, offset: 3389},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 174, col: 19, offset: 3389},
	expr: &charClassMatcher{
	pos: position{line: 174, col: 19, offset: 3389},
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
	name: "String",
	pos: position{line: 178, col: 1, offset: 3438},
	expr: &actionExpr{
	pos: position{line: 178, col: 11, offset: 3448},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 178, col: 11, offset: 3448},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 178, col: 11, offset: 3448},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 178, col: 15, offset: 3452},
	expr: &seqExpr{
	pos: position{line: 178, col: 17, offset: 3454},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 178, col: 17, offset: 3454},
	expr: &litMatcher{
	pos: position{line: 178, col: 18, offset: 3455},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 178, col: 22, offset: 3459,
},
	},
},
},
&litMatcher{
	pos: position{line: 178, col: 27, offset: 3464},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 182, col: 1, offset: 3513},
	expr: &actionExpr{
	pos: position{line: 182, col: 10, offset: 3522},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 182, col: 10, offset: 3522},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 182, col: 10, offset: 3522},
	expr: &choiceExpr{
	pos: position{line: 182, col: 11, offset: 3523},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 182, col: 11, offset: 3523},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 182, col: 17, offset: 3529},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 182, col: 23, offset: 3535},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 182, col: 31, offset: 3543},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 182, col: 35, offset: 3547},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 186, col: 1, offset: 3607},
	expr: &actionExpr{
	pos: position{line: 186, col: 12, offset: 3618},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 186, col: 12, offset: 3618},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 186, col: 12, offset: 3618},
	expr: &choiceExpr{
	pos: position{line: 186, col: 13, offset: 3619},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 186, col: 13, offset: 3619},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 186, col: 19, offset: 3625},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 186, col: 25, offset: 3631},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 190, col: 1, offset: 3693},
	expr: &choiceExpr{
	pos: position{line: 190, col: 11, offset: 3705},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 190, col: 11, offset: 3705},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 190, col: 17, offset: 3711},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 190, col: 17, offset: 3711},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 190, col: 37, offset: 3731},
	expr: &ruleRefExpr{
	pos: position{line: 190, col: 37, offset: 3731},
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
	pos: position{line: 192, col: 1, offset: 3746},
	expr: &charClassMatcher{
	pos: position{line: 192, col: 16, offset: 3763},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 193, col: 1, offset: 3769},
	expr: &charClassMatcher{
	pos: position{line: 193, col: 23, offset: 3793},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 195, col: 1, offset: 3800},
	expr: &charClassMatcher{
	pos: position{line: 195, col: 10, offset: 3809},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 196, col: 1, offset: 3815},
	expr: &charClassMatcher{
	pos: position{line: 196, col: 18, offset: 3832},
	val: "[\\n\\r]",
	chars: []rune{'\n','\r',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 198, col: 1, offset: 3840},
	expr: &choiceExpr{
	pos: position{line: 198, col: 25, offset: 3864},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 198, col: 25, offset: 3864},
	name: "NL",
},
&litMatcher{
	pos: position{line: 198, col: 30, offset: 3869},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 198, col: 36, offset: 3875},
	name: "COMMENT",
},
	},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 199, col: 1, offset: 3884},
	expr: &zeroOrMoreExpr{
	pos: position{line: 199, col: 20, offset: 3903},
	expr: &choiceExpr{
	pos: position{line: 199, col: 21, offset: 3904},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 199, col: 21, offset: 3904},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 199, col: 29, offset: 3912},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 200, col: 1, offset: 3922},
	expr: &oneOrMoreExpr{
	pos: position{line: 200, col: 35, offset: 3956},
	expr: &choiceExpr{
	pos: position{line: 200, col: 36, offset: 3957},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 200, col: 36, offset: 3957},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 200, col: 44, offset: 3965},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 200, col: 54, offset: 3975},
	name: "NL",
},
	},
},
},
},
{
	name: "COMMENT",
	pos: position{line: 202, col: 1, offset: 3981},
	expr: &seqExpr{
	pos: position{line: 202, col: 12, offset: 3992},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 202, col: 12, offset: 3992},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 202, col: 17, offset: 3997},
	expr: &seqExpr{
	pos: position{line: 202, col: 19, offset: 3999},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 202, col: 19, offset: 3999},
	expr: &litMatcher{
	pos: position{line: 202, col: 20, offset: 4000},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 202, col: 25, offset: 4005,
},
	},
},
},
&choiceExpr{
	pos: position{line: 202, col: 31, offset: 4011},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 202, col: 31, offset: 4011},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 202, col: 38, offset: 4018},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 204, col: 1, offset: 4024},
	expr: &notExpr{
	pos: position{line: 204, col: 8, offset: 4031},
	expr: &anyMatcher{
	line: 204, col: 9, offset: 4032,
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

func (c *current) onUSE_RULE1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonUSE_RULE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE_RULE1()
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
	return string(c.text), nil
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
	return string(c.text), nil
}

func (p *parser) callonIDENT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIDENT1()
}

func (c *current) onIDENT_WITH_DOT1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIDENT_WITH_DOT1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIDENT_WITH_DOT1()
}

func (c *current) onString1() (interface{}, error) {
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onFloat1() (interface{}, error) {
	return strconv.ParseFloat(string(c.text), 64)
}

func (p *parser) callonFloat1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFloat1()
}

func (c *current) onInteger1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 64)
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

