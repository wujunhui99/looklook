Update(ctx context.Context, session sqlx.Session,{{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error
Trans(ctx context.Context,fn func(context context.Context,session sqlx.Session) error) error
