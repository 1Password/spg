## creates Go code defining a string slice with list of words

# Usage
# awk -f goify_words.awk < /path/to/AgileWords.txt > agilewords.go

BEGIN {
    print "// Code generated from testdata/agsyllables.txt; DO NOT EDIT.\n";
    print "package spg\n";
    print "// AgileWords is the list of words used by the 1Password strong password generator";
    print "var AgileWords = []string{";
    OFS = ""
     }

    { print "    \"", $0, "\"," }

END { print "}" }
