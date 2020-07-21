
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
	expr: &zeroOrMoreExpr{
	pos: position{line: 61, col: 36, offset: 1232},
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
&choiceExpr{
	pos: position{line: 65, col: 47, offset: 1337},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 65, col: 47, offset: 1337},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 47, offset: 1337},
	name: "LS",
},
&zeroOrMoreExpr{
	pos: position{line: 65, col: 50, offset: 1340},
	expr: &seqExpr{
	pos: position{line: 65, col: 51, offset: 1341},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 51, offset: 1341},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 54, offset: 1344},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 65, col: 57, offset: 1347},
	name: "WS",
},
	},
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 65, col: 64, offset: 1354},
	name: "LS",
},
	},
},
&ruleRefExpr{
	pos: position{line: 65, col: 68, offset: 1358},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 65, col: 71, offset: 1361},
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
	pos: position{line: 69, col: 1, offset: 1417},
	expr: &actionExpr{
	pos: position{line: 69, col: 14, offset: 1430},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 69, col: 14, offset: 1430},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 69, col: 14, offset: 1430},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 17, offset: 1433},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 69, col: 33, offset: 1449},
	name: "WS",
},
&litMatcher{
	pos: position{line: 69, col: 36, offset: 1452},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 69, col: 40, offset: 1456},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 69, col: 43, offset: 1459},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 46, offset: 1462},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 69, col: 53, offset: 1469},
	label: "fn",
	expr: &zeroOrMoreExpr{
	pos: position{line: 69, col: 56, offset: 1472},
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 57, offset: 1473},
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
	pos: position{line: 73, col: 1, offset: 1519},
	expr: &actionExpr{
	pos: position{line: 73, col: 13, offset: 1531},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 73, col: 13, offset: 1531},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 73, col: 13, offset: 1531},
	name: "WS",
},
&litMatcher{
	pos: position{line: 73, col: 16, offset: 1534},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 73, col: 21, offset: 1539},
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 21, offset: 1539},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 73, col: 25, offset: 1543},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 29, offset: 1547},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 77, col: 1, offset: 1578},
	expr: &actionExpr{
	pos: position{line: 77, col: 13, offset: 1590},
	run: (*parser).callonFUNCTION1,
	expr: &choiceExpr{
	pos: position{line: 77, col: 14, offset: 1591},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 77, col: 14, offset: 1591},
	val: "no-multiplex",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 31, offset: 1608},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 42, offset: 1619},
	val: "json",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 77, col: 50, offset: 1627},
	val: "as-body",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "VALUE",
	pos: position{line: 81, col: 1, offset: 1669},
	expr: &actionExpr{
	pos: position{line: 81, col: 10, offset: 1678},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 81, col: 10, offset: 1678},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 81, col: 13, offset: 1681},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 81, col: 13, offset: 1681},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 81, col: 20, offset: 1688},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 81, col: 29, offset: 1697},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 81, col: 40, offset: 1708},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 85, col: 1, offset: 1744},
	expr: &actionExpr{
	pos: position{line: 85, col: 9, offset: 1752},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 85, col: 9, offset: 1752},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 85, col: 12, offset: 1755},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 85, col: 12, offset: 1755},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 85, col: 25, offset: 1768},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 89, col: 1, offset: 1804},
	expr: &actionExpr{
	pos: position{line: 89, col: 15, offset: 1818},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 89, col: 15, offset: 1818},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 89, col: 15, offset: 1818},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 89, col: 19, offset: 1822},
	name: "WS",
},
&litMatcher{
	pos: position{line: 89, col: 22, offset: 1825},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 93, col: 1, offset: 1857},
	expr: &actionExpr{
	pos: position{line: 93, col: 19, offset: 1875},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 93, col: 19, offset: 1875},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 93, col: 19, offset: 1875},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 93, col: 23, offset: 1879},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 93, col: 26, offset: 1882},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 28, offset: 1884},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 93, col: 34, offset: 1890},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 93, col: 37, offset: 1893},
	expr: &seqExpr{
	pos: position{line: 93, col: 38, offset: 1894},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 93, col: 38, offset: 1894},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 93, col: 41, offset: 1897},
	expr: &ruleRefExpr{
	pos: position{line: 93, col: 41, offset: 1897},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 45, offset: 1901},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 93, col: 48, offset: 1904},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 93, col: 56, offset: 1912},
	name: "WS",
},
&litMatcher{
	pos: position{line: 93, col: 59, offset: 1915},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 97, col: 1, offset: 1947},
	expr: &actionExpr{
	pos: position{line: 97, col: 11, offset: 1957},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 97, col: 11, offset: 1957},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 97, col: 14, offset: 1960},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 97, col: 14, offset: 1960},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 97, col: 26, offset: 1972},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 101, col: 1, offset: 2007},
	expr: &actionExpr{
	pos: position{line: 101, col: 14, offset: 2020},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 101, col: 14, offset: 2020},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 101, col: 14, offset: 2020},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 101, col: 18, offset: 2024},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 101, col: 21, offset: 2027},
	expr: &ruleRefExpr{
	pos: position{line: 101, col: 21, offset: 2027},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 101, col: 25, offset: 2031},
	name: "WS",
},
&litMatcher{
	pos: position{line: 101, col: 28, offset: 2034},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 105, col: 1, offset: 2068},
	expr: &actionExpr{
	pos: position{line: 105, col: 18, offset: 2085},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 105, col: 18, offset: 2085},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 105, col: 18, offset: 2085},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 22, offset: 2089},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 25, offset: 2092},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 25, offset: 2092},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 29, offset: 2096},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 105, col: 32, offset: 2099},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 36, offset: 2103},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 105, col: 47, offset: 2114},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 105, col: 51, offset: 2118},
	expr: &seqExpr{
	pos: position{line: 105, col: 52, offset: 2119},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 105, col: 52, offset: 2119},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 55, offset: 2122},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 105, col: 59, offset: 2126},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 62, offset: 2129},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 62, offset: 2129},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 66, offset: 2133},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 105, col: 69, offset: 2136},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 81, offset: 2148},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 105, col: 84, offset: 2151},
	expr: &ruleRefExpr{
	pos: position{line: 105, col: 84, offset: 2151},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 105, col: 88, offset: 2155},
	name: "WS",
},
&litMatcher{
	pos: position{line: 105, col: 91, offset: 2158},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 109, col: 1, offset: 2203},
	expr: &actionExpr{
	pos: position{line: 109, col: 14, offset: 2216},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 109, col: 14, offset: 2216},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 109, col: 14, offset: 2216},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 109, col: 17, offset: 2219},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 109, col: 17, offset: 2219},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 109, col: 26, offset: 2228},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 109, col: 33, offset: 2235},
	name: "WS",
},
&litMatcher{
	pos: position{line: 109, col: 36, offset: 2238},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 109, col: 40, offset: 2242},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 109, col: 43, offset: 2245},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 109, col: 46, offset: 2248},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 113, col: 1, offset: 2289},
	expr: &actionExpr{
	pos: position{line: 113, col: 14, offset: 2302},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 113, col: 14, offset: 2302},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 113, col: 17, offset: 2305},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 113, col: 17, offset: 2305},
	name: "Null",
},
&ruleRefExpr{
	pos: position{line: 113, col: 24, offset: 2312},
	name: "Boolean",
},
&ruleRefExpr{
	pos: position{line: 113, col: 34, offset: 2322},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 113, col: 43, offset: 2331},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 113, col: 51, offset: 2339},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 113, col: 61, offset: 2349},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 119, col: 1, offset: 2387},
	expr: &actionExpr{
	pos: position{line: 119, col: 14, offset: 2400},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 119, col: 14, offset: 2400},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 14, offset: 2400},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 119, col: 22, offset: 2408},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 119, col: 29, offset: 2415},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 119, col: 37, offset: 2423},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 119, col: 40, offset: 2426},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 119, col: 48, offset: 2434},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 119, col: 51, offset: 2437},
	expr: &seqExpr{
	pos: position{line: 119, col: 52, offset: 2438},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 52, offset: 2438},
	name: "WS",
},
&notExpr{
	pos: position{line: 119, col: 55, offset: 2441},
	expr: &choiceExpr{
	pos: position{line: 119, col: 57, offset: 2443},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 57, offset: 2443},
	name: "FLAGS_RULE",
},
&seqExpr{
	pos: position{line: 119, col: 70, offset: 2456},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 70, offset: 2456},
	name: "BS",
},
&ruleRefExpr{
	pos: position{line: 119, col: 73, offset: 2459},
	name: "BLOCK",
},
	},
},
	},
},
},
&choiceExpr{
	pos: position{line: 119, col: 81, offset: 2467},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 119, col: 81, offset: 2467},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 81, offset: 2467},
	name: "LS",
},
&zeroOrMoreExpr{
	pos: position{line: 119, col: 84, offset: 2470},
	expr: &seqExpr{
	pos: position{line: 119, col: 85, offset: 2471},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 85, offset: 2471},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 119, col: 88, offset: 2474},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 119, col: 91, offset: 2477},
	name: "WS",
},
	},
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 119, col: 98, offset: 2484},
	name: "LS",
},
	},
},
&ruleRefExpr{
	pos: position{line: 119, col: 102, offset: 2488},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 119, col: 105, offset: 2491},
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
	pos: position{line: 123, col: 1, offset: 2528},
	expr: &actionExpr{
	pos: position{line: 123, col: 11, offset: 2538},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 123, col: 11, offset: 2538},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 123, col: 11, offset: 2538},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 123, col: 14, offset: 2541},
	name: "FILTER_VALUE",
},
},
&labeledExpr{
	pos: position{line: 123, col: 28, offset: 2555},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 123, col: 32, offset: 2559},
	expr: &ruleRefExpr{
	pos: position{line: 123, col: 32, offset: 2559},
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
	pos: position{line: 127, col: 1, offset: 2602},
	expr: &actionExpr{
	pos: position{line: 127, col: 17, offset: 2618},
	run: (*parser).callonFILTER_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 127, col: 17, offset: 2618},
	label: "fv",
	expr: &choiceExpr{
	pos: position{line: 127, col: 21, offset: 2622},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 127, col: 21, offset: 2622},
	name: "IDENT_WITH_DOT",
},
&litMatcher{
	pos: position{line: 127, col: 38, offset: 2639},
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
	pos: position{line: 131, col: 1, offset: 2676},
	expr: &actionExpr{
	pos: position{line: 131, col: 15, offset: 2690},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 131, col: 15, offset: 2690},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 15, offset: 2690},
	name: "WS",
},
&litMatcher{
	pos: position{line: 131, col: 18, offset: 2693},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 131, col: 23, offset: 2698},
	name: "WS",
},
&litMatcher{
	pos: position{line: 131, col: 26, offset: 2701},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 131, col: 36, offset: 2711},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 131, col: 40, offset: 2715},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 131, col: 45, offset: 2720},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 131, col: 53, offset: 2728},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 135, col: 1, offset: 2769},
	expr: &actionExpr{
	pos: position{line: 135, col: 12, offset: 2780},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 135, col: 12, offset: 2780},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 135, col: 12, offset: 2780},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 135, col: 20, offset: 2788},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 135, col: 30, offset: 2798},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 135, col: 38, offset: 2806},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 135, col: 41, offset: 2809},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 135, col: 49, offset: 2817},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 135, col: 52, offset: 2820},
	expr: &seqExpr{
	pos: position{line: 135, col: 53, offset: 2821},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 135, col: 53, offset: 2821},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 135, col: 56, offset: 2824},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 135, col: 59, offset: 2827},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 135, col: 62, offset: 2830},
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
	pos: position{line: 139, col: 1, offset: 2870},
	expr: &actionExpr{
	pos: position{line: 139, col: 11, offset: 2880},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 139, col: 11, offset: 2880},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 139, col: 11, offset: 2880},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 139, col: 14, offset: 2883},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 139, col: 21, offset: 2890},
	name: "WS",
},
&litMatcher{
	pos: position{line: 139, col: 24, offset: 2893},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 139, col: 28, offset: 2897},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 139, col: 31, offset: 2900},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 139, col: 34, offset: 2903},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 139, col: 34, offset: 2903},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 139, col: 45, offset: 2914},
	name: "CHAIN",
},
&ruleRefExpr{
	pos: position{line: 139, col: 53, offset: 2922},
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
	pos: position{line: 143, col: 1, offset: 2959},
	expr: &actionExpr{
	pos: position{line: 143, col: 16, offset: 2974},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 143, col: 16, offset: 2974},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 143, col: 16, offset: 2974},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 143, col: 24, offset: 2982},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 147, col: 1, offset: 3016},
	expr: &actionExpr{
	pos: position{line: 147, col: 12, offset: 3027},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 147, col: 12, offset: 3027},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 147, col: 12, offset: 3027},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 147, col: 20, offset: 3035},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 147, col: 30, offset: 3045},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 147, col: 38, offset: 3053},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 147, col: 41, offset: 3056},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 147, col: 41, offset: 3056},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 147, col: 52, offset: 3067},
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
	pos: position{line: 151, col: 1, offset: 3103},
	expr: &actionExpr{
	pos: position{line: 151, col: 12, offset: 3114},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 151, col: 12, offset: 3114},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 151, col: 12, offset: 3114},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 151, col: 20, offset: 3122},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 151, col: 30, offset: 3132},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 151, col: 38, offset: 3140},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 151, col: 41, offset: 3143},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 151, col: 41, offset: 3143},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 151, col: 52, offset: 3154},
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
	pos: position{line: 155, col: 1, offset: 3189},
	expr: &actionExpr{
	pos: position{line: 155, col: 14, offset: 3202},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 155, col: 14, offset: 3202},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 155, col: 14, offset: 3202},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 155, col: 22, offset: 3210},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 155, col: 34, offset: 3222},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 155, col: 42, offset: 3230},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 155, col: 45, offset: 3233},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 155, col: 45, offset: 3233},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 155, col: 56, offset: 3244},
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
	pos: position{line: 159, col: 1, offset: 3280},
	expr: &actionExpr{
	pos: position{line: 159, col: 15, offset: 3294},
	run: (*parser).callonFLAGS_RULE1,
	expr: &seqExpr{
	pos: position{line: 159, col: 15, offset: 3294},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 159, col: 15, offset: 3294},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 159, col: 23, offset: 3302},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 159, col: 25, offset: 3304},
	name: "IGNORE_FLAG",
},
},
&labeledExpr{
	pos: position{line: 159, col: 37, offset: 3316},
	label: "is",
	expr: &zeroOrMoreExpr{
	pos: position{line: 159, col: 40, offset: 3319},
	expr: &seqExpr{
	pos: position{line: 159, col: 41, offset: 3320},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 159, col: 41, offset: 3320},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 159, col: 44, offset: 3323},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 159, col: 47, offset: 3326},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 159, col: 50, offset: 3329},
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
	pos: position{line: 163, col: 1, offset: 3372},
	expr: &actionExpr{
	pos: position{line: 163, col: 16, offset: 3387},
	run: (*parser).callonIGNORE_FLAG1,
	expr: &litMatcher{
	pos: position{line: 163, col: 16, offset: 3387},
	val: "ignore-errors",
	ignoreCase: false,
},
},
},
{
	name: "CHAIN",
	pos: position{line: 167, col: 1, offset: 3434},
	expr: &actionExpr{
	pos: position{line: 167, col: 10, offset: 3443},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 167, col: 10, offset: 3443},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 167, col: 10, offset: 3443},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 167, col: 13, offset: 3446},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 167, col: 27, offset: 3460},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 167, col: 30, offset: 3463},
	expr: &seqExpr{
	pos: position{line: 167, col: 31, offset: 3464},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 167, col: 31, offset: 3464},
	expr: &litMatcher{
	pos: position{line: 167, col: 31, offset: 3464},
	val: ".",
	ignoreCase: false,
},
},
&ruleRefExpr{
	pos: position{line: 167, col: 36, offset: 3469},
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
	pos: position{line: 171, col: 1, offset: 3513},
	expr: &actionExpr{
	pos: position{line: 171, col: 17, offset: 3529},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 171, col: 17, offset: 3529},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 171, col: 21, offset: 3533},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 171, col: 21, offset: 3533},
	name: "PATH_VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 171, col: 37, offset: 3549},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "PATH_VARIABLE",
	pos: position{line: 175, col: 1, offset: 3584},
	expr: &actionExpr{
	pos: position{line: 175, col: 18, offset: 3601},
	run: (*parser).callonPATH_VARIABLE1,
	expr: &seqExpr{
	pos: position{line: 175, col: 18, offset: 3601},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 175, col: 18, offset: 3601},
	expr: &litMatcher{
	pos: position{line: 175, col: 18, offset: 3601},
	val: "[",
	ignoreCase: false,
},
},
&litMatcher{
	pos: position{line: 175, col: 23, offset: 3606},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 175, col: 27, offset: 3610},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 175, col: 30, offset: 3613},
	name: "IDENT",
},
},
&zeroOrOneExpr{
	pos: position{line: 175, col: 37, offset: 3620},
	expr: &litMatcher{
	pos: position{line: 175, col: 37, offset: 3620},
	val: "]",
	ignoreCase: false,
},
},
	},
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 179, col: 1, offset: 3662},
	expr: &actionExpr{
	pos: position{line: 179, col: 13, offset: 3674},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 179, col: 13, offset: 3674},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 179, col: 13, offset: 3674},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 179, col: 17, offset: 3678},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 179, col: 20, offset: 3681},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 183, col: 1, offset: 3725},
	expr: &actionExpr{
	pos: position{line: 183, col: 10, offset: 3734},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 183, col: 10, offset: 3734},
	expr: &charClassMatcher{
	pos: position{line: 183, col: 10, offset: 3734},
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
	pos: position{line: 187, col: 1, offset: 3780},
	expr: &actionExpr{
	pos: position{line: 187, col: 19, offset: 3798},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 187, col: 19, offset: 3798},
	expr: &charClassMatcher{
	pos: position{line: 187, col: 19, offset: 3798},
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
	pos: position{line: 191, col: 1, offset: 3845},
	expr: &actionExpr{
	pos: position{line: 191, col: 9, offset: 3853},
	run: (*parser).callonNull1,
	expr: &litMatcher{
	pos: position{line: 191, col: 9, offset: 3853},
	val: "null",
	ignoreCase: false,
},
},
},
{
	name: "Boolean",
	pos: position{line: 195, col: 1, offset: 3883},
	expr: &actionExpr{
	pos: position{line: 195, col: 12, offset: 3894},
	run: (*parser).callonBoolean1,
	expr: &choiceExpr{
	pos: position{line: 195, col: 13, offset: 3895},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 195, col: 13, offset: 3895},
	val: "true",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 195, col: 22, offset: 3904},
	val: "false",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "String",
	pos: position{line: 199, col: 1, offset: 3945},
	expr: &actionExpr{
	pos: position{line: 199, col: 11, offset: 3955},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 199, col: 11, offset: 3955},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 199, col: 11, offset: 3955},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 199, col: 15, offset: 3959},
	expr: &seqExpr{
	pos: position{line: 199, col: 17, offset: 3961},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 199, col: 17, offset: 3961},
	expr: &litMatcher{
	pos: position{line: 199, col: 18, offset: 3962},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 199, col: 22, offset: 3966,
},
	},
},
},
&litMatcher{
	pos: position{line: 199, col: 27, offset: 3971},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 203, col: 1, offset: 4006},
	expr: &actionExpr{
	pos: position{line: 203, col: 10, offset: 4015},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 203, col: 10, offset: 4015},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 203, col: 10, offset: 4015},
	expr: &choiceExpr{
	pos: position{line: 203, col: 11, offset: 4016},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 203, col: 11, offset: 4016},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 203, col: 17, offset: 4022},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 203, col: 23, offset: 4028},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 203, col: 31, offset: 4036},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 203, col: 35, offset: 4040},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 207, col: 1, offset: 4078},
	expr: &actionExpr{
	pos: position{line: 207, col: 12, offset: 4089},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 207, col: 12, offset: 4089},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 207, col: 12, offset: 4089},
	expr: &choiceExpr{
	pos: position{line: 207, col: 13, offset: 4090},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 207, col: 13, offset: 4090},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 207, col: 19, offset: 4096},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 207, col: 25, offset: 4102},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 211, col: 1, offset: 4142},
	expr: &choiceExpr{
	pos: position{line: 211, col: 11, offset: 4154},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 211, col: 11, offset: 4154},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 211, col: 17, offset: 4160},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 211, col: 17, offset: 4160},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 211, col: 37, offset: 4180},
	expr: &ruleRefExpr{
	pos: position{line: 211, col: 37, offset: 4180},
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
	pos: position{line: 213, col: 1, offset: 4195},
	expr: &charClassMatcher{
	pos: position{line: 213, col: 16, offset: 4212},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 214, col: 1, offset: 4218},
	expr: &charClassMatcher{
	pos: position{line: 214, col: 23, offset: 4242},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 216, col: 1, offset: 4249},
	expr: &charClassMatcher{
	pos: position{line: 216, col: 10, offset: 4258},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 217, col: 1, offset: 4264},
	expr: &oneOrMoreExpr{
	pos: position{line: 217, col: 35, offset: 4298},
	expr: &choiceExpr{
	pos: position{line: 217, col: 36, offset: 4299},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 36, offset: 4299},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 217, col: 44, offset: 4307},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 217, col: 54, offset: 4317},
	name: "NL",
},
	},
},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 218, col: 1, offset: 4322},
	expr: &zeroOrMoreExpr{
	pos: position{line: 218, col: 20, offset: 4341},
	expr: &choiceExpr{
	pos: position{line: 218, col: 21, offset: 4342},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 218, col: 21, offset: 4342},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 218, col: 29, offset: 4350},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 219, col: 1, offset: 4360},
	expr: &choiceExpr{
	pos: position{line: 219, col: 25, offset: 4384},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 219, col: 25, offset: 4384},
	name: "NL",
},
&litMatcher{
	pos: position{line: 219, col: 30, offset: 4389},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 219, col: 36, offset: 4395},
	name: "COMMENT",
},
	},
},
},
{
	name: "BS",
	displayName: "\"block-separator\"",
	pos: position{line: 220, col: 1, offset: 4404},
	expr: &oneOrMoreExpr{
	pos: position{line: 220, col: 25, offset: 4428},
	expr: &seqExpr{
	pos: position{line: 220, col: 26, offset: 4429},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 220, col: 26, offset: 4429},
	name: "WS",
},
&choiceExpr{
	pos: position{line: 220, col: 30, offset: 4433},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 220, col: 30, offset: 4433},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 220, col: 35, offset: 4438},
	name: "COMMENT",
},
	},
},
&ruleRefExpr{
	pos: position{line: 220, col: 44, offset: 4447},
	name: "WS",
},
	},
},
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 221, col: 1, offset: 4452},
	expr: &litMatcher{
	pos: position{line: 221, col: 18, offset: 4469},
	val: "\n",
	ignoreCase: false,
},
},
{
	name: "COMMENT",
	pos: position{line: 223, col: 1, offset: 4475},
	expr: &seqExpr{
	pos: position{line: 223, col: 12, offset: 4486},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 223, col: 12, offset: 4486},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 223, col: 17, offset: 4491},
	expr: &seqExpr{
	pos: position{line: 223, col: 19, offset: 4493},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 223, col: 19, offset: 4493},
	expr: &litMatcher{
	pos: position{line: 223, col: 20, offset: 4494},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 223, col: 25, offset: 4499,
},
	},
},
},
&choiceExpr{
	pos: position{line: 223, col: 31, offset: 4505},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 223, col: 31, offset: 4505},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 223, col: 38, offset: 4512},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 225, col: 1, offset: 4518},
	expr: &notExpr{
	pos: position{line: 225, col: 8, offset: 4525},
	expr: &anyMatcher{
	line: 225, col: 9, offset: 4526,
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

func (c *current) onPATH_VARIABLE1(i interface{}) (interface{}, error) {
	return newChainPathVariable(i)
}

func (p *parser) callonPATH_VARIABLE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPATH_VARIABLE1(stack["i"])
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

