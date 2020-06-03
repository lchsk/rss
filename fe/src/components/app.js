var m = require("mithril");

const NavBar = require("./nav_bar");
const MainView = require("./main_view");
const User = require("../actions/user");
const Config = require("../config");
const LoginComponent = require("./login");

const App = {
  oninit: node => {
    m.request({
      method: "GET",
      url: Config.api_url + "/users/current",
      withCredentials: true,
      responseType: "json",
      extract: function(xhr, options) {
        if (xhr.status === 401) {
          LoginComponent.signedIn = false;
          return false;
        }

        LoginComponent.signedIn = true;
        return true;
      }
    }).then(function(signedIn) {
      console.log("signed-in", signedIn);
    });
  },
  view: node => {
    return m(".app", [m(NavBar), m(MainView, node.children)]);
  }
};

module.exports = App;
