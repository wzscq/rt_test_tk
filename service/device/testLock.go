package device

type TestLock interface {
	GetLock()(bool)
	ReleaseLock()
}