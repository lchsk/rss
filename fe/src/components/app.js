const m = require("mithril");

const NavBar = require("./nav_bar");
const MainView = require("./main_view");
const User = require("../actions/user");
const Config = require("../config");
const LoginComponent = require("./login");
const Footer = require("./Footer");

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
        } else if (xhr.status === 200) {
          const response = options.deserialize(xhr.response);
          LoginComponent.username = response['username'];
          console.log("here", options.deserialize(xhr.response));
          LoginComponent.signedIn = true;
          return true;
        }

        LoginComponent.signedIn = false;
        return false;
      }
    }).then(function(signedIn) {
    });
  },
  view: node => {
    return (
        <div class="app">
          <NavBar>

          </NavBar>
          <MainView>
            {node.children}
          </MainView>
          <Footer></Footer>
        </div>
    );
  }
};

module.exports = App;
