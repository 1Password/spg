## creates Go code defining a string slice with list of syllables

# Usage
# awk -f goify_syllables.awk < /path/to/AgileSyllables.txt > agilesyllables.go

BEGIN {
    print "// Code generated from testdata/agwordlist.txt; DO NOT EDIT.\n";
    print "package spg\n";
    print "// AgileSyllables is the list of syllables used by the 1Password strong password generator";
    print "var AgileSyllables = []string{";
    OFS = ""
     }

    { print "    \"", $0, "\"," }

END { print "}" }
