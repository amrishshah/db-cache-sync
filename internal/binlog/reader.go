package binlog

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/amrishkshah/db-cache-sync/internal/cache"
	"github.com/amrishkshah/db-cache-sync/internal/config"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type Handler struct {
	cache *cache.RedisClient
	cfg   *config.Config
}

func (h *Handler) OnRow(e *canal.RowsEvent) error {
	table := e.Table.Name
	action := e.Action

	if !h.cfg.TableSet[table] {
		return nil
	}

	for idx, row := range e.Rows {
		if action == canal.UpdateAction {
			idx = idx/2*2 + 1 // updated row is second
		}

		log.Println(e.Table.GetPKColumn(0).Name)

		id := fmt.Sprintf("%v", row[0]) // assume first column is 'id'

		key := fmt.Sprintf("%s:%s", table, id)

		if action == canal.DeleteAction {
			err := h.cache.Del(key)
			if err != nil {
				log.Printf("[ERROR] Failed to delete key %s: %v", key, err)
			} else {
				log.Printf("[DEL] key=%s", key)
			}
			continue
		}

		// Build JSON object
		rowMap := make(map[string]interface{})
		for i, column := range e.Table.Columns {
			rowMap[column.Name] = row[i]
		}

		jsonBytes, _ := json.Marshal(rowMap)

		err := h.cache.Set(key, string(jsonBytes))
		if err != nil {
			log.Printf("[ERROR] Failed to set key %s: %v", key, err)
		} else {
			log.Printf("[SET] key=%s value=%s", key, string(jsonBytes))
		}
	}

	return nil
}

func (h *Handler) String() string { return "Handler" }

func (h *Handler) OnRowsQueryEvent(rowsQueryEvent *replication.RowsQueryEvent) error {
	return nil
}

func (h *Handler) OnRotate(header *replication.EventHeader, rotateEvent *replication.RotateEvent) error {
	return nil
}

func (h *Handler) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	return nil
}

// ⚡ Corrected OnDDL
func (h *Handler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	return nil
}

func (h *Handler) OnXID(header *replication.EventHeader, nextPos mysql.Position) error { return nil }

func (h *Handler) OnGTID(header *replication.EventHeader, gtid mysql.BinlogGTIDEvent) error {
	return nil
}

func (h *Handler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return nil
}

func StartBinlogReader(cfg config.Config, redisClient *cache.RedisClient) error {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = cfg.MySQLAddr
	canalCfg.User = cfg.MySQLUser
	canalCfg.Password = cfg.MySQLPassword
	canalCfg.Dump.TableDB = ""
	canalCfg.Dump.Tables = nil
	canalCfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(canalCfg)
	if err != nil {
		return err
	}

	handler := &Handler{cache: redisClient,
		cfg: &cfg}

	c.SetEventHandler(handler)

	return c.Run()
}
