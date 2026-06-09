# alegra-cli demo — rendered with VHS (https://github.com/charmbracelet/vhs).
# Regenerate locally:  make build && vhs .demo/demo.tape
# Uses only credential-free commands (version, --help, --dry-run).

Output assets/demo.gif

Set Shell bash
Set FontSize 16
Set Width 1180
Set Height 660
Set Padding 18
Set Theme "Dracula"
Set TypingSpeed 35ms

# Setup (hidden): use the locally-built binary; a demo email so dry-run renders.
Hide
Type "export PATH=$PWD/bin:$PATH; export ALEGRA_EMAIL=demo@acme.co; clear"
Enter
Show

Sleep 700ms
Type "alegra version"
Sleep 400ms
Enter
Sleep 1800ms

Type "alegra invoices list --status open --since this-month --dry-run"
Sleep 400ms
Enter
Sleep 2600ms

Type "alegra contacts create --set name='Acme S.A.S' --set kindOfPerson=LEGAL_ENTITY --dry-run"
Sleep 400ms
Enter
Sleep 2600ms

Type "alegra invoices --help"
Sleep 400ms
Enter
Sleep 3200ms
