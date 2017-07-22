package passwords

import (
	"net/http"

	"github.com/keratin/authn-server/api"
	"github.com/keratin/authn-server/services"
)

func getPasswordReset(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := app.AccountStore.FindByUsername(r.FormValue("username"))
		if err != nil {
			panic(err)
		}

		// run in the background so that a timing attack can't enumerate usernames
		go func() {
			err := services.PasswordResetSender(app.Config, account)
			if err != nil {
				// TODO: report and continue
			}
		}()

		w.WriteHeader(http.StatusOK)
	}
}
