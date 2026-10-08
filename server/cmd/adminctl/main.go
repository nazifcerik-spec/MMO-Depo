// Command adminctl bootstraps admin accounts (DD-213). The first admin can
// only be created here; later ones can be created the same way and then
// managed through /internal/admin/users.
//
//	go run ./cmd/adminctl -email ops@example.com -name "Ops" -role admin
//
// The password is read from ADMINCTL_PASSWORD (never a flag, so it stays out
// of shell history). The TOTP secret is printed once — add it to an
// authenticator app.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/otter/server/internal/admin"
	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/platform/db"
)

func main() {
	email := flag.String("email", "", "admin email (required)")
	name := flag.String("name", "", "display name")
	roles := flag.String("role", "viewer", "comma-separated roles: viewer,support_l1,support_l2,game_designer,ops_engineer,admin")
	flag.Parse()

	password := os.Getenv("ADMINCTL_PASSWORD")
	if *email == "" || len(password) < 12 {
		fmt.Fprintln(os.Stderr, "usage: ADMINCTL_PASSWORD=<12+ chars> adminctl -email <email> [-name <name>] [-role <roles>]")
		os.Exit(2)
	}
	cfg, err := config.Load()
	if err != nil {
		fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := db.Connect(ctx, cfg.DatabaseURL, 2)
	if err != nil {
		fail(err)
	}
	defer pool.Close()

	svc := admin.NewService(pool, nil, nil)
	id, secret, err := svc.CreateAdmin(ctx, *email, *name, password, strings.Split(*roles, ","))
	if err != nil {
		fail(err)
	}
	fmt.Printf("admin created: %s\nTOTP secret (store it now, shown once): %s\notpauth URL: %s\n",
		id, secret, admin.OTPAuthURL("New World MMO Admin", *email, secret))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "adminctl:", err)
	os.Exit(1)
}
