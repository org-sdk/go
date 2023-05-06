package storage

import (
    "github.com/dgraph-io/badger/v3"
    "github.com/dgraph-io/badger/v3/options"
    "log"
    "time"
)

type BadgerStore struct {
    custom      bool
    snapshotsDB *badger.DB
    cacheDB     *badger.DB
    closing     bool
}

func NewBadgerStore(custom bool, dir string) (*BadgerStore, error) {
    snapshotsDB, err := openDB(dir+"/snapshots", true, custom)
    if err != nil {
        return nil, err
    }
    cacheDB, err := openDB(dir+"/cache", false, custom)
    if err != nil {
        return nil, err
    }
    return &BadgerStore{
        custom:      true,
        snapshotsDB: snapshotsDB,
        cacheDB:     cacheDB,
        closing:     false,
    }, nil
}

func (store *BadgerStore) Close() error {
    store.closing = true
    err := store.snapshotsDB.Close()
    if err != nil {
        return err
    }
    return store.cacheDB.Close()
}

func openDB(dir string, sync bool, custom bool) (*badger.DB, error) {
    opts := badger.DefaultOptions(dir)
    opts = opts.WithSyncWrites(sync)
    opts = opts.WithCompression(options.None)
    opts = opts.WithBlockCacheSize(0)
    opts = opts.WithIndexCacheSize(0)
    opts = opts.WithMetricsEnabled(false)
    opts = opts.WithLoggingLevel(badger.ERROR)
    db, err := badger.Open(opts)
    if err != nil {
        return nil, err
    }

    if custom {
        go func() {
            for {
                lsm, vlog := db.Size()
                log.Printf("Badger LSM %d VLOG %d\n", lsm, vlog)
                if lsm > 1024*1024*8 || vlog > 1024*1024*32 {
                    err := db.RunValueLogGC(0.5)
                    log.Printf("Badger RunValueLogGC %v\n", err)
                }
                time.Sleep(5 * time.Minute)
            }
        }()
    }

    return db, nil
}
