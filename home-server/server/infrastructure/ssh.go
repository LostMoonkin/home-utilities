package infrastructure

import (
	"context"
	"fmt"
	"homeserver/common"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/semaphore"
)

const DEFAULT_LOCK_TIMEOUT = time.Duration(5) * time.Second
const DEFAULT_IDLE_TIMEOUT = time.Duration(10) * time.Minute

type SSHClientWrapper struct {
	*ssh.Client
	Key            string
	lastUsed       time.Time
	MaxIdleTimeout time.Duration
}

var clientMap sync.Map
var lockMap sync.Map

func init() {
	clientMap = sync.Map{}
	lockMap = sync.Map{}
}

func (s *SSHClientWrapper) keepAlive() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rawLock, _ := lockMap.LoadOrStore(s.Key, semaphore.NewWeighted(1))
		lock := rawLock.(*semaphore.Weighted)
		// if locked, skip try keep alive.
		if lock.TryAcquire(1) {
			defer lock.Release(1)
			if _, _, err := s.SendRequest("keepalive@golang.org", true, nil); err != nil {
				common.Log.Warn().Err(err).Str("key", s.Key).Msg("Keepalive failed")
				_ = s.Close()
				clientMap.Delete(s.Key)
				return
			}
			if time.Since(s.lastUsed) > s.MaxIdleTimeout {
				common.Log.Info().Str("key", s.Key).Msg("client idle timeout.")
				_ = s.Close()
				clientMap.Delete(s.Key)
				return
			}
		}
	}
}

func GetClient(ctx context.Context, user, host string, port int, privateKey []byte) (*SSHClientWrapper, error) {
	key := fmt.Sprintf("%s@%s:%d", user, host, port)
	rawLock, _ := lockMap.LoadOrStore(key, semaphore.NewWeighted(1))
	lock := rawLock.(*semaphore.Weighted)

	// acquire lock with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, DEFAULT_LOCK_TIMEOUT)
	defer cancel()
	if err := lock.Acquire(timeoutCtx, 1); err != nil {
		// timeout
		common.Log.Warn().Str("key", key).Err(err).Msg("acquire lock error.")
		return nil, err
	}
	defer lock.Release(1)

	val, ok := clientMap.Load(key)
	if ok {
		// check client alive
		client := val.(*SSHClientWrapper)
		_, _, err := client.SendRequest("keepalive@golang.org", true, nil)
		if err == nil {
			return client, nil
		}
		common.Log.Warn().Str("key", key).Err(err).Msg("Connection is dead.")
		_ = client.Close()
		clientMap.Delete(key)
	}
	// create new client
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		common.Log.Error().Err(err).Str("key", key).Msg("Parse privete key error.")
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		common.Log.Error().Err(err).Str("key", key).Msg("Dial ssh client error.")
		return nil, err
	}

	// save client
	warpper := &SSHClientWrapper{
		Key:            key,
		Client:         client,
		lastUsed:       time.Now(),
		MaxIdleTimeout: DEFAULT_IDLE_TIMEOUT,
	}
	clientMap.Store(key, warpper)
	go warpper.keepAlive()
	return warpper, nil
}
