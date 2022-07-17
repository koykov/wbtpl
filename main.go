package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ps   = string(os.PathSeparator)
	rtpl []byte
	fdb  = flag.String("db", fmt.Sprintf("local%sorg.csv", ps), "path to local database")
	ftpl = flag.String("tpl", fmt.Sprintf("local%stpl.html", ps), "path to template file")
	fout = flag.String("out", fmt.Sprintf(".%sout", ps), "path to output directory")
	fday = flag.Uint("days", 1, "How many days need to print")

	bOrgName   = []byte("{ORG_NAME}")
	bOrgIDNO   = []byte("{ORG_IDNO}")
	bOrgAddr   = []byte("{ORG_ADDR}")
	bOrgPhone  = []byte("{ORG_PHONE}")
	bSeria     = []byte("{SERIA}")
	bNumber    = []byte("{NUMBER}")
	bDateDay   = []byte("{DATE_DAY}")
	bDateMonth = []byte("{DATE_MONTH}")
	bDateYear  = []byte("{DATE_YEAR}")
	bCarModel  = []byte("{CAR_MODEL}")
	bCarNumber = []byte("{CAR_NUMBER}")
	bDriver    = []byte("{DRIVER_NAME}")
)

func init() {
	flag.Parse()
	if !fileExists(*fdb) {
		log.Fatalf("local database not exists %s\n", *fdb)
	}
	if !fileExists(*ftpl) {
		log.Fatalf("template file not exists %s\n", *ftpl)
	}
	var err error
	if rtpl, err = os.ReadFile(*ftpl); err != nil {
		log.Fatalf("failed to read template file: %s", err.Error())
	}
	if err = dirProbe(*fout); err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := os.Open(*fdb)
	if err != nil {
		log.Fatalf("coudn't open local database '%s': %s", *fdb, err)
	}
	defer func() { _ = f.Close() }()

	csvr := csv.NewReader(f)
	csvr.Comma = ';'
	rows, err := csvr.ReadAll()
	if err != nil {
		log.Fatalf("couldn't parse database file '%s': %s", *fdb, err)
	}
	if len(rows) <= 1 {
		log.Fatal("local database is empty")
	}
	rows = rows[1:]
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		for i := uint(0); i < *fday; i++ {
			date := time.Now().AddDate(0, 0, int(i))
			if err = generate(row[0], row[1], row[2], row[3], row[4], date); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func generate(org, slug, idno, addr, phone string, date time.Time) error {
	odir := filepath.Dir(*fdb)
	odb := fmt.Sprintf("%s%s%s.csv", odir, ps, slug)
	f, err := os.Open(odb)
	if err != nil {
		log.Fatalf("coudn't open company database '%s': %s", *fdb, err)
	}
	defer func() { _ = f.Close() }()

	outdir := fmt.Sprintf("%s%s%s%s%s", *fout, ps, slug, ps, date.Format("2006-01-02"))
	if err = os.MkdirAll(outdir, 0755); err != nil {
		return err
	}

	csvr := csv.NewReader(f)
	csvr.Comma = ';'
	rows, err := csvr.ReadAll()
	if err != nil {
		log.Fatalf("couldn't parse database file '%s': %s", *fdb, err)
	}
	if len(rows) <= 1 {
		log.Fatal("company database is empty")
	}
	rows = rows[1:]
	var buf []byte
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		dn := strings.ReplaceAll(row[4], " ", "_")
		out := fmt.Sprintf("%s%s%s.html", outdir, ps, dn)
		log.Printf("processing '%s/%s' to '%s", slug, row[4], out)
		buf = append(buf[:0], rtpl...)
		buf = bytes.ReplaceAll(buf, bOrgName, []byte(org))
		buf = bytes.ReplaceAll(buf, bOrgIDNO, []byte(idno))
		buf = bytes.ReplaceAll(buf, bOrgAddr, []byte(addr))
		buf = bytes.ReplaceAll(buf, bOrgPhone, []byte(phone))
		buf = bytes.ReplaceAll(buf, bSeria, []byte(row[0]))
		buf = bytes.ReplaceAll(buf, bNumber, []byte(row[1]))
		buf = bytes.ReplaceAll(buf, bDateDay, []byte(date.Format("2")))
		buf = bytes.ReplaceAll(buf, bDateMonth, []byte(date.Format("1")))
		buf = bytes.ReplaceAll(buf, bDateYear, []byte(date.Format("2006")))
		buf = bytes.ReplaceAll(buf, bCarModel, []byte(row[2]))
		buf = bytes.ReplaceAll(buf, bCarNumber, []byte(row[3]))
		buf = bytes.ReplaceAll(buf, bDriver, []byte(row[4]))
		if err = os.WriteFile(out, buf, 0777); err != nil {
			log.Fatalf("failed to write file '%s: %s", out, err.Error())
		}
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func dirProbe(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}
