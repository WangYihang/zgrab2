package zgrab2

import _ "embed"

//go:generate sh -c "printf %s $(git rev-parse HEAD) > .commit"
//go:embed .commit
var Commit string
