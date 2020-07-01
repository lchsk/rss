const m = require("mithril");

const Login = require("../actions/login");
const Config = require("../config");
const getLink = require("./link");

const LoginComponent = {
  oninit: () => {
    m.request({
      method: "GET",
      url: Config.api_url + "/users/current",
      withCredentials: true,
      responseType: "json",
      extract: function(xhr, options) {
        if (xhr.status === 401) {
          LoginComponent.signedIn = false;
          return false;
        } else if (xhr.status === 200) {
          LoginComponent.signedIn = true;
          return true;
        }
        LoginComponent.signedIn = false;
        return false;
      }
    }).then(function(signedIn) {});
  },
  submitForm: e => {
    e.preventDefault();

    const emailInput = document.getElementById("form-email");
    const passwordInput = document.getElementById("form-password");

    Login.current.email = emailInput.value;
    Login.current.password = passwordInput.value;

    Login.submit();
  },
  view: node => {
    if (LoginComponent.signedIn === true) {
      m.route.set("/");
      return;
    }

    if (LoginComponent.signedIn === false) {
      return (
        <form class="form-major" onsubmit={LoginComponent.submitForm}>
          <div class="text-center mb-3">
            <img src="../assets/text2011.png" alt="rss" />
          </div>
          <div class="login-error">{Login.current.error}</div>
          <input
            id="form-email"
            type="email"
            placeholder="Email"
            class="form-control together-top"
          />
          <input
            id="form-password"
            type="password"
            placeholder="Password"
            class="form-control together-bottom"
          />
          <div class="text-center">
            <button type="submit" class="btn btn-primary">
              Sign in
            </button>
          </div>
          <div className="text-center mt-3">
            {getLink(
              ".mb-0",
              "/signup",
              "Sign up if you don't have an account"
            )}
          </div>
        </form>
      );
    }
  }
};

module.exports = LoginComponent;
