package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Studs []Stud

//addStud delStud edtStud pAdd pDel pSer (ser)

func serStud(mod Stud) int {
	// ret: 0 for all paired, 1 for username paired, 2 for nothing paired
	nMod := Stud{}
	err = db.QueryRow("SELECT * FROM stud_table WHERE sid = ?", mod.SID).Scan(&nMod.SID, &nMod.Name, &nMod.Email, &nMod.Tel)
	if err != nil {
		log.Println(err)
		return 1
	}
	if nMod.SID != "" && (nMod.SID == mod.SID) {
		return 0
	}
	return 1
}

func addStud(qStr []byte) int {
	mod := Stud{}
	err := json.Unmarshal(qStr, &mod)
	if err != nil {
		log.Println(err)
		return 1
	}
	ret := serStud(mod)
	if ret != 1 {
		return 4
	}

	stmt, err := db.Prepare("INSERT stud_table SET sid=?, name=?, email=?, tel=?")
	if err != nil {
		log.Println(err)
		return 2
	}
	_, err = stmt.Exec(mod.SID, mod.Name, mod.Email, mod.Tel)
	//log.Println(res)
	if err != nil {
		log.Println(err)
		return 3
	}

	return 0
}

func delStud(qStr []byte) int {
	mod := Stud{}
	fmt.Println("123", qStr, "456")
	err := json.Unmarshal(qStr, &mod)
	if err != nil {
		log.Println(err)
		return 1
	}

	ret := serStud(mod)
	if ret != 0 {
		return 4
	}

	stmt, err := db.Prepare("DELETE FROM stud_table WHERE sid=?")
	log.Println(stmt, err)
	if err != nil {
		log.Println(err)
		return 2
	}
	res, err := stmt.Exec(mod.SID)
	log.Println(res, err)
	if err != nil {
		log.Println(err)
		return 3
	}
	return 0
}

func edtStud(qStr []byte) int {
	mod := StudIn{}
	err := json.Unmarshal(qStr, &mod)
	nMod := Stud{}
	nMod.SID = mod.OriSID
	if err != nil {
		log.Println(err)
		return 1
	}
	fmt.Println(nMod.SID)
	ret := serStud(nMod)
	if ret != 0 {
		return 4
	}

	stmt, err := db.Prepare("UPDATE stud_table SET sid=?, name=?, email=?, tel=? WHERE sid = ?")
	log.Println(stmt)
	if err != nil {
		log.Println(err)
		return 2
	}
	res, err := stmt.Exec(mod.SID, mod.Name, mod.Email, mod.Tel, mod.OriSID)
	log.Println(res)
	if err != nil {
		log.Println(err)
		return 3
	}

	return 0
}

func pAddStud(qStr []byte) int {
	var mods []Stud
	err := json.Unmarshal(qStr, &mods)
	if err != nil {
		return -1
	}
	for ind, mod := range mods {
		res := serStud(mod)
		if res != 1 {
			return ind + 1
		}
		log.Println(ind, mod.SID, mod.Name, mod.Email, mod.Tel)
		stmt, err := db.Prepare("INSERT stud_table SET sid=?, name=?, email=?, tel=?")
		if err != nil {
			return ind + 1
		}
		_, err = stmt.Exec(mod.SID, mod.Name, mod.Email, mod.Tel)
		if err != nil {
			return ind + 1
		}
	}
	return 0
}

func pDelStud(qStr []byte) int {
	var mods []Stud
	err := json.Unmarshal(qStr, &mods)
	if err != nil {
		return -1
	}
	for ind, mod := range mods {
		res := serStud(mod)
		if res != 0 {
			return ind + 1
		}
		log.Println(ind, mod.SID, mod.Name, mod.Email, mod.Tel)
		stmt, err := db.Prepare("DELETE FROM stud_table WHERE sid=?, name=?, tel=?, email=?")
		if err != nil {
			return ind + 1
		}
		_, err = stmt.Exec(mod.SID, mod.Name, mod.Email, mod.Tel)
		if err != nil {
			return ind + 1
		}
	}
	return 0
}

func pSerStud(qStr []byte) (int, []byte) {
	var tmpstuds []Stud
	Studs = tmpstuds
	mod := Stud{}
	err := json.Unmarshal(qStr, &mod)
	rows, err := db.Query("SELECT * FROM stud_table WHERE sid LIKE ? AND name LIKE ? AND email LIKE ?", "%"+mod.SID+"%", "%"+mod.Name+"%", "%"+mod.Email+"%") //tel for char
	var xxxx []byte
	if err != nil {
		log.Println(err)
		return 1, xxxx
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&mod.SID, &mod.Name, &mod.Email, &mod.Tel)
		Studs = append(Studs, mod)
	}
	retData, err := json.Marshal(Studs)
	if err != nil {
		log.Println(err)
		return 1, retData
	}
	return 0, retData
}

func retAll() (int, []byte) {
	var tmpstuds []Stud
	Studs = tmpstuds
	rows, err := db.Query("SELECT * FROM stud_table")
	var xxxx []byte
	if err != nil {
		return 1, xxxx
	}
	defer rows.Close()
	//it := 0
	for rows.Next() {
		mod := Stud{}
		rows.Scan(&mod.SID, &mod.Name, &mod.Email, &mod.Tel)
		Studs = append(Studs, mod)
	}
	retData, err := json.Marshal(Studs)
	if err != nil {
		return 1, retData
	}
	return 0, retData
}

func sdb(op string, qStr []byte) (int, []byte) {
	studb, err := sql.Open("mysql", "root:bjwdttz@tcp(127.0.0.1:3306)/tst?charset=utf8")
	retCode := 0
	var retData []byte
	if err != nil {
		return 10, retData
	}
	switch op {
	case "add":
		retCode = addStud(qStr)
	case "del":
		retCode = delStud(qStr)
	case "edt":
		retCode = edtStud(qStr)
	case "pAdd":
		retCode = pAddStud(qStr)
	case "pDel":
		retCode = pDelStud(qStr)
	case "pSer":
		retCode, retData = pSerStud(qStr)
	case "ret":
		retCode, retData = retAll()
	}
	log.Println(retCode, retData)
	defer studb.Close()
	return retCode, retData
}
