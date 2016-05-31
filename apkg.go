package apkg

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	//"strings"
)

type akgInfo struct {
	path string
	db *sql.DB
}

func Create(path string) (*akgInfo){
	err := os.Mkdir(path, 0777);
	if err != nil {
		log.Fatal(err);
	}

	db, err := sql.Open("sqlite3", path + "/collection.anki2");
	if err != nil {
		log.Fatal(err);
	}

	var pkg = new(akgInfo);

	tx, err := db.Begin();
	if err != nil {
		log.Fatal(err);
	}

	//based on: http://decks.wikia.com/wiki/Anki_APKG_format_documentation
	_, err = db.Exec(`
			CREATE TABLE col (
				id     INTEGER PRIMARY KEY,
				crt    INTEGER NOT NULL,
				mod    INTEGER NOT NULL,
				scm    INTEGER NOT NULL,
				ver    INTEGER NOT NULL,
				dty    INTEGER NOT NULL,
				usn    INTEGER NOT NULL,
				ls     INTEGER NOT NULL,
				conf   TEXT    NOT NULL,
				models TEXT    NOT NULL,
				decks  TEXT    NOT NULL,
				dconf  TEXT    NOT NULL,
				tags   TEXT    NOT NULL
			);
	`);

	_, err = db.Exec(`
		CREATE TABLE notes (
			id    INTEGER PRIMARY KEY,
			guid  TEXT    NOT NULL,
			mid   INTEGER NOT NULL,
			mod   INTEGER NOT NULL,
			usn   INTEGER NOT NULL,
			tags  TEXT    NOT NULL,
			flds  TEXT    NOT NULL,
			sfld  INTEGER NOT NULL,
			csum  INTEGER NOT NULL,
			flags INTEGER NOT NULL,
			data  TEXT    NOT NULL
		);
	`);

	_, err = db.Exec(`
		CREATE TABLE cards (
			id     INTEGER PRIMARY KEY,
			nid    INTEGER NOT NULL,
			did    INTEGER NOT NULL,
			ord    INTEGER NOT NULL,
			mod    INTEGER NOT NULL,
			usn    INTEGER NOT NULL,
			type   INTEGER NOT NULL,
			queue  INTEGER NOT NULL,
			due    INTEGER NOT NULL,
			ivl    INTEGER NOT NULL,
			factor INTEGER NOT NULL,
			reps   INTEGER NOT NULL,
			lapses INTEGER NOT NULL,
			left   INTEGER NOT NULL,
			odue   INTEGER NOT NULL,
			odid   INTEGER NOT NULL,
			flags  INTEGER NOT NULL,
			data   TEXT    NOT NULL
		);
	`);

	_, err = db.Exec(`
		CREATE TABLE revlog (
			id      INTEGER PRIMARY KEY,
			cid     INTEGER NOT NULL,
			usn     INTEGER NOT NULL,
			ease    INTEGER NOT NULL,
			ivl     INTEGER NOT NULL,
			lastIvl INTEGER NOT NULL,
			factor  INTEGER NOT NULL,
			time    INTEGER NOT NULL,
			type    INTEGER NOT NULL
		);
	`);

	_, err = db.Exec(`
		CREATE TABLE graves (
			usn  INTEGER NOT NULL,
			oid  INTEGER NOT NULL,
			type INTEGER NOT NULL
		);
	`);

	_, err = db.Exec(`
		ANALYZE sqlite_master;
		INSERT INTO "sqlite_stat1" VALUES ('col', NULL, '1');
		CREATE INDEX ix_notes_usn ON notes (usn);
		CREATE INDEX ix_cards_usn ON cards (usn);
		CREATE INDEX ix_revlog_usn ON revlog (usn);
		CREATE INDEX ix_cards_nid ON cards (nid);
		CREATE INDEX ix_cards_sched ON cards (did, queue, due);
		CREATE INDEX ix_revlog_cid ON revlog (cid);
		CREATE INDEX ix_notes_csum ON notes (csum);
	`);

	tx.Commit();

	pkg.db = db;
	pkg.path = path;
	return pkg;
}

func (pkg *akgInfo) Close() {
	pkg.db.Close();
}