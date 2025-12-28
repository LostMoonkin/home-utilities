package infrastructure

import (
	"context"
	"fmt"
	"homeserver/common"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/semaphore"
)

const DEFAULT_LOCK_TIMEOUT = time.Duration(5) * time.Second
const DEFAULT_IDLE_TIMEOUT = time.Duration(10) * time.Minute

type SSHClientWrapper struct {
	*ssh.Client
	SFTPClient     *sftp.Client
	Key            string
	lastUsed       time.Time
	MaxIdleTimeout time.Duration
	lock           *semaphore.Weighted
}

var clientMap sync.Map
var lockMap sync.Map

func init() {
	clientMap = sync.Map{}
	lockMap = sync.Map{}
}

func (s *SSHClientWrapper) Close() {
	s.lock.Acquire(context.Background(), 1)
	defer s.lock.Release(1)
	s.closeInternal()
}

func (s *SSHClientWrapper) closeInternal() {
	if s.Client != nil {
		s.Client.Close()
		s.Client = nil
	}
	if s.SFTPClient != nil {
		s.SFTPClient.Close()
		s.SFTPClient = nil
	}
	clientMap.Delete(s.Key)
}

func (s *SSHClientWrapper) keepAlive() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// if locked, skip try to keep alive.
		if s.lock.TryAcquire(1) {
			if _, err := s.SFTPClient.Getwd(); err != nil {
				common.Log.Warn().Err(err).Str("key", s.Key).Msg("Keepalive failed")
				s.closeInternal()
				s.lock.Release(1)
				return
			}
			if time.Since(s.lastUsed) > s.MaxIdleTimeout {
				common.Log.Info().Str("key", s.Key).Msg("client idle timeout.")
				s.closeInternal()
				s.lock.Release(1)
				return
			}
			s.lock.Release(1)
		}
	}
}

func GetSSHClient(ctx context.Context, user, address string, privateKey []byte) (*SSHClientWrapper, error) {
	key := fmt.Sprintf("%s@%s", user, address)

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
	// get client if exists
	val, ok := clientMap.Load(key)
	if ok {
		// check client alive
		client := val.(*SSHClientWrapper)
		_, err := client.SFTPClient.Getwd()
		if err == nil {
			return client, nil
		}
		common.Log.Warn().Str("key", key).Err(err).Msg("Connection is dead.")
		client.Close()
	}

	// create new SSH client
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
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		common.Log.Error().Err(err).Str("key", key).Msg("Dial ssh client error.")
		return nil, err
	}

	// create new sftp client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		_ = client.Close() // Cleanup SSH if SFTP fails
		return nil, fmt.Errorf("sftp creation error: %w", err)
	}

	// save client
	warpper := &SSHClientWrapper{
		Client:         client,
		SFTPClient:     sftpClient,
		Key:            key,
		lastUsed:       time.Now(),
		MaxIdleTimeout: DEFAULT_IDLE_TIMEOUT,
		lock:           lock,
	}
	clientMap.Store(key, warpper)
	go warpper.keepAlive()
	return warpper, nil
}
