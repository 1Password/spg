## creates Go code defining a string slice with list of syllables

# Usage
# awk -f goify_syllables.awk < /path/to/AgileSyllables.txt > agilesyllables.go

BEGIN {
    print "package spg\n";
    print "// This file is automatically generated. Do not edit\n";
    print "// AgileSyllables is the list of syllables used by the 1Password strong password generator";
    print "var AgileSyllables = []string{";
    OFS = ""
     }

    { print "    \"", $0, "\"," }

END { print "}" }