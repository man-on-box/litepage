package data

type Data struct {
	Domain string
}

func New(domain string) *Data {
	return &Data{
		Domain: domain,
	}
}
