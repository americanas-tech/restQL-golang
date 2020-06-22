
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
	expr: &choiceExpr{
	pos: position{line: 17, col: 48, offset: 165},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 48, offset: 165},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 17, col: 53, offset: 170},
	name: "COMMENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 63, offset: 180},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 17, col: 66, offset: 183},
	label: "firstBlock",
	expr: &ruleRefExpr{
	pos: position{line: 17, col: 77, offset: 194},
	name: "BLOCK",
},
},
&labeledExpr{
	pos: position{line: 17, col: 83, offset: 200},
	label: "otherBlocks",
	expr: &zeroOrMoreExpr{
	pos: position{line: 17, col: 95, offset: 212},
	expr: &seqExpr{
	pos: position{line: 17, col: 96, offset: 213},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 96, offset: 213},
	name: "BS",
},
&ruleRefExpr{
	pos: position{line: 17, col: 99, offset: 216},
	name: "BLOCK",
},
	},
},
},
},
&zeroOrMoreExpr{
	pos: position{line: 17, col: 107, offset: 224},
	expr: &choiceExpr{
	pos: position{line: 17, col: 108, offset: 225},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 17, col: 108, offset: 225},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 17, col: 113, offset: 230},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 17, col: 121, offset: 238},
	name: "COMMENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 17, col: 131, offset: 248},
	name: "EOF",
},
	},
},
},
},
{
	name: "USE",
	pos: position{line: 21, col: 1, offset: 303},
	expr: &actionExpr{
	pos: position{line: 21, col: 8, offset: 310},
	run: (*parser).callonUSE1,
	expr: &seqExpr{
	pos: position{line: 21, col: 8, offset: 310},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 21, col: 8, offset: 310},
	val: "use",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 21, col: 14, offset: 316},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 21, col: 22, offset: 324},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 25, offset: 327},
	name: "USE_ACTION",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 37, offset: 339},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 21, col: 40, offset: 342},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 43, offset: 345},
	name: "USE_VALUE",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 54, offset: 356},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 21, col: 57, offset: 359},
	expr: &ruleRefExpr{
	pos: position{line: 21, col: 57, offset: 359},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 21, col: 61, offset: 363},
	name: "WS",
},
	},
},
},
},
{
	name: "USE_ACTION",
	pos: position{line: 25, col: 1, offset: 392},
	expr: &actionExpr{
	pos: position{line: 25, col: 15, offset: 406},
	run: (*parser).callonUSE_ACTION1,
	expr: &choiceExpr{
	pos: position{line: 25, col: 16, offset: 407},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 25, col: 16, offset: 407},
	val: "timeout",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 28, offset: 419},
	val: "max-age",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 25, col: 40, offset: 431},
	val: "s-max-age",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "USE_VALUE",
	pos: position{line: 29, col: 1, offset: 475},
	expr: &actionExpr{
	pos: position{line: 29, col: 14, offset: 488},
	run: (*parser).callonUSE_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 29, col: 14, offset: 488},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 29, col: 17, offset: 491},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 29, col: 17, offset: 491},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 29, col: 26, offset: 500},
	name: "Integer",
},
	},
},
},
},
},
{
	name: "BLOCK",
	pos: position{line: 33, col: 1, offset: 537},
	expr: &actionExpr{
	pos: position{line: 33, col: 10, offset: 546},
	run: (*parser).callonBLOCK1,
	expr: &seqExpr{
	pos: position{line: 33, col: 10, offset: 546},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 33, col: 10, offset: 546},
	label: "action",
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 18, offset: 554},
	name: "ACTION_RULE",
},
},
&labeledExpr{
	pos: position{line: 33, col: 31, offset: 567},
	label: "m",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 34, offset: 570},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 34, offset: 570},
	name: "MODIFIER_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 50, offset: 586},
	label: "w",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 53, offset: 589},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 53, offset: 589},
	name: "WITH_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 65, offset: 601},
	label: "f",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 67, offset: 603},
	expr: &choiceExpr{
	pos: position{line: 33, col: 68, offset: 604},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 33, col: 68, offset: 604},
	name: "HIDDEN_RULE",
},
&ruleRefExpr{
	pos: position{line: 33, col: 82, offset: 618},
	name: "ONLY_RULE",
},
	},
},
},
},
&labeledExpr{
	pos: position{line: 33, col: 94, offset: 630},
	label: "fl",
	expr: &zeroOrOneExpr{
	pos: position{line: 33, col: 98, offset: 634},
	expr: &ruleRefExpr{
	pos: position{line: 33, col: 98, offset: 634},
	name: "FLAGS_RULE",
},
},
},
&ruleRefExpr{
	pos: position{line: 33, col: 111, offset: 647},
	name: "WS",
},
	},
},
},
},
{
	name: "ACTION_RULE",
	pos: position{line: 37, col: 1, offset: 693},
	expr: &actionExpr{
	pos: position{line: 37, col: 16, offset: 708},
	run: (*parser).callonACTION_RULE1,
	expr: &seqExpr{
	pos: position{line: 37, col: 16, offset: 708},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 37, col: 16, offset: 708},
	label: "m",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 19, offset: 711},
	name: "METHOD",
},
},
&ruleRefExpr{
	pos: position{line: 37, col: 27, offset: 719},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 37, col: 35, offset: 727},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 38, offset: 730},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 37, col: 45, offset: 737},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 48, offset: 740},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 48, offset: 740},
	name: "ALIAS",
},
},
},
&labeledExpr{
	pos: position{line: 37, col: 56, offset: 748},
	label: "i",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 59, offset: 751},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 59, offset: 751},
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
	pos: position{line: 41, col: 1, offset: 795},
	expr: &actionExpr{
	pos: position{line: 41, col: 11, offset: 805},
	run: (*parser).callonMETHOD1,
	expr: &choiceExpr{
	pos: position{line: 41, col: 12, offset: 806},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 41, col: 12, offset: 806},
	val: "from",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 21, offset: 815},
	val: "to",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 28, offset: 822},
	val: "into",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 36, offset: 830},
	val: "update",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 41, col: 47, offset: 841},
	val: "delete",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "ALIAS",
	pos: position{line: 45, col: 1, offset: 882},
	expr: &actionExpr{
	pos: position{line: 45, col: 10, offset: 891},
	run: (*parser).callonALIAS1,
	expr: &seqExpr{
	pos: position{line: 45, col: 10, offset: 891},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 45, col: 10, offset: 891},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 45, col: 18, offset: 899},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 45, col: 23, offset: 904},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 45, col: 31, offset: 912},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 45, col: 34, offset: 915},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IN",
	pos: position{line: 49, col: 1, offset: 942},
	expr: &actionExpr{
	pos: position{line: 49, col: 7, offset: 948},
	run: (*parser).callonIN1,
	expr: &seqExpr{
	pos: position{line: 49, col: 7, offset: 948},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 7, offset: 948},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 49, col: 15, offset: 956},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 49, col: 20, offset: 961},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 49, col: 28, offset: 969},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 49, col: 31, offset: 972},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "MODIFIER_RULE",
	pos: position{line: 53, col: 1, offset: 1010},
	expr: &actionExpr{
	pos: position{line: 53, col: 18, offset: 1027},
	run: (*parser).callonMODIFIER_RULE1,
	expr: &labeledExpr{
	pos: position{line: 53, col: 18, offset: 1027},
	label: "m",
	expr: &oneOrMoreExpr{
	pos: position{line: 53, col: 20, offset: 1029},
	expr: &choiceExpr{
	pos: position{line: 53, col: 21, offset: 1030},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 53, col: 21, offset: 1030},
	name: "HEADERS",
},
&ruleRefExpr{
	pos: position{line: 53, col: 31, offset: 1040},
	name: "TIMEOUT",
},
&ruleRefExpr{
	pos: position{line: 53, col: 41, offset: 1050},
	name: "MAX_AGE",
},
&ruleRefExpr{
	pos: position{line: 53, col: 51, offset: 1060},
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
	pos: position{line: 57, col: 1, offset: 1092},
	expr: &actionExpr{
	pos: position{line: 57, col: 14, offset: 1105},
	run: (*parser).callonWITH_RULE1,
	expr: &seqExpr{
	pos: position{line: 57, col: 14, offset: 1105},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 57, col: 14, offset: 1105},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 57, col: 22, offset: 1113},
	val: "with",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 57, col: 29, offset: 1120},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 57, col: 37, offset: 1128},
	label: "pb",
	expr: &zeroOrOneExpr{
	pos: position{line: 57, col: 40, offset: 1131},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 40, offset: 1131},
	name: "PARAMETER_BODY",
},
},
},
&labeledExpr{
	pos: position{line: 57, col: 56, offset: 1147},
	label: "kvs",
	expr: &zeroOrOneExpr{
	pos: position{line: 57, col: 60, offset: 1151},
	expr: &ruleRefExpr{
	pos: position{line: 57, col: 60, offset: 1151},
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
	pos: position{line: 61, col: 1, offset: 1197},
	expr: &actionExpr{
	pos: position{line: 61, col: 19, offset: 1215},
	run: (*parser).callonPARAMETER_BODY1,
	expr: &seqExpr{
	pos: position{line: 61, col: 19, offset: 1215},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 61, col: 19, offset: 1215},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 61, col: 23, offset: 1219},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 26, offset: 1222},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 61, col: 33, offset: 1229},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 61, col: 37, offset: 1233},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 37, offset: 1233},
	name: "APPLY_FN",
},
},
},
&ruleRefExpr{
	pos: position{line: 61, col: 48, offset: 1244},
	name: "WS",
},
&zeroOrOneExpr{
	pos: position{line: 61, col: 51, offset: 1247},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 51, offset: 1247},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 61, col: 55, offset: 1251},
	name: "WS",
},
	},
},
},
},
{
	name: "KEY_VALUE_LIST",
	pos: position{line: 65, col: 1, offset: 1291},
	expr: &actionExpr{
	pos: position{line: 65, col: 19, offset: 1309},
	run: (*parser).callonKEY_VALUE_LIST1,
	expr: &seqExpr{
	pos: position{line: 65, col: 19, offset: 1309},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 65, col: 19, offset: 1309},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 25, offset: 1315},
	name: "KEY_VALUE",
},
},
&labeledExpr{
	pos: position{line: 65, col: 35, offset: 1325},
	label: "others",
	expr: &zeroOrMoreExpr{
	pos: position{line: 65, col: 42, offset: 1332},
	expr: &seqExpr{
	pos: position{line: 65, col: 43, offset: 1333},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 43, offset: 1333},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 46, offset: 1336},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 49, offset: 1339},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 52, offset: 1342},
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
	pos: position{line: 69, col: 1, offset: 1398},
	expr: &actionExpr{
	pos: position{line: 69, col: 14, offset: 1411},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 69, col: 14, offset: 1411},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 69, col: 14, offset: 1411},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 17, offset: 1414},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 69, col: 33, offset: 1430},
	name: "WS",
},
&litMatcher{
	pos: position{line: 69, col: 36, offset: 1433},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 69, col: 40, offset: 1437},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 69, col: 43, offset: 1440},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 46, offset: 1443},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 69, col: 53, offset: 1450},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 69, col: 57, offset: 1454},
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 57, offset: 1454},
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
	pos: position{line: 73, col: 1, offset: 1500},
	expr: &actionExpr{
	pos: position{line: 73, col: 13, offset: 1512},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 73, col: 13, offset: 1512},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 73, col: 13, offset: 1512},
	name: "WS",
},
&litMatcher{
	pos: position{line: 73, col: 16, offset: 1515},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 73, col: 21, offset: 1520},
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 21, offset: 1520},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 73, col: 25, offset: 1524},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 29, offset: 1528},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 77, col: 1, offset: 1559},
	expr: &actionExpr{
	pos: position{line: 77, col: 13, offset: 1571},
	run: (*parser).callonFUNCTION1,
	expr: &choiceExpr{
	pos: position{line: 77, col: 14, offset: 1572},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 77, col: 14, offset: 1572},
	val: "flatten",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 26, offset: 1584},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 37, offset: 1595},
	val: "json",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VALUE",
	pos: position{line: 81, col: 1, offset: 1634},
	expr: &actionExpr{
	pos: position{line: 81, col: 10, offset: 1643},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 81, col: 10, offset: 1643},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 81, col: 13, offset: 1646},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 81, col: 13, offset: 1646},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 81, col: 20, offset: 1653},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 81, col: 29, offset: 1662},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 81, col: 40, offset: 1673},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 85, col: 1, offset: 1709},
	expr: &actionExpr{
	pos: position{line: 85, col: 9, offset: 1717},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 85, col: 9, offset: 1717},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 85, col: 12, offset: 1720},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 85, col: 12, offset: 1720},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 85, col: 25, offset: 1733},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 89, col: 1, offset: 1769},
	expr: &actionExpr{
	pos: position{line: 89, col: 15, offset: 1783},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 89, col: 15, offset: 1783},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 15, offset: 1783},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 89, col: 19, offset: 1787},
	name: "WS",
},
&litMatcher{
	pos: position{line: 89, col: 22, offset: 1790},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 93, col: 1, offset: 1822},
	expr: &actionExpr{
	pos: position{line: 93, col: 19, offset: 1840},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 93, col: 19, offset: 1840},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 93, col: 19, offset: 1840},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 93, col: 23, offset: 1844},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 93, col: 26, offset: 1847},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 28, offset: 1849},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 93, col: 34, offset: 1855},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 93, col: 37, offset: 1858},
	expr: &seqExpr{
	pos: position{line: 93, col: 38, offset: 1859},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 93, col: 38, offset: 1859},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 93, col: 41, offset: 1862},
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 41, offset: 1862},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 45, offset: 1866},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 93, col: 48, offset: 1869},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 56, offset: 1877},
	name: "WS",
},
&litMatcher{
	pos: position{line: 93, col: 59, offset: 1880},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 97, col: 1, offset: 1912},
	expr: &actionExpr{
	pos: position{line: 97, col: 11, offset: 1922},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 97, col: 11, offset: 1922},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 97, col: 14, offset: 1925},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 97, col: 14, offset: 1925},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 97, col: 26, offset: 1937},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 101, col: 1, offset: 1972},
	expr: &actionExpr{
	pos: position{line: 101, col: 14, offset: 1985},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 101, col: 14, offset: 1985},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 101, col: 14, offset: 1985},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 101, col: 18, offset: 1989},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 101, col: 21, offset: 1992},
	expr: &ruleRefExpr{
	pos: position{line: 101, col: 21, offset: 1992},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 101, col: 25, offset: 1996},
	name: "WS",
},
&litMatcher{
	pos: position{line: 101, col: 28, offset: 1999},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 105, col: 1, offset: 2033},
	expr: &actionExpr{
	pos: position{line: 105, col: 18, offset: 2050},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 105, col: 18, offset: 2050},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 105, col: 18, offset: 2050},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 22, offset: 2054},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 25, offset: 2057},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 25, offset: 2057},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 29, offset: 2061},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 105, col: 32, offset: 2064},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 36, offset: 2068},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 105, col: 47, offset: 2079},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 105, col: 51, offset: 2083},
	expr: &seqExpr{
	pos: position{line: 105, col: 52, offset: 2084},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 105, col: 52, offset: 2084},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 55, offset: 2087},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 59, offset: 2091},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 62, offset: 2094},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 62, offset: 2094},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 66, offset: 2098},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 105, col: 69, offset: 2101},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 81, offset: 2113},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 84, offset: 2116},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 84, offset: 2116},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 88, offset: 2120},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 91, offset: 2123},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 109, col: 1, offset: 2168},
	expr: &actionExpr{
	pos: position{line: 109, col: 14, offset: 2181},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 109, col: 14, offset: 2181},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 109, col: 14, offset: 2181},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 109, col: 17, offset: 2184},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 109, col: 17, offset: 2184},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 109, col: 26, offset: 2193},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 33, offset: 2200},
	name: "WS",
},
&litMatcher{
	pos: position{line: 109, col: 36, offset: 2203},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 109, col: 40, offset: 2207},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 109, col: 43, offset: 2210},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 46, offset: 2213},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 113, col: 1, offset: 2254},
	expr: &actionExpr{
	pos: position{line: 113, col: 14, offset: 2267},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 113, col: 14, offset: 2267},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 113, col: 17, offset: 2270},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2270},
	name: "Null",
},
&ruleRefExpr{
	pos: position{line: 113, col: 24, offset: 2277},
	name: "Boolean",
},
&ruleRefExpr{
	pos: position{line: 113, col: 34, offset: 2287},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 113, col: 43, offset: 2296},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 113, col: 51, offset: 2304},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 113, col: 61, offset: 2314},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "CHAIN",
	pos: position{line: 117, col: 1, offset: 2350},
	expr: &actionExpr{
	pos: position{line: 117, col: 10, offset: 2359},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 117, col: 10, offset: 2359},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 117, col: 10, offset: 2359},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 117, col: 13, offset: 2362},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 117, col: 27, offset: 2376},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 117, col: 30, offset: 2379},
	expr: &seqExpr{
	pos: position{line: 117, col: 31, offset: 2380},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 117, col: 31, offset: 2380},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 117, col: 35, offset: 2384},
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
	pos: position{line: 121, col: 1, offset: 2428},
	expr: &actionExpr{
	pos: position{line: 121, col: 17, offset: 2444},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 121, col: 17, offset: 2444},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 121, col: 21, offset: 2448},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 121, col: 21, offset: 2448},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 121, col: 32, offset: 2459},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 125, col: 1, offset: 2494},
	expr: &actionExpr{
	pos: position{line: 125, col: 14, offset: 2507},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 125, col: 14, offset: 2507},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 14, offset: 2507},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 125, col: 22, offset: 2515},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 125, col: 29, offset: 2522},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 125, col: 37, offset: 2530},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 125, col: 40, offset: 2533},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 125, col: 48, offset: 2541},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 125, col: 51, offset: 2544},
	expr: &seqExpr{
	pos: position{line: 125, col: 52, offset: 2545},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 52, offset: 2545},
	name: "WS",
},
&notExpr{
	pos: position{line: 125, col: 55, offset: 2548},
	expr: &choiceExpr{
	pos: position{line: 125, col: 57, offset: 2550},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 57, offset: 2550},
	name: "FLAGS_RULE",
},
&seqExpr{
	pos: position{line: 125, col: 70, offset: 2563},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 70, offset: 2563},
	name: "BS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 73, offset: 2566},
	name: "BLOCK",
},
	},
},
	},
},
},
&choiceExpr{
	pos: position{line: 125, col: 81, offset: 2574},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 125, col: 81, offset: 2574},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 81, offset: 2574},
	name: "LS",
},
&zeroOrMoreExpr{
	pos: position{line: 125, col: 84, offset: 2577},
	expr: &seqExpr{
	pos: position{line: 125, col: 85, offset: 2578},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 125, col: 85, offset: 2578},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 88, offset: 2581},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 125, col: 91, offset: 2584},
	name: "WS",
},
	},
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 125, col: 98, offset: 2591},
	name: "LS",
},
	},
},
&ruleRefExpr{
	pos: position{line: 125, col: 102, offset: 2595},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 125, col: 105, offset: 2598},
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
	pos: position{line: 129, col: 1, offset: 2635},
	expr: &actionExpr{
	pos: position{line: 129, col: 11, offset: 2645},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 129, col: 11, offset: 2645},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 129, col: 11, offset: 2645},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 14, offset: 2648},
	name: "FILTER_VALUE",
},
},
&labeledExpr{
	pos: position{line: 129, col: 28, offset: 2662},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 129, col: 32, offset: 2666},
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 32, offset: 2666},
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
	pos: position{line: 133, col: 1, offset: 2709},
	expr: &actionExpr{
	pos: position{line: 133, col: 17, offset: 2725},
	run: (*parser).callonFILTER_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 133, col: 17, offset: 2725},
	label: "fv",
	expr: &choiceExpr{
	pos: position{line: 133, col: 21, offset: 2729},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 133, col: 21, offset: 2729},
	name: "IDENT_WITH_DOT",
},
&litMatcher{
	pos: position{line: 133, col: 38, offset: 2746},
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
	pos: position{line: 137, col: 1, offset: 2783},
	expr: &actionExpr{
	pos: position{line: 137, col: 15, offset: 2797},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 137, col: 15, offset: 2797},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 137, col: 15, offset: 2797},
	name: "WS",
},
&litMatcher{
	pos: position{line: 137, col: 18, offset: 2800},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 137, col: 23, offset: 2805},
	name: "WS",
},
&litMatcher{
	pos: position{line: 137, col: 26, offset: 2808},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 137, col: 36, offset: 2818},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 137, col: 40, offset: 2822},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 137, col: 45, offset: 2827},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 137, col: 53, offset: 2835},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 141, col: 1, offset: 2876},
	expr: &actionExpr{
	pos: position{line: 141, col: 12, offset: 2887},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 141, col: 12, offset: 2887},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 141, col: 12, offset: 2887},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 141, col: 20, offset: 2895},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 141, col: 30, offset: 2905},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 141, col: 38, offset: 2913},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 141, col: 41, offset: 2916},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 141, col: 49, offset: 2924},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 141, col: 52, offset: 2927},
	expr: &seqExpr{
	pos: position{line: 141, col: 53, offset: 2928},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 141, col: 53, offset: 2928},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 141, col: 56, offset: 2931},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 141, col: 59, offset: 2934},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 141, col: 62, offset: 2937},
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
	pos: position{line: 145, col: 1, offset: 2977},
	expr: &actionExpr{
	pos: position{line: 145, col: 11, offset: 2987},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 145, col: 11, offset: 2987},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 145, col: 11, offset: 2987},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 145, col: 14, offset: 2990},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 145, col: 21, offset: 2997},
	name: "WS",
},
&litMatcher{
	pos: position{line: 145, col: 24, offset: 3000},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 145, col: 28, offset: 3004},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 145, col: 31, offset: 3007},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 145, col: 34, offset: 3010},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 145, col: 34, offset: 3010},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 145, col: 45, offset: 3021},
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
	pos: position{line: 149, col: 1, offset: 3058},
	expr: &actionExpr{
	pos: position{line: 149, col: 16, offset: 3073},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 149, col: 16, offset: 3073},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 149, col: 16, offset: 3073},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 149, col: 24, offset: 3081},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 153, col: 1, offset: 3115},
	expr: &actionExpr{
	pos: position{line: 153, col: 12, offset: 3126},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 153, col: 12, offset: 3126},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 12, offset: 3126},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 153, col: 20, offset: 3134},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 153, col: 30, offset: 3144},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 153, col: 38, offset: 3152},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 153, col: 41, offset: 3155},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 153, col: 41, offset: 3155},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 153, col: 52, offset: 3166},
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
	pos: position{line: 157, col: 1, offset: 3202},
	expr: &actionExpr{
	pos: position{line: 157, col: 12, offset: 3213},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 157, col: 12, offset: 3213},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 12, offset: 3213},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 157, col: 20, offset: 3221},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 157, col: 30, offset: 3231},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 157, col: 38, offset: 3239},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 157, col: 41, offset: 3242},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 157, col: 41, offset: 3242},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 157, col: 52, offset: 3253},
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
	pos: position{line: 161, col: 1, offset: 3288},
	expr: &actionExpr{
	pos: position{line: 161, col: 14, offset: 3301},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 161, col: 14, offset: 3301},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 14, offset: 3301},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 161, col: 22, offset: 3309},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 161, col: 34, offset: 3321},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 161, col: 42, offset: 3329},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 161, col: 45, offset: 3332},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 161, col: 45, offset: 3332},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 161, col: 56, offset: 3343},
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
	pos: position{line: 165, col: 1, offset: 3379},
	expr: &actionExpr{
	pos: position{line: 165, col: 15, offset: 3393},
	run: (*parser).callonFLAGS_RULE1,
	expr: &seqExpr{
	pos: position{line: 165, col: 15, offset: 3393},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 165, col: 15, offset: 3393},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 165, col: 23, offset: 3401},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 165, col: 25, offset: 3403},
	name: "IGNORE_FLAG",
},
},
&labeledExpr{
	pos: position{line: 165, col: 37, offset: 3415},
	label: "is",
	expr: &zeroOrMoreExpr{
	pos: position{line: 165, col: 40, offset: 3418},
	expr: &seqExpr{
	pos: position{line: 165, col: 41, offset: 3419},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 165, col: 41, offset: 3419},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 165, col: 44, offset: 3422},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 165, col: 47, offset: 3425},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 165, col: 50, offset: 3428},
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
	pos: position{line: 169, col: 1, offset: 3471},
	expr: &actionExpr{
	pos: position{line: 169, col: 16, offset: 3486},
	run: (*parser).callonIGNORE_FLAG1,
	expr: &litMatcher{
	pos: position{line: 169, col: 16, offset: 3486},
	val: "ignore-errors",
	ignoreCase: false,
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 173, col: 1, offset: 3533},
	expr: &actionExpr{
	pos: position{line: 173, col: 13, offset: 3545},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 173, col: 13, offset: 3545},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 173, col: 13, offset: 3545},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 173, col: 17, offset: 3549},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 173, col: 20, offset: 3552},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 177, col: 1, offset: 3587},
	expr: &actionExpr{
	pos: position{line: 177, col: 10, offset: 3596},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 177, col: 10, offset: 3596},
	expr: &charClassMatcher{
	pos: position{line: 177, col: 10, offset: 3596},
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
	pos: position{line: 181, col: 1, offset: 3642},
	expr: &actionExpr{
	pos: position{line: 181, col: 19, offset: 3660},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 181, col: 19, offset: 3660},
	expr: &charClassMatcher{
	pos: position{line: 181, col: 19, offset: 3660},
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
	pos: position{line: 185, col: 1, offset: 3707},
	expr: &actionExpr{
	pos: position{line: 185, col: 9, offset: 3715},
	run: (*parser).callonNull1,
	expr: &litMatcher{
	pos: position{line: 185, col: 9, offset: 3715},
	val: "null",
	ignoreCase: false,
},
},
},
{
	name: "Boolean",
	pos: position{line: 189, col: 1, offset: 3745},
	expr: &actionExpr{
	pos: position{line: 189, col: 12, offset: 3756},
	run: (*parser).callonBoolean1,
	expr: &choiceExpr{
	pos: position{line: 189, col: 13, offset: 3757},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 189, col: 13, offset: 3757},
	val: "true",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 189, col: 22, offset: 3766},
	val: "false",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "String",
	pos: position{line: 193, col: 1, offset: 3807},
	expr: &actionExpr{
	pos: position{line: 193, col: 11, offset: 3817},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 193, col: 11, offset: 3817},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 193, col: 11, offset: 3817},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 193, col: 15, offset: 3821},
	expr: &seqExpr{
	pos: position{line: 193, col: 17, offset: 3823},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 193, col: 17, offset: 3823},
	expr: &litMatcher{
	pos: position{line: 193, col: 18, offset: 3824},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 193, col: 22, offset: 3828,
},
	},
},
},
&litMatcher{
	pos: position{line: 193, col: 27, offset: 3833},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 197, col: 1, offset: 3868},
	expr: &actionExpr{
	pos: position{line: 197, col: 10, offset: 3877},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 197, col: 10, offset: 3877},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 197, col: 10, offset: 3877},
	expr: &choiceExpr{
	pos: position{line: 197, col: 11, offset: 3878},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 197, col: 11, offset: 3878},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 197, col: 17, offset: 3884},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 197, col: 23, offset: 3890},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 197, col: 31, offset: 3898},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 197, col: 35, offset: 3902},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 201, col: 1, offset: 3940},
	expr: &actionExpr{
	pos: position{line: 201, col: 12, offset: 3951},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 201, col: 12, offset: 3951},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 201, col: 12, offset: 3951},
	expr: &choiceExpr{
	pos: position{line: 201, col: 13, offset: 3952},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 201, col: 13, offset: 3952},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 201, col: 19, offset: 3958},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 201, col: 25, offset: 3964},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 205, col: 1, offset: 4004},
	expr: &choiceExpr{
	pos: position{line: 205, col: 11, offset: 4016},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 205, col: 11, offset: 4016},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 205, col: 17, offset: 4022},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 205, col: 17, offset: 4022},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 205, col: 37, offset: 4042},
	expr: &ruleRefExpr{
	pos: position{line: 205, col: 37, offset: 4042},
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
	pos: position{line: 207, col: 1, offset: 4057},
	expr: &charClassMatcher{
	pos: position{line: 207, col: 16, offset: 4074},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 208, col: 1, offset: 4080},
	expr: &charClassMatcher{
	pos: position{line: 208, col: 23, offset: 4104},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 210, col: 1, offset: 4111},
	expr: &charClassMatcher{
	pos: position{line: 210, col: 10, offset: 4120},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 211, col: 1, offset: 4126},
	expr: &oneOrMoreExpr{
	pos: position{line: 211, col: 35, offset: 4160},
	expr: &choiceExpr{
	pos: position{line: 211, col: 36, offset: 4161},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 211, col: 36, offset: 4161},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 211, col: 44, offset: 4169},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 211, col: 54, offset: 4179},
	name: "NL",
},
	},
},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 212, col: 1, offset: 4184},
	expr: &zeroOrMoreExpr{
	pos: position{line: 212, col: 20, offset: 4203},
	expr: &choiceExpr{
	pos: position{line: 212, col: 21, offset: 4204},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 212, col: 21, offset: 4204},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 212, col: 29, offset: 4212},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 213, col: 1, offset: 4222},
	expr: &choiceExpr{
	pos: position{line: 213, col: 25, offset: 4246},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 213, col: 25, offset: 4246},
	name: "NL",
},
&litMatcher{
	pos: position{line: 213, col: 30, offset: 4251},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 213, col: 36, offset: 4257},
	name: "COMMENT",
},
	},
},
},
{
	name: "BS",
	displayName: "\"block-separator\"",
	pos: position{line: 214, col: 1, offset: 4266},
	expr: &oneOrMoreExpr{
	pos: position{line: 214, col: 25, offset: 4290},
	expr: &seqExpr{
	pos: position{line: 214, col: 26, offset: 4291},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 214, col: 26, offset: 4291},
	name: "WS",
},
&choiceExpr{
	pos: position{line: 214, col: 30, offset: 4295},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 214, col: 30, offset: 4295},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 214, col: 35, offset: 4300},
	name: "COMMENT",
},
	},
},
&ruleRefExpr{
	pos: position{line: 214, col: 44, offset: 4309},
	name: "WS",
},
	},
},
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 215, col: 1, offset: 4314},
	expr: &litMatcher{
	pos: position{line: 215, col: 18, offset: 4331},
	val: "\n",
	ignoreCase: false,
},
},
{
	name: "COMMENT",
	pos: position{line: 217, col: 1, offset: 4337},
	expr: &seqExpr{
	pos: position{line: 217, col: 12, offset: 4348},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 217, col: 12, offset: 4348},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 217, col: 17, offset: 4353},
	expr: &seqExpr{
	pos: position{line: 217, col: 19, offset: 4355},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 217, col: 19, offset: 4355},
	expr: &litMatcher{
	pos: position{line: 217, col: 20, offset: 4356},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 217, col: 25, offset: 4361,
},
	},
},
},
&choiceExpr{
	pos: position{line: 217, col: 31, offset: 4367},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 217, col: 31, offset: 4367},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 217, col: 38, offset: 4374},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 219, col: 1, offset: 4380},
	expr: &notExpr{
	pos: position{line: 219, col: 8, offset: 4387},
	expr: &anyMatcher{
	line: 219, col: 9, offset: 4388,
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

