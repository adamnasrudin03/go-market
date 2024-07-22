package repository

// Repositories all repo object injected here
type Repositories struct {
	Auth    AuthRepository
	Cache   CacheRepository
	Product ProductRepository
}
