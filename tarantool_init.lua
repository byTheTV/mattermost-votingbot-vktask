-- Инициализация конфигурации Tarantool
box.cfg{
    listen = 3301,
    wal_mode = 'none',
    memtx_memory = 128 * 1024 * 1024 -- 128MB
}

-- polls
box.schema.space.create('polls', { if_not_exists = true })
box.space.polls:format({
    { name = 'id',         type = 'string' },
    { name = 'question',   type = 'string' },
    { name = 'options',    type = 'string' },
    { name = 'created_by', type = 'string' },
    { name = 'channel_id', type = 'string' },
    { name = 'active',     type = 'boolean' },
})
box.space.polls:create_index('primary', {
    parts = {'id'},
    unique = true,
    if_not_exists = true
})

box.space.polls:create_index('channel', {
    parts = {'channel_id'},
    if_not_exists = true,
})

-- votes
box.schema.space.create('votes', { if_not_exists = true })
box.space.votes:format({
    { name = 'poll_id',  type = 'string' },
    { name = 'user_id',  type = 'string' },
    { name = 'option_idx', type = 'unsigned' },
})
box.space.votes:create_index('primary', {
    parts = {'poll_id', 'user_id'},
    unique = true,
    if_not_exists = true
})

box.space.votes:create_index('poll_id', {parts = {'poll_id'}, if_not_exists = true}) 

print("Tarantool initialization completed!")