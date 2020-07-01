const m = require("mithril");

const Config = require("../config");

const Account = {
  submitPasswordChange: e => {
    e.preventDefault();

    const currentPassword = document.getElementById("current-password").value;
    const newPassword1 = document.getElementById("new-password-1").value;
    const newPassword2 = document.getElementById("new-password-2").value;

    if (newPassword1 !== newPassword2) {
      alert("Passwords are not the same");
    }

    m.request({
      method: "PATCH",
      url: Config.api_url + "/users/current/password",
      data: {
        current: currentPassword,
        new: newPassword1
      },
      withCredentials: true
    })
      .then(result => {
        m.route.set("/logout");
      })
      .catch(e => {
        alert("Error when trying to change the password");
      });
  },
  view: node => {
    return (
      <div class="container">
        <h1>Change your password</h1>
        <form onsubmit={Account.submitPasswordChange}>
          <div className="form-group">
            <input
              type="password"
              placeholder="Current password"
              className="form-control"
              id="current-password"
            />
          </div>
          <div className="form-group">
            <input
              type="password"
              placeholder="New password"
              className="form-control"
              id="new-password-1"
            />
          </div>
          <div className="form-group">
            <input
              type="password"
              placeholder="Repeat your new password"
              className="form-control"
              id="new-password-2"
            />
          </div>
          <button type="submit" className="btn btn-primary">
            Save
          </button>
        </form>
      </div>
    );
  }
};

module.exports = Account;
