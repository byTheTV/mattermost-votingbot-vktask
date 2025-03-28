box.cfg {
    listen = 3301,
    log_level = 5,
    wal_dir = '/var/lib/tarantool',
    memtx_dir = '/var/lib/tarantool',
    vinyl_dir = '/var/lib/tarantool'
}

-- Создание пользователя для подключения
box.once('init', function()
    -- Создание пользователя с правами доступа
    box.schema.user.create('user', {password = 'password'}, {if_not_exists = true})
    box.schema.user.grant('user', 'read,write,execute', 'universe')
end)

print('Tarantool инициализирован и готов к работе')

