package randomidgenerator

type RandomIdGenerator interface {
	Create(prefix string, length int) string
}
