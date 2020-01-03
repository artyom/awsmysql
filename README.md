Package awsmysql provides shortcut functions to retrieve MySQL credentials
from AWS Secrets Manager profile.

Usage example:

    cfg, err := awsmysql.Config(ctx, "production/dbhost")
    if err != nil {
        return err
    }
    // adjust config setting timeouts, database name, etc:
    // see https://godoc.org/github.com/go-sql-driver/mysql#Config
    cfg.DBName = "data"
    connector, err := mysql.NewConnector(cfg) // package github.com/go-sql-driver/mysql
    if err != nil {
        return err
    }
    db := sql.OpenDB(connector)

More details at: <https://pkg.go.dev/github.com/artyom/awsmysql>