# alegra-cli demo — rendered with VHS (https://github.com/charmbracelet/vhs).
# Shows REAL output against an authenticated demo account.
# Regenerate locally (needs `alegra auth login` first):
#   make build && vhs .demo/demo.tape

Output assets/demo.gif

Set Shell bash
Set FontSize 15
Set Width 1200
Set Height 760
Set Padding 18
Set Theme "Dracula"
Set TypingSpeed 30ms

# Setup (hidden): use the locally-built binary.
Hide
Type "export PATH=$PWD/bin:$PATH; clear"
Enter
Show

Sleep 600ms
Type "alegra version"
Sleep 300ms
Enter
Sleep 1400ms

Type "alegra doctor"
Sleep 300ms
Enter
Sleep 3000ms

Type "alegra contacts create -f .demo/contact.json -o json"
Sleep 300ms
Enter
Sleep 2400ms

Type "alegra items create -f .demo/item.json -o json"
Sleep 300ms
Enter
Sleep 2200ms

Type "alegra contacts list -o table --columns id,name,type"
Sleep 300ms
Enter
Sleep 2400ms

Type "alegra taxes list -o table"
Sleep 300ms
Enter
Sleep 3000ms
