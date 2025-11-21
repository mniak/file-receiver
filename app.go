package main

type AppParams struct {
	HttpParams
}

func NewApp(params AppParams) Service {
	return &MultiService{
		&HttpService{
			HttpParams: params.HttpParams,
		},
	}
}
