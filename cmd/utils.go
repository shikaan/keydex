package cmd

import (
	"bufio"
	"os"
	"strings"
)

// If zero value reference is passed, reads from stdin to get the value
func ReadReferenceFromStdin(maybeReference string) (string, error) {
	if maybeReference != "" {
		return maybeReference, nil
	}

	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}


// cases:
// no database (but environment) no reference -> only database
// no database (no environment) no reference -> throw
// no database (no environment) reference -> use ref as database?
// no database (but environment) reference -> use last arg as reference
// database (and environment) no reference -> database from args


// Cases: 
//   DATABASE=lol kpcli copy -> gets ref from stdin (blank ref)
//   kpcli copy database -> gets ref from stdin (blank ref)
//   DATABASE=lol kpcli copy /ref -> OK (db: lol, ref: /ref)
//   kpcli copy database /ref -> OK (db: database, ref: /ref)
//   DATABASE=lol kpcli copy database /ref -> (db database, /ref)
//   kpcli copy /ref -> uses /ref as db and fails

//   DATABASE=lol kpcli open -> opens list (blank ref)
//   kpcli copy database -> gets ref from stdin (blank ref)
//   DATABASE=lol kpcli open /ref -> OK (db: lol, ref: /ref)
//   kpcli open database /ref -> OK (db: database, ref: /ref)
//   DATABASE=lol kpcli open database /ref -> (db database, /ref)
//   kpcli open /ref -> uses /ref as db and then fails
func ReadDatabaseArguments(args []string) (string, string) {
  var reference, database string;

  if len(args) == 0 {
    reference = ""
    database = os.Getenv("KPCLI_DATABASE")
  }

  if len(args) == 1 {
    reference = ""
    database = args[0]
  }

  if len(args) == 2 {
    database = args[0]
    reference = args[1]
  }

	return database, reference
}
