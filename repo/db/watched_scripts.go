package db

import (
	"database/sql"
	"encoding/hex"
	"sync"
)

// WatchScriptsDB type definition.
// Sets a pointer to SQL database and syncs reader/writer mutex-based lock.
type WatchedScriptsDB struct {
	db   *sql.DB
	lock sync.RWMutex
}

// WatchdScriptsDB Put method insert and replace operations based on watched script public keys.
func (w *WatchedScriptsDB) Put(scriptPubKey []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	tx, _ := w.db.Begin()
	stmt, err := tx.Prepare("insert or replace into watchedscripts(scriptPubKey) values(?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(hex.EncodeToString(scriptPubKey))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// WatchedScriptsDB GetAll() method Returns all watched script public keys.
func (w *WatchedScriptsDB) GetAll() ([][]byte, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	var ret [][]byte
	stm := "select scriptPubKey from watchedscripts"
	rows, err := w.db.Query(stm)
	if err != nil {
		return ret, err
	}
	defer rows.Close()
	for rows.Next() {
		var scriptHex string
		if err := rows.Scan(&scriptHex); err != nil {
			continue
		}
		scriptPubKey, err := hex.DecodeString(scriptHex)
		if err != nil {
			continue
		}
		ret = append(ret, scriptPubKey)
	}
	return ret, nil
}

// WatchedScriptsDB Delete method deletes a watched script based on its public key.
func (w *WatchedScriptsDB) Delete(scriptPubKey []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	_, err := w.db.Exec("delete from watchedscripts where scriptPubKey=?", hex.EncodeToString(scriptPubKey))
	if err != nil {
		return err
	}
	return nil
}
