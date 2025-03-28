package repository

import (
    "fmt"
    "log"
    "github.com/tarantool/go-tarantool/v2"
)

type SchemaManager struct {
    conn *tarantool.Connection
}

func (sm *SchemaManager) Init() error {
    if err := sm.createPollSpace(); err != nil {
        return err
    }
    if err := sm.createVoteSpace(); err != nil {
        return err
    }
    return nil
}

func (sm *SchemaManager) createPollSpace() error {
    // Формат полей через map (без использования SpaceField)
    _, err := sm.conn.Do(
        tarantool.NewCallRequest("box.schema.space.create").
            Args([]interface{}{
                "polls",
                map[string]interface{}{
                    "if_not_exists": true,
                    "format": []map[string]string{
                        {"name": "id", "type": "string"},
                        {"name": "question", "type": "string"},
                        {"name": "options", "type": "string"},
                        {"name": "created_by", "type": "string"},
                        {"name": "channel_id", "type": "string"},
                        {"name": "active", "type": "boolean"},
                    },
                },
            }),
    ).Get()
    if err != nil {
        return fmt.Errorf("create polls space failed: %v", err)
    }
    log.Println("Space 'polls' created")

    // Создание индексов
    _, err = sm.conn.Do(
        tarantool.NewCallRequest("box.space.polls:create_index").
            Args([]interface{}{
                "primary",
                map[string]interface{}{
                    "type":   "hash",
                    "parts":  []interface{}{"id"},
                    "unique": true,
                },
            }),
    ).Get()
    return err
}

func (sm *SchemaManager) createVoteSpace() error {
    // Аналогично для пространства votes
    _, err := sm.conn.Do(
        tarantool.NewCallRequest("box.schema.space.create").
            Args([]interface{}{
                "votes",
                map[string]interface{}{
                    "if_not_exists": true,
                    "format": []map[string]string{
                        {"name": "poll_id", "type": "string"},
                        {"name": "user_id", "type": "string"},
                        {"name": "option_idx", "type": "unsigned"},
                    },
                },
            }),
    ).Get()
    if err != nil {
        return fmt.Errorf("create votes space failed: %v", err)
    }
    log.Println("Space 'votes' created")

    // Создание индексов
    _, err = sm.conn.Do(
        tarantool.NewCallRequest("box.space.votes:create_index").
            Args([]interface{}{
                "primary",
                map[string]interface{}{
                    "type":   "hash",
                    "parts":  []interface{}{"poll_id", "user_id"},
                    "unique": true,
                },
            }),
    ).Get()
    return err
}