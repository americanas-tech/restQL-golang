
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
	expr: &choiceExpr{
	pos: position{line: 21, col: 8, offset: 310},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 21, col: 8, offset: 310},
	run: (*parser).callonUSE2,
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
&actionExpr{
	pos: position{line: 23, col: 5, offset: 393},
	run: (*parser).callonUSE15,
	expr: &seqExpr{
	pos: position{line: 23, col: 5, offset: 393},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 23, col: 5, offset: 393},
	val: "use",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 23, col: 11, offset: 399},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 23, col: 19, offset: 407},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 23, col: 22, offset: 410},
	name: "USE_ACTION",
},
},
&ruleRefExpr{
	pos: position{line: 23, col: 34, offset: 422},
	name: "WS",
},
&litMatcher{
	pos: position{line: 23, col: 37, offset: 425},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 23, col: 41, offset: 429},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 23, col: 44, offset: 432},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 23, col: 47, offset: 435},
	name: "USE_VALUE",
},
},
&ruleRefExpr{
	pos: position{line: 23, col: 58, offset: 446},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 23, col: 61, offset: 449},
	expr: &ruleRefExpr{
	pos: position{line: 23, col: 61, offset: 449},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 23, col: 65, offset: 453},
	name: "WS",
},
	},
},
},
	},
},
},
{
	name: "USE_ACTION",
	pos: position{line: 27, col: 1, offset: 531},
	expr: &choiceExpr{
	pos: position{line: 27, col: 15, offset: 545},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 27, col: 15, offset: 545},
	run: (*parser).callonUSE_ACTION2,
	expr: &choiceExpr{
	pos: position{line: 27, col: 16, offset: 546},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 27, col: 16, offset: 546},
	val: "timeout",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 27, col: 28, offset: 558},
	val: "max-age",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 27, col: 40, offset: 570},
	val: "s-max-age",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 29, col: 5, offset: 615},
	run: (*parser).callonUSE_ACTION7,
	expr: &litMatcher{
	pos: position{line: 29, col: 5, offset: 615},
	val: "cache-control",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "USE_VALUE",
	pos: position{line: 33, col: 1, offset: 702},
	expr: &actionExpr{
	pos: position{line: 33, col: 14, offset: 715},
	run: (*parser).callonUSE_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 33, col: 14, offset: 715},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 33, col: 17, offset: 718},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 33, col: 17, offset: 718},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 33, col: 26, offset: 727},
	name: "Integer",
},
	},
},
},
},
},
{
	name: "BLOCK",
	pos: position{line: 37, col: 1, offset: 764},
	expr: &actionExpr{
	pos: position{line: 37, col: 10, offset: 773},
	run: (*parser).callonBLOCK1,
	expr: &seqExpr{
	pos: position{line: 37, col: 10, offset: 773},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 37, col: 10, offset: 773},
	label: "action",
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 18, offset: 781},
	name: "ACTION_RULE",
},
},
&labeledExpr{
	pos: position{line: 37, col: 31, offset: 794},
	label: "m",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 34, offset: 797},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 34, offset: 797},
	name: "MODIFIER_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 37, col: 50, offset: 813},
	label: "w",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 53, offset: 816},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 53, offset: 816},
	name: "WITH_RULE",
},
},
},
&labeledExpr{
	pos: position{line: 37, col: 65, offset: 828},
	label: "f",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 67, offset: 830},
	expr: &choiceExpr{
	pos: position{line: 37, col: 68, offset: 831},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 37, col: 68, offset: 831},
	name: "HIDDEN_RULE",
},
&ruleRefExpr{
	pos: position{line: 37, col: 82, offset: 845},
	name: "ONLY_RULE",
},
	},
},
},
},
&labeledExpr{
	pos: position{line: 37, col: 94, offset: 857},
	label: "fl",
	expr: &zeroOrOneExpr{
	pos: position{line: 37, col: 98, offset: 861},
	expr: &ruleRefExpr{
	pos: position{line: 37, col: 98, offset: 861},
	name: "FLAGS_RULE",
},
},
},
&ruleRefExpr{
	pos: position{line: 37, col: 111, offset: 874},
	name: "WS",
},
	},
},
},
},
{
	name: "ACTION_RULE",
	pos: position{line: 41, col: 1, offset: 920},
	expr: &actionExpr{
	pos: position{line: 41, col: 16, offset: 935},
	run: (*parser).callonACTION_RULE1,
	expr: &seqExpr{
	pos: position{line: 41, col: 16, offset: 935},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 41, col: 16, offset: 935},
	label: "m",
	expr: &ruleRefExpr{
	pos: position{line: 41, col: 19, offset: 938},
	name: "METHOD",
},
},
&ruleRefExpr{
	pos: position{line: 41, col: 27, offset: 946},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 41, col: 35, offset: 954},
	label: "r",
	expr: &ruleRefExpr{
	pos: position{line: 41, col: 38, offset: 957},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 41, col: 45, offset: 964},
	label: "a",
	expr: &zeroOrOneExpr{
	pos: position{line: 41, col: 48, offset: 967},
	expr: &ruleRefExpr{
	pos: position{line: 41, col: 48, offset: 967},
	name: "ALIAS",
},
},
},
&labeledExpr{
	pos: position{line: 41, col: 56, offset: 975},
	label: "i",
	expr: &zeroOrOneExpr{
	pos: position{line: 41, col: 59, offset: 978},
	expr: &ruleRefExpr{
	pos: position{line: 41, col: 59, offset: 978},
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
	pos: position{line: 45, col: 1, offset: 1022},
	expr: &actionExpr{
	pos: position{line: 45, col: 11, offset: 1032},
	run: (*parser).callonMETHOD1,
	expr: &choiceExpr{
	pos: position{line: 45, col: 12, offset: 1033},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 45, col: 12, offset: 1033},
	val: "from",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 45, col: 21, offset: 1042},
	val: "to",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 45, col: 28, offset: 1049},
	val: "into",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 45, col: 36, offset: 1057},
	val: "update",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 45, col: 47, offset: 1068},
	val: "delete",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "ALIAS",
	pos: position{line: 49, col: 1, offset: 1109},
	expr: &actionExpr{
	pos: position{line: 49, col: 10, offset: 1118},
	run: (*parser).callonALIAS1,
	expr: &seqExpr{
	pos: position{line: 49, col: 10, offset: 1118},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 49, col: 10, offset: 1118},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 49, col: 18, offset: 1126},
	val: "as",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 49, col: 23, offset: 1131},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 49, col: 31, offset: 1139},
	label: "a",
	expr: &ruleRefExpr{
	pos: position{line: 49, col: 34, offset: 1142},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IN",
	pos: position{line: 53, col: 1, offset: 1169},
	expr: &actionExpr{
	pos: position{line: 53, col: 7, offset: 1175},
	run: (*parser).callonIN1,
	expr: &seqExpr{
	pos: position{line: 53, col: 7, offset: 1175},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 53, col: 7, offset: 1175},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 53, col: 15, offset: 1183},
	val: "in",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 53, col: 20, offset: 1188},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 53, col: 28, offset: 1196},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 53, col: 31, offset: 1199},
	name: "IDENT_WITH_DOT",
},
},
	},
},
},
},
{
	name: "MODIFIER_RULE",
	pos: position{line: 57, col: 1, offset: 1237},
	expr: &actionExpr{
	pos: position{line: 57, col: 18, offset: 1254},
	run: (*parser).callonMODIFIER_RULE1,
	expr: &labeledExpr{
	pos: position{line: 57, col: 18, offset: 1254},
	label: "m",
	expr: &oneOrMoreExpr{
	pos: position{line: 57, col: 20, offset: 1256},
	expr: &choiceExpr{
	pos: position{line: 57, col: 21, offset: 1257},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 57, col: 21, offset: 1257},
	name: "HEADERS",
},
&ruleRefExpr{
	pos: position{line: 57, col: 31, offset: 1267},
	name: "TIMEOUT",
},
&ruleRefExpr{
	pos: position{line: 57, col: 41, offset: 1277},
	name: "MAX_AGE",
},
&ruleRefExpr{
	pos: position{line: 57, col: 51, offset: 1287},
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
	pos: position{line: 61, col: 1, offset: 1319},
	expr: &actionExpr{
	pos: position{line: 61, col: 14, offset: 1332},
	run: (*parser).callonWITH_RULE1,
	expr: &seqExpr{
	pos: position{line: 61, col: 14, offset: 1332},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 61, col: 14, offset: 1332},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 61, col: 22, offset: 1340},
	val: "with",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 61, col: 29, offset: 1347},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 61, col: 37, offset: 1355},
	label: "pb",
	expr: &zeroOrOneExpr{
	pos: position{line: 61, col: 40, offset: 1358},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 40, offset: 1358},
	name: "PARAMETER_BODY",
},
},
},
&labeledExpr{
	pos: position{line: 61, col: 56, offset: 1374},
	label: "kvs",
	expr: &zeroOrOneExpr{
	pos: position{line: 61, col: 60, offset: 1378},
	expr: &ruleRefExpr{
	pos: position{line: 61, col: 60, offset: 1378},
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
	pos: position{line: 65, col: 1, offset: 1424},
	expr: &actionExpr{
	pos: position{line: 65, col: 19, offset: 1442},
	run: (*parser).callonPARAMETER_BODY1,
	expr: &seqExpr{
	pos: position{line: 65, col: 19, offset: 1442},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 65, col: 19, offset: 1442},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 65, col: 23, offset: 1446},
	label: "t",
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 26, offset: 1449},
	name: "IDENT",
},
},
&labeledExpr{
	pos: position{line: 65, col: 33, offset: 1456},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 65, col: 37, offset: 1460},
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 37, offset: 1460},
	name: "APPLY_FN",
},
},
},
&ruleRefExpr{
	pos: position{line: 65, col: 48, offset: 1471},
	name: "WS",
},
&zeroOrOneExpr{
	pos: position{line: 65, col: 51, offset: 1474},
	expr: &ruleRefExpr{
	pos: position{line: 65, col: 51, offset: 1474},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 65, col: 55, offset: 1478},
	name: "WS",
},
	},
},
},
},
{
	name: "KEY_VALUE_LIST",
	pos: position{line: 69, col: 1, offset: 1518},
	expr: &actionExpr{
	pos: position{line: 69, col: 19, offset: 1536},
	run: (*parser).callonKEY_VALUE_LIST1,
	expr: &seqExpr{
	pos: position{line: 69, col: 19, offset: 1536},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 69, col: 19, offset: 1536},
	label: "first",
	expr: &ruleRefExpr{
	pos: position{line: 69, col: 25, offset: 1542},
	name: "KEY_VALUE",
},
},
&labeledExpr{
	pos: position{line: 69, col: 35, offset: 1552},
	label: "others",
	expr: &zeroOrMoreExpr{
	pos: position{line: 69, col: 42, offset: 1559},
	expr: &seqExpr{
	pos: position{line: 69, col: 43, offset: 1560},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 69, col: 43, offset: 1560},
	name: "WS",
},
&choiceExpr{
	pos: position{line: 69, col: 47, offset: 1564},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 69, col: 47, offset: 1564},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 69, col: 47, offset: 1564},
	name: "LS",
},
&zeroOrMoreExpr{
	pos: position{line: 69, col: 50, offset: 1567},
	expr: &seqExpr{
	pos: position{line: 69, col: 51, offset: 1568},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 69, col: 51, offset: 1568},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 69, col: 54, offset: 1571},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 69, col: 57, offset: 1574},
	name: "WS",
},
	},
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 69, col: 64, offset: 1581},
	name: "LS",
},
	},
},
&ruleRefExpr{
	pos: position{line: 69, col: 68, offset: 1585},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 69, col: 71, offset: 1588},
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
	pos: position{line: 73, col: 1, offset: 1644},
	expr: &actionExpr{
	pos: position{line: 73, col: 14, offset: 1657},
	run: (*parser).callonKEY_VALUE1,
	expr: &seqExpr{
	pos: position{line: 73, col: 14, offset: 1657},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 73, col: 14, offset: 1657},
	label: "k",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 17, offset: 1660},
	name: "IDENT_WITH_DOT",
},
},
&ruleRefExpr{
	pos: position{line: 73, col: 33, offset: 1676},
	name: "WS",
},
&litMatcher{
	pos: position{line: 73, col: 36, offset: 1679},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 73, col: 40, offset: 1683},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 73, col: 43, offset: 1686},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 46, offset: 1689},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 73, col: 53, offset: 1696},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 73, col: 57, offset: 1700},
	expr: &ruleRefExpr{
	pos: position{line: 73, col: 57, offset: 1700},
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
	pos: position{line: 77, col: 1, offset: 1746},
	expr: &actionExpr{
	pos: position{line: 77, col: 13, offset: 1758},
	run: (*parser).callonAPPLY_FN1,
	expr: &seqExpr{
	pos: position{line: 77, col: 13, offset: 1758},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 77, col: 13, offset: 1758},
	name: "WS",
},
&litMatcher{
	pos: position{line: 77, col: 16, offset: 1761},
	val: "->",
	ignoreCase: false,
},
&zeroOrOneExpr{
	pos: position{line: 77, col: 21, offset: 1766},
	expr: &ruleRefExpr{
	pos: position{line: 77, col: 21, offset: 1766},
	name: "WS",
},
},
&labeledExpr{
	pos: position{line: 77, col: 25, offset: 1770},
	label: "fn",
	expr: &ruleRefExpr{
	pos: position{line: 77, col: 29, offset: 1774},
	name: "FUNCTION",
},
},
	},
},
},
},
{
	name: "FUNCTION",
	pos: position{line: 81, col: 1, offset: 1805},
	expr: &choiceExpr{
	pos: position{line: 81, col: 13, offset: 1817},
	alternatives: []interface{}{
&actionExpr{
	pos: position{line: 81, col: 13, offset: 1817},
	run: (*parser).callonFUNCTION2,
	expr: &choiceExpr{
	pos: position{line: 81, col: 14, offset: 1818},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 81, col: 14, offset: 1818},
	val: "flatten",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 81, col: 26, offset: 1830},
	val: "base64",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 81, col: 37, offset: 1841},
	val: "json",
	ignoreCase: false,
},
	},
},
},
&actionExpr{
	pos: position{line: 83, col: 5, offset: 1881},
	run: (*parser).callonFUNCTION7,
	expr: &litMatcher{
	pos: position{line: 83, col: 5, offset: 1881},
	val: "contract",
	ignoreCase: false,
},
},
	},
},
},
{
	name: "VALUE",
	pos: position{line: 87, col: 1, offset: 1956},
	expr: &actionExpr{
	pos: position{line: 87, col: 10, offset: 1965},
	run: (*parser).callonVALUE1,
	expr: &labeledExpr{
	pos: position{line: 87, col: 10, offset: 1965},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 87, col: 13, offset: 1968},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 87, col: 13, offset: 1968},
	name: "LIST",
},
&ruleRefExpr{
	pos: position{line: 87, col: 20, offset: 1975},
	name: "OBJECT",
},
&ruleRefExpr{
	pos: position{line: 87, col: 29, offset: 1984},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 87, col: 40, offset: 1995},
	name: "PRIMITIVE",
},
	},
},
},
},
},
{
	name: "LIST",
	pos: position{line: 91, col: 1, offset: 2031},
	expr: &actionExpr{
	pos: position{line: 91, col: 9, offset: 2039},
	run: (*parser).callonLIST1,
	expr: &labeledExpr{
	pos: position{line: 91, col: 9, offset: 2039},
	label: "l",
	expr: &choiceExpr{
	pos: position{line: 91, col: 12, offset: 2042},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 91, col: 12, offset: 2042},
	name: "EMPTY_LIST",
},
&ruleRefExpr{
	pos: position{line: 91, col: 25, offset: 2055},
	name: "POPULATED_LIST",
},
	},
},
},
},
},
{
	name: "EMPTY_LIST",
	pos: position{line: 95, col: 1, offset: 2091},
	expr: &actionExpr{
	pos: position{line: 95, col: 15, offset: 2105},
	run: (*parser).callonEMPTY_LIST1,
	expr: &seqExpr{
	pos: position{line: 95, col: 15, offset: 2105},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 95, col: 15, offset: 2105},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 95, col: 19, offset: 2109},
	name: "WS",
},
&litMatcher{
	pos: position{line: 95, col: 22, offset: 2112},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_LIST",
	pos: position{line: 99, col: 1, offset: 2144},
	expr: &actionExpr{
	pos: position{line: 99, col: 19, offset: 2162},
	run: (*parser).callonPOPULATED_LIST1,
	expr: &seqExpr{
	pos: position{line: 99, col: 19, offset: 2162},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 99, col: 19, offset: 2162},
	val: "[",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 99, col: 23, offset: 2166},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 99, col: 26, offset: 2169},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 99, col: 28, offset: 2171},
	name: "VALUE",
},
},
&labeledExpr{
	pos: position{line: 99, col: 34, offset: 2177},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 99, col: 37, offset: 2180},
	expr: &seqExpr{
	pos: position{line: 99, col: 38, offset: 2181},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 99, col: 38, offset: 2181},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 99, col: 41, offset: 2184},
	expr: &ruleRefExpr{
	pos: position{line: 99, col: 41, offset: 2184},
	name: "LS",
},
},
&ruleRefExpr{
	pos: position{line: 99, col: 45, offset: 2188},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 99, col: 48, offset: 2191},
	name: "VALUE",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 99, col: 56, offset: 2199},
	name: "WS",
},
&litMatcher{
	pos: position{line: 99, col: 59, offset: 2202},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJECT",
	pos: position{line: 103, col: 1, offset: 2234},
	expr: &actionExpr{
	pos: position{line: 103, col: 11, offset: 2244},
	run: (*parser).callonOBJECT1,
	expr: &labeledExpr{
	pos: position{line: 103, col: 11, offset: 2244},
	label: "o",
	expr: &choiceExpr{
	pos: position{line: 103, col: 14, offset: 2247},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 103, col: 14, offset: 2247},
	name: "EMPTY_OBJ",
},
&ruleRefExpr{
	pos: position{line: 103, col: 26, offset: 2259},
	name: "POPULATED_OBJ",
},
	},
},
},
},
},
{
	name: "EMPTY_OBJ",
	pos: position{line: 107, col: 1, offset: 2294},
	expr: &actionExpr{
	pos: position{line: 107, col: 14, offset: 2307},
	run: (*parser).callonEMPTY_OBJ1,
	expr: &seqExpr{
	pos: position{line: 107, col: 14, offset: 2307},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 107, col: 14, offset: 2307},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 107, col: 18, offset: 2311},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 107, col: 21, offset: 2314},
	expr: &ruleRefExpr{
	pos: position{line: 107, col: 21, offset: 2314},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 107, col: 25, offset: 2318},
	name: "WS",
},
&litMatcher{
	pos: position{line: 107, col: 28, offset: 2321},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "POPULATED_OBJ",
	pos: position{line: 111, col: 1, offset: 2355},
	expr: &actionExpr{
	pos: position{line: 111, col: 18, offset: 2372},
	run: (*parser).callonPOPULATED_OBJ1,
	expr: &seqExpr{
	pos: position{line: 111, col: 18, offset: 2372},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 111, col: 18, offset: 2372},
	val: "{",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 111, col: 22, offset: 2376},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 111, col: 25, offset: 2379},
	expr: &ruleRefExpr{
	pos: position{line: 111, col: 25, offset: 2379},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 111, col: 29, offset: 2383},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 111, col: 32, offset: 2386},
	label: "oe",
	expr: &ruleRefExpr{
	pos: position{line: 111, col: 36, offset: 2390},
	name: "OBJ_ENTRY",
},
},
&labeledExpr{
	pos: position{line: 111, col: 47, offset: 2401},
	label: "oes",
	expr: &zeroOrMoreExpr{
	pos: position{line: 111, col: 51, offset: 2405},
	expr: &seqExpr{
	pos: position{line: 111, col: 52, offset: 2406},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 111, col: 52, offset: 2406},
	name: "WS",
},
&litMatcher{
	pos: position{line: 111, col: 55, offset: 2409},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 111, col: 59, offset: 2413},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 111, col: 62, offset: 2416},
	expr: &ruleRefExpr{
	pos: position{line: 111, col: 62, offset: 2416},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 111, col: 66, offset: 2420},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 111, col: 69, offset: 2423},
	name: "OBJ_ENTRY",
},
	},
},
},
},
&ruleRefExpr{
	pos: position{line: 111, col: 81, offset: 2435},
	name: "WS",
},
&zeroOrMoreExpr{
	pos: position{line: 111, col: 84, offset: 2438},
	expr: &ruleRefExpr{
	pos: position{line: 111, col: 84, offset: 2438},
	name: "NL",
},
},
&ruleRefExpr{
	pos: position{line: 111, col: 88, offset: 2442},
	name: "WS",
},
&litMatcher{
	pos: position{line: 111, col: 91, offset: 2445},
	val: "}",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OBJ_ENTRY",
	pos: position{line: 115, col: 1, offset: 2490},
	expr: &actionExpr{
	pos: position{line: 115, col: 14, offset: 2503},
	run: (*parser).callonOBJ_ENTRY1,
	expr: &seqExpr{
	pos: position{line: 115, col: 14, offset: 2503},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 115, col: 14, offset: 2503},
	label: "k",
	expr: &choiceExpr{
	pos: position{line: 115, col: 17, offset: 2506},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 115, col: 17, offset: 2506},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 115, col: 26, offset: 2515},
	name: "IDENT",
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 115, col: 33, offset: 2522},
	name: "WS",
},
&litMatcher{
	pos: position{line: 115, col: 36, offset: 2525},
	val: ":",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 115, col: 40, offset: 2529},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 115, col: 43, offset: 2532},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 115, col: 46, offset: 2535},
	name: "VALUE",
},
},
	},
},
},
},
{
	name: "PRIMITIVE",
	pos: position{line: 119, col: 1, offset: 2576},
	expr: &actionExpr{
	pos: position{line: 119, col: 14, offset: 2589},
	run: (*parser).callonPRIMITIVE1,
	expr: &labeledExpr{
	pos: position{line: 119, col: 14, offset: 2589},
	label: "p",
	expr: &choiceExpr{
	pos: position{line: 119, col: 17, offset: 2592},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 119, col: 17, offset: 2592},
	name: "Null",
},
&ruleRefExpr{
	pos: position{line: 119, col: 24, offset: 2599},
	name: "Boolean",
},
&ruleRefExpr{
	pos: position{line: 119, col: 34, offset: 2609},
	name: "String",
},
&ruleRefExpr{
	pos: position{line: 119, col: 43, offset: 2618},
	name: "Float",
},
&ruleRefExpr{
	pos: position{line: 119, col: 51, offset: 2626},
	name: "Integer",
},
&ruleRefExpr{
	pos: position{line: 119, col: 61, offset: 2636},
	name: "CHAIN",
},
	},
},
},
},
},
{
	name: "CHAIN",
	pos: position{line: 123, col: 1, offset: 2672},
	expr: &actionExpr{
	pos: position{line: 123, col: 10, offset: 2681},
	run: (*parser).callonCHAIN1,
	expr: &seqExpr{
	pos: position{line: 123, col: 10, offset: 2681},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 123, col: 10, offset: 2681},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 123, col: 13, offset: 2684},
	name: "CHAINED_ITEM",
},
},
&labeledExpr{
	pos: position{line: 123, col: 27, offset: 2698},
	label: "ii",
	expr: &zeroOrMoreExpr{
	pos: position{line: 123, col: 30, offset: 2701},
	expr: &seqExpr{
	pos: position{line: 123, col: 31, offset: 2702},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 123, col: 31, offset: 2702},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 123, col: 35, offset: 2706},
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
	pos: position{line: 127, col: 1, offset: 2750},
	expr: &actionExpr{
	pos: position{line: 127, col: 17, offset: 2766},
	run: (*parser).callonCHAINED_ITEM1,
	expr: &labeledExpr{
	pos: position{line: 127, col: 17, offset: 2766},
	label: "ci",
	expr: &choiceExpr{
	pos: position{line: 127, col: 21, offset: 2770},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 127, col: 21, offset: 2770},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 127, col: 32, offset: 2781},
	name: "IDENT",
},
	},
},
},
},
},
{
	name: "ONLY_RULE",
	pos: position{line: 131, col: 1, offset: 2816},
	expr: &actionExpr{
	pos: position{line: 131, col: 14, offset: 2829},
	run: (*parser).callonONLY_RULE1,
	expr: &seqExpr{
	pos: position{line: 131, col: 14, offset: 2829},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 14, offset: 2829},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 131, col: 22, offset: 2837},
	val: "only",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 131, col: 29, offset: 2844},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 131, col: 37, offset: 2852},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 131, col: 40, offset: 2855},
	name: "FILTER",
},
},
&labeledExpr{
	pos: position{line: 131, col: 48, offset: 2863},
	label: "fs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 131, col: 51, offset: 2866},
	expr: &seqExpr{
	pos: position{line: 131, col: 52, offset: 2867},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 52, offset: 2867},
	name: "WS",
},
&notExpr{
	pos: position{line: 131, col: 55, offset: 2870},
	expr: &choiceExpr{
	pos: position{line: 131, col: 57, offset: 2872},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 57, offset: 2872},
	name: "FLAGS_RULE",
},
&seqExpr{
	pos: position{line: 131, col: 70, offset: 2885},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 70, offset: 2885},
	name: "BS",
},
&ruleRefExpr{
	pos: position{line: 131, col: 73, offset: 2888},
	name: "BLOCK",
},
	},
},
	},
},
},
&choiceExpr{
	pos: position{line: 131, col: 81, offset: 2896},
	alternatives: []interface{}{
&seqExpr{
	pos: position{line: 131, col: 81, offset: 2896},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 81, offset: 2896},
	name: "LS",
},
&zeroOrMoreExpr{
	pos: position{line: 131, col: 84, offset: 2899},
	expr: &seqExpr{
	pos: position{line: 131, col: 85, offset: 2900},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 131, col: 85, offset: 2900},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 131, col: 88, offset: 2903},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 131, col: 91, offset: 2906},
	name: "WS",
},
	},
},
},
	},
},
&ruleRefExpr{
	pos: position{line: 131, col: 98, offset: 2913},
	name: "LS",
},
	},
},
&ruleRefExpr{
	pos: position{line: 131, col: 102, offset: 2917},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 131, col: 105, offset: 2920},
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
	pos: position{line: 135, col: 1, offset: 2957},
	expr: &actionExpr{
	pos: position{line: 135, col: 11, offset: 2967},
	run: (*parser).callonFILTER1,
	expr: &seqExpr{
	pos: position{line: 135, col: 11, offset: 2967},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 135, col: 11, offset: 2967},
	label: "f",
	expr: &ruleRefExpr{
	pos: position{line: 135, col: 14, offset: 2970},
	name: "FILTER_VALUE",
},
},
&labeledExpr{
	pos: position{line: 135, col: 28, offset: 2984},
	label: "fn",
	expr: &zeroOrOneExpr{
	pos: position{line: 135, col: 32, offset: 2988},
	expr: &ruleRefExpr{
	pos: position{line: 135, col: 32, offset: 2988},
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
	pos: position{line: 139, col: 1, offset: 3031},
	expr: &actionExpr{
	pos: position{line: 139, col: 17, offset: 3047},
	run: (*parser).callonFILTER_VALUE1,
	expr: &labeledExpr{
	pos: position{line: 139, col: 17, offset: 3047},
	label: "fv",
	expr: &choiceExpr{
	pos: position{line: 139, col: 21, offset: 3051},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 139, col: 21, offset: 3051},
	name: "IDENT_WITH_DOT",
},
&litMatcher{
	pos: position{line: 139, col: 38, offset: 3068},
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
	pos: position{line: 143, col: 1, offset: 3105},
	expr: &actionExpr{
	pos: position{line: 143, col: 15, offset: 3119},
	run: (*parser).callonMATCHES_FN1,
	expr: &seqExpr{
	pos: position{line: 143, col: 15, offset: 3119},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 143, col: 15, offset: 3119},
	name: "WS",
},
&litMatcher{
	pos: position{line: 143, col: 18, offset: 3122},
	val: "->",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 143, col: 23, offset: 3127},
	name: "WS",
},
&litMatcher{
	pos: position{line: 143, col: 26, offset: 3130},
	val: "matches",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 143, col: 36, offset: 3140},
	val: "(",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 143, col: 40, offset: 3144},
	label: "arg",
	expr: &ruleRefExpr{
	pos: position{line: 143, col: 45, offset: 3149},
	name: "String",
},
},
&litMatcher{
	pos: position{line: 143, col: 53, offset: 3157},
	val: ")",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "HEADERS",
	pos: position{line: 147, col: 1, offset: 3198},
	expr: &actionExpr{
	pos: position{line: 147, col: 12, offset: 3209},
	run: (*parser).callonHEADERS1,
	expr: &seqExpr{
	pos: position{line: 147, col: 12, offset: 3209},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 147, col: 12, offset: 3209},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 147, col: 20, offset: 3217},
	val: "headers",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 147, col: 30, offset: 3227},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 147, col: 38, offset: 3235},
	label: "h",
	expr: &ruleRefExpr{
	pos: position{line: 147, col: 41, offset: 3238},
	name: "HEADER",
},
},
&labeledExpr{
	pos: position{line: 147, col: 49, offset: 3246},
	label: "hs",
	expr: &zeroOrMoreExpr{
	pos: position{line: 147, col: 52, offset: 3249},
	expr: &seqExpr{
	pos: position{line: 147, col: 53, offset: 3250},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 147, col: 53, offset: 3250},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 147, col: 56, offset: 3253},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 147, col: 59, offset: 3256},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 147, col: 62, offset: 3259},
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
	pos: position{line: 151, col: 1, offset: 3299},
	expr: &actionExpr{
	pos: position{line: 151, col: 11, offset: 3309},
	run: (*parser).callonHEADER1,
	expr: &seqExpr{
	pos: position{line: 151, col: 11, offset: 3309},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 151, col: 11, offset: 3309},
	label: "n",
	expr: &ruleRefExpr{
	pos: position{line: 151, col: 14, offset: 3312},
	name: "IDENT",
},
},
&ruleRefExpr{
	pos: position{line: 151, col: 21, offset: 3319},
	name: "WS",
},
&litMatcher{
	pos: position{line: 151, col: 24, offset: 3322},
	val: "=",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 151, col: 28, offset: 3326},
	name: "WS",
},
&labeledExpr{
	pos: position{line: 151, col: 31, offset: 3329},
	label: "v",
	expr: &choiceExpr{
	pos: position{line: 151, col: 34, offset: 3332},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 151, col: 34, offset: 3332},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 151, col: 45, offset: 3343},
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
	pos: position{line: 155, col: 1, offset: 3380},
	expr: &actionExpr{
	pos: position{line: 155, col: 16, offset: 3395},
	run: (*parser).callonHIDDEN_RULE1,
	expr: &seqExpr{
	pos: position{line: 155, col: 16, offset: 3395},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 155, col: 16, offset: 3395},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 155, col: 24, offset: 3403},
	val: "hidden",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "TIMEOUT",
	pos: position{line: 159, col: 1, offset: 3437},
	expr: &actionExpr{
	pos: position{line: 159, col: 12, offset: 3448},
	run: (*parser).callonTIMEOUT1,
	expr: &seqExpr{
	pos: position{line: 159, col: 12, offset: 3448},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 159, col: 12, offset: 3448},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 159, col: 20, offset: 3456},
	val: "timeout",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 159, col: 30, offset: 3466},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 159, col: 38, offset: 3474},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 159, col: 41, offset: 3477},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 159, col: 41, offset: 3477},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 159, col: 52, offset: 3488},
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
	pos: position{line: 163, col: 1, offset: 3524},
	expr: &actionExpr{
	pos: position{line: 163, col: 12, offset: 3535},
	run: (*parser).callonMAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 163, col: 12, offset: 3535},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 163, col: 12, offset: 3535},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 163, col: 20, offset: 3543},
	val: "max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 163, col: 30, offset: 3553},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 163, col: 38, offset: 3561},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 163, col: 41, offset: 3564},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 163, col: 41, offset: 3564},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 163, col: 52, offset: 3575},
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
	pos: position{line: 167, col: 1, offset: 3610},
	expr: &actionExpr{
	pos: position{line: 167, col: 14, offset: 3623},
	run: (*parser).callonS_MAX_AGE1,
	expr: &seqExpr{
	pos: position{line: 167, col: 14, offset: 3623},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 167, col: 14, offset: 3623},
	name: "WS_MAND",
},
&litMatcher{
	pos: position{line: 167, col: 22, offset: 3631},
	val: "s-max-age",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 167, col: 34, offset: 3643},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 167, col: 42, offset: 3651},
	label: "t",
	expr: &choiceExpr{
	pos: position{line: 167, col: 45, offset: 3654},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 167, col: 45, offset: 3654},
	name: "VARIABLE",
},
&ruleRefExpr{
	pos: position{line: 167, col: 56, offset: 3665},
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
	pos: position{line: 171, col: 1, offset: 3701},
	expr: &actionExpr{
	pos: position{line: 171, col: 15, offset: 3715},
	run: (*parser).callonFLAGS_RULE1,
	expr: &seqExpr{
	pos: position{line: 171, col: 15, offset: 3715},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 171, col: 15, offset: 3715},
	name: "WS_MAND",
},
&labeledExpr{
	pos: position{line: 171, col: 23, offset: 3723},
	label: "i",
	expr: &ruleRefExpr{
	pos: position{line: 171, col: 25, offset: 3725},
	name: "IGNORE_FLAG",
},
},
&labeledExpr{
	pos: position{line: 171, col: 37, offset: 3737},
	label: "is",
	expr: &zeroOrMoreExpr{
	pos: position{line: 171, col: 40, offset: 3740},
	expr: &seqExpr{
	pos: position{line: 171, col: 41, offset: 3741},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 171, col: 41, offset: 3741},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 171, col: 44, offset: 3744},
	name: "LS",
},
&ruleRefExpr{
	pos: position{line: 171, col: 47, offset: 3747},
	name: "WS",
},
&ruleRefExpr{
	pos: position{line: 171, col: 50, offset: 3750},
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
	pos: position{line: 175, col: 1, offset: 3793},
	expr: &actionExpr{
	pos: position{line: 175, col: 16, offset: 3808},
	run: (*parser).callonIGNORE_FLAG1,
	expr: &litMatcher{
	pos: position{line: 175, col: 16, offset: 3808},
	val: "ignore-errors",
	ignoreCase: false,
},
},
},
{
	name: "VARIABLE",
	pos: position{line: 179, col: 1, offset: 3855},
	expr: &actionExpr{
	pos: position{line: 179, col: 13, offset: 3867},
	run: (*parser).callonVARIABLE1,
	expr: &seqExpr{
	pos: position{line: 179, col: 13, offset: 3867},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 179, col: 13, offset: 3867},
	val: "$",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 179, col: 17, offset: 3871},
	label: "v",
	expr: &ruleRefExpr{
	pos: position{line: 179, col: 20, offset: 3874},
	name: "IDENT",
},
},
	},
},
},
},
{
	name: "IDENT",
	pos: position{line: 183, col: 1, offset: 3909},
	expr: &actionExpr{
	pos: position{line: 183, col: 10, offset: 3918},
	run: (*parser).callonIDENT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 183, col: 10, offset: 3918},
	expr: &charClassMatcher{
	pos: position{line: 183, col: 10, offset: 3918},
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
	pos: position{line: 187, col: 1, offset: 3964},
	expr: &actionExpr{
	pos: position{line: 187, col: 19, offset: 3982},
	run: (*parser).callonIDENT_WITH_DOT1,
	expr: &oneOrMoreExpr{
	pos: position{line: 187, col: 19, offset: 3982},
	expr: &charClassMatcher{
	pos: position{line: 187, col: 19, offset: 3982},
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
	pos: position{line: 191, col: 1, offset: 4029},
	expr: &actionExpr{
	pos: position{line: 191, col: 9, offset: 4037},
	run: (*parser).callonNull1,
	expr: &litMatcher{
	pos: position{line: 191, col: 9, offset: 4037},
	val: "null",
	ignoreCase: false,
},
},
},
{
	name: "Boolean",
	pos: position{line: 195, col: 1, offset: 4067},
	expr: &actionExpr{
	pos: position{line: 195, col: 12, offset: 4078},
	run: (*parser).callonBoolean1,
	expr: &choiceExpr{
	pos: position{line: 195, col: 13, offset: 4079},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 195, col: 13, offset: 4079},
	val: "true",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 195, col: 22, offset: 4088},
	val: "false",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "String",
	pos: position{line: 199, col: 1, offset: 4129},
	expr: &actionExpr{
	pos: position{line: 199, col: 11, offset: 4139},
	run: (*parser).callonString1,
	expr: &seqExpr{
	pos: position{line: 199, col: 11, offset: 4139},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 199, col: 11, offset: 4139},
	val: "\"",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 199, col: 15, offset: 4143},
	expr: &seqExpr{
	pos: position{line: 199, col: 17, offset: 4145},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 199, col: 17, offset: 4145},
	expr: &litMatcher{
	pos: position{line: 199, col: 18, offset: 4146},
	val: "\"",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 199, col: 22, offset: 4150,
},
	},
},
},
&litMatcher{
	pos: position{line: 199, col: 27, offset: 4155},
	val: "\"",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Float",
	pos: position{line: 203, col: 1, offset: 4190},
	expr: &actionExpr{
	pos: position{line: 203, col: 10, offset: 4199},
	run: (*parser).callonFloat1,
	expr: &seqExpr{
	pos: position{line: 203, col: 10, offset: 4199},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 203, col: 10, offset: 4199},
	expr: &choiceExpr{
	pos: position{line: 203, col: 11, offset: 4200},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 203, col: 11, offset: 4200},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 203, col: 17, offset: 4206},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 203, col: 23, offset: 4212},
	name: "Natural",
},
&litMatcher{
	pos: position{line: 203, col: 31, offset: 4220},
	val: ".",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 203, col: 35, offset: 4224},
	name: "Natural",
},
	},
},
},
},
{
	name: "Integer",
	pos: position{line: 207, col: 1, offset: 4262},
	expr: &actionExpr{
	pos: position{line: 207, col: 12, offset: 4273},
	run: (*parser).callonInteger1,
	expr: &seqExpr{
	pos: position{line: 207, col: 12, offset: 4273},
	exprs: []interface{}{
&zeroOrOneExpr{
	pos: position{line: 207, col: 12, offset: 4273},
	expr: &choiceExpr{
	pos: position{line: 207, col: 13, offset: 4274},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 207, col: 13, offset: 4274},
	val: "+",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 207, col: 19, offset: 4280},
	val: "-",
	ignoreCase: false,
},
	},
},
},
&ruleRefExpr{
	pos: position{line: 207, col: 25, offset: 4286},
	name: "Natural",
},
	},
},
},
},
{
	name: "Natural",
	pos: position{line: 211, col: 1, offset: 4326},
	expr: &choiceExpr{
	pos: position{line: 211, col: 11, offset: 4338},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 211, col: 11, offset: 4338},
	val: "0",
	ignoreCase: false,
},
&seqExpr{
	pos: position{line: 211, col: 17, offset: 4344},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 211, col: 17, offset: 4344},
	name: "NonZeroDecimalDigit",
},
&zeroOrMoreExpr{
	pos: position{line: 211, col: 37, offset: 4364},
	expr: &ruleRefExpr{
	pos: position{line: 211, col: 37, offset: 4364},
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
	pos: position{line: 213, col: 1, offset: 4379},
	expr: &charClassMatcher{
	pos: position{line: 213, col: 16, offset: 4396},
	val: "[0-9]",
	ranges: []rune{'0','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "NonZeroDecimalDigit",
	pos: position{line: 214, col: 1, offset: 4402},
	expr: &charClassMatcher{
	pos: position{line: 214, col: 23, offset: 4426},
	val: "[1-9]",
	ranges: []rune{'1','9',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "SPACE",
	pos: position{line: 216, col: 1, offset: 4433},
	expr: &charClassMatcher{
	pos: position{line: 216, col: 10, offset: 4442},
	val: "[ \\t]",
	chars: []rune{' ','\t',},
	ignoreCase: false,
	inverted: false,
},
},
{
	name: "WS_MAND",
	displayName: "\"mandatory-whitespace\"",
	pos: position{line: 217, col: 1, offset: 4448},
	expr: &oneOrMoreExpr{
	pos: position{line: 217, col: 35, offset: 4482},
	expr: &choiceExpr{
	pos: position{line: 217, col: 36, offset: 4483},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 217, col: 36, offset: 4483},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 217, col: 44, offset: 4491},
	name: "COMMENT",
},
&ruleRefExpr{
	pos: position{line: 217, col: 54, offset: 4501},
	name: "NL",
},
	},
},
},
},
{
	name: "WS",
	displayName: "\"whitespace\"",
	pos: position{line: 218, col: 1, offset: 4506},
	expr: &zeroOrMoreExpr{
	pos: position{line: 218, col: 20, offset: 4525},
	expr: &choiceExpr{
	pos: position{line: 218, col: 21, offset: 4526},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 218, col: 21, offset: 4526},
	name: "SPACE",
},
&ruleRefExpr{
	pos: position{line: 218, col: 29, offset: 4534},
	name: "COMMENT",
},
	},
},
},
},
{
	name: "LS",
	displayName: "\"line-separator\"",
	pos: position{line: 219, col: 1, offset: 4544},
	expr: &choiceExpr{
	pos: position{line: 219, col: 25, offset: 4568},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 219, col: 25, offset: 4568},
	name: "NL",
},
&litMatcher{
	pos: position{line: 219, col: 30, offset: 4573},
	val: ",",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 219, col: 36, offset: 4579},
	name: "COMMENT",
},
	},
},
},
{
	name: "BS",
	displayName: "\"block-separator\"",
	pos: position{line: 220, col: 1, offset: 4588},
	expr: &oneOrMoreExpr{
	pos: position{line: 220, col: 25, offset: 4612},
	expr: &seqExpr{
	pos: position{line: 220, col: 26, offset: 4613},
	exprs: []interface{}{
&ruleRefExpr{
	pos: position{line: 220, col: 26, offset: 4613},
	name: "WS",
},
&choiceExpr{
	pos: position{line: 220, col: 30, offset: 4617},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 220, col: 30, offset: 4617},
	name: "NL",
},
&ruleRefExpr{
	pos: position{line: 220, col: 35, offset: 4622},
	name: "COMMENT",
},
	},
},
&ruleRefExpr{
	pos: position{line: 220, col: 44, offset: 4631},
	name: "WS",
},
	},
},
},
},
{
	name: "NL",
	displayName: "\"new-line\"",
	pos: position{line: 221, col: 1, offset: 4636},
	expr: &litMatcher{
	pos: position{line: 221, col: 18, offset: 4653},
	val: "\n",
	ignoreCase: false,
},
},
{
	name: "COMMENT",
	pos: position{line: 223, col: 1, offset: 4659},
	expr: &seqExpr{
	pos: position{line: 223, col: 12, offset: 4670},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 223, col: 12, offset: 4670},
	val: "//",
	ignoreCase: false,
},
&zeroOrMoreExpr{
	pos: position{line: 223, col: 17, offset: 4675},
	expr: &seqExpr{
	pos: position{line: 223, col: 19, offset: 4677},
	exprs: []interface{}{
&notExpr{
	pos: position{line: 223, col: 19, offset: 4677},
	expr: &litMatcher{
	pos: position{line: 223, col: 20, offset: 4678},
	val: "\n",
	ignoreCase: false,
},
},
&anyMatcher{
	line: 223, col: 25, offset: 4683,
},
	},
},
},
&choiceExpr{
	pos: position{line: 223, col: 31, offset: 4689},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 223, col: 31, offset: 4689},
	val: "\n",
	ignoreCase: false,
},
&ruleRefExpr{
	pos: position{line: 223, col: 38, offset: 4696},
	name: "EOF",
},
	},
},
	},
},
},
{
	name: "EOF",
	pos: position{line: 225, col: 1, offset: 4702},
	expr: &notExpr{
	pos: position{line: 225, col: 8, offset: 4709},
	expr: &anyMatcher{
	line: 225, col: 9, offset: 4710,
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

func (c *current) onUSE2(r, v interface{}) (interface{}, error) {
	return newUse(r, v)
}

func (p *parser) callonUSE2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE2(stack["r"], stack["v"])
}

func (c *current) onUSE15(r, v interface{}) (interface{}, error) {
	return Use{}, errors.New("use of equal in use clause is deprecated")
}

func (p *parser) callonUSE15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUSE15(stack["r"], stack["v"])
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
	return nil, errors.New("cache-control use action is deprecated")
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

